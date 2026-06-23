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

package v1beta1

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
)

func TestSetDefaultMetricsCollector(t *testing.T) {
	tcs := []struct {
		name              string
		trialSpec         *unstructured.Unstructured
		explicitCollector *common.CollectorSpec
		wantCollectorKind common.CollectorKind
	}{
		{
			name:              "non-TrainJob trial defaults to StdOut collector",
			trialSpec:         newFakeUnstructured("batch/v1", "Job"),
			wantCollectorKind: common.StdOutCollector,
		},
		{
			name:              "TrainJob trial defaults to Push collector",
			trialSpec:         newFakeUnstructured("trainer.kubeflow.org/v1alpha1", "TrainJob"),
			wantCollectorKind: common.PushCollector,
		},
		{
			name:      "explicit collector is preserved for TrainJob",
			trialSpec: newFakeUnstructured("trainer.kubeflow.org/v1alpha1", "TrainJob"),
			explicitCollector: &common.CollectorSpec{
				Kind: common.StdOutCollector,
			},
			wantCollectorKind: common.StdOutCollector,
		},
		{
			name:              "nil TrialSpec defaults to StdOut collector",
			trialSpec:         nil,
			wantCollectorKind: common.StdOutCollector,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			e := &Experiment{
				Spec: ExperimentSpec{},
			}
			if tc.trialSpec != nil {
				e.Spec.TrialTemplate = &TrialTemplate{
					TrialSource: TrialSource{
						TrialSpec: tc.trialSpec,
					},
				}
			}
			if tc.explicitCollector != nil {
				e.Spec.MetricsCollectorSpec = &common.MetricsCollectorSpec{
					Collector: tc.explicitCollector,
				}
			}

			e.setDefaultMetricsCollector()

			gotKind := e.Spec.MetricsCollectorSpec.Collector.Kind
			if gotKind != tc.wantCollectorKind {
				t.Errorf("got collector kind %q, want %q", gotKind, tc.wantCollectorKind)
			}
		})
	}
}

func newFakeUnstructured(apiVersion, kind string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.SetAPIVersion(apiVersion)
	u.SetKind(kind)
	return u
}
