FROM tensorflow/tensorflow:2.8.0

ADD test/functional/v1beta1/suggestion/pbt /opt/pbt
WORKDIR /opt/pbt

RUN chgrp -R 0 /opt/pbt \
  && chmod -R g+rwX /opt/pbt

ENTRYPOINT ["python3", "/opt/pbt/pbt_test.py"]
