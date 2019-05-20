package trial

import (
	"context"

	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
)

type updateStatusFunc func(instance *trialsv1alpha2.Trial) error

func (r *ReconcileTrial) updateStatus(instance *trialsv1alpha2.Trial) error {
	err := r.Status().Update(context.TODO(), instance)
	if err != nil {
		return err
	}
	return nil
}
