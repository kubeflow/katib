# Build the Katib UI.
FROM node:12.18.1 AS npm-build

# Build frontend.
ADD /pkg/ui/v1beta1/frontend /frontend
RUN cd /frontend && npm ci
RUN cd /frontend && npm run build
RUN rm -rf /frontend/node_modules

# Build backend.
FROM golang:alpine AS go-build

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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -o katib-ui ./cmd/ui/v1beta1

# Copy the backend and frontend into a thin image.
FROM alpine:3.15
WORKDIR /app
COPY --from=go-build /go/src/github.com/kubeflow/katib/katib-ui /app/
COPY --from=npm-build /frontend/build /app/build
ENTRYPOINT ["./katib-ui"]
