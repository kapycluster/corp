package magicnode

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"kapycluster.com/corp/controller/controller/magicnode/google"
	"kapycluster.com/corp/log"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MagicNode struct {
	image     string
	namespace string
	gsaEmail  string
	projectID string
}

func New(image string, namespace string, gsaEmail string, projectID string) *MagicNode {
	return &MagicNode{
		image:     image,
		namespace: namespace,
		gsaEmail:  gsaEmail,
		projectID: projectID,
	}
}

func (m *MagicNode) deployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "magicnode",
			Namespace: m.namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To[int32](1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "magicnode",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "magicnode",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "magicnode",
					Containers: []corev1.Container{
						{
							Name:            "magicnode",
							Image:           m.image,
							ImagePullPolicy: corev1.PullAlways,
							Env: []corev1.EnvVar{
								{
									Name:  "CURRENT_KUBECONFIG",
									Value: "/data/kubeconfig",
								},
								{
									Name:  "MAGICNODE_NATS_URL",
									Value: "nats://nats.nats:4222",
								},
								{
									Name: "MAGICNODE_CLUSTER_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "kubeconfig-volume",
									MountPath: "/data/kubeconfig",
									SubPath:   "value",
								},
							},
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: "regcred",
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "kubeconfig-volume",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "kubeconfig",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (m *MagicNode) serviceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "magicnode",
			Namespace: m.namespace,
			Annotations: map[string]string{
				"iam.gke.io/gcp-service-account": m.gsaEmail,
			},
		},
	}
}

func (m *MagicNode) Install(ctx context.Context, c client.Client) error {
	log := log.FromContext(ctx).With("component", "magicnode")

	sa := m.serviceAccount()
	if err := c.Create(ctx, sa); client.IgnoreAlreadyExists(err) != nil {
		return fmt.Errorf("failed to create magicnode serviceaccount: %w", err)
	}
	log.Info("created deployment for magicnode", "namespace", m.namespace)

	deployment := m.deployment()
	if err := c.Create(ctx, deployment); client.IgnoreAlreadyExists(err) != nil {
		return fmt.Errorf("failed to create magicnode deployment: %w", err)
	}
	log.Info("created serviceaccount for magicnode", "namespace", "m.namespace")

	iam, err := google.NewIAM(ctx, m.gsaEmail, m.projectID, "magicnode", m.namespace)
	if err != nil {
		return err
	}

	err := iam.CreateIAMPolicyBinding(ctx, "roles/iam.workloadIdentityUser")
	if err != nil {
		return err
	}

	return nil
}
