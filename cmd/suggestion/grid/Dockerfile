FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/suggestion/grid
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
        apk --update add gcc musl-dev  && \
        go build -o grid; \
    else \
        go build -o grid; \
    fi

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/suggestion/grid /app/
ENTRYPOINT ["./grid"]
