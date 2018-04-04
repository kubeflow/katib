FROM golang
RUN : && \
    go get google.golang.org/grpc && \
    :
ADD api $GOPATH/src/github.com/kubeflow/hp-tuning/api
ADD db $GOPATH/src/github.com/kubeflow/hp-tuning/db
ADD manager $GOPATH/src/github.com/kubeflow/hp-tuning/manager
ADD suggestion $GOPATH/src/github.com/kubeflow/hp-tuning/suggestion
WORKDIR $GOPATH/src/github.com/kubeflow/hp-tuning/suggestion/hyperband
RUN go build -o hyperband
