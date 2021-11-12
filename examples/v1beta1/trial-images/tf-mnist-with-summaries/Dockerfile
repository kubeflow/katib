FROM tensorflow/tensorflow:2.7.0

ADD examples/v1beta1/trial-images/tf-mnist-with-summaries /opt/tf-mnist-with-summaries
WORKDIR /opt/tf-mnist-with-summaries

# Add folder for the logs.
RUN mkdir /katib

RUN chgrp -R 0 /opt/tf-mnist-with-summaries \
  && chmod -R g+rwX /opt/tf-mnist-with-summaries \
  && chgrp -R 0 /katib \
  && chmod -R g+rwX /katib

ENTRYPOINT ["python3", "/opt/tf-mnist-with-summaries/mnist.py"]
