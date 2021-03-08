FROM mxnet/python:latest_cpu_native_py3

ADD . /opt/mxnet-mnist
WORKDIR /opt/mxnet-mnist

RUN chgrp -R 0 /opt/mxnet-mnist \
  && chmod -R g+rwX /opt/mxnet-mnist

ENTRYPOINT ["python3", "/opt/mxnet-mnist/mnist.py"]
