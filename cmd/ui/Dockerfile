FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/manager
RUN go build -o vizier-manager

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/manager/vizier-manager /app/
COPY --from=build-env /go/src/github.com/kubeflow/katib/pkg/manager/visualise /
ENTRYPOINT ["./vizier-manager"]
CMD ["-w", "kubernetes"]
