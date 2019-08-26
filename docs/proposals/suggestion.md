# Suggestion CRD Design Proposal

Table of Contents
=================

   * [Suggestion CRD Design Proposal](#suggestion-crd-design-proposal)
      * [Background](#background)
      * [Goals](#goals)
      * [Non-Goals](#non-goals)
      * [Design](#design)
         * [Kubernetes API](#kubernetes-api)
         * [GRPC API (TODO)](#grpc-api-todo)
         * [Workflow](#workflow)
            * [Example](#example)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)


## Background

## Goals

## Non-Goals

## Design

### Kubernetes API

```go
type SuggestionSpec struct {
	//Name of the algorithm
	AlgorithmName string `json:"algorithm_name"`

	// Number of suggestions requested
	Suggestions int `json:"suggestions"`

	//Algorithm settings set by the user in the experiment config
	AlgorithmSettings []AlgorithmSetting `json:"algorithm_settings,omitempty"`
}

type AlgorithmSetting struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type SuggestionStatus struct {
	// Suggestion results
	Assignments []TrialAssignment `json:"assignments,omitempty"`

	Conditions []SuggestionCondition `json:"conditions,omitempty"`
	// include all common fields

}

type TrialAssignment struct {
	// Suggestion results
	Assignments []ParameterAssignment `json:"assignments,omitempty"`
	Name        string                `json:"name,omitempty"`
}

type ParameterAssignment struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Suggestion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SuggestionSpec   `json:"spec,omitempty"`
	Status SuggestionStatus `json:"status,omitempty"`
}
```

### GRPC API (TODO)

### Workflow

![](../images/suggestion-workflow.png)

When the user creates a Experiment, we will create a Suggestion for the Experiment. When the Experiment needs some suggestions, Experiment controller updates the `Suggestions`, then Suggestion controller communicates with the Suggestion to get parameter assignments and set them in Suggestion status.

#### Example

Now the workflow will be illustrated with an example.

```yaml
apiVersion: "kubeflow.org/v1alpha2"
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

Now, we will create a Suggestion for the Experiment:

```yaml
apiVersion: "kubeflow.org/v1alpha2"
kind: Suggestion
metadata:
  namespace: kubeflow
  name: random-experiment
spec:
  algorithmName: random
  suggestions: 0
```

Then, Experiment controller needs 3 parallel trials to run. It updates the Suggestions:

```yaml
apiVersion: "kubeflow.org/v1alpha2"
kind: Suggestion
metadata:
  namespace: kubeflow
  name: random-experiment
spec:
  algorithmName: random
  suggestions: 3
```

After that, Suggestion controller communicates with the Suggestion and updates the status:

```yaml
apiVersion: "kubeflow.org/v1alpha2"
kind: Suggestion
metadata:
  namespace: kubeflow
  name: random-experiment
spec:
  algorithmName: random
  suggestions: 3
status:
  assignments:
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

Then Experiment controller creates the trials and set Suggestion status (Optional, not decided):

```yaml
apiVersion: "kubeflow.org/v1alpha2"
kind: Suggestion
metadata:
  namespace: kubeflow
  name: random-experiment
spec:
  algorithmName: random
  suggestions: 3
status:
  assignments:
    - name: random-experiment-fsa2f
      assignments:
        - name: --lr
          value: 0.02
        - name: --num-layers
          value: 4
        - name: --optimizer
          value: sgd
    - name: random-experiment-hoj53
      assignments:
        - name: --lr
          value: 0.021
        - name: --num-layers
          value: 3
        - name: --optimizer
          value: adam
    - name: random-experiment-12lkj
      assignments:
        - name: --lr
          value: 0.03
        - name: --num-layers
          value: 5
        - name: --optimizer
          value: adam
```

When there is one trial finished, Experiment controller will ask Suggestion controller for a new suggestion:

```yaml
apiVersion: "kubeflow.org/v1alpha2"
kind: Suggestion
metadata:
  namespace: kubeflow
  name: random-experiment
spec:
  algorithmName: random
  suggestions: 4
status:
  assignments:
    - name: random-experiment-fsa2f
      assignments:
        - name: --lr
          value: 0.02
        - name: --num-layers
          value: 4
        - name: --optimizer
          value: sgd
    - name: random-experiment-hoj53
      assignments:
        - name: --lr
          value: 0.021
        - name: --num-layers
          value: 3
        - name: --optimizer
          value: adam
    - name: random-experiment-12lkj
      assignments:
        - name: --lr
          value: 0.03
        - name: --num-layers
          value: 5
        - name: --optimizer
          value: adam
    - name: random-experiment-ybfd3
      assignments:
        - name: --lr
          value: 0.012
        - name: --num-layers
          value: 4
        - name: --optimizer
          value: adam
```