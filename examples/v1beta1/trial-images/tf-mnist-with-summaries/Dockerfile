FROM python:3.10-slim

ARG TARGETARCH

ADD examples/v1beta1/trial-images/tf-mnist-with-summaries /opt/tf-mnist-with-summaries

WORKDIR /opt/tf-mnist-with-summaries

RUN if [ "${TARGETARCH}" = "arm64" ]; then \
    apt-get -y update && \
    apt-get -y install gfortran libpcre3 libpcre3-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*; \
    fi

RUN pip install --prefer-binary --no-cache-dir -r requirements.txt
RUN chgrp -R 0 /opt/tf-mnist-with-summaries \
  && chmod -R g+rwX /opt/tf-mnist-with-summaries

ENTRYPOINT ["python3", "/opt/tf-mnist-with-summaries/mnist.py"]
