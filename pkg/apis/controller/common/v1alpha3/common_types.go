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

package v1alpha3

// +k8s:deepcopy-gen=true
type AlgorithmSpec struct {
	AlgorithmName string `json:"algorithmName,omitempty"`
	// Key-value pairs representing settings for suggestion algorithms.
	AlgorithmSettings []AlgorithmSetting `json:"algorithmSettings"`
	EarlyStopping     *EarlyStoppingSpec `json:"earlyStopping,omitempty"`
}

type AlgorithmSetting struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type EarlyStoppingSpec struct {
	// TODO
}

// +k8s:deepcopy-gen=true
type ObjectiveSpec struct {
	Type                ObjectiveType `json:"type,omitempty"`
	Goal                *float64      `json:"goal,omitempty"`
	ObjectiveMetricName string        `json:"objectiveMetricName,omitempty"`
	// This can be empty if we only care about the objective metric.
	// Note: If we adopt a push instead of pull mechanism, this can be omitted completely.
	AdditionalMetricNames []string `json:"additionalMetricNames,omitempty"`
}

type ObjectiveType string

const (
	ObjectiveTypeUnknown  ObjectiveType = ""
	ObjectiveTypeMinimize ObjectiveType = "minimize"
	ObjectiveTypeMaximize ObjectiveType = "maximize"
)

type ParameterAssignment struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Metric struct {
	Name  string  `json:"name,omitempty"`
	Value float64 `json:"value,omitempty"`
}

// +k8s:deepcopy-gen=true
type Observation struct {
	// Key-value pairs for metric names and values
	Metrics []Metric `json:"metrics"`
}
