package v1beta1

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ghodss/yaml"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	experimentv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	api_pb_v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibclient"
)

func NewKatibUIHandler(dbManagerAddr string) *KatibUIHandler {
	kclient, err := katibclient.NewClient(client.Options{})
	if err != nil {
		log.Printf("NewClient for Katib failed: %v", err)
		panic(err)
	}
	return &KatibUIHandler{
		katibClient:   kclient,
		dbManagerAddr: dbManagerAddr,
	}
}

// ServeIndex will return index.html for any non-API URL
func (k *KatibUIHandler) ServeIndex(buildDir string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fp := filepath.Join(buildDir, "static/index.html")
		log.Printf("Sending file %s for url: %s", fp, r.URL)

		// never cache index.html
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")

		// return the contents of index.html
		http.ServeFile(w, r, fp)
	}
}

func (k *KatibUIHandler) connectManager() (*grpc.ClientConn, api_pb_v1beta1.DBManagerClient) {
	conn, err := grpc.Dial(k.dbManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Dial to GRPC failed: %v", err)
		return nil, nil
	}
	c := api_pb_v1beta1.NewDBManagerClient(conn)
	return conn, c
}

func (k *KatibUIHandler) CreateExperiment(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	var data map[string]interface{}

	json.NewDecoder(r.Body).Decode(&data)
	dataJSON, ok := data["postData"]
	if !ok {
		msg := "Couldn't load the 'postData' field of the request's data"
		log.Printf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	jsonbody, err := json.Marshal(dataJSON)
	if err != nil {
		log.Printf("Marshal data for experiment failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	job := experimentv1beta1.Experiment{}
	if err := json.Unmarshal(jsonbody, &job); err != nil {
		log.Printf("Unmarshal experiment failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dataMap := dataJSON.(map[string]interface{})
	job.TypeMeta = metav1.TypeMeta{
		APIVersion: "kubeflow.org/v1beta1",
		Kind:       "Experiment",
	}
	job.ObjectMeta = metav1.ObjectMeta{
		Name:      dataMap["metadata"].(map[string]interface{})["name"].(string),
		Namespace: dataMap["metadata"].(map[string]interface{})["namespace"].(string),
	}
	err = k.katibClient.CreateRuntimeObject(&job)
	if err != nil {
		log.Printf("CreateRuntimeObject from parameters failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FetchAllExperiments gets HP and NAS experiments in all namespaces.
func (k *KatibUIHandler) FetchAllExperiments(w http.ResponseWriter, r *http.Request) {
	// At first, try to list experiments in cluster scope
	experiments, err := k.getExperiments([]string{""})
	if err != nil {
		// If failed, just try to list experiments from own namespace
		experiments, err = k.getExperiments([]string{})
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(experiments)
	if err != nil {
		log.Printf("Marshal experiments failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
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
	err = k.katibClient.DeleteRuntimeObject(experiment)
	if err != nil {
		log.Printf("DeleteRuntimeObject failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isExperimentDeleted := false

	var experiments []ExperimentView

	// Waiting until experiment will be deleted
	for !isExperimentDeleted {
		// At first, try to list experiments in cluster scope
		experiments, err = k.getExperiments([]string{""})
		if err != nil {
			// If failed, just try to list experiments from own namespace
			experiments, err = k.getExperiments([]string{})
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isExperimentDeleted = true
		for _, experiment := range experiments {
			if experiment.Name == experimentName {
				isExperimentDeleted = false
				break
			}
		}
	}

	response, err := json.Marshal(experiments)
	if err != nil {
		log.Printf("Marshal HP and NAS experiments failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// FetchTrialTemplates gets all trial templates in all namespaces
func (k *KatibUIHandler) FetchTrialTemplates(w http.ResponseWriter, r *http.Request) {

	trialTemplatesViewList, err := k.getTrialTemplatesViewList()
	if err != nil {
		log.Printf("getTrialTemplatesViewList failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	TrialTemplatesResponse := TrialTemplatesResponse{
		Data: trialTemplatesViewList,
	}
	response, err := json.Marshal(TrialTemplatesResponse)
	if err != nil {
		log.Printf("Marshal templates failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

//AddTemplate adds template to ConfigMap
func (k *KatibUIHandler) AddTemplate(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	updatedConfigMapNamespace := data["updatedConfigMapNamespace"].(string)
	updatedConfigMapName := data["updatedConfigMapName"].(string)
	updatedConfigMapPath := data["updatedConfigMapPath"].(string)
	updatedTemplateYaml := data["updatedTemplateYaml"].(string)

	newTemplates, err := k.updateTrialTemplates(updatedConfigMapNamespace, updatedConfigMapName, "", updatedConfigMapPath, updatedTemplateYaml, ActionTypeAdd)
	if err != nil {
		log.Printf("updateTrialTemplates failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	TrialTemplatesResponse := TrialTemplatesResponse{
		Data: newTemplates,
	}
	response, err := json.Marshal(TrialTemplatesResponse)
	if err != nil {
		log.Printf("Marhal failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)

}

// EditTemplate edits template in ConfigMap
func (k *KatibUIHandler) EditTemplate(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	updatedConfigMapNamespace := data["updatedConfigMapNamespace"].(string)
	updatedConfigMapName := data["updatedConfigMapName"].(string)
	configMapPath := data["configMapPath"].(string)
	updatedConfigMapPath := data["updatedConfigMapPath"].(string)
	updatedTemplateYaml := data["updatedTemplateYaml"].(string)

	newTemplates, err := k.updateTrialTemplates(updatedConfigMapNamespace, updatedConfigMapName, configMapPath, updatedConfigMapPath, updatedTemplateYaml, ActionTypeEdit)
	if err != nil {
		log.Printf("updateTrialTemplates failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	TrialTemplatesResponse := TrialTemplatesResponse{
		Data: newTemplates,
	}
	response, err := json.Marshal(TrialTemplatesResponse)
	if err != nil {
		log.Printf("Marhal failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// DeleteTemplate deletes template in ConfigMap
func (k *KatibUIHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	updatedConfigMapNamespace := data["updatedConfigMapNamespace"].(string)
	updatedConfigMapName := data["updatedConfigMapName"].(string)
	updatedConfigMapPath := data["updatedConfigMapPath"].(string)

	newTemplates, err := k.updateTrialTemplates(updatedConfigMapNamespace, updatedConfigMapName, "", updatedConfigMapPath, "", ActionTypeDelete)
	if err != nil {
		log.Printf("updateTrialTemplates failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	TrialTemplatesResponse := TrialTemplatesResponse{
		Data: newTemplates,
	}

	response, err := json.Marshal(TrialTemplatesResponse)
	if err != nil {
		log.Printf("Marhal failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (k *KatibUIHandler) FetchNamespaces(w http.ResponseWriter, r *http.Request) {

	// Get all available namespaces
	namespaces, err := k.getAvailableNamespaces()
	if err != nil {
		log.Printf("GetAvailableNamespaces failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(namespaces)
	if err != nil {
		log.Printf("Marshal namespaces failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// FetchExperiment gets experiment in specific namespace.
func (k *KatibUIHandler) FetchExperiment(w http.ResponseWriter, r *http.Request) {
	experimentName := r.URL.Query()["experimentName"][0]
	namespace := r.URL.Query()["namespace"][0]

	experiment, err := k.katibClient.GetExperiment(experimentName, namespace)
	if err != nil {
		log.Printf("GetExperiment failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(experiment)
	if err != nil {
		log.Printf("Marshal Experiment failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// FetchSuggestion gets suggestion in specific namespace
func (k *KatibUIHandler) FetchSuggestion(w http.ResponseWriter, r *http.Request) {
	suggestionName := r.URL.Query()["suggestionName"][0]
	namespace := r.URL.Query()["namespace"][0]

	suggestion, err := k.katibClient.GetSuggestion(suggestionName, namespace)
	if err != nil {
		log.Printf("GetSuggestion failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(suggestion)
	if err != nil {
		log.Printf("Marshal Suggestion failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
