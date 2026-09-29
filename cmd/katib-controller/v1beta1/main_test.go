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

package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	health_pb "github.com/kubeflow/katib/pkg/apis/manager/health"
)

// fakeServer lets us control what StartedChecker() reports without
// standing up a real webhook HTTPS listener.
type fakeServer struct {
	webhook.Server
	startedErr error
}

func (f *fakeServer) StartedChecker() healthz.Checker {
	return func(req *http.Request) error {
		return f.startedErr
	}
}

func TestReadyzCheck(t *testing.T) {
	errNotStarted := errors.New("webhook server has not been started yet")

	cases := map[string]struct {
		startedErr  error
		cacheSynced bool
		wantErr     bool
	}{
		"webhook server not started": {
			startedErr:  errNotStarted,
			cacheSynced: true,
			wantErr:     true,
		},
		"webhook server started but caches not synced": {
			startedErr:  nil,
			cacheSynced: false,
			wantErr:     true,
		},
		"webhook server started and caches synced": {
			startedErr:  nil,
			cacheSynced: true,
			wantErr:     false,
		},
		"webhook server not started and caches not synced": {
			startedErr:  errNotStarted,
			cacheSynced: false,
			wantErr:     true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			hookServer := &fakeServer{startedErr: tc.startedErr}
			cacheSynced := func(ctx context.Context) bool {
				return tc.cacheSynced
			}

			check := readyzCheck(hookServer, cacheSynced)
			req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
			err := check(req)

			if tc.wantErr && err == nil {
				t.Errorf("readyzCheck() = nil, want an error")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("readyzCheck() = %v, want nil", err)
			}
		})
	}
}

// TestReadyzCheck_BoundsCacheSyncContext guards against regressing to an
// unbounded wait: cacheSynced must be given a context with its own
// deadline, not the bare (deadline-less) request context, so a single
// /readyz poll can't block forever while the caches never sync.
func TestReadyzCheck_BoundsCacheSyncContext(t *testing.T) {
	hookServer := &fakeServer{startedErr: nil}
	var gotDeadline bool
	cacheSynced := func(ctx context.Context) bool {
		_, gotDeadline = ctx.Deadline()
		return true
	}

	check := readyzCheck(hookServer, cacheSynced)
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	if err := check(req); err != nil {
		t.Fatalf("readyzCheck() = %v, want nil", err)
	}

	if !gotDeadline {
		t.Errorf("cacheSynced was called with a context that has no deadline; want a bounded context")
	}
}

// fakeAPIReader lets us control what a Get through the manager's uncached
// API reader reports, without a real apiserver.
type fakeAPIReader struct {
	client.Reader
	getErr      error
	gotDeadline bool
}

func (f *fakeAPIReader) Get(ctx context.Context, _ client.ObjectKey, _ client.Object, _ ...client.GetOption) error {
	_, f.gotDeadline = ctx.Deadline()
	return f.getErr
}

func TestApiserverCheck(t *testing.T) {
	errUnreachable := errors.New("dial tcp: connection refused")

	cases := map[string]struct {
		getErr  error
		wantErr bool
	}{
		"apiserver reachable": {
			getErr:  nil,
			wantErr: false,
		},
		"apiserver unreachable": {
			getErr:  errUnreachable,
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			reader := &fakeAPIReader{getErr: tc.getErr}
			check := apiserverCheck(reader, "kubeflow")
			req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
			err := check(req)

			if tc.wantErr && err == nil {
				t.Errorf("apiserverCheck() = nil, want an error")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("apiserverCheck() = %v, want nil", err)
			}
			if !reader.gotDeadline {
				t.Errorf("Get was called with a context that has no deadline; want a bounded context")
			}
		})
	}
}

// fakeHealthClient lets us control what db-manager's Check RPC reports
// without a real katib-db-manager connection.
type fakeHealthClient struct {
	resp        *health_pb.HealthCheckResponse
	err         error
	gotDeadline bool
}

func (f *fakeHealthClient) Check(ctx context.Context, _ *health_pb.HealthCheckRequest, _ ...grpc.CallOption) (*health_pb.HealthCheckResponse, error) {
	_, f.gotDeadline = ctx.Deadline()
	return f.resp, f.err
}

// TestCheckDBManager covers checkDBManager directly rather than through
// runDBManagerDiagnostics's ticker loop: the interesting logic (how a Check
// RPC result maps to healthy/unhealthy) lives entirely in checkDBManager,
// and it is deliberately no longer wired into AddReadyzCheck (see
// runDBManagerDiagnostics's comment for why), so there is no
// healthz.Checker/*http.Request shape to exercise here anymore.
func TestCheckDBManager(t *testing.T) {
	errUnreachable := errors.New("dial tcp: connection refused")

	cases := map[string]struct {
		resp    *health_pb.HealthCheckResponse
		err     error
		wantErr bool
	}{
		"db-manager serving": {
			resp:    &health_pb.HealthCheckResponse{Status: health_pb.HealthCheckResponse_SERVING},
			wantErr: false,
		},
		"db-manager not serving": {
			resp:    &health_pb.HealthCheckResponse{Status: health_pb.HealthCheckResponse_NOT_SERVING},
			wantErr: true,
		},
		"db-manager unreachable": {
			err:     errUnreachable,
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			healthClient := &fakeHealthClient{resp: tc.resp, err: tc.err}
			err := checkDBManager(context.Background(), healthClient)

			if tc.wantErr && err == nil {
				t.Errorf("checkDBManager() = nil, want an error")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("checkDBManager() = %v, want nil", err)
			}
		})
	}
}
