package trial

import (
	"bytes"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	managerclientmock "github.com/kubeflow/katib/pkg/mock/v1alpha3/trial/managerclient"
)

const (
	trialName = "foo"
	namespace = "default"

	timeout = time.Second * 40
)

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: trialName, Namespace: namespace}}
var expectedResult = reconcile.Result{Requeue: true}
var tfJobKey = types.NamespacedName{Name: "test", Namespace: namespace}

func init() {
	logf.SetLogger(logf.ZapLogger(true))
}

func TestCreateTFJobTrial(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := newFakeTrialWithTFJob()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mc := managerclientmock.NewMockManagerClient(mockCtrl)

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	recFn := SetupTestReconcile(&ReconcileTrial{
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		ManagerClient: mc,
		recorder:      mgr.GetRecorder(ControllerName),
		updateStatusHandler: func(instance *trialsv1alpha3.Trial) error {
			if !instance.IsCreated() {
				t.Errorf("Expected got condition created")
			}
			return nil
		},
	})
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the Trial object and expect the Reconcile and Deployment to be created
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(c.Delete(context.TODO(), instance)).NotTo(gomega.HaveOccurred())
	g.Eventually(func() bool {
		return apierrors.IsNotFound(c.Get(context.TODO(),
			expectedRequest.NamespacedName, instance))
	}, timeout).Should(gomega.BeTrue())
}

func TestReconcileTFJobTrial(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := newFakeTrialWithTFJob()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mc := managerclientmock.NewMockManagerClient(mockCtrl)
	mc.EXPECT().GetTrialObservationLog(gomock.Any()).Return(&api_pb.GetObservationLogReply{
		ObservationLog: nil,
	}, nil).AnyTimes()
	mc.EXPECT().DeleteTrialObservationLog(gomock.Any()).Return(&api_pb.DeleteObservationLogReply{}, nil).AnyTimes()

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	r := &ReconcileTrial{
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		ManagerClient: mc,
		recorder:      mgr.GetRecorder(ControllerName),
		collector:     NewTrialsCollector(mgr.GetCache(), prometheus.NewRegistry()),
	}

	r.updateStatusHandler = func(instance *trialsv1alpha3.Trial) error {
		if !instance.IsCreated() {
			t.Errorf("Expected got condition created")
		}
		return r.updateStatus(instance)
	}

	recFn := SetupTestReconcile(r)
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the Trial object and expect the Reconcile and Deployment to be created
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())

	tfJob := &unstructured.Unstructured{}
	bufSize := 1024
	buf := bytes.NewBufferString(instance.Spec.RunSpec)
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, bufSize).Decode(tfJob); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	g.Eventually(func() error { return c.Get(context.TODO(), tfJobKey, tfJob) }, timeout).
		Should(gomega.Succeed())

	// Delete the TFJob and expect Reconcile to be called for TFJob deletion
	g.Expect(c.Delete(context.TODO(), tfJob)).NotTo(gomega.HaveOccurred())
	g.Eventually(func() error { return c.Get(context.TODO(), tfJobKey, tfJob) }, timeout).
		Should(gomega.Succeed())

	// Manually delete TFJob since GC isn't enabled in the test control plane
	g.Eventually(func() error { return c.Delete(context.TODO(), tfJob) }, timeout).
		Should(gomega.MatchError("tfjobs.kubeflow.org \"test\" not found"))
	g.Expect(c.Delete(context.TODO(), instance)).NotTo(gomega.HaveOccurred())
}

func TestReconcileCompletedTFJobTrial(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := newFakeTrialWithTFJob()
	instance.Name = "tfjob-trial"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mc := managerclientmock.NewMockManagerClient(mockCtrl)
	mc.EXPECT().GetTrialObservationLog(gomock.Any()).Return(&api_pb.GetObservationLogReply{
		ObservationLog: nil,
	}, nil).AnyTimes()
	mc.EXPECT().DeleteTrialObservationLog(gomock.Any()).Return(&api_pb.DeleteObservationLogReply{}, nil).AnyTimes()

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	r := &ReconcileTrial{
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		ManagerClient: mc,
		recorder:      mgr.GetRecorder(ControllerName),
		collector:     NewTrialsCollector(mgr.GetCache(), prometheus.NewRegistry()),
	}

	r.updateStatusHandler = func(instance *trialsv1alpha3.Trial) error {
		if !instance.IsCreated() {
			t.Errorf("Expected got condition created")
		}
		return r.updateStatus(instance)
	}

	recFn := SetupTestReconcile(r)
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the Trial object and expect the Reconcile and Deployment to be created
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), instance)

	g.Eventually(func() error {
		return c.Get(context.TODO(), expectedRequest.NamespacedName, instance)
	}, timeout).
		Should(gomega.Succeed())
	instance.MarkTrialStatusSucceeded(corev1.ConditionTrue, "", "")
	g.Expect(c.Status().Update(context.TODO(), instance)).NotTo(gomega.HaveOccurred())
	g.Eventually(func() bool {
		err := c.Get(context.TODO(), expectedRequest.NamespacedName, instance)
		if err == nil && instance.IsCompleted() {
			return true
		}
		return false
	}, timeout).
		Should(gomega.BeTrue())
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

	getMetricsFromLogs := func(strategies map[string]commonv1alpha3.MetricStrategy) (*commonv1alpha3.Metric, *commonv1alpha3.Metric, error) {
		observation, err := getMetrics(metricLogs, strategies)
		if err != nil {
			return nil, nil, err
		}
		var errMetric, accMetric *commonv1alpha3.Metric
		for _, metric := range observation.Metrics {
			if metric.Name == "error" {
				errMetric = &metric
			} else if metric.Name == "accuracy" {
				accMetric = &metric
			}
		}
		return errMetric, accMetric, nil
	}

	metricStrategies := map[string]commonv1alpha3.MetricStrategy{
		"error":    commonv1alpha3.ExtractByMin,
		"accuracy": commonv1alpha3.ExtractByMax,
	}
	errMetric, accMetric, err := getMetricsFromLogs(metricStrategies)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	g.Expect(errMetric.Latest).To(gomega.Equal(0.07))
	g.Expect(errMetric.Max).To(gomega.Equal(0.1))
	g.Expect(errMetric.Min).To(gomega.Equal(0.01))
	g.Expect(accMetric.Latest).To(gomega.Equal(0.67))
	g.Expect(accMetric.Max).To(gomega.Equal(0.72))
	g.Expect(accMetric.Min).To(gomega.Equal(0.6))
}

func newFakeTrialWithTFJob() *trialsv1alpha3.Trial {
	objectiveSpec := commonv1alpha3.ObjectiveSpec{ObjectiveMetricName: "test"}
	t := &trialsv1alpha3.Trial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trialName,
			Namespace: namespace,
		},
		Spec: trialsv1alpha3.TrialSpec{
			Objective: &objectiveSpec,
			RunSpec: `apiVersion: "kubeflow.org/v1"
kind: "TFJob"
metadata:
  name: "test"
  namespace: "default"
spec:
  tfReplicaSpecs:
    PS:
      replicas: 2
      restartPolicy: Never
      template:
        spec:
          containers:
            - name: tensorflow
              image: kubeflow/tf-dist-mnist-test:1.0
    Worker:
      replicas: 4
      restartPolicy: Never
      template:
        spec:
          containers:
            - name: tensorflow
              image: kubeflow/tf-dist-mnist-test:1.0
`,
		},
	}
	return t
}
