# Build the manager binary
ARG GOARCH=amd64
FROM golang:alpine AS build-env

# Copy in the go src
ADD . /go/src/github.com/kubeflow/katib

WORKDIR /go/src/github.com/kubeflow/katib/cmd/katib-controller
# Build
ARG CGO_ENABLED=0
RUN CGO_ENABLED=${CGO_ENABLED} GOARCH=${GOARCH} go build -a -o katib-controller ./v1alpha3

# Copy the controller-manager into a thin image
FROM ${GOARCH}/alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/katib-controller/katib-controller .
ENTRYPOINT ["./katib-controller"]
