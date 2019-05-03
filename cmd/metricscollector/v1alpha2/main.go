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
	"log"

	api "github.com/kubeflow/katib/pkg/api/v1alpha2"
	"github.com/kubeflow/katib/pkg/util/v1alpha2/metricscollector"

	"google.golang.org/grpc"
)

var experimentName = flag.String("e", "", "Experiment Name")
var trialName = flag.String("t", "", "Trial Name")
var jobKind = flag.String("k", "", "Job Kind")
var namespace = flag.String("n", "", "NameSpace")
var managerSerivce = flag.String("m", "", "Katib Manager service")

func main() {
	flag.Parse()
	log.Printf("Experiment Name: %s, Trial Name: %s, Job Kind: %s", *experimentName, *trialName, *jobKind)
	conn, err := grpc.Dial(*managerSerivce, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)
	mc, err := metricscollector.NewMetricsCollector()
	if err != nil {
		log.Fatalf("Failed to create MetricsCollector: %v", err)
	}
	ctx := context.Background()
	getreq := &api.GetExperimentRequest{
		ExperimentName: *experimentName,
	}
	getrep, err := c.GetExperiment(ctx, getreq)
	if err != nil {
		log.Fatalf("Failed to GetExperiment: %v", err)
	}
	olog, err := mc.CollectObservationLog(*trialName, *jobKind, getrep.Experiment.ExperimentSpec.Objective.ObjectiveMetricName, getrep.Experiment.ExperimentSpec.Objective.AdditionalMetricsNames, *namespace)
	if err != nil {
		log.Fatalf("Failed to collect logs: %v", err)
	}
	reportreq := &api.ReportObservationLogRequest{
		TrialName:      *trialName,
		ObservationLog: olog,
	}
	_, err = c.ReportObservationLog(ctx, reportreq)
	if err != nil {
		log.Fatalf("Failed to Report logs: %v", err)
	}
	log.Printf("Metrics reported. :\n%v", olog)
	return
}
