FROM python:3.9-slim

ADD examples/v1beta1/trial-images/mxnet-mnist /opt/mxnet-mnist
WORKDIR /opt/mxnet-mnist

RUN pip install mxnet==1.9.0
RUN chgrp -R 0 /opt/mxnet-mnist \
  && chmod -R g+rwX /opt/mxnet-mnist

ENTRYPOINT ["python3", "/opt/mxnet-mnist/mnist.py"]
