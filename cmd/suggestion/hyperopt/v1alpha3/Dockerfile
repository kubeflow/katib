ARG ARCH=amd64
FROM multiarch/ubuntu-debootstrap:${ARCH}-bionic

RUN apt-get -y update && \
    apt-get -y install software-properties-common && \
    add-apt-repository universe && \
    apt-get -y update && \
    apt-get -y install apt-utils && \
    apt-get -y install python3 && \
    apt-get -y install python3-pip && \
    apt-get -y install libxml2-dev libxslt1-dev python-lxml

RUN if [ "${ARCH}" = "ppc64le" ] || [ "${ARCH}" = "ppc64el" ]; then \
        apt-get -y update && \
        apt-get -y install gfortran libblas-dev liblapack-dev && \
        pip3 install cython lxml; \
    fi

RUN if [ "${ARCH}" = "ppc64le" ] || [ "${ARCH}" = "ppc64el" ]; then \
        GRPC_ARCH="ppc64le"; \
    else \
        GRPC_ARCH="amd64"; \
    fi && \
    GRPC_HEALTH_PROBE_VERSION=v0.3.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-${GRPC_ARCH} --no-check-certificate && \
    chmod +x /bin/grpc_health_probe

ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/cmd/suggestion/hyperopt/v1alpha3
RUN pip3 install --no-cache-dir -r requirements.txt

ENV PYTHONPATH /usr/src/app/github.com/kubeflow/katib:/usr/src/app/github.com/kubeflow/katib/pkg/apis/manager/v1alpha3/python:/usr/src/app/github.com/kubeflow/katib/pkg/apis/manager/health/python

ENTRYPOINT ["python", "main.py"]
