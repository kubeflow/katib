package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	commonv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	"github.com/kubeflow/katib/pkg/util/v1alpha3/katibclient"
)

const (
	timeout = 30 * time.Minute
)

func verifyResult(exp *experimentsv1alpha3.Experiment) (*commonv1alpha3.Metric, error) {
	if len(exp.Status.CurrentOptimalTrial.ParameterAssignments) == 0 {
		return nil, fmt.Errorf("Best parameter assignments not updated in status")
	}

	if len(exp.Status.CurrentOptimalTrial.Observation.Metrics) == 0 {
		return nil, fmt.Errorf("Best metrics not updated in status")
	}

	for _, metric := range exp.Status.CurrentOptimalTrial.Observation.Metrics {
		if metric.Name == exp.Spec.Objective.ObjectiveMetricName {
			return &metric, nil
		}
	}

	return nil, fmt.Errorf("Best objective metric not updated in status")
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
	if exp.Spec.Algorithm.AlgorithmName != "hyperband" {
		// Hyperband will validate the parallel trial count,
		// thus we should not change it.
		var maxtrials int32 = 3
		var paralleltrials int32 = 2
		exp.Spec.MaxTrialCount = &maxtrials
		exp.Spec.ParallelTrialCount = &paralleltrials
	}
	err = kclient.CreateExperiment(exp)
	if err != nil {
		log.Fatal("CreateExperiment from YAML failed: ", err)
	}

	for endTime := time.Now().Add(timeout); time.Now().Before(endTime); {
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

	metric, err := verifyResult(exp)
	if err != nil {
		log.Fatal(err)
	}
	if metric == nil {
		log.Fatal("Metric value in CurrentOptimalTrial not populated")
	}

	objectiveType := exp.Spec.Objective.Type
	var goal float64
	if exp.Spec.Objective.Goal != nil {
		goal = *exp.Spec.Objective.Goal
	}
	if (exp.Spec.Objective.Goal != nil && objectiveType == commonv1alpha3.ObjectiveTypeMinimize && metric.Min < goal) ||
		(exp.Spec.Objective.Goal != nil && objectiveType == commonv1alpha3.ObjectiveTypeMaximize && metric.Max > goal) {
		log.Print("Objective Goal reached")
	} else {

		if exp.Status.Trials != *exp.Spec.MaxTrialCount {
			log.Fatal("All trials are not run in the experiment ", exp.Status.Trials, exp.Spec.MaxTrialCount)
		}

		if exp.Status.TrialsSucceeded != *exp.Spec.MaxTrialCount {
			log.Fatal("All trials are not successful ", exp.Status.TrialsSucceeded, *exp.Spec.MaxTrialCount)
		}
	}
	log.Printf("Experiment has recorded best current Optimal Trial %v", exp.Status.CurrentOptimalTrial)
}
