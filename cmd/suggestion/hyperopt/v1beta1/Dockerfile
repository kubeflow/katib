FROM python:3.10-slim

ARG TARGETARCH
ENV TARGET_DIR /opt/katib
ENV SUGGESTION_DIR cmd/suggestion/hyperopt/v1beta1
ENV PYTHONPATH ${TARGET_DIR}:${TARGET_DIR}/pkg/apis/manager/v1beta1/python:${TARGET_DIR}/pkg/apis/manager/health/python

RUN if [ "${TARGETARCH}" = "ppc64le" ]; then \
    apt-get -y update && \
    apt-get -y install gfortran libopenblas-dev liblapack-dev libssl-dev gcc g++ pkg-config git && \
    export REPO_ROOT=grpc && git clone -b v1.60.1 https://github.com/grpc/grpc $REPO_ROOT && \
    cd $REPO_ROOT && git submodule update --init && \
    pip install -r requirements.txt && GRPC_PYTHON_BUILD_SYSTEM_OPENSSL=1 pip install . && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*; \
    fi

ADD ./pkg/ ${TARGET_DIR}/pkg/
ADD ./${SUGGESTION_DIR}/ ${TARGET_DIR}/${SUGGESTION_DIR}/

WORKDIR  ${TARGET_DIR}/${SUGGESTION_DIR}

RUN pip install --prefer-binary --no-cache-dir -r requirements.txt
RUN chgrp -R 0 ${TARGET_DIR} \
  && chmod -R g+rwX ${TARGET_DIR}

ENTRYPOINT ["python", "main.py"]
