
# Kubernetes Custom Resource and Operator for PyTorch jobs

[![Build Status](https://travis-ci.org/kubeflow/pytorch-operator.svg?branch=master)](https://travis-ci.org/kubeflow/pytorch-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/pytorch-operator)](https://goreportcard.com/report/github.com/kubeflow/pytorch-operator)

## Overview

This repository contains the specification and implementation of `PyTorchJob` custom resource definition. Using this custom resource, users can create and manage PyTorch jobs like other built-in resources in Kubernetes. See [CRD definition](https://github.com/kubeflow/kubeflow/blob/master/kubeflow/pytorch-job/pytorch-operator.libsonnet#L11)
  
## Prerequisites

- Kubernetes >= 1.8
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)

## Installing PyTorch Operator

  Please refer to the installation instructions in the [Kubeflow user guide](https://www.kubeflow.org/docs/started/getting-started/). This installs `pytorchjob` CRD and `pytorch-operator` controller to manage the lifecycle of PyTorch jobs.

## Creating a PyTorch Job

You can create PyTorch Job by defining a PyTorchJob config file. See the manifests for the [distributed MNIST example](./examples/mnist/). You may change the config file based on your requirements.

```
cat examples/mnist/v1beta1/pytorch_job_mnist_gloo.yaml
```
Deploy the PyTorchJob resource to start training:

```
kubectl create -f examples/mnist/v1beta1/pytorch_job_mnist_gloo.yaml
```
You should now be able to see the created pods matching the specified number of replicas.

```
kubectl get pods -l pytorch_job_name=pytorch-dist-mnist
```
Training should run for about 10 epochs and takes 5-10 minutes on a cpu cluster. Logs can be inspected to see its training progress. 

```
PODNAME=$(kubectl get pods -l pytorch_job_name=pytorch-dist-mnist,task_index=0 -o name)
kubectl logs -f ${PODNAME}
```
## Monitoring a PyTorch Job

```
kubectl get -o yaml pytorchjobs pytorch-dist-mnist
```
See the status section to monitor the job status. Here is sample output when the job is successfully completed.

```
apiVersion: v1
items:
- apiVersion: kubeflow.org/v1beta1
  kind: PyTorchJob
  metadata:
    creationTimestamp: 2019-01-11T00:51:48Z
    generation: 1
    name: pytorch-dist-mnist
    namespace: kubeflow
    resourceVersion: "2146573"
    selfLink: /apis/kubeflow.org/v1beta1/namespaces/kubeflow/pytorchjobs/pytorch-dist-mnist
    uid: 13ad0e7f-153b-11e9-b5c1-42010a80001e
  spec:
    cleanPodPolicy: None
    pytorchReplicaSpecs:
      Master:
        replicas: 1
        restartPolicy: OnFailure
        template:
          metadata:
            creationTimestamp: null
          spec:
            containers:
            - args:
              - --backend
              - gloo
              image: gcr.io/tzaman/pytorch-dist-mnist-test:1.0
              name: pytorch
              ports:
              - containerPort: 23456
                name: pytorchjob-port
              resources:
                limits:
                  nvidia.com/gpu: "1"
      Worker:
        replicas: 1
        restartPolicy: OnFailure
        template:
          metadata:
            creationTimestamp: null
          spec:
            containers:
            - args:
              - --backend
              - gloo
              image: gcr.io/tzaman/pytorch-dist-mnist-test:1.0
              name: pytorch
              ports:
              - containerPort: 23456
                name: pytorchjob-port
              resources:
                limits:
                  nvidia.com/gpu: "1"
  status:
    completionTime: 2019-01-11T01:03:15Z
    conditions:
    - lastTransitionTime: 2019-01-11T00:51:48Z
      lastUpdateTime: 2019-01-11T00:51:48Z
      message: PyTorchJob pytorch-dist-mnist is created.
      reason: PyTorchJobCreated
      status: "True"
      type: Created
    - lastTransitionTime: 2019-01-11T00:57:22Z
      lastUpdateTime: 2019-01-11T00:57:22Z
      message: PyTorchJob pytorch-dist-mnist is running.
      reason: PyTorchJobRunning
      status: "False"
      type: Running
    - lastTransitionTime: 2019-01-11T01:03:15Z
      lastUpdateTime: 2019-01-11T01:03:15Z
      message: PyTorchJob pytorch-dist-mnist is successfully completed.
      reason: PyTorchJobSucceeded
      status: "True"
      type: Succeeded
    replicaStatuses:
      Master: {}
      Worker: {}
    startTime: 2019-01-11T00:57:22Z
```
