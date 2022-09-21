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

package suggestion

import (
	"context"

	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	"k8s.io/apimachinery/pkg/api/equality"
)

const (
	SuggestionCreatedReason      = "SuggestionCreated"
	SuggestionDeploymentReady    = "DeploymentReady"
	SuggestionDeploymentNotReady = "DeploymentNotReady"
	SuggestionRunningReason      = "SuggestionRunning"
	SuggestionFailedReason       = "SuggestionFailed"
)

func (r *ReconcileSuggestion) updateStatus(s *suggestionsv1beta1.Suggestion, oldS *suggestionsv1beta1.Suggestion) error {
	if !equality.Semantic.DeepEqual(s.Status, oldS.Status) {
		if err := r.Status().Update(context.TODO(), s); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReconcileSuggestion) updateStatusCondition(s *suggestionsv1beta1.Suggestion, oldS *suggestionsv1beta1.Suggestion) error {
	if !equality.Semantic.DeepEqual(s.Status.Conditions, oldS.Status.Conditions) {
		newConditions := s.Status.Conditions
		s.Status = oldS.Status
		s.Status.Conditions = newConditions
		if err := r.Status().Update(context.TODO(), s); err != nil {
			return err
		}
	}
	return nil
}
