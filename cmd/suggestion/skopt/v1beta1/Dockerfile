FROM alpine:3.15 AS downloader
ENV GRPC_HEALTH_PROBE_VERSION v0.4.11
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-ppc64le; \
    elif [ "$(uname -m)" = "aarch64" ]; then \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-arm64; \
    else \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64; \
    fi && \
    chmod +x /bin/grpc_health_probe

FROM python:3.9-slim
ENV TARGET_DIR /opt/katib
ENV SUGGESTION_DIR cmd/suggestion/skopt/v1beta1

RUN if [ "$(uname -m)" = "ppc64le" ] || [ "$(uname -m)" = "aarch64" ]; then \
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

ENV PYTHONPATH ${TARGET_DIR}:${TARGET_DIR}/pkg/apis/manager/v1beta1/python:${TARGET_DIR}/pkg/apis/manager/health/python

ENTRYPOINT ["python", "main.py"]
