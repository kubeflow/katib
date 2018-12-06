FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/ui
RUN go build -o katib-ui

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/ui/katib-ui /app/
COPY cmd/ui/static /app/static
COPY cmd/ui/template /app/template
ENTRYPOINT ["./katib-ui"]
