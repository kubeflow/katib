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

package sidecarmetricscollector

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCollectObservationLog(t *testing.T) {
	testCases := map[string]struct {
		fileName   string
		testData   string
		metrics    []string
		filters    []string
		fileFormat commonv1beta1.FileFormat
		wantError  error
		expected   *v1beta1.ObservationLog
	}{
		"Positive case for logs in JSON format": {
			fileName: "good.json",
			testData: `{"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": 1638422847.28721, "trial": "0"}
{"acc": "0.9349666833877563", "checkpoint_path": "", "global_step": "0", "timestamp": 1638422847.287801, "trial": "0"}
{"checkpoint_path": "", "global_step": "1", "loss": "0.1414974331855774", "timestamp": "2021-12-02T14:27:50.000035161Z", "trial": "0"}
{"acc": "0.9586416482925415", "checkpoint_path": "", "global_step": "1", "timestamp": "2021-12-02T14:27:50.000037459Z", "trial": "0"}
{"checkpoint_path": "", "global_step": "2", "loss": "0.10683439671993256", "trial": "0"}`,
			metrics:    []string{"acc", "loss"},
			fileFormat: commonv1beta1.JsonFormat,
			expected: &v1beta1.ObservationLog{
				MetricLogs: []*v1beta1.MetricLog{
					{
						TimeStamp: "2021-12-02T05:27:27.000028721Z",
						Metric: &v1beta1.Metric{
							Name:  "loss",
							Value: "0.22082142531871796",
						},
					},
					{
						TimeStamp: "2021-12-02T05:27:27.000287801Z",
						Metric: &v1beta1.Metric{
							Name:  "acc",
							Value: "0.9349666833877563",
						},
					},
					{
						TimeStamp: "2021-12-02T14:27:50.000035161Z",
						Metric: &v1beta1.Metric{
							Name:  "loss",
							Value: "0.1414974331855774",
						},
					},
					{
						TimeStamp: "2021-12-02T14:27:50.000037459Z",
						Metric: &v1beta1.Metric{
							Name:  "acc",
							Value: "0.9586416482925415",
						},
					},
					{
						TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
						Metric: &v1beta1.Metric{
							Name:  "loss",
							Value: "0.10683439671993256",
						},
					},
				},
			},
		},
		"Positive case for logs in TEXT format": {
			fileName: "good.log",
			testData: `2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 0.8078};{metricName: loss, metricValue: 0.5183}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 0.6752}
2024-03-04T17:55:08Z INFO     {metricName: loss, metricValue: 0.3634}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 100}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 888.333}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: -0.4759}
{metricName: loss, metricValue: 0.8671}`,
			metrics:    []string{"accuracy", "loss"},
			filters:    []string{"{metricName: ([\\w|-]+), metricValue: ((-?\\d+)(\\.\\d+)?)}"},
			fileFormat: commonv1beta1.TextFormat,
			expected: &v1beta1.ObservationLog{
				MetricLogs: []*v1beta1.MetricLog{
					{
						TimeStamp: "2024-03-04T17:55:08Z",
						Metric: &v1beta1.Metric{
							Name:  "accuracy",
							Value: "0.8078",
						},
					},
					{
						TimeStamp: "2024-03-04T17:55:08Z",
						Metric: &v1beta1.Metric{
							Name:  "loss",
							Value: "0.5183",
						},
					},
					{
						TimeStamp: "2024-03-04T17:55:08Z",
						Metric: &v1beta1.Metric{
							Name:  "accuracy",
							Value: "0.6752",
						},
					},
					{
						TimeStamp: "2024-03-04T17:55:08Z",
						Metric: &v1beta1.Metric{
							Name:  "loss",
							Value: "0.3634",
						},
					},
					{
						TimeStamp: "2024-03-04T17:55:08Z",
						Metric: &v1beta1.Metric{
							Name:  "accuracy",
							Value: "100",
						},
					},
					{
						TimeStamp: "2024-03-04T17:55:08Z",
						Metric: &v1beta1.Metric{
							Name:  "accuracy",
							Value: "888.333",
						},
					},
					{
						TimeStamp: "2024-03-04T17:55:08Z",
						Metric: &v1beta1.Metric{
							Name:  "accuracy",
							Value: "-0.4759",
						},
					},
					{
						TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
						Metric: &v1beta1.Metric{
							Name:  "loss",
							Value: "0.8671",
						},
					},
				},
			},
		},
		"Invalid case for logs in TEXT format": {
			fileName: "invalid-value.log",
			testData: `2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: .333}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: -.333}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: - 345.333}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 888.}`,
			filters:    []string{"{metricName: ([\\w|-]+), metricValue: ((-?\\d+)(\\.\\d+)?)}"},
			metrics:    []string{"accuracy", "loss"},
			fileFormat: commonv1beta1.TextFormat,
			expected: &v1beta1.ObservationLog{
				MetricLogs: []*v1beta1.MetricLog{
					{
						TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
						Metric: &v1beta1.Metric{
							Name:  "accuracy",
							Value: consts.UnavailableMetricValue,
						},
					},
				},
			},
		},
		"Invalid file name": {
			fileName:   "invalid",
			fileFormat: commonv1beta1.JsonFormat,
			wantError:  errOpenFile,
		},
		"Invalid file format": {
			fileName:   "good.log",
			fileFormat: "invalid",
			wantError:  errFileFormat,
		},
		"Invalid formatted file for logs in JSON format": {
			fileName: "invalid-format.json",
			testData: `"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": 1638422847.28721, "trial": "0"
{"acc": "0.9349666833877563", "checkpoint_path": "", "global_step": "0", "timestamp": 1638422847.287801, "trial": "0`,
			fileFormat: commonv1beta1.JsonFormat,
			wantError:  errParseJson,
		},
		"Invalid formatted file for logs in TEXT format": {
			fileName: "invalid-format.log",
			testData: `2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 0.6752
2024-03-04T17:55:08Z INFO     {metricName: loss, metricValue: 0.3634}`,
			filters:    []string{"{metricName: ([\\w|-]+), metricValue: ((-?\\d+)(\\.\\d+)?)}"},
			metrics:    []string{"accuracy", "loss"},
			fileFormat: commonv1beta1.TextFormat,
			expected: &v1beta1.ObservationLog{
				MetricLogs: []*v1beta1.MetricLog{
					{
						TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
						Metric: &v1beta1.Metric{
							Name:  "accuracy",
							Value: consts.UnavailableMetricValue,
						},
					},
				},
			},
		},
		"Invalid timestamp for logs in JSON format": {
			fileName: "invalid-timestamp.json",
			testData: `{"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": "invalid", "trial": "0"}
{"acc": "0.9349666833877563", "checkpoint_path": "", "global_step": "0", "timestamp": 1638422847, "trial": "0"}`,
			fileFormat: commonv1beta1.JsonFormat,
			metrics:    []string{"acc", "loss"},
			expected: &v1beta1.ObservationLog{
				MetricLogs: []*v1beta1.MetricLog{
					{
						TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
						Metric: &v1beta1.Metric{
							Name:  "loss",
							Value: "0.22082142531871796",
						},
					},
					{
						TimeStamp: "2021-12-02T05:27:27Z",
						Metric: &v1beta1.Metric{
							Name:  "acc",
							Value: "0.9349666833877563",
						},
					},
				},
			},
		},
		"Invalid timestamp for logs in TEXT format": {
			fileName: "invalid-timestamp.log",
			testData: `2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 0.6752}
invalid INFO     {metricName: loss, metricValue: 0.3634}`,
			metrics:    []string{"accuracy", "loss"},
			filters:    []string{"{metricName: ([\\w|-]+), metricValue: ((-?\\d+)(\\.\\d+)?)}"},
			fileFormat: commonv1beta1.TextFormat,
			expected: &v1beta1.ObservationLog{
				MetricLogs: []*v1beta1.MetricLog{
					{
						TimeStamp: "2024-03-04T17:55:08Z",
						Metric: &v1beta1.Metric{
							Name:  "accuracy",
							Value: "0.6752",
						},
					},
					{
						TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
						Metric: &v1beta1.Metric{
							Name:  "loss",
							Value: "0.3634",
						},
					},
				},
			},
		},
		"Missing objective metric in JSON training logs": {
			fileName: "missing-objective-metric.json",
			testData: `{"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": 1638422847.28721, "trial": "0"}
{"checkpoint_path": "", "global_step": "1", "loss": "0.1414974331855774", "timestamp": "2021-12-02T14:27:50.000035161+09:00", "trial": "0"}
{"checkpoint_path": "", "global_step": "2", "loss": "0.10683439671993256", "trial": "0"}`,
			fileFormat: commonv1beta1.JsonFormat,
			metrics:    []string{"acc", "loss"},
			expected: &v1beta1.ObservationLog{
				MetricLogs: []*v1beta1.MetricLog{
					{
						TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
						Metric: &v1beta1.Metric{
							Name:  "acc",
							Value: consts.UnavailableMetricValue,
						},
					},
				},
			},
		},
		"Missing objective metric in TEXT training logs": {
			fileName: "missing-objective-metric.log",
			testData: `2024-03-04T17:55:08Z INFO     {metricName: loss, metricValue: 0.3634}
2024-03-04T17:55:08Z INFO     {metricName: loss, metricValue: 0.8671}`,
			fileFormat: commonv1beta1.TextFormat,
			metrics:    []string{"accuracy", "loss"},
			expected: &v1beta1.ObservationLog{
				MetricLogs: []*v1beta1.MetricLog{
					{
						TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
						Metric: &v1beta1.Metric{
							Name:  "accuracy",
							Value: consts.UnavailableMetricValue,
						},
					},
				},
			},
		},
	}

	tmpDir := t.TempDir()
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			if test.testData != "" {
				if err := os.WriteFile(filepath.Join(tmpDir, test.fileName), []byte(test.testData), 0600); err != nil {
					t.Fatalf("failed to write test data: %v", err)
				}
			}
			actual, err := CollectObservationLog(filepath.Join(tmpDir, test.fileName), test.metrics, test.filters, test.fileFormat)
			if diff := cmp.Diff(test.wantError, err, cmpopts.EquateErrors()); len(diff) != 0 {
				t.Errorf("Unexpected error (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(test.expected, actual, protocmp.Transform()); len(diff) != 0 {
				t.Errorf("Unexpected parsed result (-want,+got):\n%s", diff)
			}
		})
	}
}
