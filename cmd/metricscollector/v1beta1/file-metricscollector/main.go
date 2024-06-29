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
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	psutil "github.com/shirou/gopsutil/v3/process"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	api "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
	filemc "github.com/kubeflow/katib/pkg/metricscollector/v1beta1/file-metricscollector"
)

type stopRulesFlag []commonv1beta1.EarlyStoppingRule

func (flag *stopRulesFlag) String() string {
	stopRuleStrings := []string{}
	for _, r := range *flag {
		stopRuleStrings = append(stopRuleStrings, r.Name)
		stopRuleStrings = append(stopRuleStrings, r.Value)
		stopRuleStrings = append(stopRuleStrings, string(r.Comparison))
		stopRuleStrings = append(stopRuleStrings, strconv.Itoa(r.StartStep))
	}
	return strings.Join(stopRuleStrings, ";")
}

func (flag *stopRulesFlag) Set(value string) error {
	stopRuleParsed := strings.Split(value, ";")
	if len(stopRuleParsed) != 4 {
		return fmt.Errorf("Invalid Early Stopping rule: %v", value)
	}

	// Get int start step.
	startStep, err := strconv.Atoi(stopRuleParsed[3])
	if err != nil {
		klog.Fatalf("Parse start step: %v to int error: %v", stopRuleParsed[3], err)
	}

	// For each stop rule this order: 1 - metric name, 2 - metric value, 3 - comparison type, 4 - start step.
	// Start step is equal to 0, if it's not defined.
	stopRule := commonv1beta1.EarlyStoppingRule{
		Name:       stopRuleParsed[0],
		Value:      stopRuleParsed[1],
		Comparison: commonv1beta1.ComparisonType(stopRuleParsed[2]),
		StartStep:  startStep,
	}

	*flag = append(*flag, stopRule)
	return nil
}

var (
	dbManagerServiceAddr = flag.String("s-db", "", "Katib DB Manager service endpoint")
	earlyStopServiceAddr = flag.String("s-earlystop", "", "Katib Early Stopping service endpoint")
	trialName            = flag.String("t", "", "Trial Name")
	metricsFilePath      = flag.String("path", "", "Metrics File Path")
	metricsFileFormat    = flag.String("format", "", "Metrics File Format")
	metricNames          = flag.String("m", "", "Metric names")
	objectiveType        = flag.String("o-type", "", "Objective type")
	metricFilters        = flag.String("f", "", "Metric filters")
	pollInterval         = flag.Duration("p", common.DefaultPollInterval, "Poll interval between running processes check")
	timeout              = flag.Duration("timeout", common.DefaultTimeout, "Timeout before invoke error during running processes check")
	waitAllProcesses     = flag.String("w", common.DefaultWaitAllProcesses, "Whether wait for all other main process of container exiting")
	stopRules            stopRulesFlag
	isEarlyStopped       = false
)

func checkMetricFile(mFile string) {
	for {
		_, err := os.Stat(mFile)
		if err == nil {
			break
		} else if os.IsNotExist(err) {
			continue
		} else {
			klog.Fatalf("Could not watch metrics file: %v", err)
		}
	}
}

func printMetricsFile(mFile string) {

	// Check that metric file exists.
	checkMetricFile(mFile)

	// Print lines from metrics file.
	t, _ := tail.TailFile(mFile, tail.Config{Follow: true})
	for line := range t.Lines {
		klog.Info(line.Text)
	}
}

func watchMetricsFile(mFile string, stopRules stopRulesFlag, filters []string, fileFormat commonv1beta1.FileFormat) {
	// First metric is objective in metricNames array.
	objMetric := strings.Split(*metricNames, ";")[0]
	objType := commonv1beta1.ObjectiveType(*objectiveType)

	rules, err := filemc.NewRuleSet(objMetric, objType, stopRules)
	if err != nil {
		klog.Fatalf("NewRuleSet failed: %v", err)
	}

	// Check that metric file exists.
	checkMetricFile(mFile)

	// Get Main process.
	// Extract the metric file dir path based on the file name.
	mDirPath, _ := filepath.Split(mFile)
	_, mainProcPid, err := common.GetMainProcesses(mDirPath)
	if err != nil {
		klog.Fatalf("GetMainProcesses failed: %v", err)
	}
	mainProc, err := psutil.NewProcess(int32(mainProcPid))
	if err != nil {
		klog.Fatalf("Failed to create new Process from pid %v, error: %v", mainProcPid, err)
	}

	// Get list of regural expressions from filters.
	metricRegList := filemc.GetFilterRegexpList(filters)

	liveRuleMetrics := rules.LiveMetrics()

	// Start watch log lines.
	t, _ := tail.TailFile(mFile, tail.Config{Follow: true})
	for line := range t.Lines {
		logText := line.Text
		// Print log line
		klog.Info(logText)

		switch fileFormat {
		case commonv1beta1.TextFormat:
			// Check if log line contains metric from stop rules.
			isRuleLine := false
			for _, name := range liveRuleMetrics {
				if strings.Contains(logText, name) {
					isRuleLine = true
					break
				}
			}
			// If log line doesn't contain appropriate metric, continue track file.
			if !isRuleLine {
				continue
			}

			// If log line contains appropriate metric, find all submatches from metric filters.
			for _, metricReg := range metricRegList {
				matchStrings := metricReg.FindAllStringSubmatch(logText, -1)
				for _, subMatchList := range matchStrings {
					if len(subMatchList) < 3 {
						continue
					}
					// Submatch must have metric name and float value
					metricName := strings.TrimSpace(subMatchList[1])
					metricValue, err := strconv.ParseFloat(strings.TrimSpace(subMatchList[2]), 64)
					if err != nil {
						klog.Fatalf("Unable to parse value %v to float for metric %v", metricValue, metricName)
					}

					// liveRuleMetrics contains array of EarlyStoppingRules Name that has not been reached yet.
					for _, name := range liveRuleMetrics {
						if metricName != name {
							continue
						}
						err = rules.UpdateMetric(name, metricValue)
						if err != nil {
							klog.Fatalf("Unable to UpdateMetric %s %v", name, err)
						}
					}
				}
			}
		case commonv1beta1.JsonFormat:
			var logJsonObj map[string]any
			if err = json.Unmarshal([]byte(logText), &logJsonObj); err != nil {
				klog.Fatalf("Failed to unmarshal logs in %v format, log: %s, error: %v", commonv1beta1.JsonFormat, logText, err)
			}

			// liveRuleMetrics contains array of EarlyStoppingRules Name that has not been reached yet.
			for _, name := range liveRuleMetrics {
				value, exist := logJsonObj[name].(string)
				if !exist {
					continue
				}
				metricValue, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
				if err != nil {
					klog.Fatalf("Unable to parse value %v to float for metric %v", metricValue, name)
				}
				err = rules.UpdateMetric(name, metricValue)
				if err != nil {
					klog.Fatalf("Unable to UpdateMetric %s %v", name, err)
				}
			}
		default:
			klog.Fatalf("Format must be set to %v or %v", commonv1beta1.TextFormat, commonv1beta1.JsonFormat)
		}

		liveRuleMetrics = rules.LiveMetrics()
		// If liveRuleMetrics array is empty, Trial is early stopped.
		if len(liveRuleMetrics) == 0 {
			klog.Info("Training container is early stopped")
			isEarlyStopped = true

			// Create ".pid" file with "early-stopped" line.
			// Which means that training is early stopped and Trial status is updated.
			markFile := filepath.Join(filepath.Dir(mFile), fmt.Sprintf("%d.pid", mainProcPid))
			_, err := os.Create(markFile)
			if err != nil {
				klog.Fatalf("Create mark file %v error: %v", markFile, err)
			}

			err = os.WriteFile(markFile, []byte(common.TrainingEarlyStopped), 0)
			if err != nil {
				klog.Fatalf("Write to file %v error: %v", markFile, err)
			}

			// Get child process from main PID.
			childProc, err := mainProc.Children()
			if err != nil {
				klog.Fatalf("Get children proceses for main PID: %v failed: %v", mainProcPid, err)
			}

			// TODO (andreyvelich): Currently support only single child process.
			if len(childProc) != 1 {
				klog.Fatalf("Multiple children processes are not supported. Children processes: %v", childProc)
			}

			// Terminate the child process.
			err = childProc[0].Terminate()
			if err != nil {
				klog.Fatalf("Unable to terminate child process %v, error: %v", childProc[0], err)
			}

			// Report metrics to DB.
			reportMetrics(filters, fileFormat)

			// Wait until main process is completed.
			timeout := 60 * time.Second
			endTime := time.Now().Add(timeout)
			isProcRunning := true
			for isProcRunning && time.Now().Before(endTime) {
				isProcRunning, err = mainProc.IsRunning()
				// Ignore "no such file error". It means that process is complete.
				if err != nil && !os.IsNotExist(err) {
					klog.Fatalf("Check process status for main PID: %v failed: %v", mainProcPid, err)
				}
			}

			// Create connection and client for Early Stopping service.
			conn, err := grpc.Dial(*earlyStopServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				klog.Fatalf("Could not connect to Early Stopping service, error: %v", err)
			}
			c := api.NewEarlyStoppingClient(conn)

			setTrialStatusReq := &api.SetTrialStatusRequest{
				TrialName: *trialName,
			}

			// Send request to change Trial status to early stopped.
			_, err = c.SetTrialStatus(context.Background(), setTrialStatusReq)
			if err != nil {
				klog.Fatalf("Set Trial status error: %v", err)
			}
			conn.Close()

			klog.Infof("Trial status is successfully updated")
		}
	}
}

func main() {
	flag.Var(&stopRules, "stop-rule", "The list of early stopping stop rules")
	flag.Parse()
	klog.Infof("Trial Name: %s", *trialName)

	var filters []string
	if len(*metricFilters) != 0 {
		filters = strings.Split(*metricFilters, ";")
	}

	fileFormat := commonv1beta1.FileFormat(*metricsFileFormat)

	// If stop rule is set we need to parse metrics during run.
	if len(stopRules) != 0 {
		go watchMetricsFile(*metricsFilePath, stopRules, filters, fileFormat)
	} else {
		go printMetricsFile(*metricsFilePath)
	}

	waitAll, _ := strconv.ParseBool(*waitAllProcesses)

	wopts := common.WaitPidsOpts{
		PollInterval:           *pollInterval,
		Timeout:                *timeout,
		WaitAll:                waitAll,
		CompletedMarkedDirPath: filepath.Dir(*metricsFilePath),
	}
	if err := common.WaitMainProcesses(wopts); err != nil {
		klog.Fatalf("Failed to wait for worker container: %v", err)
	}

	// If training was not early stopped, report the metrics.
	if !isEarlyStopped {
		reportMetrics(filters, fileFormat)
	}
}

func reportMetrics(filters []string, fileFormat commonv1beta1.FileFormat) {

	conn, err := grpc.Dial(*dbManagerServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Fatalf("Could not connect to DB manager service, error: %v", err)
	}
	defer conn.Close()
	c := api.NewDBManagerClient(conn)
	ctx := context.Background()
	var metricList []string
	if len(*metricNames) != 0 {
		metricList = strings.Split(*metricNames, ";")
	}
	olog, err := filemc.CollectObservationLog(*metricsFilePath, metricList, filters, fileFormat)
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
