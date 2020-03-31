<h1 align="center">
    <img src="./docs/images/Katib_Logo.png" alt="logo" width="200">
  <br>
</h1>

[![Build Status](https://travis-ci.org/kubeflow/katib.svg?branch=master)](https://travis-ci.org/kubeflow/katib)
[![Coverage Status](https://coveralls.io/repos/github/kubeflow/katib/badge.svg?branch=master)](https://coveralls.io/github/kubeflow/katib?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/katib)](https://goreportcard.com/report/github.com/kubeflow/katib)

Katib is a Kubernetes-based system for [Hyperparameter Tuning][1] and [Neural Architecture Search][2]. Katib supports a number of ML frameworks, including TensorFlow, Apache MXNet, PyTorch, XGBoost, and others.

Table of Contents
=================

  * [Getting Started](#getting-started)
  * [Name](#name)
  * [Concepts in Katib](#concepts-in-katib)
      * [Experiment](#experiment)
      * [Suggestion](#suggestion)
      * [Trial](#trial)
      * [Worker Job](#worker-job)
  * [Components in Katib](#components-in-katib)
  * [Web UI](#web-ui)
  * [API documentation](#api-documentation)
  * [Installation](#installation)
      * [TF operator](#tf-operator)
      * [Pytorch operator](#pytorch-operator)
      * [Katib](#katib)
      * [Running examples](#running-examples)
      * [Cleanups](#cleanups)
  * [Quick Start](#quick-start)
  * [Who are using Katib?](#who-are-using-katib)
  * [CONTRIBUTING](#contributing)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)

## Getting Started

See the [getting-started 
guide](https://www.kubeflow.org/docs/components/hyperparameter-tuning/hyperparameter/)
on the Kubeflow website.

## Name

Katib stands for `secretary` in Arabic.

## Concepts in Katib

For a detailed description of the concepts in Katib, hyperparameter tuning, and
neural architecture search, see the [Kubeflow 
documentation](https://www.kubeflow.org/docs/components/hyperparameter-tuning/overview/).

Katib has the concepts of Experiment, Trial, Job and Suggestion.

### Experiment

`Experiment` represents a single optimization run over a feasible space.
Each `Experiment` contains a configuration 
1. Objective: What we are trying to optimize
2. Search Space: Constraints for configurations describing the feasible space.
3. Search Algorithm: How to find the optimal configurations

`Experiment` is defined as a CRD. See the detailed guide to [configuring and running a Katib 
experiment](https://kubeflow.org/docs/components/hyperparameter-tuning/experiment/)
in the Kubeflow docs.

### Suggestion

A Suggestion is a proposed solution to the optimization problem which is one set of hyperparameter values or a list of parameter assignments. Then a `Trial` will be created to evaluate the parameter assignments.

`Suggestion` is defined as a CRD

### Trial

A `Trial` is one iteration of the optimization process, which is one `worker job` instance with a list of parameter assignments(corresponding to a suggestion).

`Trial` is defined as a CRD

### Worker Job 

A `Worker Job` refers to a process responsible for evaluating a `Trial` and calculating its objective value. 

The worker kind can be [Kubernetes Job](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/) which is a non distributed execution, [Kubeflow TFJob](https://www.kubeflow.org/docs/guides/components/tftraining/) or [Kubeflow PyTorchJob](https://www.kubeflow.org/docs/guides/components/pytorch/) which are distributed executions.
Thus, Katib supports multiple frameworks with the help of different job kinds. 

Currently Katib supports the following exploration algorithms:

#### Hyperparameter Tuning

* [Random Search](https://en.wikipedia.org/wiki/Hyperparameter_optimization#Random_search)
* [Tree of Parzen Estimators (TPE)](https://papers.nips.cc/paper/4443-algorithms-for-hyper-parameter-optimization.pdf)
* [Grid Search](https://en.wikipedia.org/wiki/Hyperparameter_optimization#Grid_search)
* [Hyperband](https://arxiv.org/pdf/1603.06560.pdf)
* [Bayesian Optimization](https://arxiv.org/pdf/1012.2599.pdf)
* [CMA Evolution Strategy](https://arxiv.org/abs/1604.00772)

#### Neural Architecture Search

* [Reinforcement Learning](https://github.com/kubeflow/katib/tree/master/pkg/suggestion/v1alpha3/NAS_Reinforcement_Learning)


## Components in Katib

Katib consists of several components as shown below. Each component is running on k8s as a deployment.
Each component communicates with others via GRPC and the API is defined at `pkg/apis/manager/v1alpha3/api.proto`.

- Katib main components:
  - katib-db-manager: GRPC API server of Katib which is the DB Interface.
  - katib-mysql: Data storage backend of Katib using mysql.
  - katib-ui: User interface of Katib.
  - katib-controller: Controller for Katib CRDs in Kubernetes.

## Web UI

Katib provides a Web UI.
You can visualize general trend of Hyper parameter space and each training history. You can use
[random-example](https://github.com/kubeflow/katib/blob/master/examples/v1alpha3/random-example.yaml) or
[other examples](https://github.com/kubeflow/katib/blob/master/examples/v1alpha3) to generate a similar UI.
![katibui](./docs/images/katib-ui.png)

## API documentation

See the [Katib API reference docs](https://www.kubeflow.org/docs/reference/katib/).

## Installation

For standard installation of Katib with support for all job operators, 
install Kubeflow. See the documentation:

* [Kubeflow installation 
guide](https://www.kubeflow.org/docs/started/getting-started/)
* [Kubeflow hyperparameter tuning 
guides](https://www.kubeflow.org/docs/components/hyperparameter-tuning/).

If you install Katib with other Kubeflow components, you can't submit Katib jobs in Kubeflow namespace.

Alternatively, if you want to install Katib manually, follow these steps:

```
git clone git@github.com:kubeflow/manifests.git
Set `MANIFESTS_DIR` to the cloned folder.

```

### TF operator

For installing tfjob operator, run the following

```
cd "${MANIFESTS_DIR}/tf-training/tf-job-crds/base"
kustomize build . | kubectl apply -f -
cd "${MANIFESTS_DIR}/tf-training/tf-job-operator/base"
kustomize build . | kubectl apply -n kubeflow -f -

```

### Pytorch operator
For installing pytorch operator, run the following

```
cd "${MANIFESTS_DIR}/pytorch-job/pytorch-job-crds/base"
kustomize build . | kubectl apply -f -
cd "${MANIFESTS_DIR}/pytorch-job/pytorch-operator/base/"
kustomize build . | kubectl apply -n kubeflow -f -
```

### Katib

Finally, you can install Katib

```
cd "${MANIFESTS_DIR}/katib/katib-crds/base"
kustomize build . | kubectl apply -f -
cd "${MANIFESTS_DIR}/katib/katib-controller/base"
kustomize build . | kubectl apply -f -

```

If you want to use Katib in a cluster that doesn't have a StorageClass for dynamic volume provisioning at your cluster, you have to create persistent volume manually to bound your persistent volume claim.

This is sample yaml file for creating a persistent volume

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: katib-mysql
  labels:
    type: local
    app: katib
spec:
  storageClassName: katib
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /tmp/katib
```

Create this pv after deploying Katib package

### Running examples

After deploy everything, you can run examples to verify the installation.

This is example for tfjob operator

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/v1alpha3/tfjob-example.yaml
```

This is example for pytorch operator

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/v1alpha3/pytorchjob-example.yaml
```

You can check status of experiment 

```yaml
$ kubectl describe experiment tfjob-example -n kubeflow


Name:         tfjob-example
Namespace:    kubeflow
Labels:       <none>
Annotations:  <none>
API Version:  kubeflow.org/v1alpha3
Kind:         Experiment
Metadata:
  Creation Timestamp:  2019-10-06T12:25:44Z
  Generation:          1
  Resource Version:    2110410
  Self Link:           /apis/kubeflow.org/v1alpha3/namespaces/kubeflow/experiments/tfjob-example
  UID:                 6b2bef2d-e834-11e9-93ee-42010aa00075
Spec:
  Algorithm:
    Algorithm Name:        random
  Max Failed Trial Count:  3
  Max Trial Count:         12
  Metrics Collector Spec:
    Collector:
      Kind:  TensorFlowEvent
    Source:
      File System Path:
        Kind:  Directory
        Path:  /train
  Objective:
    Goal:                   0.99
    Objective Metric Name:  accuracy_1
    Type:                   maximize
  Parallel Trial Count:     3
  Parameters:
    Feasible Space:
      Max:           0.05
      Min:           0.01
    Name:            --learning_rate
    Parameter Type:  double
    Feasible Space:
      Max:           200
      Min:           100
    Name:            --batch_size
    Parameter Type:  int
  Trial Template:
    Go Template:
      Raw Template:  apiVersion: "kubeflow.org/v1"
kind: TFJob
metadata:
  name: {{.Trial}}
  namespace: {{.NameSpace}}
spec:
 tfReplicaSpecs:
  Worker:
    replicas: 1 
    restartPolicy: OnFailure
    template:
      spec:
        containers:
          - name: tensorflow 
            image: gcr.io/kubeflow-ci/tf-mnist-with-summaries:1.0
            imagePullPolicy: Always
            command:
              - "python"
              - "/var/tf_mnist/mnist_with_summaries.py"
              - "--log_dir=/train/metrics"
              {{- with .HyperParameters}}
              {{- range .}}
              - "{{.Name}}={{.Value}}"
              {{- end}}
              {{- end}}
Status:
  Completion Time:  2019-10-06T12:28:50Z
  Conditions:
    Last Transition Time:  2019-10-06T12:25:44Z
    Last Update Time:      2019-10-06T12:25:44Z
    Message:               Experiment is created
    Reason:                ExperimentCreated
    Status:                True
    Type:                  Created
    Last Transition Time:  2019-10-06T12:28:50Z
    Last Update Time:      2019-10-06T12:28:50Z
    Message:               Experiment is running
    Reason:                ExperimentRunning
    Status:                False
    Type:                  Running
    Last Transition Time:  2019-10-06T12:28:50Z
    Last Update Time:      2019-10-06T12:28:50Z
    Message:               Experiment has succeeded because Objective goal has reached
    Reason:                ExperimentSucceeded
    Status:                True
    Type:                  Succeeded
  Current Optimal Trial:
    Observation:
      Metrics:
        Name:   accuracy_1
        Value:  1
    Parameter Assignments:
      Name:          --learning_rate
      Value:         0.018532845700535087
      Name:          --batch_size
      Value:         109
  Start Time:        2019-10-06T12:25:44Z
  Trials:            4
  Trials Running:    2
  Trials Succeeded:  2
Events:              <none>
```

When the spec.Status.Condition becomes ```Succeeded```, the experiment is finished.

You can monitor your results in Katib UI. 
Access Katib UI via Kubeflow dashboard if you have used standard installation or port-forward the `katib-ui` service if you have installed manually.

```
kubectl -n kubeflow port-forward svc/katib-ui 8080:80
```

You can access the Katib UI using this URL: ```http://localhost:8080/katib/```.

### Cleanups

Delete installed components using `kubectl delete -f` on the respective folders. 

## Quick Start

Please see [Quick Start Guide](./docs/quick-start.md)

## Who are using Katib?

Please see [adopters.md](./docs/community/adopters.md)

## CONTRIBUTING

Please feel free to test the system! [developer-guide.md](./docs/developer-guide.md) is a good starting point for developers.

[1]: https://en.wikipedia.org/wiki/Hyperparameter_optimization
[2]: https://en.wikipedia.org/wiki/Neural_architecture_search
[3]: https://static.googleusercontent.com/media/research.google.com/ja//pubs/archive/bcb15507f4b52991a0783013df4222240e942381.pdf
