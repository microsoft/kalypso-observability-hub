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

type GitRepoSpec struct {
	//+kubebuilder:validation:MinLength=0
	Repo string `json:"repo"`

	//+kubebuilder:validation:MinLength=0
	Branch string `json:"branch"`

	//+kubebuilder:validation:MinLength=0
	Path string `json:"path"`
}

type WorkspaceSpec struct {
	//+kubebuilder:validation:MinLength=0
	Name string `json:"name"`
}

type ApplicationSpec struct {
	//+kubebuilder:validation:MinLength=0
	Name string `json:"name"`

	Workspace WorkspaceSpec `json:"workspace"`
}

type WorkloadSpec struct {
	//+kubebuilder:validation:MinLength=0
	Name string `json:"name"`

	Source GitRepoSpec `json:"source"`

	Application ApplicationSpec `json:"application"`
}

type DeploymentTargetSpec struct {
	//+kubebuilder:validation:MinLength=0
	Name string `json:"name"`

	//+kubebuilder:validation:MinLength=0
	Environment string `json:"environment"`

	Manifests GitRepoSpec `json:"manifests"`
}

type WorkloadVersionSpec struct {
	//+kubebuilder:validation:MinLength=0
	Version string `json:"version"`

	//+kubebuilder:validation:MinLength=0
	Build string `json:"build"`

	//+kubebuilder:validation:MinLength=0
	Commit string `json:"commit"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	BuildTime metav1.Time `json:"buildTime"`
}

// DeploymentDescriptorSpec defines the desired state of DeploymentDescriptor
type DeploymentDescriptorSpec struct {
	Workload WorkloadSpec `json:"workload"`

	DeploymentTarget DeploymentTargetSpec `json:"deploymentTarget"`

	WorkloadVersion WorkloadVersionSpec `json:"workloadVersion"`
}

// DeploymentDescriptorStatus defines the observed state of DeploymentDescriptor
type DeploymentDescriptorStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DeploymentDescriptor is the Schema for the deploymentdescriptors API
type DeploymentDescriptor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploymentDescriptorSpec   `json:"spec,omitempty"`
	Status DeploymentDescriptorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeploymentDescriptorList contains a list of DeploymentDescriptor
type DeploymentDescriptorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeploymentDescriptor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeploymentDescriptor{}, &DeploymentDescriptorList{})
}
