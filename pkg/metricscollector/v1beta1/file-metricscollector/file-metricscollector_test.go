/*
Copyright 2021 The Kubeflow Authors.

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
	"path"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

func TestCollectObservationLog(t *testing.T) {
	const testJsonDataPath = "testdata/JSON"

	// TODO (tenzen-y): We should add tests for logs in TEXT format.
	// Ref: https://github.com/kubeflow/katib/issues/1756
	testCases := []struct {
		description string
		fileName    string
		metrics     []string
		filters     []string
		fileFormat  commonv1beta1.FileSystemFileFormat
		err         bool
		expected    *v1beta1.ObservationLog
	}{
		{
			description: "Positive case for logs in JSON format",
			fileName:    path.Join(testJsonDataPath, "good.json"),
			metrics:     []string{"acc", "loss"},
			fileFormat:  commonv1beta1.JsonFormat,
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
		{
			description: "Invalid file name",
			fileName:    "invalid",
			err:         true,
		},
		{
			description: "Invalid file format",
			fileName:    path.Join(testJsonDataPath, "good.json"),
			fileFormat:  "invalid",
			err:         true,
		},
		{
			description: "Invalid formatted file for logs in JSON format",
			fileName:    path.Join(testJsonDataPath, "invalid-format.json"),
			fileFormat:  commonv1beta1.JsonFormat,
			err:         true,
		},
		{
			description: "Invalid timestamp for logs in JSON format",
			fileName:    path.Join(testJsonDataPath, "invalid-timestamp.json"),
			fileFormat:  commonv1beta1.JsonFormat,
			metrics:     []string{"acc", "loss"},
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
		{
			description: "Missing objective metric in training logs",
			fileName:    path.Join(testJsonDataPath, "missing-objective-metric.json"),
			fileFormat:  commonv1beta1.JsonFormat,
			metrics:     []string{"acc", "loss"},
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
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			actual, err := CollectObservationLog(test.fileName, test.metrics, test.filters, test.fileFormat)
			if (err != nil) != test.err {
				t.Errorf("\nGOT: \n%v\nWANT: %v\n", err, test.err)
			} else {
				if diff := cmp.Diff(actual, test.expected); diff != "" {
					t.Errorf("\nDIFF: \n%v\n", diff)
				}
			}
		})
	}
}
