FROM python:3.9-slim

ADD examples/v1beta1/trial-images/simple-pbt /opt/pbt
WORKDIR /opt/pbt

RUN python3 -m pip install -r requirements.txt

RUN chgrp -R 0 /opt/pbt \
  && chmod -R g+rwX /opt/pbt

ENTRYPOINT ["python3", "/opt/pbt/pbt_test.py"]
