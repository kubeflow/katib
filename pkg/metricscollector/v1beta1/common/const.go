/*

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
	"time"

	v1beta1common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
)

const (
	// DefaultPollInterval is the default value for interval between running processes check
	DefaultPollInterval = time.Second
	// DefaultTimeout is the default value for timeout before invoke error during running processes check
	// To run without timeout set value to 0
	DefaultTimeout = 0
	// DefaultWaitAll is the default value whether wait for all other main process of container exiting
	DefaultWaitAll = true
	// TrainingCompleted is the job finished marker in $$$$.pid file when main process is completed
	TrainingCompleted = "completed"

	// DefaultFilter is the default metrics collector filter to parse the metrics.
	// Metrics must be printed this way
	// loss=0.3
	// accuracy=0.98
	DefaultFilter = `([\w|-]+)\s*=\s*((-?\d+)(\.\d+)?)`

	// TODO (andreyvelich): Do we need to maintain 2 names? Should we leave only 1?
	MetricCollectorContainerName       = "metrics-collector"
	MetricLoggerCollectorContainerName = "metrics-logger-and-collector"
)

var (
	AutoInjectMetricsCollecterList = [...]v1beta1common.CollectorKind{
		v1beta1common.StdOutCollector,
		v1beta1common.TfEventCollector,
		v1beta1common.FileCollector,
		v1beta1common.PrometheusMetricCollector,
	}
)
