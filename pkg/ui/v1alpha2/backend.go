package ui

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kubeflow/katib/pkg"
	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	api_pb_v1alpha2 "github.com/kubeflow/katib/pkg/api/v1alpha2"

	"github.com/kubeflow/katib/pkg/util/v1alpha2/katibclient"
	"google.golang.org/grpc"
	restclient "k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gographviz "github.com/awalterschulze/gographviz"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	namespace      = "default"
	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
)

type Decoder struct {
	Layers     int            `json:"num_layers"`
	InputSize  []int          `json:"input_size"`
	OutputSize []int          `json:"output_size"`
	Embedding  map[int]*Block `json:"embedding"`
}

type Block struct {
	Id    int    `json:"opt_id"`
	Type  string `json:"opt_type"`
	Param Option `json:"opt_params"`
}

type Option struct {
	FilterNumber string `json:"num_filter"`
	FilterSize   string `json:"filter_size"`
	Stride       string `json:"stride"`
}

func get_node_string(block *Block) string {
	var node_string string
	switch block.Type {
	case "convolution":
		node_string += block.Param.FilterSize + "x" + block.Param.FilterSize
		node_string += " conv\n"
		node_string += block.Param.FilterSize + " channels"
	case "separable_convolution":
		node_string += block.Param.FilterSize + "x" + block.Param.FilterSize
		node_string += " sep_conv\n"
		node_string += block.Param.FilterSize + " channels"
	case "depthwise_convolution":
		node_string += block.Param.FilterSize + "x" + block.Param.FilterSize
		node_string += " depth_conv\n"
	case "reduction":
		// fix this
		node_string += "3x3 max_pooling"
	}
	return strconv.Quote(node_string)
}

func generate_nn_image(architecture string, decoder string) string {

	var architecture_int [][]int

	if err := json.Unmarshal([]byte(architecture), &architecture_int); err != nil {
		panic(err)
	}
	/*
		Always has num_layers, input_size, output_size and embeding
		Embeding is a map: int to Parameter
		Parameter has id, type, Option

		Beforehand substite all ' to " and wrap the string in `
	*/

	replaced_decoder := strings.Replace(decoder, `'`, `"`, -1)
	var decoder_parsed Decoder

	err := json.Unmarshal([]byte(replaced_decoder), &decoder_parsed)
	if err != nil {
		panic(err)
	}

	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	graph.AddNode("G", "0", map[string]string{"label": strconv.Quote("Input")})
	var i int
	for i = 0; i < len(architecture_int); i++ {
		graph.AddNode("G", strconv.Itoa(i+1), map[string]string{"label": get_node_string(decoder_parsed.Embedding[architecture_int[i][0]])})
		graph.AddEdge(strconv.Itoa(i), strconv.Itoa(i+1), true, nil)
		for j := 1; j < i+1; j++ {
			if architecture_int[i][j] == 1 {
				graph.AddEdge(strconv.Itoa(j-1), strconv.Itoa(i+1), true, nil)
			}
		}
	}
	graph.AddNode("G", strconv.Itoa(i+1), map[string]string{"label": strconv.Quote("GlobalAvgPool")})
	graph.AddEdge(strconv.Itoa(i), strconv.Itoa(i+1), true, nil)
	graph.AddNode("G", strconv.Itoa(i+2), map[string]string{"label": strconv.Quote("FullConnect\nSoftmax")})
	graph.AddEdge(strconv.Itoa(i+1), strconv.Itoa(i+2), true, nil)
	graph.AddNode("G", strconv.Itoa(i+3), map[string]string{"label": strconv.Quote("Output")})
	graph.AddEdge(strconv.Itoa(i+2), strconv.Itoa(i+3), true, nil)
	s := graph.String()
	return s
}

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

	// For local testing
	// var kubeconfig *string
	// if home := homeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	// flag.Parse()
	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// if err != nil {
	// 	log.Fatalf("getClusterConfig: %v", err)
	// }
	// return config

	// For production
	config, err := restclient.InClusterConfig()
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
	Name   string
	Status string
}

type TemplateView struct {
	Name string
	Yaml string
}

type KatibUIHandler struct {
	katibClient *katibclient.KatibClient
}

type TemplateResponse struct {
	TemplateType string
	Data         []TemplateView
}

func NewKatibUIHandler() *KatibUIHandler {
	kclient, err := katibclient.NewClient(client.Options{})
	if err != nil {
		log.Printf("New Katib client failed: %v", err)
		panic(err)
	}
	return &KatibUIHandler{
		katibClient: kclient,
	}
}

func (k *KatibUIHandler) connectManager() (*grpc.ClientConn, api_pb_v1alpha2.ManagerClient) {
	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Dial to GRPC failed: %v", err)
		return nil, nil
	}
	c := api_pb_v1alpha2.NewManagerClient(conn)
	return conn, c
}

func (k *KatibUIHandler) FetchHPJobs(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	jobs := make([]JobView, 0)

	el, err := k.katibClient.GetExperimentList()
	if err != nil {
		log.Printf("Get Experiment List for HP failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, experiment := range el.Items {
		if experiment.Spec.Parameters != nil {
			jobs = append(jobs, JobView{
				Name: experiment.GetName(),
				// TODO: Delete from frontend
				// ID:     experiment.Status.StudyID,
				// TODO: Parse it in frontend
				Status: string(getExperimentCurrentCondition(&experiment)),
			})
		}
	}

	response, err := json.Marshal(jobs)
	if err != nil {
		log.Printf("Marshal HP jobs failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)

}

func (k *KatibUIHandler) FetchNASJobs(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	jobs := make([]JobView, 0)

	el, err := k.katibClient.GetExperimentList()
	if err != nil {
		log.Printf("Get Experiment List for NAS failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, experiment := range el.Items {
		if experiment.Spec.NasConfig != nil {
			jobs = append(jobs, JobView{
				Name: experiment.GetName(),
				// TODO: Delete from frontend
				// ID:     experiment.Status.StudyID,
				// TODO: Parse it in frontend
				Status: string(getExperimentCurrentCondition(&experiment)),
			})
		}
	}

	response, err := json.Marshal(jobs)
	if err != nil {
		log.Printf("Marshal NAS jobs failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)

}

func (k *KatibUIHandler) SubmitYamlJob(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)

	job := experimentv1alpha2.Experiment{}
	if yamlContent, ok := data["yaml"].(string); ok {
		err := yaml.Unmarshal([]byte(yamlContent), &job)
		if err != nil {
			log.Printf("Unmarshal YAML content failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = k.katibClient.CreateExperiment(&job)
		if err != nil {
			log.Printf("Create Experiment from YAML failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func (k *KatibUIHandler) SubmitHPJob(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)
	if data, ok := data["postData"]; ok {
		jsonbody, err := json.Marshal(data)
		if err != nil {
			log.Printf("Marshal data for HP job failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		job := experimentv1alpha2.Experiment{}
		//TODO: Add new fields of experiment to frontend
		if err := json.Unmarshal(jsonbody, &job); err != nil {
			log.Printf("Unmarshal HP job failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dataMap := data.(map[string]interface{})
		job.TypeMeta = metav1.TypeMeta{
			APIVersion: "kubeflow.org/v1alpha2",
			Kind:       "Experiment",
		}
		job.ObjectMeta = metav1.ObjectMeta{
			Name:      dataMap["metadata"].(map[string]interface{})["name"].(string),
			Namespace: dataMap["metadata"].(map[string]interface{})["namespace"].(string),
		}
		err = k.katibClient.CreateExperiment(&job)
		if err != nil {
			log.Printf("Create Experiment for HP failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (k *KatibUIHandler) SubmitNASJob(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)
	if data, ok := data["postData"]; ok {
		jsonbody, err := json.Marshal(data)
		if err != nil {
			log.Printf("Marshal data for NAS job failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		job := experimentv1alpha2.Experiment{}
		if err := json.Unmarshal(jsonbody, &job); err != nil {
			log.Printf("Unmarshal NAS job failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// mapstructure.Decode(data, &job)
		// // think of a better way
		dataMap := data.(map[string]interface{})

		job.TypeMeta = metav1.TypeMeta{
			APIVersion: "kubeflow.org/v1alpha2",
			Kind:       "Experiment",
		}
		job.ObjectMeta = metav1.ObjectMeta{
			Name:      dataMap["metadata"].(map[string]interface{})["name"].(string),
			Namespace: dataMap["metadata"].(map[string]interface{})["namespace"].(string),
		}

		err = k.katibClient.CreateExperiment(&job)
		if err != nil {
			log.Printf("Create Experiment for NAS failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

//TODO:
// 1. Add delete job to Katib Client
// 2. Change id to name in frontend
func (k *KatibUIHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	experimentName := r.URL.Query()["name"][0]
	log.Printf("Experiment Name: %v", experimentName)

}

func (k *KatibUIHandler) FetchHPJobInfo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	//TODO: Change id to name in frontend
	experimentName := r.URL.Query()["name"][0]

	conn, c := k.connectManager()
	defer conn.Close()
	retText := "TrialName"
	expResp, err := c.GetExperiment(
		context.Background(),
		&api_pb_v1alpha2.GetExperimentRequest{
			ExperimentName: experimentName,
		},
	)
	if err != nil {
		log.Printf("Get Experiment from HP job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Got Experiment")
	metricsList := map[string]int{}
	metricsName := expResp.Experiment.ExperimentSpec.Objective.ObjectiveMetricName
	retText += "," + metricsName
	metricsList[metricsName] = 0
	for i, m := range expResp.Experiment.ExperimentSpec.Objective.AdditionalMetricsNames {
		retText += "," + m
		metricsList[m] = i + 1
	}
	log.Printf("Got metrics names")
	paramList := map[string]int{}
	for i, p := range expResp.Experiment.ExperimentSpec.ParameterSpecs.Parameters {
		retText += "," + p.Name
		paramList[p.Name] = i + len(metricsList)
	}
	log.Printf("Got Parameters names")
	trialListResp, err := c.GetTrialList(
		context.Background(),
		&api_pb_v1alpha2.GetTrialListRequest{
			ExperimentName: expResp.Experiment.GetName(),
			Filter:         "",
		},
	)
	if err != nil {
		log.Printf("Get Trial List from HP job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Got Trial List")
	obsLogResp, err := c.GetObservationLog(
		context.Background(),
		&api_pb_v1alpha2.GetObservationLogRequest{
			TrialName:       trialListResp.,
			OnlyLatestLog: true,
		},
	)
	log.Printf("Got Parameters results")
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
	log.Printf("Parsed logs")
	log.Printf("%v", retText)
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
		&api_pb_v1alpha2.GetWorkerFullInfoRequest{
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
	studyID := r.URL.Query()["id"][0]

	conn, c := k.connectManager()

	defer conn.Close()

	gtrep, err := c.GetTrials(
		context.Background(),
		&api_pb_v1alpha2.GetTrialsRequest{
			StudyId: studyID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	type NNView struct {
		Name         string
		Architecture string
		MetricsName  []string
		MetricsValue []string
	}

	response_raw := make([]NNView, 0)

	var architecture string
	var decoder string

	gwfirep, err := c.GetWorkerFullInfo(
		context.Background(),
		&api_pb_v1alpha2.GetWorkerFullInfoRequest{
			StudyId:       studyID,
			OnlyLatestLog: true,
		},
	)

	for i, tr := range gtrep.Trials {
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
			Name:         "Generation " + strconv.Itoa(i),
			Architecture: generate_nn_image(architecture, decoder),
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
	wt, err := k.studyjobClient.GetWorkerTemplates()
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
	wt, err := k.studyjobClient.GetMetricsCollectorTemplates()
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

func (k *KatibUIHandler) AddEditTemplate(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)
	var wt map[string]string
	var err error
	if data["kind"].(string) == "collector" {
		wt, err = k.studyjobClient.GetMetricsCollectorTemplates()
	} else {
		wt, err = k.studyjobClient.GetWorkerTemplates()
	}

	if err != nil {
		log.Printf("fail to GetWorkerTemplates %v", err)
	}
	wt[data["name"].(string)] = data["yaml"].(string)
	if data["kind"].(string) == "collector" {
		err = k.studyjobClient.UpdateMetricsCollectorTemplates(wt)
	} else {
		err = k.studyjobClient.UpdateWorkerTemplates(wt)
	}
	if err != nil {
		log.Printf("fail to update template %v", err)
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
	var wt map[string]string
	var err error
	if data["kind"].(string) == "collector" {
		wt, err = k.studyjobClient.GetMetricsCollectorTemplates()
	} else {
		wt, err = k.studyjobClient.GetWorkerTemplates()
	}

	if err != nil {
		log.Printf("fail to GetWorkerTemplates %v", err)
	}
	delete(wt, data["name"].(string))
	if data["kind"].(string) == "collector" {
		err = k.studyjobClient.UpdateMetricsCollectorTemplates(wt)
	} else {
		err = k.studyjobClient.UpdateWorkerTemplates(wt)
	}
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
