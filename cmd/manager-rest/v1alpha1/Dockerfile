FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/manager-rest
RUN go build -o vizier-manager-rest ./v1alpha1

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/manager-rest/vizier-manager-rest /app/
ENTRYPOINT ["./vizier-manager-rest"]
