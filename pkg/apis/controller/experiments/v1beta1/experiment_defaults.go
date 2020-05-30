// +build !ignore_autogenerated

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
// Code generated by main. DO NOT EDIT.

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
	if obj.MetricStrategies == nil {
		obj.MetricStrategies = make(map[string]common.MetricStrategy)
	}
	// set default strategy of objective according to ObjectiveType
	if _, ok := obj.MetricStrategies[obj.ObjectiveMetricName]; !ok {
		switch e.Spec.Objective.Type {
		case common.ObjectiveTypeMinimize:
			obj.MetricStrategies[obj.ObjectiveMetricName] = common.ExtractByMin
		case common.ObjectiveTypeMaximize:
			obj.MetricStrategies[obj.ObjectiveMetricName] = common.ExtractByMax
		case common.ObjectiveTypeUnknown:
			obj.MetricStrategies[obj.ObjectiveMetricName] = common.ExtractByLatest
		default:
			obj.MetricStrategies[obj.ObjectiveMetricName] = common.ExtractByLatest
		}
	}
	// set default strategy of additional metrics to ExtractByLatest
	for _, name := range obj.AdditionalMetricNames {
		if _, ok := obj.MetricStrategies[name]; !ok {
			obj.MetricStrategies[name] = common.ExtractByLatest
		}
	}
}

func (e *Experiment) setDefaultTrialTemplate() {
	t := e.Spec.TrialTemplate
	if t == nil {
		t = &TrialTemplate{
			Retain: true,
		}
	}
	if t.TrialSource.TrialSpec == nil && t.TrialSource.ConfigMap == nil && t.TrialParameters == nil {
		t.TrialSource = TrialSource{
			ConfigMap: &ConfigMapSource{
				ConfigMapNamespace: consts.DefaultKatibNamespace,
				ConfigMapName:      DefaultTrialConfigMapName,
				TemplatePath:       DefaultTrialTemplatePath,
			},
		}
		t.TrialParameters = []TrialParameterSpec{
			{
				Name:        "learningRate",
				Description: "Learning rate for the training model",
				Reference:   "lr",
			},
			{
				Name:        "numberLayers",
				Description: "Number of training model layers",
				Reference:   "num-layers",
			},
			{
				Name:        "optimizer",
				Description: "Training model optimizer (sdg, adam or ftrl)",
				Reference:   "optimizer",
			},
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
