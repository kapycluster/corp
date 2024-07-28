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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kapyv1 "github.com/decantor/corpy/controller/api/v1"
	"github.com/decantor/corpy/controller/internal/controlplane/resources"
	"github.com/decantor/corpy/controller/internal/scope"
)

// KapyClusterReconciler reconciles a KapyCluster object
type KapyClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cluster.kapy.sh,resources=kapyclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.kapy.sh,resources=kapyclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.kapy.sh,resources=kapyclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KapyCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *KapyClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	kc := kapyv1.KapyCluster{}
	err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, &kc)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	scope := scope.NewKapyScope(&kc, r.Client)

	if !kc.ObjectMeta.DeletionTimestamp.IsZero() {
		// reconcile delete here
	}

	if !controllerutil.ContainsFinalizer(&kc, kapyv1.KapyClusterFinalizer) {
		controllerutil.AddFinalizer(&kc, kapyv1.KapyClusterFinalizer)

		err := r.Update(ctx, &kc)
		if err != nil {
			return ctrl.Result{}, err
		}

		l.Info("added finalizer to KapyCluster object")
	}

	if kc.Status.Ready {
		l.Info("KapyCluster already reconciled", "name", scope.Name())
	}

	l.Info("creating service for control plane")
	svc := resources.NewService(r.Client, scope)
	if err := svc.Create(ctx); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create service for control plane %s: %w", scope.Name(), err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KapyClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kapyv1.KapyCluster{}).
		Complete(r)
}
