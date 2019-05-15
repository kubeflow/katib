FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/manager
RUN go build -o vizier-manager ./v1alpha1

FROM alpine:3.7
WORKDIR /app
RUN GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/manager/vizier-manager /app/
#COPY --from=build-env /go/src/github.com/kubeflow/katib/pkg/manager/v1alpha1/visualise /
ENTRYPOINT ["./vizier-manager"]
CMD ["-w", "kubernetes"]
