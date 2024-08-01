package scope

import (
	kapyv1 "github.com/kapycluster/corpy/controller/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type KapyScope struct {
	kc     *kapyv1.ControlPlane
	client client.Client
}

func NewKapyScope(kc *kapyv1.ControlPlane, client client.Client) *KapyScope {
	return &KapyScope{kc: kc, client: client}
}

func (k *KapyScope) Object() metav1.Object {
	return k.kc
}

func (k *KapyScope) Name() string {
	return k.kc.Name
}

func (k *KapyScope) Namespace() string {
	return k.kc.Namespace
}

func (k *KapyScope) ServerImage() string {
	return k.kc.Spec.Server.Image
}

func (k *KapyScope) SetControllerReference(child metav1.Object) error {
	return controllerutil.SetControllerReference(k.kc, child, k.client.Scheme())
}

func (k *KapyScope) ServerCommonLabels() map[string]string {
	return map[string]string{
		"cluster.kapy.sh/name":      k.Name(),
		"cluster.kapy.sh/component": "kapy-server",
	}
}
