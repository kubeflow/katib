FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/hp-tuning
WORKDIR /go/src/github.com/kubeflow/hp-tuning/suggestion/hyperband
RUN go build -o hyperband

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/hp-tuning/suggestion/hyperband /app/
ENTRYPOINT ["./hyperband"]
