package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	controllerUtil "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibclient"
)

const (
	timeout = 30 * time.Minute
)

func verifyResult(exp *experimentsv1beta1.Experiment) (*commonv1beta1.Metric, error) {
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
	exp := &experimentsv1beta1.Experiment{}
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

	var maxtrials int32 = 7
	var paralleltrials int32 = 3
	if exp.Spec.Algorithm.AlgorithmName != "hyperband" && exp.Spec.Algorithm.AlgorithmName != "darts" {
		// Hyperband will validate the parallel trial count,
		// thus we should not change it.
		// Not necessary to test parallel Trials for Darts
		exp.Spec.MaxTrialCount = &maxtrials
		exp.Spec.ParallelTrialCount = &paralleltrials
	}
	err = kclient.UpdateRuntimeObject(exp)
	if err != nil {
		log.Fatal("UpdateRuntimeObject failed: ", err)
	}
	endTime := time.Now().Add(timeout)
	for time.Now().Before(endTime) {
		log.Printf("Waiting for Experiment %s to start running.", exp.Name)
		exp, err = kclient.GetExperiment(exp.Name, exp.Namespace)
		if err != nil {
			log.Fatal("Get Experiment error ", err)
		}
		if exp.IsRunning() {
			log.Printf("Experiment %v started running with %v MaxTrialCount", exp.Name, maxtrials)
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
	// If min metric is set, max be set also
	minMetric, err := strconv.ParseFloat(metric.Min, 64)
	maxMetric, _ := strconv.ParseFloat(metric.Max, 64)
	if err == nil &&
		((exp.Spec.Objective.Goal != nil && objectiveType == commonv1beta1.ObjectiveTypeMinimize && minMetric < goal) ||
			(exp.Spec.Objective.Goal != nil && objectiveType == commonv1beta1.ObjectiveTypeMaximize && maxMetric > goal)) {
	} else {

		if exp.Status.Trials != *exp.Spec.MaxTrialCount {
			log.Fatal("All trials are not run in the experiment ", exp.Status.Trials, exp.Spec.MaxTrialCount)
		}

		if exp.Status.TrialsSucceeded != *exp.Spec.MaxTrialCount {
			log.Fatal("All trials are not successful ", exp.Status.TrialsSucceeded, *exp.Spec.MaxTrialCount)
		}
	}

	sug, err := kclient.GetSuggestion(exp.Name, exp.Namespace)
	if exp.Spec.ResumePolicy == experimentsv1beta1.LongRunning {
		if sug.IsSucceeded() {
			log.Fatal("Suggestion is succeeded while ResumePolicy = LongRunning")
		}
	}
	if exp.Spec.ResumePolicy == experimentsv1beta1.NeverResume || exp.Spec.ResumePolicy == experimentsv1beta1.FromVolume {
		if sug.IsRunning() {
			log.Fatalf("Suggestion is still running while ResumePolicy = %v", exp.Spec.ResumePolicy)
		}

		namespacedName := types.NamespacedName{Name: controllerUtil.GetSuggestionServiceName(sug), Namespace: sug.Namespace}
		err := kclient.GetClient().Get(context.TODO(), namespacedName, &corev1.Service{})
		if err == nil || !errors.IsNotFound(err) {
			log.Fatalf("Suggestion service is still alive while ResumePolicy = %v", exp.Spec.ResumePolicy)
		}
		log.Printf("Suggestion service %v has been deleted", controllerUtil.GetSuggestionServiceName(sug))

		namespacedName = types.NamespacedName{Name: controllerUtil.GetSuggestionDeploymentName(sug), Namespace: sug.Namespace}
		err = kclient.GetClient().Get(context.TODO(), namespacedName, &appsv1.Deployment{})
		if err == nil || !errors.IsNotFound(err) {
			log.Fatalf("Suggestion deployment is still alive while ResumePolicy = %v", exp.Spec.ResumePolicy)
		}
		log.Printf("Suggestion deployment %v has been deleted", controllerUtil.GetSuggestionDeploymentName(sug))

		if exp.Spec.ResumePolicy == experimentsv1beta1.FromVolume {
			namespacedName = types.NamespacedName{Name: controllerUtil.GetSuggestionPersistentVolumeClaimName(sug), Namespace: sug.Namespace}
			err = kclient.GetClient().Get(context.TODO(), namespacedName, &corev1.PersistentVolumeClaim{})
			if err != nil {
				log.Fatalf("Suggestion persistent volume claim is not alive while ResumePolicy = %v, error: %v", experimentsv1beta1.FromVolume, err)
			}

			namespacedName = types.NamespacedName{Name: controllerUtil.GetSuggestionPersistentVolumeClaimName(sug)}
			err = kclient.GetClient().Get(context.TODO(), namespacedName, &corev1.PersistentVolume{})
			if err != nil {
				log.Fatalf("Suggestion persistent volume is not alive while ResumePolicy = %v, error: %v", experimentsv1beta1.FromVolume, err)
			}
		}
	}

	log.Printf("Experiment has recorded best current Optimal Trial %v", exp.Status.CurrentOptimalTrial)
}
