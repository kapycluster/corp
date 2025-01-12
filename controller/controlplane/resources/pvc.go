package resources

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"kapycluster.com/corp/controller/scope"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const pvcName = "kapyserver-storage"

type PersistentVolumeClaim struct {
	Client client.Client
	k8stypes.NamespacedName
	scope *scope.ControlPlaneScope
}

func NewPersistentVolumeClaim(client client.Client, scope *scope.ControlPlaneScope) *PersistentVolumeClaim {
	return &PersistentVolumeClaim{
		Client:         client,
		NamespacedName: k8stypes.NamespacedName{Name: pvcName, Namespace: scope.Namespace()},
		scope:          scope,
	}
}

func (p *PersistentVolumeClaim) Create(ctx context.Context) error {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},

			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	err := p.Client.Create(ctx, pvc)
	if err != nil {
		return fmt.Errorf("failed to create pvc: %w", err)
	}

	return nil
}
