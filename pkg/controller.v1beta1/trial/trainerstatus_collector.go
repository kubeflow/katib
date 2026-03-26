package trial

import (
	"encoding/json"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

const (
	defaultTrainerStatusReportInterval = 5 * time.Second
	maxTrainerStatusMetrics            = 128
)

func (r *ReconcileTrial) UpdateTrialStatusTrainerStatus(instance *trialsv1beta1.Trial, deployedJob *unstructured.Unstructured) error {
	if deployedJob == nil {
		return nil
	}
	if instance.Spec.MetricsCollector.Collector.Kind != commonapiv1beta1.TrainerStatusCollector {
		return nil
	}

	gvk := schema.FromAPIVersionAndKind(deployedJob.GetAPIVersion(), deployedJob.GetKind())
	if gvk.Group != "trainer.kubeflow.org" || deployedJob.GetKind() != "TrainJob" {
		return nil
	}

	statusMap, found, err := unstructured.NestedMap(deployedJob.Object, "status", "trainerStatus")
	if err != nil {
		return fmt.Errorf("get TrainJob.status.trainerStatus: %w", err)
	}
	if !found {
		return nil
	}

	payloadBytes, err := json.Marshal(statusMap)
	if err != nil {
		return fmt.Errorf("marshal TrainJob.status.trainerStatus: %w", err)
	}

	var incoming trialsv1beta1.TrainerStatus
	if err := json.Unmarshal(payloadBytes, &incoming); err != nil {
		return fmt.Errorf("unmarshal TrainJob.status.trainerStatus: %w", err)
	}

	if incoming.LastUpdatedTime == nil {
		// Trainer requires lastUpdatedTime when trainerStatus is present; if missing, ignore.
		return nil
	}

	if instance.Status.TrainerStatus != nil && instance.Status.TrainerStatus.LastUpdatedTime != nil {
		if instance.Status.TrainerStatus.LastUpdatedTime.Time.Equal(incoming.LastUpdatedTime.Time) {
			return nil
		}
	}

	if instance.Status.TrainerStatusLastReportedTime != nil {
		if time.Since(instance.Status.TrainerStatusLastReportedTime.Time) < defaultTrainerStatusReportInterval {
			return nil
		}
	}

	// Report TrainerStatus metrics to Katib DB so Katib can build metric history.
	if len(incoming.Metrics) > 0 {
		metrics := incoming.Metrics
		if len(metrics) > maxTrainerStatusMetrics {
			metrics = metrics[:maxTrainerStatusMetrics]
		}

		timestamp := incoming.LastUpdatedTime.Time.UTC().Format(time.RFC3339Nano)
		metricLogs := make([]*api_pb.MetricLog, 0, len(metrics))
		for _, m := range metrics {
			if m.Name == "" {
				continue
			}
			metricLogs = append(metricLogs, &api_pb.MetricLog{
				TimeStamp: timestamp,
				Metric: &api_pb.Metric{
					Name:  m.Name,
					Value: m.Value,
				},
			})
		}

		if len(metricLogs) > 0 {
			observationLog := &api_pb.ObservationLog{MetricLogs: metricLogs}
			if _, err := r.ReportTrialObservationLog(instance, observationLog); err != nil {
				return fmt.Errorf("%w: %w", errReportMetricsFailed, err)
			}
		}
	}

	instance.Status.TrainerStatus = &incoming
	now := metav1.Now()
	instance.Status.TrainerStatusLastReportedTime = &now

	return nil
}
