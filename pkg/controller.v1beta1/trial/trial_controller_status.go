package trial

import (
	"context"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
)

const (
	// For Trials
	TrialCreatedReason            = "TrialCreated"
	TrialRunningReason            = "TrialRunning"
	TrialSucceededReason          = "TrialSucceeded"
	TrialMetricsUnavailableReason = "MetricsUnavailable"
	TrialFailedReason             = "TrialFailed"

	// For Jobs
	JobCreatedReason            = "JobCreated"
	JobDeletedReason            = "JobDeleted"
	JobSucceededReason          = "JobSucceeded"
	JobMetricsUnavailableReason = "MetricsUnavailable"
	JobFailedReason             = "JobFailed"
	JobRunningReason            = "JobRunning"
)

type updateStatusFunc func(instance *trialsv1beta1.Trial) error

func (r *ReconcileTrial) updateStatus(instance *trialsv1beta1.Trial) error {
	err := r.Status().Update(context.TODO(), instance)
	if err != nil {
		return err
	}
	return nil
}
