package pod

import (
	"testing"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
)

func TestWrapWorkerContainer(t *testing.T) {
	testCases := []struct {
		Pod           *v1.Pod
		Namespace     string
		JobKind       string
		MetricsFile   string
		PathKind      common.FileSystemKind
		MC            common.MetricsCollectorSpec
		Expected      *v1.Pod
		ExpectedError error
	}{}

	for _, c := range testCases {
		err := wrapWorkerContainer(c.Pod, c.Namespace, c.JobKind, c.MetricsFile, c.PathKind, c.MC)
		if err != c.ExpectedError {
			t.Errorf("Expected error %v, got %v", c.ExpectedError, err)
		}
		if err == nil {
			if !equality.Semantic.DeepEqual(c.Pod.Spec, c.Expected.Spec) {
				t.Errorf("Expected pod %v, got %v", c.Pod.Spec, c.Expected.Spec)
			}
		}
	}
}
