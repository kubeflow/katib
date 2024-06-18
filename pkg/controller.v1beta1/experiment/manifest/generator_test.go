/*
Copyright 2022 The Kubeflow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package manifest

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	katibclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/util/katibclient"
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
							Image: "docker.io/kubeflowkatib/pytorch-mnist-cpu",
							Command: []string{
								"python3",
								"/opt/pytorch-mnist/mnist.py",
								"--epochs=1",
								"--batch-size=16",
								"--lr=0.05",
								"--momentum=0.9",
							},
							Env: []v1.EnvVar{
								{Name: consts.TrialTemplateMetaKeyOfName, Value: "trial-name"},
								{Name: consts.TrialTemplateMetaKeyOfNamespace, Value: "trial-namespace"},
								{Name: consts.TrialTemplateMetaKeyOfKind, Value: "Job"},
								{Name: consts.TrialTemplateMetaKeyOfAPIVersion, Value: "batch/v1"},
							},
						},
					},
				},
			},
		},
	}

	expectedRunSpec, err := util.ConvertObjectToUnstructured(expectedJob)
	if err != nil {
		t.Errorf("ConvertObjectToUnstructured failed: %v", err)
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
		// Parameter from assignments not found in Trial parameters
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
          image: docker.io/kubeflowkatib/pytorch-mnist-cpu
          command:
            - "python3"
            - "/opt/pytorch-mnist/mnist.py"
            - "--epochs=1"
            - "--batch-size=16"
            - "--lr=${trialParameters.learningRate}"
            - "--momentum=${trialParameters.momentum}"`

	invalidTrialSpec := `apiVersion: batch/v1
kind: Job
spec:
  template:
    spec:
      containers:
        - name: training-container
          image: docker.io/kubeflowkatib/pytorch-mnist-cpu
          command:
            - python3
            - /opt/pytorch-mnist/mnist.py
            - --epochs=1
            - --batch-size=16
            - --lr=${trialParameters.learningRate}
            - --momentum=${trialParameters.momentum}
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
          image: docker.io/kubeflowkatib/pytorch-mnist-cpu
          command:
            - "python3"
            - "/opt/pytorch-mnist/mnist.py"
            - "--epochs=1"
            - "--batch-size=16"
            - "--lr=0.05"
            - "--momentum=0.9"`

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
		// validGetConfigMap1 case
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
		// invalidConfigMapName case
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
		// validGetConfigMap3 case
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
							Image: "docker.io/kubeflowkatib/pytorch-mnist-cpu",
							Command: []string{
								"python3",
								"/opt/pytorch-mnist/mnist.py",
								"--epochs=1",
								"--batch-size=16",
								"--lr=${trialParameters.learningRate}",
								"--momentum=${trialParameters.momentum}",
							},
							Env: []v1.EnvVar{
								{Name: consts.TrialTemplateMetaKeyOfName, Value: "${trialParameters.trialName}"},
								{Name: consts.TrialTemplateMetaKeyOfNamespace, Value: "${trialParameters.trialNamespace}"},
								{Name: consts.TrialTemplateMetaKeyOfKind, Value: "${trialParameters.jobKind}"},
								{Name: consts.TrialTemplateMetaKeyOfAPIVersion, Value: "${trialParameters.jobAPIVersion}"},
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
						Name:        "momentum",
						Description: "Momentum for the training model",
						Reference:   "momentum",
					},
					{
						Name:        "trialName",
						Description: "name of current trial",
						Reference:   "${trialSpec.Name}",
					},
					{
						Name:        "trialNamespace",
						Description: "namespace of current trial",
						Reference:   "${trialSpec.Namespace}",
					},
					{
						Name:        "jobKind",
						Description: "job kind of current trial",
						Reference:   "${trialSpec.Kind}",
					},
					{
						Name:        "jobAPIVersion",
						Description: "job API Version of current trial",
						Reference:   "${trialSpec.APIVersion}",
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
			Name:  "momentum",
			Value: "0.9",
		},
	}
}
