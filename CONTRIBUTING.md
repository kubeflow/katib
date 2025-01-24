# Developer Guide

This developer guide is for people who want to contribute to the Katib project.
If you're interesting in using Katib in your machine learning project,
see the following guides:

- [Getting started with Katib](https://kubeflow.org/docs/components/katib/hyperparameter/).
- [How to configure Katib Experiment](https://kubeflow.org/docs/components/katib/experiment/).
- [Katib architecture and concepts](https://www.kubeflow.org/docs/components/katib/reference/architecture/)
  for hyperparameter tuning and neural architecture search.

## Requirements

- [Go](https://golang.org/) (1.22 or later)
- [Docker](https://docs.docker.com/) (24.0 or later)
- [Docker Buildx](https://docs.docker.com/build/buildx/) (0.8.0 or later)
- [Java](https://docs.oracle.com/javase/8/docs/technotes/guides/install/install_overview.html) (8 or later)
- [Python](https://www.python.org/) (3.11 or later)
- [kustomize](https://kustomize.io/) (4.0.5 or later)
- [pre-commit](https://pre-commit.com/)

## Build from source code

**Note** that your Docker Desktop should
[enable containerd image store](https://docs.docker.com/desktop/containerd/#enable-the-containerd-image-store)
to build multi-arch images. Check source code as follows:

```bash
make build REGISTRY=<image-registry> TAG=<image-tag>
```

If you are using an Apple Silicon machine and encounter the "rosetta error: bss_size overflow," go to Docker Desktop -> General and uncheck "Use Rosetta for x86_64/amd64 emulation on Apple Silicon."

To use your custom images for the Katib components, modify
[Kustomization file](https://github.com/kubeflow/katib/blob/master/manifests/v1beta1/installs/katib-standalone/kustomization.yaml)
and [Katib Config](https://github.com/kubeflow/katib/blob/master/manifests/v1beta1/installs/katib-standalone/katib-config.yaml)

You can deploy Katib v1beta1 manifests into a Kubernetes cluster as follows:

```bash
make deploy
```

You can undeploy Katib v1beta1 manifests from a Kubernetes cluster as follows:

```bash
make undeploy
```

## Technical and style guide

The following guidelines apply primarily to Katib,
but other projects like [Training Operator](https://github.com/kubeflow/training-operator) might also adhere to them.

## Go Development

When coding:

- Follow [effective go](https://go.dev/doc/effective_go) guidelines.
- Run locally [`make check`](https://github.com/kubeflow/katib/blob/46173463027e4fd2e604e25d7075b2b31a702049/Makefile#L31)
  to verify if changes follow best practices before submitting PRs.

Testing:

- Use [`cmp.Diff`](https://pkg.go.dev/github.com/google/go-cmp/cmp#Diff) instead of `reflect.Equal`, to provide useful comparisons.
- Define test cases as maps instead of slices to avoid dependencies on the running order.
  Map key should be equal to the test case name.

## Modify controller APIs

If you want to modify Katib controller APIs, you have to
generate deepcopy, clientset, listers, informers, open-api and Python SDK with the changed APIs.
You can update the necessary files as follows:

```bash
make generate
```

## Controller Flags

Below is a list of command-line flags accepted by Katib controller:

| Name         | Type   | Default | Description                                                                                                                      |
| ------------ | ------ | ------- | -------------------------------------------------------------------------------------------------------------------------------- |
| katib-config | string | ""      | The katib-controller will load its initial configuration from this file. Omit this flag to use the default configuration values. |

## DB Manager Flags

Below is a list of command-line flags accepted by Katib DB Manager:

| Name            | Type          | Default      | Description                                                         |
| --------------- | ------------- | -------------| ------------------------------------------------------------------- |
| connect-timeout | time.Duration | 60s          | Timeout before calling error during database connection             |
| listen-address  | string        | 0.0.0.0:6789 | The network interface or IP address to receive incoming connections |

## Katib admission webhooks

Katib uses three [Kubernetes admission webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/).

1. `validator.experiment.katib.kubeflow.org` -
   [Validating admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#validatingadmissionwebhook)
   to validate the Katib Experiment before the creation.

1. `defaulter.experiment.katib.kubeflow.org` -
   [Mutating admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook)
   to set the [default values](../pkg/apis/controller/experiments/v1beta1/experiment_defaults.go)
   in the Katib Experiment before the creation.

1. `mutator.pod.katib.kubeflow.org` - Mutating admission webhook to inject the metrics
   collector sidecar container to the training pod. Learn more about the Katib's
   metrics collector in the
   [Kubeflow documentation](https://www.kubeflow.org/docs/components/katib/user-guides/metrics-collector/).

You can find the YAMLs for the Katib webhooks
[here](../manifests/v1beta1/components/webhook/webhooks.yaml).

**Note:** If you are using a private Kubernetes cluster, you have to allow traffic
via `TCP:8443` by specifying the firewall rule and you have to update the master
plane CIDR source range to use the Katib webhooks

### Katib cert generator

Katib Controller has the internal `cert-generator` to generate certificates for the webhooks.

Once Katib is deployed in the Kubernetes cluster, the `cert-generator` follows these steps:

- Generate the self-signed certificate and private key.

- Update a Kubernetes Secret with the self-signed TLS certificate and private key.
- Patch the webhooks with the `CABundle`.

Once the `cert-generator` finished, the Katib controller starts to register controllers such as `experiment-controller` to the manager.

You can find the `cert-generator` source code [here](../pkg/certgenerator/v1beta1).

NOTE: the Katib also supports the [cert-manager](https://cert-manager.io/) to generate certs for the admission webhooks instead of using cert-generator.
You can find the installation with the cert-manager [here](../manifests/v1beta1/installs/katib-cert-manager).

## Implement a new algorithm and use it in Katib

Please see [new-algorithm-service.md](./new-algorithm-service.md).

## Katib UI documentation

Please see [Katib UI README](../pkg/ui/v1beta1).

## Design proposals

Please see [proposals](./proposals).

## Code Style

### pre-commit

Make sure to install [pre-commit](https://pre-commit.com/) (`pip install
pre-commit`) and run `pre-commit install` from the root of the repository at
least once before creating git commits.

The pre-commit [hooks](../.pre-commit-config.yaml) ensure code quality and
consistency. They are executed in CI. PRs that fail to comply with the hooks
will not be able to pass the corresponding CI gate. The hooks are only executed
against staged files unless you run `pre-commit run --all`, in which case,
they'll be executed against every file in the repository.

Specific programmatically generated files listed in the `exclude` field in
[.pre-commit-config.yaml](../.pre-commit-config.yaml) are deliberately excluded
from the hooks.
