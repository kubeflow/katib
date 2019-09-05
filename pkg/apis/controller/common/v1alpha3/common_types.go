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

import (
	v1 "k8s.io/api/core/v1"
)

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

// +k8s:deepcopy-gen=true
type SourceSpec struct {
	// Model-train source code can expose metrics by http, such as HTTP endpoint in
	// prometheus metric format
	HttpGet *v1.HTTPGetAction `json:"httpGet,omitempty"`
	// During training model, metrics may be persisted into local file in source
	// code, such as tfEvent use case
	FileSystemPath *FileSystemPath `json:"fileSystemPath,omitempty"`
	// Default metric output format is {"metric": "<metric_name>",
	// "value": <int_or_float>, "epoch": <int>, "step": <int>}, but if the output doesn't
	// follow default format, please extend it here
	Filter *FilterSpec `json:"filter,omitempty"`
}

// +k8s:deepcopy-gen=true
type FilterSpec struct {
	// When the metrics output follows format as this field specified, metricsCollector
	// collects it and reports to metrics server, it can be "<metric_name>: <float>" or else
	MetricsFormat []string `json:"metricsFormat,omitempty"`
}

type FileSystemKind string

const (
	DirectoryKind FileSystemKind = "Directory"
	FileKind      FileSystemKind = "File"
)

// +k8s:deepcopy-gen=true
type FileSystemPath struct {
	Path string         `json:"path,omitempty"`
	Kind FileSystemKind `json:"kind,omitempty"`
}

type CollectorKind string

const (
	StdOutCollector           CollectorKind = "StdOut"
	FileCollector             CollectorKind = "File"
	TfEventCollector          CollectorKind = "TensorflowEvent"
	PrometheusMetricCollector CollectorKind = "PrometheusMetric"
	CustomCollector           CollectorKind = "Custom"
	// When model training source code persists metrics into persistent layer
	// directly, metricsCollector isn't in need, and its kind is "noneCollector"
	NoneCollector CollectorKind = "None"
)

// +k8s:deepcopy-gen=true
type CollectorSpec struct {
	Kind CollectorKind `json:"kind"`
	// When kind is "customCollector", this field will be used
	CustomCollector *v1.Container `json:"customCollector,omitempty"`
}
