FROM pytorch/pytorch:1.0-cuda10.0-cudnn7-runtime

ADD examples/v1beta1/pytorch-mnist /opt/pytorch-mnist
WORKDIR /opt/pytorch-mnist

# Add folder for the logs.
RUN mkdir /katib

RUN chgrp -R 0 /opt/pytorch-mnist \
  && chmod -R g+rwX /opt/pytorch-mnist \
  && chgrp -R 0 /katib \
  && chmod -R g+rwX /katib

ENTRYPOINT ["python3", "/opt/pytorch-mnist/mnist.py"]
