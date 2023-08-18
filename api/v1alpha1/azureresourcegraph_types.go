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

// AzureResourceGraphSpec defines the desired state of AzureResourceGraph
type AzureResourceGraphSpec struct {
	//+kubebuilder:validation:MinLength=0
	// +required
	Subscription string `json:"subscription"`

	//+kubebuilder:validation:MinLength=0
	// +required
	Tenant string `json:"tenant"`

	//+kubebuilder:validation:MinLength=0
	// +optional
	ManagedIdentiy string `json:"managedIdentity"`

	// +required
	Interval metav1.Duration `json:"interval"`
}

// AzureResourceGraphStatus defines the observed state of AzureResourceGraph
type AzureResourceGraphStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AzureResourceGraph is the Schema for the azureresourcegraphs API
type AzureResourceGraph struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureResourceGraphSpec   `json:"spec,omitempty"`
	Status AzureResourceGraphStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AzureResourceGraphList contains a list of AzureResourceGraph
type AzureResourceGraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureResourceGraph `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureResourceGraph{}, &AzureResourceGraphList{})
}
