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
	trial "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ExperimentSpec struct {
	// List of hyperparameter configurations.
	Parameters []ParameterSpec `json:"parameters,omitempty"`

	// Describes the objective of the experiment.
	Objective *ObjectiveSpec `json:"objective,omitempty"`

	// Describes the suggestion algorithm.
	Algorithm *AlgorithmSpec `json:"algorithm,omitempty"`

	// Template for each run of the trial.
	TrialTemplate *TrialTemplate `json:"trialTemplate,omitempty"`

	// How many trials can be processed in parallel.
	ParallelTrialCount int `json:"parallelTrialCount,omitempty"`

	// Total number of trials to run.
	MaxTrialCount      int `json:"maxTrialCount,omitempty"`

	// TODO - figure out what to do with metric collectors
	MetricsCollectorType string `json:"metricsCollectorSpec,omitempty"`

	// TODO - Other fields, exact format is TBD. Will add these back during implementation.
	// - NAS
	// - Early stopping
	// - Resume experiment
}

type ExperimentStatus struct {
	// Represents time when the Experiment was acknowledged by the Experiment controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime         *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the Experiment was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime    *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the Experiment was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`

	// List of observed runtime conditions for this Experiment.
	Conditions []ExperimentCondition `json:"conditions,omitempty"`

	// Current optimal trial parameters and observations.
	CurrentOptimalTrial OptimalTrial `json:"currentOptimalTrial,omitempty"`

	// How many trials have successfully completed.
	TrialsCompleted int `json:"trialsCompleted,omitempty"`

	// How many trials have failed.
	TrialsFailed int `json:"trialsFailed,omitempty"`

	// How many trials have been killed.
	TrialsKilled int `json:"trialsKilled,omitempty"`

	// How many trials are currently pending.
	TrialsPending int `json:"trialsPending,omitempty"`
}

type OptimalTrial struct {
	// Key-value pairs for hyperparameters and assignment values.
	ParameterAssignments []trial.ParameterAssignment `json:"parameterAssignments"`

	// Observation for this trial
	Observation trial.Observation `json:"observation,omitempty"`
}

// +k8s:deepcopy-gen=true
// ExperimentCondition describes the state of the experiment at a certain point.
type ExperimentCondition struct {
	// Type of experiment condition.
	Type ExperimentConditionType `json:"type"`

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

// ExperimentConditionType defines the state of an Experiment.
type ExperimentConditionType string

const (
	ExperimentCreated    ExperimentConditionType = "Created"
	ExperimentRunning    ExperimentConditionType = "Running"
	ExperimentRestarting ExperimentConditionType = "Restarting"
	ExperimentSucceeded  ExperimentConditionType = "Succeeded"
	ExperimentFailed     ExperimentConditionType = "Failed"
)

type ParameterSpec struct {
	Name          string        `json:"name,omitempty"`
	ParameterType ParameterType `json:"parameterType,omitempty"`
	FeasibleSpace FeasibleSpace `json:"feasibleSpace,omitempty"`
}

type ParameterType string

const (
	ParameterTypeUnknown     ParameterType = "unknown"
	ParameterTypeDouble      ParameterType = "double"
	ParameterTypeInt         ParameterType = "int"
	ParameterTypeDiscrete    ParameterType = "discrete"
	ParameterTypeCategorical ParameterType = "categorical"
)

type FeasibleSpace struct {
	Max  string   `json:"max,omitempty"`
	Min  string   `json:"min,omitempty"`
	List []string `json:"list,omitempty"`
	Step string   `json:"step,omitempty"`
}

type ObjectiveSpec struct {
	Type                ObjectiveType `json:"type,omitempty"`
	Goal                float64       `json:"goal,omitempty"`
	ObjectiveMetricName string        `json:"objectiveMetricName,omitempty"`
	// This can be empty if we only care about the objective metric.
	// Note: If we adopt a push instead of pull mechanism, this can be omitted completely.
	AdditionalMetricsNames []string `json:"additionalMetricsNames,omitempty"`
}

type ObjectiveType string

const (
	ObjectiveTypeUnknown  ObjectiveType = ""
	ObjectiveTypeMinimize ObjectiveType = "minimize"
	ObjectiveTypeMaximize ObjectiveType = "maximize"
)

type AlgorithmSpec struct {
	AlgorithmName string		 `json:"algorithmName,omitempty"`
	// Key-value pairs for hyperparameters and assignment values.
	ParameterAssignments []trial.ParameterAssignment `json:"parameterAssignments"`
	EarlyStopping *EarlyStoppingSpec `json:"earlyStopping,omitempty"`
}

type EarlyStoppingSpec struct {
	// TODO
}

type TrialTemplate struct {
	Retain     bool        `json:"retain,omitempty"`
	GoTemplate *GoTemplate `json:"goTemplate,omitempty"`
}

type TemplateSpec struct {
	ConfigMapName      string `json:"configMapName,omitempty"`
	ConfigMapNamespace string `json:"configMapNamespace,omitempty"`
	TemplatePath       string `json:"templatePath,omitempty"`
}

type GoTemplate struct {
	TemplateSpec *TemplateSpec `json:"templateSpec,omitempty"`
	RawTemplate  string        `json:"rawTemplate,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Structure of the Experiment custom resource.
// +k8s:openapi-gen=true
type Experiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExperimentSpec   `json:"spec,omitempty"`
	Status ExperimentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExperimentList contains a list of Experiments
type ExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Experiment `json:"items"`
}

// TODO - enable this during API implementation.
//func init() {
//	SchemeBuilder.Register(&Experiment{}, &ExperimentList{})
//}
