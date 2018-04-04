FROM golang
RUN : && \
    go get k8s.io/api/apps/v1beta1 && \
    go get k8s.io/api/core/v1 && \
    go get k8s.io/apimachinery/pkg/apis/meta/v1 && \
    go get k8s.io/client-go/kubernetes && \
    go get k8s.io/client-go/tools/clientcmd && \
    go get k8s.io/client-go/util/homedir && \
    go get k8s.io/client-go/util/retry && \
    go get k8s.io/client-go/rest && \
    go get k8s.io/apimachinery/pkg/runtime/serializer && \
    go get database/sql && \
    go get github.com/golang/protobuf/jsonpb && \
    go get github.com/go-sql-driver/mysql && \
    go get github.com/mattn/go-sqlite3 && \
    go get google.golang.org/grpc && \
    go get gopkg.in/yaml.v2 && \
    :
ADD api $GOPATH/src/github.com/mlkube/katib/api
ADD cli $GOPATH/src/github.com/mlkube/katib/cli
WORKDIR $GOPATH/src/github.com/mlkube/katib/cli
RUN go build -o katib-cli
