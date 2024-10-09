# KEP-507: Suggestion CRD Design Document

# Table of Contents

- [Suggestion CRD Design Document](#suggestion-crd-design-document)
- [Table of Contents](#table-of-contents)
  - [Background](#background)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
  - [Design](#design)
    - [Kubernetes API](#kubernetes-api)
    - [GRPC API](#grpc-api)
    - [Workflow](#workflow)
      - [Example](#example)
  - [Algorithm Supports](#algorithm-supports)
    - [Random](#random)
    - [Grid](#grid)
    - [Bayes Optimization](#bayes-optimization)
    - [HyperBand](#hyperband)
    - [BOHB](#bohb)
    - [TPE](#tpe)
    - [SMAC](#smac)
    - [CMA-ES](#cma-es)
    - [Sobol](#sobol)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)

## Background

Katib makes suggestions long-running in v1alpha3. And the suggestions need to communicate with Katib DB manager to get experiments and trials from Katib db driver. This design hurts high availability.

Thus we proposed a new design to implement a CRD for suggestion and remove Katib db communication from main workflow. The new design simplifies the implementation of experiment and trial controller, and makes Katib Kubernetes native.

This document is to illustrate the details of the new design.

## Goals

- Propose the Suggestion CRD.
- Propose new GRPC API for Suggestion service.
- Suggest the approaches to implement suggestion algorithms.

## Non-Goals

- Metrics collection (See [Metrics Collector Design Document](./metrics-collector.md))
- Database-related refactor

## Design

### Kubernetes API

```go
// SuggestionSpec defines the desired state of Suggestion
type SuggestionSpec struct {
	AlgorithmName string `json:"algorithmName"`
	// Number of suggestions requested
	Requests int32 `json:"requests,omitempty"`
}

// SuggestionStatus defines the observed state of Suggestion
type SuggestionStatus struct {
	// Algorithmsettings set by the algorithm services.
	AlgorithmSettings []common.AlgorithmSetting `json:"algorithmSettings,omitempty"`

	// Number of suggestion results
	SuggestionCount int32 `json:"suggestionCount,omitempty"`

	// Suggestion results
	Suggestions []TrialAssignment `json:"suggestions,omitempty"`
}

// TrialAssignment is the assignment for one trial.
type TrialAssignment struct {
	// Suggestion results
	ParameterAssignments []common.ParameterAssignment `json:"parameterAssignments,omitempty"`

	//Name of the suggestion
	Name string `json:"name,omitempty"`
}
```

### GRPC API

```protobuf
syntax = "proto3";

package api.v1.alpha3;

import "google/api/annotations.proto";

service Suggestion {
    rpc GetSuggestions(GetSuggestionsRequest) returns (GetSuggestionsReply);
}

message GetSuggestionsRequest {
    Experiment experiment = 1;
    repeated Trial trials = 2; // all completed trials owned by the experiment.
    int32 request_number = 3; ///The number of Suggestion you request at one time. When you set 3 to request_number, you can get three Suggestions at one time.
}

message GetSuggestionsReply {
    message ParameterAssignments{
        repeated ParameterAssignment assignments = 1;
    }
    repeated ParameterAssignments parameter_assignments = 1;
    AlgorithmSpec algorithm = 2;
}

message Experiment {
    string name = 1;
    ExperimentSpec experiment_spec = 2;
}

message ExperimentSpec {
   AlgorithmSpec algorithm = 3;
   ParameterSpecs parameter_specs = 1;
   ObjectiveSpec objective = 2;
}

message ParameterSpecs {
    repeated ParameterSpec parameters = 1;
}

message AlgorithmSpec {
    string algorithm_name = 1;
    repeated AlgorithmSetting algorithm_settings = 2;
}

message AlgorithmSetting {
    string name = 1;
    string value = 2;
}

message ParameterSpec {
    string name = 1; /// Name of the parameter.
    ParameterType parameter_type = 2; /// Type of the parameter.
    FeasibleSpace feasible_space = 3; /// FeasibleSpace for the parameter.
}

message FeasibleSpace {
    string max = 1; /// Max Value
    string min = 2; /// Minimum Value
    repeated string list = 3; /// List of Values.
    string step = 4; /// Step for double or int parameter
}

enum ParameterType {
    UNKNOWN_TYPE = 0; /// Undefined type and not used.
    DOUBLE = 1; /// Double float type. Use "Max/Min".
    INT = 2; /// Int type. Use "Max/Min".
    DISCRETE = 3; /// Discrete number type. Use "List" as float.
    CATEGORICAL = 4; /// Categorical type. Use "List" as string.
}

enum ObjectiveType {
    UNKNOWN = 0; /// Undefined type and not used.
    MINIMIZE = 1; /// Minimize
    MAXIMIZE = 2; /// Maximize
}

message ObjectiveSpec {
    ObjectiveType type = 1;
    double goal = 2;
    string objective_metric_name = 3;
}

message Trial {
   string name = 1;
   TrialSpec spec = 2;
   TrialStatus status = 3;
}

message TrialSpec {
   ParameterAssignments parameter_assignments = 2;
   string run_spec = 3;
}

message ParameterAssignments {
    repeated ParameterAssignment assignments = 1;
}

message ParameterAssignment {
    string name = 1;
    string value = 2;
}

message TrialStatus {
   Observation observation = 4; // The best observation in logs.
}

message Observation {
    repeated Metric metrics = 1;
}

message Metric {
    string name = 1;
    string value = 2;
}
```

### Workflow

![](../images/katib-workflow.png)

When the user creates a Experiment, we will create a Suggestion for the Experiment. When the Experiment needs some suggestions, Experiment controller updates the `Suggestions`, then Suggestion controller communicates with the Suggestion to get parameter assignments and set them in Suggestion status.

#### Example

Now the workflow will be illustrated with an example.

```yaml
apiVersion: "kubeflow.org/v1alpha3"
kind: Experiment
metadata:
  namespace: kubeflow
  name: random-experiment
spec:
  parallelTrialCount: 3
  maxTrialCount: 12
  maxFailedTrialCount: 3
  objective:
    type: maximize
    goal: 0.99
    objectiveMetricName: Validation-accuracy
    additionalMetricNames:
      - accuracy
  algorithm:
    algorithmName: random
  trialTemplate:
    goTemplate:
      rawTemplate: |-
        apiVersion: batch/v1
        kind: Job
        metadata:
          name: {{.Trial}}
          namespace: {{.NameSpace}}
        spec:
          template:
            spec:
              containers:
              - name: {{.Trial}}
                image: katib/mxnet-mnist-example
                command:
                - "python"
                - "/mxnet/example/image-classification/train_mnist.py"
                - "--batch-size=64"
                {{- with .HyperParameters}}
                {{- range .}}
                - "{{.Name}}={{.Value}}"
                {{- end}}
                {{- end}}
              restartPolicy: Never
  parameters:
    - name: --lr
      parameterType: double
      feasibleSpace:
        min: "0.01"
        max: "0.03"
    - name: --num-layers
      parameterType: int
      feasibleSpace:
        min: "2"
        max: "5"
    - name: --optimizer
      parameterType: categorical
      feasibleSpace:
        list:
          - sgd
          - adam
          - ftrl
```

Then, Experiment controller needs 3 parallel trials to run. It creates the Suggestions:

```yaml
apiVersion: "kubeflow.org/v1alpha3"
kind: Suggestion
metadata:
  namespace: kubeflow
  name: random-experiment
spec:
  algorithmName: random
  requests: 3
```

After that, Suggestion controller communicates with the Suggestion via GRPC and updates the status:

```yaml
apiVersion: "kubeflow.org/v1alpha3"
kind: Suggestion
metadata:
  namespace: kubeflow
  name: random-experiment
spec:
  algorithmName: random
  requests: 3
status:
  suggestions:
    - assignments:
        - name: --lr
          value: 0.02
        - name: --num-layers
          value: 4
        - name: --optimizer
          value: sgd
    - assignments:
        - name: --lr
          value: 0.021
        - name: --num-layers
          value: 3
        - name: --optimizer
          value: adam
    - assignments:
        - name: --lr
          value: 0.03
        - name: --num-layers
          value: 5
        - name: --optimizer
          value: adam
```

Then Experiment controller creates the trial. When there is one trial finished, Experiment controller will ask Suggestion controller for a new suggestion:

```yaml
apiVersion: "kubeflow.org/v1alpha3"
kind: Suggestion
metadata:
  namespace: kubeflow
  name: random-experiment
spec:
  algorithmName: random
  requests: 4
status:
  suggestions:
    - assignments:
        - name: --lr
          value: 0.02
        - name: --num-layers
          value: 4
        - name: --optimizer
          value: sgd
    - assignments:
        - name: --lr
          value: 0.021
        - name: --num-layers
          value: 3
        - name: --optimizer
          value: adam
    - assignments:
        - name: --lr
          value: 0.03
        - name: --num-layers
          value: 5
        - name: --optimizer
          value: adam
    - assignments:
        - name: --lr
          value: 0.012
        - name: --num-layers
          value: 4
        - name: --optimizer
          value: adam
```

## Algorithm Supports

### Random

We can use the implementation in Katib or [hyperopt](https://github.com/hyperopt/hyperopt).

### Grid

We can use the length of the trials to know which grid we are in. Please refer to the [implementation in advisor](https://github.com/tobegit3hub/advisor/blob/master/advisor_server/suggestion/algorithm/grid_search.py).

Or we can use [chocolate](https://github.com/AIworx-Labs/chocolate).

### Bayes Optimization

We can use [skopt](https://github.com/scikit-optimize/scikit-optimize) to run bayes optimization.

### HyperBand

We can use [HpBandSter](https://github.com/automl/HpBandSter) to run HyperBand.

### BOHB

We can use [HpBandSter](https://github.com/automl/HpBandSter) to run BOHB.

### TPE

We can use [hyperopt](https://github.com/hyperopt/hyperopt) to run TPE.

### SMAC

We can use [SMAC3](https://github.com/automl/SMAC3) to run SMAC.

### CMA-ES

We can use [goptuna](https://github.com/c-bata/goptuna) to run CMA-ES.

### Sobol

We can use [goptuna](https://github.com/c-bata/goptuna) to run Sobol.
