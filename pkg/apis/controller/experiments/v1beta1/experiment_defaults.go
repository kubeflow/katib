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

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

func (e *Experiment) SetDefault() {
	e.setDefaultParallelTrialCount()
	e.setDefaultResumePolicy()
	e.setDefaultObjective()
	e.setDefaultTrialTemplate()
	e.setDefaultMetricsCollector()
}

func (e *Experiment) setDefaultParallelTrialCount() {
	if e.Spec.ParallelTrialCount == nil {
		e.Spec.ParallelTrialCount = new(int32)
		*e.Spec.ParallelTrialCount = DefaultTrialParallelCount
	}
}

func (e *Experiment) setDefaultResumePolicy() {
	if e.Spec.ResumePolicy == "" {
		e.Spec.ResumePolicy = DefaultResumePolicy
	}
}

func (e *Experiment) setDefaultObjective() {
	obj := e.Spec.Objective
	if obj != nil {
		if obj.MetricStrategies == nil {
			obj.MetricStrategies = make([]common.MetricStrategy, 0)
		}
		objectiveHasDefault := false
		metricsWithDefault := make(map[string]int)
		for _, strategy := range obj.MetricStrategies {
			if strategy.Name == obj.ObjectiveMetricName {
				objectiveHasDefault = true
				continue
			}
			metricsWithDefault[strategy.Name] = 1
		}

		// Set default strategy of objective according to ObjectiveType.
		if !objectiveHasDefault {
			var strategy common.MetricStrategy
			switch e.Spec.Objective.Type {
			case common.ObjectiveTypeMinimize:
				strategy = common.MetricStrategy{Name: obj.ObjectiveMetricName, Value: common.ExtractByMin}
			case common.ObjectiveTypeMaximize:
				strategy = common.MetricStrategy{Name: obj.ObjectiveMetricName, Value: common.ExtractByMax}
			default:
				strategy = common.MetricStrategy{Name: obj.ObjectiveMetricName, Value: common.ExtractByLatest}
			}
			obj.MetricStrategies = append(obj.MetricStrategies, strategy)
		}

		// Set default strategy of additional metrics to ExtractByLatest.
		for _, metricName := range obj.AdditionalMetricNames {
			if _, ok := metricsWithDefault[metricName]; !ok {
				var strategy common.MetricStrategy
				switch e.Spec.Objective.Type {
				case common.ObjectiveTypeMinimize:
					strategy = common.MetricStrategy{Name: metricName, Value: common.ExtractByMin}
				case common.ObjectiveTypeMaximize:
					strategy = common.MetricStrategy{Name: metricName, Value: common.ExtractByMax}
				default:
					strategy = common.MetricStrategy{Name: metricName, Value: common.ExtractByLatest}
				}
				obj.MetricStrategies = append(obj.MetricStrategies, strategy)
			}
		}
	}
}

func (e *Experiment) setDefaultTrialTemplate() {
	t := e.Spec.TrialTemplate

	// Set default values for Job and Kubeflow Training Job if TrialSpec is not nil
	if t != nil && t.TrialSource.TrialSpec != nil {
		jobKind := t.TrialSource.TrialSpec.GetKind()
		if jobKind == consts.JobKindJob {
			if t.SuccessCondition == "" {
				t.SuccessCondition = DefaultJobSuccessCondition
			}
			if t.FailureCondition == "" {
				t.FailureCondition = DefaultJobFailureCondition
			}
		} else if _, ok := KubeflowJobKinds[jobKind]; ok {
			if t.SuccessCondition == "" {
				t.SuccessCondition = DefaultKubeflowJobSuccessCondition
			}
			if t.FailureCondition == "" {
				t.FailureCondition = DefaultKubeflowJobFailureCondition
			}
			// For Kubeflow Job also set default PrimaryPodLabels
			if len(t.PrimaryPodLabels) == 0 {
				t.PrimaryPodLabels = DefaultKubeflowJobPrimaryPodLabels
			}
		}
	}
	e.Spec.TrialTemplate = t
}

func (e *Experiment) setDefaultMetricsCollector() {
	if e.Spec.MetricsCollectorSpec == nil {
		e.Spec.MetricsCollectorSpec = &common.MetricsCollectorSpec{}
	}
	if e.Spec.MetricsCollectorSpec.Collector == nil {
		e.Spec.MetricsCollectorSpec.Collector = &common.CollectorSpec{
			Kind: common.StdOutCollector,
		}
	}
	switch e.Spec.MetricsCollectorSpec.Collector.Kind {
	case common.PrometheusMetricCollector:
		if e.Spec.MetricsCollectorSpec.Source == nil {
			e.Spec.MetricsCollectorSpec.Source = &common.SourceSpec{}
		}
		if e.Spec.MetricsCollectorSpec.Source.HttpGet == nil {
			e.Spec.MetricsCollectorSpec.Source.HttpGet = &v1.HTTPGetAction{}
		}
		if e.Spec.MetricsCollectorSpec.Source.HttpGet.Path == "" {
			e.Spec.MetricsCollectorSpec.Source.HttpGet.Path = common.DefaultPrometheusPath
		}
		if e.Spec.MetricsCollectorSpec.Source.HttpGet.Port.String() == "0" {
			e.Spec.MetricsCollectorSpec.Source.HttpGet.Port = intstr.FromInt(common.DefaultPrometheusPort)
		}
	case common.FileCollector:
		if e.Spec.MetricsCollectorSpec.Source == nil {
			e.Spec.MetricsCollectorSpec.Source = &common.SourceSpec{}
		}
		if e.Spec.MetricsCollectorSpec.Source.FileSystemPath == nil {
			e.Spec.MetricsCollectorSpec.Source.FileSystemPath = &common.FileSystemPath{}
		}
		if e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Kind == "" {
			e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Kind = common.FileKind
		}
		if e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Path == "" {
			e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Path = common.DefaultFilePath
		}
		if e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Format == "" {
			e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Format = common.TextFormat
		}
	case common.TfEventCollector:
		if e.Spec.MetricsCollectorSpec.Source == nil {
			e.Spec.MetricsCollectorSpec.Source = &common.SourceSpec{}
		}
		if e.Spec.MetricsCollectorSpec.Source.FileSystemPath == nil {
			e.Spec.MetricsCollectorSpec.Source.FileSystemPath = &common.FileSystemPath{}
		}
		if e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Kind == "" {
			e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Kind = common.DirectoryKind
		}
		if e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Path == "" {
			e.Spec.MetricsCollectorSpec.Source.FileSystemPath.Path = common.DefaultTensorflowEventDirPath
		}
	}
}
