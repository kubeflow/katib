ARG GOARCH=amd64
FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/manager
RUN if [ "${GOARCH}" = "ppc64le" ]; then \
        apk --update add git gcc musl-dev && \
        go build -o katib-manager ./v1alpha3  && \
        go get github.com/grpc-ecosystem/grpc-health-probe && \
        mv $GOPATH/bin/grpc-health-probe /bin/grpc_health_probe && \
        chmod +x /bin/grpc_health_probe; \
    else \
        go build -o katib-manager ./v1alpha3 && \
        GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
            wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
            chmod +x /bin/grpc_health_probe; \
    fi

FROM ${GOARCH}/alpine:3.7
WORKDIR /app
COPY --from=build-env /bin/grpc_health_probe /bin/
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/manager/katib-manager /app/
ENTRYPOINT ["./katib-manager"]
CMD ["-w", "kubernetes"]
