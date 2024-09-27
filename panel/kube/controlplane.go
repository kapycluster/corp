package kube

import (
	kapyv1 "github.com/kapycluster/corpy/controller/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ControlPlaneStatus string

const (
	ControlPlaneStatusCreating    ControlPlaneStatus = "Creating"
	ControlPlaneStatusInitialized ControlPlaneStatus = "Initialized"
	ControlPlaneStatusReady       ControlPlaneStatus = "Ready"

	labelUserID = "controlplanes.kapy.sh/user"
)

type ControlPlane struct {
	// Name is the name of the control plane
	Name string

	// ID is the namespace of the control plane
	ID string

	// UserID is the user who created the control plane
	UserID string

	// Version is the K8s version of the control plane
	Version string

	// Status is the status of the control plane
	Status string
}

// ToKubeObject converts a ControlPlane object to a kapyv1.ControlPlane object
func (cp *ControlPlane) ToKubeObject() *kapyv1.ControlPlane {
	if cp.Version == "" {
		cp.Version = "1.30"
	}

	if cp.Status == "" {
		cp.Status = string(ControlPlaneStatusCreating)
	}

	return &kapyv1.ControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cp.Name,
			Namespace: cp.ID,
			Labels: map[string]string{
				labelUserID: cp.UserID,
			},
		},
		Spec: kapyv1.ControlPlaneSpec{
			Version: cp.Version,
			Server: kapyv1.KapyServer{
				Token: "dummy",
			},
		},
	}
}

func FromKubeObject(kcp *kapyv1.ControlPlane) *ControlPlane {
	cp := &ControlPlane{
		Name:    kcp.Name,
		ID:      kcp.Namespace,
		UserID:  kcp.Labels[labelUserID],
		Version: kcp.Spec.Version,
	}

	switch {
	case kcp.Status.Ready:
		cp.Status = string(ControlPlaneStatusReady)
	case kcp.Status.Initialized:
		cp.Status = string(ControlPlaneStatusInitialized)
	default:
		cp.Status = string(ControlPlaneStatusCreating)
	}

	return cp
}
