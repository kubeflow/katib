# Table of Contents

- [Table of Contents](#table-of-contents)
- [Developer Guide](#developer-guide)
  - [Requirements](#requirements)
  - [Build from source code](#build-from-source-code)
  - [Modify controller APIs](#modify-controller-apis)
  - [Controller Flags](#controller-flags)
  - [Workflow design](#workflow-design)
  - [Implement a new algorithm and use it in Katib](#implement-a-new-algorithm-and-use-it-in-katib)
  - [Algorithm settings documentation](#algorithm-settings-documentation)
  - [Katib UI documentation](#katib-ui-documentation)
  - [Design proposals](#design-proposals)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)

# Developer Guide

This developer guide is for people who want to contribute to the Katib project.
If you're interesting in using Katib in your machine learning project,
see the following user guides:

- [Concepts](https://www.kubeflow.org/docs/components/katib/overview/)
  in Katib, hyperparameter tuning, and neural architecture search.
- [Getting started with Katib](https://kubeflow.org/docs/components/katib/hyperparameter/).
- Detailed guide to [configuring and running a Katib
  experiment](https://kubeflow.org/docs/components/katib/experiment/).

## Requirements

- [Go](https://golang.org/) (1.13 or later)
- [Docker](https://docs.docker.com/) (17.05 or later.)
- [kustomize](https://kustomize.io/) (3.2 or later)

## Build from source code

Check source code as follows:

```bash
make build REGISTRY=<image-registry> TAG=<image-tag>
```

You can deploy Katib v1beta1 manifests into a k8s cluster as follows:

```bash
make deploy
```

You can undeploy Katib v1beta1 manifests from a k8s cluster as follows:

```bash
make undeploy
```

## Modify controller APIs

If you want to modify Katib controller APIs you have to
generate deepcopy, clientset, listers, informers, open-api and python SDK with changed APIs.
You can update necessary files as follows:

```bash
make generate
```

## Controller Flags

Below is a list of command-line flags accepted by Katib controller:

| Name                            | Type                      | Default   | Description                                                                                                            |
| ------------------------------- | ------------------------- | --------- | ---------------------------------------------------------------------------------------------------------------------- |
| enable-grpc-probe-in-suggestion | bool                      | true      | Enable grpc probe in suggestions                                                                                       |
| experiment-suggestion-name      | string                    | "default" | The implementation of suggestion interface in experiment controller                                                    |
| metrics-addr                    | string                    | ":8080"   | The address the metric endpoint binds to                                                                               |
| trial-resources                 | []schema.GroupVersionKind | null      | The list of resources that can be used as trial template, in the form: Kind.version.group (e.g. TFJob.v1.kubeflow.org) |
| webhook-inject-securitycontext  | bool                      | false     | Inject the securityContext of container[0] in the sidecar                                                              |
| webhook-port                    | int                       | 8443      | The port number to be used for admission webhook server                                                                |

## Workflow design

Please see [workflow-design.md](./workflow-design.md).

## Implement a new algorithm and use it in Katib

Please see [new-algorithm-service.md](./new-algorithm-service.md).

## Algorithm settings documentation

Please see [algorithm-settings.md](./algorithm-settings.md).

## Katib UI documentation

Please see [Katib UI README](https://github.com/kubeflow/katib/tree/master/pkg/ui/v1beta1).

## Design proposals

Please see [proposals](./proposals).
