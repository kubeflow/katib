# How Katib v1beta1 tunes hyperparameter automatically in a Kubernetes native way

See the following guides in the Kubeflow documentation:

* [Concepts](https://www.kubeflow.org/docs/components/hyperparameter-tuning/overview/) 
  in Katib, hyperparameter tuning, and neural architecture search.
* [Getting started with Katib](https://kubeflow.org/docs/components/hyperparameter-tuning/hyperparameter/).
* Detailed guide to [configuring and running a Katib 
  experiment](https://kubeflow.org/docs/components/hyperparameter-tuning/experiment/).

## Example and Illustration

After install Katib v1beta1, you can run `kubectl apply -f katib/examples/v1beta1/random-example.yaml` to try the first example of Katib.
Then you can get the new `Experiment` as below. Katib concepts will be introduced based on this example.

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
              image: docker.io/kubeflowkatib/mxnet-mnist
              name: training-container
            restartPolicy: Never
status:
  ...
```
#### Experiment

When you want to tune hyperparameters for your machine learning model before 
training it further, you just need to create an `Experiment` CR like above. To
learn what fields are included in the `Experiment.spec`, see
the detailed guide to [configuring and running a Katib 
experiment](https://kubeflow.org/docs/components/hyperparameter-tuning/experiment/).

#### Trial

For each set of hyperparameters, Katib will internally generate a `Trial` CR with the hyperparameters key-value pairs, job manifest string with parameters instantiated and some other fields like below. `Trial` CR is used for internal logic control, and end user can even ignore it.

```yaml
$ kubectl get trial -n kubeflow

NAME                      TYPE        STATUS   AGE
random-example-58tbx6xc   Succeeded   True     14m
random-example-5nkb2gz2   Succeeded   True     21m
random-example-88bdbkzr   Succeeded   True     20m
random-example-9tgjl9nt   Succeeded   True     17m
random-example-dqzjb2r9   Succeeded   True     19m
random-example-gjfdgxxn   Succeeded   True     20m
random-example-nhrx8tb8   Succeeded   True     15m
random-example-nkv76z8z   Succeeded   True     18m
random-example-pcnmzl76   Succeeded   True     21m
random-example-spmk57dw   Succeeded   True     14m
random-example-tvxz667x   Succeeded   True     16m
random-example-xpw8wnjc   Succeeded   True     21m

$ kubectl get trial random-example-gjfdgxxn -o yaml -n kubeflow

apiVersion: kubeflow.org/v1beta1
kind: Trial
metadata:
  ...
  name: random-example-gjfdgxxn
  namespace: kubeflow
  ownerReferences:
  - apiVersion: kubeflow.org/v1beta1
    blockOwnerDeletion: true
    controller: true
    kind: Experiment
    name: random-example
    uid: 34349cb7-c6af-11ea-90dd-42010a9a0020
  ...
spec:
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
    value: "0.012171302435678337"
  - name: num-layers
    value: "3"
  - name: optimizer
    value: adam
  runSpec:
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: random-example-gjfdgxxn
      namespace: kubeflow
    spec:
      template:
        spec:
          containers:
          - command:
            - python3
            - /opt/mxnet-mnist/mnist.py
            - --batch-size=64
            - --lr=0.012171302435678337
            - --num-layers=3
            - --optimizer=adam
            image: docker.io/kubeflowkatib/mxnet-mnist
            name: training-container
          restartPolicy: Never
status:
  completionTime: "2020-07-15T15:29:00Z"
  conditions:
  - lastTransitionTime: "2020-07-15T15:25:16Z"
    lastUpdateTime: "2020-07-15T15:25:16Z"
    message: Trial is created
    reason: TrialCreated
    status: "True"
    type: Created
  - lastTransitionTime: "2020-07-15T15:29:00Z"
    lastUpdateTime: "2020-07-15T15:29:00Z"
    message: Trial is running
    reason: TrialRunning
    status: "False"
    type: Running
  - lastTransitionTime: "2020-07-15T15:29:00Z"
    lastUpdateTime: "2020-07-15T15:29:00Z"
    message: Trial has succeeded
    reason: TrialSucceeded
    status: "True"
    type: Succeeded
  observation:
    metrics:
    - latest: "0.959594"
      max: "0.960490"
      min: "0.940585"
      name: Validation-accuracy
    - latest: "0.959022"
      max: "0.959188"
      min: "0.921658"
      name: Train-accuracy
  startTime: "2020-07-15T15:25:16Z"
```

#### Suggestion

Katib will internally create a `Suggestion` CR for each `Experiment` CR. `Suggestion` CR includes the hyperparameter algorithm name by `algorithmName` field and how many sets of hyperparameter Katib asks to be generated by `requests` field. The CR also traces all already generated sets of hyperparameter in `status.suggestions`. Same as `Trial`, `Suggestion` CR is used for internal logic control and end user can even ignore it.

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
    uid: 34349cb7-c6af-11ea-90dd-42010a9a0020
  ...
spec:
  algorithmName: random
  requests: 12
status:
  suggestionCount: 12
  suggestions:
  ...
  - name: random-example-gjfdgxxn
    parameterAssignments:
    - name: lr
      value: "0.012171302435678337"
    - name: num-layers
      value: "3"
    - name: optimizer
      value: adam
  - name: random-example-88bdbkzr
    parameterAssignments:
    - name: lr
      value: "0.013408352284328112"
    - name: num-layers
      value: "4"
    - name: optimizer
      value: ftrl
  - name: random-example-dqzjb2r9
    parameterAssignments:
    - name: lr
      value: "0.028873905258692753"
    - name: num-layers
      value: "3"
    - name: optimizer
      value: adam
  ...
```

## What happens after an `Experiment` CR created

When a user created an `Experiment` CR, Katib controllers including experiment controller, trial controller and suggestion controller will work together to achieve hyperparameters tuning for user Machine learning model.
<center>
<img width="100%" alt="image" src="images/katib-workflow.png">
</center>

1. A `Experiment` CR is submitted to Kubernetes API server, Katib experiment mutating and validating webhook will be called to set default value for the `Experiment` CR and validate the CR separately.
2. Experiment controller creates a `Suggestion` CR.
3. Suggestion controller creates the algorithm deployment and service based on the new `Suggestion` CR.
4. When Suggestion controller verifies that the algorithm service is ready, it calls the service to generate `spec.request - len(status.suggestions)` sets of hyperparamters and append them into `status.suggestions`
5. Experiment controller finds that `Suggestion` CR had been updated, then generate each `Trial` for each new hyperparamters set. 
6. Trial controller generates job based on `trialSpec` manifest with the new hyperparamters set.
7. Related job controller (Kubernetes batch Job, Kubeflow PyTorchJob or Kubeflow TFJob) generates Pods.
8. Katib Pod mutating webhook is called to inject metrics collector sidecar container to the candidate Pod.
9. During the ML model container runs, metrics collector container in the same Pod tries to collect metrics from it and persists them into Katib DB backend.
10. When the ML model Job ends, Trial controller will update status of the corresponding `Trial` CR.
11. When a `Trial` CR goes to end, Experiment controller will increase `request` field of corresponding 
`Suggestion` CR if it is needed, then everything goes to `step 4` again. Of course, if `Trial` CRs meet one of `end` condition (exceeds `maxTrialCount`, `maxFailedTrialCount` or `goal`), Experiment controller will take everything done.
