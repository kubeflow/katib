FROM ubuntu:16.04

RUN apt-get update && \
    apt-get install -y wget python3-dev gcc && \
    wget https://bootstrap.pypa.io/get-pip.py && \
    python3 get-pip.py

RUN pip3 install mxnet

ADD . /opt/mxnet-mnist
WORKDIR /opt/mxnet-mnist

ENTRYPOINT ["python3", "/opt/mxnet-mnist/mnist.py"]
