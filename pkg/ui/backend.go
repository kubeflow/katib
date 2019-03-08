package ui

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kubeflow/katib/pkg"
	"github.com/kubeflow/katib/pkg/api"
	"github.com/kubeflow/katib/pkg/manager/studyjobclient"

	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	"google.golang.org/grpc"
	restclient "k8s.io/client-go/rest"

	"github.com/ghodss/yaml"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespace      = "default"
	port           = "9303"
	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "text/html; charset=utf-8")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	(*w).Header().Set("Access-Control-Expose-Headers", "Access-Control-*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

var config = parseKubernetesConfig()

func parseKubernetesConfig() *restclient.Config {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}
	return config
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

const maxMsgSize = 1<<31 - 1

type IDList struct {
	StudyId  string
	WorkerId string
	TrialId  string
}

type JobView struct {
	ID     string
	Name   string
	Status string
}

type TemplateView struct {
	Name string
	Yaml string
}

type KatibUIHandler struct {
	studyjobClient *studyjobclient.StudyjobClient
}

type TemplateResponse struct {
	TemplateType string
	Data         []TemplateView
}

func NewKatibUIHandler() *KatibUIHandler {
	sjc, err := studyjobclient.NewStudyjobClient(config)
	if err != nil {
		panic(err)
	}
	return &KatibUIHandler{
		studyjobClient: sjc,
	}
}

// func (k *KatibUIHandler) connectManager() (*grpc.ClientConn, api.ManagerClient, error) {
// 	opts := []grpc.DialOption{
// 		grpc.WithInsecure(),
// 		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
// 	}
// 	conn, err := grpc.Dial(pkg.ManagerAddr, opts...)
// 	if err != nil {
// 		log.Printf("Connect katib manager error %v", err)
// 		return nil, nil, err
// 	}
// 	c := api.NewManagerClient(conn)
// 	return conn, c, nil
// }

func (k *KatibUIHandler) connectManager() (*grpc.ClientConn, katibapi.ManagerClient) {
	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
	if err != nil {
		return nil, nil
	}
	c := katibapi.NewManagerClient(conn)
	return conn, c
}

func (k *KatibUIHandler) FetchHPJobs(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	jobs := make([]JobView, 0)

	sl, err := k.studyjobClient.GetStudyJobList()
	if err != nil {
		log.Printf("Get Study list failed %v", err)
		return
	}
	for _, sj := range sl.Items {
		if sj.Spec.ParameterConfigs != nil {
			jobs = append(jobs, JobView{
				Name:   sj.Spec.StudyName,
				ID:     sj.Status.StudyID,
				Status: string(sj.Status.Condition),
			})
		}
	}

	response, err := json.Marshal(jobs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)

}

func (k *KatibUIHandler) FetchNASJobs(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	jobs := make([]JobView, 0)

	sl, err := k.studyjobClient.GetStudyJobList()
	if err != nil {
		log.Printf("Get Study list failed %v", err)
		return
	}
	for _, sj := range sl.Items {
		if sj.Spec.NasConfig != nil {
			jobs = append(jobs, JobView{
				Name:   sj.Spec.StudyName,
				ID:     sj.Status.StudyID,
				Status: string(sj.Status.Condition),
			})
		}
	}

	response, err := json.Marshal(jobs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)

}

func (k *KatibUIHandler) SubmitYamlJob(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)

	job := katibv1alpha1.StudyJob{}
	if yamlContent, ok := data["yaml"].(string); ok {
		err := yaml.Unmarshal([]byte(yamlContent), &job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = k.studyjobClient.CreateStudyJob(&job, namespace)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func (k *KatibUIHandler) FetchJobInfo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	studyID := r.URL.Query()["id"][0]

	conn, c := k.connectManager()
	defer conn.Close()
	retText := "WorkerID,TrialID"
	gsrep, err := c.GetStudy(
		context.Background(),
		&api.GetStudyRequest{
			StudyId: studyID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	metricsList := map[string]int{}
	for i, m := range gsrep.StudyConfig.Metrics {
		retText += "," + m
		metricsList[m] = i
	}
	paramList := map[string]int{}
	for i, p := range gsrep.StudyConfig.ParameterConfigs.Configs {
		retText += "," + p.Name
		paramList[p.Name] = i + len(metricsList)
	}
	gwfirep, err := c.GetWorkerFullInfo(
		context.Background(),
		&api.GetWorkerFullInfoRequest{
			StudyId:       studyID,
			OnlyLatestLog: true,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	retText += "\n"
	for _, wfi := range gwfirep.WorkerFullInfos {
		restext := make([]string, len(metricsList)+len(paramList))
		for _, m := range wfi.MetricsLogs {
			if len(m.Values) > 0 {
				restext[metricsList[m.Name]] = m.Values[len(m.Values)-1].Value
			}
		}
		for _, p := range wfi.ParameterSet {
			restext[paramList[p.Name]] = p.Value
		}
		retText += wfi.Worker.WorkerId + "," + wfi.Worker.TrialId + "," + strings.Join(restext, ",") + "\n"
	}
	response, err := json.Marshal(retText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)

}

func (k *KatibUIHandler) FetchWorkerInfo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	studyID := r.URL.Query()["studyID"][0]
	workerID := r.URL.Query()["workerID"][0]
	conn, c := k.connectManager()
	defer conn.Close()

	defer conn.Close()
	retText := "symbol,time,value\n"
	gwfirep, err := c.GetWorkerFullInfo(
		context.Background(),
		&api.GetWorkerFullInfoRequest{
			StudyId:  studyID,
			WorkerId: workerID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(gwfirep.WorkerFullInfos) > 0 {
		for _, m := range gwfirep.WorkerFullInfos[0].MetricsLogs {
			pvtime := ""
			for _, v := range m.Values {
				mvtime, _ := time.Parse(time.RFC3339Nano, v.Time)
				ctime := mvtime.Format("2006-01-02T15:4:5")
				if pvtime != ctime {
					retText += m.Name + "," + ctime + "," + v.Value + "\n"
					pvtime = ctime
				}
			}
		}
	}

	response, err := json.Marshal(retText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (k *KatibUIHandler) FetchNASJobInfo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	studyID := r.URL.Query()["studyID"][0]
	// trialID := r.URL.Query()["trialID"][0]

	conn, c := k.connectManager()

	defer conn.Close()

	gtrep, err := c.GetTrials(
		context.Background(),
		&api.GetTrialsRequest{
			StudyId: studyID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	type NNView struct {
		Architecture string
		Decoder      string
		MetricsName  []string
		MetricsValue []string
	}

	response_raw := make([]NNView, 0)

	var architecture string
	var decoder string

	gwfirep, err := c.GetWorkerFullInfo(
		context.Background(),
		&api.GetWorkerFullInfoRequest{
			StudyId:       studyID,
			OnlyLatestLog: true,
		},
	)

	for _, tr := range gtrep.Trials {
		for _, parameter := range tr.ParameterSet {
			if parameter.Name == "architecture" {
				architecture = parameter.Value
			}
			if parameter.Name == "nn_config" {
				decoder = parameter.Value
			}
		}

		metricsName := make([]string, 0)
		metricsValue := make([]string, 0)
		trialID := tr.TrialId
		for _, wfi := range gwfirep.WorkerFullInfos {
			if wfi.Worker.TrialId == trialID {
				for _, metrics := range wfi.MetricsLogs {
					metricsName = append(metricsName, metrics.Name)
					metricsValue = append(metricsValue, metrics.Values[0].GetValue())
				}
			}
		}
		response_raw = append(response_raw, NNView{
			Architecture: architecture,
			Decoder:      decoder,
			MetricsName:  metricsName,
			MetricsValue: metricsValue,
		})
	}

	response, err := json.Marshal(response_raw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (k *KatibUIHandler) FetchWorkerTemplates(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	wt, err := k.studyjobClient.GetWorkerTemplates(namespace)
	if err != nil {
		log.Printf("GetWorkerTemplates err %v", err)
	}
	templates := make([]TemplateView, 0)

	for key := range wt {
		templates = append(templates, TemplateView{Name: key, Yaml: wt[key]})
	}
	response, err := json.Marshal(templates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (k *KatibUIHandler) FetchCollectorTemplates(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	wt, err := k.studyjobClient.GetMetricsCollectorTemplates(namespace)
	if err != nil {
		log.Printf("GetWorkerTemplates err %v", err)
	}
	templates := make([]TemplateView, 0)

	for key := range wt {
		fmt.Println(key)
		templates = append(templates, TemplateView{Name: key, Yaml: wt[key]})
	}
	response, err := json.Marshal(templates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (k *KatibUIHandler) AddEditTemplate(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)
	wt, err := k.studyjobClient.GetWorkerTemplates(namespace)

	if data["kind"].(string) == "collector" {
		wt, err = k.studyjobClient.GetMetricsCollectorTemplates(namespace)
	}

	if err != nil {
		log.Printf("fail to GetWorkerTemplates %v", err)
	}
	wt[data["name"].(string)] = data["yaml"].(string)
	err = k.studyjobClient.UpdateWorkerTemplates(wt, namespace)
	if err != nil {
		log.Printf("fail to UpdateWorkerTemplate %v", err)
	}

	templates := make([]TemplateView, 0)

	for key := range wt {
		templates = append(templates, TemplateView{Name: key, Yaml: wt[key]})
	}
	response_raw := TemplateResponse{
		Data:         templates,
		TemplateType: data["kind"].(string),
	}
	response, err := json.Marshal(response_raw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (k *KatibUIHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)
	fmt.Println(data["name"].(string))
	wt, err := k.studyjobClient.GetWorkerTemplates(namespace)

	if data["kind"].(string) == "collector" {
		wt, err = k.studyjobClient.GetMetricsCollectorTemplates(namespace)
	}

	if err != nil {
		log.Printf("fail to GetWorkerTemplates %v", err)
	}
	delete(wt, data["name"].(string))
	err = k.studyjobClient.UpdateWorkerTemplates(wt, namespace)
	if err != nil {
		log.Printf("fail to UpdateWorkerTemplate %v", err)
	}

	templates := make([]TemplateView, 0)

	for key := range wt {
		templates = append(templates, TemplateView{Name: key, Yaml: wt[key]})
	}
	response_raw := TemplateResponse{
		Data:         templates,
		TemplateType: data["kind"].(string),
	}
	response, err := json.Marshal(response_raw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
