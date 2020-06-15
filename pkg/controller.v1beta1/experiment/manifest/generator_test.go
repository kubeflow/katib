package manifest

import (
	"errors"
	"math"
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
		Instance             *experimentsv1beta1.Experiment
		ParameterAssignments []commonapiv1beta1.ParameterAssignment
		expectedRunSpec      *unstructured.Unstructured
		Err                  bool
		testDescription      string
	}{
		// Valid run
		{
			Instance:             newFakeInstance(),
			ParameterAssignments: newFakeParameterAssignment(),
			expectedRunSpec:      expectedRunSpec,
			Err:                  false,
			testDescription:      "Run with valid parameters",
		},
		// Invalid JSON in unstructured
		{
			Instance: func() *experimentsv1beta1.Experiment {
				i := newFakeInstance()
				trialSpec := i.Spec.TrialTemplate.TrialSource.TrialSpec
				trialSpec.Object = map[string]interface{}{
					"invalidJSON": math.NaN(),
				}
				return i
			}(),
			ParameterAssignments: newFakeParameterAssignment(),
			Err:                  true,
			testDescription:      "Invalid JSON in Trial template",
		},
		// len(parameterAssignment) != len(trialParameters)
		{
			Instance: newFakeInstance(),
			ParameterAssignments: func() []commonapiv1beta1.ParameterAssignment {
				pa := newFakeParameterAssignment()
				pa = pa[1:]
				return pa
			}(),
			Err:             true,
			testDescription: "Number of parameter assignments is not equal to number of Trial parameters",
		},
		// Parameter from assignments not found in Trial paramters
		{
			Instance: newFakeInstance(),
			ParameterAssignments: func() []commonapiv1beta1.ParameterAssignment {
				pa := newFakeParameterAssignment()
				pa[0] = commonapiv1beta1.ParameterAssignment{
					Name:  "invalid-name",
					Value: "invalid-value",
				}
				return pa
			}(),
			Err:             true,
			testDescription: "Trial parameters don't have parameter from assignments",
		},
	}

	for _, tc := range tcs {
		actualRunSpec, err := p.GetRunSpecWithHyperParameters(tc.Instance, "trial-name", "trial-namespace", tc.ParameterAssignments)

		if tc.Err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		} else if !tc.Err {
			if err != nil {
				t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
			} else if !reflect.DeepEqual(tc.expectedRunSpec, actualRunSpec) {
				t.Errorf("Case: %v failed. Expected %v\n got %v", tc.testDescription, tc.expectedRunSpec.Object, actualRunSpec.Object)
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

	templatePath := "trial-template-path"

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

	invalidTrialSpec := `apiVersion: batch/v1
kind: Job
spec:
  template:
    spec:
      containers:
        - name: training-container
          image: docker.io/kubeflowkatib/mxnet-mnist
          command:
            - python3
            - /opt/mxnet-mnist/mnist.py
            - --lr=${trialParameters.learningRate}
            - --num-layers=${trialParameters.numberLayers}
            - --invalidParameter={'num_layers': 2, 'input_sizes': [32, 32, 3]}`

	validGetConfigMap1 := c.EXPECT().GetConfigMap(gomock.Any(), gomock.Any()).Return(
		map[string]string{templatePath: trialSpec}, nil)

	invalidConfigMapName := c.EXPECT().GetConfigMap(gomock.Any(), gomock.Any()).Return(
		nil, errors.New("Unable to get ConfigMap"))

	validGetConfigMap3 := c.EXPECT().GetConfigMap(gomock.Any(), gomock.Any()).Return(
		map[string]string{templatePath: trialSpec}, nil)

	invalidTemplate := c.EXPECT().GetConfigMap(gomock.Any(), gomock.Any()).Return(
		map[string]string{templatePath: invalidTrialSpec}, nil)

	gomock.InOrder(
		validGetConfigMap1,
		invalidConfigMapName,
		validGetConfigMap3,
		invalidTemplate,
	)

	// We can't compare structures, because in ConfigMap trialSpec is a string and creationTimestamp was not added
	expectedStr := `apiVersion: batch/v1
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

	expectedRunSpec, err := util.ConvertStringToUnstructured(expectedStr)
	if err != nil {
		t.Errorf("ConvertStringToUnstructured failed: %v", err)
	}

	tcs := []struct {
		Instance             *experimentsv1beta1.Experiment
		ParameterAssignments []commonapiv1beta1.ParameterAssignment
		Err                  bool
		testDescription      string
	}{
		// Valid run
		// validGetConfigMap case
		{
			Instance: func() *experimentsv1beta1.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate.TrialSource = experimentsv1beta1.TrialSource{
					ConfigMap: &experimentsv1beta1.ConfigMapSource{
						ConfigMapName:      "config-map-name",
						ConfigMapNamespace: "config-map-namespace",
						TemplatePath:       "trial-template-path",
					},
				}
				return i
			}(),
			ParameterAssignments: newFakeParameterAssignment(),
			Err:                  false,
			testDescription:      "Run with valid parameters",
		},
		// Invalid ConfigMap name
		{
			Instance: func() *experimentsv1beta1.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate.TrialSource = experimentsv1beta1.TrialSource{
					ConfigMap: &experimentsv1beta1.ConfigMapSource{
						ConfigMapName: "invalid-name",
					},
				}
				return i
			}(),
			ParameterAssignments: newFakeParameterAssignment(),
			Err:                  true,
			testDescription:      "Invalid ConfigMap name",
		},
		// Invalid template path in ConfigMap name
		{
			Instance: func() *experimentsv1beta1.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate.TrialSource = experimentsv1beta1.TrialSource{
					ConfigMap: &experimentsv1beta1.ConfigMapSource{
						ConfigMapName:      "config-map-name",
						ConfigMapNamespace: "config-map-namespace",
						TemplatePath:       "invalid-path",
					},
				}
				return i
			}(),
			ParameterAssignments: newFakeParameterAssignment(),
			Err:                  true,
			testDescription:      "Invalid template path in ConfigMap",
		},
		// Invalid Trial template spec in ConfigMap
		// Trial template is a string in ConfigMap
		// Because of that, user can specify not valid unstructured template
		// invalidTemplate case
		{
			Instance: func() *experimentsv1beta1.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate.TrialSource = experimentsv1beta1.TrialSource{
					ConfigMap: &experimentsv1beta1.ConfigMapSource{
						ConfigMapName:      "config-map-name",
						ConfigMapNamespace: "config-map-namespace",
						TemplatePath:       "trial-template-path",
					},
				}
				return i
			}(),
			ParameterAssignments: newFakeParameterAssignment(),
			Err:                  true,
			testDescription:      "Invalid Trial spec in ConfigMap",
		},
	}

	for _, tc := range tcs {
		actualRunSpec, err := p.GetRunSpecWithHyperParameters(tc.Instance, "trial-name", "trial-namespace", tc.ParameterAssignments)
		if tc.Err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		} else if !tc.Err {
			if err != nil {
				t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
			} else if !reflect.DeepEqual(expectedRunSpec, actualRunSpec) {
				t.Errorf("Case: %v failed. Expected %v\n got %v", tc.testDescription, expectedRunSpec.Object, actualRunSpec.Object)
			}
		}
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
