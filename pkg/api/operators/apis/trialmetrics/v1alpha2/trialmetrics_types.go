/*

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

package v1alpha2

import (
	common "github.com/kubeflow/katib/pkg/api/operators/apis/common/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TrialMetrics struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TrialMetricsSpec   `json:"spec,omitempty"`
	Status            TrialMetricsStatus `json:"status,omitempty"`
}

type TrialMetricsSpec struct {
	MetricNames []string `json:"metricNames,omitempty"`
}

type TrialMetricsStatus struct {
	StartTime         *metav1.Time        `json:"startTime,omitempty"`
	CompletionTime    *metav1.Time        `json:"completionTime,omitempty"`
	LastReconcileTime *metav1.Time        `json:"lastReconcileTime,omitempty"`
	Observation       *common.Observation `json:"observation,omitempty"`
	Epoch             *int32              `json:"epoch,omitempty"`
	Step              *int32              `json:"step,omitempty"`
}

// TrialMetricsList contains a list of Trials
type TrialMetricsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TrialMetrics `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrialMetrics{}, &TrialMetricsList{})
}
