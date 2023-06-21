# Using Katib with Kubeflow Pipelines

The following examples show how to use Katib with
[Kubeflow Pipelines](https://github.com/kubeflow/pipelines).

Two different aspects are illustrated here:
A) How to orchestrate Katib experiments from Kubeflow pipelines using the Katib Kubeflow Component (Example 1 & 2)
B) How to use Katib to tune parameters of Kubeflow pipelines

You can find the Katib Component source code for the Kubeflow Pipelines
[here](https://github.com/kubeflow/pipelines/tree/master/components/kubeflow/katib-launcher).

## Prerequisites

You have to install the following Python SDK to run these examples:

- [`kfp`](https://pypi.org/project/kfp/) >= 1.8.12
- [`kubeflow-katib`](https://pypi.org/project/kubeflow-katib/) >= 0.13.0

In order to run parameter tuning over Kubeflow pipelines, additionally Katib needs to be setup to run with Argo workflow tasks. The setup is described within the example notebook (3).

## Multi-User Pipelines Setup

The Notebooks examples run Pipelines in multi-user mode and your Kubeflow Notebook
must have the appropriate `PodDefault` with the `pipelines.kubeflow.org` audience.

Please follow [this guide](https://www.kubeflow.org/docs/components/pipelines/sdk/connect-api/#multi-user-mode)
to give an access Kubeflow Notebook to run Kubeflow Pipelines.

## List of Examples

The following Pipelines are deployed from Kubeflow Notebook:

1) [Kubeflow E2E MNIST](kubeflow-e2e-mnist.ipynb)

2) [Katib Experiment with Early Stopping](early-stopping.ipynb)

3) [Tune parameters of a `MNIST` kubeflow pipeline with Katib](kubeflow-kfpv1-opt-mnist.ipynb)

The following Pipelines have to be compiled and uploaded to the Kubeflow Pipelines UI for examples 1 & 2:

- [MPIJob Horovod](mpi-job-horovod.py)
