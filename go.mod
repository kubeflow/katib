module github.com/kubeflow/katib

go 1.15

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/awalterschulze/gographviz v2.0.3+incompatible
	github.com/c-bata/goptuna v0.8.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-openapi/spec v0.19.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/go-containerregistry v0.4.1-0.20210128200529-19c2b639fab1
	github.com/google/go-containerregistry/pkg/authn/k8schain v0.0.0-20210224013640-6928f6d356ab
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/hpcloud/tail v1.0.1-0.20180514194441-a1dbeea552b7
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a
	github.com/onsi/gomega v1.10.3
	github.com/prometheus/client_golang v1.9.0
	github.com/shirou/gopsutil v2.20.7+incompatible
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/tidwall/gjson v1.6.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	google.golang.org/grpc v1.38.0
	gopkg.in/fsnotify/fsnotify.v1 v1.4.7 // indirect
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20210216185858-15cd8face8d6
	k8s.io/utils v0.0.0-20210707171843-4b05e18ac7d9
	sigs.k8s.io/controller-runtime v0.8.2
)

replace k8s.io/code-generator => k8s.io/code-generator v0.20.4
