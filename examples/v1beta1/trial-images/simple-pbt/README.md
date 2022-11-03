# Toy PBT problem for benchmarking adaptive learning rate

The goal is to optimize this trainable's accuracy. The accuracy increases
fastest at the optimal lr, which is a function of the current accuracy.
The optimal lr schedule for this problem is the triangle wave as follows.
Note that many lr schedules for real models also follow this shape:

```
 best lr
  ^
  |    /\
  |   /  \
  |  /    \
  | /      \
  ------------> accuracy
```

In this problem, using PBT with a population of 2-4 is sufficient to
roughly approximate this lr schedule. Higher population sizes will yield
faster convergence. Training will not converge without PBT.

If you want to read more about this example, vist the 
[ray](https://github.com/ray-project/ray/blob/7f1bacc7dc9caf6d0ec042e39499bbf1d9a7d065/python/ray/tune/examples/README.rst) 
documentation.

Katib uses this training container in some Experiments, for instance in the
[PBT example](../../hp-tuning/simple-pbt.yaml#L44-L52).
