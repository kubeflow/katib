FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/earlystopping/medianstopping
RUN go build -o medianstopping

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/earlystopping/medianstopping /app/
ENTRYPOINT ["./medianstopping"]
