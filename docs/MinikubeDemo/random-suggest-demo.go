package main

import (
	"context"
	"log"
	"time"

	"github.com/kubeflow/katib/pkg/api"
	"google.golang.org/grpc"
)

const (
	manager = "192.168.99.100:30678"
)

var studyConfig = api.StudyConfig{
	Name:                          "random-demo",
	Owner:                         "katib",
	OptimizationType:              api.OptimizationType_MAXIMIZE,
	OptimizationGoal:              0.99,
	DefaultSuggestionAlgorithm:    "random",
	DefaultEarlyStoppingAlgorithm: "medianstopping",
	ObjectiveValueName:            "Validation-accuracy",
	Metrics: []string{
		"accuracy",
	},
	ParameterConfigs: &api.StudyConfig_ParameterConfigs{
		Configs: []*api.ParameterConfig{
			&api.ParameterConfig{
				Name:          "--lr",
				ParameterType: api.ParameterType_DOUBLE,
				Feasible: &api.FeasibleSpace{
					Min: "0.03",
					Max: "0.07",
				},
			},
		},
	},
}

func main() {
	conn, err := grpc.Dial(manager, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)
	createStudyreq := &api.CreateStudyRequest{
		StudyConfig: &studyConfig,
	}
	createStudyreply, err := c.CreateStudy(context.Background(), createStudyreq)
	if err != nil {
		log.Fatalf("StudyConfig Error %v", err)
	}
	studyId := createStudyreply.StudyId
	log.Printf("Study ID %s", studyId)
	getStudyreq := &api.GetStudyRequest{
		StudyId: studyId,
	}
	getStudyReply, err := c.GetStudy(context.Background(), getStudyreq)
	if err != nil {
		log.Fatalf("GetConfig Error %v", err)
	}
	log.Printf("Study ID %s StudyConf%v", studyId, getStudyReply.StudyConfig)
	getRandomSuggestRequest := &api.GetSuggestionsRequest{
		StudyId:             studyId,
		SuggestionAlgorithm: "random",
		RequestNumber:       2,
	}
	getRandomSuggestReply, err := c.GetSuggestions(context.Background(), getRandomSuggestRequest)
	if err != nil {
		log.Fatalf("GetSuggestion Error %v", err)
	}
	log.Printf("Get Random Suggestions %v", getRandomSuggestReply.Trials)
	workerIds := make([]string, len(getRandomSuggestReply.Trials))
	workerParameter := make(map[string][]*api.Parameter)
	for i, t := range getRandomSuggestReply.Trials {
		rtr := &api.RunTrialRequest{
			StudyId: studyId,
			TrialId: t.TrialId,
			Runtime: "kubernetes",
			WorkerConfig: &api.WorkerConfig{
				Image: "mxnet/python",
				Command: []string{
					"python",
					"/mxnet/example/image-classification/train_mnist.py",
					"--batch-size=64",
				},
				Gpu:       0,
				Scheduler: "default-scheduler",
			},
		}
		for _, p := range t.ParameterSet {
			rtr.WorkerConfig.Command = append(rtr.WorkerConfig.Command, p.Name)
			rtr.WorkerConfig.Command = append(rtr.WorkerConfig.Command, p.Value)
		}
		workerReply, err := c.RunTrial(context.Background(), rtr)
		if err != nil {
			log.Fatalf("RunTrial Error %v", err)
		}
		workerIds[i] = workerReply.WorkerId
		workerParameter[workerReply.WorkerId] = t.ParameterSet
		saveModelRequest := &api.SaveModelRequest{
			Model: &api.ModelInfo{
				StudyName:  studyConfig.Name,
				WorkerId:   workerReply.WorkerId,
				Parameters: t.ParameterSet,
				Metrics:    []*api.Metrics{},
				ModelPath:  "pvc:/Path/to/Model",
			},
			DataSet: &api.DataSetInfo{
				Name: "Mnist",
				Path: "/path/to/data",
			},
		}
		_, err = c.SaveModel(context.Background(), saveModelRequest)
		if err != nil {
			log.Fatalf("SaveModel Error %v", err)
		}
		log.Printf("WorkerID %s start\n", workerReply.WorkerId)
	}
	for true {
		time.Sleep(10 * time.Second)
		getMetricsRequest := &api.GetMetricsRequest{
			StudyId:   studyId,
			WorkerIds: workerIds,
		}
		getMetricsReply, err := c.GetMetrics(context.Background(), getMetricsRequest)
		if err != nil {
			log.Printf("GetMetErr %v", err)
			continue
		}
		for _, mls := range getMetricsReply.MetricsLogSets {
			if len(mls.MetricsLogs) > 0 {
				log.Printf("WorkerID %s :", mls.WorkerId)
				//Only Metrics can be updated.
				saveModelRequest := &api.SaveModelRequest{
					Model: &api.ModelInfo{
						StudyName: studyConfig.Name,
						WorkerId:  mls.WorkerId,
						Metrics:   []*api.Metrics{},
					},
				}
				for _, ml := range mls.MetricsLogs {
					if len(ml.Values) > 0 {
						log.Printf("\t Metrics Name %s Value %v", ml.Name, ml.Values[len(ml.Values)-1])
						saveModelRequest.Model.Metrics = append(saveModelRequest.Model.Metrics, &api.Metrics{Name: ml.Name, Value: ml.Values[len(ml.Values)-1]})
					}
				}
				_, err = c.SaveModel(context.Background(), saveModelRequest)
				if err != nil {
					log.Fatalf("SaveModel Error %v", err)
				}
			}
		}
		getWorkerRequest := &api.GetWorkersRequest{StudyId: studyId}
		getWorkerReply, err := c.GetWorkers(context.Background(), getWorkerRequest)
		if err != nil {
			log.Fatalf("GetWorker Error %v", err)
		}
		completeCount := 0
		for _, w := range getWorkerReply.Workers {
			if w.Status == api.State_COMPLETED {
				completeCount++
			}
		}
		if completeCount == len(getWorkerReply.Workers) {
			log.Printf("All Worker Completed!")
			break
		}
	}
}
