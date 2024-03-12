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
	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

const (
	validJSONTestFile            = "good.json"
	invalidFormatJSONTestFile    = "invalid-format.json"
	invalidTimestampJSONTestFile = "invalid-timestamp.json"
	missingMetricJSONTestFile    = "missing-objective-metric.json"

	validTEXTTestFile            = "good.log"
	invalidValueTEXTTestFile     = "invalid-value.log"
	invalidFormatTEXTTestFile    = "invalid-format.log"
	invalidTimestampTEXTTestFile = "invalid-timestamp.log"
	missingMetricTEXTTestFile    = "missing-objective-metric.log"
)

var (
	testJsonDataPath = filepath.Join("testdata", "JSON")
	testTextDataPath = filepath.Join("testdata", "TEXT")
)

func TestCollectObservationLog(t *testing.T) {
	if err := generateTestDirs(); err != nil {
		t.Fatal(err)
	}
	if err := generateJSONTestFiles(); err != nil {
		t.Fatal(err)
	}
	if err := generateTEXTTestFiles(); err != nil {
		t.Fatal(err)
	}
	defer deleteTestDirs()

	testCases := map[string]struct {
		filePath   string
		metrics    []string
		filters    []string
		fileFormat commonv1beta1.FileFormat
		err        bool
		expected   *v1beta1.ObservationLog
	}{
		"Positive case for logs in JSON format": {
			filePath:   filepath.Join(testJsonDataPath, validJSONTestFile),
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
			filePath:   filepath.Join(testTextDataPath, validTEXTTestFile),
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
			filePath:   filepath.Join(testTextDataPath, invalidValueTEXTTestFile),
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
			filePath: "invalid",
			err:      true,
		},
		"Invalid file format": {
			filePath:   filepath.Join(testTextDataPath, validTEXTTestFile),
			fileFormat: "invalid",
			err:        true,
		},
		"Invalid formatted file for logs in JSON format": {
			filePath:   filepath.Join(testJsonDataPath, invalidFormatJSONTestFile),
			fileFormat: commonv1beta1.JsonFormat,
			err:        true,
		},
		"Invalid formatted file for logs in TEXT format": {
			filePath:   filepath.Join(testTextDataPath, invalidFormatTEXTTestFile),
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
			filePath:   filepath.Join(testJsonDataPath, invalidTimestampJSONTestFile),
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
			filePath:   filepath.Join(testTextDataPath, invalidTimestampTEXTTestFile),
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
			filePath:   filepath.Join(testJsonDataPath, missingMetricJSONTestFile),
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
			filePath:   filepath.Join(testTextDataPath, missingMetricTEXTTestFile),
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

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			actual, err := CollectObservationLog(test.filePath, test.metrics, test.filters, test.fileFormat)
			if (err != nil) != test.err {
				t.Errorf("\nGOT: \n%v\nWANT: %v\n", err, test.err)
			} else {
				if diff := cmp.Diff(test.expected, actual); diff != "" {
					t.Errorf("Unexpected parsed result (-want,+got):\n%s", diff)
				}
			}
		})
	}
}

func generateTestDirs() error {
	// Generate JSON files' dir
	if _, err := os.Stat(testJsonDataPath); err != nil {
		if err = os.MkdirAll(testJsonDataPath, 0700); err != nil {
			return err
		}
	}

	// Generate TEXT files' dir
	if _, err := os.Stat(testTextDataPath); err != nil {
		if err = os.MkdirAll(testTextDataPath, 0700); err != nil {
			return err
		}
	}

	return nil
}

func deleteTestDirs() error {
	if err := os.RemoveAll(filepath.Dir(testJsonDataPath)); err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Dir(testTextDataPath)); err != nil {
		return err
	}

	return nil
}

func generateJSONTestFiles() error {
	testData := []struct {
		fileName string
		data     string
	}{
		{
			fileName: validJSONTestFile,
			data: `{"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": 1638422847.28721, "trial": "0"}
{"acc": "0.9349666833877563", "checkpoint_path": "", "global_step": "0", "timestamp": 1638422847.287801, "trial": "0"}
{"checkpoint_path": "", "global_step": "1", "loss": "0.1414974331855774", "timestamp": "2021-12-02T14:27:50.000035161Z", "trial": "0"}
{"acc": "0.9586416482925415", "checkpoint_path": "", "global_step": "1", "timestamp": "2021-12-02T14:27:50.000037459Z", "trial": "0"}
{"checkpoint_path": "", "global_step": "2", "loss": "0.10683439671993256", "trial": "0"}
`,
		},
		{
			fileName: invalidFormatJSONTestFile,
			data: `"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": 1638422847.28721, "trial": "0"
{"acc": "0.9349666833877563", "checkpoint_path": "", "global_step": "0", "timestamp": 1638422847.287801, "trial": "0
`,
		},
		{
			fileName: invalidTimestampJSONTestFile,
			data: `{"checkpoint_path": "", "global_step": "0", "loss": "0.22082142531871796", "timestamp": "invalid", "trial": "0"}
{"acc": "0.9349666833877563", "checkpoint_path": "", "global_step": "0", "timestamp": 1638422847, "trial": "0"}
`,
		}, {
			fileName: missingMetricJSONTestFile,
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

func generateTEXTTestFiles() error {
	testData := []struct {
		fileName string
		data     string
	}{
		{
			fileName: validTEXTTestFile,
			data: `2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 0.8078};{metricName: loss, metricValue: 0.5183}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 0.6752}
2024-03-04T17:55:08Z INFO     {metricName: loss, metricValue: 0.3634}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 100}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 888.333}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: -0.4759}
{metricName: loss, metricValue: 0.8671}`,
		},
		{
			fileName: invalidValueTEXTTestFile,
			data: `2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: .333}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: -.333}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: - 345.333}
2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 888.}`,
		},
		{
			fileName: invalidFormatTEXTTestFile,
			data: `2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 0.6752
2024-03-04T17:55:08Z INFO     {metricName: loss, metricValue: 0.3634}`,
		},
		{
			fileName: invalidTimestampTEXTTestFile,
			data: `2024-03-04T17:55:08Z INFO     {metricName: accuracy, metricValue: 0.6752}
invalid INFO     {metricName: loss, metricValue: 0.3634}`,
		},
		{
			fileName: missingMetricTEXTTestFile,
			data: `2024-03-04T17:55:08Z INFO     {metricName: loss, metricValue: 0.3634}
2024-03-04T17:55:08Z INFO     {metricName: loss, metricValue: 0.8671}`,
		},
	}

	for _, td := range testData {
		filePath := filepath.Join(testTextDataPath, td.fileName)
		if err := os.WriteFile(filePath, []byte(td.data), 0600); err != nil {
			return err
		}
	}

	return nil
}
