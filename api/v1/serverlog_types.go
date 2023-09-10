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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// ServerLogSpec defines the desired state of ServerLog
type ServerLogSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:MinLength=2
	Dir        string `json:"dir,omitempty"`
	FileFilter string `json:"fileFilter,omitempty"`
	Pattern    string `json:"pattern,omitempty"`
}

// ServerLogStatus defines the observed state of ServerLog
type ServerLogStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase string `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName={sl}
// ServerLog is the Schema for the serverlogs API
type ServerLog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerLogSpec   `json:"spec,omitempty"`
	Status ServerLogStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServerLogList contains a list of ServerLog
type ServerLogList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServerLog `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServerLog{}, &ServerLogList{})
}
