# Katib

[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/katib)](https://goreportcard.com/report/github.com/kubeflow/katib)

<img src="./img/Katib_Logo.png" width="320px">

Hyperparameter Tuning on Kubernetes.
This project is inspired by [Google vizier](https://static.googleusercontent.com/media/research.google.com/ja//pubs/archive/bcb15507f4b52991a0783013df4222240e942381.pdf). Katib is a scalable and flexible hyperparameter tuning framework and is tightly integrated with kubernetes. Also it does not depend on a specific Deep Learning framework e.g. TensorFlow, MXNet, and PyTorch).

## Name

Katib stands for `secretary` in Arabic. As `Vizier` stands for a high official or a prime minister in Arabic, this project Katib is named in the honor of Vizier.

## Concepts in Google Vizier

As in Google Vizier, Katib also has the concepts of Study, Trial and Suggestion.

### Study

Represents a single optimization run over a feasible space. Each Study contains a configuration describing the feasible space, as well as a set of Trials. It is assumed that objective function f(x) does not change in the course of a Study.

### Trial

A Trial is a list of parameter values, x, that will lead to a single evaluation of f(x). A Trial can be “Completed”, which means that it has been evaluated and the objective value f(x) has been assigned to it, otherwise it is “Pending”.
One trial corresponds to one k8s Job.

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

## CLI Documentation

Please refer to [katib-cli.md](./docs/CLI/katib-cli.md).

## CONTRIBUTING

Please feel free to test the system! [developer-guide.md](./docs/developer-guide.md) is a good starting point for developers.

## TODOs

* Integrate KubeFlow (TensorFlow, Caffe2 and PyTorch operators)
* Support Early Stopping
* Enrich the GUI
