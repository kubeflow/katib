<h1 align="center">
    <img src="./docs/images/Katib_Logo.png" alt="logo" width="200">
  <br>
</h1>

[![Build Status](https://travis-ci.com/kubeflow/katib.svg?branch=master)](https://travis-ci.com/kubeflow/katib)
[![Coverage Status](https://coveralls.io/repos/github/kubeflow/katib/badge.svg?branch=master)](https://coveralls.io/github/kubeflow/katib?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/katib)](https://goreportcard.com/report/github.com/kubeflow/katib)
[![Releases](https://img.shields.io/github/release-pre/kubeflow/katib.svg?sort=semver)](https://github.com/kubeflow/katib/releases)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://kubeflow.slack.com/archives/C018PMV53NW)

Katib is a Kubernetes-native project for automated machine learning (AutoML).
Katib supports
[Hyperparameter Tuning](https://en.wikipedia.org/wiki/Hyperparameter_optimization),
[Early Stopping](https://en.wikipedia.org/wiki/Early_stopping) and
[Neural Architecture Search](https://en.wikipedia.org/wiki/Neural_architecture_search)

Katib is the project which is agnostic to machine learning (ML) frameworks.
It can tune hyperparameters of applications written in any language of the
users’ choice and natively supports many ML frameworks, such as TensorFlow,
MXNet, PyTorch, XGBoost, and others.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

# Table of Contents

- [Getting Started](#getting-started)
- [Name](#name)
- [Concepts in Katib](#concepts-in-katib)
  - [Experiment](#experiment)
  - [Suggestion](#suggestion)
  - [Trial](#trial)
  - [Worker Job](#worker-job)
  - [Search Algorithms](#search-algorithms)
    - [Hyperparameter Tuning](#hyperparameter-tuning)
    - [Neural Architecture Search](#neural-architecture-search)
- [Components in Katib](#components-in-katib)
- [Web UI](#web-ui)
- [GRPC API documentation](#grpc-api-documentation)
- [Installation](#installation)
  - [TF operator](#tf-operator)
  - [PyTorch operator](#pytorch-operator)
  - [Katib](#katib)
  - [Running examples](#running-examples)
  - [Katib SDK](#katib-sdk)
  - [Cleanups](#cleanups)
- [Quick Start](#quick-start)
- [Who are using Katib?](#who-are-using-katib)
- [CONTRIBUTING](#contributing)
- [Citation](#citation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

Created by [doctoc](https://github.com/thlorenz/doctoc).

## Getting Started

Follow the
[getting-started guide](https://www.kubeflow.org/docs/components/katib/hyperparameter/)
on the Kubeflow website.

## Name

Katib stands for `secretary` in Arabic.

## Concepts in Katib

For a detailed description of the concepts in Katib and AutoML, check the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/overview/).

Katib has the concepts of `Experiment`, `Suggestion`, `Trial` and `Worker Job`.

### Experiment

An `Experiment` represents a single optimization run over a feasible space.
Each `Experiment` contains a configuration:

1. **Objective**: What you want to optimize.
2. **Search Space**: Constraints for configurations describing the feasible space.
3. **Search Algorithm**: How to find the optimal configurations.

Katib `Experiment` is defined as a CRD. Check the detailed guide to
[configuring and running a Katib `Experiment`](https://kubeflow.org/docs/components/katib/experiment/)
in the Kubeflow docs.

### Suggestion

A `Suggestion` is a set of hyperparameter values that the hyperparameter tuning
process has proposed. Katib creates a `Trial` to evaluate
the suggested set of values.

Katib `Suggestion` is defined as a CRD.

### Trial

A `Trial` is one iteration of the hyperparameter tuning process.
A `Trial` corresponds to one worker job instance with a list of parameter
assignments. The list of parameter assignments corresponds to a `Suggestion`.

Each `Experiment` runs several `Trials`. The `Experiment` runs the `Trials` until
it reaches either the objective or the configured maximum number of `Trials`.

Katib `Trial` is defined as a CRD.

### Worker Job

The `Worker Job` is the process that runs to evaluate a `Trial` and calculate
its objective value.

The `Worker Job` can be any type of Kubernetes resource or
[Kubernetes CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).
Follow the [`Trial` template guide](https://www.kubeflow.org/docs/components/katib/trial-template/#custom-resource)
to support your own Kubernetes resource in Katib.

Katib has these CRD examples in upstream:

- [Kubernetes `Job`](https://kubernetes.io/docs/concepts/workloads/controllers/job/)

- [Kubeflow `TFJob`](https://www.kubeflow.org/docs/components/training/tftraining/)

- [Kubeflow `PyTorchJob`](https://www.kubeflow.org/docs/components/training/pytorch/)

- [Kubeflow `MPIJob`](https://www.kubeflow.org/docs/components/training/mpi/)

- [Tekton `Pipeline`](https://github.com/tektoncd/pipeline)

Thus, Katib supports multiple frameworks with the help of different job kinds.

### Search Algorithms

Katib currently supports several search algorithms. Follow the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/experiment/#search-algorithms-in-detail)
to know more about each algorithm.

#### Hyperparameter Tuning

- [Random Search](https://en.wikipedia.org/wiki/Hyperparameter_optimization#Random_search)
- [Tree of Parzen Estimators (TPE)](https://papers.nips.cc/paper/4443-algorithms-for-hyper-parameter-optimization.pdf)
- [Grid Search](https://en.wikipedia.org/wiki/Hyperparameter_optimization#Grid_search)
- [Hyperband](https://arxiv.org/pdf/1603.06560.pdf)
- [Bayesian Optimization](https://arxiv.org/pdf/1012.2599.pdf)
- [Covariance Matrix Adaptation Evolution Strategy (CMA-ES)](https://arxiv.org/abs/1604.00772)

#### Neural Architecture Search

- [Efficient Neural Architecture Search (ENAS)](https://github.com/kubeflow/katib/tree/master/pkg/suggestion/v1beta1/nas/enas)
- [Differentiable Architecture Search (DARTS)](https://github.com/kubeflow/katib/tree/master/pkg/suggestion/v1beta1/nas/darts)

## Components in Katib

Katib consists of several components as shown below. Each component is running
on Kubernetes as a deployment. Each component communicates with others via GRPC
and the API is defined at `pkg/apis/manager/v1beta1/api.proto`.

- Katib main components:
  - `katib-db-manager` - the GRPC API server of Katib which is the DB Interface.
  - `katib-mysql` - the data storage backend of Katib using mysql.
  - `katib-ui` - the user interface of Katib.
  - `katib-controller` - the controller for the Katib CRDs in Kubernetes.

## Web UI

Katib provides a Web UI.
You can visualize general trend of Hyper parameter space and
each training history. You can use
[random-example](https://github.com/kubeflow/katib/blob/master/examples/v1beta1/random-example.yaml)
or
[other examples](https://github.com/kubeflow/katib/blob/master/examples/v1beta1)
to generate a similar UI. Follow the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/hyperparameter/#katib-ui)
to access the Katib UI.
![katibui](./docs/images/katib-ui.png)

## GRPC API documentation

Check the [Katib v1beta1 API reference docs](https://www.kubeflow.org/docs/reference/katib/v1beta1/katib/).

## Installation

For standard installation of Katib with support for all job operators,
install Kubeflow.
Follow the documentation:

- [Kubeflow installation guide](https://www.kubeflow.org/docs/started/getting-started/)
- [Kubeflow Katib guides](https://www.kubeflow.org/docs/components/katib/).

If you install Katib with other Kubeflow components,
you can't submit Katib jobs in Kubeflow namespace. Check the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/hyperparameter/#example-using-random-algorithm)
to know more about it.

Alternatively, if you want to install Katib manually with TF and PyTorch
operators support, follow these steps:

Create Kubeflow namespace:

```
kubectl create namespace kubeflow
```

Clone Kubeflow manifest repository:

```
git clone git@github.com:kubeflow/manifests.git
Set `MANIFESTS_DIR` to the cloned folder.
export MANIFESTS_DIR=<cloned-folder>
```

### TF operator

For installing TF operator, run the following:

```
cd "${MANIFESTS_DIR}/tf-training/tf-job-crds/base"
kustomize build . | kubectl apply -f -
cd "${MANIFESTS_DIR}/tf-training/tf-job-operator/base"
kustomize build . | kubectl apply -f -
```

### PyTorch operator

For installing PyTorch operator, run the following:

```
cd "${MANIFESTS_DIR}/pytorch-job/pytorch-job-crds/base"
kustomize build . | kubectl apply -f -
cd "${MANIFESTS_DIR}/pytorch-job/pytorch-operator/base/"
kustomize build . | kubectl apply -f -
```

### Katib

Finally, you can install Katib:

```
git clone git@github.com:kubeflow/katib.git
make deploy
```

Check if all components are running successfully:

```
kubectl get pods -n kubeflow
```

Expected output:

```
NAME                                READY   STATUS    RESTARTS   AGE
katib-controller-858d6cc48c-df9jc   1/1     Running   1          20m
katib-db-manager-7966fbdf9b-w2tn8   1/1     Running   0          20m
katib-mysql-7f8bc6956f-898f9        1/1     Running   0          20m
katib-ui-7cf9f967bf-nm72p           1/1     Running   0          20m
pytorch-operator-55f966b548-9gq9v   1/1     Running   0          20m
tf-job-operator-796b4747d8-4fh82    1/1     Running   0          21m
```

### Running examples

After deploy everything, you can run examples to verify the installation.

This is an example for TF operator:

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/v1beta1/tfjob-example.yaml
```

This is an example for PyTorch operator:

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/v1beta1/pytorchjob-example.yaml
```

Check the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/hyperparameter/#example-using-random-algorithm)
how to monitor your `Experiment` status.

You can view your results in Katib UI.
If you used standard installation, access the Katib UI via Kubeflow dashboard.
Otherwise, port-forward the `katib-ui`:

```
kubectl -n kubeflow port-forward svc/katib-ui 8080:80
```

You can access the Katib UI using this URL: `http://localhost:8080/katib/`.

### Katib SDK

Katib supports Python SDK:

- Check the [Katib v1beta1 SDK documentation](https://github.com/kubeflow/katib/tree/master/sdk/python/v1beta1).

Run `make generate` to update Katib SDK.

### Cleanups

To delete installed TF and PyTorch operator run `kubectl delete -f`
on the respective folders.

To delete Katib run `make undeploy`.

## Quick Start

Please follow the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/hyperparameter/#examples)
to submit your first Katib experiment.

## Who are using Katib?

Please see [ADOPTERS.md](ADOPTERS.md).

## CONTRIBUTING

Please feel free to test the system!
[developer-guide.md](./docs/developer-guide.md) is a good starting point
for developers.

## Citation

If you use Katib in a scientific publication, we would appreciate
citations to the following paper:

[A Scalable and Cloud-Native Hyperparameter Tuning System](https://arxiv.org/abs/2006.02085), George _et al._, arXiv:2006.02085, 2020.

Bibtex entry:

```
@misc{george2020katib,
    title={A Scalable and Cloud-Native Hyperparameter Tuning System},
    author={Johnu George and Ce Gao and Richard Liu and Hou Gang Liu and Yuan Tang and Ramdoot Pydipaty and Amit Kumar Saha},
    year={2020},
    eprint={2006.02085},
    archivePrefix={arXiv},
    primaryClass={cs.DC}
}
```
