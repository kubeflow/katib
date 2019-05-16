FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/ui
RUN go build -o katib-ui ./v1alpha1

FROM alpine:3.7
WORKDIR /app
# v1alpha1 source code
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/ui/katib-ui /app/
COPY cmd/ui/v1alpha1/static /app/static
COPY cmd/ui/v1alpha1/template /app/template

ENTRYPOINT ["./katib-ui"]