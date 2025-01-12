package resources

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"kapycluster.com/corp/controller/scope"
	"kapycluster.com/corp/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Deployment struct {
	Client client.Client
	k8stypes.NamespacedName
	scope *scope.ControlPlaneScope
}

func NewDeployment(client client.Client, scope *scope.ControlPlaneScope) *Deployment {
	return &Deployment{
		Client:         client,
		NamespacedName: k8stypes.NamespacedName{Name: "kapyserver", Namespace: scope.Namespace()},
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
					ImagePullSecrets: []corev1.LocalObjectReference{{Name: "regcred"}},
					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "kapyserver-storage",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "kapyserver",
							Image:           d.scope.ServerImage(),
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{"/kapyserver"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/data",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  types.KapyServerClusterCIDR,
									Value: "10.11.0.0/16",
								},
								{
									Name:  types.KapyServerDatastore,
									Value: d.scope.Persistence(),
								},
								{
									Name:  types.KapyServerKubeConfigPath,
									Value: "/data/kubeconfig",
								},
								{
									Name:  types.KapyServerDataDir,
									Value: "/data",
								},
								{
									Name:  types.KapyServerLoadBalancerAddress,
									Value: d.scope.LoadBalancerAddress(),
								},
								{
									Name: types.KapyServerAdvertiseIP,
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name:  types.KapyServerToken,
									Value: d.scope.Token(),
								},
								{
									Name:  types.KapyServerGRPCAddress,
									Value: "0.0.0.0:54545",
								},
							},
							// TODO(icy): add live/readiness probes back; they need auth
						},
					},
				},
			},
		},
	}
}
