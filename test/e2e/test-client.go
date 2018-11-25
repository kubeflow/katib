package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	//	"time"

	"github.com/kubeflow/katib/pkg/api"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
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
	readConfigs()
	conn, err := grpc.Dial(*managerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)

	//CreateStudy
	studyID := CreateStudy(c)

	//SetSuggestParam
	paramID := setSuggestionParam(c, studyID)

	//GetSuggestion
	if *suggestArgo == "hyperband" {
		for true {
			getSuggestReply := getSuggestion(c, studyID, paramID)
			checkSuggestions(getSuggestReply)
			if len(getSuggestReply.Trials) == 0 {
				log.Printf("Hyperband ended")
				break
			}
			//RunTrials
			//workerIds := runTrials(c, studyId, getSuggestReply)

			////GetWorkersCheck
			//getWorker(c, studyId, getSuggestReply, workerIds)
			//log.Println("GetWorkers checked!")

			//iter := 0

			//for !isCompletedAllWorker(c, studyId) {
			//	if iter > TimeOut {
			//		log.Fatal("GetMetrics Timeout.")
			//	}
			//	time.Sleep(1 * time.Second)
			//	getMetricsRequest := &api.GetMetricsRequest{
			//		StudyId:   studyId,
			//		WorkerIds: workerIds,
			//	}
			//	//GetMetrics
			//	getMetricsReply, err := c.GetMetrics(ctx, getMetricsRequest)
			//	if err != nil {
			//		continue
			//	}
			//	//Save or Update model on ModelDB
			//	SaveOrUpdateModel(c, getMetricsReply)
			//	iter++
			//}
			//checkWorkersResult(c, studyId)

		}
	} else {
		getSuggestReply := getSuggestion(c, studyID, paramID)
		checkSuggestions(getSuggestReply)
	}
	conn.Close()
	log.Println("E2E test OK!")
}

func readConfigs() {
	flag.Parse()
	buf, err := ioutil.ReadFile("study-config.yml")
	if err != nil {
		log.Fatalf("fail to read study-config yaml")
	}
	err = yaml.Unmarshal(buf, &studyConfig)
	if err != nil {
		log.Fatalf("fail to Unmarshal yaml")
	}
	studyConfig.Name += *suggestArgo

	if *suggestionConfFile != "" {
		buf, err = ioutil.ReadFile(*suggestionConfFile)
		if err != nil {
			log.Fatalf("fail to read suggestion-config yaml")
		}
	}
	err = yaml.Unmarshal(buf, &suggestionConfig)
	if err != nil {
		log.Fatalf("fail to Unmarshal yaml")
	}
}

func CreateStudy(c api.ManagerClient) string {
	ctx := context.Background()
	createStudyreq := &api.CreateStudyRequest{
		StudyConfig: &studyConfig,
	}
	createStudyreply, err := c.CreateStudy(ctx, createStudyreq)
	if err != nil {
		log.Fatalf("StudyConfig Error %v", err)
	}
	studyID := createStudyreply.StudyId
	log.Printf("Study ID %s", studyID)
	getStudyreq := &api.GetStudyRequest{
		StudyId: studyID,
	}
	getStudyReply, err := c.GetStudy(ctx, getStudyreq)
	if err != nil {
		log.Fatalf("GetConfig Error %v", err)
	}
	log.Printf("Study ID %s StudyConf %v", studyID, getStudyReply.StudyConfig)
	return studyID
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
			log.Fatalf("SetConfig Error %v", err)
		}
		log.Printf("Grid suggestion prameter ID %s", setSuggesitonParameterReply.ParamId)
		return setSuggesitonParameterReply.ParamId
	case "hyperband":
		suggestionConfig.StudyId = studyID
		setSuggesitonParameterReply, err := c.SetSuggestionParameters(ctx, &suggestionConfig)
		if err != nil {
			log.Fatalf("SetConfig Error %v", err)
		}
		log.Printf("HyperBand suggestion prameter ID %s", setSuggesitonParameterReply.ParamId)
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
		log.Fatalf("GetSuggestion Error %v \nRequest %v", err, getSuggestRequest)
	}
	log.Println("Get " + *suggestArgo + " Suggestions:")
	for _, t := range getSuggestReply.Trials {
		log.Printf("%v", t)
	}
	return getSuggestReply
}

func checkSuggestions(getSuggestReply *api.GetSuggestionsReply) bool {
	switch *suggestArgo {
	case "random":
		if len(getSuggestReply.Trials) != *requestnum {
			log.Fatalf("Number of Random suggestion incorrect. Expected %d Got %d", *requestnum, len(getSuggestReply.Trials))
		}
	case "grid":
		if len(getSuggestReply.Trials) != 4 {
			log.Fatalf("Number of Grid suggestion incorrect. Expected %d Got %d", 4, len(getSuggestReply.Trials))
		}
	}
	log.Println("Check suggestion passed!")
	return true
}

//func runTrials(c api.ManagerClient, studyId string, getSuggestReply *api.GetSuggestionsReply) []string {
//	ctx := context.Background()
//	workerIds := make([]string, len(getSuggestReply.Trials))
//	workerParameter := make(map[string][]*api.Parameter)
//	for i, t := range getSuggestReply.Trials {
//		wc := workerConfig
//		rtr := &api.RunTrialRequest{
//			StudyId:      studyId,
//			TrialId:      t.TrialId,
//			Runtime:      "kubernetes",
//			WorkerConfig: &wc,
//		}
//		for _, p := range t.ParameterSet {
//			rtr.WorkerConfig.Command = append(rtr.WorkerConfig.Command, p.Name+"="+p.Value)
//		}
//		workerReply, err := c.RunTrial(ctx, rtr)
//		if err != nil {
//			log.Fatalf("RunTrial Error %v", err)
//		}
//		workerIds[i] = workerReply.WorkerId
//		workerParameter[workerReply.WorkerId] = t.ParameterSet
//		saveModelRequest := &api.SaveModelRequest{
//			Model: &api.ModelInfo{
//				StudyName:  studyConfig.Name,
//				WorkerId:   workerReply.WorkerId,
//				Parameters: t.ParameterSet,
//				Metrics:    []*api.Metrics{},
//				ModelPath:  "pvc:/Path/to/Model",
//			},
//			DataSet: &api.DataSetInfo{
//				Name: "Mnist",
//				Path: "/path/to/data",
//			},
//		}
//		_, err = c.SaveModel(ctx, saveModelRequest)
//		if err != nil {
//			log.Fatalf("SaveModel Error %v", err)
//		}
//		log.Printf("WorkerID %s start\n", workerReply.WorkerId)
//		trials[workerReply.WorkerId] = t
//	}
//	return workerIds
//}
//
//func getWorker(c api.ManagerClient, studyId string, getSuggestReply *api.GetSuggestionsReply, workerIds []string) {
//	ctx := context.Background()
//	getByTID := []*api.Worker{}
//	getByWID := []*api.Worker{}
//	gwr := &api.GetWorkersRequest{
//		StudyId: studyId,
//	}
//	_, err := c.GetWorkers(ctx, gwr)
//	if err != nil {
//		log.Fatalf("GwtWorker by Study ID Error %v", err)
//	}
//	for _, t := range getSuggestReply.Trials {
//		gwr := &api.GetWorkersRequest{
//			StudyId: studyId,
//			TrialId: t.TrialId,
//		}
//		gwrep, err := c.GetWorkers(ctx, gwr)
//		if err != nil {
//			log.Fatalf("GwtWorker by Trial ID Error %v", err)
//		}
//		getByTID = append(getByTID, gwrep.Workers...)
//	}
//	for _, w := range workerIds {
//		gwr = &api.GetWorkersRequest{
//			StudyId:  studyId,
//			WorkerId: w,
//		}
//		gwrep, err := c.GetWorkers(ctx, gwr)
//		if err != nil {
//			log.Fatalf("GwtWorker by Worker ID Error %v", err)
//		}
//		getByWID = append(getByWID, gwrep.Workers...)
//	}
//	if len(getByTID) != len(getByWID) {
//		log.Fatalf("GwtWorker by Worker ID is mismatch with GwtWorker by Trial ID \n by Trial ID %v \n by Worker ID %v", getByTID, getByWID)
//	}
//}

//func SaveOrUpdateModel(c api.ManagerClient, getMetricsReply *api.GetMetricsReply) {
//	ctx := context.Background()
//	for _, mls := range getMetricsReply.MetricsLogSets {
//		if len(mls.MetricsLogs) > 0 {
//			log.Printf("WorkerID %s :", mls.WorkerId)
//			//Only Metrics can be updated.
//			saveModelRequest := &api.SaveModelRequest{
//				Model: &api.ModelInfo{
//					StudyName: studyConfig.Name,
//					WorkerId:  mls.WorkerId,
//					Metrics:   []*api.Metrics{},
//				},
//			}
//			for _, ml := range mls.MetricsLogs {
//				if len(ml.Values) > 0 {
//					log.Printf("\t Metrics Name %s Value %v", ml.Name, ml.Values[len(ml.Values)-1])
//					saveModelRequest.Model.Metrics = append(saveModelRequest.Model.Metrics, &api.Metrics{Name: ml.Name, Value: ml.Values[len(ml.Values)-1]})
//				}
//			}
//			_, err := c.SaveModel(ctx, saveModelRequest)
//			if err != nil {
//				log.Fatalf("SaveModel Error %v", err)
//			}
//		}
//	}
//}

func isCompletedAllWorker(c api.ManagerClient, studyID string) bool {
	ctx := context.Background()
	getWorkerRequest := &api.GetWorkersRequest{StudyId: studyID}
	getWorkerReply, err := c.GetWorkers(ctx, getWorkerRequest)
	if err != nil {
		log.Fatalf("GetWorker Error %v", err)
	}
	for _, w := range getWorkerReply.Workers {
		if w.Status != api.State_COMPLETED {
			return false
		}
	}
	log.Println("All Worker Completed")
	return true
}

func checkWorkersResult(c api.ManagerClient, studyID string) bool {
	ctx := context.Background()
	getMetricsRequest := &api.GetMetricsRequest{
		StudyId: studyID,
	}
	//GetMetrics
	getMetricsReply, err := c.GetMetrics(ctx, getMetricsRequest)
	if err != nil {
		log.Fatalf("Fataled to Get Metrics")
	}

	for _, mls := range getMetricsReply.MetricsLogSets {
		for _, p := range trials[mls.WorkerId].ParameterSet {
			for _, ml := range mls.MetricsLogs {
				if p.Name == ml.Name {
					if p.Value != ml.Values[len(ml.Values)-1].Value {
						log.Fatalf("Output %s is mismuched to Input %s", ml.Values[len(ml.Values)-1], p.Value)
						return false
					}
				}
			}
		}
	}
	log.Println("Input Output check passed")
	return true
}
