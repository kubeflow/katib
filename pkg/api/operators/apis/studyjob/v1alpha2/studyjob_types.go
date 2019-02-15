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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//runtime "k8s.io/apimachinery/pkg/runtime"

	pb "github.com/kubeflow/katib/pkg/api"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StudyJobSpec defines the desired state of StudyJob
type StudyJobSpec struct {
	Owner                string                `json:"owner,omitempty"`
	OptimizationType     OptimizationType      `json:"optimizationType,omitempty"`
	OptimizationGoal     *float64              `json:"optimizationGoal,omitempty"`
	ObjectiveValueName   string                `json:"objectiveValueName,omitempty"`
	MaxSuggestionCount   int                   `json:"maxSuggestionCount,omitempty"`
	MetricsNames         []string              `json:"metricsNames,omitempty"`
	ParameterConfigs     []ParameterConfig     `json:"parameterConfigs,omitempty"`
	WorkerSpec           *WorkerSpec           `json:"workerSpec,omitempty"`
	SuggestionSpec       *SuggestionSpec       `json:"suggestionSpec,omitempty"`
	EarlyStoppingSpec    *EarlyStoppingSpec    `json:"earlyStoppingSpec,omitempty"`
	MetricsCollectorSpec *MetricsCollectorSpec `json:"metricsCollectorSpec,omitempty"`

	// NAS is still in early development; its API design will be a separate discussion.
	//NasConfig            *NasConfig            `json:"nasConfig,omitempty"`

	// See #352
	ReuseStudyID         string                `json:"reuseStudyId,omitempty"`
}

// +k8s:deepcopy-gen=true
// JobCondition describes the state of the job at a certain point.
type JobCondition struct {
	// Type of job condition.
	Type JobConditionType `json:"type"`
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

// JobConditionType defines all kinds of types of JobStatus.
type JobConditionType string

const (
	JobCreated JobConditionType = "Created"
	JobRunning JobConditionType = "Running"
	JobRestarting JobConditionType = "Restarting"
	JobSucceeded JobConditionType = "Succeeded"
	JobFailed JobConditionType = "Failed"
)

// StudyJobStatus defines the observed state of StudyJob
type StudyJobStatus struct {
	// Represents time when the StudyJob was acknowledged by the StudyJob controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the StudyJob was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the StudyJob was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`

	Conditions               []JobCondition  `json:"conditions,omitempty"`
	StudyID                  string     `json:"studyId,omitempty"`
	SuggestionParameterID    string     `json:"suggestionParameterId,omitempty"`
	EarlyStoppingParameterID string     `json:"earlyStoppingParameterId,omitempty"`
	Trials                   []Trial    `json:"trials,omitempty"`
	BestObjectiveValue       *float64   `json:"bestObjectiveValue,omitempty"`
	BestTrialID              string     `json:"bestTrialId,omitempty"`
	BestWorkerID             string     `json:"bestWorkerId,omitempty"`
	BestParameterValues      []ParameterValue    `json:"bestParameterValues,omitempty"`
	CurrentSuggestionCount   int        `json:"currentSuggestionCount,omitempty"`
}

// +k8s:deepcopy-gen=true
// WorkerCondition describes the state of the worker at a certain point.
type WorkerCondition struct {
	// Type of job condition.
	Type WorkerConditionType `json:"type"`
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

// WorkerConditionType defines all kinds of types of WorkerStatus.
type WorkerConditionType string

const (
	WorkerPending WorkerConditionType = "Pending"
	WorkerRunning WorkerConditionType = "Running"
	WorkerCompleted WorkerConditionType = "Completed"
	WorkerKilled WorkerConditionType = "Killed"
	WorkerFailed WorkerConditionType = "Failed"
)

type WorkerMetadata struct {
	WorkerID       string      `json:"workerId,omitempty"`
	Kind           string      `json:"kind,omitempty"`
	Conditions     []WorkerCondition   `json:"conditions,omitempty"`
	ObjectiveValue *float64    `json:"objectiveValue,omitempty"`
	StartTime      metav1.Time `json:"startTime,omitempty"`
	CompletionTime metav1.Time `json:"completionTime,omitempty"`
}

type Trial struct {
	TrialID    string            `json:"trialId,omitempty"`
	Workers    []WorkerMetadata `json:"workers,omitempty"`
}

type ParameterConfig struct {
	Name          string        `json:"name,omitempty"`
	ParameterType ParameterType `json:"parameterType,omitempty"`
	Feasible      FeasibleSpace `json:"feasible,omitempty"`
}

type ParameterValue struct {
	Name          string        `json:"name,omitempty"`
	ParameterType ParameterType `json:"parameterType,omitempty"`
	Value         string `json:"value,omitempty"`
}

type FeasibleSpace struct {
	Max  string   `json:"max,omitempty"`
	Min  string   `json:"min,omitempty"`
	List []string `json:"list,omitempty"`
	Step string   `json:"step,omitempty"`
}

type ParameterType string

const (
	ParameterTypeUnknown     ParameterType = "unknown"
	ParameterTypeDouble      ParameterType = "double"
	ParameterTypeInt         ParameterType = "int"
	ParameterTypeDiscrete    ParameterType = "discrete"
	ParameterTypeCategorical ParameterType = "categorical"
)

type OptimizationType string

const (
	OptimizationTypeUnknown  OptimizationType = ""
	OptimizationTypeMinimize OptimizationType = "minimize"
	OptimizationTypeMaximize OptimizationType = "maximize"
)

type TemplateSpec struct {
	ConfigMapName string `json:"configMapName,omitempty"`
	ConfigMapNamespace string `json:"configMapNamespace,omitempty"`
	TemplatePath string `json:"templatePath,omitempty"`
}

type GoTemplate struct {
	TemplateSpec *TemplateSpec `json:"templateSpec,omitempty"`
	RawTemplate  string `json:"rawTemplate,omitempty"`
}

type WorkerSpec struct {
	Retain     bool       `json:"retain,omitempty"`
	GoTemplate *GoTemplate `json:"goTemplate,omitempty"`
}

type MetricsCollectorSpec struct {
	Retain     bool       `json:"retain,omitempty"`
	GoTemplate *GoTemplate `json:"goTemplate,omitempty"`
}

type SuggestionSpec struct {
	SuggestionAlgorithm  string                   `json:"suggestionAlgorithm,omitempty"`
	SuggestionParameters []pb.SuggestionParameter `json:"suggestionParameters"`
	NumParallelTrials    int                      `json:"numParallelTrials,omitempty"`
}

type EarlyStoppingSpec struct {
	EarlyStoppingAlgorithm  string                      `json:"earlyStoppingAlgorithm,omitempty"`
	EarlyStoppingParameters []pb.EarlyStoppingParameter `json:"earlyStoppingParameters"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StudyJob is the Schema for the studyjob API
// +k8s:openapi-gen=true
type StudyJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StudyJobSpec   `json:"spec,omitempty"`
	Status StudyJobStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StudyJobList contains a list of StudyJob
type StudyJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StudyJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StudyJob{}, &StudyJobList{})
}
