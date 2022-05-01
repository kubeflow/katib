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
)

// AlgorithmSpec is the specification for a HP or NAS algorithm.
// +k8s:deepcopy-gen=true
type AlgorithmSpec struct {
	// HP or NAS algorithm name.
	AlgorithmName string `json:"algorithmName,omitempty"`

	// Key-value pairs representing settings for suggestion algorithms.
	AlgorithmSettings []AlgorithmSetting `json:"algorithmSettings,omitempty"`
}

// AlgorithmSetting represents key-value pair for HP or NAS algorithm settings.
type AlgorithmSetting struct {
	// Name is setting name.
	Name string `json:"name,omitempty"`

	// Value is the setting value.
	Value string `json:"value,omitempty"`
}

// EarlyStoppingSpec is the specification for a early stopping algorithm.
// +k8s:deepcopy-gen=true
type EarlyStoppingSpec struct {
	// Early stopping algorithm name.
	AlgorithmName string `json:"algorithmName,omitempty"`

	// Key-value pairs representing settings for early stopping algorithm.
	AlgorithmSettings []EarlyStoppingSetting `json:"algorithmSettings,omitempty"`
}

// EarlyStoppingSetting represents key-value pair for early stopping algorithm settings.
type EarlyStoppingSetting struct {
	// Name is the setting name.
	Name string `json:"name,omitempty"`

	// Value is the setting value.
	Value string `json:"value,omitempty"`
}

// EarlyStoppingRule represents each rule for early stopping.
type EarlyStoppingRule struct {
	// Name contains metric name for the rule.
	Name string `json:"name,omitempty"`

	// Value contains metric value for the rule.
	Value string `json:"value,omitempty"`

	// Comparison defines correlation between name and value.
	Comparison ComparisonType `json:"comparison,omitempty"`

	// StartStep defines quantity of intermediate results
	// that should be received before applying the rule.
	// If start step is empty, rule is applied from the first recorded metric.
	StartStep int `json:"startStep,omitempty"`
}

// ComparisonType is the type of comparison, one of equal, less or greater.
type ComparisonType string

const (
	// ComparisonTypeEqual means that metric value = early stopping rule value.
	ComparisonTypeEqual ComparisonType = "equal"

	// ComparisonTypeLess means that metric value < early stopping rule value.
	ComparisonTypeLess ComparisonType = "less"

	// ComparisonTypeGreater means that metric value > early stopping rule value.
	ComparisonTypeGreater ComparisonType = "greater"
)

// ObjectiveSpec represents Experiment's objective specification.
// +k8s:deepcopy-gen=true
type ObjectiveSpec struct {
	// Type for Experiment optimization.
	Type ObjectiveType `json:"type,omitempty"`

	// Goal is the Experiment's objective goal that should be reached.
	// In case of empty goal, Experiment is running until MaxTrialCount = TrialsSucceeded.
	Goal *float64 `json:"goal,omitempty"`

	// ObjectiveMetricName represents primary Experiment's metric to optimize.
	ObjectiveMetricName string `json:"objectiveMetricName,omitempty"`

	// AdditionalMetricNames represents metrics that should be collected from Trials.
	// This can be empty if we only care about the objective metric.
	// Note: If we adopt a push instead of pull mechanism, this can be omitted completely.
	AdditionalMetricNames []string `json:"additionalMetricNames,omitempty"`

	// MetricStrategies defines various rules (min, max or latest) to extract metrics values.
	// This field is allowed to missing, experiment defaulter (webhook) will fill it.
	MetricStrategies []MetricStrategy `json:"metricStrategies,omitempty"`
}

// ObjectiveType is the type of Experiment optimization, one of minimize or maximize.
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

// MetricStrategyType describes the various approaches to extract objective value from metrics.
type MetricStrategyType string

const (
	ExtractByMin    MetricStrategyType = "min"
	ExtractByMax    MetricStrategyType = "max"
	ExtractByLatest MetricStrategyType = "latest"
)

type MetricStrategy struct {
	Name  string             `json:"name,omitempty"`
	Value MetricStrategyType `json:"value,omitempty"`
}

type Metric struct {
	Name   string `json:"name,omitempty"`
	Min    string `json:"min,omitempty"`
	Max    string `json:"max,omitempty"`
	Latest string `json:"latest,omitempty"`
}

// +k8s:deepcopy-gen=true
type Observation struct {
	// Key-value pairs for metric names and values
	Metrics []Metric `json:"metrics,omitempty"`
}

// +k8s:deepcopy-gen=true
type MetricsCollectorSpec struct {
	Source    *SourceSpec    `json:"source,omitempty"`
	Collector *CollectorSpec `json:"collector,omitempty"`
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
	InvalidKind   FileSystemKind = "Invalid"
)

type FileFormat string

const (
	TextFormat FileFormat = "TEXT"
	JsonFormat FileFormat = "JSON"
)

// +k8s:deepcopy-gen=true
type FileSystemPath struct {
	Path   string         `json:"path,omitempty"`
	Kind   FileSystemKind `json:"kind,omitempty"`
	Format FileFormat     `json:"format,omitempty"`
}

type CollectorKind string

const (
	StdOutCollector CollectorKind = "StdOut"

	FileCollector   CollectorKind = "File"
	DefaultFilePath string        = "/var/log/katib/metrics.log"

	TfEventCollector              CollectorKind = "TensorFlowEvent"
	DefaultTensorflowEventDirPath string        = "/var/log/katib/tfevent/"

	PrometheusMetricCollector CollectorKind = "PrometheusMetric"
	DefaultPrometheusPath     string        = "/metrics"
	DefaultPrometheusPort     int           = 8080

	CustomCollector CollectorKind = "Custom"

	// When model training source code persists metrics into persistent layer
	// directly, metricsCollector isn't in need, and its kind is "noneCollector"
	NoneCollector CollectorKind = "None"

	MetricsVolume = "metrics-volume"
)

// +k8s:deepcopy-gen=true
type CollectorSpec struct {
	Kind CollectorKind `json:"kind,omitempty"`
	// When kind is "customCollector", this field will be used
	CustomCollector *v1.Container `json:"customCollector,omitempty"`
}
