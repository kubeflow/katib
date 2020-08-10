package trial

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	"github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	trialutil "github.com/kubeflow/katib/pkg/controller.v1beta1/trial/util"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	util "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	managerclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/trial/managerclient"
	kubeflowcommonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
)

const (
	trialName    = "test-trial"
	namespace    = "default"
	tfJobName    = "test-tfjob"
	batchJobName = "test-job"

	timeout = time.Second * 40
)

var trialKey = types.NamespacedName{Name: trialName, Namespace: namespace}
var tfJobKey = types.NamespacedName{Name: tfJobName, Namespace: namespace}
var batchJobKey = types.NamespacedName{Name: batchJobName, Namespace: namespace}

func init() {
	logf.SetLogger(logf.ZapLogger(true))
}

func TestAdd(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

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
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	r := &ReconcileTrial{
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		ManagerClient: mockManagerClient,
		recorder:      mgr.GetRecorder(ControllerName),
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
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Empty result for GetTrialObservationLog
	observationLog := &api_pb.GetObservationLogReply{
		ObservationLog: &api_pb.ObservationLog{
			MetricLogs: []*api_pb.MetricLog{},
		},
	}

	mockManagerClient.EXPECT().GetTrialObservationLog(gomock.Any()).Return(observationLog, nil).AnyTimes()
	mockManagerClient.EXPECT().DeleteTrialObservationLog(gomock.Any()).Return(nil, nil).AnyTimes()

	// Test - Regural Trial run with TFJob
	trial := newFakeTrial(newFakeTFJob())
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
	g.Eventually(func() bool {
		c.Get(context.TODO(), tfJobKey, tfJob)
		tfJob.Status = kubeflowcommonv1.JobStatus{
			Conditions: []kubeflowcommonv1.JobCondition{
				{
					Type:    kubeflowcommonv1.JobSucceeded,
					Status:  corev1.ConditionTrue,
					Message: "TFJob succeeded test message",
					Reason:  "TFJob succeeded test reason",
				},
			},
		}
		// For TFJob we use c.Update() instead of c.Status().Update() to update status
		c.Update(context.TODO(), tfJob)

		c.Get(context.TODO(), trialKey, trial)
		isConditionCorrect := false
		for _, cond := range trial.Status.Conditions {
			if cond.Type == trialsv1beta1.TrialSucceeded && cond.Status == corev1.ConditionFalse && cond.Reason == TrialMetricsUnavailableReason {
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
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	r := &ReconcileTrial{
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		ManagerClient: mockManagerClient,
		recorder:      mgr.GetRecorder(ControllerName),
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
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Result for GetTrialObservationLog
	observationLog := &api_pb.GetObservationLogReply{
		ObservationLog: &api_pb.ObservationLog{
			MetricLogs: []*api_pb.MetricLog{
				{
					TimeStamp: "2020-08-10T14:47:38+08:00",
					Metric: &api_pb.Metric{
						Name:  "accuracy",
						Value: "0.99",
					},
				},
				{
					TimeStamp: "2020-08-10T14:50:38+08:00",
					Metric: &api_pb.Metric{
						Name:  "accuracy",
						Value: "0.11",
					},
				},
			},
		},
	}

	mockManagerClient.EXPECT().GetTrialObservationLog(gomock.Any()).Return(observationLog, nil).AnyTimes()
	mockManagerClient.EXPECT().DeleteTrialObservationLog(gomock.Any()).Return(nil, nil).AnyTimes()

	// Test 1 - Regural Trial run with BatchJob
	trial := newFakeTrial(newFakeBatchJob())
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
	trial = newFakeTrial(newFakeBatchJob())
	g.Expect(c.Create(context.TODO(), trial)).NotTo(gomega.HaveOccurred())

	// Expect that Trial status is succeeded and metrics are properly populated
	// Metrics available because GetTrialObservationLog returns something
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
		{TimeStamp: "2020-04-13T14:47:38+08:00", Metric: &api_pb.Metric{Name: "accuracy", Value: "0.7"}},
		{TimeStamp: "2020-04-13T14:47:39+08:00", Metric: &api_pb.Metric{Name: "accuracy", Value: "0.71"}},
		{TimeStamp: "2020-04-13T14:47:40+08:00", Metric: &api_pb.Metric{Name: "accuracy", Value: "0.72"}},
		{TimeStamp: "2020-04-13T14:47:41+08:00", Metric: &api_pb.Metric{Name: "accuracy", Value: "0.68"}},
		{TimeStamp: "2020-04-13T14:47:41+08:00", Metric: &api_pb.Metric{Name: "accuracy", Value: "0.69"}},
		{TimeStamp: "2020-04-13T14:47:41+08:00", Metric: &api_pb.Metric{Name: "accuracy", Value: "0.67"}},
		{TimeStamp: "2020-04-12T14:47:42+08:00", Metric: &api_pb.Metric{Name: "accuracy", Value: "0.6"}},
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
			} else if metric.Name == "accuracy" {
				accMetric = &observation.Metrics[index]
			}
		}
		return errMetric, accMetric, nil
	}

	metricStrategies := []commonv1beta1.MetricStrategy{
		{Name: "error", Value: commonv1beta1.ExtractByMin},
		{Name: "accuracy", Value: commonv1beta1.ExtractByMax},
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
		{TimeStamp: "2020-08-10T14:47:42", Metric: &api_pb.Metric{Name: "accuracy", Value: "0.77"}},
	}
	_, err = getMetrics(invalidLogs, metricStrategies)
	g.Expect(err).To(gomega.HaveOccurred())
}

func newFakeTFJob() *tfv1.TFJob {
	return &tfv1.TFJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeflow.org/v1",
			Kind:       "TFJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tfJobName,
			Namespace: namespace,
		},
		Spec: tfv1.TFJobSpec{
			TFReplicaSpecs: map[tfv1.TFReplicaType]*kubeflowcommonv1.ReplicaSpec{
				tfv1.TFReplicaTypePS: {
					Replicas:      func() *int32 { i := int32(2); return &i }(),
					RestartPolicy: kubeflowcommonv1.RestartPolicyNever,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  "tensorflow",
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
					RestartPolicy: kubeflowcommonv1.RestartPolicyNever,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  "tensorflow",
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
}

func newFakeBatchJob() *batchv1.Job {
	return &batchv1.Job{
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
							Name:  "training-container",
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
}

func newFakeTrial(runObject interface{}) *trialsv1beta1.Trial {

	runSpec, _ := util.ConvertObjectToUnstructured(runObject)

	t := &trialsv1beta1.Trial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trialName,
			Namespace: namespace,
		},
		Spec: trialsv1beta1.TrialSpec{
			Objective: &commonv1beta1.ObjectiveSpec{
				ObjectiveMetricName: "accuracy",
				MetricStrategies: []commonv1beta1.MetricStrategy{
					{
						Name:  "accuracy",
						Value: commonv1beta1.ExtractByMax,
					},
				},
			},
			RunSpec: runSpec,
		},
	}
	return t
}
