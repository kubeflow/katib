FROM python:3.9-slim

ENV TARGET_DIR /opt/katib
ENV METRICS_COLLECTOR_DIR cmd/metricscollector/v1beta1/tfevent-metricscollector

ADD ./pkg/ ${TARGET_DIR}/pkg/
ADD ./${METRICS_COLLECTOR_DIR}/ ${TARGET_DIR}/${METRICS_COLLECTOR_DIR}/
WORKDIR  ${TARGET_DIR}/${METRICS_COLLECTOR_DIR}

RUN if [ "$(uname -m)" = "aarch64" ]; then \
    apt-get -y update && \
    apt-get -y install gfortran libpcre3 libpcre3-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*; \
    fi

RUN pip install --no-cache-dir -r requirements.txt

RUN chgrp -R 0 ${TARGET_DIR} \
  && chmod -R g+rwX ${TARGET_DIR}

ENV PYTHONPATH ${TARGET_DIR}:${TARGET_DIR}/pkg/apis/manager/v1beta1/python:${TARGET_DIR}/pkg/metricscollector/v1beta1/tfevent-metricscollector/::${TARGET_DIR}/pkg/metricscollector/v1beta1/common/

ENTRYPOINT ["python", "main.py"]
