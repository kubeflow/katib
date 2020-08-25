/*
Copyright 2018 The Kubeflow Authors

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

/*
MetricsCollector is a default metricscollector for worker.
It will collect metrics from pod log.
You should print metrics in {{MetricsName}}={{MetricsValue}} format.
For example, the objective value name is F1 and the metrics are loss, your training code should print like below.
     ---
     epoch 1:
     batch1 loss=0.8
     batch2 loss=0.6

     F1=0.4

     epoch 2:
     batch1 loss=0.4
     batch2 loss=0.2

     F1=0.7
     ---
The metrics collector will collect all logs of metrics.
*/

package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/hpcloud/tail"
	"google.golang.org/grpc"
	"k8s.io/klog"

	api "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
	filemc "github.com/kubeflow/katib/pkg/metricscollector/v1beta1/file-metricscollector"
)

var (
	managerServiceAddr = flag.String("s", "", "Katib Manager service")
	trialName          = flag.String("t", "", "Trial Name")
	metricsFilePath    = flag.String("path", "", "Metrics File Path")
	metricNames        = flag.String("m", "", "Metric names")
	metricFilters      = flag.String("f", "", "Metric filters")
	pollInterval       = flag.Duration("p", common.DefaultPollInterval, "Poll interval between running processes check")
	timeout            = flag.Duration("timeout", common.DefaultTimeout, "Timeout before invoke error during running processes check")
	waitAll            = flag.Bool("w", common.DefaultWaitAll, "Whether wait for all other main process of container exiting")
)

func printMetricsFile(mFile string) {
	for {
		_, err := os.Stat(mFile)
		if err == nil {
			break
		} else if os.IsNotExist(err) {
			continue
		} else {
			klog.Fatalf("could not watch metrics file: %v", err)
		}
	}

	t, _ := tail.TailFile(mFile, tail.Config{Follow: true})
	for line := range t.Lines {
		klog.Info(line.Text)
	}
}

func main() {
	flag.Parse()
	klog.Infof("Trial Name: %s", *trialName)

	go printMetricsFile(*metricsFilePath)
	wopts := common.WaitPidsOpts{
		PollInterval:           *pollInterval,
		Timeout:                *timeout,
		WaitAll:                *waitAll,
		CompletedMarkedDirPath: filepath.Dir(*metricsFilePath),
	}
	if err := common.WaitMainProcesses(wopts); err != nil {
		klog.Fatalf("Failed to wait for worker container: %v", err)
	}

	conn, err := grpc.Dial(*managerServiceAddr, grpc.WithInsecure())
	if err != nil {
		klog.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewDBManagerClient(conn)
	ctx := context.Background()
	var metricList []string
	if len(*metricNames) != 0 {
		metricList = strings.Split(*metricNames, ";")
	}
	var filterList []string
	if len(*metricFilters) != 0 {
		filterList = strings.Split(*metricFilters, ";")
	}
	olog, err := filemc.CollectObservationLog(*metricsFilePath, metricList, filterList)
	if err != nil {
		klog.Fatalf("Failed to collect logs: %v", err)
	}
	reportreq := &api.ReportObservationLogRequest{
		TrialName:      *trialName,
		ObservationLog: olog,
	}
	_, err = c.ReportObservationLog(ctx, reportreq)
	if err != nil {
		klog.Fatalf("Failed to Report logs: %v", err)
	}
	klog.Infof("Metrics reported. :\n%v", olog)
}
