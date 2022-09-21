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

package util

import (
	"context"

	"github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

type ExperimentsCollector struct {
	store           cache.Cache
	expDeleteCount  *prometheus.CounterVec
	expCreateCount  *prometheus.CounterVec
	expSucceedCount *prometheus.CounterVec
	expFailCount    *prometheus.CounterVec
	expCurrent      *prometheus.GaugeVec
}

func NewExpsCollector(store cache.Cache, registerer prometheus.Registerer) *ExperimentsCollector {
	c := &ExperimentsCollector{
		store: store,
		expDeleteCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "katib_experiment_deleted_total",
			Help: "The total number of deleted experiments",
		}, []string{"namespace"}),

		expCreateCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "katib_experiment_created_total",
			Help: "The total number of created experiments",
		}, []string{"namespace"}),

		expSucceedCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "katib_experiment_succeeded_total",
			Help: "The total number of succeeded experiments",
		}, []string{"namespace"}),

		expFailCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "katib_experiment_failed_total",
			Help: "The total number of failed experiments",
		}, []string{"namespace"}),

		expCurrent: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "katib_experiments_current",
			Help: "The number of current katib experiments in the cluster",
		}, []string{"namespace", "status"}),
	}
	registerer.MustRegister(c)
	return c
}

// Describe implements the prometheus.Collector interface.
func (m *ExperimentsCollector) Describe(ch chan<- *prometheus.Desc) {
	m.expDeleteCount.Describe(ch)
	m.expSucceedCount.Describe(ch)
	m.expFailCount.Describe(ch)
	m.expCreateCount.Describe(ch)
	m.expCurrent.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (m *ExperimentsCollector) Collect(ch chan<- prometheus.Metric) {
	m.collect()
	m.expDeleteCount.Collect(ch)
	m.expSucceedCount.Collect(ch)
	m.expFailCount.Collect(ch)
	m.expCreateCount.Collect(ch)
	m.expCurrent.Collect(ch)
}

func (c *ExperimentsCollector) IncreaseExperimentsDeletedCount(ns string) {
	c.expDeleteCount.WithLabelValues(ns).Inc()
}

func (c *ExperimentsCollector) IncreaseExperimentsCreatedCount(ns string) {
	c.expCreateCount.WithLabelValues(ns).Inc()
}

func (c *ExperimentsCollector) IncreaseExperimentsSucceededCount(ns string) {
	c.expSucceedCount.WithLabelValues(ns).Inc()
}

func (c *ExperimentsCollector) IncreaseExperimentsFailedCount(ns string) {
	c.expFailCount.WithLabelValues(ns).Inc()
}

// collect gets the current experiments from cache.
func (c *ExperimentsCollector) collect() {
	var (
		conditionType v1beta1.ExperimentConditionType
		status        string
		err           error
	)
	expLists := &v1beta1.ExperimentList{}
	if err = c.store.List(context.TODO(), expLists); err != nil {
		return
	}

	expCache := map[string]map[string]int{}
	for _, exp := range expLists.Items {
		conditionType, err = exp.GetLastConditionType()
		status = string(conditionType)
		// If experiment doesn't have any condition, use unknown.
		if err != nil {
			status = "Unknown"
		}

		if _, ok := expCache[exp.Namespace]; !ok {
			expCache[exp.Namespace] = make(map[string]int)
		}
		expCache[exp.Namespace][status] += 1
	}

	c.expCurrent.Reset()
	for ns, v := range expCache {
		for status, count := range v {
			c.expCurrent.WithLabelValues(ns, status).Set(float64(count))
		}
	}
}
