/*
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
package studyjob

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"

	batchv1beta "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func validateWorkerResource(wkind string) error {
	for _, crd := range invalidCRDResources {
		if crd == wkind {
			return fmt.Errorf("Cannot support %s; Please install the CRD and restart studyjob-controller", wkind)
		}
	}
	return nil
}

func isFatalWatchError(err error, job string) bool {
	if err == nil {
		return false
	}
	if meta.IsNoMatchError(err) {
		invalidCRDResources = append(invalidCRDResources, job)
		log.Printf("Fail to watch CRD: %v; Please install the CRD and restart studyjob-controller to support %s worker", err, job)
		return false
	} else {
		return true
	}
}

func getWorkerKind(workerSpec *katibv1alpha1.WorkerSpec) (*schema.GroupVersionKind, error) {
	var typeChecker interface{}
	BUFSIZE := 1024
	_, m, err := getWorkerManifest(
		nil,
		"validation",
		&katibapi.Trial{
			TrialId:      "validation",
			ParameterSet: []*katibapi.Parameter{},
		},
		workerSpec,
		"",
		"",
		true,
	)
	if err != nil {
		return nil, err
	}
	if err := k8syaml.NewYAMLOrJSONDecoder(m, BUFSIZE).Decode(&typeChecker); err != nil {
		log.Printf("Yaml decode validation error %v", err)
		return nil, err
	}
	tcMap, ok := typeChecker.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Cannot get kind of worker %v", typeChecker)
	}
	wkind, ok := tcMap["kind"]
	if !ok {
		return nil, fmt.Errorf("Cannot get kind of worker %v", typeChecker)
	}
	wkindS, ok := wkind.(string)
	if !ok {
		return nil, fmt.Errorf("Cannot get kind of worker %v", typeChecker)
	}
	apiVersion, ok := tcMap["apiVersion"]
	if !ok {
		return nil, fmt.Errorf("Cannot get apiVersion of worker %v", typeChecker)
	}
	apiVersionS, ok := apiVersion.(string)
	if !ok {
		return nil, fmt.Errorf("Cannot get apiVersion of worker %v", typeChecker)
	}
	for _, kind := range ValidWorkerKindList {
		if kind == wkindS {
			workerGVK := schema.FromAPIVersionAndKind(apiVersionS, kind)
			return &workerGVK, validateWorkerResource(kind)
		}
	}
	return nil, fmt.Errorf("Invalid kind of worker %v", typeChecker)
}

func validateStudy(instance *katibv1alpha1.StudyJob) error {
	if instance.Spec.SuggestionSpec == nil {
		return fmt.Errorf("No Spec.SuggestionSpec specified.")
	}
	namespace := instance.Namespace
	BUFSIZE := 1024
	wkind, err := getWorkerKind(instance.Spec.WorkerSpec)
	if err != nil {
		log.Printf("getWorkerKind error %v", err)
		return err
	}

	studyID := "studyID4Validation"
	trialID := "trialID4Validation"
	workerID, wm, err := getWorkerManifest(
		nil,
		studyID,
		&katibapi.Trial{
			TrialId:      trialID,
			ParameterSet: []*katibapi.Parameter{},
		},
		instance.Spec.WorkerSpec,
		wkind.Kind,
		namespace,
		true,
	)
	if err != nil {
		return err
	}

	job := &unstructured.Unstructured{}
	if err := k8syaml.NewYAMLOrJSONDecoder(wm, BUFSIZE).Decode(job); err != nil {
		log.Printf("Yaml decode error %v", err)
		return err
	}

	if job.GetNamespace() != namespace || job.GetName() != workerID {
		return fmt.Errorf("Invalid worker template.")
	}

	var mcjob batchv1beta.CronJob
	mcm, err := getMetricsCollectorManifest(studyID, trialID, workerID, wkind.Kind, namespace, instance.Spec.MetricsCollectorSpec)
	if err != nil {
		log.Printf("getMetricsCollectorManifest error %v", err)
		return err
	}

	if err := k8syaml.NewYAMLOrJSONDecoder(mcm, BUFSIZE).Decode(&mcjob); err != nil {
		log.Printf("MetricsCollector Yaml decode error %v", err)
		return err
	}

	if mcjob.GetNamespace() != namespace || mcjob.GetName() != workerID {
		return fmt.Errorf("Invalid metricsCollector template.")
	}
	return nil
}

func checkGoalAndUpdateObject(curValue float64, instance *katibv1alpha1.StudyJob, workerId string) bool {
	optTypeFuncMap := map[katibv1alpha1.OptimizationType]func(float64, float64) bool{
		katibv1alpha1.OptimizationTypeMinimize: func(a, b float64) bool { return a < b },
		katibv1alpha1.OptimizationTypeMaximize: func(a, b float64) bool { return a > b },
	}
	goal := false
	if optTypeFuncMap[instance.Spec.OptimizationType] == nil {
		return false
	}
	var trialId string
OuterLoop:
	for i := range instance.Status.Trials {
		for j := range instance.Status.Trials[i].WorkerList {
			if instance.Status.Trials[i].WorkerList[j].WorkerID == workerId {
				instance.Status.Trials[i].WorkerList[j].ObjectiveValue = &curValue
				trialId = instance.Status.Trials[i].TrialID
				break OuterLoop
			}
		}
	}
	opFunc := optTypeFuncMap[instance.Spec.OptimizationType]
	if opFunc(curValue, *instance.Spec.OptimizationGoal) {
		goal = true
	}
	if instance.Status.BestObjectiveValue != nil {
		if opFunc(curValue, *instance.Status.BestObjectiveValue) {
			instance.Status.BestObjectiveValue = &curValue
			instance.Status.BestTrialID = trialId
			instance.Status.BestWorkerID = workerId
		}
	} else {
		instance.Status.BestObjectiveValue = &curValue
		instance.Status.BestTrialID = trialId
		instance.Status.BestWorkerID = workerId
	}

	return goal
}

func contains(l []string, s string) bool {
	for _, elem := range l {
		if elem == s {
			return true
		}
	}
	return false
}

func getMyNamespace() string {
	data, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	return strings.TrimSpace(string(data))
}
