package sidecarmetricscollector

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	v1alpha3 "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"github.com/kubeflow/katib/pkg/metricscollector/v1alpha3/common"
	"k8s.io/klog"
)

type FileMetricsCollector struct {
}

func NewFileMetricsCollector() (*FileMetricsCollector, error) {
	return &FileMetricsCollector{}, nil
}

func (d *FileMetricsCollector) CollectObservationLog(fileName string, metrics []string, filters []string) (*v1alpha3.ObservationLog, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	logs := string(content)
	olog, err := d.parseLogs(strings.Split(logs, "\n"), metrics, filters)
	return olog, err
}

func (d *FileMetricsCollector) parseLogs(logs []string, metrics []string, filters []string) (*v1alpha3.ObservationLog, error) {
	var lasterr error
	olog := &v1alpha3.ObservationLog{}
	mlogs := []*v1alpha3.MetricLog{}
	metricRegList := getFilterRegexpList(filters)

	for _, logline := range logs {
		if logline == "" {
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
					mlogs = append(mlogs, &v1alpha3.MetricLog{
						TimeStamp: timestamp,
						Metric: &v1alpha3.Metric{
							Name:  name,
							Value: value,
						},
					})
					break
				}
			}
		}
	}
	olog.MetricLogs = mlogs
	if lasterr != nil {
		return olog, lasterr
	}
	return olog, nil
}

func getFilterRegexpList(filters []string) []*regexp.Regexp {
	regexpList := []*regexp.Regexp{}
	if len(filters) == 0 {
		filters = append(filters, common.DefaultFilter)
	}
	for _, filter := range filters {
		reg, _ := regexp.Compile(filter)
		regexpList = append(regexpList, reg)
	}
	return regexpList
}
