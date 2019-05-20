package experiment

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	managerclientmock "github.com/kubeflow/katib/pkg/mock/v1alpha2/experiment/managerclient"
	manifestmock "github.com/kubeflow/katib/pkg/mock/v1alpha2/experiment/manifest"
	suggestionmock "github.com/kubeflow/katib/pkg/mock/v1alpha2/experiment/suggestion"
)

const (
	experimentName = "foo"
	namespace      = "default"

	timeout = time.Second * 20
)

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: experimentName, Namespace: namespace}}
var trialKey = types.NamespacedName{Name: "test", Namespace: namespace}

func init() {
	logf.SetLogger(logf.ZapLogger(true))
}

func TestCreateExperiment(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := newFakeInstance()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mc := managerclientmock.NewMockManagerClient(mockCtrl)
	mc.EXPECT().CreateExperimentInDB(gomock.Any()).Return(nil).AnyTimes()
	mc.EXPECT().UpdateExperimentStatusInDB(gomock.Any()).Return(nil).AnyTimes()

	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()
	suggestion := suggestionmock.NewMockSuggestion(mockCtrl)

	mockCtrl3 := gomock.NewController(t)
	defer mockCtrl3.Finish()
	generator := manifestmock.NewMockGenerator(mockCtrl)

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	recFn, requests := SetupTestReconcile(&ReconcileExperiment{
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		ManagerClient: mc,
		Suggestion:    suggestion,
		Generator:     generator,
		updateStatusHandler: func(instance *experimentsv1alpha2.Experiment) error {
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
	defer c.Delete(context.TODO(), instance)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
}

func newFakeInstance() *experimentsv1alpha2.Experiment {
	return &experimentsv1alpha2.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      experimentName,
			Namespace: namespace,
		},
	}
}
