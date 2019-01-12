# Katib

[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/katib)](https://goreportcard.com/report/github.com/kubeflow/katib)

<img src="./img/Katib_Logo.png" width="320px">

Hyperparameter Tuning on Kubernetes.
This project is inspired by [Google vizier](https://static.googleusercontent.com/media/research.google.com/ja//pubs/archive/bcb15507f4b52991a0783013df4222240e942381.pdf). Katib is a scalable and flexible hyperparameter tuning framework and is tightly integrated with kubernetes. Also it does not depend on a specific Deep Learning framework (e.g. TensorFlow, MXNet, and PyTorch).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Name](#name)
- [Concepts in Google Vizier](#concepts-in-google-vizier)
  - [Study](#study)
  - [Trial](#trial)
  - [Suggestion](#suggestion)
- [Components in Katib](#components-in-katib)
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
- [TODOs](#todos)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Name

Katib stands for `secretary` in Arabic. As `Vizier` stands for a high official or a prime minister in Arabic, this project Katib is named in the honor of Vizier.

## Concepts in Google Vizier

As in Google Vizier, Katib also has the concepts of Study, Trial and Suggestion.

### Study

Represents a single optimization run over a feasible space. Each Study contains a configuration describing the feasible space, as well as a set of Trials. It is assumed that objective function f(x) does not change in the course of a Study.

### Trial

A Trial is a list of parameter values, x, that will lead to a single evaluation of f(x). A Trial can be “Completed”, which means that it has been evaluated and the objective value f(x) has been assigned to it, otherwise it is “Pending”.
One trial corresponds to one job, and the job kind can be [k8s Job](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/), [TFJob](https://www.kubeflow.org/docs/guides/components/tftraining/) or [PyTorchJob](https://www.kubeflow.org/docs/guides/components/pytorch/), which depends on the Study's worker kind.

### Suggestion

A Suggestion is an algorithm to construct a parameter set. Currently Katib supports the following exploration algorithms:

* random
* grid
* [hyperband](https://arxiv.org/pdf/1603.06560.pdf)
* [bayesian optimization](https://arxiv.org/pdf/1012.2599.pdf)

## Components in Katib

Katib consists of several components as shown below. Each component is running on k8s as a deployment.
Each component communicates with others via GRPC and the API is defined at `pkg/api/api.proto`.

- vizier: main components.
    - vizier-core : API server of vizier.
    - vizier-db
- suggestion : implementation of each exploration algorithm.
    - vizier-suggestion-random
    - vizier-suggestion-grid
    - vizier-suggestion-hyperband
    - vizier-suggestion-bayesianoptimization
- modeldb : WebUI
    - modeldb-frontend
    - modeldb-backend
    - modeldb-db

## Getting Started

Please see [MinikubeDemo.md](./examples/MinikubeDemo.md) for more details.

## Web UI

Katib provides a Web UI.
You can visualize general trend of Hyper parameter space and each training history.
![katibui](https://user-images.githubusercontent.com/10014831/48778081-a4388b80-ed17-11e8-938b-fc59a5d2e574.gif)

## API Documentation

Please refer to [api.md](./pkg/api/gen-doc/api.md).

## Quickstart to run tfjob and pytorch operator jobs in Katib

For running tfjob and pytorch operator jobs in Katib you have to install their packages.

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

After this you have to install volume for tfjob operator.

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
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/tfevent-volume/tfevent-pvc.yaml

kubectl create -f https://raw.githubusercontent.com/andreyvelich/katib/example-doc-pytorch-tfjob-313/examples/tfevent-volume/tfevent-pv.yaml
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
kubectl create -f katib-mysql-pv.yaml
```

### Running examples

After deploy everything you can run examples.

This is example for tfjob operator

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/pytorchjob-example.yaml
```

This is example for pytorch operator
```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/tfjob-example.yaml
```


### Cleanups

Delete installed components

```
ks delete ${KF_ENV} -c katib
ks delete ${KF_ENV} -c pytorch-operator
ks delete ${KF_ENV} -c tf-job-operator
```

If you create pv for Katib delete it

```
kubectl delete -f katib-mysql-pv.yaml
```

## CONTRIBUTING

Please feel free to test the system! [developer-guide.md](./docs/developer-guide.md) is a good starting point for developers.

## TODOs

* Integrate KubeFlow (TensorFlow, Caffe2 and PyTorch operators)
* Support Early Stopping
* Enrich the GUI
