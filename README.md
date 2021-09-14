<h1 align="center">
    <img src="./docs/images/logo-title.png" alt="logo" width="200">
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
[Neural Architecture Search](https://en.wikipedia.org/wiki/Neural_architecture_search).

Katib is the project which is agnostic to machine learning (ML) frameworks.
It can tune hyperparameters of applications written in any language of the
users’ choice and natively supports many ML frameworks, such as
[TensorFlow](https://www.tensorflow.org/), [Apache MXNet](https://mxnet.apache.org/),
[PyTorch](https://pytorch.org/), [XGBoost](https://xgboost.readthedocs.io/en/latest/), and others.

Katib can perform training jobs using any Kubernetes
[Custom Resources](https://www.kubeflow.org/docs/components/katib/trial-template/)
with out of the box support for [Kubeflow Training Operators](https://github.com/kubeflow/tf-operator),
[Argo Workflows](https://github.com/argoproj/argo-workflows), [Tekton Pipelines](https://github.com/tektoncd/pipeline)
and many more.

Katib stands for `secretary` in Arabic.

# Installation

## Prerequisites

- Kubernetes >= 1.17
- `kubectl` >= 1.21

## Latest Version

Install Katib with the single command:

```
kubectl apply -k "github.com/kubeflow/katib.git/manifests/installs/katib-standalone?ref=master"
```

## Release Version

For the specific release (for example `v0.11.0`) run this command:

```
kubectl apply -k "github.com/kubeflow/katib.git/manifests/installs/katib-standalone?ref=v0.11.0"
```

Learn more about various Katib installs in the
[Kubeflow guides](https://www.kubeflow.org/docs/components/katib/hyperparameter/#katib-setup)

# Search Algorithms

Katib currently supports several search algorithms. Follow the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/experiment/#search-algorithms-in-detail)
to know more about each algorithm.

<table>
  <tbody>
    <tr align="center">
      <td>
        <b>Hyperparameter Tuning</b>
      </td>
      <td>
        <b>Neural Architecture Search</b>
      </td>
      <td>
        <b>Early Stopping</b>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#random-search">Random Search</a>
      </td>
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#random-search">ENAS</a>
      </td>
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#random-search">Median Stop</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#random-search">Grid Search</a>
      </td>
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#random-search">DARTS</a>
      </td>
      <td>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#random-search">Bayesian Optimization</a>
      </td>
      <td>
      </td>
      <td>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#random-search">TPE</a>
      </td>
      <td>
      </td>
      <td>
      </td>
    </tr>
  </tbody>
</table>

## Hyperparameter Tuning

- [Random Search](https://en.wikipedia.org/wiki/Hyperparameter_optimization#Random_search)
- [Tree of Parzen Estimators (TPE)](https://papers.nips.cc/paper/4443-algorithms-for-hyper-parameter-optimization.pdf)
- [Multivariate TPE](https://tech.preferred.jp/en/blog/multivariate-tpe-makes-optuna-even-more-powerful/)
- [Grid Search](https://en.wikipedia.org/wiki/Hyperparameter_optimization#Grid_search)
- [Hyperband](https://arxiv.org/pdf/1603.06560.pdf)
- [Bayesian Optimization](https://arxiv.org/pdf/1012.2599.pdf)
- [Covariance Matrix Adaptation Evolution Strategy (CMA-ES)](https://arxiv.org/abs/1604.00772)
- [Sobol's Quasirandom Sequence](https://dl.acm.org/doi/10.1145/641876.641879)

## Neural Architecture Search

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
During 1.3 we've worked on a new iteration of the UI, which is rewritten in
Angular and is utilizing the common code of the other Kubeflow [dashboards](https://github.com/kubeflow/kubeflow/tree/master/components/crud-web-apps).

The users are currently able to list, delete and create Experiments in their
cluster via this new UI as well as inspect the owned Trials. One important
missing functionalities are the ability to edit the Trial templates ConfigMaps
and view Neural Architecture Search models. Check [this Project](https://github.com/kubeflow/katib/projects/1)
to monitor the current progress.

![katibui](./docs/images/katib-ui.png)

To use the old Katib UI you can update the Katib image `newName` with the previous
image tag `docker.io/kubeflowkatib/katib-ui:v0.11.1` in the [Kustomize](./manifests/v1beta1/installs/katib-standalone/kustomization.yaml#L29)
manifests.

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
git clone -b v1.2-branch git@github.com:kubeflow/manifests.git
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

Note that your [kustomize](https://kustomize.io/) version should be >= 3.2.
To install Katib run:

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

## Community

We are always growing our community and invite new users and AutoML enthusiasts
to contribute to the Katib project. The following links provide information
about getting involved in the community:

- If you use Katib, please update [the adopters list](ADOPTERS.md).

- Subscribe
  [to the calendar](https://calendar.google.com/calendar/u/0/r?cid=ZDQ5bnNpZWZzbmZna2Y5MW8wdThoMmpoazRAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
  to attend the AutoML WG community meeting.

- Check
  [the AutoML WG meeting notes](https://docs.google.com/document/d/1MChKfzrKAeFRtYqypFbMXL6ZIc_OgijjkvbqmwRV-64/edit).

- Join
  [the AutoML WG Slack channel](https://kubeflow.slack.com/archives/C018PMV53NW).

- Learn more about Katib in
  [the presentations and demos list](./docs/presentations.md).

### Blog posts

- [Kubeflow Katib: Scalable, Portable and Cloud Native System for AutoML](https://blog.kubeflow.org/katib/)
  (by Andrey Velichkevich)

### Events

- [AutoML and Training WG Summit. 16th of July 2021](https://docs.google.com/document/d/1vGluSPHmAqEr8k9Dmm82RcQ-MVnqbYYSfnjMGB-aPuo/edit?usp=sharing)

## Contributing

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
