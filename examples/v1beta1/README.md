# Simple Minikube Demo

You can deploy katib components and try a simple mnist demo on your laptop!

## Requirement

- [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)
- kubectl

## Deploy katib

Start Katib on Minikube with [deploy.sh](./MinikubeDemo/deploy.sh).
A Minikube cluster and Katib components will be deployed! You can check them with `kubectl -n kubeflow get pods`.

Then, start port-forward for katib UI `8080 -> UI`.

kubectl v1.10~:

```
$ kubectl -n kubeflow port-forward svc/katib-ui 8080:80
```

kubectl ~v1.9:

```
& kubectl -n kubeflow port-forward $(kubectl -n kubeflow get pod -o=name | grep katib-ui | sed -e "s@pods\/@@") 8080:80
```

## Create Experiment

#### Random Suggestion Demo

```
$ kubectl apply -f random-example.yaml
```

#### Grid Suggestion Demo

```
$ kubectl apply -f grid-example.yaml
```

#### Bayesian Optimization Suggestion Demo

```
$ kubectl apply -f bayesianoptimization-example.yaml
```

#### Hyperband Suggestion Demo

```
$ kubectl apply -f hyperband-example.yaml
```

#### Run trial evaluation job by [PyTorchJob](https://github.com/kubeflow/pytorch-operator)

```
$ kubectl apply -f pytorchjob-example.yaml
```

#### Run trial evaluation job by [TFJob](https://github.com/kubeflow/tf-operator)

```
$ kubectl apply -f tfjob-example.yaml
```

## Monitor Experiment

#### CLI

You can submit a new Experiment or check your Experiment results with `kubectl` CLI.  
List experiments:

```
# kubectl get experiment -n kubeflow

NAME             STATUS      AGE
random-example   Succeeded   3h
```

Check experiment result:

```yaml
$ kubectl get experiment random-experiment -n kubeflow -o yaml

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
  completionTime: "2020-07-15T15:32:56Z"
  conditions:
  - lastTransitionTime: "2020-07-15T15:23:58Z"
    lastUpdateTime: "2020-07-15T15:23:58Z"
    message: Experiment is created
    reason: ExperimentCreated
    status: "True"
    type: Created
  - lastTransitionTime: "2020-07-15T15:32:56Z"
    lastUpdateTime: "2020-07-15T15:32:56Z"
    message: Experiment is running
    reason: ExperimentRunning
    status: "False"
    type: Running
  - lastTransitionTime: "2020-07-15T15:32:56Z"
    lastUpdateTime: "2020-07-15T15:32:56Z"
    message: Experiment has succeeded because max trial count has reached
    reason: ExperimentMaxTrialsReached
    status: "True"
    type: Succeeded
  currentOptimalTrial:
    bestTrialName: random-example-tvxz667x
    observation:
      metrics:
      - latest: "0.975816"
        max: "0.978901"
        min: "0.955812"
        name: Validation-accuracy
      - latest: "0.993970"
        max: "0.993970"
        min: "0.913713"
        name: Train-accuracy
    parameterAssignments:
    - name: lr
      value: "0.021031758718972005"
    - name: num-layers
      value: "2"
    - name: optimizer
      value: sgd
  startTime: "2020-07-15T15:23:58Z"
  succeededTrialList:
  - random-example-58tbx6xc
  - random-example-5nkb2gz2
  - random-example-88bdbkzr
  - random-example-9tgjl9nt
  - random-example-dqzjb2r9
  - random-example-gjfdgxxn
  - random-example-nhrx8tb8
  - random-example-nkv76z8z
  - random-example-pcnmzl76
  - random-example-spmk57dw
  - random-example-tvxz667x
  - random-example-xpw8wnjc
  trials: 12
  trialsSucceeded: 12
```

List trials:

```
# kubectl get trials -n kubeflow

NAME                      TYPE        STATUS   AGE
random-example-58tbx6xc   Succeeded   True     48m
random-example-5nkb2gz2   Succeeded   True     54m
random-example-88bdbkzr   Succeeded   True     53m
random-example-9tgjl9nt   Succeeded   True     50m
random-example-dqzjb2r9   Succeeded   True     52m
random-example-gjfdgxxn   Succeeded   True     53m
random-example-nhrx8tb8   Succeeded   True     49m
random-example-nkv76z8z   Succeeded   True     51m
random-example-pcnmzl76   Succeeded   True     54m
random-example-spmk57dw   Succeeded   True     48m
random-example-tvxz667x   Succeeded   True     49m
random-example-xpw8wnjc   Succeeded   True     54m
```

Check trial details:

```yaml
$ kubectl get trials random-example-58tbx6xc -o yaml -n kubeflow

apiVersion: kubeflow.org/v1beta1
kind: Trial
metadata:
  ...
  name: random-example-58tbx6xc
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
    value: "0.011911183432583596"
  - name: num-layers
    value: "3"
  - name: optimizer
    value: ftrl
  runSpec:
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: random-example-58tbx6xc
      namespace: kubeflow
    spec:
      template:
        spec:
          containers:
          - command:
            - python3
            - /opt/mxnet-mnist/mnist.py
            - --batch-size=64
            - --lr=0.011911183432583596
            - --num-layers=3
            - --optimizer=ftrl
            image: docker.io/kubeflowkatib/mxnet-mnist
            name: training-container
          restartPolicy: Never
status:
  completionTime: "2020-07-15T15:32:12Z"
  conditions:
  - lastTransitionTime: "2020-07-15T15:30:42Z"
    lastUpdateTime: "2020-07-15T15:30:42Z"
    message: Trial is created
    reason: TrialCreated
    status: "True"
    type: Created
  - lastTransitionTime: "2020-07-15T15:32:12Z"
    lastUpdateTime: "2020-07-15T15:32:12Z"
    message: Trial is running
    reason: TrialRunning
    status: "False"
    type: Running
  - lastTransitionTime: "2020-07-15T15:32:12Z"
    lastUpdateTime: "2020-07-15T15:32:12Z"
    message: Trial has succeeded
    reason: TrialSucceeded
    status: "True"
    type: Succeeded
  observation:
    metrics:
    - latest: "0.113854"
      max: "0.113854"
      min: "0.113854"
      name: Validation-accuracy
    - latest: "0.112390"
      max: "0.112423"
      min: "0.111907"
      name: Train-accuracy
  startTime: "2020-07-15T15:30:42Z"
```

#### UI

You can submit a new Experiment or check your Experiment results with Web UI.
Access to `http://127.0.0.1:8080/katib`.

## Clean

Clean up with [destroy.sh](./MinikubeDemo/destroy.sh) script.
It will stop port-forward process and delete minikube cluster.

# List of current Katib training container images

- Mxnet mnist example with collecting metrics time, [source](https://github.com/kubeflow/katib/blob/master/examples/v1beta1/mxnet-mnist/mnist.py).

```
docker.io/kubeflowkatib/mxnet-mnist
```

- Pytorch mnist example with saving metrics to the file or print them to the StdOut,
  [source](https://github.com/kubeflow/katib/blob/master/examples/v1beta1/pytorch-mnist/mnist.py).

```
docker.io/kubeflowkatib/pytorch-mnist
```

- Keras cifar10 CNN example for ENAS with gpu support, [source](https://github.com/kubeflow/katib/blob/master/examples/v1beta1/nas/enas-cnn-cifar10/Dockerfile.gpu).

```
docker.io/kubeflowkatib/enas-cnn-cifar10-gpu
```

- Keras cifar10 CNN example for ENAS with cpu support, [source](https://github.com/kubeflow/katib/blob/master/examples/v1beta1/nas/enas-cnn-cifar10/Dockerfile.cpu).

```
docker.io/kubeflowkatib/enas-cnn-cifar10-cpu
```

- Pytorch cifar10 CNN example for DARTS, [source](https://github.com/kubeflow/katib/blob/master/examples/v1beta1/nas/darts-cnn-cifar10/Dockerfile)

```
docker.io/kubeflowkatib/darts-cnn-cifar10
```

- TF operator mnist example with writing summary data,
  [source](https://github.com/kubeflow/tf-operator/blob/master/examples/v1/mnist_with_summaries/mnist_with_summaries.py).

```
gcr.io/kubeflow-ci/tf-mnist-with-summaries
```

- FPGA XGBoost Parameter Tuning example, [source](https://github.com/inaccel/jupyter/blob/master/lab/dot/XGBoost/parameter-tuning.py).

```
docker.io/inaccel/jupyter:lab
```

- MPI operator horovod mnist example, [source](https://github.com/kubeflow/mpi-operator/tree/master/examples/horovod).

```
docker.io/kubeflow/mpi-horovod-mnist
```
