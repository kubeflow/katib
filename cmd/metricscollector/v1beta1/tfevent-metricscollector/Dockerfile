FROM python:3.9

ENV TARGET_DIR /opt/katib
ENV METRICS_COLLECTOR_DIR cmd/metricscollector/v1beta1/tfevent-metricscollector
# tensorflow community build for aarch64
# https://github.com/tensorflow/build#tensorflow-builds
ENV PIP_EXTRA_INDEX_URL https://snapshots.linaro.org/ldcg/python-cache/

ADD ./pkg/ ${TARGET_DIR}/pkg/
ADD ./${METRICS_COLLECTOR_DIR}/ ${TARGET_DIR}/${METRICS_COLLECTOR_DIR}/
WORKDIR  ${TARGET_DIR}/${METRICS_COLLECTOR_DIR}

RUN if [ "$(uname -m)" = "aarch64" ]; then \
    pip install tensorflow-aarch64==2.7.0; \
    else \
    pip install tensorflow==2.7.0; \
    fi;
RUN pip install --no-cache-dir -r requirements.txt

RUN chgrp -R 0 ${TARGET_DIR} \
  && chmod -R g+rwX ${TARGET_DIR}

ENV PYTHONPATH ${TARGET_DIR}:${TARGET_DIR}/pkg/apis/manager/v1beta1/python:${TARGET_DIR}/pkg/metricscollector/v1beta1/tfevent-metricscollector/::${TARGET_DIR}/pkg/metricscollector/v1beta1/common/

ENTRYPOINT ["python", "main.py"]
