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

	"github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

type TrialsCollector struct {
	store                   cache.Cache
	trialDeleteCount        *prometheus.CounterVec
	trialCreateCount        *prometheus.CounterVec
	trialSucceedCount       *prometheus.CounterVec
	trialFailCount          *prometheus.CounterVec
	trialMetricsUnavailable *prometheus.CounterVec
	trialCurrent            *prometheus.GaugeVec
}

func NewTrialsCollector(store cache.Cache, registerer prometheus.Registerer) *TrialsCollector {
	c := &TrialsCollector{
		store: store,
		trialDeleteCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "katib_trial_deleted_total",
			Help: "The total number of deleted trials",
		}, []string{"namespace"}),

		trialCreateCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "katib_trial_created_total",
			Help: "The total number of created trials",
		}, []string{"namespace"}),

		trialSucceedCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "katib_trial_succeeded_total",
			Help: "The total number of succeeded trials",
		}, []string{"namespace"}),

		trialFailCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "katib_trial_failed_total",
			Help: "The total number of failed trials",
		}, []string{"namespace"}),

		trialMetricsUnavailable: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "katib_trial_metrics_unavailable_total",
			Help: "The total number of metrics unavailable trials",
		}, []string{"namespace"}),

		trialCurrent: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "katib_trials_current",
			Help: "The number of current katib trials in the cluster",
		}, []string{"namespace", "status"}),
	}
	registerer.MustRegister(c)
	return c
}

// Describe implements the prometheus.Collector interface.
func (m *TrialsCollector) Describe(ch chan<- *prometheus.Desc) {
	m.trialDeleteCount.Describe(ch)
	m.trialSucceedCount.Describe(ch)
	m.trialFailCount.Describe(ch)
	m.trialCreateCount.Describe(ch)
	m.trialMetricsUnavailable.Describe(ch)
	m.trialCurrent.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (m *TrialsCollector) Collect(ch chan<- prometheus.Metric) {
	m.collect()
	m.trialDeleteCount.Collect(ch)
	m.trialSucceedCount.Collect(ch)
	m.trialFailCount.Collect(ch)
	m.trialCreateCount.Collect(ch)
	m.trialMetricsUnavailable.Collect(ch)
	m.trialCurrent.Collect(ch)
}

func (c *TrialsCollector) IncreaseTrialsDeletedCount(ns string) {
	c.trialDeleteCount.WithLabelValues(ns).Inc()
}

func (c *TrialsCollector) IncreaseTrialsCreatedCount(ns string) {
	c.trialCreateCount.WithLabelValues(ns).Inc()
}

func (c *TrialsCollector) IncreaseTrialsSucceededCount(ns string) {
	c.trialSucceedCount.WithLabelValues(ns).Inc()
}

func (c *TrialsCollector) IncreaseTrialsFailedCount(ns string) {
	c.trialFailCount.WithLabelValues(ns).Inc()
}

func (c *TrialsCollector) IncreaseTrialsMetricsUnavailableCount(ns string) {
	c.trialMetricsUnavailable.WithLabelValues(ns).Inc()
}

// collect gets the current experiments from cache.
func (c *TrialsCollector) collect() {
	var (
		conditionType v1beta1.TrialConditionType
		status        string
		err           error
	)
	trialLists := &v1beta1.TrialList{}
	if err = c.store.List(context.TODO(), trialLists); err != nil {
		return
	}

	trialCache := map[string]map[string]int{}
	for _, trial := range trialLists.Items {
		conditionType, err = trial.GetLastConditionType()
		status = string(conditionType)
		// If trial doesn't have any condition, use unknown.
		if err != nil {
			status = "Unknown"
		}

		if _, ok := trialCache[trial.Namespace]; !ok {
			trialCache[trial.Namespace] = make(map[string]int)
		}
		trialCache[trial.Namespace][status] += 1
	}

	c.trialCurrent.Reset()
	for ns, v := range trialCache {
		for status, count := range v {
			c.trialCurrent.WithLabelValues(ns, status).Set(float64(count))
		}
	}
}
