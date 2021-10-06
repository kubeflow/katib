# MXNet MNIST Image Classification Example

This is MXNet MNIST image classification training container with recording time
of the metrics. It uses only simple multilayer perceptron network (mlp).

If you want to read more about this example, visit the official
[incubator-mxnet](https://github.com/apache/incubator-mxnet/tree/1cf2fe5f8753042951bc0aacb6c95ddd3a904395/example/image-classification)
GitHub repository.

Katib uses this training container in some Experiments, for instance in the
[random search](../../hp-tuning/random.yaml#L55-L64).
