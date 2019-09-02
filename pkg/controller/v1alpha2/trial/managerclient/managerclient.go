package managerclient

import (
	"fmt"

	commonv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/common/v1alpha2"
	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	common "github.com/kubeflow/katib/pkg/common/v1alpha2"
)

// ManagerClient is the interface for katib manager client in trial controller.
type ManagerClient interface {
	CreateTrialInDB(instance *trialsv1alpha2.Trial) error
	UpdateTrialStatusInDB(instance *trialsv1alpha2.Trial) error
	GetTrialObservation(instance *trialsv1alpha2.Trial) error
	GetTrialObservationLog(
		instance *trialsv1alpha2.Trial) (*api_pb.GetObservationLogReply, error)
	GetTrialConf(instance *trialsv1alpha2.Trial) *api_pb.Trial
}

// DefaultClient implements the Client interface.
type DefaultClient struct {
}

// New creates a new ManagerClient.
func New() ManagerClient {
	return &DefaultClient{}
}

func (d *DefaultClient) CreateTrialInDB(instance *trialsv1alpha2.Trial) error {
	trial := d.GetTrialConf(instance)
	request := &api_pb.RegisterTrialRequest{
		Trial: trial,
	}
	if _, err := common.RegisterTrial(request); err != nil {
		return err
	}
	return nil
}

func (d *DefaultClient) UpdateTrialStatusInDB(instance *trialsv1alpha2.Trial) error {
	newStatus := &api_pb.TrialStatus{
		StartTime:      common.ConvertTime2RFC3339(instance.Status.StartTime),
		CompletionTime: common.ConvertTime2RFC3339(instance.Status.CompletionTime),
		Condition:      getCondition(instance),
	}
	if instance.Status.Observation != nil {
		observation := &api_pb.Observation{
			Metrics: []*api_pb.Metric{},
		}
		for _, m := range instance.Status.Observation.Metrics {
			metric := &api_pb.Metric{
				Name:  m.Name,
				Value: fmt.Sprintf("%f", m.Value),
			}
			observation.Metrics = append(observation.Metrics, metric)
		}
		newStatus.Observation = observation
	}
	request := &api_pb.UpdateTrialStatusRequest{
		NewStatus: newStatus,
		TrialName: instance.Name,
	}
	if _, err := common.UpdateTrialStatus(request); err != nil {
		return err
	}
	return nil
}

func (d *DefaultClient) GetTrialObservationLog(
	instance *trialsv1alpha2.Trial) (*api_pb.GetObservationLogReply, error) {
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

func (d *DefaultClient) GetTrialObservation(instance *trialsv1alpha2.Trial) error {
	return nil
}

func (d *DefaultClient) GetTrialConf(instance *trialsv1alpha2.Trial) *api_pb.Trial {
	trial := &api_pb.Trial{
		Spec: &api_pb.TrialSpec{
			Objective: &api_pb.ObjectiveSpec{
				AdditionalMetricNames: []string{},
			},
			ParameterAssignments: &api_pb.TrialSpec_ParameterAssignments{
				Assignments: []*api_pb.ParameterAssignment{},
			},
		},
		Status: &api_pb.TrialStatus{
			StartTime:      common.ConvertTime2RFC3339(instance.Status.StartTime),
			CompletionTime: common.ConvertTime2RFC3339(instance.Status.CompletionTime),
			Condition:      getCondition(instance),
		},
	}
	trial.Name = instance.Name

	trial.Spec.ExperimentName = instance.Labels["experiment"]

	//Populate Objective
	switch instance.Spec.Objective.Type {
	case commonv1alpha2.ObjectiveTypeMaximize:
		trial.Spec.Objective.Type = api_pb.ObjectiveType_MAXIMIZE
	case commonv1alpha2.ObjectiveTypeMinimize:
		trial.Spec.Objective.Type = api_pb.ObjectiveType_MINIMIZE
	default:
		trial.Spec.Objective.Type = api_pb.ObjectiveType_UNKNOWN

	}
	trial.Spec.Objective.Goal = float64(*instance.Spec.Objective.Goal)
	trial.Spec.Objective.ObjectiveMetricName = instance.Spec.Objective.ObjectiveMetricName
	for _, m := range instance.Spec.Objective.AdditionalMetricNames {
		trial.Spec.Objective.AdditionalMetricNames = append(trial.Spec.Objective.AdditionalMetricNames, m)
	}

	//Populate Parameter Assignments
	for _, p := range instance.Spec.ParameterAssignments {
		trial.Spec.ParameterAssignments.Assignments = append(
			trial.Spec.ParameterAssignments.Assignments,
			&api_pb.ParameterAssignment{
				Name:  p.Name,
				Value: p.Value,
			})
	}

	trial.Spec.RunSpec = instance.Spec.RunSpec

	trial.Spec.MetricsCollectorSpec = instance.Spec.MetricsCollectorSpec

	return trial
}

func getCondition(inst *trialsv1alpha2.Trial) api_pb.TrialStatus_TrialConditionType {
	condition, _ := inst.GetLastConditionType()
	switch condition {
	case trialsv1alpha2.TrialCreated:
		return api_pb.TrialStatus_CREATED
	case trialsv1alpha2.TrialRunning:
		return api_pb.TrialStatus_RUNNING
	case trialsv1alpha2.TrialKilled:
		return api_pb.TrialStatus_KILLED
	case trialsv1alpha2.TrialSucceeded:
		return api_pb.TrialStatus_SUCCEEDED
	case trialsv1alpha2.TrialFailed:
		return api_pb.TrialStatus_FAILED
	default:
		return api_pb.TrialStatus_UNKNOWN
	}
}
