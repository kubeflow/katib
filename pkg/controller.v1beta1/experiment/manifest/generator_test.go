package manifest

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	util "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	katibclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/util/katibclient"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestGetRunSpecWithHP(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := katibclientmock.NewMockClient(mockCtrl)

	p := &DefaultGenerator{
		client: c,
	}

	expectedJob := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trial-name",
			Namespace: "trial-namespace",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "training-container",
							Image: "docker.io/kubeflowkatib/mxnet-mnist",
							Command: []string{
								"python3",
								"/opt/mxnet-mnist/mnist.py",
								"--lr=0.05",
								"--num-layers=5",
							},
						},
					},
				},
			},
		},
	}

	expectedRunSpec, err := util.ConvertObjectToUnstructured(expectedJob)
	if err != nil {
		t.Errorf("ConvertStringToUnstructured failed: %v", err)
	}

	tcs := []struct {
		Instance            *experimentsv1beta1.Experiment
		ParameterAssignment []commonapiv1beta1.ParameterAssignment
		expectedRunSpec     *unstructured.Unstructured
		Err                 bool
	}{
		{
			Instance:            newFakeInstance(),
			ParameterAssignment: newFakeParameterAssignment(),
			expectedRunSpec:     expectedRunSpec,
			Err:                 true,
		},
	}

	for _, tc := range tcs {
		actualRunSpec, err := p.GetRunSpecWithHyperParameters(tc.Instance, "trial-name", "trial-namespace", tc.ParameterAssignment)

		if tc.Err && err == nil {
			t.Errorf("Expected err, got nil")
		} else if !tc.Err {
			if err != nil {
				t.Errorf("Expected nil, got %v", err)
			} else if !reflect.DeepEqual(tc.expectedRunSpec, actualRunSpec) {
				t.Errorf("Expected %v\n got %v", tc.expectedRunSpec.Object, actualRunSpec.Object)
			}
		}
	}
}

func TestGetRunSpecWithHPConfigMap(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := katibclientmock.NewMockClient(mockCtrl)

	p := &DefaultGenerator{
		client: c,
	}

	templatePath := "trial-template.yaml"

	trialSpec := `apiVersion: batch/v1
kind: Job
spec:
  template:
    spec:
      containers:
        - name: training-container
          image: docker.io/kubeflowkatib/mxnet-mnist
          command:
            - "python3"
            - "/opt/mxnet-mnist/mnist.py"
            - "--lr=${trialParameters.learningRate}"
            - "--num-layers=${trialParameters.numberLayers}"`

	c.EXPECT().GetConfigMap(gomock.Any(), gomock.Any()).Return(map[string]string{
		templatePath: trialSpec,
	}, nil)

	instance := newFakeInstance()
	instance.Spec.TrialTemplate.TrialSource.ConfigMap = &experimentsv1beta1.ConfigMapSource{
		TemplatePath: templatePath,
	}
	instance.Spec.TrialTemplate.TrialSource.TrialSpec = nil
	actual, err := p.GetRunSpecWithHyperParameters(instance, "trial-name", "trial-namespace", []commonapiv1beta1.ParameterAssignment{
		{
			Name:  "lr",
			Value: "0.05",
		},
		{
			Name:  "num-layers",
			Value: "5",
		},
	})
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	// We can't compare structures, because trialSpec is a string and creationTimestamp was not added
	expectedJob := `apiVersion: batch/v1
kind: Job
metadata:
  name: trial-name
  namespace: trial-namespace
spec:
  template:
    spec:
      containers:
        - name: training-container
          image: docker.io/kubeflowkatib/mxnet-mnist
          command:
            - "python3"
            - "/opt/mxnet-mnist/mnist.py"
            - "--lr=0.05"
            - "--num-layers=5"`

	expected, err := util.ConvertStringToUnstructured(expectedJob)
	if err != nil {
		t.Errorf("ConvertStringToUnstructured failed: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %s\n got %s", expected.Object, actual.Object)
	}
}

func newFakeInstance() *experimentsv1beta1.Experiment {

	trialTemplateJob := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "training-container",
							Image: "docker.io/kubeflowkatib/mxnet-mnist",
							Command: []string{
								"python3",
								"/opt/mxnet-mnist/mnist.py",
								"--lr=${trialParameters.learningRate}",
								"--num-layers=${trialParameters.numberLayers}",
							},
						},
					},
				},
			},
		},
	}
	trialSpec, _ := util.ConvertObjectToUnstructured(trialTemplateJob)

	return &experimentsv1beta1.Experiment{
		Spec: experimentsv1beta1.ExperimentSpec{
			TrialTemplate: &experimentsv1beta1.TrialTemplate{
				TrialSource: experimentsv1beta1.TrialSource{
					TrialSpec: trialSpec,
				},
				TrialParameters: []experimentsv1beta1.TrialParameterSpec{
					{
						Name:        "learningRate",
						Description: "Learning Rate",
						Reference:   "lr",
					},
					{
						Name:        "numberLayers",
						Description: "Number of layers",
						Reference:   "num-layers",
					},
				},
			},
		},
	}
}

func newFakeParameterAssignment() []commonapiv1beta1.ParameterAssignment {
	return []commonapiv1beta1.ParameterAssignment{

		{
			Name:  "lr",
			Value: "0.05",
		},
		{
			Name:  "num-layers",
			Value: "5",
		},
	}
}
