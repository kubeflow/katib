# Kubeflow Katib

[![Build Status](https://github.com/kubeflow/katib/actions/workflows/test-go.yaml/badge.svg?branch=master)](https://github.com/kubeflow/katib/actions/workflows/test-go.yaml?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/kubeflow/katib/badge.svg?branch=master)](https://coveralls.io/github/kubeflow/katib?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/katib)](https://goreportcard.com/report/github.com/kubeflow/katib)
[![Releases](https://img.shields.io/github/release-pre/kubeflow/katib.svg?sort=semver)](https://github.com/kubeflow/katib/releases)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://www.kubeflow.org/docs/about/community/#kubeflow-slack-channels)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/9941/badge)](https://www.bestpractices.dev/projects/9941)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeflow%2Fkatib.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeflow%2Fkatib?ref=badge_shield)

<h1 align="center">
    <img src="./docs/images/logo-title.png" alt="logo" width="200">
  <br>
</h1>

Kubeflow Katib is a Kubernetes-native project for automated machine learning (AutoML).
Katib supports
[Hyperparameter Tuning](https://en.wikipedia.org/wiki/Hyperparameter_optimization),
[Early Stopping](https://en.wikipedia.org/wiki/Early_stopping) and
[Neural Architecture Search](https://en.wikipedia.org/wiki/Neural_architecture_search).

Katib is the project which is agnostic to machine learning (ML) frameworks.
It can tune hyperparameters of applications written in any language of the
users’ choice and natively supports many ML frameworks, such as
[TensorFlow](https://www.tensorflow.org/), [PyTorch](https://pytorch.org/), [XGBoost](https://xgboost.readthedocs.io/en/latest/), and others.

Katib can perform training jobs using any Kubernetes
[Custom Resources](https://www.kubeflow.org/docs/components/katib/trial-template/)
with out of the box support for [Kubeflow Training Operator](https://github.com/kubeflow/training-operator),
[Argo Workflows](https://github.com/argoproj/argo-workflows), [Tekton Pipelines](https://github.com/tektoncd/pipeline)
and many more.

Katib stands for `secretary` in Arabic.

## Search Algorithms

Katib supports several search algorithms. Follow the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/user-guides/hp-tuning/configure-algorithm/#hp-tuning-algorithms)
to know more about each algorithm and check the
[this guide](https://www.kubeflow.org/docs/components/katib/user-guides/hp-tuning/configure-algorithm/#use-custom-algorithm-in-katib)
to implement your custom algorithm.

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
    <tr align="center">
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#pbt">Population Based Training</a>
      </td>
      <td>
      </td>
      <td>
      </td>
    </tr>
  </tbody>
</table>

To perform the above algorithms Katib supports the following frameworks:

- [Goptuna](https://github.com/c-bata/goptuna)
- [Hyperopt](https://github.com/hyperopt/hyperopt)
- [Optuna](https://github.com/optuna/optuna)
- [Scikit Optimize](https://github.com/scikit-optimize/scikit-optimize)

## Prerequisites

Please check [the official Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/installation/#prerequisites)
for prerequisites to install Katib.

## Installation

Please follow [the Kubeflow Katib guide](https://www.kubeflow.org/docs/components/katib/installation/#installing-katib)
for the detailed instructions on how to install Katib.

### Installing the Control Plane

Run the following command to install the latest stable release of Katib control plane:

```
kubectl apply -k "github.com/kubeflow/katib.git/manifests/v1beta1/installs/katib-standalone?ref=v0.17.0"
```

Run the following command to install the latest changes of Katib control plane:

```
kubectl apply -k "github.com/kubeflow/katib.git/manifests/v1beta1/installs/katib-standalone?ref=master"
```

For the Katib Experiments check the [complete examples list](./examples/v1beta1).

### Installing the Python SDK

Katib implements [a Python SDK](https://pypi.org/project/kubeflow-katib/) to simplify creation of
hyperparameter tuning jobs for Data Scientists.

Run the following command to install the latest stable release of Katib SDK:

```sh
pip install -U kubeflow-katib
```

## Getting Started

Please refer to [the getting started guide](https://www.kubeflow.org/docs/components/katib/getting-started/#getting-started-with-katib-python-sdk)
to quickly create your first hyperparameter tuning Experiment using the Python SDK.

## Community

The following links provide information on how to get involved in the community:

- Attend [the bi-weekly AutoML and Training Working Group](https://bit.ly/2PWVCkV)
  community meeting.
- Join our [`#kubeflow-katib`](https://www.kubeflow.org/docs/about/community/#kubeflow-slack-channels)
  Slack channel.
- Check out [who is using Katib](ADOPTERS.md) and [presentations about Katib project](docs/presentations.md).

## Contributing

Please refer to the [CONTRIBUTING guide](CONTRIBUTING.md).

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


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeflow%2Fkatib.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeflow%2Fkatib?ref=badge_large)