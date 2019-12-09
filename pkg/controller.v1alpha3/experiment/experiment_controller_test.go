package experiment

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/util"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonapiv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	manifestmock "github.com/kubeflow/katib/pkg/mock/v1alpha3/experiment/manifest"
	suggestionmock "github.com/kubeflow/katib/pkg/mock/v1alpha3/experiment/suggestion"
)

const (
	experimentName = "foo"
	namespace      = "default"

	timeout = time.Second * 40
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

	recFn := SetupTestReconcile(&ReconcileExperiment{
		Client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		Suggestion: suggestion,
		Generator:  generator,
		updateStatusHandler: func(instance *experimentsv1alpha3.Experiment) error {
			if !instance.IsCreated() {
				t.Errorf("Expected got condition created")
			}
			return nil
		},
	})
	g.Expect(addForTestPurpose(mgr, recFn)).NotTo(gomega.HaveOccurred())

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
		return errors.IsNotFound(c.Get(context.TODO(),
			expectedRequest.NamespacedName, instance))
	}, timeout).Should(gomega.BeTrue())
}

func TestReconcileExperiment(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	testName := "tn"
	instance := newFakeInstance()
	instance.Name = testName

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()
	suggestion := suggestionmock.NewMockSuggestion(mockCtrl)
	suggestion.EXPECT().GetOrCreateSuggestion(gomock.Any(), gomock.Any()).Return(
		&suggestionsv1alpha3.Suggestion{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instance.Name,
				Namespace: instance.Namespace,
			},
			Status: suggestionsv1alpha3.SuggestionStatus{
				Suggestions: []suggestionsv1alpha3.TrialAssignment{
					{
						Name: trialKey.Name,
						ParameterAssignments: []commonapiv1alpha3.ParameterAssignment{
							{
								Name:  "--lr",
								Value: "0.5",
							},
						},
					},
				},
			},
		}, nil).AnyTimes()
	suggestion.EXPECT().UpdateSuggestion(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockCtrl3 := gomock.NewController(t)
	defer mockCtrl3.Finish()
	generator := manifestmock.NewMockGenerator(mockCtrl)
	generator.EXPECT().GetRunSpecWithHyperParameters(gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(`apiVersion: "kubeflow.org/v1"
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
              image: kubeflow/tf-dist-mnist-test:1.0`, nil).AnyTimes()

	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	r := &ReconcileExperiment{
		Client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		Suggestion: suggestion,
		Generator:  generator,
		collector:  util.NewExpsCollector(mgr.GetCache()),
	}
	r.updateStatusHandler = func(instance *experimentsv1alpha3.Experiment) error {
		if !instance.IsCreated() {
			t.Errorf("Expected got condition created")
		}
		return r.updateStatus(instance)
	}

	recFn := SetupTestReconcile(r)
	g.Expect(addForTestPurpose(mgr, recFn)).NotTo(gomega.HaveOccurred())

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

	trials := &trialsv1alpha3.TrialList{}
	g.Eventually(func() int {
		label := labels.Set{
			consts.LabelExperimentName: testName,
		}
		c.List(context.TODO(), &client.ListOptions{
			LabelSelector: label.AsSelector(),
		}, trials)
		return len(trials.Items)
	}, timeout).
		Should(gomega.Equal(1))

	g.Expect(c.Delete(context.TODO(), instance)).NotTo(gomega.HaveOccurred())
	g.Eventually(func() bool {
		return errors.IsNotFound(c.Get(context.TODO(),
			types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}, instance))
	}, timeout).Should(gomega.BeTrue())
}

func newFakeInstance() *experimentsv1alpha3.Experiment {
	var parallelCount int32 = 1
	var goal float64 = 99.9
	return &experimentsv1alpha3.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      experimentName,
			Namespace: namespace,
		},
		Spec: experimentsv1alpha3.ExperimentSpec{
			ParallelTrialCount: &parallelCount,
			MaxTrialCount:      &parallelCount,
			Objective: &commonapiv1alpha3.ObjectiveSpec{
				Type:                commonapiv1alpha3.ObjectiveTypeMaximize,
				Goal:                &goal,
				ObjectiveMetricName: "accuracy",
			},
			TrialTemplate: &experimentsv1alpha3.TrialTemplate{
				GoTemplate: &experimentsv1alpha3.GoTemplate{
					RawTemplate: `apiVersion: "kubeflow.org/v1"
kind: TFJob
metadata:
  name: {{.Trial}}
  namespace: {{.NameSpace}}
spec:
  tfReplicaSpecs:
  Worker:
    replicas: 1 
    restartPolicy: OnFailure
    template:
      spec:
        containers:
          - name: tensorflow 
            image: gcr.io/kubeflow-ci/tf-mnist-with-summaries:1.0
            imagePullPolicy: Always
            command:
              - "python"
              - "/var/tf_mnist/mnist_with_summaries.py"
              - "--log_dir=/train/{{.Trial}}"
              {{- with .HyperParameters}}
              {{- range .}}
              - "{{.Name}}={{.Value}}"
              {{- end}}
              {{- end}}
            volumeMounts:
              - mountPath: "/train"
                name: "train"
        volumes:
          - name: "train"
            persistentVolumeClaim:
              claimName: "tfevent-volume"`,
				},
			},
		},
	}
}
