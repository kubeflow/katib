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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//runtime "k8s.io/apimachinery/pkg/runtime"

	pb "github.com/kubeflow/katib/pkg/api"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StudyJobSpec defines the desired state of StudyJob
type StudyJobSpec struct {
	StudySpec            *StudySpec            `json:"studySpec,omitempty"`
	WorkerSpec           *WorkerSpec           `json:"workerSpec,omitempty"`
	SuggestionSpec       *SuggestionSpec       `json:"suggestionSpec,omitempty"`
	EarlyStoppingSpec    *EarlyStoppingSpec    `json:"earlyStoppingSpec,omitempty"`
	MetricsCollectorSpec *MetricsCollectorSpec `json:"metricsCollectorSpec,omitempty"`
}

// StudyJobStatus defines the observed state of StudyJob
type StudyJobStatus struct {
	// Represents time when the TFJob was acknowledged by the TFJob controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the TFJob was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the TFJob was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`

	Condition                Condition  `json:"conditon,omitempty"`
	StudyId                  string     `json:"studyid,omitempty"`
	SuggestionParameterId    string     `json:"suggestionParameterId"`
	EarlyStoppingParameterId string     `json:"earlyStoppingParameterId"`
	Trials                   []TrialSet `json:"trials,omitempty"`
}

type WorkerCondition struct {
	WorkerId  string    `json:"workerid,omitempty"`
	Condition Condition `json:"conditon,omitempty"`
}

type TrialSet struct {
	TrialId    string            `json:"trialid,omitempty"`
	WorkerList []WorkerCondition `json:"workeridlist,omitempty"`
}

type StudySpec struct {
	Name               string            `json:"name,omitempty"`
	Owner              string            `json:"owner,omitempty"`
	OptimizationType   OptimizationType  `json:"optimizationtype,omitempty"`
	OptimizationGoal   *float64          `json:"optimizationgoal,omitempty"`
	ObjectiveValueName string            `json:"objectivevaluename,omitempty"`
	MetricsNames       []string          `json:"metricsnames,omitempty"`
	ParameterConfigs   []ParameterConfig `json:"parameterconfigs,omitempty"`
}

type ParameterConfig struct {
	Name          string        `json:"name,omitempty"`
	ParameterType ParameterType `json:"parametertype,omitempty"`
	Feasible      FeasibleSpace `json:"feasible,omitempty"`
}

type FeasibleSpace struct {
	Max  string   `json:"max,omitempty"`
	Min  string   `json:"min,omitempty"`
	List []string `json:"list,omitempty"`
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

type WorkerSpec struct {
	WorkerTemplatePath string            `json:"workerTemplatePath,omitempty"`
	WorkerType         string            `json:"workerType,omitempty"`
	Image              string            `json:"image,omitempty"`
	Command            []string          `json:"command,omitempty"`
	VolumeConfigs      []VolumeConfig    `json:"volumeConfigs,omitempty"`
	WorkerParameters   map[string]string `json:"workerParameters,omitempty"`
}

type VolumeConfig struct {
	Name      string `json:"name,omitempty"`
	PvcName   string `json:"pvcName,omitempty"`
	MountPath string `json:"mountPath,omitempty"`
}

type MetricsCollectorSpec struct {
	MetricsCollectorTemplatePath string `json:"metricsCollectorTemplatePath,omitempty"`
	ServiceAccount               string `json"serviceAccount,omitempty"`
}

type ServiceParameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type SuggestionSpec struct {
	SuggestionAlgorithm  string                   `json:"suggestionAlgorithm,omitempty"`
	SuggestionParameters []pb.SuggestionParameter `json:"suggestionParameters"`
	RequestNumber        int                      `json:"requestNumber,omitempty"`
}

type EarlyStoppingSpec struct {
	EarlyStoppingAlgorithm  string                      `json:"earlyStoppingAlgorithm,omitempty"`
	EarlyStoppingParameters []pb.EarlyStoppingParameter `json:"earlyStoppingParameters"`
}

type ParameterEmbedding string

const (
	ParameterEmbeddingUndefined        ParameterEmbedding = ""
	ParameterEmbeddingArgument         ParameterEmbedding = "args"
	ParameterEmbeddingEnvironmentValue ParameterEmbedding = "env"
)

type Condition string

const (
	ConditionUnknown   Condition = "Unknown"
	ConditionCreated   Condition = "Created"
	ConditionRunning   Condition = "Running"
	ConditionCompleted Condition = "Completed"
	ConditionFailed    Condition = "Failed"
)

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
