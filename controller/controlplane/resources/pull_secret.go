package resources

import (
	"context"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"kapycluster.com/corp/controller/scope"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PullSecret struct {
	Client client.Client
	types.NamespacedName
	scope *scope.ControlPlaneScope
}

func NewPullSecret(client client.Client, scope *scope.ControlPlaneScope) *PullSecret {
	return &PullSecret{
		Client:         client,
		NamespacedName: types.NamespacedName{Name: "regcred", Namespace: scope.Namespace()},
		scope:          scope,
	}
}

// Create copies the regcred secret from the current namespace into the control plane namespace
func (s *PullSecret) Create(ctx context.Context) error {
	sourceSecret := &corev1.Secret{}
	if err := s.Client.Get(ctx, types.NamespacedName{Name: "regcred", Namespace: os.Getenv("POD_NAMESPACE")}, sourceSecret); err != nil {
		return err
	}

	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
			Labels:    s.scope.ServerCommonLabels(),
		},
		Type: sourceSecret.Type,
		Data: sourceSecret.Data,
	}

	if err := s.Client.Create(ctx, newSecret); client.IgnoreAlreadyExists(err) != nil {
		return err
	}

	if err := s.scope.SetControllerReference(ctx, newSecret); err != nil {
		return err
	}

	return nil
}
