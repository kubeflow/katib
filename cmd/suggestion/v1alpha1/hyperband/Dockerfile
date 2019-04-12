FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/suggestion/v1alpha1/hyperband
RUN go build -o hyperband .

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/suggestion/v1alpha1/hyperband /app/
ENTRYPOINT ["./hyperband"]
