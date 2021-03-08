FROM pytorch/pytorch:1.0-cuda10.0-cudnn7-runtime

RUN mkdir -p /katib /opt/mnist/src

WORKDIR /opt/mnist/src
ADD mnist.py /opt/mnist/src

RUN chgrp -R 0 /opt/mnist \
  && chmod -R g+rwX /opt/mnist \
  && chgrp -R 0 /katib \
  && chmod -R g+rwX /katib

ENTRYPOINT ["python", "/opt/mnist/src/mnist.py"]
