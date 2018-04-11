FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/hp-tuning
WORKDIR /go/src/github.com/kubeflow/hp-tuning/manager
RUN go build -o vizier-manager

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/hp-tuning/manager/vizier-manager /app/
COPY --from=build-env /go/src/github.com/kubeflow/hp-tuning/manager/visualise /
ENTRYPOINT ["./vizier-manager"]
CMD ["-w", "dlk"]
