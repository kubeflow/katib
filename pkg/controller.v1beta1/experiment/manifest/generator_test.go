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
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	katibclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/util/katibclient"
)

func TestGetRunSpecWithHP(t *testing.T) {

	scheme := runtime.NewScheme()
	v1.AddToScheme(scheme)
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	katibClient := katibclientmock.NewFakeClient(fakeClient)
	p := &DefaultGenerator{
		client: katibClient,
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
							Image: "ghcr.io/kubeflow/katib/pytorch-mnist-cpu",
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

	cases := map[string]struct {
		instance                       *experimentsv1beta1.Experiment
		parameterAssignments           []commonapiv1beta1.ParameterAssignment
		wantRunSpecWithHyperParameters *unstructured.Unstructured
		wantError                      error
	}{
		"Run with valid parameters": {
			instance:                       newFakeInstance(),
			parameterAssignments:           newFakeParameterAssignment(),
			wantRunSpecWithHyperParameters: expectedRunSpec,
		},
		"Invalid JSON in Unstructured Trial template": {
			instance: func() *experimentsv1beta1.Experiment {
				i := newFakeInstance()
				trialSpec := i.Spec.TrialTemplate.TrialSource.TrialSpec
				trialSpec.Object = map[string]interface{}{
					"invalidJSON": math.NaN(),
				}
				return i
			}(),
			parameterAssignments: newFakeParameterAssignment(),
			wantError:            errConvertUnstructuredToStringFailed,
		},
		"Non-meta parameter from TrialParameters not found in ParameterAssignment": {
			instance: newFakeInstance(),
			parameterAssignments: func() []commonapiv1beta1.ParameterAssignment {
				pa := newFakeParameterAssignment()
				pa[0] = commonapiv1beta1.ParameterAssignment{
					Name:  "invalid-name",
					Value: "invalid-value",
				}
				return pa
			}(),
			wantError: errParamNotFoundInParameterAssignment,
		},
		// case in which the lengths of trial parameters and parameter assignments are different
		"Parameter from ParameterAssignment not found in TrialParameters": {
			instance: newFakeInstance(),
			parameterAssignments: func() []commonapiv1beta1.ParameterAssignment {
				pa := newFakeParameterAssignment()
				pa = append(pa, commonapiv1beta1.ParameterAssignment{
					Name:  "extra-name",
					Value: "extra-value",
				})
				return pa
			}(),
			wantError: errParamNotFoundInTrialParameters,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := p.GetRunSpecWithHyperParameters(tc.instance, "trial-name", "trial-namespace", tc.parameterAssignments)
			if diff := cmp.Diff(tc.wantError, err, cmpopts.EquateErrors()); len(diff) != 0 {
				t.Errorf("Unexpected error from GetRunSpecWithHyperParameters (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantRunSpecWithHyperParameters, got); len(diff) != 0 {
				t.Errorf("Unexpected run spec from GetRunSpecWithHyperParameters (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestGetRunSpecWithHPConfigMap(t *testing.T) {
	// Mocking the ConfigMap
	templatePath := "trial-template-path"

	trialSpec := `apiVersion: batch/v1
kind: Job
spec:
  template:
    spec:
      containers:
        - name: training-container
          image: ghcr.io/kubeflow/katib/pytorch-mnist-cpu
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
          image: ghcr.io/kubeflow/katib/pytorch-mnist-cpu
          command:
            - python3
            - /opt/pytorch-mnist/mnist.py
            - --epochs=1
            - --batch-size=16
            - --lr=${trialParameters.learningRate}
            - --momentum=${trialParameters.momentum}
            - --invalidParameter={'num_layers': 2, 'input_sizes': [32, 32, 3]}`

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
          image: ghcr.io/kubeflow/katib/pytorch-mnist-cpu
          command:
            - "python3"
            - "/opt/pytorch-mnist/mnist.py"
            - "--epochs=1"
            - "--batch-size=16"
            - "--lr=0.05"
            - "--momentum=0.9"`

	expectedRunSpec, err := util.ConvertStringToUnstructured(expectedStr)
	if diff := cmp.Diff(nil, err, cmpopts.EquateErrors()); len(diff) != 0 {
		t.Errorf("ConvertStringToUnstructured failed (-want,+got):\n%s", diff)
	}

	cases := map[string]struct {
		objects                        []runtime.Object
		instance                       *experimentsv1beta1.Experiment
		parameterAssignments           []commonapiv1beta1.ParameterAssignment
		wantRunSpecWithHyperParameters *unstructured.Unstructured
		wantError                      error
	}{
		"Run with valid parameters": {
			objects: []runtime.Object{
				&v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "config-map-name",
						Namespace: "config-map-namespace",
					},
					Data: map[string]string{
						templatePath: trialSpec,
					},
				},
			},
			instance: func() *experimentsv1beta1.Experiment {
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
			parameterAssignments:           newFakeParameterAssignment(),
			wantRunSpecWithHyperParameters: expectedRunSpec,
		},
		"Invalid ConfigMap name": {
			objects: []runtime.Object{
				&v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "config-map-name",
						Namespace: "config-map-namespace",
					},
					Data: map[string]string{
						templatePath: trialSpec,
					},
				},
			},
			instance: func() *experimentsv1beta1.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate.TrialSource = experimentsv1beta1.TrialSource{
					ConfigMap: &experimentsv1beta1.ConfigMapSource{
						ConfigMapName: "invalid-name",
					},
				}
				return i
			}(),
			parameterAssignments: newFakeParameterAssignment(),
			wantError:            errConfigMapNotFound,
		},
		"Invalid template path in ConfigMap name": {
			objects: []runtime.Object{
				&v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "config-map-name",
						Namespace: "config-map-namespace",
					},
					Data: map[string]string{
						templatePath: trialSpec, //No templatePath
					},
				},
			},
			instance: func() *experimentsv1beta1.Experiment {
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
			parameterAssignments: newFakeParameterAssignment(),
			wantError:            errTrialTemplateNotFound,
		},
		// Trial template is a string in ConfigMap
		// Because of that, user can specify not valid unstructured template
		"Invalid trial spec in ConfigMap": {
			objects: []runtime.Object{
				&v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "config-map-name",
						Namespace: "config-map-namespace",
					},
					Data: map[string]string{
						templatePath: invalidTrialSpec,
					},
				},
			},
			instance: func() *experimentsv1beta1.Experiment {
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
			parameterAssignments: newFakeParameterAssignment(),
			wantError:            errConvertStringToUnstructuredFailed,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			v1.AddToScheme(scheme)
			fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(tc.objects...).Build()
			katibClient := katibclientmock.NewFakeClient(fakeClient)
			p := &DefaultGenerator{
				client: katibClient,
			}
			got, err := p.GetRunSpecWithHyperParameters(tc.instance, "trial-name", "trial-namespace", tc.parameterAssignments)
			if diff := cmp.Diff(tc.wantError, err, cmpopts.EquateErrors()); len(diff) != 0 {
				t.Errorf("Unexpected error from GetRunSpecWithHyperParameters (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantRunSpecWithHyperParameters, got); len(diff) != 0 {
				t.Errorf("Unexpected run spec from GetRunSpecWithHyperParameters (-want,+got):\n%s", diff)
			}
		})
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
							Image: "ghcr.io/kubeflow/katib/pytorch-mnist-cpu",
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
