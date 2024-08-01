package resources

import (
	"context"

	"github.com/kapycluster/corpy/controller/internal/scope"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Deployment struct {
	Client client.Client
	types.NamespacedName
	scope *scope.KapyScope
}

func NewDeployment(client client.Client, scope *scope.KapyScope) *Deployment {
	return &Deployment{
		Client:         client,
		NamespacedName: types.NamespacedName{Name: "kapy-server", Namespace: scope.Namespace()},
		scope:          scope,
	}
}

func (d *Deployment) Create(ctx context.Context) error {
	deploy := d.deployment()

	if err := d.Client.Create(ctx, deploy); err != nil {
		return err
	}

	if err := d.scope.SetControllerReference(deploy); err != nil {
		return err
	}

	return nil
}

func (d *Deployment) deployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   d.Name,
			Labels: d.scope.ServerCommonLabels(),
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
							Args:    []string{"--disable-agent"},
							LivenessProbe: &corev1.Probe{
								FailureThreshold:    3,
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								TimeoutSeconds:      60,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Host:   "127.0.0.1",
										Path:   "/healthz",
										Port:   intstr.FromInt(6443),
										Scheme: corev1.URISchemeHTTP,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
