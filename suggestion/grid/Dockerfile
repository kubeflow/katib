FROM golang
RUN : && \
    go get google.golang.org/grpc && \
    :
ADD api $GOPATH/src/github.com/mlkube/katib/api
ADD db $GOPATH/src/github.com/mlkube/katib/db
ADD manager $GOPATH/src/github.com/mlkube/katib/manager
ADD suggestion $GOPATH/src/github.com/mlkube/katib/suggestion
WORKDIR $GOPATH/src/github.com/mlkube/katib/suggestion/grid
RUN go build -o grid
