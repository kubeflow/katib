/*
Copyright 2022 The Kubeflow Authors.

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

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"

	experiment "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
)

// SuggestionSpec is the specification of a Suggestion.
type SuggestionSpec struct {
	// Algorithm describes HP or NAS algorithm that suggestion is used.
	Algorithm *common.AlgorithmSpec `json:"algorithm,omitempty"`

	// EarlyStopping describes early stopping algorithm that suggestion is used.
	EarlyStopping *common.EarlyStoppingSpec `json:"earlyStopping,omitempty"`

	// Number of suggestions requested.
	Requests int32 `json:"requests,omitempty"`

	// ResumePolicy describes resuming policy which usually take effect after experiment terminated.
	// Default value is Never.
	ResumePolicy experiment.ResumePolicyType `json:"resumePolicy,omitempty"`
}

// SuggestionStatus is the current status of a Suggestion.
type SuggestionStatus struct {
	// AlgorithmSettings defines HP or NAS algorithm settings which suggestion gRPC service returns.
	// These settings overwrites Experiment's settings before the gRPC request.
	// It can be empty if settings haven't been changed.
	AlgorithmSettings []common.AlgorithmSetting `json:"algorithmSettings,omitempty"`

	// Number of suggestion results
	SuggestionCount int32 `json:"suggestionCount,omitempty"`

	// Suggestion results
	Suggestions []TrialAssignment `json:"suggestions,omitempty"`

	// Represents time when the Suggestion was acknowledged by the Suggestion controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the Suggestion was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the Suggestion was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`

	// List of observed runtime conditions for this Suggestion.
	Conditions []SuggestionCondition `json:"conditions,omitempty"`
}

// TrialAssignment is the assignment for one trial.
type TrialAssignment struct {
	// Suggestion results with Trial parameters
	ParameterAssignments []common.ParameterAssignment `json:"parameterAssignments,omitempty"`

	// Name of the suggestion
	Name string `json:"name,omitempty"`

	// Rules for early stopping techniques
	// Contains rule name, value and comparison type
	EarlyStoppingRules []common.EarlyStoppingRule `json:"earlyStoppingRules,omitempty"`

	// Suggestion label metadata to attach to Trial job
	Labels map[string]string `json:"labels,omitempty"`
}

// SuggestionCondition describes the state of the Suggestion at a certain point.
// +k8s:deepcopy-gen=true
type SuggestionCondition struct {
	// Type of Suggestion condition.
	Type SuggestionConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`

	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`

	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`

	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// SuggestionConditionType defines the state of a Suggestion.
type SuggestionConditionType string

const (
	SuggestionCreated         SuggestionConditionType = "Created"
	SuggestionDeploymentReady SuggestionConditionType = "DeploymentReady"
	SuggestionRunning         SuggestionConditionType = "Running"
	SuggestionSucceeded       SuggestionConditionType = "Succeeded"
	SuggestionFailed          SuggestionConditionType = "Failed"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Suggestion represents the structure of a Suggestion resource.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Suggestion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SuggestionSpec   `json:"spec,omitempty"`
	Status SuggestionStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SuggestionList contains a list of Suggestion
type SuggestionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Suggestion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Suggestion{}, &SuggestionList{})
}
