FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
# FOR NOW IT IS NEW-UI
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/new-ui
RUN go build -o katib-ui

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/new-ui/katib-ui /app/
COPY --from=build-env /go/src/github.com/kubeflow/katib/pkg/ui/frontend/build/ /app/build
ENTRYPOINT ["/app/katib-ui"]
EXPOSE 80
