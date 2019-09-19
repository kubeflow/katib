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

	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	"github.com/kubeflow/katib/pkg/util/v1alpha3/katibclient"
)

const (
	timeout = 20 * time.Minute
)

func verifyResult(exp *experimentsv1alpha3.Experiment) error {
	if len(exp.Status.CurrentOptimalTrial.ParameterAssignments) == 0 {
		return fmt.Errorf("Best parameter assignments not updated in status")
	}

	if len(exp.Status.CurrentOptimalTrial.Observation.Metrics) == 0 {
		return fmt.Errorf("Bst metrics not updated in status")
	}

	metric := exp.Status.CurrentOptimalTrial.Observation.Metrics[0]
	if metric.Name != exp.Spec.Objective.ObjectiveMetricName {
		return fmt.Errorf("Best objective metric not updated in status")
	}
	return nil
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
	var maxtrials int32 = 3
	var paralleltrials int32 = 2
	exp.Spec.MaxTrialCount = &maxtrials
	exp.Spec.ParallelTrialCount = &paralleltrials
	err = kclient.CreateExperiment(exp)
	if err != nil {
		log.Fatal("CreateExperiment from YAML failed: ", err)
	}

	for endTime := time.Now().Add(timeout); time.Now().Before(endTime); {
		exp, err = kclient.GetExperiment(exp.Name, exp.Namespace)
		if err != nil {
			log.Fatal("Get Experiment error ", err)
		}
		if exp.IsCompleted() {
			log.Printf("Experiment %v finished", exp.Name)
			break
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
		time.Sleep(20 * time.Second)
	}

	if !exp.IsCompleted() {
		log.Fatal("Experiment run timed out")
	}

	if exp.Status.Trials != *exp.Spec.MaxTrialCount {
		log.Fatal("All trials are not run in the experiment ", exp.Status.Trials, exp.Spec.MaxTrialCount)
	}

	if exp.Status.TrialsSucceeded != *exp.Spec.MaxTrialCount {
		log.Fatal("All trials are not successful ", exp.Status.TrialsSucceeded, *exp.Spec.MaxTrialCount)
	}
	err = verifyResult(exp)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Experiment has recorded best current Optimal Trial %v", exp.Status.CurrentOptimalTrial)
	err = kclient.DeleteExperiment(exp)
	if err != nil {
		log.Printf("CreateExperiment from YAML failed: %v", err)
		return
	}
	log.Printf("Experiment %v deleted", exp.Name)
}
