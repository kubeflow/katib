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

package pod

import (
	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
)

const (
	MasterRole = "master"
	BatchJob   = "Job"
	// TrialKind is the name of Trial kind
	TrialKind = "Trial"
	// TrialAPIVersion is the name of Trial API Version
	TrialAPIVersion = "kubeflow.org/v1beta1"
)

var (
	NeedWrapWorkerMetricsCollectorList = [...]common.CollectorKind{
		common.StdOutCollector,
		common.TfEventCollector,
		common.FileCollector,
	}
)
