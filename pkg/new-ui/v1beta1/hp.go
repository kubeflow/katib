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
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	trialv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb_v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

const kfpRunIDAnnotation = "kubeflow-kale.org/kfp-run-uuid"

func (k *KatibUIHandler) FetchHPJobInfo(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)

	namespaces, ok := r.URL.Query()["namespace"]
	if !ok {
		log.Printf("No namespace provided in Query parameters! Provide a 'namespace' param")
		err := errors.New("no 'namespace' provided")
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
		log.Printf("The user: %s is not authorized to get experiments from namespace: %s \n", user, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	log.Printf("Start FetchHPJobInfo for Experiment: %v in namespace: %v", experimentName, namespace)

	conn, c := k.connectManager()
	defer conn.Close()

	resultText := "Status,trialName"
	experiment, err := k.katibClient.GetExperiment(experimentName, namespace)
	if err != nil {
		log.Printf("GetExperiment from HP job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Got Experiment %v", experimentName)
	metricsList := map[string]int{}
	metricsName := experiment.Spec.Objective.ObjectiveMetricName
	resultText += "," + metricsName
	metricsList[metricsName] = 0
	for i, m := range experiment.Spec.Objective.AdditionalMetricNames {
		resultText += "," + m
		metricsList[m] = i + 1
	}
	log.Printf("Got metrics names")
	paramList := map[string]int{}
	for i, p := range experiment.Spec.Parameters {
		resultText += "," + p.Name
		paramList[p.Name] = i + len(metricsList)
	}
	log.Printf("Got Parameters names")

	_, err = IsAuthorized(consts.ActionTypeList, namespace, consts.PluralTrial, "", "", trialv1beta1.SchemeGroupVersion, k.katibClient.GetClient(), r)
	if user == "" && err != nil {
		log.Printf("No user provided in kubeflow-userid header.")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("The user: %s is not authorized to list trials from namespace: %s \n", user, namespace)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	trialList, err := k.katibClient.GetTrialList(experimentName, namespace)
	if err != nil {
		log.Printf("GetTrialList from HP job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Got Trial List - Count: %v", len(trialList.Items))

	// append a column for the Pipeline UID associated with the Trial
	if havePipelineUID(trialList.Items) {
		resultText += ",KFP Run"
	}

	foundPipelineUID := false
	for _, t := range trialList.Items {
		runUid, ok := t.GetAnnotations()[kfpRunIDAnnotation]
		if !ok {
			log.Printf("Trial %s has no pipeline run.", t.Name)
			runUid = ""
		} else {
			foundPipelineUID = true
		}

		var lastTrialCondition string

		// Take only the latest condition
		if len(t.Status.Conditions) > 0 {
			lastTrialCondition = string(t.Status.Conditions[len(t.Status.Conditions)-1].Type)
		}

		trialResText := make([]string, len(metricsList)+len(paramList))

		if t.IsSucceeded() || t.IsEarlyStopped() {
			obsLogResp, err := c.GetObservationLog(
				context.Background(),
				&api_pb_v1beta1.GetObservationLogRequest{
					TrialName: t.Name,
					StartTime: "",
					EndTime:   "",
				},
			)
			if err != nil {
				log.Printf("GetObservationLog from HP job failed: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, m := range obsLogResp.ObservationLog.MetricLogs {
				if trialResText[metricsList[m.Metric.Name]] == "" {
					trialResText[metricsList[m.Metric.Name]] = m.Metric.Value
				} else {
					currentValue, _ := strconv.ParseFloat(m.Metric.Value, 64)
					bestValue, _ := strconv.ParseFloat(trialResText[metricsList[m.Metric.Name]], 64)
					if t.Spec.Objective.Type == commonv1beta1.ObjectiveTypeMinimize && currentValue < bestValue {
						trialResText[metricsList[m.Metric.Name]] = m.Metric.Value
					} else if t.Spec.Objective.Type == commonv1beta1.ObjectiveTypeMaximize && currentValue > bestValue {
						trialResText[metricsList[m.Metric.Name]] = m.Metric.Value
					}
				}
			}
		}
		for _, trialParam := range t.Spec.ParameterAssignments {
			trialResText[paramList[trialParam.Name]] = trialParam.Value
		}
		resultText += "\n" + lastTrialCondition + "," + t.Name + "," + strings.Join(trialResText, ",")
		if foundPipelineUID {
			resultText += "," + runUid
		}
	}
	log.Printf("Logs parsed, results:\n %v", resultText)
	response, err := json.Marshal(resultText)
	if err != nil {
		log.Printf("Marshal result text for HP job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(response); err != nil {
		log.Printf("Write result text for HP job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FetchHPJobTrialInfo returns all metrics for the HP Job Trial
func (k *KatibUIHandler) FetchHPJobTrialInfo(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)

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

	conn, c := k.connectManager()
	defer conn.Close()

	user, err := IsAuthorized(consts.ActionTypeList, namespace, consts.PluralTrial, "", trialName, trialv1beta1.SchemeGroupVersion, k.katibClient.GetClient(), r)
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
		log.Printf("GetTrial from HP job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	objectiveType := trial.Spec.Objective.Type

	// resultArray - array of arrays, where [i][0] - metricName, [i][1] - metricTime, [i][2] - metricValue
	var resultArray [][]string
	resultArray = append(resultArray, strings.Split("metricName,time,value", ","))
	obsLogResp, err := c.GetObservationLog(
		context.Background(),
		&api_pb_v1beta1.GetObservationLogRequest{
			TrialName: trialName,
			StartTime: "",
			EndTime:   "",
		},
	)
	if err != nil {
		log.Printf("GetObservationLog failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// prevMetricTimeValue is the dict, where key = metric name,
	// value = array, where [0] - Last metric time, [1] - Best metric value for this time
	prevMetricTimeValue := make(map[string][]string)
	for _, m := range obsLogResp.ObservationLog.MetricLogs {
		parsedCurrentTime, _ := time.Parse(time.RFC3339Nano, m.TimeStamp)
		formatCurrentTime := parsedCurrentTime.Format("2006-01-02T15:04:05")
		if _, found := prevMetricTimeValue[m.Metric.Name]; !found {
			prevMetricTimeValue[m.Metric.Name] = []string{"", ""}

		}

		newMetricValue, err := strconv.ParseFloat(m.Metric.Value, 64)
		if err != nil {
			log.Printf("ParseFloat for new metric value: %v failed: %v", m.Metric.Value, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var prevMetricValue float64
		if prevMetricTimeValue[m.Metric.Name][1] != "" {
			prevMetricValue, err = strconv.ParseFloat(prevMetricTimeValue[m.Metric.Name][1], 64)
			if err != nil {
				log.Printf("ParseFloat for prev metric value: %v failed: %v", prevMetricTimeValue[m.Metric.Name][1], err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if formatCurrentTime == prevMetricTimeValue[m.Metric.Name][0] &&
			((objectiveType == commonv1beta1.ObjectiveTypeMinimize &&
				newMetricValue < prevMetricValue) ||
				(objectiveType == commonv1beta1.ObjectiveTypeMaximize &&
					newMetricValue > prevMetricValue)) {

			prevMetricTimeValue[m.Metric.Name][1] = m.Metric.Value
			for i := len(resultArray) - 1; i >= 0; i-- {
				if resultArray[i][0] == m.Metric.Name {
					resultArray[i][2] = m.Metric.Value
					break
				}
			}
		} else if formatCurrentTime != prevMetricTimeValue[m.Metric.Name][0] {
			resultArray = append(resultArray, []string{m.Metric.Name, formatCurrentTime, m.Metric.Value})
			prevMetricTimeValue[m.Metric.Name][0] = formatCurrentTime
			prevMetricTimeValue[m.Metric.Name][1] = m.Metric.Value
		}
	}

	var resultText string
	for _, metric := range resultArray {
		resultText += strings.Join(metric, ",") + "\n"
	}

	response, err := json.Marshal(resultText)
	if err != nil {
		log.Printf("Marshal result text in Trial info failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(response); err != nil {
		log.Printf("Write result text in Trial info failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
