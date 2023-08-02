# Build the Katib DB manager.
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
RUN CGO_ENABLED=0 GOOS=linux GOARCH="${TARGETARCH}" go build -a -o katib-db-manager ./cmd/db-manager/v1beta1

# Copy the db-manager into a thin image.
FROM alpine:3.15
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/katib-db-manager /app/
ENTRYPOINT ["./katib-db-manager"]
