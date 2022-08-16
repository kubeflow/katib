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

package experiment

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
)

const (
	updatePrometheusMetrics = "update-prometheus-metrics"
)

func (r *ReconcileExperiment) getTrialInstance(expInstance *experimentsv1beta1.Experiment, trialAssignment *suggestionsv1beta1.TrialAssignment) (*trialsv1beta1.Trial, error) {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: expInstance.GetName(), Namespace: expInstance.GetNamespace()})

	trial := &trialsv1beta1.Trial{}
	trial.Name = trialAssignment.Name
	trial.Namespace = expInstance.GetNamespace()
	trial.Labels = util.TrialLabels(expInstance)

	if trialAssignment.Labels != nil {
		for k, v := range trialAssignment.Labels {
			trial.Labels[k] = v
		}
	}

	if err := controllerutil.SetControllerReference(expInstance, trial, r.scheme); err != nil {
		logger.Error(err, "Set controller reference error")
		return nil, err
	}

	trial.Spec.Objective = expInstance.Spec.Objective

	hps := trialAssignment.ParameterAssignments
	trial.Spec.ParameterAssignments = trialAssignment.ParameterAssignments

	if expInstance.Spec.EarlyStopping != nil {
		trial.Spec.EarlyStoppingRules = trialAssignment.EarlyStoppingRules
	}

	runSpec, err := r.GetRunSpecWithHyperParameters(expInstance, trial.Name, trial.Namespace, hps)
	if err != nil {
		logger.Error(err, "Fail to get RunSpec from experiment", expInstance.Name)
		return nil, err
	}

	trial.Spec.RunSpec = runSpec
	if expInstance.Spec.TrialTemplate != nil {
		trial.Spec.RetainRun = expInstance.Spec.TrialTemplate.Retain
	}

	if expInstance.Spec.MetricsCollectorSpec != nil {
		trial.Spec.MetricsCollector = *expInstance.Spec.MetricsCollectorSpec
	}

	if expInstance.Spec.TrialTemplate.PrimaryPodLabels != nil {
		trial.Spec.PrimaryPodLabels = expInstance.Spec.TrialTemplate.PrimaryPodLabels
	}

	if expInstance.Spec.TrialTemplate.PrimaryContainerName != "" {
		trial.Spec.PrimaryContainerName = expInstance.Spec.TrialTemplate.PrimaryContainerName
	}

	if expInstance.Spec.TrialTemplate.SuccessCondition != "" && expInstance.Spec.TrialTemplate.FailureCondition != "" {
		trial.Spec.SuccessCondition = expInstance.Spec.TrialTemplate.SuccessCondition
		trial.Spec.FailureCondition = expInstance.Spec.TrialTemplate.FailureCondition
	}

	return trial, nil
}

func needUpdateFinalizers(exp *experimentsv1beta1.Experiment) (bool, []string) {
	deleted := !exp.ObjectMeta.DeletionTimestamp.IsZero()
	pendingFinalizers := exp.GetFinalizers()
	contained := false
	for _, elem := range pendingFinalizers {
		if elem == updatePrometheusMetrics {
			contained = true
			break
		}
	}

	if !deleted && !contained {
		finalizers := append(pendingFinalizers, updatePrometheusMetrics)
		return true, finalizers
	}
	if deleted && contained {
		finalizers := []string{}
		for _, pendingFinalizer := range pendingFinalizers {
			if pendingFinalizer != updatePrometheusMetrics {
				finalizers = append(finalizers, pendingFinalizer)
			}
		}
		return true, finalizers
	}
	return false, []string{}
}

func (r *ReconcileExperiment) updateFinalizers(instance *experimentsv1beta1.Experiment, finalizers []string) (reconcile.Result, error) {
	instance.SetFinalizers(finalizers)
	if err := r.Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	} else {
		if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
			r.collector.IncreaseExperimentsDeletedCount(instance.Namespace)
		} else {
			r.collector.IncreaseExperimentsCreatedCount(instance.Namespace)
		}
		// Need to requeue because finalizer update does not change metadata.generation
		return reconcile.Result{Requeue: true}, err
	}
}

func (r *ReconcileExperiment) cleanupSuggestionResources(instance *experimentsv1beta1.Experiment) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	original := &suggestionsv1beta1.Suggestion{}
	err := r.Get(context.TODO(),
		types.NamespacedName{Namespace: instance.GetNamespace(), Name: instance.GetName()}, original)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// If Suggestion is completed or Suggestion is restarting not needed to terminate Suggestion
	if original.IsCompleted() || original.IsRestarting() {
		return nil
	}

	logger.Info("Start cleanup suggestion resources")
	suggestion := original.DeepCopy()

	reason := "Experiment is succeeded"
	// If ResumePolicy = Never, mark suggestion status succeeded, can't be restarted
	if instance.Spec.ResumePolicy == experimentsv1beta1.NeverResume {
		msg := "Suggestion is succeeded, can't be restarted"
		suggestion.MarkSuggestionStatusSucceeded(reason, msg)
		logger.Info("Mark suggestion succeeded, can't be restarted")

		// If ResumePolicy = FromVolume, mark suggestion status succeeded, can be restarted
	} else if instance.Spec.ResumePolicy == experimentsv1beta1.FromVolume {
		msg := "Suggestion is succeeded, suggestion volume is not deleted, can be restarted"
		suggestion.MarkSuggestionStatusSucceeded(reason, msg)
		logger.Info("Mark suggestion succeeded, can be restarted")
	}

	if err := r.UpdateSuggestionStatus(suggestion); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileExperiment) restartSuggestion(instance *experimentsv1beta1.Experiment) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	original := &suggestionsv1beta1.Suggestion{}
	err := r.Get(context.TODO(),
		types.NamespacedName{Namespace: instance.GetNamespace(), Name: instance.GetName()}, original)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	// If Suggestion is restarting not needed to restart Suggestion
	if original.IsRestarting() {
		return nil
	}

	logger.Info("Suggestion is restarting, suggestion Running status is false")
	suggestion := original.DeepCopy()
	msg := "Suggestion is not running"
	// Mark suggestion status not running because experiment is restarting and suggestion deployment is not ready
	suggestion.MarkSuggestionStatusRunning(corev1.ConditionFalse, suggestionsv1beta1.SuggestionRestartReason, msg)

	if err := r.UpdateSuggestionStatus(suggestion); err != nil {
		return err
	}
	return nil
}
