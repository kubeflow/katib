package validator

import (
	"testing"

	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	manifestmock "github.com/kubeflow/katib/pkg/mock/v1alpha3/experiment/manifest"
)

func init() {
	logf.SetLogger(logf.ZapLogger(false))
}

func TestValidateTFJobTrialTemplate(t *testing.T) {
	trialTFJobTemplate := `apiVersion: "kubeflow.org/v1"
kind: "TFJob"
metadata:
    name: "dist-mnist-for-e2e-test"
spec:
    tfReplicaSpecs:
        Worker:
            template:
                spec:
                    containers:
                      - name: tensorflow
                        image: gaocegege/mnist:1`

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()

	p := manifestmock.NewMockGenerator(mockCtrl)
	g := New(p)

	p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(trialTFJobTemplate, nil)

	instance := newFakeInstance()
	if err := g.(*DefaultValidator).validateTrialTemplate(instance); err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestValidateJobTrialTemplate(t *testing.T) {
	trialJobTemplate := `apiVersion: "batch/v1"
kind: "Job"
metadata:
  name: "fake-trial"
  namespace: fakens`

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()

	p := manifestmock.NewMockGenerator(mockCtrl)
	g := New(p)

	p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(trialJobTemplate, nil)

	instance := newFakeInstance()
	if err := g.(*DefaultValidator).validateTrialTemplate(instance); err != nil {
		t.Errorf("Expected nil, got err %v", err)
	}
}

func TestValidateExperiment(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()

	p := manifestmock.NewMockGenerator(mockCtrl)
	g := New(p)

	trialJobTemplate := `apiVersion: "batch/v1"
kind: "Job"
metadata:
  name: "fake-trial"
  namespace: fakens`

	suggestionConfigData := map[string]string{}
	suggestionConfigData[consts.LabelSuggestionImageTag] = "algorithmImage"

	p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(trialJobTemplate, nil).AnyTimes()
	p.EXPECT().GetSuggestionConfigData(gomock.Any()).Return(suggestionConfigData, nil).AnyTimes()
	p.EXPECT().GetMetricsCollectorImage(gomock.Any()).Return("metricsCollectorImage", nil).AnyTimes()

	tcs := []struct {
		Instance *experimentsv1alpha3.Experiment
		Err      bool
	}{
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.Objective = nil
				return i
			}(),
			Err: true,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.Algorithm = nil
				return i
			}(),
			Err: true,
		},
		{
			Instance: newFakeInstance(),
			Err:      false,
		},
	}

	for _, tc := range tcs {
		err := g.ValidateExperiment(tc.Instance, nil)
		if !tc.Err && err != nil {
			t.Errorf("Expected nil, got %v", err)
		} else if tc.Err && err == nil {
			t.Errorf("Expected err, got nil")
		}
	}
}

func newFakeInstance() *experimentsv1alpha3.Experiment {
	goal := 0.11
	return &experimentsv1alpha3.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "fakens",
		},
		Spec: experimentsv1alpha3.ExperimentSpec{
			MetricsCollectorSpec: &commonv1alpha3.MetricsCollectorSpec{
				Collector: &commonv1alpha3.CollectorSpec{
					Kind: commonv1alpha3.StdOutCollector,
				},
			},
			Objective: &commonv1alpha3.ObjectiveSpec{
				Type:                commonv1alpha3.ObjectiveTypeMaximize,
				Goal:                &goal,
				ObjectiveMetricName: "testme",
			},
			Algorithm: &commonv1alpha3.AlgorithmSpec{
				AlgorithmName: "test",
				AlgorithmSettings: []commonv1alpha3.AlgorithmSetting{
					{
						Name:  "test1",
						Value: "value1",
					},
				},
			},
			Parameters: []experimentsv1alpha3.ParameterSpec{
				{
					Name:          "test",
					ParameterType: experimentsv1alpha3.ParameterTypeCategorical,
					FeasibleSpace: experimentsv1alpha3.FeasibleSpace{
						List: []string{"1", "2"},
					},
				},
			},
		},
	}
}
