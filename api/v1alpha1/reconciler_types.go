/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type ManifestsStorageType string

const (
	Git ManifestsStorageType = "git"
)

// ReconcilerSpec defines the desired state of Reconciler
type ReconcilerSpec struct {
	// +required
	HostName string `json:"hostName"`

	// +required
	ReconcilerName string `json:"reconcilerName"`

	// +optional
	Type string `json:"type"`

	// +optional
	ManifestsStorageType ManifestsStorageType `json:"manifestsStorageType"`

	// +required
	ManifestsEndpoint string `json:"manifestsEndpoint"`
}

// ReconcilerStatus defines the observed state of Reconciler
type ReconcilerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Reconciler is the Schema for the reconcilers API
type Reconciler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReconcilerSpec   `json:"spec,omitempty"`
	Status ReconcilerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReconcilerList contains a list of Reconciler
type ReconcilerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Reconciler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Reconciler{}, &ReconcilerList{})
}
