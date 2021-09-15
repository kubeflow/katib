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

Katib supports several search algorithms. Follow the
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
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#neural-architecture-search-based-on-enas">ENAS</a>
      </td>
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/early-stopping/#median-stopping-rule">Median Stop</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#grid-search">Grid Search</a>
      </td>
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#differentiable-architecture-search-darts">DARTS</a>
      </td>
      <td>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#bayesian-optimization">Bayesian Optimization</a>
      </td>
      <td>
      </td>
      <td>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#tree-of-parzen-estimators-tpe">TPE</a>
      </td>
      <td>
      </td>
      <td>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#multivariate-tpe">Multivariate TPE</a>
      </td>
      <td>
      </td>
      <td>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#covariance-matrix-adaptation-evolution-strategy-cma-es">CMA-ES</a>
      </td>
      <td>
      </td>
      <td>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#sobols-quasirandom-sequence">Sobol's Quasirandom Sequence</a>
      </td>
      <td>
      </td>
      <td>
      </td>
    </tr>
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#hyperband">HyperBand</a>
      </td>
      <td>
      </td>
      <td>
      </td>
    </tr>
  </tbody>
</table>

Follow [this guide](/docs/new-algorithm-service.md) to implement your custom algorithm
service in Katib.

# Documentation

Run your first Katib Experiment in the
[getting started guide](https://www.kubeflow.org/docs/components/katib/hyperparameter/#example-using-random-algorithm).

Learn about Katib concepts in the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/overview/#katib-concepts).

Know more about various Katib interfaces in the
[introduction guide](https://www.kubeflow.org/docs/components/katib/overview/#katib-interfaces).

TODO (andreyvelich): Add Katib components section from the website.

Learn more about Katib in the [presentations and demos list](./docs/presentations.md).

# Community

We are always growing our community and invite new users and AutoML enthusiasts
to contribute to the Katib project. The following links provide information
about getting involved in the community:

- Subscribe
  [to the calendar](https://calendar.google.com/calendar/u/0/r?cid=ZDQ5bnNpZWZzbmZna2Y5MW8wdThoMmpoazRAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
  to attend the AutoML WG community meeting.

- Check
  [the AutoML WG meeting notes](https://docs.google.com/document/d/1MChKfzrKAeFRtYqypFbMXL6ZIc_OgijjkvbqmwRV-64/edit).

- If you use Katib, please update [the adopters list](ADOPTERS.md).

## Contributing

Please feel free to test the system!
[developer-guide.md](./docs/developer-guide.md) is a good starting point
for developers.

## Blog posts

- [Kubeflow Katib: Scalable, Portable and Cloud Native System for AutoML](https://blog.kubeflow.org/katib/)
  (by Andrey Velichkevich)

## Events

- [AutoML and Training WG Summit. 16th of July 2021](https://docs.google.com/document/d/1vGluSPHmAqEr8k9Dmm82RcQ-MVnqbYYSfnjMGB-aPuo/edit?usp=sharing)

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
