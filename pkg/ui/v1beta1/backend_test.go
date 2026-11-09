/*
Copyright 2026 The Kubeflow Authors.

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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	experimentv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	katibclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/util/katibclient"
)

func newExperiment(name, namespace string) experimentv1beta1.Experiment {
	return experimentv1beta1.Experiment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Status: experimentv1beta1.ExperimentStatus{
			Conditions: []experimentv1beta1.ExperimentCondition{
				{Type: experimentv1beta1.ExperimentRunning},
			},
		},
	}
}

func TestDeleteExperimentListsOnlyRequestNamespace(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	katibClient := katibclientmock.NewMockClient(mockCtrl)
	katibClient.EXPECT().GetClient().Return(nil).AnyTimes()
	k := &KatibUIHandler{katibClient: katibClient}

	deleted := newExperiment("exp", "team-a")
	katibClient.EXPECT().GetExperiment("exp", "team-a").Return(&deleted, nil)
	katibClient.EXPECT().DeleteRuntimeObject(&deleted).Return(nil)
	katibClient.EXPECT().GetExperimentList(gomock.Any()).DoAndReturn(
		func(namespace ...string) (*experimentv1beta1.ExperimentList, error) {
			list := &experimentv1beta1.ExperimentList{}
			if len(namespace) == 1 && namespace[0] == "team-a" {
				list.Items = append(list.Items, newExperiment("other", "team-a"))
			} else {
				// Cluster-wide listing exposes Experiments from other namespaces.
				list.Items = append(list.Items, newExperiment("other", "team-a"), newExperiment("secret", "team-b"))
			}
			return list, nil
		}).AnyTimes()

	req := httptest.NewRequest(http.MethodDelete, "/katib/delete_experiment/?experimentName=exp&namespace=team-a", nil)
	rec := httptest.NewRecorder()
	k.DeleteExperiment(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Unexpected status code: %v, body: %s", rec.Code, rec.Body.String())
	}
	var experiments []ExperimentView
	if err := json.Unmarshal(rec.Body.Bytes(), &experiments); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	for _, e := range experiments {
		if e.Namespace != "team-a" {
			t.Errorf("DeleteExperiment response leaked Experiment %s/%s", e.Namespace, e.Name)
		}
	}
}
