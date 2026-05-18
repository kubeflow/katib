# --- Clone the kubeflow/notebooks code ---
FROM alpine/git AS fetch-kubeflow-notebooks

WORKDIR /kf
COPY ./pkg/ui/v1beta1/frontend/COMMIT ./
RUN git clone https://github.com/kubeflow/notebooks.git && \
    COMMIT=$(cat ./COMMIT) && \
    cd notebooks && \
    git checkout $COMMIT

# --- Build the frontend kubeflow library ---
FROM node:16-alpine AS frontend-kubeflow-lib

WORKDIR /src

ARG LIB=/kf/notebooks/components/crud-web-apps/common/frontend/kubeflow-common-lib
COPY --from=fetch-kubeflow-notebooks $LIB/package*.json ./
RUN npm config set fetch-retry-mintimeout 200000 && \
    npm config set fetch-retry-maxtimeout 1200000 && \
    npm config get registry && \
    npm config set registry https://registry.npmjs.org/ && \
    npm config delete https-proxy && \
    npm config set loglevel verbose && \
    npm cache clean --force && \
    npm ci --force --prefer-offline --no-audit

COPY --from=fetch-kubeflow-notebooks $LIB/ ./
RUN npm run build

# --- Build the frontend ---
FROM node:16-alpine AS frontend

WORKDIR /src
COPY ./pkg/ui/v1beta1/frontend/package*.json ./
RUN npm config set fetch-retry-mintimeout 200000 && \
    npm config set fetch-retry-maxtimeout 1200000 && \
    npm config get registry && \
    npm config set registry https://registry.npmjs.org/ && \
    npm config delete https-proxy && \
    npm config set loglevel verbose && \
    npm cache clean --force && \
    npm ci --force --prefer-offline --no-audit

COPY ./pkg/ui/v1beta1/frontend/ .
COPY --from=frontend-kubeflow-lib /src/dist/kubeflow/ ./node_modules/kubeflow/

RUN npm run build:prod

# --- Build the backend ---
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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -o katib-ui  ./cmd/ui/v1beta1

# --- Compose the web app ---
FROM alpine:3.15
WORKDIR /app
COPY --from=go-build /go/src/github.com/kubeflow/katib/katib-ui /app/
COPY --from=frontend /src/dist/static /app/build/static/
ENTRYPOINT ["./katib-ui"]
