package experiment

import (
	"context"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
)

type updateStatusFunc func(instance *experimentsv1alpha2.Experiment) error

func (r *ReconcileExperiment) updateStatus(instance *experimentsv1alpha2.Experiment) error {
	err := r.Status().Update(context.TODO(), instance)
	if err != nil {
		return err
	}
	return nil
}
