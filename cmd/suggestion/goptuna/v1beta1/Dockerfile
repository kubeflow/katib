# Build the Goptuna Suggestion.
FROM golang:alpine AS build-env

ARG TARGETARCH
ENV GRPC_HEALTH_PROBE_VERSION v0.4.11

WORKDIR /go/src/github.com/kubeflow/katib

# Download packages.
COPY go.mod .
COPY go.sum .
RUN go mod download -x

# Copy sources.
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -o goptuna-suggestion ./cmd/suggestion/goptuna/v1beta1

# Add GRPC health probe.
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-${TARGETARCH} \
  && chmod +x /bin/grpc_health_probe

# Copy the Goptuna suggestion into a thin image.
FROM alpine:3.15

ENV TARGET_DIR /opt/katib

WORKDIR ${TARGET_DIR}

COPY --from=build-env /bin/grpc_health_probe /bin/
COPY --from=build-env /go/src/github.com/kubeflow/katib/goptuna-suggestion ${TARGET_DIR}/

RUN chgrp -R 0 ${TARGET_DIR} \
  && chmod -R g+rwX ${TARGET_DIR}

ENTRYPOINT ["./goptuna-suggestion"]
