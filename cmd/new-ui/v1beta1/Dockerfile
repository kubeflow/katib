# --- Clone the kubeflow/kubeflow code ---
FROM ubuntu AS fetch-kubeflow-kubeflow

RUN apt-get update && apt-get install git -y

WORKDIR /kf
RUN git clone https://github.com/kubeflow/kubeflow.git && \
    cd kubeflow && \
    git checkout a349284

# --- Build the frontend kubeflow library ---
FROM node:12 AS frontend-kubeflow-lib

WORKDIR /src

ARG LIB=/kf/kubeflow/components/crud-web-apps/common/frontend/kubeflow-common-lib
COPY --from=fetch-kubeflow-kubeflow $LIB/package*.json ./
RUN npm ci

COPY --from=fetch-kubeflow-kubeflow $LIB/ ./
RUN npm run build

# --- Build the frontend ---
FROM node:12 AS frontend

WORKDIR /src
COPY ./pkg/new-ui/v1beta1/frontend/package*.json ./
RUN npm ci

COPY ./pkg/new-ui/v1beta1/frontend/ .
COPY --from=frontend-kubeflow-lib /src/dist/kubeflow/ ./node_modules/kubeflow/

RUN npm run build:prod

# --- Build the backend ---
FROM golang:alpine AS go-build

# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/new-ui
RUN if [ "$(uname -m)" = "ppc64le" ] || [ "$(uname -m)" = "aarch64" ]; then \
    apk --update add gcc musl-dev  && \
    go build -o katib-ui ./v1beta1; \
    else \
    go build -o katib-ui ./v1beta1; \
    fi

# --- compose the web app ---
FROM alpine:3.7
WORKDIR /app
COPY --from=go-build /go/src/github.com/kubeflow/katib/cmd/new-ui/katib-ui /app/
COPY --from=frontend /src/dist/static /app/build/static/

USER 1000
ENTRYPOINT ["./katib-ui"]
