package trial

import (
	"context"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
)

type updateStatusFunc func(instance *trialsv1beta1.Trial) error

func (r *ReconcileTrial) updateStatus(instance *trialsv1beta1.Trial) error {
	err := r.Status().Update(context.TODO(), instance)
	if err != nil {
		return err
	}
	return nil
}
