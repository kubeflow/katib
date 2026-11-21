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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
	authzv1 "k8s.io/api/authorization/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	katibclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/util/katibclient"
)

func TestFetchNASJobInfoAuthorization(t *testing.T) {
	origDisableAuth := DISABLE_AUTH
	DISABLE_AUTH = "false"
	defer func() { DISABLE_AUTH = origDisableAuth }()

	cases := map[string]struct {
		userHeader     string
		allowed        bool
		wantStatusCode int
		wantTrialList  bool
	}{
		"missing user header": {
			wantStatusCode: http.StatusUnauthorized,
		},
		"user not authorized in namespace": {
			userHeader:     "alice@example.com",
			wantStatusCode: http.StatusForbidden,
		},
		"user authorized in namespace": {
			userHeader:     "alice@example.com",
			allowed:        true,
			wantStatusCode: http.StatusOK,
			wantTrialList:  true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			sarClient := fake.NewClientBuilder().
				WithScheme(clientgoscheme.Scheme).
				WithInterceptorFuncs(interceptor.Funcs{
					Create: func(_ context.Context, _ client.WithWatch, obj client.Object, _ ...client.CreateOption) error {
						if sar, ok := obj.(*authzv1.SubjectAccessReview); ok {
							sar.Status.Allowed = tc.allowed
						}
						return nil
					},
				}).
				Build()

			katibClient := katibclientmock.NewMockClient(mockCtrl)
			katibClient.EXPECT().GetClient().Return(sarClient).AnyTimes()
			if tc.wantTrialList {
				katibClient.EXPECT().GetTrialList("nas-exp", "other-ns").Return(&trialsv1beta1.TrialList{}, nil)
			}

			k := &KatibUIHandler{
				katibClient:   katibClient,
				dbManagerAddr: "localhost:6789",
			}

			req := httptest.NewRequest(http.MethodGet, "/katib/fetch_nas_job_info/?experimentName=nas-exp&namespace=other-ns", nil)
			if tc.userHeader != "" {
				req.Header.Set(USER_HEADER, tc.userHeader)
			}
			rec := httptest.NewRecorder()

			k.FetchNASJobInfo(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("FetchNASJobInfo() status = %d, want %d, body: %s", rec.Code, tc.wantStatusCode, rec.Body.String())
			}
		})
	}
}
