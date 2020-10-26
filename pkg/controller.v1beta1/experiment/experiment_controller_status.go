package experiment

import (
	"context"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
)

type updateStatusFunc func(instance *experimentsv1beta1.Experiment) error

func (r *ReconcileExperiment) updateStatus(instance *experimentsv1beta1.Experiment) error {
	err := r.Status().Update(context.TODO(), instance)
	if err != nil {
		return err
	}
	return nil
}
