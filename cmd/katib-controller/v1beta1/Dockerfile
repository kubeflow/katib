# Build the manager binary
FROM golang:alpine AS build-env

# Copy in the go src
ADD . /go/src/github.com/kubeflow/katib

WORKDIR /go/src/github.com/kubeflow/katib/cmd/katib-controller
# Build
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
    CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -a -o katib-controller  ./v1beta1; \
    elif [ "$(uname -m)" = "aarch64" ]; then \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o katib-controller  ./v1beta1; \
    else \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o katib-controller  ./v1beta1; \
    fi
# Copy the controller-manager into a thin image
FROM alpine:3.7
WORKDIR /app
RUN apk update && apk add ca-certificates
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/katib-controller/katib-controller .
USER 1000
ENTRYPOINT ["./katib-controller"]
