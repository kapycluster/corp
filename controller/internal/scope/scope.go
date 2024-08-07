package scope

import (
	"context"
	"fmt"

	kapyv1 "github.com/kapycluster/corpy/controller/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ControlPlaneScope struct {
	kc     *kapyv1.ControlPlane
	client client.Client
}

func NewControlPlaneScope(kc *kapyv1.ControlPlane, client client.Client) *ControlPlaneScope {
	return &ControlPlaneScope{kc: kc, client: client}
}

func (k *ControlPlaneScope) Object() metav1.Object {
	return k.kc
}

func (k *ControlPlaneScope) Name() string {
	return k.kc.Name
}

func (k *ControlPlaneScope) Namespace() string {
	return k.kc.Namespace
}

func (k *ControlPlaneScope) ServerImage() string {
	return k.kc.Spec.Server.Image
}

func (k *ControlPlaneScope) Persistence() string {
	return k.kc.Spec.Server.Persistence
}

func (k *ControlPlaneScope) Token() string {
	return k.kc.Spec.Server.Token
}

func (k *ControlPlaneScope) LoadBalancerAddress() string {
	return k.kc.Spec.Network.LoadBalancerAddress
}

func (k *ControlPlaneScope) SetControllerReference(ctx context.Context, child metav1.Object) error {
	if err := controllerutil.SetControllerReference(k.kc, child, k.client.Scheme()); err != nil {
		return err
	}

	obj, ok := child.(client.Object)
	if !ok {
		return fmt.Errorf("child is not a client.Object")
	}

	return k.client.Update(ctx, obj)
}

func (k *ControlPlaneScope) UpdateStatus(ctx context.Context, kc *kapyv1.ControlPlane) error {
	return k.client.Status().Update(ctx, kc)
}

func (k *ControlPlaneScope) ServerCommonLabels() map[string]string {
	return map[string]string{
		"controlplane.kapy.sh/name":      k.Name(),
		"controlplane.kapy.sh/component": "kapy-server",
	}
}
