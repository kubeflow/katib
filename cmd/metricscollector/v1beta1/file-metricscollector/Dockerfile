# Build the Katib file metrics collector.
FROM golang:alpine AS build-env

ARG TARGETARCH

WORKDIR /go/src/github.com/kubeflow/katib

# Download packages.
COPY go.mod .
COPY go.sum .
RUN go mod download -x

# Copy sources.
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -o file-metricscollector ./cmd/metricscollector/v1beta1/file-metricscollector

# Copy the file metrics collector into a thin image.
FROM alpine:3.15
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/file-metricscollector .
ENTRYPOINT ["./file-metricscollector"]
