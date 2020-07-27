package experiment

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
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
)

const (
	experimentName = "test-experiment"
	trialName      = "test-trial"
	namespace      = "default"

	timeout = time.Second * 40
)

func init() {
	logf.SetLogger(logf.ZapLogger(true))
}

type statusMatcher struct {
	x *suggestionsv1beta1.Suggestion
}

func (statusM statusMatcher) Matches(x interface{}) bool {
	// Cast interface to suggestion object
	s := x.(*suggestionsv1beta1.Suggestion)
	// Verify that status is correct
	return (s.Status.Conditions[0].Type == statusM.x.Status.Conditions[0].Type &&
		s.Status.Conditions[0].Reason == statusM.x.Status.Conditions[0].Reason &&
		s.Status.Conditions[0].Message == statusM.x.Status.Conditions[0].Message)
}

func (statusM statusMatcher) String() string {
	return fmt.Sprintf("status is equal %v", statusM.x.Status.Conditions)
}

func TestAdd(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// Test - Try to add experiment controller to the manager
	g.Expect(Add(mgr)).NotTo(gomega.HaveOccurred())
}

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSuggestion := suggestionmock.NewMockSuggestion(mockCtrl)

	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()
	mockGenerator := manifestmock.NewMockGenerator(mockCtrl)

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	r := &ReconcileExperiment{
		Client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		Suggestion: mockSuggestion,
		Generator:  mockGenerator,
		collector:  experimentUtil.NewExpsCollector(mgr.GetCache(), prometheus.NewRegistry()),
	}
	r.updateStatusHandler = func(instance *experimentsv1beta1.Experiment) error {
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

	returnedTFJob := newFakeTFJob()

	returnedUnstructured, err := util.ConvertObjectToUnstructured(returnedTFJob)
	if err != nil {
		t.Errorf("ConvertObjectToUnstructured failed: %v", err)
	}

	mockGenerator.EXPECT().GetRunSpecWithHyperParameters(gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any()).Return(
		returnedUnstructured,
		nil).AnyTimes()

	suggestion := newFakeSuggestion()

	mockSuggestion.EXPECT().GetOrCreateSuggestion(gomock.Any(), gomock.Any()).Return(
		suggestion, nil).AnyTimes()

	mockSuggestion.EXPECT().UpdateSuggestion(gomock.Any()).Return(nil).AnyTimes()

	// Test 1 - Regural experiment run
	instance := newFakeInstance()

	// Create the experiment object
	g.Expect(c.Create(context.TODO(), instance)).NotTo(gomega.HaveOccurred())

	// Check that experiment was properly created
	g.Eventually(func() bool {
		experiment := &experimentsv1beta1.Experiment{}
		c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: experimentName}, experiment)
		return experiment.Name == instance.Name
	}, timeout).Should(gomega.BeTrue())

	// Check that 2 trials were properly created, 1 should be deleted because ParallelTrialCount=2
	trials := &trialsv1beta1.TrialList{}
	label := labels.Set{
		consts.LabelExperimentName: experimentName,
	}
	g.Eventually(func() int {
		c.List(context.TODO(), &client.ListOptions{
			LabelSelector: label.AsSelector(),
		}, trials)
		return len(trials.Items)
	}, timeout).
		Should(gomega.Equal(2))

	// Delete the experiment object
	g.Expect(c.Delete(context.TODO(), instance)).NotTo(gomega.HaveOccurred())

	// Check that experiment was properly deleted
	g.Eventually(func() bool {
		return errors.IsNotFound(c.Get(context.TODO(),
			types.NamespacedName{Namespace: namespace, Name: experimentName}, &experimentsv1beta1.Experiment{}))
	}, timeout).Should(gomega.BeTrue())

	// Test 2 - Update status for empty experiment object
	g.Expect(r.updateStatus(&experimentsv1beta1.Experiment{})).To(gomega.HaveOccurred())

	// Test 3 - Cleanup suggestion resources without deployed suggestion
	g.Expect(r.cleanupSuggestionResources(instance)).NotTo(gomega.HaveOccurred())
}

func TestCompleteExperiment(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSuggestion := suggestionmock.NewMockSuggestion(mockCtrl)

	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()
	mockGenerator := manifestmock.NewMockGenerator(mockCtrl)

	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	r := &ReconcileExperiment{
		Client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		Suggestion: mockSuggestion,
		Generator:  mockGenerator,
		collector:  experimentUtil.NewExpsCollector(mgr.GetCache(), prometheus.NewRegistry()),
	}
	r.updateStatusHandler = func(instance *experimentsv1beta1.Experiment) error {
		// Manually mark experiment succeeded after created to test cleanup + resuming experiment
		if instance.IsCreated() {
			instance.MarkExperimentStatusSucceeded(experimentUtil.ExperimentMaxTrialsReachedReason, "Experiment is succeeded")
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

	suggestionRestartNo := newFakeSuggestion()
	failedSuggestion := newFakeSuggestion()

	// After restarting experiment, GetOrCreateSuggestion is executed
	failedSuggestion.MarkSuggestionStatusFailed("test", "test")
	mockSuggestion.EXPECT().GetOrCreateSuggestion(gomock.Any(), gomock.Any()).Return(
		failedSuggestion, nil).AnyTimes()

	suggestionRestartYes := newFakeSuggestion()

	suggestionRestarting := newFakeSuggestion()

	reason := "Experiment is succeeded"
	msg := "Suggestion is succeeded, can't be restarted"
	suggestionRestartNo.MarkSuggestionStatusSucceeded(reason, msg)

	msg = "Suggestion is succeeded, suggestion volume is not deleted, can be restarted"
	suggestionRestartYes.MarkSuggestionStatusSucceeded(reason, msg)

	reason = "Experiment is restarting"
	msg = "Suggestion is not running"
	suggestionRestarting.MarkSuggestionStatusRunning(corev1.ConditionFalse, reason, msg)

	mockSuggestion.EXPECT().UpdateSuggestionStatus(statusMatcher{suggestionRestartNo}).Return(nil).MinTimes(1)
	mockSuggestion.EXPECT().UpdateSuggestionStatus(statusMatcher{suggestionRestartYes}).Return(nil).MinTimes(1)
	mockSuggestion.EXPECT().UpdateSuggestionStatus(statusMatcher{suggestionRestarting}).Return(nil).MinTimes(1)

	nameRestartNo := "test-never-resume"
	nameRestartYes := "test-from-volume"
	instanceRestartNo := newFakeInstance()
	instanceRestartNo.Name = nameRestartNo

	instanceRestartNo.Spec.ResumePolicy = experimentsv1beta1.NeverResume
	suggestionRestartNo.Name = nameRestartNo

	instanceRestartYes := newFakeInstance()
	instanceRestartYes.Name = nameRestartYes
	instanceRestartYes.Spec.ResumePolicy = experimentsv1beta1.FromVolume
	suggestionRestartYes.Name = nameRestartYes

	// Test - Mark correct status for suggestion object after experiment is completed.
	// Create experiment and suggestion objects for different resume policies
	g.Expect(c.Create(context.TODO(), instanceRestartNo)).NotTo(gomega.HaveOccurred())
	g.Expect(c.Create(context.TODO(), suggestionRestartNo)).NotTo(gomega.HaveOccurred())
	g.Expect(c.Create(context.TODO(), instanceRestartYes)).NotTo(gomega.HaveOccurred())
	g.Expect(c.Create(context.TODO(), suggestionRestartYes)).NotTo(gomega.HaveOccurred())

	// Check that experiment with NeverResume is succeeded
	g.Eventually(func() bool {
		experiment := &experimentsv1beta1.Experiment{}
		c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: nameRestartNo}, experiment)
		return len(experiment.Status.Conditions) > 1 && experiment.Status.Conditions[1].Type == experimentsv1beta1.ExperimentSucceeded
	}, timeout).Should(gomega.BeTrue())

	// Check that experiment with FromVolume is succeeded
	g.Eventually(func() bool {
		experiment := &experimentsv1beta1.Experiment{}
		c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: nameRestartYes}, experiment)
		return len(experiment.Status.Conditions) > 1 && experiment.Status.Conditions[1].Type == experimentsv1beta1.ExperimentSucceeded
	}, timeout).Should(gomega.BeTrue())

	// Check that experiment with FromVolume is restarting
	g.Eventually(func() bool {
		experiment := &experimentsv1beta1.Experiment{}
		c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: nameRestartYes}, experiment)
		return len(experiment.Status.Conditions) > 1 && experiment.Status.Conditions[1].Type == experimentsv1beta1.ExperimentRestarting
	}, timeout).Should(gomega.BeTrue())

}

func newFakeInstance() *experimentsv1beta1.Experiment {
	var parallelCount int32 = 2
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
			Algorithm: &commonapiv1beta1.AlgorithmSpec{
				AlgorithmName: "random",
			},
			MetricsCollectorSpec: &commonapiv1beta1.MetricsCollectorSpec{
				Collector: &commonapiv1beta1.CollectorSpec{
					Kind: commonapiv1beta1.StdOutCollector,
				},
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

func newFakeSuggestion() *suggestionsv1beta1.Suggestion {
	return &suggestionsv1beta1.Suggestion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      experimentName,
			Namespace: namespace,
		},
		Status: suggestionsv1beta1.SuggestionStatus{
			Suggestions: []suggestionsv1beta1.TrialAssignment{
				{
					Name: trialName + "-1",
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
				{
					Name: trialName + "-2",
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
				{
					Name: trialName + "-3",
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
	}
}

func newFakeTFJob() *tfv1.TFJob {
	return &tfv1.TFJob{
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
								},
							},
						},
					},
				},
			},
		},
	}
}
