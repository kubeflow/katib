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

package v1beta1

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	experimentv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	api_pb_v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibclient"
	"github.com/pkg/errors"
)

func NewKatibUIHandler(dbManagerAddr string) *KatibUIHandler {
	kclient, err := katibclient.NewClient(client.Options{})
	if err != nil {
		log.Printf("NewClient for Katib failed: %v", err)
		panic(err)
	}
	// create a new client-go client for sending SAR objects in the API-SERVER
	conf, err := config.GetConfig()
	sarclient, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Printf("SarClient for Katib failes: %v", err)
		panic(err)
	}
	return &KatibUIHandler{
		katibClient:   kclient,
		sarClient:     *sarclient,
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
	conn, err := grpc.Dial(k.dbManagerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Dial to GRPC failed: %v", err)
		return nil, nil
	}
	c := api_pb_v1beta1.NewDBManagerClient(conn)
	return conn, c
}

func (k *KatibUIHandler) CreateExperiment(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	// check if user's username is provided in request Header
	user, err := GetUsername(r)
	if err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dataJSON, ok := data["postData"]
	if !ok {
		msg := "Couldn't load the 'postData' field of the request's data"
		log.Println(msg)
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

	namespace := job.ObjectMeta.Namespace
	expName := job.ObjectMeta.Name

	err = IsAuthorized(user, "create", namespace, "kubeflow.org", "v1beta1", "experiments", "", "", &k.sarClient)
	if err != nil {
		log.Printf("The user: %s is not authorized to create experiment: %s from namespace: %s \n", user, expName, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = k.katibClient.CreateRuntimeObject(&job)
	if err != nil {
		log.Printf("CreateRuntimeObject from parameters failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (k *KatibUIHandler) FetchNamespacedExperiments(w http.ResponseWriter, r *http.Request) {

	// check if user's username is provided in request Header
	user, err := GetUsername(r)
	if err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No 'namespace' query parameter was provided.")
		err := errors.New("No 'namespace' query parameter was provided.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	namespace := namespaces[0]

	err = IsAuthorized(user, "list", namespace, "kubeflow.org", "v1beta1", "experiments", "", "", &k.sarClient)
	if err != nil {
		log.Printf("The user: %s is not authorized to list experiments from namespace: %s \n", user, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	experiments, err := k.getExperiments([]string{namespace})
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

	if _, err = w.Write(response); err != nil {
		log.Printf("Write experiments failed: %v", err)
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
	if _, err = w.Write(response); err != nil {
		log.Printf("Write experiments failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (k *KatibUIHandler) DeleteExperiment(w http.ResponseWriter, r *http.Request) {

	// check if user's username is provided in request Header
	user, err := GetUsername(r)
	if err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No 'namespace' query parameter was provided.")
		err := errors.New("No 'namespace' query parameter was provided.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	experimentNames, ok := r.URL.Query()["experimentName"]
	if !ok {
		log.Printf("No experimentName provided in Query parameteres! Provide an 'experimentName' param")
		err := errors.New("No experimentName provided!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	experimentName := experimentNames[0]
	namespace := namespaces[0]

	err = IsAuthorized(user, "delete", namespace, "kubeflow.org", "v1beta1", "experiments", "", experimentName, &k.sarClient)
	if err != nil {
		log.Printf("The user: %s is not authorized to delete experiment: %s from namespace: %s \n", user, experimentName, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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
	if _, err = w.Write(response); err != nil {
		log.Printf("Write HP and NAS experiments failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FetchTrialTemplates gets all trial templates in all namespaces
func (k *KatibUIHandler) FetchTrialTemplates(w http.ResponseWriter, r *http.Request) {

	// check if user's username is provided in request Header
	user, err := GetUsername(r)
	if err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	trialTemplatesViewList, err := k.getTrialTemplatesViewList(user)
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
	if _, err = w.Write(response); err != nil {
		log.Printf("Write template failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//AddTemplate adds template to ConfigMap
func (k *KatibUIHandler) AddTemplate(w http.ResponseWriter, r *http.Request) {

	// check if user's username is provided in request Header
	user, err := GetUsername(r)
	if err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedConfigMapNamespace := data["updatedConfigMapNamespace"].(string)
	updatedConfigMapName := data["updatedConfigMapName"].(string)
	updatedConfigMapPath := data["updatedConfigMapPath"].(string)
	updatedTemplateYaml := data["updatedTemplateYaml"].(string)

	err = IsAuthorized(user, "add", updatedConfigMapNamespace, "", "v1", "configmaps", "", updatedConfigMapName, &k.sarClient)
	if err != nil {
		log.Printf("The user: %s is not authorized to add configmap: %s from namespace: %s \n", user, updatedConfigMapName, updatedConfigMapNamespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	newTemplates, err := k.updateTrialTemplates(updatedConfigMapNamespace, updatedConfigMapName, "", updatedConfigMapPath, updatedTemplateYaml, ActionTypeAdd, user)
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
	if _, err = w.Write(response); err != nil {
		log.Printf("Write failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// EditTemplate edits template in ConfigMap
func (k *KatibUIHandler) EditTemplate(w http.ResponseWriter, r *http.Request) {

	// check if user's username is provided in request Header
	user, err := GetUsername(r)
	if err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedConfigMapNamespace := data["updatedConfigMapNamespace"].(string)
	updatedConfigMapName := data["updatedConfigMapName"].(string)
	configMapPath := data["configMapPath"].(string)
	updatedConfigMapPath := data["updatedConfigMapPath"].(string)
	updatedTemplateYaml := data["updatedTemplateYaml"].(string)

	err = IsAuthorized(user, "edit", updatedConfigMapNamespace, "", "v1", "configmaps", "", updatedConfigMapName, &k.sarClient)
	if err != nil {
		log.Printf("The user: %s is not authorized to edit configmap: %s from namespace: %s \n", user, updatedConfigMapName, updatedConfigMapNamespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	newTemplates, err := k.updateTrialTemplates(updatedConfigMapNamespace, updatedConfigMapName, configMapPath, updatedConfigMapPath, updatedTemplateYaml, ActionTypeEdit, user)
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
	if _, err = w.Write(response); err != nil {
		log.Printf("Write failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteTemplate deletes template in ConfigMap
func (k *KatibUIHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {

	// check if user's username is provided in request Header
	user, err := GetUsername(r)
	if err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedConfigMapNamespace := data["updatedConfigMapNamespace"].(string)
	updatedConfigMapName := data["updatedConfigMapName"].(string)
	updatedConfigMapPath := data["updatedConfigMapPath"].(string)

	err = IsAuthorized(user, "delete", updatedConfigMapNamespace, "", "v1", "configmaps", "", updatedConfigMapName, &k.sarClient)
	if err != nil {
		log.Printf("The user: %s is not authorized to delete configmap: %s from namespace: %s \n", user, updatedConfigMapName, updatedConfigMapNamespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	newTemplates, err := k.updateTrialTemplates(updatedConfigMapNamespace, updatedConfigMapName, "", updatedConfigMapPath, "", ActionTypeDelete, user)
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
	if _, err = w.Write(response); err != nil {
		log.Printf("Write failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	if _, err = w.Write(response); err != nil {
		log.Printf("Write namespaces failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FetchExperiment gets experiment in specific namespace.
func (k *KatibUIHandler) FetchExperiment(w http.ResponseWriter, r *http.Request) {

	user, err := GetUsername(r)
	if err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No 'namespace' query parameter was provided.")
		err := errors.New("No 'namespace' query parameter was provided.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	experimentNames, ok := r.URL.Query()["experimentName"]
	if !ok {
		log.Printf("No experimentName provided in Query parameteres! Provide an 'experimentName' param")
		err := errors.New("No experimentName provided!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	experimentName := experimentNames[0]
	namespace := namespaces[0]

	err = IsAuthorized(user, "get", namespace, "kubeflow.org", "v1beta1", "experiments", "", experimentName, &k.sarClient)
	if err != nil {
		log.Printf("The user: %s is not authorized to get experiment: %s from namespace: %s \n", user, experimentName, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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
	if _, err = w.Write(response); err != nil {
		log.Printf("Write Experiment failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FetchSuggestion gets suggestion in specific namespace
func (k *KatibUIHandler) FetchSuggestion(w http.ResponseWriter, r *http.Request) {

	user, err := GetUsername(r)
	if err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No namespace provided in Query parameters! Provide a 'namespace' param")
		err := errors.New("No 'namespace' provided!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	suggestionNames, ok := r.URL.Query()["suggestionName"]
	if !ok {
		log.Printf("No experimentName provided in Query parameteres! Provide an 'experimentName' param")
		err := errors.New("No experimentName provided!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	suggestionName := suggestionNames[0]
	namespace := namespaces[0]

	err = IsAuthorized(user, "get", namespace, "kubeflow.org", "v1beta1", "suggestions", "", suggestionName, &k.sarClient)
	if err != nil {
		log.Printf("The user: %s is not authorized to get suggestion: %s from namespace: %s \n", user, suggestionName, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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
	if _, err = w.Write(response); err != nil {
		log.Printf("Write Suggestion failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FetchTrial gets trial in specific namespace.
func (k *KatibUIHandler) FetchTrial(w http.ResponseWriter, r *http.Request) {
	trialName := r.URL.Query()["trialName"][0]
	namespace := r.URL.Query()["namespace"][0]

	trial, err := k.katibClient.GetTrial(trialName, namespace)
	if err != nil {
		log.Printf("GetTrial failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(trial)
	if err != nil {
		log.Printf("Marshal Trial failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(response); err != nil {
		log.Printf("Write trial failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
