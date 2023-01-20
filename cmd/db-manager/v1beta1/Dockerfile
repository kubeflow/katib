# Build the Katib DB manager.
FROM golang:alpine AS build-env

ARG TARGETARCH
ENV GRPC_HEALTH_PROBE_VERSION v0.4.15

WORKDIR /go/src/github.com/kubeflow/katib

# Download packages.
COPY go.mod .
COPY go.sum .
RUN go mod download -x

# Copy sources.
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH="${TARGETARCH}" go build -a -o katib-db-manager ./cmd/db-manager/v1beta1

# Add GRPC health probe.
RUN wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-${TARGETARCH} \
  && chmod +x /bin/grpc_health_probe

# Copy the db-manager into a thin image.
FROM alpine:3.15
WORKDIR /app
COPY --from=build-env /bin/grpc_health_probe /bin/
COPY --from=build-env /go/src/github.com/kubeflow/katib/katib-db-manager /app/
ENTRYPOINT ["./katib-db-manager"]
