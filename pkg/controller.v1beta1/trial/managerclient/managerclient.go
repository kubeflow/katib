package managerclient

import (
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	common "github.com/kubeflow/katib/pkg/common/v1beta1"
)

// ManagerClient is the interface for katib manager client in trial controller.
type ManagerClient interface {
	GetTrialObservationLog(
		instance *trialsv1beta1.Trial) (*api_pb.GetObservationLogReply, error)
	DeleteTrialObservationLog(
		instance *trialsv1beta1.Trial) (*api_pb.DeleteObservationLogReply, error)
}

// DefaultClient implements the Client interface.
type DefaultClient struct {
}

// New creates a new ManagerClient.
func New() ManagerClient {
	return &DefaultClient{}
}

func (d *DefaultClient) GetTrialObservationLog(
	instance *trialsv1beta1.Trial) (*api_pb.GetObservationLogReply, error) {
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
	if reply.ObservationLog == nil {
		reply.ObservationLog = &api_pb.ObservationLog{}
	}
	if reply.ObservationLog.MetricLogs == nil {
		reply.ObservationLog.MetricLogs = []*api_pb.MetricLog{}
	}
	// fetch additional metrics if exists
	metricLogs := reply.ObservationLog.MetricLogs
	for _, metricName := range instance.Spec.Objective.AdditionalMetricNames {
		request := &api_pb.GetObservationLogRequest{
			TrialName: instance.Name, MetricName: metricName,
		}
		reply, err := common.GetObservationLog(request)
		if err != nil {
			return nil, err
		}
		if reply.ObservationLog == nil || reply.ObservationLog.MetricLogs == nil {
			continue
		}
		for _, log := range reply.ObservationLog.MetricLogs {
			metricLogs = append(metricLogs, log)
		}
	}
	reply.ObservationLog.MetricLogs = metricLogs

	return reply, nil
}

func (d *DefaultClient) DeleteTrialObservationLog(
	instance *trialsv1beta1.Trial) (*api_pb.DeleteObservationLogReply, error) {
	request := &api_pb.DeleteObservationLogRequest{
		TrialName: instance.Name,
	}
	reply, err := common.DeleteObservationLog(request)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
