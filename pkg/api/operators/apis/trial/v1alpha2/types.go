/*

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

package v1alpha2

import (
	pb "github.com/kubeflow/katib/pkg/api"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TrialSpec struct {
	// Key-value pairs for Parameters and Assignment values
	SuggestionParameters []pb.SuggestionParameter `json:"suggestionParameters"`

	// Raw text for the trial run template
	RawTemplate string `json:"rawTemplate,omitempty"`
}

type TrialStatus struct {
	StartTime         *metav1.Time     `json:"startTime,omitempty"`
	CompletionTime    *metav1.Time     `json:"completionTime,omitempty"`
	LastReconcileTime *metav1.Time     `json:"lastReconcileTime,omitempty"`
	Conditions        []TrialCondition `json:"conditions,omitempty"`

	// Results - objectives and other metrics values
	Observation Observation `json:"observation,omitempty"`
}

type Observation struct {
	ObjectiveName  string  `json:"objectiveName,omitempty"`
	ObjectiveValue float64 `json:"objectiveValue,omitempty"`
}

// +k8s:deepcopy-gen=true
// TrialCondition describes the state of the trial at a certain point.
type TrialCondition struct {
	// Type of trial condition.
	Type TrialConditionType `json:"type"`

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

// TrialConditionType describes the various conditions a Trial can be in.
type TrialConditionType string

const (
	TrialPending   TrialConditionType = "Pending"
	TrialRunning   TrialConditionType = "Running"
	TrialCompleted TrialConditionType = "Completed"
	TrialKilled    TrialConditionType = "Killed"
	TrialFailed    TrialConditionType = "Failed"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Trial is the Schema for the trial API
// +k8s:openapi-gen=true
type Trial struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TrialSpec   `json:"spec,omitempty"`
	Status TrialStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TrialList contains a list of Trials
type TrialList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Trial `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Trial{}, &TrialList{})
}
