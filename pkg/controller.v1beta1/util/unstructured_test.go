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

package util

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestInjectLabelsToPodTemplate(t *testing.T) {
	labels := map[string]string{
		"katib.kubeflow.org/trial":      "test-trial",
		"katib.kubeflow.org/experiment": "test-experiment",
	}

	cases := map[string]struct {
		obj      *unstructured.Unstructured
		expected *unstructured.Unstructured
	}{
		"Standard Job": {
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "Job",
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"labels": map[string]interface{}{
									"unaffected": "label",
								},
							},
						},
					},
				},
			},
			expected: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "Job",
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"katib.kubeflow.org/trial":      "test-trial",
							"katib.kubeflow.org/experiment": "test-experiment",
						},
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"labels": map[string]interface{}{
									"unaffected":                    "label",
									"katib.kubeflow.org/trial":      "test-trial",
									"katib.kubeflow.org/experiment": "test-experiment",
								},
							},
						},
					},
				},
			},
		},
		"CronJob": {
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "CronJob",
					"spec": map[string]interface{}{
						"jobTemplate": map[string]interface{}{
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"metadata": map[string]interface{}{
										"labels": map[string]interface{}{},
									},
								},
							},
						},
					},
				},
			},
			expected: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "CronJob",
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"katib.kubeflow.org/trial":      "test-trial",
							"katib.kubeflow.org/experiment": "test-experiment",
						},
					},
					"spec": map[string]interface{}{
						"jobTemplate": map[string]interface{}{
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"metadata": map[string]interface{}{
										"labels": map[string]interface{}{
											"katib.kubeflow.org/trial":      "test-trial",
											"katib.kubeflow.org/experiment": "test-experiment",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"TFJob (Training Operator)": {
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "kubeflow.org/v1",
					"kind":       "TFJob",
					"spec": map[string]interface{}{
						"replicaSpecs": map[string]interface{}{
							"Worker": map[string]interface{}{
								"template": map[string]interface{}{
									"metadata": map[string]interface{}{
										"labels": map[string]interface{}{},
									},
								},
							},
							"PS": map[string]interface{}{
								"template": map[string]interface{}{
									"metadata": map[string]interface{}{
										"labels": map[string]interface{}{},
									},
								},
							},
						},
					},
				},
			},
			expected: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "kubeflow.org/v1",
					"kind":       "TFJob",
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"katib.kubeflow.org/trial":      "test-trial",
							"katib.kubeflow.org/experiment": "test-experiment",
						},
					},
					"spec": map[string]interface{}{
						"replicaSpecs": map[string]interface{}{
							"Worker": map[string]interface{}{
								"template": map[string]interface{}{
									"metadata": map[string]interface{}{
										"labels": map[string]interface{}{
											"katib.kubeflow.org/trial":      "test-trial",
											"katib.kubeflow.org/experiment": "test-experiment",
										},
									},
								},
							},
							"PS": map[string]interface{}{
								"template": map[string]interface{}{
									"metadata": map[string]interface{}{
										"labels": map[string]interface{}{
											"katib.kubeflow.org/trial":      "test-trial",
											"katib.kubeflow.org/experiment": "test-experiment",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"Job with no initial labels": {
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "Job",
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{},
						},
					},
				},
			},
			expected: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "Job",
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"katib.kubeflow.org/trial":      "test-trial",
							"katib.kubeflow.org/experiment": "test-experiment",
						},
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"labels": map[string]interface{}{
									"katib.kubeflow.org/trial":      "test-trial",
									"katib.kubeflow.org/experiment": "test-experiment",
								},
							},
							"spec": map[string]interface{}{},
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := InjectLabelsToPodTemplate(tc.obj, labels)
			if err != nil {
				t.Fatalf("InjectLabelsToPodTemplate failed: %v", err)
			}

			gotLabels, _, _ := unstructured.NestedStringMap(tc.obj.Object, "metadata", "labels")
			for k, v := range labels {
				if gotLabels[k] != v {
					t.Errorf("%s: metadata.labels[%s] = %v, want %v", name, k, gotLabels[k], v)
				}
			}

			if tc.obj.GetKind() == "Job" {
				gotLabels, _, _ = unstructured.NestedStringMap(tc.obj.Object, "spec", "template", "metadata", "labels")
				for k, v := range labels {
					if gotLabels[k] != v {
						t.Errorf("%s: spec.template.metadata.labels[%s] = %v, want %v", name, k, gotLabels[k], v)
					}
				}
			}

			if tc.obj.GetKind() == "CronJob" {
				gotLabels, _, _ = unstructured.NestedStringMap(tc.obj.Object, "spec", "jobTemplate", "spec", "template", "metadata", "labels")
				for k, v := range labels {
					if gotLabels[k] != v {
						t.Errorf("%s: spec.jobTemplate.spec.template.metadata.labels[%s] = %v, want %v", name, k, gotLabels[k], v)
					}
				}
			}

			if tc.obj.GetKind() == "TFJob" {
				replicaSpecs, _, _ := unstructured.NestedMap(tc.obj.Object, "spec", "replicaSpecs")
				for rName, rSpec := range replicaSpecs {
					rs := rSpec.(map[string]interface{})
					gotLabels, _, _ = unstructured.NestedStringMap(rs, "template", "metadata", "labels")
					for k, v := range labels {
						if gotLabels[k] != v {
							t.Errorf("%s: replica %s labels[%s] = %v, want %v", name, rName, k, gotLabels[k], v)
						}
					}
				}
			}
		})
	}
}
