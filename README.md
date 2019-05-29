<h1 align="center">
    <img src="./docs/images/Katib_Logo.png" alt="logo" width="200">
  <br>
</h1>

[![Build Status](https://travis-ci.org/kubeflow/katib.svg?branch=master)](https://travis-ci.org/kubeflow/katib)
[![Coverage Status](https://coveralls.io/repos/github/kubeflow/katib/badge.svg?branch=master)](https://coveralls.io/github/kubeflow/katib?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/katib)](https://goreportcard.com/report/github.com/kubeflow/katib)

Katib is a Kubernetes Native System for [Hyperparameter Tuning][1] and [Neural Architecture Search][2].
The system is inspired by [Google vizier][3] and supports multiple ML/DL frameworks (e.g. TensorFlow, MXNet, and PyTorch).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Name](#name)
- [Concepts in Katib](#concepts-in-katib)
  - [Experiment](#experiment)
  - [Trial](#trial)
  - [Job](#job)
  - [Suggestion](#suggestion)
- [Components in Katib](#components-in-katib)
  - [v1alpha1](#v1alpha1)
  - [v1alpha2](#v1alpha2)
- [Getting Started](#getting-started)
- [Web UI](#web-ui)
- [API Documentation](#api-documentation)
- [Quickstart to run tfjob and pytorch operator jobs in Katib](#quickstart-to-run-tfjob-and-pytorch-operator-jobs-in-katib)
  - [TFjob operator](#tfjob-operator)
  - [Pytorch operator](#pytorch-operator)
  - [Katib](#katib)
  - [Running examples](#running-examples)
  - [Cleanups](#cleanups)
- [CONTRIBUTING](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Name

Katib stands for `secretary` in Arabic. As `Vizier` stands for a high official or a prime minister in Arabic, this project Katib is named in the honor of Vizier.

## Concepts in Katib

Katib has the concepts of Experiment, Trial, Job and Suggestion.

### Experiment

`Experiment` represents a single optimization run over a feasible space.
Each `Experiment` contains a configuration describing the feasible space, as well as a set of Trials. 
It is assumed that objective function f(x) does not change in the course of a `Experiment`.

In v1alpha1, `Experiment` is defined as a CRD `StudyJob` in Kubernetes.
In v1alpha2, `Experiment` is defined as a CRD `Experiment`.

### Trial

A `Trial` is a list of parameter values, x, that will lead to a single evaluation of f(x). A Trial can be “Completed”, which means that it has been evaluated and the objective value f(x) has been assigned to it, otherwise it is “Pending”.

In v1alpha1, `Trial` is just a concept inside Katib and not exposed to users.
In v1alpha2, `Trial` is defined as a CRD `Trial` in Kubernetes.

### Job

A `Job` refers to a process responsible for evaluating a Pending `Trial` and calculating its objective value.

The job kind can be [Kubernetes Job](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/), [Kubeflow TFJob](https://www.kubeflow.org/docs/guides/components/tftraining/) or [Kubeflow PyTorchJob](https://www.kubeflow.org/docs/guides/components/pytorch/).
Thus Katib supports multiple frameworks with the help of different job kinds.

### Suggestion

A Suggestion is an algorithm to construct a parameter set according to the `Experiment`. Then `Trial` will be created to evaluate the parameter set.

Currently Katib supports the following exploration algorithms in v1alpha1:

* random search
* grid search
* [hyperband](https://arxiv.org/pdf/1603.06560.pdf)
* [bayesian optimization](https://arxiv.org/pdf/1012.2599.pdf)
* [NAS based on reinforcement learning](https://github.com/kubeflow/katib/tree/master/pkg/suggestion/v1alpha1/NAS_Reinforcement_Learning)
* [NAS based on EnvelopeNets](https://github.com/kubeflow/katib/tree/master/pkg/suggestion/v1alpha1/NAS_Envelopenet)

And Katib supports the following exploration algorithms in v1alpha2:

* random search

## Components in Katib

### v1alpha1

Katib consists of several components as shown below. Each component is running on k8s as a deployment.
Each component communicates with others via GRPC and the API is defined at `pkg/api/v1alpha1/api.proto`.

- vizier: main components.
  - vizier-core: GRPC API server of vizier.
  - vizier-core-rest: REST API server of vizier.
  - vizier-db: Data storage backend of vizier.
  - suggestion: implementation of each exploration algorithm.
    - suggestion-random
    - suggestion-grid
    - suggestion-hyperband
    - suggestion-bayesianoptimization
    - suggestion-nasrl
    - suggestion-nasenvelopenets
- studyjob-controller: Controller for `StudyJob` CRD in Kubernetes.
- modeldb : WebUI
  - modeldb-frontend
  - modeldb-backend
  - modeldb-db

### v1alpha2

Katib consists of several components as shown below. Each component is running on k8s as a deployment.
Each component communicates with others via GRPC and the API is defined at `pkg/api/v1alpha2/api.proto`.

- katib: main components.
  - katib-manager: GRPC API server of katib.
  - katib-manager-rest: REST API server of katib.
  - katib-db: Data storage backend of katib.
  - katib-ui: User interface of katib.
  - suggestion: implementation of each exploration algorithm.
    - suggestion-random
-  katib-controller: Controller for katib CRDs in Kubernetes.
  - experiment-controller: Controller for `Experiment` CRD in Kubernetes.
  - trial-controller: Controller for `Trial` CRD in Kubernetes.

## Getting Started

Please see [README.md](./examples/v1alpha1/README.md) for more details.

## Web UI

Katib provides a Web UI.
You can visualize general trend of Hyper parameter space and each training history. You can use
[random-example](https://github.com/kubeflow/katib/blob/master/examples/v1alpha1/random-example.yaml) or
[other examples](https://github.com/kubeflow/katib/blob/master/examples/v1alpha1) to generate a similar UI.
![katibui](https://user-images.githubusercontent.com/10014831/48778081-a4388b80-ed17-11e8-938b-fc59a5d2e574.gif)

## API Documentation

Please refer to [api.md](./pkg/api/v1alpha1/gen-doc/api.md).

## Quickstart to run tfjob and pytorch operator jobs in Katib

For running tfjob and pytorch operator jobs in Katib, you have to install their packages.

In your Ksonnet app root, run the following

```
export KF_ENV=default
ks env set ${KF_ENV} --namespace=kubeflow
ks registry add kubeflow github.com/kubeflow/kubeflow/tree/master/kubeflow
```

### TFjob operator

For installing tfjob operator, run the following

```
ks pkg install kubeflow/tf-training
ks pkg install kubeflow/common
ks generate tf-job-operator tf-job-operator
ks apply ${KF_ENV} -c tf-job-operator
```

### Pytorch operator
For installing pytorch operator, run the following

```
ks pkg install kubeflow/pytorch-job
ks generate pytorch-operator pytorch-operator
ks apply ${KF_ENV} -c pytorch-operator
```

### Katib

Finally, you can install Katib

```
ks pkg install kubeflow/katib
ks generate katib katib
ks apply ${KF_ENV} -c katib
```

If you want to use Katib not in GKE and you don't have StorageClass for dynamic volume provisioning at your cluster, you have to create persistent volume to bound your persistent volume claim.

This is yaml file for persistent volume

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: katib-mysql
  labels:
    type: local
    app: katib
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /data/katib
```

Create this pv after deploying Katib package

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/manifests/v1alpha1/pv/pv.yaml
```

### Running examples

After deploy everything, you can run examples.

To run tfjob operator example, you have to install volume for it.

If you are using GKE and default StorageClass, you have to create this pvc

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: tfevent-volume
  namespace: kubeflow
  labels:
    type: local
    app: tfjob
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

If you are not using GKE and you don't have StorageClass for dynamic volume provisioning at your cluster, you have to create pvc and pv

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/v1alpha1/tfevent-volume/tfevent-pvc.yaml

kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/v1alpha1/tfevent-volume/tfevent-pv.yaml
```

This is example for tfjob operator

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/v1alpha1/tfjob-example.yaml
```

This is example for pytorch operator

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/v1alpha1/pytorchjob-example.yaml
```

You can check status of StudyJob

```yaml
$ kubectl describe studyjob pytorchjob-example -n kubeflow

Name:         pytorchjob-example
Namespace:    kubeflow
Labels:       controller-tools.k8s.io=1.0
Annotations:  <none>
API Version:  kubeflow.org/v1alpha1
Kind:         StudyJob
Metadata:
  Cluster Name:
  Creation Timestamp:  2019-01-15T18:35:20Z
  Generation:          1
  Resource Version:    1058135
  Self Link:           /apis/kubeflow.org/v1alpha1/namespaces/kubeflow/studyjobs/pytorchjob-example
  UID:                 4fc7ad83-18f4-11e9-a6de-42010a8e0225
Spec:
  Metricsnames:
    accuracy
  Objectivevaluename:  accuracy
  Optimizationgoal:    0.99
  Optimizationtype:    maximize
  Owner:               crd
  Parameterconfigs:
    Feasible:
      Max:          0.05
      Min:          0.01
    Name:           --lr
    Parametertype:  double
    Feasible:
      Max:          0.9
      Min:          0.5
    Name:           --momentum
    Parametertype:  double
  Requestcount:     4
  Study Name:       pytorchjob-example
  Suggestion Spec:
    Request Number:        3
    Suggestion Algorithm:  random
    Suggestion Parameters:
      Name:   SuggestionCount
      Value:  0
  Worker Spec:
    Go Template:
      Raw Template:  apiVersion: "kubeflow.org/v1beta1"
kind: PyTorchJob
metadata:
  name: {{.WorkerID}}
  namespace: kubeflow
spec:
 pytorchReplicaSpecs:
  Master:
    replicas: 1
    restartPolicy: OnFailure
    template:
      spec:
        containers:
          - name: pytorch
            image: gcr.io/kubeflow-ci/pytorch-mnist-with-summary:1.0
            imagePullPolicy: Always
            command:
              - "python"
              - "/opt/pytorch_dist_mnist/dist_mnist_with_summary.py"
              {{- with .HyperParameters}}
              {{- range .}}
              - "{{.Name}}={{.Value}}"
              {{- end}}
              {{- end}}
  Worker:
    replicas: 2
    restartPolicy: OnFailure
    template:
      spec:
        containers:
          - name: pytorch
            image: gcr.io/kubeflow-ci/pytorch-mnist-with-summary:1.0
            imagePullPolicy: Always
            command:
              - "python"
              - "/opt/pytorch_dist_mnist/dist_mnist_with_summary.py"
              {{- with .HyperParameters}}
              {{- range .}}
              - "{{.Name}}={{.Value}}"
              {{- end}}
              {{- end}}
    Retain:  true
Status:
  Conditon:                     Running
  Early Stopping Parameter Id:
  Last Reconcile Time:          2019-01-15T18:35:20Z
  Start Time:                   2019-01-15T18:35:20Z
  Studyid:                      k291b444a0b68631
  Suggestion Count:             1
  Suggestion Parameter Id:      n6f17dd9ff466a2b
  Trials:
    Trialid:  o104235328003ad9
    Workeridlist:
      Completion Time:  <nil>
      Conditon:         Running
      Kind:             PyTorchJob
      Start Time:       2019-01-15T18:35:20Z
      Workerid:         b3b371c89144727f
    Trialid:            ca207b2432231de3
    Workeridlist:
      Completion Time:  <nil>
      Conditon:         Running
      Kind:             PyTorchJob
      Start Time:       2019-01-15T18:35:20Z
      Workerid:         f291b04fb27ece3c
    Trialid:            ddff69212e826432
    Workeridlist:
      Completion Time:  <nil>
      Conditon:         Running
      Kind:             PyTorchJob
      Start Time:       2019-01-15T18:35:20Z
      Workerid:         ncbed67bbcd4a8ed
Events:                 <none>
```

When the spec.Status.Condition becomes ```Completed```, the StudyJob is finished.

You can monitor your results in Katib UI. For accessing to Katib UI, you have to install Ambassador.

In your Ksonnet app root, run the following

```
ks generate ambassador ambassador
ks apply ${KF_ENV} -c ambassador
```

After this, you have to port-forward Ambassador service

```
kubectl port-forward svc/ambassador -n kubeflow 8080:80
```

Finally, you can access to Katib UI using this URL: ```http://localhost:8080/katib/```.

### Cleanups

Delete installed components

```
ks delete ${KF_ENV} -c katib
ks delete ${KF_ENV} -c pytorch-operator
ks delete ${KF_ENV} -c tf-job-operator
```

If you create pv for Katib, delete it

```
kubectl delete -f https://raw.githubusercontent.com/kubeflow/katib/master/manifests/v1alpha1/pv/pv.yaml
```

If you deploy Ambassador, delete it

```
ks delete ${KF_ENV} -c ambassador
```

## CONTRIBUTING

Please feel free to test the system! [developer-guide.md](./docs/developer-guide.md) is a good starting point for developers.

[1]: https://en.wikipedia.org/wiki/Hyperparameter_optimization
[2]: https://en.wikipedia.org/wiki/Neural_architecture_search
[3]: https://static.googleusercontent.com/media/research.google.com/ja//pubs/archive/bcb15507f4b52991a0783013df4222240e942381.pdf
