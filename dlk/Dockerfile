FROM golang:alpine AS build-env
# The GOPATH in the image is /go.
ADD . /go/src/github.com/kubeflow/hp-tuning
WORKDIR /go/src/github.com/kubeflow/hp-tuning/dlk/dlkmanager
RUN go build -o dlkmanager 

FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /go/src/github.com/kubeflow/hp-tuning/dlk/dlkmanager/dlkmanager /app/
ENTRYPOINT ["./dlkmanager"]
