# Build the Katib Cert Generator.
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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -o katib-cert-generator ./cmd/cert-generator/v1beta1

# Copy the cert-generator into a thin image.
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/katib-cert-generator /app/
USER 65532:65532
ENTRYPOINT ["./katib-cert-generator"]
