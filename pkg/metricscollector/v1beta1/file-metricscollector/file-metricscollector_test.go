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
	"reflect"
	"testing"
	"time"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

var testJsonDataPath = filepath.Join("testdata", "JSON")

func TestCollectObservationLog(t *testing.T) {

	if err := generateTestFiles(); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(filepath.Dir(testJsonDataPath))

	// TODO (tenzen-y): We should add tests for logs in TEXT format.
	// Ref: https://github.com/kubeflow/katib/issues/1756
	testCases := []struct {
		description string
		filePath    string
		metrics     []string
		filters     []string
		fileFormat  commonv1beta1.FileFormat
		err         bool
		expected    *v1beta1.ObservationLog
	}{
		{
			description: "Positive case for logs in JSON format",
			filePath:    filepath.Join(testJsonDataPath, "good.json"),
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
			filePath:    "invalid",
			err:         true,
		},
		{
			description: "Invalid file format",
			filePath:    filepath.Join(testJsonDataPath, "good.json"),
			fileFormat:  "invalid",
			err:         true,
		},
		{
			description: "Invalid formatted file for logs in JSON format",
			filePath:    filepath.Join(testJsonDataPath, "invalid-format.json"),
			fileFormat:  commonv1beta1.JsonFormat,
			err:         true,
		},
		{
			description: "Invalid timestamp for logs in JSON format",
			filePath:    filepath.Join(testJsonDataPath, "invalid-timestamp.json"),
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
			filePath:    filepath.Join(testJsonDataPath, "missing-objective-metric.json"),
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
			actual, err := CollectObservationLog(test.filePath, test.metrics, test.filters, test.fileFormat)
			if (err != nil) != test.err {
				t.Errorf("\nGOT: \n%v\nWANT: %v\n", err, test.err)
			} else {
				if !reflect.DeepEqual(actual, test.expected) {
					t.Errorf("Expected %v\n got %v", test.expected, actual)
				}
			}
		})
	}
}

func generateTestFiles() error {
	if _, err := os.Stat(testJsonDataPath); err != nil {
		if err = os.MkdirAll(testJsonDataPath, 0700); err != nil {
			return err
		}
	}

	testData := []struct {
		fileName string
		data     string
	}{
		{
			fileName: "good.json",
			data: `{"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": 1638422847.28721, "trial": "0"}
{"acc": "0.9349666833877563", "checkpoint_path": "", "global_step": "0", "timestamp": 1638422847.287801, "trial": "0"}
{"checkpoint_path": "", "global_step": "1", "loss": "0.1414974331855774", "timestamp": "2021-12-02T14:27:50.000035161Z", "trial": "0"}
{"acc": "0.9586416482925415", "checkpoint_path": "", "global_step": "1", "timestamp": "2021-12-02T14:27:50.000037459Z", "trial": "0"}
{"checkpoint_path": "", "global_step": "2", "loss": "0.10683439671993256", "trial": "0"}
`,
		},
		{
			fileName: "invalid-format.json",
			data: `"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": 1638422847.28721, "trial": "0"
{"acc": "0.9349666833877563", "checkpoint_path": "", "global_step": "0", "timestamp": 1638422847.287801, "trial": "0
`,
		},
		{
			fileName: "invalid-timestamp.json",
			data: `{"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": "invalid", "trial": "0"}
{"acc": "0.9349666833877563", "checkpoint_path": "", "global_step": "0", "timestamp": 1638422847, "trial": "0"}
`,
		}, {
			fileName: "missing-objective-metric.json",
			data: `{"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": 1638422847.28721, "trial": "0"}
{"checkpoint_path": "", "global_step": "1", "loss": "0.1414974331855774", "timestamp": "2021-12-02T14:27:50.000035161+09:00", "trial": "0"}
{"checkpoint_path": "", "global_step": "2", "loss": "0.10683439671993256", "trial": "0"}`,
		},
	}

	for _, td := range testData {
		filePath := filepath.Join(testJsonDataPath, td.fileName)
		if err := os.WriteFile(filePath, []byte(td.data), 0600); err != nil {
			return err
		}
	}

	return nil
}
