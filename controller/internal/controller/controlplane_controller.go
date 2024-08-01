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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kapyv1 "github.com/kapycluster/corpy/controller/api/v1"
	"github.com/kapycluster/corpy/controller/internal/controlplane"
	"github.com/kapycluster/corpy/controller/internal/scope"
)

// ControlPlaneReconciler reconciles a ControlPlane object
type ControlPlaneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kapy.sh,resources=controlplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kapy.sh,resources=controlplanes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kapy.sh,resources=controlplanes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ControlPlane object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
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

	// if !controllerutil.ContainsFinalizer(&kcp, kapyv1.ControlPlaneFinalizer) {
	// 	controllerutil.AddFinalizer(&kcp, kapyv1.ControlPlaneFinalizer)

	// 	err := r.Update(ctx, &kcp)
	// 	if err != nil {
	// 		return ctrl.Result{}, err
	// 	}

	// 	l.Info("added finalizer to ControlPlane object")
	// }

	if kcp.Status.Ready {
		l.Info("ControlPlane already reconciled", "name", scope.Name())
		return ctrl.Result{}, nil
	}

	l.Info("creating control plane", "name", scope.Name(), "namespace", scope.Namespace())
	if err := controlplane.Create(ctx, r.Client, scope); err != nil {
		return ctrl.Result{}, fmt.Errorf("creation failed: %w", err)
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
