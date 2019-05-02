FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/katib
WORKDIR /go/src/github.com/kubeflow/katib/cmd/ui
RUN go build -o katib-ui ./v1alpha1
RUN go build -o katib-ui.v1alpha2 ./v1alpha2

FROM alpine:3.7
WORKDIR /app
# v1alpha1 source code
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/ui/katib-ui /app/
COPY cmd/ui/v1alpha1/static /app/static
COPY cmd/ui/v1alpha1/template /app/template
# v1alpha2 source code
COPY --from=build-env /go/src/github.com/kubeflow/katib/cmd/ui/katib-ui.v1alpha2 /app/
COPY --from=build-env /go/src/github.com/kubeflow/katib/pkg/ui/v1alpha2/frontend/build/ /app/build
