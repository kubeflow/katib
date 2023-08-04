<h1 align="center">
    <img src="./docs/images/logo-title.png" alt="logo" width="200">
  <br>
</h1>

[![Build Status](https://github.com/kubeflow/katib/actions/workflows/test-go.yaml/badge.svg?branch=master)](https://github.com/kubeflow/katib/actions/workflows/test-go.yaml?branch=master)
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
with out of the box support for [Kubeflow Training Operator](https://github.com/kubeflow/training-operator),
[Argo Workflows](https://github.com/argoproj/argo-workflows), [Tekton Pipelines](https://github.com/tektoncd/pipeline)
and many more.

Katib stands for `secretary` in Arabic.

# Search Algorithms

Katib supports several search algorithms. Follow the
[Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/experiment/#search-algorithms-in-detail)
to know more about each algorithm and check the
[Suggestion service guide](/docs/new-algorithm-service.md) to implement your
custom algorithm.

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

To perform above algorithms Katib supports the following frameworks:

- [Goptuna](https://github.com/c-bata/goptuna)
- [Hyperopt](https://github.com/hyperopt/hyperopt)
- [Optuna](https://github.com/optuna/optuna)
- [Scikit Optimize](https://github.com/scikit-optimize/scikit-optimize)

# Installation

For the various Katib installs check the
[Kubeflow guide](https://www.kubeflow.org/docs/components/katib/hyperparameter/#katib-setup).
Follow the next steps to install Katib standalone.

## Prerequisites

This is the minimal requirements to install Katib:

- Kubernetes >= 1.25
- `kubectl` >= 1.25

## Latest Version

For the latest Katib version run this command:

```
kubectl apply -k "github.com/kubeflow/katib.git/manifests/v1beta1/installs/katib-standalone?ref=master"
```

## Release Version

For the specific Katib release (for example `v0.14.0`) run this command:

```
kubectl apply -k "github.com/kubeflow/katib.git/manifests/v1beta1/installs/katib-standalone?ref=v0.14.0"
```

Make sure that all Katib components are running:

```
$ kubectl get pods -n kubeflow

NAME                                READY   STATUS      RESTARTS   AGE
katib-controller-566595bdd8-hbxgf   1/1     Running     0          36s
katib-db-manager-57cd769cdb-4g99m   1/1     Running     0          36s
katib-mysql-7894994f88-5d4s5        1/1     Running     0          36s
katib-ui-5767cfccdc-pwg2x           1/1     Running     0          36s
```

For the Katib Experiments check the [complete examples list](./examples/v1beta1).

# Quickstart

You can run your first HyperParameter Tuning Experiment using [Katib Python SDK](./sdk/python/v1beta1).

In the following example we are going to maximize a simple objective function:
$F(a,b) = 4a - b^2$. The bigger $a$ and the lesser $b$ value, the bigger the function value $F$.

```python
import kubeflow.katib as katib

# Step 1. Create an objective function.
def objective(parameters):
    # Import required packages.
    import time
    time.sleep(5)
    # Calculate objective function.
    result = 4 * int(parameters["a"]) - float(parameters["b"]) ** 2
    # Katib parses metrics in this format: <metric-name>=<metric-value>.
    print(f"result={result}")

# Step 2. Create HyperParameter search space.
parameters = {
    "a": katib.search.int(min=10, max=20),
    "b": katib.search.double(min=0.1, max=0.2)
}

# Step 3. Create Katib Experiment.
katib_client = katib.KatibClient()
name = "tune-experiment"
katib_client.tune(
    name=name,
    objective=objective,
    parameters=parameters,
    objective_metric_name="result",
    max_trial_count=12
)

# Step 4. Get the best HyperParameters.
print(katib_client.get_optimal_hyperparameters(name))
```

# Documentation

- Check
  [the Katib getting started guide](https://www.kubeflow.org/docs/components/katib/hyperparameter/#example-using-random-search-algorithm).

- Learn about Katib **Concepts** in this
  [guide](https://www.kubeflow.org/docs/components/katib/overview/#katib-concepts).

- Learn about Katib **Interfaces** in this
  [guide](https://www.kubeflow.org/docs/components/katib/overview/#katib-interfaces).

- Learn about Katib **Components** in this
  [guide](https://www.kubeflow.org/docs/components/katib/hyperparameter/#katib-components).

- Know more about Katib in the [presentations and demos list](./docs/presentations.md).

# Community

We are always growing our community and invite new users and AutoML enthusiasts
to contribute to the Katib project. The following links provide information
about getting involved in the community:

- Subscribe to the
  [AutoML calendar](https://calendar.google.com/calendar/u/0/r?cid=ZDQ5bnNpZWZzbmZna2Y5MW8wdThoMmpoazRAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
  to attend Working Group bi-weekly community meetings.

- Check the
  [AutoML and Training Working Group meeting notes](https://docs.google.com/document/d/1MChKfzrKAeFRtYqypFbMXL6ZIc_OgijjkvbqmwRV-64/edit).

- If you use Katib, please update [the adopters list](ADOPTERS.md).

## Contributing

Please feel free to test the system! [Developer guide](./docs/developer-guide.md)
is a good starting point for our developers.

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
