FROM node:12.18.1 AS npm-build

ADD /pkg/ui/v1beta1/frontend /frontend
RUN cd /frontend && npm install
RUN cd /frontend && npm run build
RUN rm -rf /frontend/node_modules

FROM golang:alpine AS go-build
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/ui
RUN if [ "$(uname -m)" = "ppc64le" ] || [ "$(uname -m)" = "aarch64" ]; then \
    apk --update add gcc musl-dev  && \
    go build -o katib-ui ./v1beta1; \
    else \
    go build -o katib-ui ./v1beta1; \
    fi

FROM alpine:3.7
WORKDIR /app
# v1beta1 source code
COPY --from=go-build /go/src/github.com/kubeflow/katib/cmd/ui/katib-ui /app/
COPY --from=npm-build /frontend/build /app/build

USER 1000

ENTRYPOINT ["./katib-ui"]
