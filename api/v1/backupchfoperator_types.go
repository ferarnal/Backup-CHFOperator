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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BackupCHFOperatorSpec defines the desired state of BackupCHFOperator
type BackupCHFOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Command  string   `json:"command,omitempty"`
	Args     []string `json:"args,omitempty"`
	Interval int      `json:"interval,omitempty"`
}

// BackupCHFOperatorStatus defines the observed state of BackupCHFOperator
type BackupCHFOperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LastExecutionTime metav1.Time `json:"lastExecutionTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// BackupCHFOperator is the Schema for the backupchfoperators API
type BackupCHFOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupCHFOperatorSpec   `json:"spec,omitempty"`
	Status BackupCHFOperatorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// BackupCHFOperatorList contains a list of BackupCHFOperator
type BackupCHFOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupCHFOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupCHFOperator{}, &BackupCHFOperatorList{})
}
