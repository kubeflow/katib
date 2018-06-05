# Katib

[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/katib)](https://goreportcard.com/report/github.com/kubeflow/katib)

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
Each component communicates with others via GRPC and the API is defined at `api/api.proto`.

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

## StudyConfig

In the Study config file, we define the feasible space of parameters and configuration of a kubernetes job. Examples of such Study configs are in the `conf` directory. The configuration items are as follows:

- name: Study name
- owner: Owner
- objectivevaluename: Name of the objective value. Your evaluated software should be print log `{objectivevaluename}={objective value}` in std-io.
- optimizationtype: Optimization direction of the objective value. 1=maximize 2=minimize
- suggestalgorithm: [random, grid, hyperband] now
- suggestionparameters: Parameter of the algorithm. Set name-value style.
    - In random suggestion
        - SuggestionNum: How many suggestions will Katib create.
        - MaxParallel: Max number of run on kubernetes
    - In grid suggestion
        - MaxParallel: Max number of run on kubernetes
        - GridDefault: default number of grid
        - name: [parameter name] grid number of specified parameter.
- metrics: The value you want to save to modeldb besides objectivevaluename.
- image: docker image name
- mount
    - pvc: pvc
    - path: MountPath in container
- pullsecret: Name of Image pull secret
- gpu: number of GPU (If you want to run cpu task, set 0 or delete this parameter)
- command: commands
- parameterconfigs: define feasible space
    - configs
        - name : parameter space
        - parametertype: 1=float, 2=int, 4=categorical
        - feasible
            - min
            - max
            - list (for categorical)

## Web UI

Katib provides a Web UI based on ModelDB(https://github.com/mitdbg/modeldb). The ingress setting is defined in [`manifests/modeldb/frontend/ingress.yaml`](manifests/modeldb/frontend/ingress.yaml).

![katib-demo](https://user-images.githubusercontent.com/10014831/38241910-64fb0646-376e-11e8-8b98-c26e577f3935.gif)

## CONTRIBUTING

Please feel free to test the system! [developer-guide.md](./docs/developer-guide.md) is a good starting point for  developers.

## TODOs

* Integrate KubeFlow (TensorFlow, Caffe2 and PyTorch operators)
* Support Early Stopping
* Enrich the GUI
