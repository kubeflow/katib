package trial

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
)

const (
	trialName = "foo"
	namespace = "default"

	timeout = time.Second * 5
)

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: trialName, Namespace: namespace}}

func TestCreate(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := newFakeTrial()

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	recFn, requests := SetupTestReconcile(newReconciler(mgr))
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

func newFakeTrial() *trialsv1alpha2.Trial {
	t := &trialsv1alpha2.Trial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trialName,
			Namespace: namespace,
		},
		Spec: trialsv1alpha2.TrialSpec{
			RunSpec: `
apiVersion: "kubeflow.org/v1beta1"
kind: "TFJob"
metadata:
	name: "dist-mnist-for-e2e-test"
spec:
	tfReplicaSpecs:
		Worker:
			replicas: 1
			restartPolicy: Never
			template:
				spec:
					containers:
					  - name: tensorflow
						image: gaocegege/mnist:1`,
		},
	}
	return t
}
