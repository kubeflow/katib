Table of Contents
=================

   * [Developer Guide](#developer-guide)
      * [Requirements](#requirements)
      * [Build from source code](#build-from-source-code)
      * [Implement a new algorithm and use it in katib](#implement-a-new-algorithm-and-use-it-in-katib)
         * [Implement the algorithm](#implement-the-algorithm)
         * [Make the algorithm a GRPC server](#make-the-algorithm-a-grpc-server)
         * [Use the algorithm in katib.](#use-the-algorithm-in-katib)
         * [Contribute the algorithm to katib](#contribute-the-algorithm-to-katib)
            * [Unit Test](#unit-test)
            * [E2E Test (Optional)](#e2e-test-optional)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)

# Developer Guide

## Requirements

- [Go](https://golang.org/)
- [Dep](https://golang.github.io/dep/)
- [Docker](https://docs.docker.com/) (17.05 or later.)

## Build from source code

Check source code as follows:

```bash
make build
```

You can deploy katib v1alpha3 manifests into a k8s cluster as follows:

```bash
make deploy
```

You can undeploy katib v1alpha3 manifests from a k8s cluster as follows:

```bash
make undeploy
```
