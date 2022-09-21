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

package experiment

import (
	"context"
	stdlog "log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	apis "github.com/kubeflow/katib/pkg/apis/controller"
)

var (
	cfg                      *rest.Config
	ctx                      context.Context
	cancel                   context.CancelFunc
	controlPlaneStartTimeout = 60 * time.Second
	controlPlaneStopTimeout  = 60 * time.Second
)

func TestMain(m *testing.M) {
	// To avoid the `timeout waiting for process kube-apiserver to stop` error,
	// we must use the `context.WithCancel`.
	// Ref: https://github.com/kubernetes-sigs/controller-runtime/issues/1571#issuecomment-945535598
	ctx, cancel = context.WithCancel(context.TODO())

	t := &envtest.Environment{
		ControlPlaneStartTimeout: controlPlaneStartTimeout,
		ControlPlaneStopTimeout:  controlPlaneStopTimeout,
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "manifests", "v1beta1", "components", "crd"),
		},
	}
	var err error
	if err = apis.AddToScheme(scheme.Scheme); err != nil {
		stdlog.Fatal(err)
	}

	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}

	code := m.Run()
	cancel()
	if err = t.Stop(); err != nil {
		stdlog.Fatal(err)
	}
	os.Exit(code)
}

// SetupTestReconcile returns a reconcile.Reconcile implementation that delegates to inner.
func SetupTestReconcile(inner reconcile.Reconciler) reconcile.Reconciler {
	fn := reconcile.Func(func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(ctx, req)
		return result, err
	})
	return fn
}
