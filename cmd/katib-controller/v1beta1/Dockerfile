# Build the Katib controller.
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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -o katib-controller ./cmd/katib-controller/v1beta1

# Copy the controller-manager into a thin image.
FROM alpine:3.15
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/katib-controller .
ENTRYPOINT ["./katib-controller"]
