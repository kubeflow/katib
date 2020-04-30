package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"os"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	commonv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	controllerUtil "github.com/kubeflow/katib/pkg/controller.v1alpha3/util"
	"github.com/kubeflow/katib/pkg/util/v1alpha3/katibclient"
)

const (
	timeout = 30 * time.Minute
)

func verifyResult(exp *experimentsv1alpha3.Experiment) (*float64, error) {
	if len(exp.Status.CurrentOptimalTrial.ParameterAssignments) == 0 {
		return nil, fmt.Errorf("Best parameter assignments not updated in status")
	}

	if len(exp.Status.CurrentOptimalTrial.Observation.Metrics) == 0 {
		return nil, fmt.Errorf("Best metrics not updated in status")
	}

	metric := exp.Status.CurrentOptimalTrial.Observation.Metrics[0]
	if metric.Name != exp.Spec.Objective.ObjectiveMetricName {
		return nil, fmt.Errorf("Best objective metric not updated in status")
	}
	return &metric.Value, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Experiment name is missing")
	}
	expName := os.Args[1]
	b, err := ioutil.ReadFile(expName)
	if err != nil {
		log.Fatal("Error in reading file ", err)
	}
	exp := &experimentsv1alpha3.Experiment{}
	buf := bytes.NewBufferString(string(b))
	if err = k8syaml.NewYAMLOrJSONDecoder(buf, 1024).Decode(exp); err != nil {
		log.Fatal("Yaml decode error ", err)
	}
	kclient, err := katibclient.NewClient(client.Options{})
	if err != nil {
		log.Fatal("NewClient for Katib failed: ", err)
	}
	exp, err = kclient.GetExperiment(exp.Name, exp.Namespace)
	if err != nil {
		log.Fatal("Get Experiment error. Experiment not created yet ", err)
	}
	if exp.Spec.Algorithm.AlgorithmName != "hyperband" {
		// Hyperband will validate the parallel trial count,
		// thus we should not change it.
		var maxtrials int32 = 7
		var paralleltrials int32 = 3
		exp.Spec.MaxTrialCount = &maxtrials
		exp.Spec.ParallelTrialCount = &paralleltrials
	}
	err = kclient.UpdateExperiment(exp)
	if err != nil {
		log.Fatal("UpdateExperiment from YAML failed: ", err)
	}
	endTime := time.Now().Add(timeout)
	for time.Now().Before(endTime) {
		log.Printf("Waiting for Experiment %s to start running.", exp.Name)
		exp, err = kclient.GetExperiment(exp.Name, exp.Namespace)
		if err != nil {
			log.Fatal("Get Experiment error ", err)
		}
		if exp.IsRunning() {
			log.Printf("Experiment %v started running", exp.Name)
			break
		}
		time.Sleep(5 * time.Second)
	}

	for time.Now().Before(endTime) {
		exp, err = kclient.GetExperiment(exp.Name, exp.Namespace)
		if err != nil {
			log.Fatal("Get Experiment error ", err)
		}
		log.Printf("Waiting for Experiment %s to finish.", exp.Name)
		log.Printf(`Experiment %s's trials: %d trials, %d pending trials,
%d running trials, %d killed trials, %d succeeded trials, %d failed trials.`,
			exp.Name,
			exp.Status.Trials, exp.Status.TrialsPending, exp.Status.TrialsRunning,
			exp.Status.TrialsKilled, exp.Status.TrialsSucceeded, exp.Status.TrialsFailed)
		log.Printf("Optimal Trial for Experiment %s: %v", exp.Name,
			exp.Status.CurrentOptimalTrial)
		log.Printf("Experiment %s's conditions: %v", exp.Name, exp.Status.Conditions)

		suggestion, err := kclient.GetSuggestion(exp.Name, exp.Namespace)
		if err != nil {
			log.Printf("Get Suggestion error: %v", err)
		} else {
			log.Printf("Suggestion %s's conditions: %v", suggestion.Name,
				suggestion.Status.Conditions)
			log.Printf("Suggestion %s's suggestions: %v", suggestion.Name,
				suggestion.Status.Suggestions)
		}
		if exp.IsCompleted() {
			log.Printf("Experiment %v finished", exp.Name)
			break
		}
		time.Sleep(20 * time.Second)
	}

	if !exp.IsCompleted() {
		log.Fatal("Experiment run timed out")
	}

	metricVal, err := verifyResult(exp)
	if err != nil {
		log.Fatal(err)
	}
	if metricVal == nil {
		log.Fatal("Metric value in CurrentOptimalTrial not populated")
	}

	objectiveType := exp.Spec.Objective.Type
	var goal float64
	if exp.Spec.Objective.Goal != nil {
		goal = *exp.Spec.Objective.Goal
	}
	if (exp.Spec.Objective.Goal != nil && objectiveType == commonv1alpha3.ObjectiveTypeMinimize && *metricVal < goal) ||
		(exp.Spec.Objective.Goal != nil && objectiveType == commonv1alpha3.ObjectiveTypeMaximize && *metricVal > goal) {
		log.Print("Objective Goal reached")
	} else {

		if exp.Status.Trials != *exp.Spec.MaxTrialCount {
			log.Fatal("All trials are not run in the experiment ", exp.Status.Trials, exp.Spec.MaxTrialCount)
		}

		if exp.Status.TrialsSucceeded != *exp.Spec.MaxTrialCount {
			log.Fatal("All trials are not successful ", exp.Status.TrialsSucceeded, *exp.Spec.MaxTrialCount)
		}
	}

	sug, err := kclient.GetSuggestion(exp.Name, exp.Namespace)
	if exp.Spec.ResumePolicy == experimentsv1alpha3.LongRunning {
		if sug.IsSucceeded() {
			log.Fatal("Suggestion is terminated while ResumePolicy = LongRunning")
		}
	}
	if exp.Spec.ResumePolicy == experimentsv1alpha3.NeverResume {
		if sug.IsRunning() {
			log.Fatal("Suggestion is still running while ResumePolicy = NeverResume")
		}
		namespacedName := types.NamespacedName{Name: controllerUtil.GetAlgorithmServiceName(sug), Namespace: sug.Namespace}
		service := &corev1.Service{}
		err := kclient.GetClient().Get(context.TODO(), namespacedName, service)
		if err == nil || !errors.IsNotFound(err) {
			log.Fatal("Suggestion service is still alive while ResumePolicy = NeverResume")
		}
		namespacedName = types.NamespacedName{Name: controllerUtil.GetAlgorithmDeploymentName(sug), Namespace: sug.Namespace}
		deployment := &appsv1.Deployment{}
		err = kclient.GetClient().Get(context.TODO(), namespacedName, deployment)
		if err == nil || !errors.IsNotFound(err) {
			log.Fatal("Suggestion deployment is still alive while ResumePolicy = NeverResume")
		}
	}

	log.Printf("Experiment has recorded best current Optimal Trial %v", exp.Status.CurrentOptimalTrial)
}
