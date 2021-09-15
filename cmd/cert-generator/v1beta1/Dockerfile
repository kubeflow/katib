# Build the Katib Cert Generatoe.
FROM golang:alpine AS build-env

WORKDIR /go/src/github.com/kubeflow/katib

# Download packages.
COPY go.mod .
COPY go.sum .
RUN go mod download -x

# Copy sources.
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build the binary.
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
    CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -a -o katib-cert-generator ./cmd/cert-generator/v1beta1; \
    elif [ "$(uname -m)" = "aarch64" ]; then \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o katib-cert-generator ./cmd/cert-generator/v1beta1; \
    else \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o katib-cert-generator ./cmd/cert-generator/v1beta1; \
    fi

# Copy the cert-generator into a thin image.
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/katib-cert-generator /app/
USER 65532:65532
ENTRYPOINT ["./katib-cert-generator"]
