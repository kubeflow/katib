FROM tensorflow/tensorflow:2.9.1

ADD examples/v1beta1/trial-images/tf-mnist-with-summaries /opt/tf-mnist-with-summaries
WORKDIR /opt/tf-mnist-with-summaries

RUN chgrp -R 0 /opt/tf-mnist-with-summaries \
  && chmod -R g+rwX /opt/tf-mnist-with-summaries

ENTRYPOINT ["python3", "/opt/tf-mnist-with-summaries/mnist.py"]
