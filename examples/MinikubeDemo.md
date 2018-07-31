# Simple Minikube Demo
You can deploy katib components and try a simple mnist demo on your laptop!

## Requirement
* VirtualBox
* Minikube
* kubectl

## deploy
Start Katib on Minikube with [deploy.sh](./MinikubeDemo/deploy.sh).
A Minikube cluster and Katib components will be deployed!

You can check them with `kubectl -n katib get pods`.
Don't worry if the `vizier-core` get an error. 
It will be recovered after DB will be prepared.
Wait until all components will be Running status.

Then, start port-forward for katib services `6789 -> manager` and `3000 -> UI`.

kubectl v1.10~
```
$ kubectl -n katib port-forward svc/vizier-core 6789:6789 &
$ kubectl -n katib port-forward svc/modeldb-frontend 3000:3000 &
```

kubectl ~v1.9
```
& kubectl -n katib port-forward $(kubectl -n katib get pod -o=name | grep vizier-core | sed -e "s@pods\/@@") 6789:6789 &
& kubectl -n katib port-forward $(kubectl -n katib get pod -o=name | grep modeldb-frontend | sed -e "s@pods\/@@") 3000:3000 &
```

To start HyperParameter Tuning, you need a katib client.
It will call API of Katib to create study, get suggestions, run trial, and get metrics.
The details of the system flow for the client and katib components is [here](../docs/images/SystemFlow.png).

An example of client is [here](./client-example.go).
The client will read three config files.
* study-config.yml: Define study property and feasible space of parameters.
* suggesiton-config.yml: Define suggesiton parameter for each study and suggestion service. In this file, the config is for grid suggestion service.
* worker-config.yml: Define config for evaluation worker.

## Create Study
### Random Suggestion Demo
You can run rundom suggesiton demo.
```
kubectl -n katib -f random-example.yaml
```
In this demo, 2 random parameters in
* Learning Rate (--lr) - type: double
* Number of NN Layer (--num-layers) - type: int
* optimizer (--optimizer) - type: categorical

Check the study status.
```
$ kubectl -n katib describe studyjobs random-example
Name:         random-example
Namespace:    katib
Labels:       controller-tools.k8s.io=1.0
Annotations:  kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"kubeflow.org/v1alpha1","kind":"StudyJob","metadata":{"annotations":{},"labels":{"controller-tools.k8s.io":"1.0"},"name":"random-example"...
API Version:  kubeflow.org/v1alpha1
Kind:         StudyJob
Metadata:
  Cluster Name:
  Creation Timestamp:  2018-07-31T06:01:33Z
  Generation:          0
  Resource Version:    2105851
  Self Link:           /apis/kubeflow.org/v1alpha1/namespaces/katib/studyjobs/random-example
  UID:                 2d2979d9-9487-11e8-9e34-42010a9200a6
Spec:
  Study Spec:
    Metricsnames:
      accuracy
    Name:                random-example
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
  Suggestion Spec:
    Request Number:         5
    Suggestion Algorithm:   random
    Suggestion Parameters:  <nil>
  Worker Spec:
    Command:
      python
      /mxnet/example/image-classification/train_mnist.py
      --batch-size=64
    Image:  katib/mxnet-mnist-example
    Mountconf:
Status:
  State:    Running
  Studyid:  h52473c20520ce8b
  Trials:
    Trialid:  v90be640b4ac3aa2
    Workeridlist:
      u01344a450f5f2c2
    Trialid:  w5fc56fdfdb39322
    Workeridlist:
      f391d7200e4c631e
    Trialid:  jdd0a57d44090984
    Workeridlist:
      yef518cad8fa5aec
    Trialid:  i39c63a7aef16f13
    Workeridlist:
      f08b781bcfb1505c
    Trialid:  x7a7f8bd9db02db7
    Workeridlist:
      d530be392563c354
Events:  <none>
```
When the Spec.Status.State become `Completed`, the study is completed.
You can look the result on `http://127.0.0.1:3000`.

### Grid Demo
Almost same as random suggestion.

```
kubectl -n katib describe studycontroller grid-example.yml
```

```
$ kubectl -n katib describe studyjobs grid-example
Name:         grid-example
Namespace:    katib
Labels:       controller-tools.k8s.io=1.0
Annotations:  kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"kubeflow.org/v1alpha1","kind":"StudyJob","metadata":{"annotations":{},"labels":{"controller-tools.k8s.io":"1.0"},"name":"grid-example","...
API Version:  kubeflow.org/v1alpha1
Kind:         StudyJob
Metadata:
  Cluster Name:
  Creation Timestamp:  2018-07-26T14:36:13Z
  Generation:          0
  Resource Version:    1564376
  Self Link:           /apis/kubeflow.org/v1alpha1/namespaces/katib/studyjobs/grid-example
  UID:                 3f2a9dc4-90e1-11e8-9e34-42010a9200a6
Spec:
  Study Spec:
    Metricsnames:
      accuracy
    Name:                grid-example
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
  Suggestion Spec:
    Suggestion Algorithm:  grid
    Suggestion Parameters:
      Name:   DefaultGrid
      Value:  1
      Name:   --lr
      Value:  4
  Worker Spec:
    Command:
      python
      /mxnet/example/image-classification/train_mnist.py
      --batch-size=64
    Image:  katib/mxnet-mnist-example
    Mountconf:
Status:
  State:    Completed
  Studyid:  bdf1eeb5cc8c2895
  Trials:
    Trialid:  h082169809c716bf
    Workeridlist:
      b693d6a0ef3ee82b
    Trialid:  x02b8580838c8d87
    Workeridlist:
      rdfeec2aace25331
    Trialid:  icee5229503a7c07
    Workeridlist:
      x599b4a3a6c4617f
    Trialid:  l3b02e51229332fb
    Workeridlist:
      dce5418974b944be
Events:  <none>
```
In this demo, make 4 grids for learning rate (--lr) Min 0.03 and Max 0.07.

### Hyperband Demo
Almost same as random suggestion.

```
kubectl -n katib describe studycontroller hypb-example.yml
```
In this demo, the eta is 3 and the R is 9.

## UI
You can check your Model with Web UI.
Acsess to `http://127.0.0.1:3000/`
The Results will be saved automatically.

## ModelManagement
You can export model data to yaml file with CLI.
```
katib-cli -s {{server-cli}} pull study {{study ID or name}}  -o {{filename}}
```

And you can push your existing models to Katib with CLI.
`mnist-models.yaml` is traind 22 models using random suggestion with this Parameter Config.

```
configs:
    - name: --lr
      parametertype: 1
      feasible:
        max: "0.07"
        min: "0.03"
        list: []
    - name: --lr-factor
      parametertype: 1
      feasible:
        max: "0.05"
        min: "0.005"
        list: []
    - name: --lr-step
      parametertype: 2
      feasible:
        max: "20"
        min: "5"
        list: []
    - name: --optimizer
      parametertype: 4
      feasible:
        max: ""
        min: ""
        list:
        - sgd
        - adam
        - ftrl
```
You can easy to explore the model on ModelDB.

```
katib-cli push md -f mnist-models.yaml
```

## Clean
Clean up with `./destroy.sh` script.
It will stop port-forward process and delete minikube cluster.
