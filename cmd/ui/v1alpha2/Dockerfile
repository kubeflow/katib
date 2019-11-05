FROM node:10.16.2-alpine AS npm-build

ADD /pkg/ui/v1alpha2/frontend /frontend
RUN cd /frontend && npm install
RUN cd /frontend && npm run build
RUN rm -rf /frontend/node_modules

FROM golang:alpine AS go-build
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/ui
RUN if [ "$(uname -m)" = "ppc64le" ] || [ "$(uname -m)" = "aarch64" ]; then \
        apk --update add gcc musl-dev  && \
	go build -o katib-ui ./v1alpha2; \
    else \
	go build -o katib-ui ./v1alpha2; \
    fi

FROM alpine:3.7
WORKDIR /app
# v1alpha2 source code
COPY --from=go-build /go/src/github.com/kubeflow/katib/cmd/ui/katib-ui /app/
COPY --from=npm-build /frontend/build /app/build

ENTRYPOINT ["./katib-ui"]
