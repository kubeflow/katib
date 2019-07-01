FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/ui
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
        apk --update add gcc musl-dev  && \
	go build -o katib-ui ./v1alpha2; \
    else \
	go build -o katib-ui ./v1alpha2; \
    fi

FROM alpine:3.7
WORKDIR /app
# v1alpha2 source code
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/ui/katib-ui /app/
COPY --from=build-env /go/src/github.com/kubeflow/katib/pkg/ui/v1alpha2/frontend/build/ /app/build

ENTRYPOINT ["./katib-ui"]
