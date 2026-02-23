// pkg/apis/optimizer/v1alpha1/optimizationjob_types.go

package v1alpha1

import (
	trainerv1alpha1 "github.com/kubeflow/trainer/pkg/apis/trainer/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Objective defines the metric and goal for the HPO job.
type Objective struct {
	Metric    string   `json:"metric"`
	Direction string   `json:"direction"`
	Goal      *float64 `json:"goal,omitempty"`
}

// Algorithm defines the optimization algorithm configuration.
type Algorithm struct {
	Name     string      `json:"name"`
	Settings []SettingKV `json:"settings,omitempty"`
}

// SettingKV is a key-value pair for algorithm settings.
type SettingKV struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TrialConfig controls the orchestration of the trials.
type TrialConfig struct {
	NumTrials       *int32 `json:"num_trials,omitempty"`
	ParallelTrials  *int32 `json:"parallel_trials,omitempty"`
	MaxFailedTrials *int32 `json:"max_failed_trials,omitempty"`
}

// OptimizationJobSpec defines the desired state of OptimizationJob.
type OptimizationJobSpec struct {
	Objectives []Objective `json:"objectives"`
	Algorithm  Algorithm   `json:"algorithm"`

	// Using map[string]string initially, can be refined to strict types later if needed.
	SearchSpace map[string]string `json:"searchSpace"`

	TrialConfig TrialConfig `json:"trialConfig"`

	// Tighter TrainJob Integration: Strongly typed to TrainJob rather than arbitrary CRDs.
	TrialTemplate trainerv1alpha1.TrainJob `json:"trialTemplate"`
}

// OptimizationJobStatus defines the observed state of OptimizationJob.
type OptimizationJobStatus struct {
	// Add status fields here (e.g., Conditions, BestTrial, etc.) as the controller matures.
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

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
