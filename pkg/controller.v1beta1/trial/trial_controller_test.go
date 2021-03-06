package trial

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	commonv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	"github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	trialutil "github.com/kubeflow/katib/pkg/controller.v1beta1/trial/util"
	util "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	managerclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/trial/managerclient"
)

const (
	namespace       = "default"
	trialName       = "test-trial"
	tfJobName       = "test-tfjob"
	batchJobName    = "test-job"
	objectiveMetric = "accuracy"
	timeout         = time.Second * 40
)

var trialKey = types.NamespacedName{Name: trialName, Namespace: namespace}
var tfJobKey = types.NamespacedName{Name: tfJobName, Namespace: namespace}
var batchJobKey = types.NamespacedName{Name: batchJobName, Namespace: namespace}

func init() {
	logf.SetLogger(zap.New())
}

func TestAdd(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// Set Trial resources.
	// TFJob controller is installed, MPIJob controller is missed.
	trialResources := trialutil.GvkListFlag{
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

func TestReconcileTFJob(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockManagerClient := managerclientmock.NewMockManagerClient(mockCtrl)

	// Setup the Manager and Controller. Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
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
			c.Get(context.TODO(), trialKey, updatedInstance)
			updatedInstance.Status = instance.Status
			err = r.updateStatus(updatedInstance)
		}
		return err
	}

	recFn := SetupTestReconcile(r)

	// Set TFJob resource
	trialResources := trialutil.GvkListFlag{
		{
			Group:   "kubeflow.org",
			Version: "v1",
			Kind:    "TFJob",
		},
	}

	viper.Set(consts.ConfigTrialResources, trialResources)

	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	// Start test manager.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(context.TODO())).NotTo(gomega.HaveOccurred())
	}()

	// Empty result for GetTrialObservationLog.
	// If objective metrics are not parsed, metrics collector reports "unavailable" value to DB.
	observationLog := &api_pb.GetObservationLogReply{
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

	mockManagerClient.EXPECT().GetTrialObservationLog(gomock.Any()).Return(observationLog, nil).AnyTimes()
	mockManagerClient.EXPECT().DeleteTrialObservationLog(gomock.Any()).Return(nil, nil).AnyTimes()

	// Test - Regural Trial run with TFJob
	trial := newFakeTrialTFJob()
	tfJob := &tfv1.TFJob{}

	// Create the Trial
	g.Expect(c.Create(context.TODO(), trial)).NotTo(gomega.HaveOccurred())

	// Expect that TFJob with appropriate name is created
	g.Eventually(func() error {
		return c.Get(context.TODO(), tfJobKey, tfJob)
	}, timeout).Should(gomega.Succeed())

	// Expect that Trial status is running
	g.Eventually(func() bool {
		c.Get(context.TODO(), trialKey, trial)
		return trial.IsRunning()
	}, timeout).Should(gomega.BeTrue())

	// Manually update TFJob status to succeeded
	// Expect that Trial succeeded status is false with metrics unavailable reason
	// Metrics unavailable because GetTrialObservationLog returns nil
	SucceededReason := "TFJob succeeded test reason"
	SucceededMessage := "TFJob succeeded test message"
	g.Eventually(func() bool {
		c.Get(context.TODO(), tfJobKey, tfJob)
		tfJob.Status = commonv1.JobStatus{
			Conditions: []commonv1.JobCondition{
				{
					Type:    commonv1.JobSucceeded,
					Status:  corev1.ConditionTrue,
					Message: SucceededMessage,
					Reason:  SucceededReason,
				},
			},
		}
		// For TFJob we use c.Update() instead of c.Status().Update() to update status
		c.Update(context.TODO(), tfJob)

		c.Get(context.TODO(), trialKey, trial)
		isConditionCorrect := false
		for _, cond := range trial.Status.Conditions {
			if cond.Type == trialsv1beta1.TrialSucceeded && cond.Status == corev1.ConditionFalse &&
				cond.Reason == fmt.Sprintf("%v. Job reason: %v", TrialMetricsUnavailableReason, SucceededReason) &&
				cond.Message == fmt.Sprintf("Metrics are not available. Job message: %v", SucceededMessage) {
				isConditionCorrect = true
			}
		}
		return isConditionCorrect
	}, timeout).Should(gomega.BeTrue())

	// Expect that Trial is deleted.
	// TFJob can't be deleted because GC doesn't work in envtest and TFJob stuck in termination phase.
	// Ref: https://book.kubebuilder.io/reference/testing/envtest.html#testing-considerations.
	g.Eventually(func() bool {
		// Delete the Trial
		c.Delete(context.TODO(), trial)
		return errors.IsNotFound(c.Get(context.TODO(), trialKey, &trialsv1beta1.Trial{}))
	}, timeout).Should(gomega.BeTrue())
}

func TestReconcileBatchJob(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockManagerClient := managerclientmock.NewMockManagerClient(mockCtrl)

	// Setup the Manager and Controller. Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
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
			c.Get(context.TODO(), trialKey, updatedInstance)
			updatedInstance.Status = instance.Status
			err = r.updateStatus(updatedInstance)
		}
		return err
	}

	recFn := SetupTestReconcile(r)
	// Set Job resource
	trialResources := trialutil.GvkListFlag{
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
		g.Expect(mgr.Start(context.TODO())).NotTo(gomega.HaveOccurred())
	}()

	// Result for GetTrialObservationLog
	observationLog := &api_pb.GetObservationLogReply{
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

	mockManagerClient.EXPECT().GetTrialObservationLog(gomock.Any()).Return(observationLog, nil).AnyTimes()
	mockManagerClient.EXPECT().DeleteTrialObservationLog(gomock.Any()).Return(nil, nil).AnyTimes()

	// Test 1 - Regural Trial run with BatchJob
	trial := newFakeTrialBatchJob()
	batchJob := &batchv1.Job{}

	// Create the Trial
	g.Expect(c.Create(context.TODO(), trial)).NotTo(gomega.HaveOccurred())

	// Expect that BatchJob with appropriate name is created
	g.Eventually(func() error {
		return c.Get(context.TODO(), batchJobKey, batchJob)
	}, timeout).Should(gomega.Succeed())

	// Expect that Trial status is running
	g.Eventually(func() bool {
		c.Get(context.TODO(), trialKey, trial)
		return trial.IsRunning()
	}, timeout).Should(gomega.BeTrue())

	// Manually update BatchJob status to failed
	// Expect that Trial status is failed
	g.Eventually(func() bool {
		c.Get(context.TODO(), batchJobKey, batchJob)
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
		c.Status().Update(context.TODO(), batchJob)

		c.Get(context.TODO(), trialKey, trial)
		return trial.IsFailed()
	}, timeout).Should(gomega.BeTrue())

	// Expect that Trial is deleted
	g.Eventually(func() bool {
		// Delete the Trial
		c.Delete(context.TODO(), trial)
		return errors.IsNotFound(c.Get(context.TODO(), trialKey, &trialsv1beta1.Trial{}))
	}, timeout).Should(gomega.BeTrue())

	// Get BatchJob
	g.Expect(c.Get(context.TODO(), batchJobKey, batchJob)).NotTo(gomega.HaveOccurred())
	// Add completed status
	batchJob.Status = batchv1.JobStatus{
		Conditions: []batchv1.JobCondition{
			{
				Type:    batchv1.JobComplete,
				Status:  corev1.ConditionTrue,
				Message: "BatchJob completed test message",
				Reason:  "BatchJob completed test reason",
			},
		},
	}
	// Manually update BatchJob status to completed
	g.Expect(c.Status().Update(context.TODO(), batchJob)).NotTo(gomega.HaveOccurred())

	// Create the Trial
	trial = newFakeTrialBatchJob()
	g.Expect(c.Create(context.TODO(), trial)).NotTo(gomega.HaveOccurred())

	// Expect that Trial status is succeeded and metrics are properly populated
	// Metrics available because GetTrialObservationLog returns values
	g.Eventually(func() bool {
		c.Get(context.TODO(), trialKey, trial)
		return trial.IsSucceeded() &&
			len(trial.Status.Observation.Metrics) > 0 &&
			trial.Status.Observation.Metrics[0].Max == "0.99" &&
			trial.Status.Observation.Metrics[0].Min == "0.11" &&
			trial.Status.Observation.Metrics[0].Latest == "0.11"
	}, timeout).Should(gomega.BeTrue())

	// Expect that Trial is deleted
	// BatchJob can't be deleted because GC doesn't work in envtest and BatchJob stuck in termination phase.
	// Ref: https://book.kubebuilder.io/reference/testing/envtest.html#testing-considerations.
	g.Eventually(func() bool {
		// Delete the Trial
		c.Delete(context.TODO(), trial)
		return errors.IsNotFound(c.Get(context.TODO(), trialKey, &trialsv1beta1.Trial{}))
	}, timeout).Should(gomega.BeTrue())

	// Test 2 - Update status for empty Trial
	g.Expect(r.updateStatus(&trialsv1beta1.Trial{})).To(gomega.HaveOccurred())

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

func newFakeTrialTFJob() *trialsv1beta1.Trial {
	primaryContainer := "tensorflow"

	tfJob := &tfv1.TFJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeflow.org/v1",
			Kind:       "TFJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tfJobName,
			Namespace: namespace,
		},
		Spec: tfv1.TFJobSpec{
			TFReplicaSpecs: map[commonv1.ReplicaType]*commonv1.ReplicaSpec{
				tfv1.TFReplicaTypePS: {
					Replicas:      func() *int32 { i := int32(2); return &i }(),
					RestartPolicy: commonv1.RestartPolicyNever,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  primaryContainer,
									Image: "gcr.io/kubeflow-ci/tf-mnist-with-summaries:1.0",
									Command: []string{
										"python",
										"/var/tf_mnist/mnist_with_summaries.py",
										"--log_dir=/train/metrics",
										"--lr=0.01",
										"--num-layers=5",
									},
								},
							},
						},
					},
				},
				tfv1.TFReplicaTypeWorker: {
					Replicas:      func() *int32 { i := int32(4); return &i }(),
					RestartPolicy: commonv1.RestartPolicyNever,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  primaryContainer,
									Image: "gcr.io/kubeflow-ci/tf-mnist-with-summaries:1.0",
									Command: []string{
										"python",
										"/var/tf_mnist/mnist_with_summaries.py",
										"--log_dir=/train/metrics",
										"--lr=0.01",
										"--num-layers=5",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	runSpec, _ := util.ConvertObjectToUnstructured(tfJob)

	return &trialsv1beta1.Trial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trialName,
			Namespace: namespace,
		},
		Spec: trialsv1beta1.TrialSpec{
			PrimaryPodLabels:     experimentsv1beta1.DefaultKubeflowJobPrimaryPodLabels,
			PrimaryContainerName: primaryContainer,
			SuccessCondition:     experimentsv1beta1.DefaultKubeflowJobSuccessCondition,
			FailureCondition:     experimentsv1beta1.DefaultKubeflowJobFailureCondition,
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
							Image: "docker.io/kubeflowkatib/mxnet-mnist",
							Command: []string{
								"python3",
								"/opt/mxnet-mnist/mnist.py",
								"--lr=0.01",
								"--num-layers=5",
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
