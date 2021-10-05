# PyTorch MNIST Image Classification Example

This is PyTorch MNIST image classification training container with saving metrics
to the file or printing to the StdOut. It uses convolutional neural network to
train the model.

Katib uses this training container in some Experiments, for instance in the
[file Metrics Collector example](../../metrics-collector/file-metrics-collector.yaml#L55-L64)
or in the [PyTorchJob example](../../kubeflow-training-operator/pytorchjob-mnist.yaml#L47-L54).
