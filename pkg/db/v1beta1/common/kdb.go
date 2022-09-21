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

package common

import (
	v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

type KatibDBInterface interface {
	DBInit()
	SelectOne() error

	RegisterObservationLog(trialName string, observationLog *v1beta1.ObservationLog) error
	GetObservationLog(trialName string, metricName string, startTime string, endTime string) (*v1beta1.ObservationLog, error)
	DeleteObservationLog(trialName string) error
}
