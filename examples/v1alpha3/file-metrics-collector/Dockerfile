FROM pytorch/pytorch:1.0-cuda10.0-cudnn7-runtime

WORKDIR /var
ADD mnist.py /var

ENTRYPOINT ["python", "/var/mnist.py"]
