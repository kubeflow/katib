FROM python:3.11-slim

ARG TARGETARCH
ADD examples/v1beta1/trial-images/simple-pbt /opt/pbt

WORKDIR /opt/pbt


RUN if [ "${TARGETARCH}" = "ppc64le" ]; then \
      python3 -m pip install --prefer-binary -r requirements.txt \
        --extra-index-url https://wheels.developerfirst.ibm.com/ppc64le/linux/+simple/; \
    else \
      python3 -m pip install --prefer-binary -r requirements.txt; \
    fi
RUN chgrp -R 0 /opt/pbt \
  && chmod -R g+rwX /opt/pbt

ENTRYPOINT ["python3", "/opt/pbt/pbt_test.py"]
