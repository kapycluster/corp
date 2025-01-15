/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ControlPlaneFinalizer = "controlplanes.kapy.sh/finalizer"

type KapyServer struct {
	// +kubebuilder:default="ghcr.io/kapycluster/kapyserver:master"
	Image string `json:"image"`
	// +kubebuilder:default="sqlite"
	Persistence string `json:"persistence"`
	Token       string `json:"token"`
}

type Network struct {
	LoadBalancerAddress string `json:"loadBalancerAddress,omitempty"`
}

type MagicNode struct {
	Enabled   bool   `json:"enabled"`
	GSAEmail  string `json:"gsaEmail"`
	ProjectID string `json:"projectID"`
}

// ControlPlaneSpec defines the desired state of ControlPlane
type ControlPlaneSpec struct {
	Server    KapyServer `json:"server"`
	Network   Network    `json:"network"`
	MagicNode MagicNode  `json:"magicnode"`

	// Version is the version of Kubernetes to deploy
	Version string `json:"version"`
}

// ControlPlaneStatus defines the observed state of ControlPlane
type ControlPlaneStatus struct {
	// Initialized is set when the Deployment is healthy
	Initialized bool `json:"initialized"`
	// Ready is set when the ControlPlane is ready to serve
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=cp
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="Kubernetes version"
// +kubebuilder:printcolumn:name="LB Address",type="string",JSONPath=".spec.LoadBalancerAddress",description="Load balancer address"
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready",description="Control Plane is ready"
// +kubebuilder:printcolumn:name="Initialized",type="boolean",JSONPath=".status.initialized",description="Deployment initialized"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Age"
// ControlPlane is the Schema for the controlplanes API
type ControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ControlPlaneSpec   `json:"spec,omitempty"`
	Status ControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ControlPlaneList contains a list of ControlPlane
type ControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ControlPlane `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ControlPlane{}, &ControlPlaneList{})
}
