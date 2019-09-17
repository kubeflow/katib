package manifest

import (
	"reflect"
	"testing"
	"text/template"

	"github.com/golang/mock/gomock"

	commonapiv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	katibclientmock "github.com/kubeflow/katib/pkg/mock/v1alpha3/util/katibclient"
)

const (
	rawTemplate = `apiVersion: batch/v1
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
					{{- end}}`
)

func TestGetTrialTemplateConfigMap(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := katibclientmock.NewMockClient(mockCtrl)

	p := &DefaultGenerator{
		client: c,
	}

	templatePath := "test.yaml"

	c.EXPECT().GetConfigMap(gomock.Any(), gomock.Any()).Return(map[string]string{
		templatePath: rawTemplate,
	}, nil)

	instance := newFakeInstance()
	instance.Spec.TrialTemplate.GoTemplate.TemplateSpec = &experimentsv1alpha3.TemplateSpec{
		TemplatePath: templatePath,
	}
	instance.Spec.TrialTemplate.GoTemplate.RawTemplate = ""
	actual, err := p.getTrialTemplate(instance)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	expected, err := template.New("Trial").Parse(rawTemplate)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", *expected, *actual)
	}
}

func TestGetTrialTemplate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := katibclientmock.NewMockClient(mockCtrl)

	p := &DefaultGenerator{
		client: c,
	}

	tc := newFakeInstance()

	expected, err := template.New("Trial").
		Parse(tc.Spec.TrialTemplate.GoTemplate.RawTemplate)
	if err != nil {
		t.Errorf("Failed to compose expected result")
	}

	actual, err := p.getTrialTemplate(tc)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", *expected, *actual)
	}
}

func TestGetRunSpec(t *testing.T) {
	tc := newFakeInstance()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := katibclientmock.NewMockClient(mockCtrl)

	p := &DefaultGenerator{
		client: c,
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

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := katibclientmock.NewMockClient(mockCtrl)

	p := &DefaultGenerator{
		client: c,
	}

	actual, err := p.GetRunSpecWithHyperParameters(tc, "", "test", "testns", []commonapiv1alpha3.ParameterAssignment{
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

func newFakeInstance() *experimentsv1alpha3.Experiment {
	return &experimentsv1alpha3.Experiment{
		Spec: experimentsv1alpha3.ExperimentSpec{
			TrialTemplate: &experimentsv1alpha3.TrialTemplate{
				GoTemplate: &experimentsv1alpha3.GoTemplate{
					RawTemplate: rawTemplate,
				},
			},
		},
	}
}
