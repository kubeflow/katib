package validator

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	manifestmock "github.com/kubeflow/katib/pkg/mock/v1alpha3/experiment/manifest"
	v1 "k8s.io/api/core/v1"
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

	p := manifestmock.NewMockGenerator(mockCtrl)
	g := New(p)

	p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(trialTFJobTemplate, nil)

	instance := newFakeInstance()
	if err := g.(*DefaultValidator).validateTrialTemplate(instance); err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestValidateJobTrialTemplate(t *testing.T) {
	trialJobTemplate := `apiVersion: batch/v1
kind: Job
metadata:
  name: fake-trial
  namespace: fakens
spec:
  template:
    spec:
      containers:
      - name: fake-trial
        image: test-image`

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	p := manifestmock.NewMockGenerator(mockCtrl)
	g := New(p)

	invalidYaml := strings.Replace(trialJobTemplate, "- name", "- * -", -1)
	invalidJobType := strings.Replace(trialJobTemplate, "Job", "NewJobType", -1)
	invalidNamespace := strings.Replace(trialJobTemplate, "fakens", "not-fakens", -1)
	invalidJobName := strings.Replace(trialJobTemplate, "fake-trial", "new-name", -1)

	validRun := p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(trialJobTemplate, nil)
	invalidYamlRun := p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(invalidYaml, nil)
	invalidJobTypeRun := p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(invalidJobType, nil)
	invalidNamespaceRun := p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(invalidNamespace, nil)
	invalidJobNameRun := p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(invalidJobName, nil)

	gomock.InOrder(
		validRun,
		invalidYamlRun,
		invalidJobTypeRun,
		invalidNamespaceRun,
		invalidJobNameRun,
	)

	tcs := []struct {
		Instance *experimentsv1alpha3.Experiment
		Err      bool
	}{
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate = newFakeTrialTemplate(trialJobTemplate)
				return i
			}(),
			Err: false,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate = newFakeTrialTemplate(invalidYaml)
				return i
			}(),
			Err: true,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate = newFakeTrialTemplate(invalidJobType)
				return i
			}(),
			Err: true,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate = newFakeTrialTemplate(invalidNamespace)
				return i
			}(),
			Err: true,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.TrialTemplate = newFakeTrialTemplate(invalidJobName)
				return i
			}(),
			Err: true,
		},
	}
	for _, tc := range tcs {
		err := g.(*DefaultValidator).validateTrialTemplate(tc.Instance)
		if !tc.Err && err != nil {
			t.Errorf("Expected nil, got %v", err)
		} else if tc.Err && err == nil {
			t.Errorf("Expected err, got nil")
		}
	}
}

func TestValidateExperiment(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	p := manifestmock.NewMockGenerator(mockCtrl)
	g := New(p)

	trialJobTemplate := `apiVersion: "batch/v1"
kind: "Job"
metadata:
  name: "fake-trial"
  namespace: fakens`

	suggestionConfigData := map[string]string{}
	suggestionConfigData[consts.LabelSuggestionImageTag] = "algorithmImage"
	fakeNegativeInt := int32(-1)

	p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(trialJobTemplate, nil).AnyTimes()
	p.EXPECT().GetSuggestionConfigData(gomock.Any()).Return(suggestionConfigData, nil).AnyTimes()
	p.EXPECT().GetMetricsCollectorImage(gomock.Any()).Return("metricsCollectorImage", nil).AnyTimes()

	tcs := []struct {
		Instance    *experimentsv1alpha3.Experiment
		Err         bool
		oldInstance *experimentsv1alpha3.Experiment
	}{
		//Objective
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
				i.Spec.Objective.Type = commonv1alpha3.ObjectiveTypeUnknown
				return i
			}(),
			Err: true,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.Objective.ObjectiveMetricName = ""
				return i
			}(),
			Err: true,
		},
		//Algorithm
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.Algorithm = nil
				return i
			}(),
			Err: true,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.Algorithm.AlgorithmName = ""
				return i
			}(),
			Err: true,
		},
		{
			Instance: newFakeInstance(),
			Err:      false,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MaxFailedTrialCount = &fakeNegativeInt
				return i
			}(),
			Err: true,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MaxTrialCount = &fakeNegativeInt
				return i
			}(),
			Err: true,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.ParallelTrialCount = &fakeNegativeInt
				return i
			}(),
			Err: true,
		},
		{
			Instance:    newFakeInstance(),
			Err:         false,
			oldInstance: newFakeInstance(),
		},
		{
			Instance: newFakeInstance(),
			Err:      true,
			oldInstance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.Algorithm.AlgorithmName = "not-test"
				return i
			}(),
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.Parameters = []experimentsv1alpha3.ParameterSpec{}
				i.Spec.NasConfig = nil
				return i
			}(),
			Err: true,
		},
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.NasConfig = &experimentsv1alpha3.NasConfig{
					Operations: []experimentsv1alpha3.Operation{
						{
							OperationType: "op1",
						},
					},
				}
				return i
			}(),
			Err: true,
		},
	}

	for _, tc := range tcs {
		err := g.ValidateExperiment(tc.Instance, tc.oldInstance)
		if !tc.Err && err != nil {
			t.Errorf("Expected nil, got %v", err)
		} else if tc.Err && err == nil {
			t.Errorf("Expected err, got nil")
		}
	}
}

func TestValidateMetricsCollector(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	p := manifestmock.NewMockGenerator(mockCtrl)
	g := New(p)

	p.EXPECT().GetMetricsCollectorImage(gomock.Any()).Return("metricsCollectorImage", nil).AnyTimes()

	tcs := []struct {
		Instance *experimentsv1alpha3.Experiment
		Err      bool
	}{
		// Invalid Metrics Collector Kind
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.CollectorKind("invalid-kind"),
					},
				}
				return i
			}(),
			Err: true,
		},
		// FileCollector invalid Path
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.FileCollector,
					},
					Source: &commonv1alpha3.SourceSpec{
						FileSystemPath: &commonv1alpha3.FileSystemPath{
							Path: "not/absolute/path",
						},
					},
				}
				return i
			}(),
			Err: true,
		},
		// TfEventCollector invalid Path
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.TfEventCollector,
					},
					Source: &commonv1alpha3.SourceSpec{
						FileSystemPath: &commonv1alpha3.FileSystemPath{
							Path: "not/absolute/path",
						},
					},
				}
				return i
			}(),
			Err: true,
		},
		// PrometheusMetricCollector invalid Port
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.PrometheusMetricCollector,
					},
					Source: &commonv1alpha3.SourceSpec{
						HttpGet: &v1.HTTPGetAction{
							Port: intstr.IntOrString{
								StrVal: "Port",
							},
						},
					},
				}
				return i
			}(),
			Err: true,
		},
		// PrometheusMetricCollector invalid Path
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.PrometheusMetricCollector,
					},
					Source: &commonv1alpha3.SourceSpec{
						HttpGet: &v1.HTTPGetAction{
							Port: intstr.IntOrString{
								IntVal: 8888,
							},
							Path: "not/valid/path",
						},
					},
				}
				return i
			}(),
			Err: true,
		},
		//  CustomCollector empty CustomCollector
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.CustomCollector,
					},
				}
				return i
			}(),
			Err: true,
		},
		//  CustomCollector invalid Path
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.CustomCollector,
						CustomCollector: &v1.Container{
							Name: "my-collector",
						},
					},
					Source: &commonv1alpha3.SourceSpec{
						FileSystemPath: &commonv1alpha3.FileSystemPath{
							Path: "not/absolute/path",
						},
					},
				}
				return i
			}(),
			Err: true,
		},
		// FileMetricCollector invalid regexp in metrics format
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.FileCollector,
					},
					Source: &commonv1alpha3.SourceSpec{
						Filter: &commonv1alpha3.FilterSpec{
							MetricsFormat: []string{
								"[",
							},
						},
						FileSystemPath: &commonv1alpha3.FileSystemPath{
							Path: "/absolute/path",
							Kind: commonv1alpha3.FileKind,
						},
					},
				}
				return i
			}(),
			Err: true,
		},
		// FileMetricCollector one subexpression in metrics format
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.FileCollector,
					},
					Source: &commonv1alpha3.SourceSpec{
						Filter: &commonv1alpha3.FilterSpec{
							MetricsFormat: []string{
								"{metricName: ([\\w|-]+)}",
							},
						},
						FileSystemPath: &commonv1alpha3.FileSystemPath{
							Path: "/absolute/path",
							Kind: commonv1alpha3.FileKind,
						},
					},
				}
				return i
			}(),
			Err: true,
		},
		// Valid FileMetricCollector
		{
			Instance: func() *experimentsv1alpha3.Experiment {
				i := newFakeInstance()
				i.Spec.MetricsCollectorSpec = &commonv1alpha3.MetricsCollectorSpec{
					Collector: &commonv1alpha3.CollectorSpec{
						Kind: commonv1alpha3.FileCollector,
					},
					Source: &commonv1alpha3.SourceSpec{
						FileSystemPath: &commonv1alpha3.FileSystemPath{
							Path: "/absolute/path",
							Kind: commonv1alpha3.FileKind,
						},
					},
				}
				return i
			}(),
			Err: false,
		},
	}

	for _, tc := range tcs {
		err := g.(*DefaultValidator).validateMetricsCollector(tc.Instance)
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

func newFakeTrialTemplate(template string) *experimentsv1alpha3.TrialTemplate {
	return &experimentsv1alpha3.TrialTemplate{
		Retain: false,
		GoTemplate: &experimentsv1alpha3.GoTemplate{
			RawTemplate: template,
		},
	}
}
