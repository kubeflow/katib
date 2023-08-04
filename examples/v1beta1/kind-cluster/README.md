# Katib Example with Kind Cluster

Follow this example to run Katib Experiment on your **local laptop** with
[Kind](https://github.com/kubernetes-sigs/kind/) cluster. This example doesn't
require any public or private cloud to run Katib Experiments.

## Prerequisites

Install the following tools to run the example:

- [Docker](https://docs.docker.com/get-docker) >= 20.10
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) >= 0.13
- [`kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl) >= 1.25

## Installation

Run the following command to create Kind cluster with the
[Katib components](https://www.kubeflow.org/docs/components/katib/hyperparameter/#katib-components):

```
./deploy.sh
```

If the above script was successful, Katib components will be running:

```
$ kubectl get pods -n kubeflow

NAME                                READY   STATUS      RESTARTS   AGE
katib-controller-566595bdd8-x7z6w   1/1     Running     0          67s
katib-db-manager-57cd769cdb-x4lnz   1/1     Running     0          67s
katib-mysql-7894994f88-7l8nd        1/1     Running     0          67s
katib-ui-5767cfccdc-nt6mz           1/1     Running     0          67s
```

## Run Katib Experiment

You can use various [Katib interfaces](https://www.kubeflow.org/docs/components/katib/overview/#katib-interfaces)
to run your first Katib Experiment.

For example, create Hyperparameter Tuning Katib Experiment with
[random search algorithm](https://www.kubeflow.org/docs/components/katib/experiment/#random-search)
using `kubectl`:

```
kubectl create -f https://raw.githubusercontent.com/kubeflow/katib/master/examples/v1beta1/hp-tuning/random.yaml
```

This example uses a MXNet neural network to train an image classification model
using the MNIST dataset. You can check the training container source code
[here](../trial-images/mxnet-mnist).
The Experiment runs twelve training jobs (Trials) and tunes the following hyperparameters:

- Learning Rate (`lr`).
- Number of layers (`num-layers`).
- Neural network optimizer (`optimizer`).

After creating above example, check the
[Experiment](https://www.kubeflow.org/docs/components/katib/overview/#experiment) status:

```
$ kubectl get experiment random -n kubeflow

NAME     TYPE      STATUS   AGE
random   Running   True     6m19s
```

Check the [Suggestion](https://www.kubeflow.org/docs/components/katib/overview/#suggestion) status:

```
$ kubectl get suggestion -n kubeflow

NAME     TYPE      STATUS   REQUESTED   ASSIGNED   AGE
random   Running   True     4           4          6m21s
```

Check the [Trials](https://www.kubeflow.org/docs/components/katib/overview/#trial) statuses:

```
$ kubectl get trial -n kubeflow

NAME              TYPE        STATUS   AGE
random-9hmdjqk9   Running     True     99s
random-cf7tfss2   Succeeded   True     5m21s
random-fr5lfn2x   Running     True     5m21s
random-z9wqm7xh   Running     True     5m21s
```

You can get the best hyperparameters with the following command:

```
$ kubectl get experiment random -n kubeflow -o jsonpath='{range .status.currentOptimalTrial.parameterAssignments[*]}{.name}: {.value}{"\n"}{end}'

lr: 0.028162244250364066
num-layers: 5
optimizer: sgd
```

To view created Experiment in Katib UI, follow
[this guide](https://www.kubeflow.org/docs/components/katib/hyperparameter/#accessing-the-katib-ui).

## Cleanup

To cleanup Kind cluster run:

```
kind delete cluster
```
