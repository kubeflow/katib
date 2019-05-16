FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/manager-rest
RUN go build -o katib-manager-rest ./v1alpha2

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/manager-rest/katib-manager-rest /app/
ENTRYPOINT ["./katib-manager-rest"]
