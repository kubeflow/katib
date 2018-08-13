package main

import (
	"context"
	"flag"
	"log"

	"github.com/kubeflow/katib/pkg"
	api "github.com/kubeflow/katib/pkg/api"
	"github.com/kubeflow/katib/pkg/manager/metricscollector"

	"google.golang.org/grpc"
)

var studyId = flag.String("s", "", "Study ID")
var workerId = flag.String("w", "", "Worker ID")
var namespace = flag.String("n", "", "NameSpace")

func main() {
	flag.Parse()
	log.Printf("Study ID: %s, Worker ID: %s", *studyId, *workerId)
	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
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
	screq := &api.GetStudyRequest{
		StudyId: *studyId,
	}
	screp, err := c.GetStudy(ctx, screq)
	if err != nil {
		log.Fatalf("Failed to GetStudyConf: %v", err)
	}
	mls, err := mc.CollectWorkerLog(*workerId, screp.StudyConfig.ObjectiveValueName, screp.StudyConfig.Metrics, *namespace)
	if err != nil {
		log.Fatalf("Failed to collect logs: %v", err)
	}
	rmreq := &api.ReportMetricsRequest{
		StudyId:        *studyId,
		MetricsLogSets: []*api.MetricsLogSet{mls},
	}
	_, err = c.ReportMetrics(ctx, rmreq)
	if err != nil {
		log.Fatalf("Failed to Report logs: %v", err)
	}
	log.Printf("Metrics reported. :\n%v", mls)
	return
}
