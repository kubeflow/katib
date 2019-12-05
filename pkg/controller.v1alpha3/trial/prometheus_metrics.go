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

package trial

import (
	"github.com/prometheus/client_golang/prometheus"

	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	trialsDeletedCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "katib_trial_deleted_total",
		Help: "Counts number of Trial deleted",
	})
	trialsCreatedCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "katib_trial_created_total",
		Help: "Counts number of Trial created",
	})
	trialsSucceededCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "katib_trial_succeeded_total",
		Help: "Counts number of Trial succeeded",
	})
	trialsFailedCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "katib_trial_failed_total",
		Help: "Counts number of Trial failed",
	})
)

func init() {
	metrics.Registry.MustRegister(
		trialsDeletedCount,
		trialsCreatedCount,
		trialsSucceededCount,
		trialsFailedCount)
}

func IncreaseTrialsDeletedCount() {
	trialsDeletedCount.Inc()
}

func IncreaseTrialsCreatedCount() {
	trialsCreatedCount.Inc()
}

func IncreaseTrialsSucceededCount() {
	trialsSucceededCount.Inc()
}

func IncreaseTrialsFailedCount() {
	trialsFailedCount.Inc()
}
