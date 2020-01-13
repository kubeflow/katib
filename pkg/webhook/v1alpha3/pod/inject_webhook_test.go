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
		Name          string
	}{
		{
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "tensorflow",
							Command: []string{
								"python main.py",
							},
						},
					},
				},
			},
			Namespace:   "nohere",
			JobKind:     "TFJob",
			MetricsFile: "testfile",
			PathKind:    common.FileKind,
			MC: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			Expected: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "tensorflow",
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								"python main.py 1>testfile 2>&1 && echo completed > $$$$.pid",
							},
						},
					},
				},
			},
			ExpectedError: nil,
			Name:          "tensorflow container without sh -c",
		},
		{
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "test",
							Command: []string{
								"python main.py",
							},
						},
					},
				},
			},
			Namespace:   "nohere",
			JobKind:     "TFJob",
			MetricsFile: "testfile",
			PathKind:    common.FileKind,
			MC: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			Expected: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "test",
							Command: []string{
								"python main.py",
							},
						},
					},
				},
			},
			ExpectedError: nil,
			Name:          "test container without sh -c",
		},
		{
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "tensorflow",
							Command: []string{
								"sh", "-c",
								"python main.py",
							},
						},
					},
				},
			},
			Namespace:   "nohere",
			JobKind:     "TFJob",
			MetricsFile: "testfile",
			PathKind:    common.FileKind,
			MC: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			Expected: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "tensorflow",
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								"python main.py 1>testfile 2>&1 && echo completed > $$$$.pid",
							},
						},
					},
				},
			},
			ExpectedError: nil,
			Name:          "Tensorflow container with sh -c",
		},
	}

	for _, c := range testCases {
		err := wrapWorkerContainer(c.Pod, c.Namespace, c.JobKind, c.MetricsFile, c.PathKind, c.MC)
		if err != c.ExpectedError {
			t.Errorf("Expected error %v, got %v", c.ExpectedError, err)
		}
		if err == nil {
			if !equality.Semantic.DeepEqual(c.Pod.Spec.Containers, c.Expected.Spec.Containers) {
				t.Errorf("Case %s: Expected pod %v, got %v",
					c.Name, c.Expected.Spec.Containers, c.Pod.Spec.Containers)
			}
		}
	}
}
