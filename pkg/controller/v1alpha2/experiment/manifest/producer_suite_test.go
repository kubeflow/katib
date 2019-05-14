package manifest

import (
	stdlog "log"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/kubeflow/katib/pkg/api/operators/apis"
)

var cfg *rest.Config

func TestMain(m *testing.M) {
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "..", "config", "crds", "v1alpha2"),
		},
	}
	stdlog.Println("Start adding apis")
	if err := apis.AddToScheme(scheme.Scheme); err != nil {
		stdlog.Fatal(err)
	}

	var err error
	stdlog.Println("Start server")
	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}
	stdlog.Println("Start code")
	code := m.Run()
	t.Stop()
	os.Exit(code)
}
