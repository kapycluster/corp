package resources

import (
	"context"

	"github.com/kapycluster/corpy/controller/internal/scope"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Service struct {
	Client client.Client
	types.NamespacedName
	scope *scope.KapyScope
}

func NewService(client client.Client, scope *scope.KapyScope) *Service {
	return &Service{
		Client:         client,
		NamespacedName: types.NamespacedName{Name: "kapy-server", Namespace: scope.Namespace()},
		scope:          scope,
	}
}

func (s *Service) Create(ctx context.Context) error {
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.scope.Name(),
			Namespace: s.scope.Namespace(),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:       "kapy-server",
					Port:       6443,
					TargetPort: intstr.FromInt(6443),
				},
			},
			Selector: s.scope.ServerCommonLabels(),
		},
	}

	if err := s.Client.Create(ctx, &svc); err != client.IgnoreAlreadyExists(err) {
		return err
	}

	if err := s.scope.SetControllerReference(&svc); err != nil {
		return nil
	}

	return nil
}
