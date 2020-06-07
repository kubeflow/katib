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

	util "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	managerclientmock "github.com/kubeflow/katib/pkg/mock/v1beta1/trial/managerclient"
	kubeflowcommonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
)

const (
	trialName      = "trial-name"
	trialNamespace = "trial-namespace"
	tfJobName      = "tfjob-name"

	timeout = time.Second * 40
)

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: trialName, Namespace: trialNamespace}}
var expectedResult = reconcile.Result{Requeue: true}
var tfJobKey = types.NamespacedName{Name: tfJobName, Namespace: trialNamespace}

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
		updateStatusHandler: func(instance *trialsv1beta1.Trial) error {
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

	r.updateStatusHandler = func(instance *trialsv1beta1.Trial) error {
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

	tfJob := instance.Spec.RunSpec

	g.Eventually(func() error { return c.Get(context.TODO(), tfJobKey, tfJob) }, timeout).
		Should(gomega.Succeed())

	// Delete the TFJob and expect Reconcile to be called for TFJob deletion
	g.Expect(c.Delete(context.TODO(), tfJob)).NotTo(gomega.HaveOccurred())
	g.Eventually(func() error { return c.Get(context.TODO(), tfJobKey, tfJob) }, timeout).
		Should(gomega.Succeed())

	// Manually delete TFJob since GC isn't enabled in the test control plane
	g.Eventually(func() error { return c.Delete(context.TODO(), tfJob) }, timeout).
		Should(gomega.MatchError("tfjobs.kubeflow.org \"tfjob-name\" not found"))
	g.Expect(c.Delete(context.TODO(), instance)).NotTo(gomega.HaveOccurred())
}

func TestReconcileCompletedTFJobTrial(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := newFakeTrialWithTFJob()
	// Trial name must be different to avoid errors
	instance.Name = "new-trial-name"

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

	r.updateStatusHandler = func(instance *trialsv1beta1.Trial) error {
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

func newFakeTrialWithTFJob() *trialsv1beta1.Trial {
	objectiveSpec := commonv1beta1.ObjectiveSpec{ObjectiveMetricName: "test"}
	runSpecTFJob := &tfv1.TFJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeflow.org/v1",
			Kind:       "TFJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tfJobName,
			Namespace: trialNamespace,
		},
		Spec: tfv1.TFJobSpec{
			TFReplicaSpecs: map[tfv1.TFReplicaType]*kubeflowcommonv1.ReplicaSpec{
				tfv1.TFReplicaTypePS: &kubeflowcommonv1.ReplicaSpec{
					Replicas:      func() *int32 { i := int32(2); return &i }(),
					RestartPolicy: kubeflowcommonv1.RestartPolicyNever,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
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
				tfv1.TFReplicaTypeWorker: &kubeflowcommonv1.ReplicaSpec{
					Replicas:      func() *int32 { i := int32(4); return &i }(),
					RestartPolicy: kubeflowcommonv1.RestartPolicyNever,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
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

	runSpec, _ := util.ConvertObjectToUnstructured(runSpecTFJob)

	t := &trialsv1beta1.Trial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trialName,
			Namespace: trialNamespace,
		},
		Spec: trialsv1beta1.TrialSpec{
			Objective: &objectiveSpec,
			RunSpec:   runSpec,
		},
	}
	return t
}
