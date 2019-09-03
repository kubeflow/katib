/*
Copyright 2019 The Kubernetes Authors.

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
	stdlog "log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	apis "github.com/kubeflow/katib/pkg/apis/controller"
)

var (
	cfg                      *rest.Config
	controlPlaneStartTimeout = 60 * time.Second
	controlPlaneStopTimeout  = 60 * time.Second
)

func TestMain(m *testing.M) {
	t := &envtest.Environment{
		ControlPlaneStartTimeout: controlPlaneStartTimeout,
		ControlPlaneStopTimeout:  controlPlaneStopTimeout,
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "manifests", "v1alpha2", "katib-controller"),
			filepath.Join("..", "..", "..", "..", "test", "unit", "v1alpha2", "crds"),
		},
	}
	apis.AddToScheme(scheme.Scheme)

	var err error
	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}

	code := m.Run()
	t.Stop()
	os.Exit(code)
}

// SetupTestReconcile returns a reconcile.Reconcile implementation that delegates to inner.
func SetupTestReconcile(inner reconcile.Reconciler) reconcile.Reconciler {
	fn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(req)
		return result, err
	})
	return fn
}

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager, g *gomega.GomegaWithT) (chan struct{}, *sync.WaitGroup) {
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(stop)).NotTo(gomega.HaveOccurred())
	}()
	return stop, wg
}

// addForTestPurpose adds a new Controller to mgr with r as the reconcile.Reconciler.
func addForTestPurpose(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("test-experiment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		log.Error(err, "Failed to create experiment controller for test purpose.")
		return err
	}

	if err = addWatch(mgr, c); err != nil {
		log.Error(err, "Trial watch failed")
		return err
	}

	log.Info("Experiment controller created")
	return nil
}
