package sidecarmetricscollector

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	v1alpha3 "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"k8s.io/klog"
)

type FileMetricsCollector struct {
}

func NewFileMetricsCollector() (*FileMetricsCollector, error) {
	return &FileMetricsCollector{}, nil
}

func (d *FileMetricsCollector) CollectObservationLog(fileName string, metrics []string) (*v1alpha3.ObservationLog, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	logs := string(content)
	olog, err := d.parseLogs(strings.Split(logs, "\n"), metrics)
	return olog, err
}

func (d *FileMetricsCollector) parseLogs(logs []string, metrics []string) (*v1alpha3.ObservationLog, error) {
	// TODO(hougangliu): handle custom filter string
	var lasterr error
	olog := &v1alpha3.ObservationLog{}
	mlogs := []*v1alpha3.MetricLog{}
	for _, logline := range logs {
		if logline == "" {
			continue
		}
		timestamp := time.Time{}.UTC().Format(time.RFC3339)
		parseStr := logline
		ls := strings.SplitN(logline, " ", 2)
		if len(ls) != 2 {
			klog.Warningf("Metrics will not have timestamp since %s doesn't begin with timestamp string", logline)
		} else {
			_, err := time.Parse(time.RFC3339Nano, ls[0])
			if err != nil {
				klog.Warningf("Metrics will not have timestamp since error parsing time %s: %v", ls[0], err)
			} else {
				parseStr = ls[1]
				timestamp = ls[0]
			}
		}

		kvpairs := strings.Fields(parseStr)
		for _, kv := range kvpairs {
			v := strings.Split(kv, "=")
			if len(v) > 2 {
				klog.Infof("Ignoring trailing garbage: %s", kv)
			}
			if len(v) == 1 {
				continue
			}
			metricName := ""
			for _, m := range metrics {
				if v[0] == m {
					metricName = v[0]
				}
			}
			if metricName == "" {
				continue
			}
			mlogs = append(mlogs, &v1alpha3.MetricLog{
				TimeStamp: timestamp,
				Metric: &v1alpha3.Metric{
					Name:  metricName,
					Value: v[1],
				},
			})
		}
	}
	olog.MetricLogs = mlogs
	if lasterr != nil {
		return olog, lasterr
	}
	return olog, nil
}
