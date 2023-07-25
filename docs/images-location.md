# Katib Images Location

Here you can find the location for images that are used in Katib.

## Katib Components Images

The following table shows images for the
[Katib components](https://www.kubeflow.org/docs/components/katib/hyperparameter/#katib-components).

<table>
  <tbody>
    <tr align="center">
      <td>
        <b>Image Name</b>
      </td>
      <td>
        <b>Description</b>
      </td>
      <td>
        <b>Location</b>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/katib-controller</code>
      </td>
      <td>
        Katib Controller
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/tree/master/cmd/katib-controller/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/katib-ui</code>
      </td>
      <td>
        Katib User Interface
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/tree/master/cmd/ui/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/katib-db-manager</code>
      </td>
      <td>
        Katib DB Manager
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/tree/master/cmd/db-manager/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/mysql</code>
      </td>
      <td>
        Katib MySQL DB
      </td>
      <td>
        <a href="https://github.com/docker-library/mysql/blob/c506174eab8ae160f56483e8d72410f8f1e1470f/8.0/Dockerfile.debian">Dockerfile</a>
      </td>
    </tr>
  </tbody>
</table>

## Katib Metrics Collectors Images

The following table shows images for the
[Katib Metrics Collectors](https://www.kubeflow.org/docs/components/katib/experiment/#metrics-collector).

<table>
  <tbody>
    <tr align="center">
      <td>
        <b>Image Name</b>
      </td>
      <td>
        <b>Description</b>
      </td>
      <td>
        <b>Location</b>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/file-metrics-collector</code>
      </td>
      <td>
        File Metrics Collector
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/metricscollector/v1beta1/file-metricscollector/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/tfevent-metrics-collector</code>
      </td>
      <td>
        Tensorflow Event Metrics Collector
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/metricscollector/v1beta1/tfevent-metricscollector/Dockerfile">Dockerfile</a>
      </td>
    </tr>
  </tbody>
</table>

## Katib Suggestions and Early Stopping Images

The following table shows images for the
[Katib Suggestions](https://www.kubeflow.org/docs/components/katib/experiment/#search-algorithms-in-detail)
and the [Katib Early Stopping algorithms](https://www.kubeflow.org/docs/components/katib/early-stopping/).

<table>
  <tbody>
    <tr align="center">
      <td>
        <b>Image Name</b>
      </td>
      <td>
        <b>Description</b>
      </td>
      <td>
        <b>Location</b>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/suggestion-hyperopt</code>
      </td>
      <td>
        <a href="https://github.com/hyperopt/hyperopt">Hyperopt</a> Suggestion
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/suggestion/hyperopt/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/suggestion-skopt</code>
      </td>
      <td>
        <a href="https://github.com/scikit-optimize/scikit-optimize">Skopt</a> Suggestion
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/suggestion/skopt/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/suggestion-optuna</code>
      </td>
      <td>
        <a href="https://github.com/optuna/optuna">Optuna</a> Suggestion
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/suggestion/optuna/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/suggestion-goptuna</code>
      </td>
      <td>
        <a href="https://github.com/c-bata/goptuna">Goptuna</a> Suggestion
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/suggestion/goptuna/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/suggestion-hyperband</code>
      </td>
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#hyperband">Hyperband</a> Suggestion
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/suggestion/hyperband/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/suggestion-enas</code>
      </td>
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#enas">ENAS</a> Suggestion
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/suggestion/nas/enas/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/suggestion-darts</code>
      </td>
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/experiment/#differentiable-architecture-search-darts">DARTS</a> Suggestion
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/suggestion/nas/darts/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/earlystopping-medianstop</code>
      </td>
      <td>
        <a href="https://www.kubeflow.org/docs/components/katib/early-stopping/#median-stopping-rule">Median Stopping Rule</a>
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/cmd/earlystopping/medianstop/v1beta1/Dockerfile">Dockerfile</a>
      </td>
    </tr>
  </tbody>
</table>

## Training Containers Images

The following table shows images for training containers which are used in the
[Katib Trials](https://www.kubeflow.org/docs/components/katib/experiment/#packaging-your-training-code-in-a-container-image).

<table>
  <tbody>
    <tr align="center">
      <td>
        <b>Image Name</b>
      </td>
      <td>
        <b>Description</b>
      </td>
      <td>
        <b>Location</b>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/mxnet-mnist</code>
      </td>
      <td>
        MXNet MNIST example with collecting metrics time
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/examples/v1beta1/trial-images/mxnet-mnist/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/pytorch-mnist-cpu</code>
      </td>
      <td>
        PyTorch MNIST example with printing metrics to the file or StdOut with CPU support
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/examples/v1beta1/trial-images/pytorch-mnist/Dockerfile.cpu">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/pytorch-mnist-gpu</code>
      </td>
      <td>
        PyTorch MNIST example with printing metrics to the file or StdOut with GPU support
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/examples/v1beta1/trial-images/pytorch-mnist/Dockerfile.gpu">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/tf-mnist-with-summaries</code>
      </td>
      <td>
        Tensorflow MNIST example with saving metrics in the summaries
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/examples/v1beta1/trial-images/tf-mnist-with-summaries/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/bytepsimage/mxnet</code>
      </td>
      <td>
        Distributed BytePS example for MXJob
      </td>
      <td>
        <a href="https://github.com/bytedance/byteps/blob/v0.2.5/docker/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/xgboost-lightgbm</code>
      </td>
      <td>
        Distributed LightGBM example for XGBoostJob
      </td>
      <td>
        <a href="https://github.com/kubeflow/xgboost-operator/blob/9c8c97d0125a8156f12b8ef5b93f99e709fb57ea/config/samples/lightgbm-dist/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflow/mpi-horovod-mnist</code>
      </td>
      <td>
        Distributed Horovod example for MPIJob
      </td>
      <td>
        <a href="https://github.com/kubeflow/mpi-operator/blob/947d396a9caf70d3c94bf587d5e5da32b70f0f53/examples/horovod/Dockerfile.cpu">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/inaccel/jupyter:lab</code>
      </td>
      <td>
        FPGA XGBoost with parameter tuning
      </td>
      <td>
        <a href="https://github.com/inaccel/jupyter/blob/master/lab/Dockerfile">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/enas-cnn-cifar10-gpu</code>
      </td>
      <td>
        Keras CIFAR-10 CNN example for ENAS with GPU support
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/examples/v1beta1/trial-images/enas-cnn-cifar10/Dockerfile.gpu">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/enas-cnn-cifar10-cpu</code>
      </td>
      <td>
        Keras CIFAR-10 CNN example for ENAS with CPU support
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/examples/v1beta1/trial-images/enas-cnn-cifar10/Dockerfile.cpu">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/darts-cnn-cifar10-gpu</code>
      </td>
      <td>
        PyTorch CIFAR-10 CNN example for DARTS with GPU support
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/examples/v1beta1/trial-images/darts-cnn-cifar10/Dockerfile.gpu">Dockerfile</a>
      </td>
    </tr>
    <tr align="center">
      <td>
        <code>docker.io/kubeflowkatib/darts-cnn-cifar10-cpu</code>
      </td>
      <td>
        PyTorch CIFAR-10 CNN example for DARTS with CPU support
      </td>
      <td>
        <a href="https://github.com/kubeflow/katib/blob/master/examples/v1beta1/trial-images/darts-cnn-cifar10/Dockerfile.cpu">Dockerfile</a>
      </td>
    </tr>
</table>
