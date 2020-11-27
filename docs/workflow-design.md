# How Katib v1beta1 tunes hyperparameter automatically in a Kubernetes native way

Follow the Kubeflow documentation guides:

- [Concepts](https://www.kubeflow.org/docs/components/katib/overview/)
  in Katib, hyperparameter tuning, and neural architecture search.
- [Getting started with Katib](https://kubeflow.org/docs/components/katib/hyperparameter/).
- Detailed guide to
  [configuring and running a Katib `Experiment`](https://kubeflow.org/docs/components/katib/experiment/).

## Example and Illustration

After install Katib v1beta1, you can run
`kubectl apply -f katib/examples/v1beta1/random-example.yaml` to try the first
example of Katib.

### Experiment

When you want to tune hyperparameters for your machine learning model before
training it further, you just need to create an `Experiment` CR. To
learn what fields are included in the `Experiment.spec`, follow
the detailed guide to
[configuring and running a Katib `Experiment`](https://kubeflow.org/docs/components/katib/experiment/).
Then you can get the new `Experiment` as below.
Katib concepts are introduced based on this example.

```yaml
$ kubectl get experiment random-example -n kubeflow -o yaml

apiVersion: kubeflow.org/v1beta1
kind: Experiment
metadata:
  ...
  name: random-example
  namespace: kubeflow
  ...
spec:
  algorithm:
    algorithmName: random
  maxFailedTrialCount: 3
  maxTrialCount: 12
  metricsCollectorSpec:
    collector:
      kind: StdOut
  objective:
    additionalMetricNames:
    - Train-accuracy
    goal: 0.99
    metricStrategies:
    - name: Validation-accuracy
      value: max
    - name: Train-accuracy
      value: max
    objectiveMetricName: Validation-accuracy
    type: maximize
  parallelTrialCount: 3
  parameters:
  - feasibleSpace:
      max: "0.03"
      min: "0.01"
    name: lr
    parameterType: double
  - feasibleSpace:
      max: "5"
      min: "2"
    name: num-layers
    parameterType: int
  - feasibleSpace:
      list:
      - sgd
      - adam
      - ftrl
    name: optimizer
    parameterType: categorical
  resumePolicy: LongRunning
  trialTemplate:
    failureCondition: status.conditions.#(type=="Failed")#|#(status=="True")#
    primaryContainerName: training-container
    successCondition: status.conditions.#(type=="Complete")#|#(status=="True")#
    trialParameters:
    - description: Learning rate for the training model
      name: learningRate
      reference: lr
    - description: Number of training model layers
      name: numberLayers
      reference: num-layers
    - description: Training model optimizer (sdg, adam or ftrl)
      name: optimizer
      reference: optimizer
    trialSpec:
      apiVersion: batch/v1
      kind: Job
      spec:
        template:
          spec:
            containers:
            - command:
              - python3
              - /opt/mxnet-mnist/mnist.py
              - --batch-size=64
              - --lr=${trialParameters.learningRate}
              - --num-layers=${trialParameters.numberLayers}
              - --optimizer=${trialParameters.optimizer}
              image: docker.io/kubeflowkatib/mxnet-mnist:v1beta1-e294a90
              name: training-container
            restartPolicy: Never
status:
  completionTime: "2020-11-16T20:13:02Z"
  conditions:
  - lastTransitionTime: "2020-11-16T20:00:15Z"
    lastUpdateTime: "2020-11-16T20:00:15Z"
    message: Experiment is created
    reason: ExperimentCreated
    status: "True"
    type: Created
  - lastTransitionTime: "2020-11-16T20:13:02Z"
    lastUpdateTime: "2020-11-16T20:13:02Z"
    message: Experiment is running
    reason: ExperimentRunning
    status: "False"
    type: Running
  - lastTransitionTime: "2020-11-16T20:13:02Z"
    lastUpdateTime: "2020-11-16T20:13:02Z"
    message: Experiment has succeeded because max trial count has reached
    reason: ExperimentMaxTrialsReached
    status: "True"
    type: Succeeded
  currentOptimalTrial:
    bestTrialName: random-example-gnz5nccf
    observation:
      metrics:
      - latest: "0.979299"
        max: "0.979299"
        min: "0.955115"
        name: Validation-accuracy
      - latest: "0.993503"
        max: "0.993503"
        min: "0.912413"
        name: Train-accuracy
    parameterAssignments:
    - name: lr
      value: "0.01874909352953323"
    - name: num-layers
      value: "5"
    - name: optimizer
      value: sgd
  startTime: "2020-11-16T20:00:15Z"
  succeededTrialList:
  - random-example-2fpnqfv8
  - random-example-2s9vfb9s
  - random-example-5hxm45x4
  - random-example-8xmpj4gv
  - random-example-b6gnl4cs
  - random-example-ftm2v84q
  - random-example-gnz5nccf
  - random-example-p74tn9gk
  - random-example-q6jrlshx
  - random-example-tkk46c4x
  - random-example-w5qgblgk
  - random-example-xcnrpx4x
  trials: 12
  trialsSucceeded: 12
```

### Suggestion

Katib internally creates a `Suggestion` CR for each `Experiment` CR. The
`Suggestion` CR includes the hyperparameter algorithm name by `algorithmName`
field and how many sets of hyperparameter Katib asks to be generated by
`requests` field. The `Suggestion` also traces all already generated sets of
hyperparameter in `status.suggestions`. The `Suggestion` CR is used for internal
logic control and end user can even ignore it.

```yaml
$ kubectl get suggestion random-example -n kubeflow -o yaml

apiVersion: kubeflow.org/v1beta1
kind: Suggestion
metadata:
  ...
  name: random-example
  namespace: kubeflow
  ownerReferences:
  - apiVersion: kubeflow.org/v1beta1
    blockOwnerDeletion: true
    controller: true
    kind: Experiment
    name: random-example
    uid: 302e79ae-8659-4679-9e2d-461209619883
  ...
spec:
  algorithm:
    algorithmName: random
  requests: 12
  resumePolicy: LongRunning
status:
  conditions:
  - lastTransitionTime: "2020-11-16T20:00:15Z"
    lastUpdateTime: "2020-11-16T20:00:15Z"
    message: Suggestion is created
    reason: SuggestionCreated
    status: "True"
    type: Created
  - lastTransitionTime: "2020-11-16T20:00:36Z"
    lastUpdateTime: "2020-11-16T20:00:36Z"
    message: Deployment is ready
    reason: DeploymentReady
    status: "True"
    type: DeploymentReady
  - lastTransitionTime: "2020-11-16T20:00:38Z"
    lastUpdateTime: "2020-11-16T20:00:38Z"
    message: Suggestion is running
    reason: SuggestionRunning
    status: "True"
    type: Running
  startTime: "2020-11-16T20:00:15Z"
  suggestionCount: 12
  suggestions:
  ...
  - name: random-example-2fpnqfv8
    parameterAssignments:
    - name: lr
      value: "0.021135228357807213"
    - name: num-layers
      value: "4"
    - name: optimizer
      value: sgd
  - name: random-example-xcnrpx4x
    parameterAssignments:
    - name: lr
      value: "0.02414696373094622"
    - name: num-layers
      value: "3"
    - name: optimizer
      value: adam
  - name: random-example-8xmpj4gv
    parameterAssignments:
    - name: lr
      value: "0.02471053882990492"
    - name: num-layers
      value: "4"
    - name: optimizer
      value: sgd
  ...
```

### Trial

For each set of hyperparameters, Katib internally generates a `Trial` CR
with the hyperparameters key-value pairs, `Worker Job` run specification with
parameters instantiated and some other fields like below. The `Trial` CR
is used for internal logic control and end user can even ignore it.

```yaml
$ kubectl get trial -n kubeflow

NAME                      TYPE        STATUS   AGE
random-example-2fpnqfv8   Succeeded   True     10m
random-example-2s9vfb9s   Succeeded   True     8m15s
random-example-5hxm45x4   Succeeded   True     17m
random-example-8xmpj4gv   Succeeded   True     8m44s
random-example-b6gnl4cs   Succeeded   True     12m
random-example-ftm2v84q   Succeeded   True     17m
random-example-gnz5nccf   Succeeded   True     14m
random-example-p74tn9gk   Succeeded   True     11m
random-example-q6jrlshx   Succeeded   True     17m
random-example-tkk46c4x   Succeeded   True     12m
random-example-w5qgblgk   Succeeded   True     12m
random-example-xcnrpx4x   Succeeded   True     10m

$ kubectl get trial random-example-2fpnqfv8 -o yaml -n kubeflow

apiVersion: kubeflow.org/v1beta1
kind: Trial
metadata:
  ...
  name: random-example-2fpnqfv8
  namespace: kubeflow
  ownerReferences:
  - apiVersion: kubeflow.org/v1beta1
    blockOwnerDeletion: true
    controller: true
    kind: Experiment
    name: random-example
    uid: 302e79ae-8659-4679-9e2d-461209619883
  ...
spec:
  failureCondition: status.conditions.#(type=="Failed")#|#(status=="True")#
  metricsCollector:
    collector:
      kind: StdOut
  objective:
    additionalMetricNames:
    - Train-accuracy
    goal: 0.99
    metricStrategies:
    - name: Validation-accuracy
      value: max
    - name: Train-accuracy
      value: max
    objectiveMetricName: Validation-accuracy
    type: maximize
  parameterAssignments:
  - name: lr
    value: "0.021135228357807213"
  - name: num-layers
    value: "4"
  - name: optimizer
    value: sgd
  primaryContainerName: training-container
  runSpec:
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: random-example-2fpnqfv8
      namespace: kubeflow
    spec:
      template:
        spec:
          containers:
          - command:
            - python3
            - /opt/mxnet-mnist/mnist.py
            - --batch-size=64
            - --lr=0.021135228357807213
            - --num-layers=4
            - --optimizer=sgd
            image: docker.io/kubeflowkatib/mxnet-mnist:v1beta1-e294a90
            name: training-container
          restartPolicy: Never
  successCondition: status.conditions.#(type=="Complete")#|#(status=="True")#
status:
  completionTime: "2020-11-16T20:09:33Z"
  conditions:
  - lastTransitionTime: "2020-11-16T20:07:48Z"
    lastUpdateTime: "2020-11-16T20:07:48Z"
    message: Trial is created
    reason: TrialCreated
    status: "True"
    type: Created
  - lastTransitionTime: "2020-11-16T20:09:33Z"
    lastUpdateTime: "2020-11-16T20:09:33Z"
    message: Trial is running
    reason: TrialRunning
    status: "False"
    type: Running
  - lastTransitionTime: "2020-11-16T20:09:33Z"
    lastUpdateTime: "2020-11-16T20:09:33Z"
    message: Trial has succeeded
    reason: TrialSucceeded
    status: "True"
    type: Succeeded
  observation:
    metrics:
    - latest: "0.977309"
      max: "0.978105"
      min: "0.958002"
      name: Validation-accuracy
    - latest: "0.993820"
      max: "0.993820"
      min: "0.916611"
      name: Train-accuracy
  startTime: "2020-11-16T20:07:48Z"
```

## What happens after an `Experiment` CR is created

When user creates an `Experiment` CR, Katib `Experiment` controller,
`Suggestion` controller and `Trial` controller is working together to achieve
hyperparameters tuning for user's Machine learning model. The Experiment
workflow looks as follows:

<center>
<img width="100%" alt="image" src="images/katib-workflow.png">
</center>

1. The `Experiment` CR is submitted to the Kubernetes API server. Katib
   `Experiment` mutating and validating webhook is called to set the default
   values for the `Experiment` CR and validate the CR separately.

1. The `Experiment` controller creates the `Suggestion` CR.

1. The `Suggestion` controller creates the algorithm deployment and service
   based on the new `Suggestion` CR.

1. When the `Suggestion` controller verifies that the algorithm service is
   ready, it calls the service to generate
   `spec.request - len(status.suggestions)` sets of hyperparamters and append
   them into `status.suggestions`.

1. The `Experiment` controller finds that `Suggestion` CR had been updated and
   generates each `Trial` for the each new hyperparamters set.

1. The `Trial` controller generates `Worker Job` based on the `runSpec`
   from the `Trial` CR with the new hyperparamters set.

1. The related job controller
   (Kubernetes batch Job, Kubeflow TFJob, Tekton Pipeline, etc.) generates
   Kubernetes Pods.

1. Katib Pod mutating webhook is called to inject the metrics collector sidecar
   container to the candidate Pods.

1. During the ML model container runs, the metrics collector container
   collects metrics from the injected pod and persists metrics to the Katib
   DB backend.

1. When the ML model training ends, the `Trial` controller updates status
   of the corresponding `Trial` CR.

1. When the `Trial` CR goes to end, the `Experiment` controller increases
   `request` field of the corresponding `Suggestion` CR if it is needed,
   then everything goes to `step 4` again.
   Of course, if the `Trial` CRs meet one of `end` condition
   (exceeds `maxTrialCount`, `maxFailedTrialCount` or `goal`),
   the `Experiment` controller takes everything done.
