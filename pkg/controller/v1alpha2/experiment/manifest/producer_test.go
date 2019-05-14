package manifest

import (
	"reflect"
	"testing"
	"text/template"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	apiv1alpha2 "github.com/kubeflow/katib/pkg/api/v1alpha2"
	"github.com/kubeflow/katib/pkg/util/v1alpha2/katibclient"
)

func TestGetTrialTemplate(t *testing.T) {
	p := &General{
		client: katibclient.MustNewTestClient(cfg),
	}

	tc := newFakeInstance()

	expected, err := template.New("Trial").
		Parse(tc.Spec.TrialTemplate.GoTemplate.RawTemplate)
	if err != nil {
		t.Errorf("Failed to compose expected result")
	}

	actual, err := p.getTrialTemplate(tc)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", *expected, *actual)
	}
}

func TestGetRunSpec(t *testing.T) {
	tc := newFakeInstance()

	p := &General{
		client: katibclient.MustNewTestClient(cfg),
	}

	actual, err := p.GetRunSpec(tc, "", "test", "testns")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	expected := `apiVersion: batch/v1
kind: Job
metadata:
name: test
namespace: testns
spec:
	template:
		spec:
			containers:
				- name: test
				  image: katib/mxnet-mnist-example
				  command:
					- "python"
					- "/mxnet/example/image-classification/train_mnist.py"
					- "--batch-size=64"`
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestGetRunSpecWithHP(t *testing.T) {
	tc := newFakeInstance()

	p := &General{
		client: katibclient.MustNewTestClient(cfg),
	}

	actual, err := p.GetRunSpecWithHyperParameters(tc, "", "test", "testns", []*apiv1alpha2.ParameterAssignment{
		{
			Name:  "testname",
			Value: "testvalue",
		},
	})
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	expected := `apiVersion: batch/v1
kind: Job
metadata:
name: test
namespace: testns
spec:
	template:
		spec:
			containers:
				- name: test
				  image: katib/mxnet-mnist-example
				  command:
					- "python"
					- "/mxnet/example/image-classification/train_mnist.py"
					- "--batch-size=64"
					- "testname=testvalue"`
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func newFakeInstance() *experimentsv1alpha2.Experiment {
	return &experimentsv1alpha2.Experiment{
		Spec: experimentsv1alpha2.ExperimentSpec{
			TrialTemplate: &experimentsv1alpha2.TrialTemplate{
				GoTemplate: &experimentsv1alpha2.GoTemplate{
					RawTemplate: `apiVersion: batch/v1
kind: Job
metadata:
name: {{.Trial}}
namespace: {{.NameSpace}}
spec:
	template:
		spec:
			containers:
				- name: {{.Trial}}
				  image: katib/mxnet-mnist-example
				  command:
					- "python"
					- "/mxnet/example/image-classification/train_mnist.py"
					- "--batch-size=64"
					{{- with .HyperParameters}}
					{{- range .}}
					- "{{.Name}}={{.Value}}"
					{{- end}}
					{{- end}}`,
				},
			},
		},
	}
}
