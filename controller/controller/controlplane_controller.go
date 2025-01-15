package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kapyv1 "kapycluster.com/corp/controller/api/v1"
	"kapycluster.com/corp/controller/controlplane"
	"kapycluster.com/corp/controller/scope"
	"kapycluster.com/corp/kapyclient"
)

// ControlPlaneReconciler reconciles a ControlPlane object
type ControlPlaneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kapy.sh,resources=controlplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kapy.sh,resources=controlplanes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kapy.sh,resources=controlplanes/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;delete

func (r *ControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	kcp := kapyv1.ControlPlane{}
	err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, &kcp)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	scope := scope.NewControlPlaneScope(&kcp, r.Client)

	if !kcp.ObjectMeta.DeletionTimestamp.IsZero() {
		// reconcile delete here
		return ctrl.Result{}, nil
	}

	// If already fully reconciled, nothing to do
	if kcp.Status.Ready {
		l.Info("ControlPlane already reconciled", "name", scope.Name())
		return ctrl.Result{}, nil
	}

	// Phase 1: Create resources if not created yet
	if !kcp.Status.Initialized {
		l.Info("creating control plane", "name", scope.Name(), "namespace", scope.Namespace())
		if err := controlplane.Create(ctx, r.Client, scope); err != nil {
			return ctrl.Result{}, fmt.Errorf("creation failed: %w", err)
		}

		kcp.Status.Initialized = true
		if err := scope.UpdateStatus(ctx, &kcp); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to update status: %w", err)
		}

		// Requeue to check deployment status
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	// Phase 2: Wait for kapyserver deployment to be healthy
	kapyDeploy := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: "kapyserver", Namespace: scope.Namespace()}, kapyDeploy)
	if errors.IsNotFound(err) {
		l.Info("kapyserver deployment not found, requeuing...")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	if !isDeploymentHealthy(kapyDeploy) {
		l.Info("kapyserver deployment is not healthy yet, requeuing...")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	l.Info("kapyserver deployment is healthy")

	// Phase 3: Handle gRPC requests and finish reconciliation
	if !kcp.Status.Ready {
		l.Info("letting kapyserver breathe a bit")
		time.Sleep(5 * time.Second)

		res, err := r.reconcileGRPCRequests(ctx, scope)
		if err != nil {
			return res, err
		}

		kcp.Status.Ready = true
		if err := scope.UpdateStatus(ctx, &kcp); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to update status: %w", err)
		}
	}

	return ctrl.Result{}, nil
}

// Helper function to check deployment health
func isDeploymentHealthy(deploy *appsv1.Deployment) bool {
	if deploy.Status.Conditions == nil {
		return false
	}

	available := false
	progressing := false
	for _, condition := range deploy.Status.Conditions {
		if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
			available = true
		}
		if condition.Type == appsv1.DeploymentProgressing && condition.Status == corev1.ConditionTrue {
			progressing = true
		}
	}

	return available && progressing
}

// SetupWithManager sets up the controller with the Manager.
func (r *ControlPlaneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kapyv1.ControlPlane{}).
		Complete(r)
}

func (r *ControlPlaneReconciler) reconcileGRPCRequests(ctx context.Context, scope *scope.ControlPlaneScope) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("asking kapyserver for kubeconfig")
	kapyclient, err := kapyclient.NewKapyClient(scope.ServiceAddress())
	if err != nil {
		l.Info("failed to connect to kapyserver, requeing...", "error", err.Error())
		return ctrl.Result{RequeueAfter: time.Second * 5}, nil
	}
	defer kapyclient.Close()

	kcfg, err := kapyclient.GetKubeConfig(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to fetch kubeconfig: %w", err)
	}

	kcfgData := make(map[string][]byte)
	kcfgData["value"] = []byte(kcfg)

	kubeconfigSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubeconfig",
			Namespace: scope.Namespace(),
		},
		Data: kcfgData,
	}

	l.Info("creating kubeconfig secret", "name", kubeconfigSecret.Name)
	if err := r.Client.Create(ctx, kubeconfigSecret); client.IgnoreAlreadyExists(err) != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create kubeconfig secret: %w", err)
	}
	if err := scope.SetControllerReference(ctx, kubeconfigSecret); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set controller reference: %w", err)
	}

	token, err := kapyclient.GenerateToken(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenData := make(map[string][]byte)
	tokenData["value"] = []byte(token)

	tokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    scope.Namespace(),
			GenerateName: "join-token",
		},
		Data: tokenData,
	}

	l.Info("creating token secret", "name", tokenSecret.Name)
	if err := r.Client.Create(ctx, tokenSecret); client.IgnoreAlreadyExists(err) != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create token secret: %w", err)
	}
	if err := scope.SetControllerReference(ctx, tokenSecret); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set controller reference: %w", err)
	}

	return ctrl.Result{}, nil
}
