package validator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/manifest"
	jobv1beta1 "github.com/kubeflow/katib/pkg/job/v1beta1"
	mccommon "github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
)

var log = logf.Log.WithName("experiment-validating-webhook")

type Validator interface {
	ValidateExperiment(instance, oldInst *experimentsv1beta1.Experiment) error
	InjectClient(c client.Client)
}

type DefaultValidator struct {
	manifest.Generator
}

func New(generator manifest.Generator) Validator {
	return &DefaultValidator{
		Generator: generator,
	}
}

func (g *DefaultValidator) InjectClient(c client.Client) {
	g.Generator.InjectClient(c)
}

func (g *DefaultValidator) ValidateExperiment(instance, oldInst *experimentsv1beta1.Experiment) error {
	if instance.Spec.MaxFailedTrialCount != nil && *instance.Spec.MaxFailedTrialCount < 0 {
		return fmt.Errorf("spec.maxFailedTrialCount should not be less than 0")
	}
	if instance.Spec.MaxTrialCount != nil && *instance.Spec.MaxTrialCount <= 0 {
		return fmt.Errorf("spec.maxTrialCount must be greater than 0")
	}
	if instance.Spec.ParallelTrialCount != nil && *instance.Spec.ParallelTrialCount <= 0 {
		return fmt.Errorf("spec.parallelTrialCount must be greater than 0")
	}
	if oldInst != nil {
		oldInst.Spec.MaxFailedTrialCount = instance.Spec.MaxFailedTrialCount
		oldInst.Spec.MaxTrialCount = instance.Spec.MaxTrialCount
		oldInst.Spec.ParallelTrialCount = instance.Spec.ParallelTrialCount
		if equality.Semantic.DeepEqual(instance.Spec, oldInst.Spec) {
			return nil
		} else {
			return fmt.Errorf("Only spec.parallelTrialCount, spec.maxTrialCount and spec.maxFailedTrialCount are editable.")
		}
	}
	if err := g.validateObjective(instance.Spec.Objective); err != nil {
		return err
	}
	if err := g.validateAlgorithm(instance.Spec.Algorithm); err != nil {
		return err
	}
	if err := g.validateResumePolicy(instance.Spec.ResumePolicy); err != nil {
		return err
	}

	if err := g.validateTrialTemplate(instance); err != nil {
		return err
	}

	if len(instance.Spec.Parameters) == 0 && instance.Spec.NasConfig == nil {
		return fmt.Errorf("spec.parameters or spec.nasConfig must be specified.")
	}

	if len(instance.Spec.Parameters) > 0 && instance.Spec.NasConfig != nil {
		return fmt.Errorf("Only one of spec.parameters and spec.nasConfig can be specified.")
	}

	if err := g.validateMetricsCollector(instance); err != nil {
		return err
	}
	return nil
}

func (g *DefaultValidator) validateObjective(obj *commonapiv1beta1.ObjectiveSpec) error {
	if obj == nil {
		return fmt.Errorf("No spec.objective specified.")
	}
	if obj.Type != commonapiv1beta1.ObjectiveTypeMinimize && obj.Type != commonapiv1beta1.ObjectiveTypeMaximize {
		return fmt.Errorf("spec.objective.type must be %s or %s.", commonapiv1beta1.ObjectiveTypeMinimize, commonapiv1beta1.ObjectiveTypeMaximize)
	}
	if obj.ObjectiveMetricName == "" {
		return fmt.Errorf("No spec.objective.objectiveMetricName specified.")
	}
	return nil
}

func (g *DefaultValidator) validateAlgorithm(ag *commonapiv1beta1.AlgorithmSpec) error {
	if ag == nil {
		return fmt.Errorf("No spec.algorithm specified.")
	}
	if ag.AlgorithmName == "" {
		return fmt.Errorf("No spec.algorithm.name specified.")
	}

	if _, err := g.GetSuggestionConfigData(ag.AlgorithmName); err != nil {
		return fmt.Errorf("Don't support algorithm %s: %v.", ag.AlgorithmName, err)
	}

	return nil
}

func (g *DefaultValidator) validateResumePolicy(resume experimentsv1beta1.ResumePolicyType) error {
	validTypes := map[experimentsv1beta1.ResumePolicyType]string{
		"":                             "",
		experimentsv1beta1.NeverResume: "",
		experimentsv1beta1.LongRunning: "",
	}
	if _, ok := validTypes[resume]; !ok {
		return fmt.Errorf("invalid ResumePolicyType %s", resume)
	}
	return nil
}

func (g *DefaultValidator) validateTrialTemplate(instance *experimentsv1beta1.Experiment) error {
	trialName := fmt.Sprintf("%s-trial", instance.GetName())
	runSpec, err := g.GetRunSpec(instance, instance.GetName(), trialName, instance.GetNamespace())
	if err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	bufSize := 1024
	buf := bytes.NewBufferString(runSpec)

	job := &unstructured.Unstructured{}
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, bufSize).Decode(job); err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	if err := g.validateSupportedJob(job); err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	if job.GetNamespace() != instance.GetNamespace() {
		return fmt.Errorf("Invalid spec.trialTemplate: metadata.namespace should be %s or {{.NameSpace}}", instance.GetNamespace())
	}
	if job.GetName() != trialName {
		return fmt.Errorf("Invalid spec.trialTemplate: metadata.name should be {{.Trial}}")
	}
	return nil
}

func (g *DefaultValidator) validateSupportedJob(job *unstructured.Unstructured) error {
	gvk := job.GroupVersionKind()
	supportedJobs := jobv1beta1.SupportedJobList
	for _, sJob := range supportedJobs {
		if gvk == sJob {
			return nil
		}
	}
	return fmt.Errorf("Job type %v not supported", gvk)
}

func (g *DefaultValidator) validateMetricsCollector(inst *experimentsv1beta1.Experiment) error {
	mcSpec := inst.Spec.MetricsCollectorSpec
	mcKind := mcSpec.Collector.Kind
	for _, mc := range mccommon.AutoInjectMetricsCollecterList {
		if mcKind != mc {
			continue
		}
		if _, err := g.GetMetricsCollectorImage(mcKind); err != nil {
			return fmt.Errorf("GetMetricsCollectorImage failed: %v.", err)
		}
		break
	}
	// TODO(hougangliu): log warning message if some field will not be used for the metricsCollector kind
	switch mcKind {
	case commonapiv1beta1.NoneCollector, commonapiv1beta1.StdOutCollector:
		return nil
	case commonapiv1beta1.FileCollector:
		if mcSpec.Source == nil || mcSpec.Source.FileSystemPath == nil ||
			mcSpec.Source.FileSystemPath.Kind != commonapiv1beta1.FileKind || !filepath.IsAbs(mcSpec.Source.FileSystemPath.Path) {
			return fmt.Errorf("File path where metrics file exists is required by .spec.metricsCollectorSpec.source.fileSystemPath.path")
		}
	case commonapiv1beta1.TfEventCollector:
		if mcSpec.Source == nil || mcSpec.Source.FileSystemPath == nil ||
			mcSpec.Source.FileSystemPath.Kind != commonapiv1beta1.DirectoryKind || !filepath.IsAbs(mcSpec.Source.FileSystemPath.Path) {
			return fmt.Errorf("Directory path where tensorflow event files exist is required by .spec.metricsCollectorSpec.source.fileSystemPath.path")
		}
	case commonapiv1beta1.PrometheusMetricCollector:
		i, err := strconv.Atoi(mcSpec.Source.HttpGet.Port.String())
		if err != nil || i <= 0 {
			return fmt.Errorf(".spec.metricsCollectorSpec.source.httpGet.port must be a positive integer value for metrics collector kind: %v.", mcKind)
		}
		if !strings.HasPrefix(mcSpec.Source.HttpGet.Path, "/") {
			return fmt.Errorf(".spec.metricsCollectorSpec.source.httpGet.path is invalid for metrics collector kind: %v.", mcKind)
		}
	case commonapiv1beta1.CustomCollector:
		if mcSpec.Collector.CustomCollector == nil {
			return fmt.Errorf(".spec.metricsCollectorSpec.collector.customCollector is required for metrics collector kind: %v.", mcKind)
		}
		if mcSpec.Source.FileSystemPath != nil {
			if !filepath.IsAbs(mcSpec.Source.FileSystemPath.Path) || (mcSpec.Source.FileSystemPath.Kind != commonapiv1beta1.DirectoryKind &&
				mcSpec.Source.FileSystemPath.Kind != commonapiv1beta1.FileKind) {
				return fmt.Errorf(".spec.metricsCollectorSpec.source is invalid")
			}
		}
	default:
		return fmt.Errorf("Invalid metrics collector kind: %v.", mcKind)
	}
	if mcSpec.Source != nil && mcSpec.Source.Filter != nil && len(mcSpec.Source.Filter.MetricsFormat) > 0 {
		// the filter regular expression must have two top subexpressions, the first matched one will be taken as metric name, the second one as metric value
		mustTwoBracket, _ := regexp.Compile(`.*\(.*\).*\(.*\).*`)
		for _, mFormat := range mcSpec.Source.Filter.MetricsFormat {
			if _, err := regexp.Compile(mFormat); err != nil {
				return fmt.Errorf("Invalid %q in .spec.metricsCollectorSpec.source.filter: %v.", mFormat, err)
			} else {
				if !mustTwoBracket.MatchString(mFormat) {
					return fmt.Errorf("Invalid %q in .spec.metricsCollectorSpec.source.filter: two top subexpressions are required", mFormat)
				}
			}
		}
	}

	return nil
}
