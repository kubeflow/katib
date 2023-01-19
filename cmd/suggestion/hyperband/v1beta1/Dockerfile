FROM alpine:3.15 AS downloader

ARG TARGETARCH
ENV GRPC_HEALTH_PROBE_VERSION v0.4.11

RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-${TARGETARCH} \
  && chmod +x /bin/grpc_health_probe

FROM python:3.10-slim

ARG TARGETARCH
ENV TARGET_DIR /opt/katib
ENV SUGGESTION_DIR cmd/suggestion/hyperband/v1beta1
ENV PYTHONPATH ${TARGET_DIR}:${TARGET_DIR}/pkg/apis/manager/v1beta1/python:${TARGET_DIR}/pkg/apis/manager/health/python

RUN if [ "${TARGETARCH}" = "ppc64le" ] || [ "${TARGETARCH}" = "arm64" ]; then \
    apt-get -y update && \
    apt-get -y install gfortran libopenblas-dev liblapack-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*; \
    fi

ADD ./pkg/ ${TARGET_DIR}/pkg/
ADD ./${SUGGESTION_DIR}/ ${TARGET_DIR}/${SUGGESTION_DIR}/
COPY --from=downloader /bin/grpc_health_probe /bin/grpc_health_probe

WORKDIR  ${TARGET_DIR}/${SUGGESTION_DIR}

RUN pip install --no-cache-dir -r requirements.txt
RUN chgrp -R 0 ${TARGET_DIR} \
  && chmod -R g+rwX ${TARGET_DIR}

ENTRYPOINT ["python", "main.py"]
