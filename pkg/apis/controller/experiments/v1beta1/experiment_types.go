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
	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ExperimentSpec is the specification of an Experiment.
type ExperimentSpec struct {
	// List of hyperparameter configurations.
	Parameters []ParameterSpec `json:"parameters,omitempty"`

	// Describes the objective of the experiment.
	Objective *common.ObjectiveSpec `json:"objective,omitempty"`

	// Describes the suggestion algorithm.
	Algorithm *common.AlgorithmSpec `json:"algorithm,omitempty"`

	// Describes the early stopping algorithm.
	EarlyStopping *common.EarlyStoppingSpec `json:"earlyStopping,omitempty"`

	// Template for each run of the trial.
	TrialTemplate *TrialTemplate `json:"trialTemplate,omitempty"`

	// How many trials can be processed in parallel.
	// Defaults to 3
	ParallelTrialCount *int32 `json:"parallelTrialCount,omitempty"`

	// Max completed trials to mark experiment as succeeded
	MaxTrialCount *int32 `json:"maxTrialCount,omitempty"`

	// Max failed trials to mark experiment as failed.
	MaxFailedTrialCount *int32 `json:"maxFailedTrialCount,omitempty"`

	// Describes the specification of the metrics collector
	MetricsCollectorSpec *common.MetricsCollectorSpec `json:"metricsCollectorSpec,omitempty"`

	NasConfig *NasConfig `json:"nasConfig,omitempty"`

	// Describes resuming policy which usually take effect after experiment terminated.
	// Default value is Never.
	ResumePolicy ResumePolicyType `json:"resumePolicy,omitempty"`
}

// ExperimentStatus is the current status of an Experiment.
type ExperimentStatus struct {
	// Represents time when the Experiment was acknowledged by the Experiment controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the Experiment was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the Experiment was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`

	// List of observed runtime conditions for this Experiment.
	Conditions []ExperimentCondition `json:"conditions,omitempty"`

	// Current optimal trial parameters and observations.
	CurrentOptimalTrial OptimalTrial `json:"currentOptimalTrial,omitempty"`

	// List of trial names which are running.
	RunningTrialList []string `json:"runningTrialList,omitempty"`

	// List of trial names which are pending.
	PendingTrialList []string `json:"pendingTrialList,omitempty"`

	// List of trial names which have already failed.
	FailedTrialList []string `json:"failedTrialList,omitempty"`

	// List of trial names which have already succeeded.
	SucceededTrialList []string `json:"succeededTrialList,omitempty"`

	// List of trial names which have been killed.
	KilledTrialList []string `json:"killedTrialList,omitempty"`

	// List of trial names which have been early stopped.
	EarlyStoppedTrialList []string `json:"earlyStoppedTrialList,omitempty"`

	// List of trial names which have been metrics unavailable
	MetricsUnavailableTrialList []string `json:"metricsUnavailableTrialList,omitempty"`

	// Trials is the total number of trials owned by the experiment.
	Trials int32 `json:"trials,omitempty"`

	// How many trials have succeeded.
	TrialsSucceeded int32 `json:"trialsSucceeded,omitempty"`

	// How many trials have failed.
	TrialsFailed int32 `json:"trialsFailed,omitempty"`

	// How many trials have been killed.
	TrialsKilled int32 `json:"trialsKilled,omitempty"`

	// How many trials are currently pending.
	TrialsPending int32 `json:"trialsPending,omitempty"`

	// How many trials are currently running.
	TrialsRunning int32 `json:"trialsRunning,omitempty"`

	// How many trials are currently early stopped.
	TrialsEarlyStopped int32 `json:"trialsEarlyStopped,omitempty"`

	// How many trials are currently metrics unavailable.
	TrialMetricsUnavailable int32 `json:"trialMetricsUnavailable,omitempty"`
}

// OptimalTrial is the metrics and assignments of the best trial.
type OptimalTrial struct {
	// BestTrialName is the name of the best trial.
	BestTrialName string `json:"bestTrialName,omitempty"`
	// Key-value pairs for hyperparameters and assignment values.
	ParameterAssignments []common.ParameterAssignment `json:"parameterAssignments,omitempty"`

	// Observation for this trial
	Observation common.Observation `json:"observation,omitempty"`
}

// ExperimentCondition describes the state of the experiment at a certain point.
// +k8s:deepcopy-gen=true
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

// ResumePolicyType describes how the experiment should be resumed.
// Only one of the following resume policies may be specified.
// If none of the following policies is specified, the default one is LongRunning.
type ResumePolicyType string

const (
	// NeverResume indicates that experiment can't be resumed.
	NeverResume ResumePolicyType = "Never"
	// LongRunning indicates that experiment's suggestion resources
	// (deployment and service) are always running.
	LongRunning ResumePolicyType = "LongRunning"
	// FromVolume indicates that volume is attached to experiment's
	// suggestion. Suggestion data can be retained in the volume.
	// When experiment is succeeded suggestion deployment and service are deleted.
	FromVolume ResumePolicyType = "FromVolume"
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

// TrialTemplate describes structure of trial template
type TrialTemplate struct {
	// Retain indicates that trial resources must be not cleanup
	Retain bool `json:"retain,omitempty"`

	// Source for trial template (unstructured structure or config map)
	TrialSource `json:",inline"`

	// List of parameters that are used in trial template
	TrialParameters []TrialParameterSpec `json:"trialParameters,omitempty"`

	// Labels that determines if pod needs to be injected by Katib sidecar container.
	// If PrimaryPodLabels is omitted, metrics collector wraps all Trial's pods.
	PrimaryPodLabels map[string]string `json:"primaryPodLabels,omitempty"`

	// Name of training container where actual model training is running
	PrimaryContainerName string `json:"primaryContainerName,omitempty"`

	// Condition when trial custom resource is succeeded.
	// Condition must be in GJSON format, ref https://github.com/tidwall/gjson.
	// For example for BatchJob: status.conditions.#(type=="Complete")#|#(status=="True")#
	SuccessCondition string `json:"successCondition,omitempty"`

	// Condition when trial custom resource is failed.
	// Condition must be in GJSON format, ref https://github.com/tidwall/gjson.
	// For example for BatchJob: status.conditions.#(type=="Failed")#|#(status=="True")#
	FailureCondition string `json:"failureCondition,omitempty"`
}

// TrialSource represent the source for trial template
// Only one source can be specified
type TrialSource struct {

	// TrialSpec represents trial template in unstructured format
	TrialSpec *unstructured.Unstructured `json:"trialSpec,omitempty"`

	// ConfigMap spec represents a reference to ConfigMap
	ConfigMap *ConfigMapSource `json:"configMap,omitempty"`
}

// ConfigMapSource references the config map where trial template is located
type ConfigMapSource struct {
	// Name of config map where trial template is located
	ConfigMapName string `json:"configMapName,omitempty"`

	// Namespace of config map where trial template is located
	ConfigMapNamespace string `json:"configMapNamespace,omitempty"`

	// Path in config map where trial template is located
	TemplatePath string `json:"templatePath,omitempty"`
}

// TrialParameterSpec describes parameters that must be replaced in trial template
type TrialParameterSpec struct {
	// Name of the parameter that must be replaced in trial template
	Name string `json:"name,omitempty"`

	// Description of the parameter
	Description string `json:"description,omitempty"`

	// Reference to the parameter in search space
	Reference string `json:"reference,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Structure of the Experiment custom resource.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
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

// NasConfig contains config for NAS job
type NasConfig struct {
	GraphConfig GraphConfig `json:"graphConfig,omitempty"`
	Operations  []Operation `json:"operations,omitempty"`
}

// GraphConfig contains a config of DAG
type GraphConfig struct {
	NumLayers   *int32  `json:"numLayers,omitempty"`
	InputSizes  []int32 `json:"inputSizes,omitempty"`
	OutputSizes []int32 `json:"outputSizes,omitempty"`
}

// Operation contains type of operation in DAG
type Operation struct {
	OperationType string          `json:"operationType,omitempty"`
	Parameters    []ParameterSpec `json:"parameters,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Experiment{}, &ExperimentList{})
}
