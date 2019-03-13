# Build the manager binary
FROM golang:alpine AS build-env

# Copy in the go src
ADD . /go/src/github.com/kubeflow/katib

WORKDIR /go/src/github.com/kubeflow/katib/cmd/studyjobcontroller
# Build
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
        CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -a -o studyjobcontroller; \
    else \
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o studyjobcontroller; \
    fi
# Copy the controller-manager into a thin image
FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/studyjobcontroller/studyjobcontroller .
ENTRYPOINT ["./studyjobcontroller"]
