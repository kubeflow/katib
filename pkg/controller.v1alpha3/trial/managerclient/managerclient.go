package managerclient

import (
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	common "github.com/kubeflow/katib/pkg/common/v1alpha3"
)

// ManagerClient is the interface for katib manager client in trial controller.
type ManagerClient interface {
	GetTrialObservationLog(
		instance *trialsv1alpha3.Trial) (*api_pb.GetObservationLogReply, error)
}

// DefaultClient implements the Client interface.
type DefaultClient struct {
}

// New creates a new ManagerClient.
func New() ManagerClient {
	return &DefaultClient{}
}

func (d *DefaultClient) GetTrialObservationLog(
	instance *trialsv1alpha3.Trial) (*api_pb.GetObservationLogReply, error) {
	// read GetObservationLog call and update observation field
	objectiveMetricName := instance.Spec.Objective.ObjectiveMetricName
	request := &api_pb.GetObservationLogRequest{
		TrialName:  instance.Name,
		MetricName: objectiveMetricName,
	}
	reply, err := common.GetObservationLog(request)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
