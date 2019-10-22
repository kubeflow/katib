FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
        apk --update add git gcc musl-dev && \
        go get github.com/grpc-ecosystem/grpc-health-probe && \
        mv $GOPATH/bin/grpc-health-probe /bin/grpc_health_probe && \
        chmod +x /bin/grpc_health_probe; \
    else \
        GRPC_HEALTH_PROBE_VERSION=v0.3.0 && \
            wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
            chmod +x /bin/grpc_health_probe; \
    fi

FROM python:3.6

COPY --from=build-env /bin/grpc_health_probe /bin/

RUN if [ "$(uname -m)" = "ppc64le" ]; then \
        apt-get -y update && \
        apt-get -y install gfortran libopenblas-dev liblapack-dev && \
        pip install cython; \
    fi

COPY --from=build-env /bin/grpc_health_probe /bin/

ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/cmd/suggestion/nasrl/v1alpha3
RUN pip install --no-cache-dir -r requirements.txt

ENV PYTHONPATH /usr/src/app/github.com/kubeflow/katib:/usr/src/app/github.com/kubeflow/katib/pkg/apis/manager/v1alpha3/python:/usr/src/app/github.com/kubeflow/katib/pkg/apis/manager/health/python

ENTRYPOINT ["python", "main.py"]
