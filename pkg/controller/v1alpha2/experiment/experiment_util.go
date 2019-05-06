package experiment

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	apiv1alpha2 "github.com/kubeflow/katib/pkg/api/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/experiment/util"
)

func (r *ReconcileExperiment) createTrialInstance(expInstance *experimentsv1alpha2.Experiment, trialInstance *apiv1alpha2.Trial) error {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: expInstance.GetName(), Namespace: expInstance.GetNamespace()})

	trial := &trialsv1alpha2.Trial{}
	trial.Name = fmt.Sprintf("%s-%s", expInstance.GetName(), utilrand.String(8))
	trial.Namespace = expInstance.GetNamespace()
	trial.Labels = map[string]string{"experiment": expInstance.GetName()}

	if err := controllerutil.SetControllerReference(expInstance, trial, r.scheme); err != nil {
		logger.Error(err, "Set controller reference error")
		return err
	}

	trialParams := util.TrialTemplateParams{
		Experiment: expInstance.GetName(),
		Trial:      trial.Name,
		NameSpace:  trial.Namespace,
	}
	if trialInstance.Spec != nil && trialInstance.Spec.ParameterAssignments != nil {
		for _, p := range trialInstance.Spec.ParameterAssignments.Assignments {
			trialParams.HyperParameters = append(trialParams.HyperParameters, p)
		}
	}

	runSpec, err := util.GetRunSpec(expInstance, trialParams)
	if err != nil {
		logger.Error(err, "Fail to get RunSpec from experiment", expInstance.Name)
		return err
	}

	trial.Spec.RunSpec = runSpec

	if err := r.Create(context.TODO(), trial); err != nil {
		logger.Error(err, "Trial create error", "Trial name", trial.Name)
		return err
	}
	return nil

}
