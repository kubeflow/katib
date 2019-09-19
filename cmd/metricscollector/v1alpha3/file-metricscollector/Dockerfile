# Build the manager binary
ARG GOARCH=amd64
FROM golang:alpine AS build-env

# Copy in the go src
ADD . /go/src/github.com/kubeflow/katib

WORKDIR /go/src/github.com/kubeflow/katib/cmd/metricscollector/v1alpha3/file-metricscollector/

# Build
ARG CGO_ENABLED=0
RUN CGO_ENABLED=${CGO_ENABLED} GOOS=linux GOARCH=${GOARCH} go build -a -o file-metricscollector ./

# Copy the controller-manager into a thin image
FROM ${GOARCH}/alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/metricscollector/v1alpha3/file-metricscollector/file-metricscollector .
ENTRYPOINT ["./file-metricscollector"]
