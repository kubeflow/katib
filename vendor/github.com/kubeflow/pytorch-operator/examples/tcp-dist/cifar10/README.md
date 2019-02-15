### Distributed cifar model for e2e test

This folder containers Dockerfile and distributed cufar model for e2e test.

**Build Image**

The default image name and tag is `kubeflow/pytorch-dist-cifar-test:1.0`.

```shell
docker build -f Dockerfile -t kubeflow/pytorch-dist-cifar-test:1.0 ./
```

**Create the mnist PyTorch job**

```
kubectl create -f ./pytorch_job_cifar.yaml
```
