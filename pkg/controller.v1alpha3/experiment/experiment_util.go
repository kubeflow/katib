package experiment

import (
	"bytes"
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/util"
)

func (r *ReconcileExperiment) createTrialInstance(expInstance *experimentsv1alpha3.Experiment, trialAssignment *suggestionsv1alpha3.TrialAssignment) error {
	BUFSIZE := 1024
	logger := log.WithValues("Experiment", types.NamespacedName{Name: expInstance.GetName(), Namespace: expInstance.GetNamespace()})

	trial := &trialsv1alpha3.Trial{}
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
	runSpec, err := r.GetRunSpecWithHyperParameters(expInstance, expInstance.GetName(), trial.Name, trial.Namespace, hps)
	if err != nil {
		logger.Error(err, "Fail to get RunSpec from experiment", expInstance.Name)
		return err
	}

	trial.Spec.RunSpec = runSpec
	if expInstance.Spec.TrialTemplate != nil {
		trial.Spec.RetainRun = expInstance.Spec.TrialTemplate.Retain
	}

	buf := bytes.NewBufferString(runSpec)
	job := &unstructured.Unstructured{}
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, BUFSIZE).Decode(job); err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	if expInstance.Spec.MetricsCollectorSpec != nil {
		trial.Spec.MetricsCollector.Collector = expInstance.Spec.MetricsCollectorSpec.Collector
		trial.Spec.MetricsCollector.Source = expInstance.Spec.MetricsCollectorSpec.Source
	}

	if err := r.Create(context.TODO(), trial); err != nil {
		logger.Error(err, "Trial create error", "Trial name", trial.Name)
		return err
	}
	return nil

}
