package experiment

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	experimentUtil "github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/util"
	util "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	manifestmock "github.com/kubeflow/katib/pkg/mock/v1beta1/experiment/manifest"
	suggestionmock "github.com/kubeflow/katib/pkg/mock/v1beta1/experiment/suggestion"
	kubeflowcommonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	v1 "k8s.io/api/core/v1"
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
		updateStatusHandler: func(instance *experimentsv1beta1.Experiment) error {
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

	// Create the experiment object and expect the Reconcile and Deployment to be created
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
		&suggestionsv1beta1.Suggestion{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instance.Name,
				Namespace: instance.Namespace,
			},
			Status: suggestionsv1beta1.SuggestionStatus{
				Suggestions: []suggestionsv1beta1.TrialAssignment{
					{
						Name: trialKey.Name,
						ParameterAssignments: []commonapiv1beta1.ParameterAssignment{
							{
								Name:  "lr",
								Value: "0.01",
							},
							{
								Name:  "num-layers",
								Value: "5",
							},
						},
					},
				},
			},
		}, nil).AnyTimes()
	suggestion.EXPECT().UpdateSuggestion(gomock.Any()).Return(nil).AnyTimes()
	mockCtrl3 := gomock.NewController(t)
	defer mockCtrl3.Finish()
	generator := manifestmock.NewMockGenerator(mockCtrl)

	returnedTFJob := &tfv1.TFJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeflow.org/v1",
			Kind:       "TFJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trial-name",
			Namespace: "trial-namespace",
		},
		Spec: tfv1.TFJobSpec{
			TFReplicaSpecs: map[tfv1.TFReplicaType]*kubeflowcommonv1.ReplicaSpec{
				tfv1.TFReplicaTypePS: {
					Replicas:      func() *int32 { i := int32(1); return &i }(),
					RestartPolicy: kubeflowcommonv1.RestartPolicyOnFailure,
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
									VolumeMounts: []v1.VolumeMount{
										{
											Name:      "train",
											MountPath: "/train",
										},
									},
								},
							},
							Volumes: []v1.Volume{
								{
									Name: "train",
									VolumeSource: v1.VolumeSource{
										PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
											ClaimName: "tfevent-volume",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	returnedUnstructured, err := util.ConvertObjectToUnstructured(returnedTFJob)
	if err != nil {
		t.Errorf("ConvertObjectToUnstructured failed: %v", err)
	}
	generator.EXPECT().GetRunSpecWithHyperParameters(gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any()).Return(
		returnedUnstructured,
		nil).AnyTimes()

	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	r := &ReconcileExperiment{
		Client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		Suggestion: suggestion,
		Generator:  generator,
		collector:  experimentUtil.NewExpsCollector(mgr.GetCache(), prometheus.NewRegistry()),
	}
	r.updateStatusHandler = func(instance *experimentsv1beta1.Experiment) error {
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

	// Create the Experiment object and expect the Reconcile and Deployment to be created
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())

	trials := &trialsv1beta1.TrialList{}
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

func newFakeInstance() *experimentsv1beta1.Experiment {
	var parallelCount int32 = 1
	var goal float64 = 99.9

	trialTemplateJob := &tfv1.TFJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeflow.org/v1",
			Kind:       "TFJob",
		},
		Spec: tfv1.TFJobSpec{
			TFReplicaSpecs: map[tfv1.TFReplicaType]*kubeflowcommonv1.ReplicaSpec{
				tfv1.TFReplicaTypePS: {
					Replicas:      func() *int32 { i := int32(1); return &i }(),
					RestartPolicy: kubeflowcommonv1.RestartPolicyOnFailure,
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
										"--lr=${trialParameters.learningRate}",
										"--num-layers=${trialParameters.numberLayers}",
									},
									VolumeMounts: []v1.VolumeMount{
										{
											Name:      "train",
											MountPath: "/train",
										},
									},
								},
							},
							Volumes: []v1.Volume{
								{
									Name: "train",
									VolumeSource: v1.VolumeSource{
										PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
											ClaimName: "tfevent-volume",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	trialSpec, _ := util.ConvertObjectToUnstructured(trialTemplateJob)

	return &experimentsv1beta1.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      experimentName,
			Namespace: namespace,
		},
		Spec: experimentsv1beta1.ExperimentSpec{
			ParallelTrialCount: &parallelCount,
			MaxTrialCount:      &parallelCount,
			Objective: &commonapiv1beta1.ObjectiveSpec{
				Type:                commonapiv1beta1.ObjectiveTypeMaximize,
				Goal:                &goal,
				ObjectiveMetricName: "accuracy",
			},
			TrialTemplate: &experimentsv1beta1.TrialTemplate{
				TrialParameters: []experimentsv1beta1.TrialParameterSpec{
					{
						Name:        "learningRate",
						Description: "Learning Rate",
						Reference:   "lr",
					},
					{
						Name:        "numberLayers",
						Description: "Number of layers",
						Reference:   "num-layers",
					},
				},
				TrialSource: experimentsv1beta1.TrialSource{
					TrialSpec: trialSpec,
				},
			},
		},
	}
}
