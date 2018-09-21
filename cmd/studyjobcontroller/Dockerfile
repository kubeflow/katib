# Build the manager binary
FROM golang:1.10 AS build-env

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy in the go src
ADD . /go/src/github.com/kubeflow/katib

# Install dependencies
RUN cd /go/src/github.com/kubeflow/katib && \
	dep ensure --vendor-only

WORKDIR /go/src/github.com/kubeflow/katib/cmd/studyjobcontroller

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o studyjobcontroller
# Copy the controller-manager into a thin image
FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/studyjobcontroller/studyjobcontroller .
ENTRYPOINT ["./studyjobcontroller"]
