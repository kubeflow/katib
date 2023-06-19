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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sigs.k8s.io/controller-runtime/pkg/client"

	experimentv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb_v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	consts "github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibclient"
	corev1 "k8s.io/api/core/v1"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	mccommon "github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
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
	experimentName := job.ObjectMeta.Name

	user, err := IsAuthorized(consts.ActionTypeCreate, namespace, consts.PluralExperiment, "", experimentName, experimentv1beta1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to create experiment: %s in namespace: %s \n", user, experimentName, namespace)
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

func (k *KatibUIHandler) FetchExperiments(w http.ResponseWriter, r *http.Request) {

	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No 'namespace' query parameter was provided.")
		err := errors.New("no 'namespace' query parameter was provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	namespace := namespaces[0]

	user, err := IsAuthorized(consts.ActionTypeList, namespace, consts.PluralExperiment, "", "", experimentv1beta1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to list experiments in namespace: %s \n", user, namespace)
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

func (k *KatibUIHandler) DeleteExperiment(w http.ResponseWriter, r *http.Request) {

	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No 'namespace' query parameter was provided.")
		err := errors.New("no 'namespace' query parameter was provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	experimentNames, ok := r.URL.Query()["experimentName"]
	if !ok {
		log.Printf("No experimentName provided in Query parameteres! Provide an 'experimentName' param")
		err := errors.New("no experimentName provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	experimentName := experimentNames[0]
	namespace := namespaces[0]

	user, err := IsAuthorized(consts.ActionTypeDelete, namespace, consts.PluralExperiment, "", experimentName, experimentv1beta1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to delete experiment: %s in namespace: %s \n", user, experimentName, namespace)
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

	trialTemplatesViewList, err := k.getTrialTemplatesViewList(r)
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

// AddTemplate adds template to ConfigMap
func (k *KatibUIHandler) AddTemplate(w http.ResponseWriter, r *http.Request) {

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

	user, err := IsAuthorized(consts.ActionTypeCreate, updatedConfigMapNamespace, corev1.ResourceConfigMaps.String(), "", updatedConfigMapName, corev1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to add configmap: %s in namespace: %s \n", user, updatedConfigMapName, updatedConfigMapNamespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	newTemplates, err := k.updateTrialTemplates(updatedConfigMapNamespace, updatedConfigMapName, "", updatedConfigMapPath, updatedTemplateYaml, ActionTypeAdd, r)
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

	user, err := IsAuthorized(consts.ActionTypeUpdate, updatedConfigMapNamespace, corev1.ResourceConfigMaps.String(), "", updatedConfigMapName, corev1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to edit configmap: %s in namespace: %s \n", user, updatedConfigMapName, updatedConfigMapNamespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	newTemplates, err := k.updateTrialTemplates(updatedConfigMapNamespace, updatedConfigMapName, configMapPath, updatedConfigMapPath, updatedTemplateYaml, ActionTypeEdit, r)
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

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedConfigMapNamespace := data["updatedConfigMapNamespace"].(string)
	updatedConfigMapName := data["updatedConfigMapName"].(string)
	updatedConfigMapPath := data["updatedConfigMapPath"].(string)

	user, err := IsAuthorized(consts.ActionTypeDelete, updatedConfigMapNamespace, corev1.ResourceConfigMaps.String(), "", updatedConfigMapName, corev1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to delete configmap: %s in namespace: %s \n", user, updatedConfigMapName, updatedConfigMapNamespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	newTemplates, err := k.updateTrialTemplates(updatedConfigMapNamespace, updatedConfigMapName, "", updatedConfigMapPath, "", ActionTypeDelete, r)
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

	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No 'namespace' query parameter was provided.")
		err := errors.New("no 'namespace' query parameter was provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	experimentNames, ok := r.URL.Query()["experimentName"]
	if !ok {
		log.Printf("No experimentName provided in Query parameteres! Provide an 'experimentName' param")
		err := errors.New("no experimentName provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	experimentName := experimentNames[0]
	namespace := namespaces[0]

	user, err := IsAuthorized(consts.ActionTypeGet, namespace, consts.PluralExperiment, "", experimentName, experimentv1beta1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to get experiment: %s in namespace: %s \n", user, experimentName, namespace)
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

	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No namespace provided in Query parameters! Provide a 'namespace' param")
		err := errors.New("no 'namespace' provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	suggestionNames, ok := r.URL.Query()["suggestionName"]
	if !ok {
		log.Printf("No experimentName provided in Query parameteres! Provide an 'experimentName' param")
		err := errors.New("no experimentName provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	suggestionName := suggestionNames[0]
	namespace := namespaces[0]

	user, err := IsAuthorized(consts.ActionTypeGet, namespace, consts.PluralSuggestion, "", suggestionName, suggestionv1beta1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to get suggestion: %s in namespace: %s \n", user, suggestionName, namespace)
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
	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No namespace provided in Query parameters! Provide a 'namespace' param")
		err := errors.New("no 'namespace' provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trialNames, ok := r.URL.Query()["trialName"]
	if !ok {
		log.Printf("No trialName provided in Query parameters! Provide a 'trialName' param")
		err := errors.New("no 'trialName' provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trialName := trialNames[0]
	namespace := namespaces[0]

	user, err := IsAuthorized(consts.ActionTypeGet, namespace, consts.PluralTrial, "", trialName, trialsv1beta1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to get trial: %s from namespace: %s \n", user, trialName, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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

// FetchTrialLogs fetches logs for a trial in specific namespace.
func (k *KatibUIHandler) FetchTrialLogs(w http.ResponseWriter, r *http.Request) {
	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No namespace provided in Query parameters! Provide a 'namespace' param")
		err := errors.New("no 'namespace' provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trialNames, ok := r.URL.Query()["trialName"]
	if !ok {
		log.Printf("No trialName provided in Query parameters! Provide a 'trialName' param")
		err := errors.New("no 'trialName' provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trialName := trialNames[0]
	namespace := namespaces[0]

	user, err := IsAuthorized(consts.ActionTypeGet, namespace, consts.PluralTrial, "", trialName, trialsv1beta1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to get trial: %s in namespace: %s \n", user, trialName, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	trial := &trialsv1beta1.Trial{}
	if err := k.katibClient.GetClient().Get(context.Background(), types.NamespacedName{Name: trialName, Namespace: namespace}, trial); err != nil {
		log.Printf("GetLogs failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Use controller-runtime client instead of kubernetes client to get logs, once this is available
	clientset, err := createKubernetesClientset()
	if err != nil {
		log.Printf("GetLogs failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	podName, err := fetchMasterPodName(clientset, trial)
	if err != nil {
		log.Printf("GetLogs failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err = IsAuthorized(consts.ActionTypeGet, namespace, corev1.ResourcePods.String(), "log", podName, corev1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to get pod logs: %s in namespace: %s \n", user, podName, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	podLogOpts := corev1.PodLogOptions{}
	podLogOpts.Container = trial.Spec.PrimaryContainerName
	if trial.Spec.MetricsCollector.Collector.Kind == common.StdOutCollector {
		podLogOpts.Container = mccommon.MetricLoggerCollectorContainerName
	}

	logs, err := fetchPodLogs(clientset, namespace, podName, podLogOpts)
	if err != nil {
		log.Printf("GetLogs failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(logs)
	if err != nil {
		log.Printf("Marshal logs failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(response); err != nil {
		log.Printf("Write logs failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// createKubernetesClientset returns kubernetes clientset
func createKubernetesClientset() (*kubernetes.Clientset, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// fetchMasterPodName returns name of the master pod for a trial
func fetchMasterPodName(clientset *kubernetes.Clientset, trial *trialsv1beta1.Trial) (string, error) {
	selectionLabel := consts.LabelTrialName + "=" + trial.ObjectMeta.Name
	for primaryKey, primaryValue := range trial.Spec.PrimaryPodLabels {
		selectionLabel = selectionLabel + "," + primaryKey + "=" + primaryValue
	}

	podList, err := clientset.CoreV1().Pods(trial.ObjectMeta.Namespace).List(context.Background(), metav1.ListOptions{LabelSelector: selectionLabel})
	if err != nil {
		return "", err
	}

	if len(podList.Items) == 0 {
		return "", errors.New(`Failed to find logs for this Trial. Make sure you've set "spec.trialTemplate.retain"
		field to "true" in the Experiment definition. If this error persists then the Pod's logs are not currently
		persisted in the cluster.`)
	}

	// If Pod is Running or Succeeded Pod, return it.
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodRunning {
			return pod.Name, nil
		}
	}

	// Otherwise, return the first Failed Pod.
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodFailed {
			return pod.Name, nil
		}
	}

	// Otherwise, return error since Pod is in the Pending state.
	return "", errors.New("Failed to get logs for this Trial. Pod is in the Pending or Unknown state.")
}

// fetchPodLogs returns logs of a pod for the given job name and namespace
func fetchPodLogs(clientset *kubernetes.Clientset, namespace string, podName string, podLogOpts corev1.PodLogOptions) (string, error) {
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	str := buf.String()

	return str, nil
}
