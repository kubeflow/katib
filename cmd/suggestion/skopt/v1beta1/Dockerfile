FROM python:3.10-slim

ARG TARGETARCH
ENV TARGET_DIR /opt/katib
ENV SUGGESTION_DIR cmd/suggestion/skopt/v1beta1
ENV PYTHONPATH ${TARGET_DIR}:${TARGET_DIR}/pkg/apis/manager/v1beta1/python:${TARGET_DIR}/pkg/apis/manager/health/python
# The GP surrogate works on at most a few hundred observations, so BLAS threading
# brings no speedup. Without this, OpenBLAS spawns one busy-waiting thread per host
# core, which is throttled by the Suggestion CPU limit and slows down GetSuggestions.
ENV OPENBLAS_NUM_THREADS 1
ENV OMP_NUM_THREADS 1
ENV MKL_NUM_THREADS 1

RUN if [ "${TARGETARCH}" = "ppc64le" ]; then \
    apt-get -y update && \
    apt-get -y install gfortran libopenblas-dev liblapack-dev && \
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
