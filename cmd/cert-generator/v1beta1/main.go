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
	"github.com/kubeflow/katib/pkg/cert-generator/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	kubeClient, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme})
	if err != nil {
		klog.Fatalf("Failed to create kube client.")
	}

	cmd, err := v1beta1.NewKatibCertGeneratorCmd(kubeClient)
	if err != nil {
		klog.Fatalf("Failed to generate cert: %v", err)
	}

	if err = cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
