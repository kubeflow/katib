
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

You can create PyTorch Job by defining a PyTorchJob config file. See [distributed MNIST example](https://github.com/kubeflow/pytorch-operator/blob/master/examples/tcp-dist/mnist/v1alpha2/pytorch_job_mnist.yaml) config file. You may change the config file based on your requirements.

```
cat examples/tcp-dist/mnist/v1alpha2/pytorch_job_mnist.yaml
```
Deploy the PyTorchJob resource to start training:

```
kubectl create -f examples/tcp-dist/mnist/v1alpha2/pytorch_job_mnist.yaml
```
You should now be able to see the created pods matching the specified number of replicas.

```
kubectl get pods -l pytorch_job_name=pytorch-tcp-dist-mnist
```
Training should run for about 10 epochs and takes 5-10 minutes on a cpu cluster. Logs can be inspected to see its training progress. 

```
PODNAME=$(kubectl get pods -l pytorch_job_name=pytorch-tcp-dist-mnist,task_index=0 -o name)
kubectl logs -f ${PODNAME}
```
## Monitoring a PyTorch Job

```
kubectl get -o yaml pytorchjobs pytorch-tcp-dist-mnist
```
See the status section to monitor the job status. Here is sample output when the job is successfully completed.

```
apiVersion: v1
items:
- apiVersion: kubeflow.org/v1alpha2
  kind: PyTorchJob
  metadata:
    clusterName: ""
    creationTimestamp: 2018-09-14T14:31:02Z
    generation: 1
    labels:
      app.kubernetes.io/deploy-manager: ksonnet
    name: pytorch-tcp-dist-mnist
    resourceVersion: "21599007"
    selfLink: /apis/kubeflow.org/v1alpha2/namespaces/kubeflow/pytorchjobs/pytorch-tcp-dist-mnist
    uid: ce70c8de-b82a-11e8-8b09-42010aa000d2
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
            - image: johnugeorge/pytorch-tcp-mnist:3.0
              name: pytorch
              ports:
              - containerPort: 23456
                name: pytorchjob-port
              resources: {}
      Worker:
        replicas: 3
        restartPolicy: OnFailure
        template:
          metadata:
            creationTimestamp: null
          spec:
            containers:
            - image: johnugeorge/pytorch-tcp-mnist:3.0
              name: pytorch
              ports:
              - containerPort: 23456
                name: pytorchjob-port
              resources: {}
  status:
    completionTime: 2018-09-14T14:33:13Z
    conditions:
    - lastTransitionTime: 2018-09-14T14:31:02Z
      lastUpdateTime: 2018-09-14T14:31:02Z
      message: PyTorchJob pytorch-tcp-dist-mnist is created.
      reason: PyTorchJobCreated
      status: "True"
      type: Created
    - lastTransitionTime: 2018-09-14T14:31:02Z
      lastUpdateTime: 2018-09-14T14:31:05Z
      message: PyTorchJob pytorch-tcp-dist-mnist is running.
      reason: PyTorchJobRunning
      status: "False"
      type: Running
    - lastTransitionTime: 2018-09-14T14:31:02Z
      lastUpdateTime: 2018-09-14T14:33:13Z
      message: PyTorchJob pytorch-tcp-dist-mnist is successfully completed.
      reason: PyTorchJobSucceeded
      status: "True"
      type: Succeeded
    pytorchReplicaStatuses:
      Master: {}
      Worker: {}
    startTime: 2018-09-14T14:31:04Z
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""

```
