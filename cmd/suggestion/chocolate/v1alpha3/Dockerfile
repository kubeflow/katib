FROM python:3.6

RUN if [ "$(uname -m)" = "ppc64le" ] || [ "$(uname -m)" = "aarch64" ]; then \
        apt-get -y update && \
        apt-get -y install gfortran libopenblas-dev liblapack-dev && \
        pip install cython 'numpy>=1.13.3'; \
    fi
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    if [ "$(uname -m)" = "ppc64le" ]; then \
	wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-ppc64le; \
    elif [ "$(uname -m)" = "aarch64" ]; then \
        wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-arm64; \
    else \
	wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64; \
    fi && \
    chmod +x /bin/grpc_health_probe

ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/cmd/suggestion/chocolate/v1alpha3
RUN pip install --no-cache-dir -r requirements.txt

ENV PYTHONPATH /usr/src/app/github.com/kubeflow/katib:/usr/src/app/github.com/kubeflow/katib/pkg/apis/manager/v1alpha3/python:/usr/src/app/github.com/kubeflow/katib/pkg/apis/manager/health/python

ENTRYPOINT ["python", "main.py"]
