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

package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // For GCP testing
	"sigs.k8s.io/controller-runtime/pkg/client"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	controllerUtil "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibclient"
)

func main() {
	// For AWS we should point KUBECONFIG env to correct folder.
	err := os.Setenv("KUBECONFIG", "/root/.kube/config")
	if err != nil {
		log.Fatalf("Unable to set KUBECONFIG env variable, error: %v", err)
	}

	// First argument should be Experiment yaml path.
	if len(os.Args) != 2 {
		log.Fatal("Path to Experiment yaml is missing")
	}
	expPath := os.Args[1]
	byteExp, err := os.ReadFile(expPath)
	if err != nil {
		log.Fatalf("Error in reading file: %v", err)
	}

	// Replace batch size to number of epochs for faster execution.
	strExp := strings.Replace(string(byteExp), "--batch-size=64", "--num-epochs=2", -1)

	exp := &experimentsv1beta1.Experiment{}
	buf := bytes.NewBufferString(strExp)
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, 1024).Decode(exp); err != nil {
		log.Fatal("Yaml decode error ", err)
	}

	kclient, err := katibclient.NewClient(client.Options{})
	if err != nil {
		log.Fatal("Create NewClient for Katib failed: ", err)
	}

	var maxTrials int32 = 2
	var parallelTrials int32 = 1
	// For random we test 2 parallel execution.
	if exp.Name == "random" {
		maxTrials = 3
		parallelTrials = 2
	}
	if exp.Spec.Algorithm.AlgorithmName != "hyperband" && exp.Spec.Algorithm.AlgorithmName != "darts" {
		// Hyperband will validate the parallel trial count,
		// thus we should not change it.
		// Not necessary to test parallel Trials for Darts.
		exp.Spec.MaxTrialCount = &maxTrials
		exp.Spec.ParallelTrialCount = &parallelTrials
	}
	log.Printf("Creating Experiment %v with MaxTrialCount: %v, ParallelTrialCount: %v", exp.Name, maxTrials, parallelTrials)
	err = kclient.CreateRuntimeObject(exp)
	if err != nil {
		log.Fatalf("CreateRuntimeObject failed: %v", err)
	}

	// Wait until Experiment is finished.
	exp, err = waitExperimentFinish(kclient, exp)
	if err != nil {
		// Delete Experiment in case of error.
		log.Printf("Deleting Experiment %v\n", exp.Name)
		if kclient.DeleteRuntimeObject(exp) != nil {
			log.Fatalf("Unable to delete Experiment %v, error: %v", exp.Name, err)
		}
		log.Fatalf("Wait Experiment finish failed: %v", err)
	}

	// For random example and from volume we restart Experiment.
	if exp.Name == "random" || exp.Name == "from-volume-resume" {
		// Increase parallel Trials and max Trials counts.
		parallelTrials++
		maxTrials += parallelTrials + 1
		exp.Spec.MaxTrialCount = &maxTrials
		exp.Spec.ParallelTrialCount = &parallelTrials

		// Update Experiment with new info.
		err := kclient.UpdateRuntimeObject(exp)
		if err != nil {
			log.Fatalf("UpdateRuntimeObject failed: %v", err)
			// Delete Experiment in case of error.
			log.Printf("Deleting Experiment %v\n", exp.Name)
			if kclient.DeleteRuntimeObject(exp) != nil {
				log.Fatalf("Unable to delete Experiment %v, error: %v", exp.Name, err)
			}
		}

		log.Printf("Restarting Experiment %v with MaxTrialCount: %v, ParallelTrialCount: %v\n\n",
			exp.Name, maxTrials, parallelTrials)

		// Wait until Experiment is restarted.
		timeout := 60 * time.Second
		endTime := time.Now().Add(timeout)
		for time.Now().Before(endTime) {
			exp, err = kclient.GetExperiment(exp.Name, exp.Namespace)
			if err != nil {
				log.Fatalf("Get Experiment error: %v", err)
				// Delete Experiment in case of error
				log.Printf("Deleting Experiment %v\n", exp.Name)
				if kclient.DeleteRuntimeObject(exp) != nil {
					log.Fatalf("Unable to delete Experiment %v, error: %v", exp.Name, err)
				}
			}
			// Once Experiment is restarted stop the waiting loop.
			if exp.IsRestarting() {
				break
			}
			time.Sleep(1 * time.Second)
		}
		// Check if Experiment is not restarting and is not running.
		if !exp.IsRestarting() && !exp.IsRunning() {
			log.Fatalf("Unable to restart Experiment %v, Experiment conditions: %v", exp.Name, exp.Status.Conditions)
			// Delete experiment in case of error.
			log.Printf("Deleting Experiment %v\n", exp.Name)
			if kclient.DeleteRuntimeObject(exp) != nil {
				log.Fatalf("Unable to delete Experiment %v, error: %v", exp.Name, err)
			}
		}
		// Wait until Experiment is finished.
		exp, err = waitExperimentFinish(kclient, exp)
		if err != nil {
			log.Fatalf("Wait Experiment finish failed: %v", err)
			// Delete experiment in case of error
			log.Printf("Deleting Experiment %v\n", exp.Name)
			if kclient.DeleteRuntimeObject(exp) != nil {
				log.Fatalf("Unable to delete Experiment %v, error: %v", exp.Name, err)
			}
		}
	}

	// Verify Experiment results
	err = verifyExperimentResults(kclient, exp)
	if err != nil {
		// Delete Experiment in case of error
		log.Printf("Deleting Experiment %v\n", exp.Name)
		if kclient.DeleteRuntimeObject(exp) != nil {
			log.Fatalf("Unable to delete Experiment %v, error: %v", exp.Name, err)
		}
		log.Fatalf("Verify Experiment results failed: %v", err)
	}

	// Print results.
	err = printResults(exp)
	if err != nil {
		// Delete Experiment in case of error.
		log.Printf("Deleting Experiment %v\n", exp.Name)
		if kclient.DeleteRuntimeObject(exp) != nil {
			log.Fatalf("Unable to delete Experiment %v, error: %v", exp.Name, err)
		}
		log.Fatalf("Print Experiment results failed: %v", err)
	}

	// Delete Experiment.
	log.Printf("Deleting Experiment %v\n", exp.Name)
	if kclient.DeleteRuntimeObject(exp) != nil {
		log.Fatalf("Unable to delete Experiment %v, error: %v", exp.Name, err)
	}

}

func waitExperimentFinish(kclient katibclient.Client, exp *experimentsv1beta1.Experiment) (*experimentsv1beta1.Experiment, error) {
	// Experiment should be completed before the timeout.
	timeout := 50 * time.Minute
	for endTime := time.Now().Add(timeout); time.Now().Before(endTime); {
		exp, err := kclient.GetExperiment(exp.Name, exp.Namespace)
		if err != nil {
			return exp, fmt.Errorf("Get Experiment error: %v", err)
		}

		log.Printf("Waiting for Experiment %s to finish", exp.Name)
		log.Printf(`Experiment is running: %v Trials, %v Pending Trials, %v Running Trials, %v Succeeded Trials, %v Failed Trials, %v EarlyStopped Trials`,
			exp.Status.Trials, exp.Status.TrialsPending, exp.Status.TrialsRunning, exp.Status.TrialsSucceeded, exp.Status.TrialsFailed, exp.Status.TrialsEarlyStopped)
		log.Printf("Current optimal Trial: %v", exp.Status.CurrentOptimalTrial)
		log.Printf("Experiment conditions: %v\n\n\n", exp.Status.Conditions)

		// Check if Experiment is completed.
		if exp.IsCompleted() {
			log.Printf("Experiment %v is finished", exp.Name)
			if exp.IsFailed() {
				return exp, fmt.Errorf("Experiment %v is failed", exp.Name)
			}
			// Print latest condition message.
			log.Printf("%v\n\n", exp.Status.Conditions[len(exp.Status.Conditions)-1].Message)
			// Print Suggestion conditions.
			suggestion, err := kclient.GetSuggestion(exp.Name, exp.Namespace)
			if err != nil {
				return exp, fmt.Errorf("Get Suggestion error: %v", err)
			}
			log.Printf("Suggestion %v. Conditions: %v", suggestion.Name, suggestion.Status.Conditions)
			log.Printf("Suggestion %v. Suggestions: %v\n\n", suggestion.Name, suggestion.Status.Suggestions)

			// Return succeeded Experiment.
			return exp, nil
		}
		time.Sleep(20 * time.Second)
	}

	// If loop is end, Experiment is not finished.
	return exp, fmt.Errorf("Experiment run timed out")
}

func verifyExperimentResults(kclient katibclient.Client, exp *experimentsv1beta1.Experiment) error {

	// Current optimal Trial should be set.
	if equality.Semantic.DeepEqual(exp.Status.CurrentOptimalTrial, experimentsv1beta1.OptimalTrial{}) {
		return fmt.Errorf("Current optimal Trial is empty. Experiment status: %v", exp.Status)
	}

	// Best objective metric should be set.
	var bestObjectiveMetric *commonv1beta1.Metric
	for _, metric := range exp.Status.CurrentOptimalTrial.Observation.Metrics {
		if metric.Name == exp.Spec.Objective.ObjectiveMetricName {
			bestObjectiveMetric = &metric
			break
		}
	}
	if bestObjectiveMetric == nil {
		return fmt.Errorf("Unable to get best metrics for objective: %v", exp.Spec.Objective.ObjectiveMetricName)
	}

	// Verify objective metric.
	objectiveType := exp.Spec.Objective.Type
	goal := exp.Spec.Objective.Goal
	// If min metric is set, max be set also.
	minMetric, err := strconv.ParseFloat(bestObjectiveMetric.Min, 64)
	maxMetric, _ := strconv.ParseFloat(bestObjectiveMetric.Max, 64)

	// If metrics can't be parsed to float or goal is empty, succeeded Trials should be equal to MaxTrialCount.
	if (err != nil || goal == nil) && exp.Status.TrialsSucceeded != *exp.Spec.MaxTrialCount {
		return fmt.Errorf("All trials are not successful. MaxTrialCount: %v, TrialsSucceeded: %v",
			*exp.Spec.MaxTrialCount, exp.Status.TrialsSucceeded)
	}

	trialsCompleted := exp.Status.TrialsSucceeded
	if exp.Spec.EarlyStopping != nil {
		trialsCompleted += exp.Status.TrialsEarlyStopped
	}

	// Otherwise, Goal should be reached.
	if trialsCompleted != *exp.Spec.MaxTrialCount &&
		((objectiveType == commonv1beta1.ObjectiveTypeMinimize && minMetric > *goal) ||
			(objectiveType == commonv1beta1.ObjectiveTypeMaximize && maxMetric < *goal)) {
		return fmt.Errorf(`Objective Goal is not reached and Succeeded Trials: %v != %v MaxTrialCount.
			ObjectiveType: %v, Goal: %v, MinMetric: %v, MaxMetric: %v`,
			exp.Status.TrialsSucceeded, *exp.Spec.MaxTrialCount,
			objectiveType, *goal, minMetric, maxMetric)
	}

	err = verifySuggestion(kclient, exp)
	if err != nil {
		return fmt.Errorf("Verify Suggestion failed: %v", err)
	}
	return nil
}

func verifySuggestion(kclient katibclient.Client, exp *experimentsv1beta1.Experiment) error {

	// Verify Suggestion's resources.
	sug, err := kclient.GetSuggestion(exp.Name, exp.Namespace)
	if err != nil {
		return fmt.Errorf("GetSuggestion failed: %v", err)
	}

	// When Suggestion is LongRunning, it can't be succeeded.
	if exp.Spec.ResumePolicy == experimentsv1beta1.LongRunning && sug.IsSucceeded() {
		return fmt.Errorf("Suggestion is succeeded while ResumePolicy = %v", experimentsv1beta1.LongRunning)
	}

	// Verify Suggestion with resume policy Never and FromVolume.
	if exp.Spec.ResumePolicy == experimentsv1beta1.NeverResume || exp.Spec.ResumePolicy == experimentsv1beta1.FromVolume {

		// Give controller time to delete Suggestion resources and change Suggestion status.
		// TODO (andreyvelich): Think about better way to handle this.
		time.Sleep(10 * time.Second)

		// When Suggestion has resume policy Never or FromVolume, it should be not running.
		if sug.IsRunning() {
			return fmt.Errorf("Suggestion is still running while ResumePolicy = %v", exp.Spec.ResumePolicy)
		}

		// Suggestion service should be deleted.
		serviceName := controllerUtil.GetSuggestionServiceName(sug)
		namespacedName := types.NamespacedName{Name: serviceName, Namespace: sug.Namespace}
		err = kclient.GetClient().Get(context.TODO(), namespacedName, &corev1.Service{})
		if errors.IsNotFound(err) {
			log.Printf("Suggestion service %v has been deleted", serviceName)
		} else {
			return fmt.Errorf("Suggestion service: %v is still alive while ResumePolicy: %v, error: %v", serviceName, exp.Spec.ResumePolicy, err)
		}

		// Suggestion deployment should be deleted.
		deploymentName := controllerUtil.GetSuggestionDeploymentName(sug)
		namespacedName = types.NamespacedName{Name: deploymentName, Namespace: sug.Namespace}
		err = kclient.GetClient().Get(context.TODO(), namespacedName, &appsv1.Deployment{})
		if errors.IsNotFound(err) {
			log.Printf("Suggestion deployment %v has been deleted", deploymentName)
		} else {
			return fmt.Errorf("Suggestion deployment: %v is still alive while ResumePolicy: %v, error: %v", deploymentName, exp.Spec.ResumePolicy, err)
		}

		// PVC should not be deleted for Suggestion with resume policy FromVolume.
		if exp.Spec.ResumePolicy == experimentsv1beta1.FromVolume {
			pvcName := controllerUtil.GetSuggestionPersistentVolumeClaimName(sug)
			namespacedName = types.NamespacedName{Name: pvcName, Namespace: sug.Namespace}
			err = kclient.GetClient().Get(context.TODO(), namespacedName, &corev1.PersistentVolumeClaim{})
			if errors.IsNotFound(err) {
				return fmt.Errorf("suggestion PVC: %v is not alive while ResumePolicy: %v", pvcName, exp.Spec.ResumePolicy)
			}
		}
	}
	return nil
}

func printResults(exp *experimentsv1beta1.Experiment) error {
	log.Printf("Experiment has recorded best current Optimal Trial %v\n\n", exp.Status.CurrentOptimalTrial)

	// Describe the Experiment.
	cmd := exec.Command("kubectl", "describe", "experiment", exp.Name, "-n", exp.Namespace)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Execute \"kubectl describe suggestion\" failed: %v", err)
	}
	log.Println(cmd.String())
	log.Printf("\n%v\n\n", string(out))

	// Describe the Suggestion.
	cmd = exec.Command("kubectl", "describe", "suggestion", exp.Name, "-n", exp.Namespace)
	out, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("Execute \"kubectl describe experiment\" failed: %v", err)
	}
	log.Println(cmd.String())
	log.Printf("\n%v", string(out))

	return nil
}
