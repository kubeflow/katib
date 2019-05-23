package main

import (
	"context"
	"flag"
	"io/ioutil"
	"strconv"

	api "github.com/kubeflow/katib/pkg/api/v1alpha1"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/klog"
)

var managerAddr = flag.String("s", "127.0.0.1:6789", "Endpoint of manager default 127.0.0.1:6789")
var suggestArgo = flag.String("a", "random", "Suggestion Algorithm (random, grid, hyperband)")
var requestnum = flag.Int("r", 2, "Request number for random Suggestions (default: 2)")
var suggestionConfFile = flag.String("c", "", "File path to suggestion config.")

var studyConfig = api.StudyConfig{}
var suggestionConfig = api.SetSuggestionParametersRequest{}

const TimeOut = 600

var trials = map[string]*api.Trial{}

func main() {
	ctx := context.Background()
	readConfigs()
	conn, err := grpc.Dial(*managerAddr, grpc.WithInsecure())
	if err != nil {
		klog.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)

	//CreateStudy
	studyID := CreateStudy(c)

	//SetSuggestParam
	paramID := setSuggestionParam(c, studyID)

	iter := 1

	//GetSuggestion
	if *suggestArgo == "hyperband" {
		for true {
			getSuggestReply := getSuggestion(c, studyID, paramID)
			checkSuggestions(getSuggestReply, iter)
			if len(getSuggestReply.Trials) == 0 {
				klog.Infof("Hyperband ended")
				break
			}

			//RegisterWorkers
			workerIds := registerWorkers(c, studyID, getSuggestReply)
			if workerIds == nil {
				klog.Fatalf("Register Workers error")
			}

			//UpdateWorkerState
			for _, w := range workerIds {
				uwsreq := &api.UpdateWorkerStateRequest{
					WorkerId: w,
					Status:   api.State_COMPLETED,
				}
				c.UpdateWorkerState(ctx, uwsreq)
			}

			getMetricsRequest := &api.GetMetricsRequest{
				StudyId:   studyID,
				WorkerIds: workerIds,
			}
			//GetMetrics
			getMetricsReply, err := c.GetMetrics(ctx, getMetricsRequest)
			if err != nil {
				continue
			}

			mlSet := getMetricsReply.MetricsLogSets

			//add dummy metricsValueTime
			for i := range mlSet {
				for j := range mlSet[i].MetricsLogs {
					mlSet[i].MetricsLogs[j].Values = append(mlSet[i].MetricsLogs[j].Values, &api.MetricsValueTime{Time: "2018-01-01T12:00:00.999999999Z", Value: "1.0"})
					mlSet[i].MetricsLogs[j].Values = append(mlSet[i].MetricsLogs[j].Values, &api.MetricsValueTime{Time: "2019-02-02T13:30:00.999999999Z", Value: "2.0"})
				}
			}

			//ReportMetrics
			rmlreq := &api.ReportMetricsLogsRequest{
				StudyId:        studyID,
				MetricsLogSets: mlSet,
			}
			c.ReportMetricsLogs(ctx, rmlreq)

			checkWorkersResult(c, studyID)
			iter++
		}
	} else {
		getSuggestReply := getSuggestion(c, studyID, paramID)
		checkSuggestions(getSuggestReply, iter)
	}
	DeleteStudy(c, studyID)
	conn.Close()
	klog.Infof("E2E test OK!")
}

func readConfigs() {
	flag.Parse()
	buf, err := ioutil.ReadFile("study-config.yml")
	if err != nil {
		klog.Fatalf("fail to read study-config yaml")
	}
	err = yaml.Unmarshal(buf, &studyConfig)
	if err != nil {
		klog.Fatalf("fail to Unmarshal yaml")
	}
	studyConfig.Name += *suggestArgo

	if *suggestionConfFile != "" {
		buf, err = ioutil.ReadFile(*suggestionConfFile)
		if err != nil {
			klog.Fatalf("fail to read suggestion-config yaml")
		}
	}
	err = yaml.Unmarshal(buf, &suggestionConfig)
	if err != nil {
		klog.Fatalf("fail to Unmarshal yaml")
	}
}

func CreateStudy(c api.ManagerClient) string {
	ctx := context.Background()
	createStudyreq := &api.CreateStudyRequest{
		StudyConfig: &studyConfig,
	}
	createStudyreply, err := c.CreateStudy(ctx, createStudyreq)
	if err != nil {
		klog.Fatalf("StudyConfig Error %v", err)
	}
	studyID := createStudyreply.StudyId
	klog.Infof("Study ID %s", studyID)
	getStudyreq := &api.GetStudyRequest{
		StudyId: studyID,
	}
	getStudyReply, err := c.GetStudy(ctx, getStudyreq)
	if err != nil {
		klog.Fatalf("GetConfig Error %v", err)
	}
	klog.Infof("Study ID %s StudyConf %v", studyID, getStudyReply.StudyConfig)
	return studyID
}

func DeleteStudy(c api.ManagerClient, studyID string) {
	ctx := context.Background()
	deleteStudyreq := &api.DeleteStudyRequest{
		StudyId: studyID,
	}
	if _, err := c.DeleteStudy(ctx, deleteStudyreq); err != nil {
		klog.Fatalf("DeleteStudy error %v", err)
	}
	getStudyreq := &api.GetStudyRequest{
		StudyId: studyID,
	}
	getStudyReply, _ := c.GetStudy(ctx, getStudyreq)
	if getStudyReply != nil && getStudyReply.StudyConfig != nil {
		klog.Fatalf("Failed to delete Study %s", studyID)
	}
	getTrialsRequest := &api.GetTrialsRequest{
		StudyId: studyID,
	}
	gtrep, _ := c.GetTrials(ctx, getTrialsRequest)
	if gtrep != nil && len(gtrep.Trials) > 0 {
		klog.Fatalf("Failed to delete Trials of Study %s", studyID)
	}
	getWorkersRequest := &api.GetWorkersRequest{
		StudyId: studyID,
	}
	gwrep, _ := c.GetWorkers(ctx, getWorkersRequest)
	if gwrep != nil && len(gwrep.Workers) > 0 {
		klog.Fatalf("Failed to delete Workers of Study %s", studyID)
	}
	klog.Infof("Study %s is deleted", studyID)
}

func setSuggestionParam(c api.ManagerClient, studyID string) string {
	ctx := context.Background()
	switch *suggestArgo {
	case "random":
		return ""
	case "grid":
		suggestionConfig.StudyId = studyID
		setSuggesitonParameterReply, err := c.SetSuggestionParameters(ctx, &suggestionConfig)
		if err != nil {
			klog.Fatalf("SetConfig Error %v", err)
		}
		klog.Infof("Grid suggestion prameter ID %s", setSuggesitonParameterReply.ParamId)
		return setSuggesitonParameterReply.ParamId
	case "hyperband":
		suggestionConfig.StudyId = studyID
		setSuggesitonParameterReply, err := c.SetSuggestionParameters(ctx, &suggestionConfig)
		if err != nil {
			klog.Fatalf("SetConfig Error %v", err)
		}
		klog.Infof("HyperBand suggestion prameter ID %s", setSuggesitonParameterReply.ParamId)
		return setSuggesitonParameterReply.ParamId
	}
	return ""

}

func getSuggestion(c api.ManagerClient, studyID string, paramID string) *api.GetSuggestionsReply {
	ctx := context.Background()
	var getSuggestRequest *api.GetSuggestionsRequest
	switch *suggestArgo {
	case "random":
		//Random suggestion doesn't need suggestion parameter
		getSuggestRequest = &api.GetSuggestionsRequest{
			StudyId:             studyID,
			SuggestionAlgorithm: "random",
			RequestNumber:       int32(*requestnum),
		}

	case "grid":
		getSuggestRequest = &api.GetSuggestionsRequest{
			StudyId:             studyID,
			SuggestionAlgorithm: "grid",
			RequestNumber:       0,
			//RequestNumber=0 means get all grids.
			ParamId: paramID,
		}
	case "hyperband":
		getSuggestRequest = &api.GetSuggestionsRequest{
			StudyId:             studyID,
			SuggestionAlgorithm: "hyperband",
			RequestNumber:       0,
			ParamId:             paramID,
		}
	}

	getSuggestReply, err := c.GetSuggestions(ctx, getSuggestRequest)
	if err != nil {
		klog.Fatalf("GetSuggestion Error %v \nRequest %v", err, getSuggestRequest)
	}
	klog.Infof("Get " + *suggestArgo + " Suggestions:")
	for _, t := range getSuggestReply.Trials {
		klog.Infof("%v", t)
	}
	return getSuggestReply
}

func checkSuggestions(getSuggestReply *api.GetSuggestionsReply, iter int) bool {
	switch *suggestArgo {
	case "random":
		if len(getSuggestReply.Trials) != *requestnum {
			klog.Fatalf("Number of Random suggestion incorrect. Expected %d Got %d", *requestnum, len(getSuggestReply.Trials))
		}
	case "grid":
		if len(getSuggestReply.Trials) != 4 {
			klog.Fatalf("Number of Grid suggestion incorrect. Expected %d Got %d", 4, len(getSuggestReply.Trials))
		}
		min, max := 1.0, 1.0
		for _, m := range studyConfig.ParameterConfigs.Configs {
			if m.Name == "learning-rate" {
				min, _ = strconv.ParseFloat(m.Feasible.Min, 8)
				max, _ = strconv.ParseFloat(m.Feasible.Max, 8)
			}
		}
		learningRate := 1.0
		for _, l := range suggestionConfig.SuggestionParameters {
			if l.Name == "learning-rate" {
				learningRate, _ = strconv.ParseFloat(l.Value, 8)
			}
		}
		for i, trial := range getSuggestReply.Trials {
			for _, param := range trial.ParameterSet {
				if param.Name == "learning-rate" && learningRate != 0 {
					expValue := min + (max-min)/(learningRate-1)*float64(i)
					if param.Value != strconv.FormatFloat(expValue, 'f', 4, 64) {
						klog.Infof("Grid point incorrect. Expected %v Got %v", strconv.FormatFloat(expValue, 'f', 4, 64), param.Value)
					}
				}
			}
		}
	case "hyperband":
		if iter == 1 {
			if len(getSuggestReply.Trials) != 3 {
				klog.Fatalf("Number of Hyperband suggestion incorrect. Expected %d Got %d", 3, len(getSuggestReply.Trials))
			}
		} else if iter == 2 {
			if len(getSuggestReply.Trials) != 1 {
				klog.Fatalf("Number of Hyperband suggestion incorrect. Expected %d Got %d", 1, len(getSuggestReply.Trials))
			}
		} else if iter == 3 {
			if len(getSuggestReply.Trials) != 0 {
				klog.Fatalf("Number of Hyperband suggestion incorrect. Expected %d Got %d", 0, len(getSuggestReply.Trials))
			}
		}
	}
	klog.Infof("Check suggestion passed!")
	return true
}

func registerWorkers(c api.ManagerClient, studyId string, getSuggestReply *api.GetSuggestionsReply) []string {
	ctx := context.Background()
	workerIds := make([]string, len(getSuggestReply.Trials))
	for i, t := range getSuggestReply.Trials {
		worker := &api.Worker{
			StudyId: studyId,
			TrialId: t.TrialId,
		}
		workerreq := &api.RegisterWorkerRequest{
			Worker: worker,
		}
		workerrep, err := c.RegisterWorker(ctx, workerreq)
		if err != nil {
			klog.Fatalf("RegisterWorker Error %v", err)
		}

		workerIds[i] = workerrep.WorkerId
		saveModelRequest := &api.SaveModelRequest{
			Model: &api.ModelInfo{
				StudyName:  studyConfig.Name,
				WorkerId:   workerrep.WorkerId,
				Parameters: t.ParameterSet,
				Metrics:    []*api.Metrics{},
				ModelPath:  "pvc:/Path/to/Model",
			},
			DataSet: &api.DataSetInfo{
				Name: "Mnist",
				Path: "/path/to/data",
			},
		}
		_, err = c.SaveModel(ctx, saveModelRequest)
		if err != nil {
			klog.Fatalf("SaveModel Error %v", err)
		}
		klog.Infof("WorkerID %s start \n", workerrep.WorkerId)
		trials[workerrep.WorkerId] = t
	}
	return workerIds
}

func checkWorkersResult(c api.ManagerClient, studyID string) bool {
	ctx := context.Background()
	getMetricsRequest := &api.GetMetricsRequest{
		StudyId: studyID,
	}
	//GetMetrics
	getMetricsReply, err := c.GetMetrics(ctx, getMetricsRequest)
	if err != nil {
		klog.Fatalf("Fataled to Get Metrics")
	}

	for _, mls := range getMetricsReply.MetricsLogSets {
		for _, p := range trials[mls.WorkerId].ParameterSet {
			for _, ml := range mls.MetricsLogs {
				if p.Name == ml.Name {
					if p.Value != ml.Values[len(ml.Values)-1].Value {
						klog.Fatalf("Output %s is mismuched to Input %s", ml.Values[len(ml.Values)-1], p.Value)
						return false
					}
				}
			}
		}
	}
	klog.Info("Input Output check passed")
	return true
}
