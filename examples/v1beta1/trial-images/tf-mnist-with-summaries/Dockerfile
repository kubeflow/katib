FROM python:3.11-slim

ARG TARGETARCH

ADD examples/v1beta1/trial-images/tf-mnist-with-summaries /opt/tf-mnist-with-summaries

WORKDIR /opt/tf-mnist-with-summaries

RUN if [ "${TARGETARCH}" = "arm64" ]; then \
    apt-get -y update && \
    apt-get -y install gfortran libpcre2-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*; \
    fi

RUN if [ "${TARGETARCH}" = "ppc64le" ]; then \
      pip install --prefer-binary --no-cache-dir -r requirements.txt \
        --extra-index-url https://wheels.developerfirst.ibm.com/ppc64le/linux/+simple/; \
    else \
      pip install --prefer-binary --no-cache-dir -r requirements.txt; \
    fi
RUN chgrp -R 0 /opt/tf-mnist-with-summaries \
  && chmod -R g+rwX /opt/tf-mnist-with-summaries

ENTRYPOINT ["python3", "/opt/tf-mnist-with-summaries/mnist.py"]
