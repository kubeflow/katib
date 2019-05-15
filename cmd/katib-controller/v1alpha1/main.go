/*
Copyright 2018 The Kubeflow Authors

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

/*
 StudyJobController is a controller (operator) for StudyJob
 StudyJobController create and watch workers and metricscollectors.
 The workers and metricscollectors are generated from template defined ConfigMap.
 The workers and metricscollectors are kubernetes object. The default object is a Job and CronJob.
*/
package main

import (
	"k8s.io/klog"

	"github.com/kubeflow/katib/pkg/api/operators/apis"
	controller "github.com/kubeflow/katib/pkg/controller/v1alpha1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func main() {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		klog.Info("config.GetConfig()")
		klog.Fatal(err)
	}

	// Create a new StudyJobController to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		klog.Info("manager.New")
		klog.Fatal(err)
	}

	klog.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Info("apis.AddToScheme")
		klog.Fatal(err)
	}

	// Setup StudyJobController
	if err := controller.AddToManager(mgr); err != nil {
		klog.Info("controller.AddToManager(mgr)")
		klog.Fatal(err)
	}

	klog.Info("Starting the Cmd.")

	// Starting the StudyJobController
	klog.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
