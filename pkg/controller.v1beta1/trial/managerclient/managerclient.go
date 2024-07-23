/*
Copyright 2022 The Kubeflow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package managerclient

import (
	"time"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	common "github.com/kubeflow/katib/pkg/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

// ManagerClient is the interface for katib manager client in trial controller.
type ManagerClient interface {
	GetTrialObservationLog(
		instance *trialsv1beta1.Trial) (*api_pb.GetObservationLogReply, error)
	DeleteTrialObservationLog(
		instance *trialsv1beta1.Trial) (*api_pb.DeleteObservationLogReply, error)
	ReportTrialUnavailableMetrics(
		instance *trialsv1beta1.Trial) (*api_pb.ReportObservationLogReply, error)
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
		metricLogs = append(metricLogs, reply.ObservationLog.MetricLogs...)
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

func (d *DefaultClient) ReportTrialUnavailableMetrics(
	instance *trialsv1beta1.Trial) (*api_pb.ReportObservationLogReply, error) {
	request := &api_pb.ReportObservationLogRequest{
		TrialName: instance.Name,
		ObservationLog: &api_pb.ObservationLog{
			MetricLogs: []*api_pb.MetricLog{
				{
					TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
					Metric: &api_pb.Metric{
						Name:  instance.Spec.Objective.ObjectiveMetricName,
						Value: consts.UnavailableMetricValue,
					},
				},
			},
		},
	}
	reply, err := common.ReportObservationLog(request)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
