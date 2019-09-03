package trial

import (
	"context"

	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
)

type updateStatusFunc func(instance *trialsv1alpha3.Trial) error

func (r *ReconcileTrial) updateStatus(instance *trialsv1alpha3.Trial) error {
	err := r.Status().Update(context.TODO(), instance)
	if err != nil {
		return err
	}
	return nil
}
