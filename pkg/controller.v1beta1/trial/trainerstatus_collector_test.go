package trial

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	managerclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/trial/managerclient"
)

type observationLogMatcher struct {
	expectedTimestamp string
	expectedMetrics   map[string]string
}

func (m observationLogMatcher) Matches(x any) bool {
	obs, ok := x.(*api_pb.ObservationLog)
	if !ok || obs == nil {
		return false
	}
	if len(obs.MetricLogs) != len(m.expectedMetrics) {
		return false
	}
	seen := map[string]string{}
	for _, ml := range obs.MetricLogs {
		if ml == nil || ml.Metric == nil {
			return false
		}
		if ml.TimeStamp != m.expectedTimestamp {
			return false
		}
		seen[ml.Metric.Name] = ml.Metric.Value
	}
	for k, v := range m.expectedMetrics {
		if seen[k] != v {
			return false
		}
	}
	return true
}

func (m observationLogMatcher) String() string {
	return "matches expected ObservationLog"
}

func TestUpdateTrialStatusTrainerStatus(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockManagerClient := managerclientmock.NewMockManagerClient(mockCtrl)

	r := &ReconcileTrial{ManagerClient: mockManagerClient}

	trial := &trialsv1beta1.Trial{}
	trial.Name = "trial-1"
	trial.Namespace = "default"
	trial.Spec.Objective = &commonv1beta1.ObjectiveSpec{ObjectiveMetricName: "loss"}
	trial.Spec.MetricsCollector = commonv1beta1.MetricsCollectorSpec{Collector: &commonv1beta1.CollectorSpec{Kind: commonv1beta1.TrainerStatusCollector}}

	trainJob := &unstructured.Unstructured{}
	trainJob.SetAPIVersion("trainer.kubeflow.org/v1alpha1")
	trainJob.SetKind("TrainJob")
	trainJob.SetName("trainjob-1")
	trainJob.SetNamespace("default")
	trainJob.Object = map[string]any{
		"apiVersion": "trainer.kubeflow.org/v1alpha1",
		"kind":       "TrainJob",
		"metadata": map[string]any{
			"name":      "trainjob-1",
			"namespace": "default",
		},
		"status": map[string]any{
			"trainerStatus": map[string]any{
				"progressPercentage":        int64(42),
				"estimatedRemainingSeconds": int64(120),
				"metrics": []any{
					map[string]any{"name": "loss", "value": "0.123"},
					map[string]any{"name": "accuracy", "value": "0.99"},
				},
				"lastUpdatedTime": "2024-01-02T03:04:05Z",
			},
		},
	}

	mockManagerClient.EXPECT().ReportTrialObservationLog(gomock.Any(), observationLogMatcher{
		expectedTimestamp: "2024-01-02T03:04:05Z",
		expectedMetrics: map[string]string{
			"loss":     "0.123",
			"accuracy": "0.99",
		},
	}).Return(&api_pb.ReportObservationLogReply{}, nil).Times(1)

	err := r.UpdateTrialStatusTrainerStatus(trial, trainJob)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(trial.Status.TrainerStatus).NotTo(gomega.BeNil())
	g.Expect(trial.Status.TrainerStatus.ProgressPercentage).NotTo(gomega.BeNil())
	g.Expect(*trial.Status.TrainerStatus.ProgressPercentage).To(gomega.Equal(int32(42)))
	g.Expect(trial.Status.TrainerStatus.EstimatedRemainingSeconds).NotTo(gomega.BeNil())
	g.Expect(*trial.Status.TrainerStatus.EstimatedRemainingSeconds).To(gomega.Equal(int32(120)))
	g.Expect(trial.Status.TrainerStatus.LastUpdatedTime).NotTo(gomega.BeNil())
	expectedTime, err := time.Parse(time.RFC3339, "2024-01-02T03:04:05Z")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(trial.Status.TrainerStatus.LastUpdatedTime.Time).To(gomega.Equal(expectedTime))
	g.Expect(trial.Status.TrainerStatusLastReportedTime).NotTo(gomega.BeNil())

	// Dedupe: same lastUpdatedTime should not trigger another report.
	err = r.UpdateTrialStatusTrainerStatus(trial, trainJob)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
