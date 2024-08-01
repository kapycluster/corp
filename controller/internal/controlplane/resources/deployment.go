package resources

import (
	"context"

	"github.com/kapycluster/corpy/controller/internal/scope"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Deployment struct {
	Client client.Client
	types.NamespacedName
	scope *scope.ControlPlaneScope
}

func NewDeployment(client client.Client, scope *scope.ControlPlaneScope) *Deployment {
	return &Deployment{
		Client:         client,
		NamespacedName: types.NamespacedName{Name: "kapy-server", Namespace: scope.Namespace()},
		scope:          scope,
	}
}

func (d *Deployment) Create(ctx context.Context) error {
	deploy := d.deployment()

	if err := d.Client.Create(ctx, deploy); client.IgnoreAlreadyExists(err) != nil {
		return err
	}

	if err := d.scope.SetControllerReference(ctx, deploy); err != nil {
		return err
	}

	return nil
}

func (d *Deployment) deployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: d.Namespace,
			Labels:    d.scope.ServerCommonLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: d.scope.ServerCommonLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: d.scope.ServerCommonLabels(),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "kapy-server",
							Image:           d.scope.ServerImage(),
							ImagePullPolicy: corev1.PullIfNotPresent,
							// TODO(icy): populate these
							Command: []string{"k3s", "server"},
							Args: []string{
								"--disable-agent",
								"--disable=traefik",
								"--disable=servicelb",
								"--disable=metrics-server",
							},
							// TODO(icy): add live/readiness probes back; they need auth
						},
					},
				},
			},
		},
	}
}
