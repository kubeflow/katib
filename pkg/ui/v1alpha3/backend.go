package v1alpha3

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ghodss/yaml"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	experimentv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	api_pb_v1alpha3 "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	common_v1alpha3 "github.com/kubeflow/katib/pkg/common/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	"github.com/kubeflow/katib/pkg/util/v1alpha3/katibclient"
)

func NewKatibUIHandler() *KatibUIHandler {
	kclient, err := katibclient.NewClient(client.Options{})
	if err != nil {
		log.Printf("NewClient for Katib failed: %v", err)
		panic(err)
	}
	return &KatibUIHandler{
		katibClient: kclient,
	}
}

func (k *KatibUIHandler) connectManager() (*grpc.ClientConn, api_pb_v1alpha3.ManagerClient) {
	conn, err := grpc.Dial(common_v1alpha3.KatibDBManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Dial to GRPC failed: %v", err)
		return nil, nil
	}
	c := api_pb_v1alpha3.NewManagerClient(conn)
	return conn, c
}

func (k *KatibUIHandler) SubmitYamlJob(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)

	job := experimentv1alpha3.Experiment{}
	if yamlContent, ok := data["yaml"].(string); ok {
		err := yaml.Unmarshal([]byte(yamlContent), &job)
		if err != nil {
			log.Printf("Unmarshal YAML content failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = k.katibClient.CreateExperiment(&job)
		if err != nil {
			log.Printf("CreateExperiment from YAML failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func (k *KatibUIHandler) SubmitParamsJob(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)
	if data, ok := data["postData"]; ok {
		jsonbody, err := json.Marshal(data)
		if err != nil {
			log.Printf("Marshal data for HP job failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		job := experimentv1alpha3.Experiment{}
		if err := json.Unmarshal(jsonbody, &job); err != nil {
			log.Printf("Unmarshal HP job failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dataMap := data.(map[string]interface{})
		job.TypeMeta = metav1.TypeMeta{
			APIVersion: "kubeflow.org/v1alpha3",
			Kind:       "Experiment",
		}
		job.ObjectMeta = metav1.ObjectMeta{
			Name:      dataMap["metadata"].(map[string]interface{})["name"].(string),
			Namespace: dataMap["metadata"].(map[string]interface{})["namespace"].(string),
		}
		err = k.katibClient.CreateExperiment(&job)
		if err != nil {
			log.Printf("CreateExperiment for HP failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (k *KatibUIHandler) DeleteExperiment(w http.ResponseWriter, r *http.Request) {
	experimentName := r.URL.Query()["experimentName"][0]
	namespace := r.URL.Query()["namespace"][0]

	experiment, err := k.katibClient.GetExperiment(experimentName, namespace)
	if err != nil {
		log.Printf("GetExperiment failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = k.katibClient.DeleteExperiment(experiment)
	if err != nil {
		log.Printf("DeleteExperiment failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FetchTrialTemplates gets the trial templates for the given namespace.
func (k *KatibUIHandler) FetchTrialTemplates(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	namespace := r.URL.Query()["namespace"][0]
	if namespace == "" {
		namespace = consts.DefaultKatibNamespace
	}
	trialTemplates, err := k.katibClient.GetTrialTemplates(namespace)
	if err != nil {
		log.Printf("GetTrialTemplate failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(getTemplatesView(trialTemplates))
	if err != nil {
		log.Printf("Marshal templates failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (k *KatibUIHandler) AddEditDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	//TODO: need to delete?
	if r.Method == "OPTIONS" {
		return
	}
	var data map[string]interface{}
	var err error
	var templateResponse TemplateResponse

	json.NewDecoder(r.Body).Decode(&data)
	if data["action"].(string) == "delete" {
		templateResponse, err = k.updateTemplates(data, true)
	} else {
		templateResponse, err = k.updateTemplates(data, false)
	}
	if err != nil {
		log.Printf("updateTemplates failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(templateResponse)
	if err != nil {
		log.Printf("Marhal failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (k *KatibUIHandler) FetchNamespaces(w http.ResponseWriter, r *http.Request) {

	namespaceList, err := k.katibClient.GetNamespaceList()
	if err != nil {
		log.Printf("GetNamespaceList failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var namespaces []string

	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, namespace.ObjectMeta.Name)
	}

	response, err := json.Marshal(namespaces)
	if err != nil {
		log.Printf("Marshal namespaces failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
