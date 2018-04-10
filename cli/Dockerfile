FROM golang:1.8.2 AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/hp-tuning
WORKDIR /go/src/github.com/kubeflow/hp-tuning/cli
RUN go build -o katib-cli

FROM centos:7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/hp-tuning/cli/katib-cli /app/
