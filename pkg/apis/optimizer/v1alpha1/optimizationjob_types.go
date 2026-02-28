package v1alpha1

import (
	trainerv1alpha1 "github.com/kubeflow/trainer/pkg/apis/trainer/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// ObjectiveDirection is the optimization direction for an objective metric.
// +kubebuilder:validation:Enum=minimize;maximize
type ObjectiveDirection string

const (
	ObjectiveDirectionMinimize ObjectiveDirection = "minimize"
	ObjectiveDirectionMaximize ObjectiveDirection = "maximize"
)

// Distribution defines the sampling distribution for a continuous parameter.
// +kubebuilder:validation:Enum=uniform;logUniform;normal;logNormal
type Distribution string

const (
	DistributionUniform    Distribution = "uniform"
	DistributionLogUniform Distribution = "logUniform"
	DistributionNormal     Distribution = "normal"
	DistributionLogNormal  Distribution = "logNormal"
)

// OptimizationJobConditionType defines the condition types for an OptimizationJob.
type OptimizationJobConditionType string

const (
	OptimizationJobInitializerReady OptimizationJobConditionType = "InitializerReady"
	OptimizationJobRunning          OptimizationJobConditionType = "Running"
	OptimizationJobSucceeded        OptimizationJobConditionType = "Succeeded"
	OptimizationJobFailed           OptimizationJobConditionType = "Failed"
)

// Objective defines the metric and goal for the HPO job.
type Objective struct {
	// Metric is the name of the metric to optimize (e.g., "accuracy", "loss").
	// +kubebuilder:validation:MinLength=1
	Metric string `json:"metric"`

	// Direction specifies whether to minimize or maximize the metric.
	// +kubebuilder:default=maximize
	Direction ObjectiveDirection `json:"direction"`

	// Goal is the target value for the metric. When reached, the optimization stops.
	// +optional
	Goal *float64 `json:"goal,omitempty"`
}

// Algorithm defines the optimization algorithm configuration.
type Algorithm struct {
	// Name is the optimization algorithm (e.g., "random", "bayesian", "tpe", "cmaes").
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Settings are algorithm-specific key-value parameters.
	// +optional
	// +listType=map
	// +listMapKey=name
	Settings []AlgorithmSetting `json:"settings,omitempty"`
}

// AlgorithmSetting is a key-value pair for algorithm configuration.
type AlgorithmSetting struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	Value string `json:"value"`
}

// TrialConfig controls the orchestration of the trials.
type TrialConfig struct {
	// NumTrials is the maximum number of trials to run.
	// +kubebuilder:validation:Minimum=1
	// +optional
	NumTrials *int32 `json:"numTrials,omitempty"`

	// ParallelTrials is how many trials can run concurrently.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	// +optional
	ParallelTrials *int32 `json:"parallelTrials,omitempty"`

	// MaxFailedTrials is the threshold of failures before marking the job as failed.
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxFailedTrials *int32 `json:"maxFailedTrials,omitempty"`
}

// ParameterSpec defines one hyperparameter and its search domain.
// Exactly one of Continuous, Categorical, or Discrete must be set.
type ParameterSpec struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Continuous defines a float-valued parameter with min/max bounds.
	// +optional
	Continuous *ContinuousParam `json:"continuous,omitempty"`

	// Categorical defines a parameter that takes one of a fixed set of string values.
	// +optional
	Categorical *CategoricalParam `json:"categorical,omitempty"`

	// Discrete defines a parameter that takes one of a fixed set of numeric values.
	// +optional
	Discrete *DiscreteParam `json:"discrete,omitempty"`
}

// ContinuousParam defines a float-valued search range.
type ContinuousParam struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`

	// Distribution controls how values are sampled within [min, max].
	// +kubebuilder:default=uniform
	// +optional
	Distribution Distribution `json:"distribution,omitempty"`
}

// CategoricalParam defines a set of allowed string values.
type CategoricalParam struct {
	// +kubebuilder:validation:MinItems=1
	Choices []string `json:"choices"`
}

// DiscreteParam defines a set of allowed numeric values.
type DiscreteParam struct {
	// +kubebuilder:validation:MinItems=1
	Values []float64 `json:"values"`
}

// MetricValue holds a single objective metric observation.
type MetricValue struct {
	// Metric is the name of the objective metric.
	Metric string `json:"metric"`
	// Value is the observed value.
	Value float64 `json:"value"`
}

// BestTrial tracks the best performing trial and its metrics.
type BestTrial struct {
	// Name is the name of the best-performing Trial / TrainJob.
	Name string `json:"name"`

	// Metrics are the observed objective metric values for this trial.
	Metrics []MetricValue `json:"metrics"`

	// OptimalParameters is the map of hyperparameter names to the values used by this trial.
	// +optional
	OptimalParameters map[string]string `json:"optimalParameters,omitempty"`
}

// OptimizationJobSpec defines the desired state of OptimizationJob.
type OptimizationJobSpec struct {
	// Objectives defines the metrics to optimize, their direction, and optional goal.
	// +kubebuilder:validation:MinItems=1
	Objectives []Objective `json:"objectives"`

	// Algorithm specifies the HPO algorithm and its settings.
	Algorithm Algorithm `json:"algorithm"`

	// SearchSpace defines the hyperparameter boundaries.
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=name
	SearchSpace []ParameterSpec `json:"searchSpace"`

	// TrialConfig controls parallelism, trial limits, and failure thresholds.
	TrialConfig TrialConfig `json:"trialConfig"`

	// Initializer runs once before any trials to download shared artifacts (models, datasets)
	// and stores them on a PVC that is mounted into every trial's TrainJob.
	// +optional
	Initializer *trainerv1alpha1.Initializer `json:"initializer,omitempty"`

	// TrialTemplate is the TrainJob manifest used as the template for each trial.
	// The controller substitutes search-space values using ${searchSpace.<paramName>} placeholders.
	// +kubebuilder:pruning:PreserveUnknownFields
	TrialTemplate runtime.RawExtension `json:"trialTemplate"`
}

// OptimizationJobStatus defines the observed state of OptimizationJob.
type OptimizationJobStatus struct {
	// Conditions track the overall lifecycle of the OptimizationJob.
	// Known condition types: InitializerReady, Running, Succeeded, Failed.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// StartTime is when the OptimizationJob controller first started processing this resource.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is when all trials finished (succeeded or hit failure threshold).
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Trial counters.
	Active    int32 `json:"active,omitempty"`
	Succeeded int32 `json:"succeeded,omitempty"`
	Failed    int32 `json:"failed,omitempty"`

	// BestTrial holds the best performing trial observed so far.
	// +optional
	BestTrial *BestTrial `json:"bestTrial,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=optjob
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.conditions[?(@.status=="True")].type`
// +kubebuilder:printcolumn:name="Best Metric",type=string,JSONPath=`.status.bestTrial.metrics[0].value`
// +kubebuilder:printcolumn:name="Succeeded",type=integer,JSONPath=`.status.succeeded`
// +kubebuilder:printcolumn:name="Failed",type=integer,JSONPath=`.status.failed`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// OptimizationJob is the Schema for the optimizationjobs API.
type OptimizationJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OptimizationJobSpec   `json:"spec,omitempty"`
	Status OptimizationJobStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// OptimizationJobList contains a list of OptimizationJob.
type OptimizationJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OptimizationJob `json:"items"`
}
