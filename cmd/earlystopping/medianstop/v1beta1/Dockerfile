FROM python:3.10-slim

ARG TARGETARCH
ENV TARGET_DIR /opt/katib
ENV EARLY_STOPPING_DIR cmd/earlystopping/medianstop/v1beta1
ENV PYTHONPATH ${TARGET_DIR}:${TARGET_DIR}/pkg/apis/manager/v1beta1/python

RUN if [ "${TARGETARCH}" = "ppc64le" ] || [ "${TARGETARCH}" = "arm64" ]; then \
    apt-get -y update && \
    apt-get -y install gfortran libopenblas-dev liblapack-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*; \
  fi

ADD ./pkg/ ${TARGET_DIR}/pkg/
ADD ./${EARLY_STOPPING_DIR}/ ${TARGET_DIR}/${EARLY_STOPPING_DIR}/

WORKDIR  ${TARGET_DIR}/${EARLY_STOPPING_DIR}

RUN pip install --prefer-binary --no-cache-dir -r requirements.txt
RUN chgrp -R 0 ${TARGET_DIR} \
  && chmod -R g+rwX ${TARGET_DIR}

ENTRYPOINT ["python", "main.py"]
