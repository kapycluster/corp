package kube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kapyv1 "kapycluster.com/corp/controller/api/v1"
)

type ControlPlaneStatus string

const (
	ControlPlaneStatusCreating    ControlPlaneStatus = "Creating"
	ControlPlaneStatusInitialized ControlPlaneStatus = "Initialized"
	ControlPlaneStatusReady       ControlPlaneStatus = "Ready"

	labelUserID = "controlplanes.kapy.sh/user"
	labelRegion = "controlplanes.kapy.sh/region"
)

type Network struct {
	LoadBalancerAddress string
}

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

	// Network is the network configuration of the control plane
	Network Network

	// Region is the region of the control plane
	Region string
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
				labelRegion: cp.Region,
			},
		},
		Spec: kapyv1.ControlPlaneSpec{
			Version: cp.Version,
			Server: kapyv1.KapyServer{
				Token:       "dummy",
				Persistence: "sqlite:///tmp/data/kine.db?_journal=WAL&cache=shared&_busy_timeout=30000&_txlock=immediate",
			},
			Network: kapyv1.Network{
				LoadBalancerAddress: cp.Network.LoadBalancerAddress,
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
