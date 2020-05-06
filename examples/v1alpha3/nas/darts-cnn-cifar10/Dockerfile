ARG cuda_version=10.0
ARG cudnn_version=7
FROM pytorch/pytorch:1.0-cuda${cuda_version}-cudnn${cudnn_version}-runtime


ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/examples/v1alpha3/nas/darts-cnn-cifar10

ENTRYPOINT ["python3", "-u", "run_trial.py"]
