FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/hp-tuning
WORKDIR /go/src/github.com/kubeflow/hp-tuning/suggestion/random
RUN go build -o random

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/hp-tuning/suggestion/random /app/
ENTRYPOINT ["./random"]
