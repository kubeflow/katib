FROM python:3.11-slim

ARG TARGETARCH
ENV TARGET_DIR /opt/katib
ENV SUGGESTION_DIR cmd/suggestion/pbt/v1beta1
ENV PYTHONPATH ${TARGET_DIR}:${TARGET_DIR}/pkg/apis/manager/v1beta1/python:${TARGET_DIR}/pkg/apis/manager/health/python

RUN if [ "${TARGETARCH}" = "ppc64le" ]; then \
    apt-get -y update && \
    apt-get -y install g++ gfortran libopenblas-dev liblapack-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*; \
    fi

ADD ./pkg/ ${TARGET_DIR}/pkg/
ADD ./${SUGGESTION_DIR}/ ${TARGET_DIR}/${SUGGESTION_DIR}/

WORKDIR  ${TARGET_DIR}/${SUGGESTION_DIR}

RUN if [ "${TARGETARCH}" = "ppc64le" ]; then \
      pip install --prefer-binary --no-cache-dir -r requirements.txt \
        --extra-index-url https://wheels.developerfirst.ibm.com/ppc64le/linux/+simple/; \
    else \
      pip install --prefer-binary --no-cache-dir -r requirements.txt; \
    fi
RUN chgrp -R 0 ${TARGET_DIR} \
  && chmod -R g+rwX ${TARGET_DIR}

ENTRYPOINT ["python", "main.py"]
