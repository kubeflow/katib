# Simple Minikube Demo
You can deploy katib components and try a simple mnist demo on your laptop!

## Requirement
* VirtualBox
* Minikube
* kubectl

## deploy
Start Katib on Minikube with [deploy.sh](./MinikubeDemo/deploy.sh).
A Minikube cluster and Katib components will be deployed!

You can check them with `kubectl -n kubeflow get pods`.
Don't worry if the `vizier-core` get an error. 
It will be recovered after DB will be prepared.
Wait until all components will be Running status.

Then, start port-forward for katib services `6789 -> manager` and `8000 -> UI`.

kubectl v1.10~
```
$ kubectl -n kubeflow port-forward svc/vizier-core 6789:6789 &
$ kubectl -n kubeflow port-forward svc/katib-ui 8000:80 &
```

kubectl ~v1.9

```
& kubectl -n kubeflow port-forward $(kubectl -n kubeflow get pod -o=name | grep vizier-core | sed -e "s@pods\/@@") 6789:6789 &
& kubectl -n kubeflow port-forward $(kubectl -n kubeflow get pod -o=name | grep katib-ui | sed -e "s@pods\/@@") 8000:80 &
```

## Create Study
### Random Suggestion Demo
```
$ kubectl apply -f random-example.yaml
```
Only this command, a study will start, generate hyper-parameters and save the results.
The configurations for the study(hyper-parameter feasible space, optimization parameter, optimization goal, suggestion algorithm, and so on) are defined in `random-example.yaml`,
In this demo, hyper-parameters are embedded as args.
You can embed hyper-parameters in another way(e.g. environment values) by using template.
It is defined in `WorkerSpec.GoTemplate.RawTemplate`.
It is written in [go template](https://golang.org/pkg/text/template/) format.

In this demo, 3 hyper parameters 
* Learning Rate (--lr) - type: double
* Number of NN Layer (--num-layers) - type: int
* optimizer (--optimizer) - type: categorical
are randomly generated.

```
$ kubectl -n kubeflow get studyjob
NAME             AGE
random-example   2m
```

Check the study status.

```
$ kubectl -n kubeflow describe studyjobs random-example
Name:         random-example
Namespace:    kubeflow
Labels:       controller-tools.k8s.io=1.0
Annotations:  <none>
API Version:  kubeflow.org/v1alpha1
Kind:         StudyJob
Metadata:
  Creation Timestamp:  2019-02-14T13:53:11Z
  Finalizers:
    clean-studyjob-data
  Generation:        1
  Resource Version:  5625476
  Self Link:         /apis/kubeflow.org/v1alpha1/namespaces/kubeflow/studyjobs/random-example
  UID:               de365269-305f-11e9-973d-0016ac101a86
Spec:
  Metricsnames:
    accuracy
  Objectivevaluename:  Validation-accuracy
  Optimizationgoal:    0.99
  Optimizationtype:    maximize
  Owner:               crd
  Parameterconfigs:
    Feasible:
      Max:          0.03
      Min:          0.01
    Name:           --lr
    Parametertype:  double
    Feasible:
      Max:          5
      Min:          2
    Name:           --num-layers
    Parametertype:  int
    Feasible:
      List:
        sgd
        adam
        ftrl
    Name:           --optimizer
    Parametertype:  categorical
  Requestcount:     1
  Study Name:       random-example
  Suggestion Spec:
    Request Number:        3
    Suggestion Algorithm:  random
    Suggestion Parameters:
      Name:   SuggestionCount
      Value:  0
  Worker Spec:
    Go Template:
      Raw Template:  apiVersion: batch/v1
kind: Job
metadata:
  name: {{.WorkerID}}
  namespace: kubeflow
spec:
  template:
    spec:
      containers:
      - name: {{.WorkerID}}
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
Status:
  Condition:                Running
  Last Reconcile Time:      2019-02-14T13:53:12Z
  Start Time:               2019-02-14T13:53:11Z
  Studyid:                  q267516663b357c2
  Suggestion Count:         1
  Suggestion Parameter Id:  wa4e0d9f801a5a33
  Trials:
    Trialid:  y9b54306d9d9b4d5
    Workeridlist:
      Completion Time:  <nil>
      Condition:        Running
      Kind:             Job
      Start Time:       2019-02-14T13:53:12Z
      Workerid:         ib2201d45c3df144
    Trialid:            dff87c7ef278a1e4
    Workeridlist:
      Completion Time:  <nil>
      Condition:        Running
      Kind:             Job
      Start Time:       2019-02-14T13:53:12Z
      Workerid:         cc0402c150661f3c
    Trialid:            n594a099a65d3a88
    Workeridlist:
      Completion Time:  <nil>
      Condition:        Running
      Kind:             Job
      Start Time:       2019-02-14T13:53:12Z
      Workerid:         e9eae6139a57892f
Events:                 <none>
```

When the Spec.Status.Condition becomes `Completed`, the study is completed.
You can look the result on `http://127.0.0.1:8000/katib`.

### Use ConfigMap for Worker Template
In Random example, the template for workers is defined in StudyJob manifest.
A ConfigMap is also used for worker template.
Let's use [this](./workerConfigMap.yaml) template.
```
kubectl apply -f workerConfigMap.yaml
```
This template will be shared among the three demos below(Grid, Hyperband, and GPU).

### Grid Demo
Almost same as random suggestion.

In this demo, Katib will make 4 grids for learning rate (--lr) Min 0.03 and Max 0.07.
```
kubectl apply -f grid-example.yaml
```

### Hyperband Demo
In this demo, the eta is 3 and the R is 9.
```
kubectl apply -f hypb-example.yaml
```

## UI
You can check your study results with Web UI.
Acsess to `http://127.0.0.1:8000/katib`
The Results will be saved automatically.

### Using GPU demo
You can set any configuration for your worker pods.
Here, try to set config for GPU.
The manifest of the worker pods are generated from a template.
The templates are defined in [ConfigMap](./workerConfigMap.yaml).
There are two templates: defaultWorkerTemplate.yaml and gpuWorkerTemplate.yaml.
You can add your template for worker.
Then you should specify the template in your studyjob spec.
[This example](/examples/gpu-example.yaml) uses `gpuWorkerTemplate.yaml`.
You can apply it same as other examples.
```
$ kubectl apply -f gpu-example.yaml
$ kubectl -n kubeflow get studyjob

NAME             AGE
gpu-example      1m
random-example   17m

$ kubectl -n kubeflow describe studyjob gpu-example
Name:         gpu-example
Namespace:    kubeflow
Labels:       controller-tools.k8s.io=1.0
Annotations:  <none>
API Version:  kubeflow.org/v1alpha1
Kind:         StudyJob
Metadata:
  Creation Timestamp:  2019-02-14T14:00:15Z
  Finalizers:
    clean-studyjob-data
  Generation:        1
  Resource Version:  5626905
  Self Link:         /apis/kubeflow.org/v1alpha1/namespaces/kubeflow/studyjobs/gpu-example
  UID:               daba7165-3060-11e9-973d-0016ac101a86
Spec:
  Metricsnames:
    accuracy
  Objectivevaluename:  Validation-accuracy
  Optimizationgoal:    0.99
  Optimizationtype:    maximize
  Owner:               crd
  Parameterconfigs:
    Feasible:
      Max:          0.03
      Min:          0.01
    Name:           --lr
    Parametertype:  double
    Feasible:
      Max:          3
      Min:          2
    Name:           --num-layers
    Parametertype:  int
    Feasible:
      List:
        sgd
        adam
        ftrl
    Name:           --optimizer
    Parametertype:  categorical
  Study Name:       gpu-example
  Suggestion Spec:
    Request Number:        3
    Suggestion Algorithm:  random
    Suggestion Parameters:
      Name:   SuggestionCount
      Value:  0
  Worker Spec:
    Go Template:
      Template Path:        gpuWorkerTemplate.yaml
Status:
  Condition:                Running
  Last Reconcile Time:      2019-02-14T14:00:17Z
  Start Time:               2019-02-14T14:00:15Z
  Studyid:                  g3b79d9c0ff8881f
  Suggestion Count:         1
  Suggestion Parameter Id:  z313763f77337c14
  Trials:
    Trialid:  xc63f4f77156df83
    Workeridlist:
      Completion Time:  <nil>
      Condition:        Running
      Kind:             Job
      Start Time:       2019-02-14T14:00:16Z
      Workerid:         ue4468cbb6cb2045
    Trialid:            ee8011ddd3937998
    Workeridlist:
      Completion Time:  <nil>
      Condition:        Running
      Kind:             Job
      Start Time:       2019-02-14T14:00:16Z
      Workerid:         mcff6fcf01c8f2d4
    Trialid:            f54ba014544ef2ad
    Workeridlist:
      Completion Time:  <nil>
      Condition:        Running
      Kind:             Job
      Start Time:       2019-02-14T14:00:16Z
      Workerid:         r5c011c5c9eca2e9
Events:                 <none>
```

Check if the GPU configuration works correctly.

```
$ kubectl -n kubeflow describe pod r5c011c5c9eca2e9-ftmcj
Name:               r5c011c5c9eca2e9-ftmcj
Namespace:          kubeflow
Priority:           0
PriorityClassName:  <none>
Node:               <none>
Labels:             controller-uid=db82ecf9-3060-11e9-973d-0016ac101a86
                    job-name=r5c011c5c9eca2e9
Annotations:        kubernetes.io/psp=ibm-privileged-psp
Status:             Pending
IP:                 
Controlled By:      Job/r5c011c5c9eca2e9
Containers:
  r5c011c5c9eca2e9:
    Image:      katib/mxnet-mnist-example:gpu
    Port:       <none>
    Host Port:  <none>
    Command:
      python
      /mxnet/example/image-classification/train_mnist.py
      --batch-size=64
      --lr=0.0244
      --num-layers=2
      --optimizer=ftrl
    Limits:
      nvidia.com/gpu:  1
    Requests:
      nvidia.com/gpu:  1
    Environment:       <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from default-token-vh4d7 (ro)
Conditions:
  Type           Status
  PodScheduled   False 
Volumes:
  default-token-vh4d7:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  default-token-vh4d7
    Optional:    false
QoS Class:       BestEffort
Node-Selectors:  <none>
Tolerations:     node.kubernetes.io/not-ready:NoExecute for 300s
                 node.kubernetes.io/unreachable:NoExecute for 300s
Events:
  Type     Reason            Age               From               Message
  ----     ------            ----              ----               -------
  Warning  FailedScheduling  3m (x25 over 4m)  default-scheduler  0/3 nodes are available: 3 Insufficient nvidia.com/gpu.
```

## Metrics Collection

### Design of Metrics Collector
![metricscollectordesign](https://user-images.githubusercontent.com/10014831/47256754-e32cb480-d4bf-11e8-98e9-4bbec562ad75.png)

### Default Metrics Collector

The default metrics collector collects metrics from the StdOut of workers.
It is deployed as a cronjob. It will collect and report metrics periodically.
It collects metrics through k8s pod log API.
You should print logs in {metrics name}={value} style.
In the above demo, the objective value name is *Validation-accuracy* and the metrics are [*accuracy*], so your training code should print like below.
```
epoch 1:
batch1 accuracy=0.3
batch2 accuracy=0.5

Validation-accuracy=0.4

epoch 2:
batch1 accuracy=0.7
batch2 accuracy=0.8

Validation-accuracy=0.75
```
The metrics collector will collect all logs of metrics.
The manifest of metrics collector is also generated from template and defined [here](/manifests/studyjobcontroller/metricsControllerConfigMap.yaml).
You can add your template and specify `spec.metricsCollectorSpec.metricsCollectorTemplatePath` in a studyjob manifest.

### TF Event File Metrics Collector

The TF Event file metrics collector will collect metrics from tf.event files.
It is also deployed as a cronjob.
When you use TF Event File Metrics Collector, you need to share files between the metrics collector and the worker by PVC.
There is an example for TF Event file metrics collector.
First, please create PV and PVC to share event file.
```
$ kubectl apply -f tfevent-volume/
```
Then, create a studyjob that uses TF Event file metrics collector.
```
$ kubectl apply -f tf-event_test.yaml
```

It will create a tensorflow worker from whose eventfile metrics are collected.
The code of tensorflow is [the official tutorial for mnist with summary](https://github.com/tensorflow/tensorflow/blob/master/tensorflow/examples/tutorials/mnist/mnist_with_summaries.py).
It will save event file to `/log/train` and `/log/test` directory.
They have same named metrics ('accuracy' and 'cross_entropy').
The accuracy in training and test will be saved in *train/* directory and *test/* directory respectively.
In a studyjob, please add directry name to the name of metrics as a prefix e.g. `train/accuracy`, `test/accuracy`.

## Clean
Clean up with `./destroy.sh` script.
It will stop port-forward process and delete minikube cluster.
