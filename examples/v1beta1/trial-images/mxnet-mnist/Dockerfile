FROM python:3.10-slim

ARG TARGETARCH
ENV LD_LIBRARY_PATH "${LD_LIBRARY_PATH}:/opt/arm/armpl_22.0.2_gcc-11.2/lib"

COPY examples/v1beta1/trial-images/mxnet-mnist /opt/mxnet-mnist

WORKDIR /opt/mxnet-mnist

RUN apt-get -y update \
 && apt-get -y install libgomp1 \
 && if [ "${TARGETARCH}" = "arm64" ]; then \
      apt-get -y install wget; \
    elif [ "${TARGETARCH}" = "amd64" ]; then \
      apt-get -y install libquadmath0; \
    fi \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

RUN if [ "${TARGETARCH}" = "arm64" ]; then \
      /opt/mxnet-mnist/install-arm-performance-libraries.sh; \
    fi
RUN pip install --prefer-binary -r requirements.txt
RUN chgrp -R 0 /opt/mxnet-mnist \
  && chmod -R g+rwX /opt/mxnet-mnist

ENTRYPOINT ["python3", "/opt/mxnet-mnist/mnist.py"]
