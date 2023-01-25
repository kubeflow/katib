# Katib Examples

Katib is an open source project which uses Kubernetes CRD to run Automated
Machine Learning (AutoML) tasks. To know more about Katib follow the
[official guides](https://www.kubeflow.org/docs/components/katib/overview/).

This directory contains examples of Katib Experiments in action. To install Katib on your
Kubernetes cluster check the
[setup guide](https://www.kubeflow.org/docs/components/katib/hyperparameter/#katib-setup).
You can use various [Katib interfaces](https://www.kubeflow.org/docs/components/katib/overview/#katib-interfaces)
to run these examples.

For a complete description of the Katib Experiment specification follow the
[configuration guide](https://www.kubeflow.org/docs/components/katib/experiment/#configuration-spec)

## Local Cluster Example

Get started with Katib Experiments from your **local laptop** and
[Kind](https://github.com/kubernetes-sigs/kind/) cluster by following
[this example](./kind-cluster).

## AutoML Algorithms

The following examples show various AutoML algorithms in Katib.

### Hyperparameter Tuning

Check the [Hyperparameter Tuning](https://www.kubeflow.org/docs/components/katib/overview/#hyperparameters-and-hyperparameter-tuning)
Experiments for the following algorithms:

- [Random Search](./hp-tuning/random.yaml)

- [Grid Search](./hp-tuning/grid.yaml)

- [Bayesian Optimization](./hp-tuning/bayesian-optimization.yaml)

- [Tree of Parzen Estimators (TPE)](./hp-tuning/tpe.yaml)

- [Multivariate TPE](./hp-tuning/multivariate-tpe.yaml)

- [Covariance Matrix Adaptation Evaluation Strategy (CMA-ES)](./hp-tuning/cma-es.yaml)

- [Sobol's Quasirandom Sequence](./hp-tuning/sobol.yaml)

- [HyperBand](./hp-tuning/hyperband.yaml)

- [PBT](./hp-tuning/simple-pbt.yaml)

### Neural Architecture Search

Check the [Neural Architecture Search](https://www.kubeflow.org/docs/components/katib/overview/#neural-architecture-search)
Experiments for the following algorithms:

- [Efficient Neural Architecture Search (ENAS)](./nas/enas-gpu.yaml)

- [Differentiable Architecture Search (DARTS)](./nas/darts-gpu.yaml)

### Early Stopping

Improve your Hyperparameter Tuning Experiments with the following
[Early Stopping](https://www.kubeflow.org/docs/components/katib/early-stopping/) algorithms:

- [Median Stopping Rule](./early-stopping/median-stop.yaml)

## Katib Python SDK Examples

To learn more about Katib Python SDK check [this directory](./sdk).

## Resume Katib Experiments

You can use different resume policies in Katib Experiments. Follow
[this guide](https://www.kubeflow.org/docs/components/katib/resume-experiment/)
to know more about it. Check the following examples:

- [Resume From Volume](./resume-experiment/from-volume-resume.yaml)

- [Resume Long Running Experiment](./resume-experiment/long-running-resume.yaml)

## Metrics Collector

Katib supports the various metrics collectors and metrics strategies.
Check the [official guide](https://www.kubeflow.org/docs/components/katib/experiment/#configuration-spec)
to know more about it. In this directory you can find the following examples:

- [File Metrics Collector](./metrics-collector/file-metrics-collector.yaml)

- [Custom Metrics Collector](./metrics-collector/custom-metrics-collector.yaml)

- [Metrics Collection Strategy](./metrics-collector/metrics-collection-strategy.yaml)

## Trial Template

You can specify different settings for your Trial template. To know more about it
follow [this guide](https://www.kubeflow.org/docs/components/katib/trial-template/#use-trial-template-to-submit-experiment).
Check the following examples:

- [Trial with ConfigMap Source](./trial-template/trial-configmap-source.yaml)

- [Trial with Metadata Substitution](./trial-template/trial-metadata-substitution.yaml)

## Trial Images

Check the following images for the Trial containers:

- [Tensorflow MNIST with summaries](./trial-images/tf-mnist-with-summaries)

- [MXNet MNIST](./trial-images/mxnet-mnist)

- [PyTorch MNIST](./trial-images/pytorch-mnist)

- [ENAS Keras CNN CIFAR-10](./trial-images/enas-cnn-cifar10)

- [DARTS PyTorch CNN CIFAR-10](./trial-images/darts-cnn-cifar10)

- [PBT proof of concept](./trial-images/simple-pbt)

## Katib with Kubeflow Training Operator

Katib has out of the box support for the [Kubeflow Training Operators](https://github.com/kubeflow/training-operator) to
perform Trial's [Worker job](https://www.kubeflow.org/docs/components/katib/overview/#trial).
Check the following examples for the various distributed operators:

- [TFJob MNIST with summaries](./kubeflow-training-operator/tfjob-mnist-with-summaries.yaml)

- [PyTorchJob MNIST](./kubeflow-training-operator/pytorchjob-mnist.yaml)

- [MXJob BytePS](./kubeflow-training-operator/mxjob-byteps.yaml)

- [XGBoostJob LightGBM](./kubeflow-training-operator/xgboostjob-lightgbm.yaml)

- [MPIJob Horovod](./kubeflow-training-operator/mpijob-horovod.yaml)

## Katib with Kubeflow Pipelines

To run Katib with [Kubeflow Pipelines](https://github.com/kubeflow/pipelines) check
[these examples](./kubeflow-pipelines).

## Katib with Argo Workflows

To know more about using [Argo Workflows](https://github.com/argoproj/argo-workflows)
in Katib check [this directory](./argo).

## Katib with Tekton Pipelines

To know more about using [Tekton Pipelines](https://github.com/tektoncd/pipeline)
in Katib check [this directory](./tekton).

## FPGA Support in Katib Experiments

You can run Katib Experiments on [FPGA](https://en.wikipedia.org/wiki/Field-programmable_gate_array)
based instances. For more information check [these examples](./fpga).
