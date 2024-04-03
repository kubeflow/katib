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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"k8s.io/klog"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
)

var (
	errFileFormat = fmt.Errorf("format must be set %v or %v", commonv1beta1.TextFormat, commonv1beta1.JsonFormat)
	errOpenFile   = errors.New("failed to open the file")
	errReadFile   = errors.New("failed to read the file")
	errParseJson  = errors.New("failed to parse the json object")
)

func CollectObservationLog(fileName string, metrics []string, filters []string, fileFormat commonv1beta1.FileFormat) (*v1beta1.ObservationLog, error) {
	// we should check fileFormat first in case of opening an invalid file
	if fileFormat != commonv1beta1.JsonFormat && fileFormat != commonv1beta1.TextFormat {
		return nil, errFileFormat
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errOpenFile, err.Error())
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errReadFile, err.Error())
	}
	logs := string(content)

	switch fileFormat {
	case commonv1beta1.TextFormat:
		return parseLogsInTextFormat(strings.Split(logs, "\n"), metrics, filters)
	case commonv1beta1.JsonFormat:
		return parseLogsInJsonFormat(strings.Split(logs, "\n"), metrics)
	default:
		return nil, nil
	}
}

func parseLogsInTextFormat(logs []string, metrics []string, filters []string) (*v1beta1.ObservationLog, error) {
	metricRegList := GetFilterRegexpList(filters)
	mlogs := make([]*v1beta1.MetricLog, 0, len(logs))

	for _, logline := range logs {
		// skip line which doesn't contain any metrics keywords, avoiding unnecessary pattern match
		isMetricLine := false
		for _, m := range metrics {
			if strings.Contains(logline, m) {
				isMetricLine = true
				break
			}
		}
		if !isMetricLine {
			continue
		}

		timestamp := time.Time{}.UTC().Format(time.RFC3339)
		ls := strings.SplitN(logline, " ", 2)
		if len(ls) != 2 {
			klog.Warningf("Metrics will not have timestamp since %s doesn't begin with timestamp string", logline)
		} else {
			if _, err := time.Parse(time.RFC3339Nano, ls[0]); err != nil {
				klog.Warningf("Metrics will not have timestamp since error parsing time %s: %v", ls[0], err)
			} else {
				timestamp = ls[0]
			}
		}

		for _, metricReg := range metricRegList {
			matchStrs := metricReg.FindAllStringSubmatch(logline, -1)
			for _, kevList := range matchStrs {
				if len(kevList) < 3 {
					continue
				}
				name := strings.TrimSpace(kevList[1])
				value := strings.TrimSpace(kevList[2])
				for _, m := range metrics {
					if name != m {
						continue
					}
					mlogs = append(mlogs, &v1beta1.MetricLog{
						TimeStamp: timestamp,
						Metric: &v1beta1.Metric{
							Name:  name,
							Value: value,
						},
					})
					break
				}
			}
		}
	}
	return newObservationLog(mlogs, metrics), nil
}

func parseLogsInJsonFormat(logs []string, metrics []string) (*v1beta1.ObservationLog, error) {
	mlogs := make([]*v1beta1.MetricLog, 0, len(logs))

	for _, logline := range logs {
		if len(logline) == 0 {
			continue
		}
		var jsonObj map[string]interface{}
		if err := json.Unmarshal([]byte(logline), &jsonObj); err != nil {
			return nil, fmt.Errorf("%w: %s", errParseJson, err.Error())
		}

		timestamp := time.Time{}.UTC().Format(time.RFC3339)
		timestampJsonValue, exist := jsonObj[common.TimeStampJsonKey]
		if !exist {
			klog.Warningf("Metrics will not have timestamp since %s doesn't have the key timestamp", logline)
		} else {
			if parsedTimestamp := parseTimestamp(timestampJsonValue); parsedTimestamp == "" {
				klog.Warningf("Metrics will not have timestamp since error parsing time %v", timestampJsonValue)
			} else {
				timestamp = parsedTimestamp
			}
		}

		for _, m := range metrics {
			value, exist := jsonObj[m].(string)
			if !exist {
				continue
			}
			mlogs = append(mlogs, &v1beta1.MetricLog{
				TimeStamp: timestamp,
				Metric: &v1beta1.Metric{
					Name:  m,
					Value: value,
				},
			})
		}
	}
	return newObservationLog(mlogs, metrics), nil
}

func newObservationLog(mlogs []*v1beta1.MetricLog, metrics []string) *v1beta1.ObservationLog {
	// Metrics logs must contain at least one objective metric value
	// Objective metric is located at first index
	isObjectiveMetricReported := false
	for _, mLog := range mlogs {
		if mLog.Metric.Name == metrics[0] {
			isObjectiveMetricReported = true
			break
		}
	}
	// If objective metrics were not reported, insert unavailable value in the DB
	if !isObjectiveMetricReported {
		klog.Infof("Objective metric %v is not found in training logs, %v value is reported", metrics[0], consts.UnavailableMetricValue)
		return &v1beta1.ObservationLog{
			MetricLogs: []*v1beta1.MetricLog{
				{
					TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
					Metric: &v1beta1.Metric{
						Name:  metrics[0],
						Value: consts.UnavailableMetricValue,
					},
				},
			},
		}
	}
	return &v1beta1.ObservationLog{
		MetricLogs: mlogs,
	}
}

func parseTimestamp(timestamp interface{}) string {
	if stringTimestamp, ok := timestamp.(string); ok {

		if stringTimestamp == "" {
			klog.Warningln("Timestamp is empty")
			return ""
		} else if _, err := time.Parse(time.RFC3339Nano, stringTimestamp); err != nil {
			klog.Warningf("Failed to parse timestamp since %s is not RFC3339Nano format", stringTimestamp)
			return ""
		}
		return stringTimestamp

	} else {

		floatTimestamp, ok := timestamp.(float64)
		if !ok {
			klog.Warningf("Failed to parse timestamp since the type of %v is neither string nor float64", timestamp)
			return ""
		}

		stringTimestamp = strconv.FormatFloat(floatTimestamp, 'f', -1, 64)
		t := strings.Split(stringTimestamp, ".")

		sec, err := strconv.ParseInt(t[0], 10, 64)
		if err != nil {
			klog.Warningf("Failed to parse timestamp; %v", err)
			return ""
		}

		var nanoSec int64 = 0
		if len(t) == 2 {
			nanoSec, err = strconv.ParseInt(t[1], 10, 64)
			if err != nil {
				klog.Warningf("Failed to parse timestamp; %v", err)
				return ""
			}
		}

		return time.Unix(sec, nanoSec).UTC().Format(time.RFC3339Nano)
	}
}

// GetFilterRegexpList returns Regexp array from filters string array
func GetFilterRegexpList(filters []string) []*regexp.Regexp {
	regexpList := make([]*regexp.Regexp, 0, len(filters))
	if len(filters) == 0 {
		filters = append(filters, common.DefaultFilter)
	}
	for _, filter := range filters {
		reg, _ := regexp.Compile(filter)
		regexpList = append(regexpList, reg)
	}
	return regexpList
}
