<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Developer Guide](#developer-guide)
  - [Requirements](#requirements)
  - [Build from source code](#build-from-source-code)
  - [Implement new suggestion algorithm](#implement-new-suggestion-algorithm)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Developer Guide

## Requirements

- [Go](https://golang.org/)
- [Dep](https://golang.github.io/dep/)
- [Docker](https://docs.docker.com/) (17.05 or later.)

## Build from source code

Check source code as follows:

```
make check
```

If there are some errors for go fmt, update the go fmt as follows:

```
make update
```

You can build all images from source for v1alpha2 as follows:

```bash
make images
```

You can deploy katib v1alpha2 manifests into a k8s cluster as follows:

```bash
make deploy
```

You can undeploy katib v1alpha2 manifests from a k8s cluster as follows:

```bash
make undeploy
```

## Run katib-controller natively

```bash
export KATIB_CORE_NAMESPACE=default
kubectl apply -f ./manifests/v1alpha2/katib-controller/crd-trial.yaml
kubectl apply -f ./manifests/v1alpha2/katib-controller/crd-suggestion.yaml
kubectl apply -f ./manifests/v1alpha2/katib-controller/crd-experiment.yaml
go build -o katib-controller github.com/kubeflow/katib/cmd/katib-controller/v1alpha2
mkdir -p /tmp/cert
cd /tmp/cert
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem
./katib-controller
```

## Implement new suggestion algorithm

Suggestion API is defined as GRPC service at `pkg/api/v1alpha1/api.proto`. Source code is [here](https://github.com/kubeflow/katib/blob/master/pkg/api/v1alpha1/api.proto). You can attach new algorithm easily.

- implement suggestion API
- make k8s service named `vizier-suggestion-{ algorithm-name }` and expose port 6789

And to add new suggestion service, you don't need to stop components ( vizier-core, modeldb, and anything) that are already running.
