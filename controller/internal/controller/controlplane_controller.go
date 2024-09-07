/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

	kapyv1 "github.com/kapycluster/corpy/controller/api/v1"
	"github.com/kapycluster/corpy/controller/internal/controlplane"
	"github.com/kapycluster/corpy/controller/internal/scope"
	"github.com/kapycluster/corpy/kapyclient"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
	}

	if kcp.Status.Initalized != kcp.Status.Ready {
		l.Info("ControlPlane initialized but not ready, picking up from where we left off")
		res, err := r.reconcileGRPCRequests(ctx, scope)
		if err != nil {
			return res, err
		}
	}

	if kcp.Status.Ready {
		l.Info("ControlPlane already reconciled", "name", scope.Name())
		return ctrl.Result{}, nil
	}

	l.Info("creating control plane", "name", scope.Name(), "namespace", scope.Namespace())
	if err := controlplane.Create(ctx, r.Client, scope); err != nil {
		return ctrl.Result{}, fmt.Errorf("creation failed: %w", err)
	}

	l.Info("checking if kapyserver is up...")
	kapyDeploy := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: "kapyserver", Namespace: scope.Namespace()}, kapyDeploy)
	if errors.IsNotFound(err) {
		l.Info("kapyserver deployment not found, requeuing...")
		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// TODO: make this entire for loop a non-blocking op
	if kapyDeploy.Status.Conditions == nil {
		l.Info("kapyserver deployment conditions are not available yet...")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	available := false
	progressing := false
	for _, condition := range kapyDeploy.Status.Conditions {
		if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
			available = true
		}
		if condition.Type == appsv1.DeploymentProgressing && condition.Status == corev1.ConditionTrue {
			progressing = true
		}
	}

	if !(available && progressing) {
		l.Info("kapyserver deployment is not healthy yet, requeuing...")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	l.Info("kapyserver deployment is healthy")

	// Re-fetch the Deployment to get the latest status
	err = r.Get(ctx, types.NamespacedName{
		Name:      "kapyserver",
		Namespace: scope.Namespace(),
	}, kapyDeploy)
	if err != nil {
		return ctrl.Result{}, err
	}

	kcp.Status.Initalized = true
	err = scope.UpdateStatus(ctx, &kcp)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update status: %w", err)
	}

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

	return ctrl.Result{}, nil
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
	kapyclient, err := kapyclient.NewKapyClient("127.0.0.1:54545")
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
