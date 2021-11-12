# Tensorflow MNIST Classification With Summaries Example

This is Tensorflow MNIST image classification training container that outputs TF summaries.
It uses convolutional neural network to train the model.

If you want to read more about this example, visit the official
[tensorflow](https://github.com/tensorflow/tensorflow/blob/7462dcaae1e8cfe1dfd0c62dd6083f9749a9d827/tensorflow/examples/tutorials/mnist/mnist_with_summaries.py)
GitHub repository.

Katib uses this training container in some Experiments, for instance in the
[TF Event Metrics Collector](../../metrics-collector/tfevent-metrics-collector.yaml#L55-L64).
