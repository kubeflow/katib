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
	"log"
	"net/http"
	"strconv"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb_v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

func (k *KatibUIHandler) FetchNASJobInfo(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	experimentName := r.URL.Query()["experimentName"][0]
	namespace := r.URL.Query()["namespace"][0]

	responseRaw := make([]NNView, 0)
	var architecture string
	var decoder string

	conn, c := k.connectManager()

	defer conn.Close()

	trials, err := k.katibClient.GetTrialList(experimentName, namespace)
	if err != nil {
		log.Printf("GetTrialList from NAS job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Got Trial List")

	for i, t := range trials.Items {
		succeeded := false
		for _, condition := range t.Status.Conditions {
			if condition.Type == trialsv1beta1.TrialSucceeded {
				succeeded = true
			}
		}
		if succeeded {
			obsLogResp, err := c.GetObservationLog(
				context.Background(),
				&api_pb_v1beta1.GetObservationLogRequest{
					TrialName: t.Name,
					StartTime: "",
					EndTime:   "",
				},
			)
			if err != nil {
				log.Printf("GetObservationLog from NAS job failed: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			metricsName := make([]string, 0)
			metricsValue := make([]string, 0)
			for _, m := range obsLogResp.ObservationLog.MetricLogs {
				metricsName = append(metricsName, m.Metric.Name)
				metricsValue = append(metricsValue, m.Metric.Value)

			}
			for _, trialParam := range t.Spec.ParameterAssignments {
				if trialParam.Name == "architecture" {
					architecture = trialParam.Value
				}
				if trialParam.Name == "nn_config" {
					decoder = trialParam.Value
				}
			}
			responseRaw = append(responseRaw, NNView{
				Name:         "Generation " + strconv.Itoa(i),
				TrialName:    t.Name,
				Architecture: generateNNImage(architecture, decoder),
				MetricsName:  metricsName,
				MetricsValue: metricsValue,
			})
		}
	}
	log.Printf("Logs parsed, result: %v", responseRaw)

	response, err := json.Marshal(responseRaw)
	if err != nil {
		log.Printf("Marshal result in NAS job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(response); err != nil {
		log.Printf("Write result in NAS job failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
