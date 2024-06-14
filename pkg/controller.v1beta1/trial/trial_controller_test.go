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

package trial

import (
	"sync"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.uber.org/mock/gomock"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	trialutil "github.com/kubeflow/katib/pkg/controller.v1beta1/trial/util"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	managerclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/trial/managerclient"
)

const (
	namespace       = "default"
	trialName       = "test-trial"
	batchJobName    = "test-job"
	objectiveMetric = "accuracy"
	timeout         = time.Second * 80
)

var trialKey = types.NamespacedName{Name: trialName, Namespace: namespace}
var batchJobKey = types.NamespacedName{Name: batchJobName, Namespace: namespace}

func init() {
	logf.SetLogger(zap.New(zap.UseDevMode(true)))
}

func TestAdd(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	mgr, err := manager.New(cfg, manager.Options{Metrics: metricsserver.Options{BindAddress: "0"}})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// Set Trial resources.
	trialResources := []schema.GroupVersionKind{
		{
			Group:   "kubeflow.org",
			Version: "v1",
			Kind:    "TFJob",
		},
		{
			Group:   "kubeflow.org",
			Version: "v1",
			Kind:    "MPIJob",
		},
	}

	viper.Set(consts.ConfigTrialResources, trialResources)

	// Test - Try to add Trial controller to the manager
	g.Expect(Add(mgr)).NotTo(gomega.HaveOccurred())
}

func TestReconcileBatchJob(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockManagerClient := managerclientmock.NewMockManagerClient(mockCtrl)

	// Setup the Manager and Controller. Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{Metrics: metricsserver.Options{BindAddress: "0"}})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	r := &ReconcileTrial{
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		ManagerClient: mockManagerClient,
		recorder:      mgr.GetEventRecorderFor(ControllerName),
		collector:     trialutil.NewTrialsCollector(mgr.GetCache(), prometheus.NewRegistry()),
	}

	r.updateStatusHandler = func(instance *trialsv1beta1.Trial) error {
		var err error = errors.NewBadRequest("fake-error")
		// Try to update status until it be succeeded
		for err != nil {
			updatedInstance := &trialsv1beta1.Trial{}
			if err = c.Get(ctx, trialKey, updatedInstance); err != nil {
				continue
			}
			updatedInstance.Status = instance.Status
			err = r.updateStatus(updatedInstance)
		}
		return err
	}

	recFn := SetupTestReconcile(r)
	// Set Job resource
	trialResources := []schema.GroupVersionKind{
		{
			Group:   "batch",
			Version: "v1",
			Kind:    "Job",
		},
	}

	viper.Set(consts.ConfigTrialResources, trialResources)
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	// Start test manager.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(ctx)).NotTo(gomega.HaveOccurred())
	}()

	// Result for GetTrialObservationLog with some metrics.
	observationLogAvailable := &api_pb.GetObservationLogReply{
		ObservationLog: &api_pb.ObservationLog{
			MetricLogs: []*api_pb.MetricLog{
				{
					TimeStamp: "2020-08-10T14:47:38+08:00",
					Metric: &api_pb.Metric{
						Name:  objectiveMetric,
						Value: "0.99",
					},
				},
				{
					TimeStamp: "2020-08-10T14:50:38+08:00",
					Metric: &api_pb.Metric{
						Name:  objectiveMetric,
						Value: "0.11",
					},
				},
			},
		},
	}
	// Empty result for GetTrialObservationLog.
	// If objective metrics are not parsed, metrics collector reports "unavailable" value to DB.
	observationLogUnavailable := &api_pb.GetObservationLogReply{
		ObservationLog: &api_pb.ObservationLog{
			MetricLogs: []*api_pb.MetricLog{
				{
					Metric: &api_pb.Metric{
						Name:  objectiveMetric,
						Value: consts.UnavailableMetricValue,
					},
					TimeStamp: time.Time{}.UTC().Format(time.RFC3339),
				},
			},
		},
	}

	t.Run(`Trial run with "Failed" BatchJob.`, func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		mockManagerClient.EXPECT().DeleteTrialObservationLog(gomock.Any()).Return(nil, nil)

		trial := newFakeTrialBatchJob()
		batchJob := &batchv1.Job{}

		// Create the Trial
		g.Expect(c.Create(ctx, trial)).NotTo(gomega.HaveOccurred())

		// Expect that BatchJob with appropriate name is created
		g.Eventually(func() error {
			return c.Get(ctx, batchJobKey, batchJob)
		}, timeout).Should(gomega.Succeed())

		// Expect that Trial status is running
		g.Eventually(func() bool {
			if err = c.Get(ctx, trialKey, trial); err != nil {
				return false
			}
			return trial.IsRunning()
		}, timeout).Should(gomega.BeTrue())

		// Manually update BatchJob status to failed
		// Expect that Trial status is failed
		g.Eventually(func() bool {
			if err = c.Get(ctx, batchJobKey, batchJob); err != nil {
				return false
			}
			batchJob.Status = batchv1.JobStatus{
				Conditions: []batchv1.JobCondition{
					{
						Type:    batchv1.JobFailed,
						Status:  corev1.ConditionTrue,
						Message: "BatchJob failed test message",
						Reason:  "BatchJob failed test reason",
					},
				},
			}
			if err = c.Status().Update(ctx, batchJob); err != nil {
				return false
			}

			if err = c.Get(ctx, trialKey, trial); err != nil {
				return false
			}
			return trial.IsFailed()
		}, timeout).Should(gomega.BeTrue())

		// Delete the Trial
		g.Expect(c.Delete(ctx, trial)).NotTo(gomega.HaveOccurred())

		// Expect that Trial is deleted
		// BatchJob can't be deleted because GC doesn't work in envtest and BatchJob stuck in termination phase.
		// Ref: https://book.kubebuilder.io/reference/testing/envtest.html#testing-considerations.
		g.Eventually(func() bool {
			return errors.IsNotFound(c.Get(ctx, trialKey, &trialsv1beta1.Trial{}))
		}, timeout).Should(gomega.BeTrue())
	})

	t.Run(`Trail with "Complete" BatchJob and Available metrics.`, func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		gomock.InOrder(
			mockManagerClient.EXPECT().GetTrialObservationLog(gomock.Any()).Return(observationLogAvailable, nil).MinTimes(1),
			mockManagerClient.EXPECT().DeleteTrialObservationLog(gomock.Any()).Return(nil, nil),
		)
		batchJob := &batchv1.Job{}
		batchJobCompleteMessage := "BatchJob completed test message"
		batchJobCompleteReason := "BatchJob completed test reason"
		// Update BatchJob status to Complete.
		g.Expect(c.Get(ctx, batchJobKey, batchJob)).NotTo(gomega.HaveOccurred())
		batchJob.Status = batchv1.JobStatus{
			Conditions: []batchv1.JobCondition{
				{
					Type:    batchv1.JobComplete,
					Status:  corev1.ConditionTrue,
					Message: batchJobCompleteMessage,
					Reason:  batchJobCompleteReason,
				},
			},
		}
		g.Expect(c.Status().Update(ctx, batchJob)).NotTo(gomega.HaveOccurred())

		// Create the Trial
		trial := newFakeTrialBatchJob()
		g.Expect(c.Create(ctx, trial)).NotTo(gomega.HaveOccurred())

		// Expect that Trial status is succeeded and metrics are properly populated
		// Metrics available because GetTrialObservationLog returns values
		start := time.Now()
		g.Eventually(func() bool {
			if err = c.Get(ctx, trialKey, trial); err != nil {
				t.Log(time.Since(start), err)
				return false
			}
			return trial.IsSucceeded() &&
				len(trial.Status.Observation.Metrics) > 0 &&
				trial.Status.Observation.Metrics[0].Min == "0.11" &&
				trial.Status.Observation.Metrics[0].Max == "0.99" &&
				trial.Status.Observation.Metrics[0].Latest == "0.11"
		}, timeout).Should(gomega.BeTrue())

		// Delete the Trial
		g.Expect(c.Delete(ctx, trial)).NotTo(gomega.HaveOccurred())

		// Expect that Trial is deleted
		g.Eventually(func() bool {
			return errors.IsNotFound(c.Get(ctx, trialKey, &trialsv1beta1.Trial{}))
		}, timeout).Should(gomega.BeTrue())
	})

	t.Run(`Trail with "Complete" BatchJob and Unavailable metrics.`, func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		gomock.InOrder(
			mockManagerClient.EXPECT().GetTrialObservationLog(gomock.Any()).Return(observationLogUnavailable, nil).MinTimes(1),
			mockManagerClient.EXPECT().DeleteTrialObservationLog(gomock.Any()).Return(nil, nil),
		)
		// Create the Trial
		trial := newFakeTrialBatchJob()
		g.Expect(c.Create(ctx, trial)).NotTo(gomega.HaveOccurred())

		// Expect that Trial status is succeeded with "false" status and "metrics unavailable" reason.
		// Metrics unavailable because GetTrialObservationLog returns "unavailable".
		g.Eventually(func() bool {
			if err = c.Get(ctx, trialKey, trial); err != nil {
				return false
			}
			return trial.IsMetricsUnavailable() &&
				len(trial.Status.Observation.Metrics) > 0 &&
				trial.Status.Observation.Metrics[0].Min == consts.UnavailableMetricValue &&
				trial.Status.Observation.Metrics[0].Max == consts.UnavailableMetricValue &&
				trial.Status.Observation.Metrics[0].Latest == consts.UnavailableMetricValue
		}, timeout).Should(gomega.BeTrue())

		// Delete the Trial
		g.Expect(c.Delete(ctx, trial)).NotTo(gomega.HaveOccurred())

		// Expect that Trial is deleted
		g.Eventually(func() bool {
			return errors.IsNotFound(c.Get(ctx, trialKey, &trialsv1beta1.Trial{}))
		}, timeout).Should(gomega.BeTrue())
	})

	t.Run("Update status for empty Trial", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		g.Expect(r.updateStatus(&trialsv1beta1.Trial{})).To(gomega.HaveOccurred())
	})
}

func TestGetObjectiveMetricValue(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	metricLogs := []*api_pb.MetricLog{
		{TimeStamp: "2020-04-13T14:47:38+08:00", Metric: &api_pb.Metric{Name: "error", Value: "0.03"}},
		{TimeStamp: "2020-04-13T14:47:39+08:00", Metric: &api_pb.Metric{Name: "error", Value: "0.02"}},
		{TimeStamp: "2020-04-13T14:47:40+08:00", Metric: &api_pb.Metric{Name: "error", Value: "0.01"}},
		{TimeStamp: "2020-04-13T14:47:41+08:00", Metric: &api_pb.Metric{Name: "error", Value: "0.05"}},
		{TimeStamp: "2020-04-13T14:47:41+08:00", Metric: &api_pb.Metric{Name: "error", Value: "0.06"}},
		{TimeStamp: "2020-04-13T14:47:41+08:00", Metric: &api_pb.Metric{Name: "error", Value: "0.07"}},
		{TimeStamp: "2020-04-12T14:47:42+08:00", Metric: &api_pb.Metric{Name: "error", Value: "0.1"}},
		{TimeStamp: "2020-04-13T14:47:38+08:00", Metric: &api_pb.Metric{Name: objectiveMetric, Value: "0.7"}},
		{TimeStamp: "2020-04-13T14:47:39+08:00", Metric: &api_pb.Metric{Name: objectiveMetric, Value: "0.71"}},
		{TimeStamp: "2020-04-13T14:47:40+08:00", Metric: &api_pb.Metric{Name: objectiveMetric, Value: "0.72"}},
		{TimeStamp: "2020-04-13T14:47:41+08:00", Metric: &api_pb.Metric{Name: objectiveMetric, Value: "0.68"}},
		{TimeStamp: "2020-04-13T14:47:41+08:00", Metric: &api_pb.Metric{Name: objectiveMetric, Value: "0.69"}},
		{TimeStamp: "2020-04-13T14:47:41+08:00", Metric: &api_pb.Metric{Name: objectiveMetric, Value: "0.67"}},
		{TimeStamp: "2020-04-12T14:47:42+08:00", Metric: &api_pb.Metric{Name: objectiveMetric, Value: "0.6"}},
	}

	getMetricsFromLogs := func(strategies []commonv1beta1.MetricStrategy) (*commonv1beta1.Metric, *commonv1beta1.Metric, error) {
		observation, err := getMetrics(metricLogs, strategies)
		if err != nil {
			return nil, nil, err
		}
		var errMetric, accMetric *commonv1beta1.Metric
		for index, metric := range observation.Metrics {
			if metric.Name == "error" {
				errMetric = &observation.Metrics[index]
			} else if metric.Name == objectiveMetric {
				accMetric = &observation.Metrics[index]
			}
		}
		return errMetric, accMetric, nil
	}

	metricStrategies := []commonv1beta1.MetricStrategy{
		{Name: "error", Value: commonv1beta1.ExtractByMin},
		{Name: objectiveMetric, Value: commonv1beta1.ExtractByMax},
	}
	errMetric, accMetric, err := getMetricsFromLogs(metricStrategies)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(errMetric.Latest).To(gomega.Equal("0.07"))
	g.Expect(errMetric.Max).To(gomega.Equal("0.1"))
	g.Expect(errMetric.Min).To(gomega.Equal("0.01"))
	g.Expect(accMetric.Latest).To(gomega.Equal("0.67"))
	g.Expect(accMetric.Max).To(gomega.Equal("0.72"))
	g.Expect(accMetric.Min).To(gomega.Equal("0.6"))

	invalidLogs := []*api_pb.MetricLog{
		// Add one other metric to test correct parsing
		{TimeStamp: "2020-08-10T14:47:42+08:00", Metric: &api_pb.Metric{Name: "not-accuracy", Value: "1.15"}},
		// Add metric with invalid timestamp
		{TimeStamp: "2020-08-10T14:47:42", Metric: &api_pb.Metric{Name: objectiveMetric, Value: "0.77"}},
	}
	_, err = getMetrics(invalidLogs, metricStrategies)
	g.Expect(err).To(gomega.HaveOccurred())
}

func newFakeTrialBatchJob() *trialsv1beta1.Trial {
	primaryContainer := "training-container"

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      batchJobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  primaryContainer,
							Image: "docker.io/kubeflowkatib/pytorch-mnist-cpu",
							Command: []string{
								"python3",
								"/opt/pytorch-mnist/mnist.py",
								"--epochs=1",
								"--batch-size=16",
								"--lr=0.01",
								"--momentum=0.9",
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
	runSpec, _ := util.ConvertObjectToUnstructured(job)

	return &trialsv1beta1.Trial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trialName,
			Namespace: namespace,
		},
		Spec: trialsv1beta1.TrialSpec{
			PrimaryContainerName: primaryContainer,
			SuccessCondition:     experimentsv1beta1.DefaultJobSuccessCondition,
			FailureCondition:     experimentsv1beta1.DefaultJobFailureCondition,
			Objective: &commonv1beta1.ObjectiveSpec{
				ObjectiveMetricName: objectiveMetric,
				MetricStrategies: []commonv1beta1.MetricStrategy{
					{
						Name:  objectiveMetric,
						Value: commonv1beta1.ExtractByMax,
					},
				},
			},
			RunSpec: runSpec,
		},
	}
}
