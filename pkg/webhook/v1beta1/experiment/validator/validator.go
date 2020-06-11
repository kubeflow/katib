package validator

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/manifest"
	util "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	jobv1beta1 "github.com/kubeflow/katib/pkg/job/v1beta1"
	mccommon "github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
	batchv1 "k8s.io/api/batch/v1"

	pytorchv1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
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

	trialTemplate := instance.Spec.TrialTemplate

	// Check if trialParameters exists
	if trialTemplate.TrialParameters == nil {
		return fmt.Errorf("spec.trialTemplate.trialParameters must be specified")
	}

	// Check if trialSpec or configMap exists
	if trialTemplate.TrialSource.TrialSpec == nil && trialTemplate.TrialSource.ConfigMap == nil {
		return fmt.Errorf("spec.trialTemplate.trialSpec or spec.trialTemplate.configMap must be specified")
	}

	// Check if trialSpec and configMap doesn't exist together
	if trialTemplate.TrialSource.TrialSpec != nil && trialTemplate.TrialSource.ConfigMap != nil {
		return fmt.Errorf("Only one of spec.trialTemplate.trialSpec or spec.trialTemplate.configMap can be specified")
	}

	// Check if configMap parameters are specified
	if trialTemplate.ConfigMap != nil &&
		(trialTemplate.TrialSource.ConfigMap.ConfigMapName == "" ||
			trialTemplate.TrialSource.ConfigMap.ConfigMapNamespace == "" ||
			trialTemplate.TrialSource.ConfigMap.TemplatePath == "") {
		return fmt.Errorf("For spec.trialTemplate.configMap .configMapName and .configMapNamespace and .templatePath must be specified")
	}

	// Check if Trial template can be parsed to string
	trialTemplateStr, err := g.GetTrialTemplate(instance)
	if err != nil {
		return fmt.Errorf("Unable to parse spec.trialTemplate: %v", err)
	}

	trialParametersNames := make(map[string]bool)
	trialParametersRefs := make(map[string]bool)

	for _, parameter := range trialTemplate.TrialParameters {
		// Check if all trialParameters contain name and reference. Or name contains invalid character
		if parameter.Name == "" || parameter.Reference == "" ||
			strings.Index(parameter.Name, "{") != -1 || strings.Index(parameter.Name, "}") != -1 {
			return fmt.Errorf("Invalid spec.trialTemplate.trialParameters: %v", parameter)
		}

		// Check if parameter names are not duplicated
		if _, ok := trialParametersNames[parameter.Name]; ok {
			return fmt.Errorf("Parameter name %v can't be duplicated in spec.trialTemplate.trialParameters: %v", parameter.Name, trialTemplate.TrialParameters)
		}
		// Check if parameter references are not duplicated
		if _, ok := trialParametersRefs[parameter.Reference]; ok {
			return fmt.Errorf("Parameter reference %v can't be duplicated in spec.trialTemplate.trialParameters: %v", parameter.Reference, trialTemplate.TrialParameters)
		}
		trialParametersNames[parameter.Name] = true
		trialParametersRefs[parameter.Reference] = true

		// Check if trialParameters contains all substitution for Trial template
		if strings.Index(trialTemplateStr, fmt.Sprintf(consts.TrialTemplateReplaceFormat, parameter.Name)) == -1 {
			return fmt.Errorf("Parameter name: %v in spec.trialParameters not found in spec.trialTemplate: %v", parameter.Name, trialTemplateStr)
		}

		trialTemplateStr = strings.Replace(trialTemplateStr, fmt.Sprintf(consts.TrialTemplateReplaceFormat, parameter.Name), "test-value", -1)
	}

	// Check if Trial template contains all substitution for trialParameters
	substitutionRegex := regexp.MustCompile(consts.TrialTemplateReplaceFormatRegex)
	notReplacedParams := substitutionRegex.FindAllString(trialTemplateStr, -1)
	if len(notReplacedParams) != 0 {
		return fmt.Errorf("Parameters: %v in spec.trialTemplate not found in spec.trialParameters: %v", notReplacedParams, trialTemplate.TrialParameters)
	}

	// Check if Trial template can be converted to unstructured
	runSpec, err := util.ConvertStringToUnstructured(trialTemplateStr)
	if err != nil {
		return fmt.Errorf("Unable to convert spec.trialTemplate: %v to unstructured", trialTemplateStr)
	}

	// Check if metadata.name and metatdata.namespace is omittied
	if runSpec.GetName() != "" || runSpec.GetNamespace() != "" {
		return fmt.Errorf("metadata.name and metadata.namespace in spec.trialTemplate must be omitted")
	}

	// Check if ApiVersion and Kind is specified
	if runSpec.GetAPIVersion() == "" || runSpec.GetKind() == "" {
		return fmt.Errorf("apiVersion and kind in spec.trialTemplate must be specified")
	}

	// Check if Job is supported
	// Check if Job can be converted to Batch Job/TFJob/PyTorchJob
	// Not default CRDs can be omitted later
	if err := g.validateSupportedJob(runSpec); err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v", err)
	}

	return nil
}

func (g *DefaultValidator) validateSupportedJob(runSpec *unstructured.Unstructured) error {
	gvk := runSpec.GroupVersionKind()
	supportedJobs := jobv1beta1.SupportedJobList
	for _, sJob := range supportedJobs {
		if gvk == sJob {
			switch gvk.Kind {
			case consts.JobKindJob:
				batchJob := &batchv1.Job{}
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(runSpec.Object, &batchJob)
				if err != nil {
					return fmt.Errorf("Unable to convert spec.TrialTemplate to BatchJob: %v", err)
				}
			case consts.JobKindTF:
				tfJob := &tfv1.TFJob{}
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(runSpec.Object, &tfJob)
				if err != nil {
					return fmt.Errorf("Unable to convert spec.TrialTemplate to TFJob: %v", err)
				}
			case consts.JobKindPyTorch:
				pytorchJob := &pytorchv1.PyTorchJob{}
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(runSpec.Object, &pytorchJob)
				if err != nil {
					return fmt.Errorf("Unable to convert spec.TrialTemplate to PyTorchJob: %v", err)
				}
			}
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
