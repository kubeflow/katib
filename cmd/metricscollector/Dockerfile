# Build the manager binary
FROM golang:alpine AS build-env

# Copy in the go src
ADD . /go/src/github.com/kubeflow/katib

WORKDIR /go/src/github.com/kubeflow/katib/cmd/metricscollector
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o metricscollector
# Copy the controller-manager into a thin image
FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/metricscollector/metricscollector .
