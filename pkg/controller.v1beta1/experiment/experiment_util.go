package experiment

import (
	"context"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	suggestionController "github.com/kubeflow/katib/pkg/controller.v1beta1/suggestion"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
)

const (
	updatePrometheusMetrics = "update-prometheus-metrics"
)

func (r *ReconcileExperiment) createTrialInstance(expInstance *experimentsv1beta1.Experiment, trialAssignment *suggestionsv1beta1.TrialAssignment) error {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: expInstance.GetName(), Namespace: expInstance.GetNamespace()})

	trial := &trialsv1beta1.Trial{}
	trial.Name = trialAssignment.Name
	trial.Namespace = expInstance.GetNamespace()
	trial.Labels = util.TrialLabels(expInstance)

	if err := controllerutil.SetControllerReference(expInstance, trial, r.scheme); err != nil {
		logger.Error(err, "Set controller reference error")
		return err
	}

	trial.Spec.Objective = expInstance.Spec.Objective

	hps := trialAssignment.ParameterAssignments
	trial.Spec.ParameterAssignments = trialAssignment.ParameterAssignments

	runSpec, err := r.GetRunSpecWithHyperParameters(expInstance, trial.Name, trial.Namespace, hps, buildTrialMetaForRunSpec(trial))
	if err != nil {
		logger.Error(err, "Fail to get RunSpec from experiment", expInstance.Name)
		return err
	}

	trial.Spec.RunSpec = runSpec
	if expInstance.Spec.TrialTemplate != nil {
		trial.Spec.RetainRun = expInstance.Spec.TrialTemplate.Retain
	}

	if expInstance.Spec.MetricsCollectorSpec != nil {
		trial.Spec.MetricsCollector = *expInstance.Spec.MetricsCollectorSpec
	}

	if err := r.Create(context.TODO(), trial); err != nil {
		logger.Error(err, "Trial create error", "Trial name", trial.Name)
		return err
	}
	return nil

}

func buildTrialMetaForRunSpec(trial *trialsv1beta1.Trial) map[string]string {
	return map[string]string{
		consts.TrialTemplateMetaKeyOfName:       trial.Name,
		consts.TrialTemplateMetaKeyOfNamespace:  trial.Namespace,
		consts.TrialTemplateMetaKeyOfKind:       trial.Kind,
		consts.TrialTemplateMetaKeyOfAPIVersion: trial.APIVersion,
	}
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

func (r *ReconcileExperiment) terminateSuggestion(instance *experimentsv1beta1.Experiment) error {
	original := &suggestionsv1beta1.Suggestion{}
	err := r.Get(context.TODO(),
		types.NamespacedName{Namespace: instance.GetNamespace(), Name: instance.GetName()}, original)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	// If Suggestion is failed or Suggestion is Succeeded, not needed to terminate Suggestion
	if original.IsFailed() || original.IsSucceeded() {
		return nil
	}
	log.Info("Start terminating suggestion")
	suggestion := original.DeepCopy()
	msg := "Suggestion is succeeded"
	suggestion.MarkSuggestionStatusSucceeded(suggestionController.SuggestionSucceededReason, msg)
	log.Info("Mark suggestion succeeded")

	if err := r.UpdateSuggestionStatus(suggestion); err != nil {
		return err
	}
	return nil
}
