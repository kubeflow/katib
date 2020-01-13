# Simple Minikube Demo
You can deploy katib components and try a simple mnist demo on your laptop!

## Requirement
* [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)
* kubectl

## Deploy katib
Start Katib on Minikube with [deploy.sh](./MinikubeDemo/deploy.sh).
A Minikube cluster and Katib components will be deployed! You can check them with `kubectl -n kubeflow get pods`.

Then, start port-forward for katib UI `8000 -> UI`.

kubectl v1.10~
```
$ kubectl -n kubeflow port-forward svc/katib-ui 8000:80 &
```

kubectl ~v1.9

```
& kubectl -n kubeflow port-forward $(kubectl -n kubeflow get pod -o=name | grep katib-ui | sed -e "s@pods\/@@") 8000:80 &
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
$ kubectl apply -f tfevent-volume/
$ kubectl apply -f tfjob-example.yaml
```
## Monitor Experiment
#### CLI
You can submit a new Experiment or check your Experiment results with `kubectl` CLI.  
List experiments
```
# kubectl get experiment -n kubeflow
NAME                STATUS      AGE
random-experiment   Succeeded   25m
```
Check experiment result
```
# kubectl get experiment random-experiment -n kubeflow -oyaml
apiVersion: kubeflow.org/v1alpha3
kind: Experiment
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"kubeflow.org/v1alpha3","kind":"Experiment","metadata":{"annotations":{},"name":"random-experiment","namespace":"kubeflow"},"spec":{"algorithm":{"algorithmName":"random"},"maxFailedTrialCount":3,"maxTrialCount":12,"objective":{"additionalMetricNames":["accuracy"],"goal":0.99,"objectiveMetricName":"Validation-accuracy","type":"maximize"},"parallelTrialCount":3,"parameters":[{"feasibleSpace":{"max":"0.03","min":"0.01"},"name":"--lr","parameterType":"double"},{"feasibleSpace":{"max":"5","min":"2"},"name":"--num-layers","parameterType":"int"},{"feasibleSpace":{"list":["sgd","adam","ftrl"]},"name":"--optimizer","parameterType":"categorical"}],"trialTemplate":{"goTemplate":{"rawTemplate":"apiVersion: batch/v1\nkind: Job\nmetadata:\n  name: {{.Trial}}\n  namespace: {{.NameSpace}}\nspec:\n  template:\n    spec:\n      containers:\n      - name: {{.Trial}}\n        image: katib/mxnet-mnist-example\n        command:\n        - \"python\"\n        - \"/mxnet/example/image-classification/train_mnist.py\"\n        - \"--batch-size=64\"\n        {{- with .HyperParameters}}\n        {{- range .}}\n        - \"{{.Name}}={{.Value}}\"\n        {{- end}}\n        {{- end}}\n      restartPolicy: Never"}}}}
  creationTimestamp: 2019-07-15T07:37:40Z
  finalizers:
  - clean-data-in-db
  name: random-experiment
  namespace: kubeflow
  resourceVersion: "22147879"
  selfLink: /apis/kubeflow.org/v1alpha3/namespaces/kubeflow/experiments/random-experiment
  uid: 6c8896db-a6d3-11e9-b55b-00163e01b303
spec:
  algorithm:
    algorithmName: random
    algorithmSettings: null
  maxFailedTrialCount: 3
  maxTrialCount: 12
  objective:
    additionalMetricNames:
    - accuracy
    goal: 0.99
    objectiveMetricName: Validation-accuracy
    type: maximize
  parallelTrialCount: 3
  parameters:
  - feasibleSpace:
      max: "0.03"
      min: "0.01"
    name: --lr
    parameterType: double
  - feasibleSpace:
      max: "5"
      min: "2"
    name: --num-layers
    parameterType: int
  - feasibleSpace:
      list:
      - sgd
      - adam
      - ftrl
    name: --optimizer
    parameterType: categorical
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
status:
  completionTime: 2019-07-15T07:45:56Z
  conditions:
  - lastTransitionTime: 2019-07-15T07:37:29Z
    lastUpdateTime: 2019-07-15T07:37:29Z
    message: Experiment is created
    reason: ExperimentCreated
    status: "True"
    type: Created
  - lastTransitionTime: 2019-07-15T07:45:56Z
    lastUpdateTime: 2019-07-15T07:45:56Z
    message: Experiment is running
    reason: ExperimentRunning
    status: "False"
    type: Running
  - lastTransitionTime: 2019-07-15T07:45:56Z
    lastUpdateTime: 2019-07-15T07:45:56Z
    message: Experiment has succeeded because max trial count has reached
    reason: ExperimentSucceeded
    status: "True"
    type: Succeeded
  currentOptimalTrial:
    observation:
      metrics:
      - name: Validation-accuracy
        value: 0.98119
    parameterAssignments:
    - name: --lr
      value: "0.01178778887185771"
    - name: --num-layers
      value: "4"
    - name: --optimizer
      value: sgd
  startTime: 2019-07-15T07:37:29Z
  trials: 12
  trialsSucceeded: 12
```
List trials
```
# kubectl get trials -n kubeflow
NAME                         STATUS      AGE
random-experiment-24lgqghm   Succeeded   26m
random-experiment-2vdqlqfm   Succeeded   28m
random-experiment-4xg8n48f   Succeeded   30m
random-experiment-64stflgp   Succeeded   29m
random-experiment-d9jgsm96   Succeeded   29m
random-experiment-pnrqmqdm   Succeeded   27m
random-experiment-qvcdfppz   Succeeded   27m
random-experiment-r49pflgp   Succeeded   30m
random-experiment-r7d7mcbx   Succeeded   29m
random-experiment-rwbf62k5   Succeeded   26m
random-experiment-vs8pmh2m   Succeeded   27m
random-experiment-wmnlq972   Succeeded   30m
```
Check trial detail
```
# kubectl get trials random-experiment-24lgqghm -oyaml -n kubeflow
apiVersion: kubeflow.org/v1alpha3
kind: Trial
metadata:
  creationTimestamp: 2019-07-15T07:41:38Z
  generation: 1
  labels:
    experiment: random-experiment
  name: random-experiment-24lgqghm
  namespace: kubeflow
  ownerReferences:
  - apiVersion: kubeflow.org/v1alpha3
    blockOwnerDeletion: true
    controller: true
    kind: Experiment
    name: random-experiment
    uid: 6c8896db-a6d3-11e9-b55b-00163e01b303
  resourceVersion: "22147830"
  selfLink: /apis/kubeflow.org/v1alpha3/namespaces/kubeflow/trials/random-experiment-24lgqghm
  uid: fad59cb8-a6d3-11e9-b55b-00163e01b303
spec:
  metricsCollectorSpec: |-
    apiVersion: batch/v1beta1
    kind: CronJob
    metadata:
      name: random-experiment-24lgqghm
      namespace: kubeflow
    spec:
      schedule: "*/1 * * * *"
      successfulJobsHistoryLimit: 0
      failedJobsHistoryLimit: 1
      concurrencyPolicy: Forbid
      jobTemplate:
        spec:
          backoffLimit: 0
          template:
            spec:
              serviceAccountName: metrics-collector
              containers:
              - name: random-experiment-24lgqghm
                image: gcr.io/kubeflow-images-public/katib/v1alpha3/metrics-collector
                imagePullPolicy: IfNotPresent
                command: ["./metricscollector"]
                args:
                - "-e"
                - "random-experiment"
                - "-t"
                - "random-experiment-24lgqghm"
                - "-k"
                - "Job"
                - "-n"
                - "kubeflow"
                - "-m"
                - "katib-db-manager.kubeflow:6789"
                - "-mn"
                - "Validation-accuracy;accuracy"
              restartPolicy: Never
  objective:
    additionalMetricNames:
    - accuracy
    goal: 0.99
    objectiveMetricName: Validation-accuracy
    type: maximize
  parameterAssignments:
  - name: --lr
    value: "0.017151855585117313"
  - name: --num-layers
    value: "5"
  - name: --optimizer
    value: adam
  runSpec: |-
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: random-experiment-24lgqghm
      namespace: kubeflow
    spec:
      template:
        spec:
          containers:
          - name: random-experiment-24lgqghm
            image: katib/mxnet-mnist-example
            command:
            - "python"
            - "/mxnet/example/image-classification/train_mnist.py"
            - "--batch-size=64"
            - "--lr=0.017151855585117313"
            - "--num-layers=5"
            - "--optimizer=adam"
          restartPolicy: Never
status:
  completionTime: 2019-07-15T07:45:42Z
  conditions:
  - lastTransitionTime: 2019-07-15T07:41:29Z
    lastUpdateTime: 2019-07-15T07:41:29Z
    message: Trial is created
    reason: TrialCreated
    status: "True"
    type: Created
  - lastTransitionTime: 2019-07-15T07:45:42Z
    lastUpdateTime: 2019-07-15T07:45:42Z
    message: Trial is running
    reason: TrialRunning
    status: "False"
    type: Running
  - lastTransitionTime: 2019-07-15T07:45:42Z
    lastUpdateTime: 2019-07-15T07:45:42Z
    message: Trial has succeeded
    reason: TrialSucceeded
    status: "True"
    type: Succeeded
  observation:
    metrics:
    - name: Validation-accuracy
      value: 0.969347
  startTime: 2019-07-15T07:41:29Z
```
#### UI
You can submit a new Experiment or check your Experiment results with Web UI.
Acsess to `http://127.0.0.1:8000/katib`
## Clean
Clean up with [destroy.sh](./MinikubeDemo/destroy.sh) script.
It will stop port-forward process and delete minikube cluster.

# List of current Katib training container images

- Mxnet mnist example with collecting metrics time, [source](https://github.com/kubeflow/katib/blob/master/examples/v1alpha3/mxnet-mnist/mnist.py).

```
docker.io/kubeflowkatib/mxnet-mnist
```

- Pytorch mnist example with saving metrics to the file, [source](https://github.com/kubeflow/katib/blob/master/examples/v1alpha3/file-metrics-collector/mnist.py).

```
docker.io/kubeflowkatib/pytorch-mnist
```

- Keras cifar10 example for NAS RL with gpu support, [source](https://github.com/kubeflow/katib/blob/master/examples/v1alpha3/NAS-training-containers/RL-cifar10/Dockerfile.cpu).

```
docker.io/kubeflowkatib/nasrl-cifar10-gpu
```

- Keras cifar10 example for NAS RL with cpu support, [source](https://github.com/kubeflow/katib/blob/master/examples/v1alpha3/NAS-training-containers/RL-cifar10/Dockerfile.cpu).

```
docker.io/kubeflowkatib/nasrl-cifar10-cpu
```

- Pytorch operator mnist example, [source](https://github.com/kubeflow/pytorch-operator/blob/master/examples/mnist/mnist.py).

```
gcr.io/kubeflow-ci/pytorch-dist-mnist-test
```

- Tf operator mnist example, [source](https://github.com/kubeflow/tf-operator/blob/master/examples/v1/mnist_with_summaries/mnist_with_summaries.py).

```
gcr.io/kubeflow-ci/tf-mnist-with-summaries
```
