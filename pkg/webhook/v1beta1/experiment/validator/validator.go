/*
Copyright 2022 The Kubeflow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validator

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	jsonPatch "github.com/mattbaird/jsonpatch"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/manifest"
	experimentutil "github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/util"
	util "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
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

// ValidateExperiment validates experiment for the given instance.
// oldInst is specified when experiment is edited.
func (g *DefaultValidator) ValidateExperiment(instance, oldInst *experimentsv1beta1.Experiment) error {
	namingConvention, _ := regexp.Compile("^[a-z]([-a-z0-9]*[a-z0-9])?")
	if !namingConvention.MatchString(instance.Name) {
		msg := "name must consist of lower case alphanumeric characters or '-'," +
			" start with an alphabetic character, and end with an alphanumeric character" +
			" (e.g. 'my-name', or 'abc-123', regex used for validation is '^[a-z]([-a-z0-9]*[a-z0-9])?)'"

		return fmt.Errorf(msg)
	}

	if instance.Spec.MaxFailedTrialCount != nil && *instance.Spec.MaxFailedTrialCount < 0 {
		return fmt.Errorf("spec.maxFailedTrialCount should not be less than 0")
	}
	if instance.Spec.MaxTrialCount != nil && *instance.Spec.MaxTrialCount <= 0 {
		return fmt.Errorf("spec.maxTrialCount must be greater than 0")
	}
	if instance.Spec.ParallelTrialCount != nil && *instance.Spec.ParallelTrialCount <= 0 {
		return fmt.Errorf("spec.parallelTrialCount must be greater than 0")
	}

	if instance.Spec.MaxFailedTrialCount != nil && instance.Spec.MaxTrialCount != nil {
		if *instance.Spec.MaxFailedTrialCount > *instance.Spec.MaxTrialCount {
			return fmt.Errorf("spec.maxFailedTrialCount should be less than or equal to spec.maxTrialCount")
		}
	}
	if instance.Spec.ParallelTrialCount != nil && instance.Spec.MaxTrialCount != nil {
		if *instance.Spec.ParallelTrialCount > *instance.Spec.MaxTrialCount {
			return fmt.Errorf("spec.paralelTrialCount should be less than or equal to spec.maxTrialCount")
		}
	}

	if oldInst != nil {
		// We should validate restart only if appropriate fields are changed.
		// Otherwise check below is triggered when experiment is deleted.
		isRestarting := false
		if !equality.Semantic.DeepEqual(instance.Spec, oldInst.Spec) {
			isRestarting = true
		}

		// When experiment is completed IsCompletedExperimentRestartable must return true
		if isRestarting && oldInst.IsCompleted() && !experimentutil.IsCompletedExperimentRestartable(oldInst) {
			msg := fmt.Sprintf("Experiment can be restarted if it is in succeeded state by reaching max trials and "+
				"spec.resumePolicy = %v or spec.resumePolicy = %v, when experiment is completed",
				experimentsv1beta1.LongRunning, experimentsv1beta1.FromVolume)
			return fmt.Errorf(msg)
		}

		if isRestarting && instance.Spec.MaxTrialCount != nil && *instance.Spec.MaxTrialCount <= oldInst.Status.Trials {
			return fmt.Errorf("spec.maxTrialCount: %v must be greater than status.trials count: %v",
				*instance.Spec.MaxTrialCount, oldInst.Status.Trials)
		}
		oldInst.Spec.MaxFailedTrialCount = instance.Spec.MaxFailedTrialCount
		oldInst.Spec.MaxTrialCount = instance.Spec.MaxTrialCount
		oldInst.Spec.ParallelTrialCount = instance.Spec.ParallelTrialCount
		if !equality.Semantic.DeepEqual(instance.Spec, oldInst.Spec) {
			return fmt.Errorf("only spec.parallelTrialCount, spec.maxTrialCount and spec.maxFailedTrialCount are editable")
		}
	}
	if err := g.validateObjective(instance.Spec.Objective); err != nil {
		return err
	}
	if err := g.validateAlgorithm(instance.Spec.Algorithm); err != nil {
		return err
	}
	if err := g.validateEarlyStopping(instance.Spec.EarlyStopping); err != nil {
		return err
	}
	if err := g.validateResumePolicy(instance.Spec.ResumePolicy); err != nil {
		return err
	}

	if err := g.validateTrialTemplate(instance); err != nil {
		return err
	}

	if len(instance.Spec.Parameters) == 0 && instance.Spec.NasConfig == nil {
		return fmt.Errorf("spec.parameters or spec.nasConfig must be specified")
	}

	if len(instance.Spec.Parameters) > 0 && instance.Spec.NasConfig != nil {
		return fmt.Errorf("only one of spec.parameters and spec.nasConfig can be specified")
	}

	if len(instance.Spec.Parameters) > 0 {
		if err := g.validateParameters(instance.Spec.Parameters); err != nil {
			return err
		}
	}

	if err := g.validateMetricsCollector(instance); err != nil {
		return err
	}
	return nil
}

func (g *DefaultValidator) validateObjective(obj *commonapiv1beta1.ObjectiveSpec) error {
	if obj == nil {
		return fmt.Errorf("no spec.objective specified")
	}
	if obj.Type != commonapiv1beta1.ObjectiveTypeMinimize && obj.Type != commonapiv1beta1.ObjectiveTypeMaximize {
		return fmt.Errorf("spec.objective.type must be %s or %s", commonapiv1beta1.ObjectiveTypeMinimize, commonapiv1beta1.ObjectiveTypeMaximize)
	}
	if obj.ObjectiveMetricName == "" {
		return fmt.Errorf("no spec.objective.objectiveMetricName specified")
	}
	if contains(obj.AdditionalMetricNames, obj.ObjectiveMetricName) {
		return fmt.Errorf("spec.objective.additionalMetricNames should not contain spec.objective.objectiveMetricName")
	}
	return nil
}

func (g *DefaultValidator) validateAlgorithm(ag *commonapiv1beta1.AlgorithmSpec) error {
	if ag == nil {
		return fmt.Errorf("no spec.algorithm specified")
	}
	if ag.AlgorithmName == "" {
		return fmt.Errorf("no spec.algorithm.name specified")
	}

	if _, err := g.GetSuggestionConfigData(ag.AlgorithmName); err != nil {
		return fmt.Errorf("unable to get Suggestion config data for algorithm %s: %v", ag.AlgorithmName, err)
	}

	return nil
}

func (g *DefaultValidator) validateEarlyStopping(es *commonapiv1beta1.EarlyStoppingSpec) error {
	if es == nil {
		return nil
	}
	if es.AlgorithmName == "" {
		return fmt.Errorf("no spec.earlyStopping.algorithmName specified")
	}

	if _, err := g.GetEarlyStoppingConfigData(es.AlgorithmName); err != nil {
		return fmt.Errorf("unable to get EarlyStopping config data for algorithm %s: %v", es.AlgorithmName, err)
	}

	return nil
}

func (g *DefaultValidator) validateResumePolicy(resume experimentsv1beta1.ResumePolicyType) error {
	validTypes := map[experimentsv1beta1.ResumePolicyType]string{
		"":                             "",
		experimentsv1beta1.NeverResume: "",
		experimentsv1beta1.LongRunning: "",
		experimentsv1beta1.FromVolume:  "",
	}
	if _, ok := validTypes[resume]; !ok {
		return fmt.Errorf("invalid ResumePolicyType %s", resume)
	}
	return nil
}

func (g *DefaultValidator) validateParameters(parameters []experimentsv1beta1.ParameterSpec) error {
	for i, param := range parameters {

		if param.ParameterType != experimentsv1beta1.ParameterTypeInt &&
			param.ParameterType != experimentsv1beta1.ParameterTypeDouble &&
			param.ParameterType != experimentsv1beta1.ParameterTypeCategorical &&
			param.ParameterType != experimentsv1beta1.ParameterTypeDiscrete &&
			param.ParameterType != experimentsv1beta1.ParameterTypeUnknown {
			return fmt.Errorf("parameterType: %v is not supported in spec.parameters[%v]: %v", param.ParameterType, i, param)
		}

		if equality.Semantic.DeepEqual(param.FeasibleSpace, experimentsv1beta1.FeasibleSpace{}) {
			return fmt.Errorf("feasibleSpace must be specified in spec.parameters[%v]: %v", i, param)
		}

		if param.ParameterType == experimentsv1beta1.ParameterTypeDouble || param.ParameterType == experimentsv1beta1.ParameterTypeInt {
			if len(param.FeasibleSpace.List) > 0 {
				return fmt.Errorf("feasibleSpace.list is not supported for parameterType: %v in spec.parameters[%v]: %v", param.ParameterType, i, param)
			}
			if param.FeasibleSpace.Max == "" && param.FeasibleSpace.Min == "" {
				return fmt.Errorf("feasibleSpace.max or feasibleSpace.min must be specified for parameterType: %v in spec.parameters[%v]: %v", param.ParameterType, i, param)
			}

		} else if param.ParameterType == experimentsv1beta1.ParameterTypeCategorical || param.ParameterType == experimentsv1beta1.ParameterTypeDiscrete {
			if param.FeasibleSpace.Max != "" || param.FeasibleSpace.Min != "" || param.FeasibleSpace.Step != "" {
				return fmt.Errorf("feasibleSpace .max, .min and .step is not supported for parameterType: %v in spec.parameters[%v]: %v", param.ParameterType, i, param)
			}
		}
	}

	return nil
}

func (g *DefaultValidator) validateTrialTemplate(instance *experimentsv1beta1.Experiment) error {

	trialTemplate := instance.Spec.TrialTemplate

	if trialTemplate == nil {
		return fmt.Errorf("spec.trialTemplate must be specified")
	}

	// Check if PrimaryContainerName is set
	if trialTemplate.PrimaryContainerName == "" {
		return fmt.Errorf("spec.trialTemplate.primaryContainerName must be specified")
	}

	// Check if SuccessCondition and FailureCondition is set
	if trialTemplate.SuccessCondition == "" || trialTemplate.FailureCondition == "" {
		return fmt.Errorf("spec.trialTemplate.successCondition and spec.trialTemplate.failureCondition must be specified")
	}

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
		return fmt.Errorf("only one of spec.trialTemplate.trialSpec or spec.trialTemplate.configMap can be specified")
	}

	// Check if configMap parameters are specified
	if trialTemplate.ConfigMap != nil &&
		(trialTemplate.TrialSource.ConfigMap.ConfigMapName == "" ||
			trialTemplate.TrialSource.ConfigMap.ConfigMapNamespace == "" ||
			trialTemplate.TrialSource.ConfigMap.TemplatePath == "") {
		return fmt.Errorf("for spec.trialTemplate.configMap .configMapName and .configMapNamespace and .templatePath must be specified")
	}

	// Check if Trial template can be parsed to string
	trialTemplateStr, err := g.GetTrialTemplate(instance)
	if err != nil {
		return fmt.Errorf("unable to parse spec.trialTemplate: %v", err)
	}

	experimentParameterNames := make(map[string]bool)
	for _, parameter := range instance.Spec.Parameters {
		experimentParameterNames[parameter.Name] = true
	}

	trialParametersNames := make(map[string]bool)
	trialParametersRefs := make(map[string]bool)

	for _, parameter := range trialTemplate.TrialParameters {
		// Check if all trialParameters contain name and reference. Or name contains invalid character
		if parameter.Name == "" || parameter.Reference == "" ||
			strings.Contains(parameter.Name, "{") || strings.Contains(parameter.Name, "}") {
			return fmt.Errorf("invalid spec.trialTemplate.trialParameters: %v", parameter)
		}

		// Check if parameter names are not duplicated
		if _, ok := trialParametersNames[parameter.Name]; ok {
			return fmt.Errorf("parameter name %v can't be duplicated in spec.trialTemplate.trialParameters: %v", parameter.Name, trialTemplate.TrialParameters)
		}
		// Check if parameter references are not duplicated
		if _, ok := trialParametersRefs[parameter.Reference]; ok {
			return fmt.Errorf("parameter reference %v can't be duplicated in spec.trialTemplate.trialParameters: %v", parameter.Reference, trialTemplate.TrialParameters)
		}
		trialParametersNames[parameter.Name] = true
		trialParametersRefs[parameter.Reference] = true

		// Check if parameter reference exist in experiment parameters
		if len(experimentParameterNames) > 0 {
			if !isMetaKey(parameter.Reference) {
				if _, ok := experimentParameterNames[parameter.Reference]; !ok {
					return fmt.Errorf("parameter reference %v does not exist in spec.parameters: %v", parameter.Reference, instance.Spec.Parameters)
				}
			}
		}

		// Check if trialParameters contains all substitution for Trial template
		if !strings.Contains(trialTemplateStr, fmt.Sprintf(consts.TrialTemplateParamReplaceFormat, parameter.Name)) {
			return fmt.Errorf("parameter name: %v in spec.trialParameters not found in spec.trialTemplate: %v", parameter.Name, trialTemplateStr)
		}

		trialTemplateStr = strings.Replace(trialTemplateStr, fmt.Sprintf(consts.TrialTemplateParamReplaceFormat, parameter.Name), "test-value", -1)
	}

	// Check if Trial template contains all substitution for trialParameters
	substitutionRegex := regexp.MustCompile(consts.TrialTemplateParamReplaceFormatRegex)
	notReplacedParams := substitutionRegex.FindAllString(trialTemplateStr, -1)
	if len(notReplacedParams) != 0 {
		return fmt.Errorf("parameters: %v in spec.trialTemplate not found in spec.trialParameters: %v", notReplacedParams, trialTemplate.TrialParameters)
	}

	// Check if Trial template can be converted to unstructured
	runSpec, err := util.ConvertStringToUnstructured(trialTemplateStr)
	if err != nil {
		return fmt.Errorf("unable to convert spec.trialTemplate: %v to unstructured", trialTemplateStr)
	}

	// Check if metadata.name and metatdata.namespace is omittied
	if runSpec.GetName() != "" || runSpec.GetNamespace() != "" {
		return fmt.Errorf("metadata.name and metadata.namespace in spec.trialTemplate must be omitted")
	}

	// Check if ApiVersion and Kind is specified
	if runSpec.GetAPIVersion() == "" || runSpec.GetKind() == "" {
		return fmt.Errorf("APIVersion and Kind in spec.trialTemplate must be specified")
	}

	// Check if Job can be converted to Batch Job
	// Other CRDs are not validated
	if err := g.validateTrialJob(runSpec); err != nil {
		return fmt.Errorf("invalid spec.trialTemplate: %v", err)
	}

	return nil
}

func (g *DefaultValidator) validateTrialJob(runSpec *unstructured.Unstructured) error {
	gvk := runSpec.GroupVersionKind()

	// Validate only Kubernetes Job
	if gvk.GroupVersion() != batchv1.SchemeGroupVersion || gvk.Kind != consts.JobKindJob {
		return nil
	}

	batchJob := batchv1.Job{}

	// Validate that RunSpec can be converted to Batch Job
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(runSpec.Object, &batchJob)
	if err != nil {
		return fmt.Errorf("unable to convert spec.TrialTemplate: %v to %v: %v", runSpec.Object, gvk.Kind, err)
	}

	// Try to patch runSpec to Batch Job
	// TODO (andreyvelich): Do we want to remove it completely ?
	return validatePatchJob(runSpec, batchJob, gvk.Kind)
}

func validatePatchJob(runSpec *unstructured.Unstructured, job interface{}, jobType string) error {

	// Not necessary to check error runSpec.Object must be valid JSON
	runSpecBefore, _ := json.Marshal(runSpec.Object)

	// Not necessary to check error job must be valid JSON
	runSpecAfter, _ := json.Marshal(job)

	// Create Patch on transformed Job (e.g: Job) using unstructured JSON
	runSpecPatchOperations, err := jsonPatch.CreatePatch(runSpecAfter, runSpecBefore)
	if err != nil {
		return fmt.Errorf("create patch error: %v", err)
	}

	for _, operation := range runSpecPatchOperations {
		// If operation != "remove" some values from trialTemplate were not converted
		// We can't validate /resources/limits/ because CRDs can have custom k8s resources using defice plugin
		// ref https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/
		if operation.Operation != "remove" && !strings.Contains(operation.Path, "/resources/limits/") && !strings.Contains(operation.Path, "/resources/requests/") {
			return fmt.Errorf("unable to convert: %v - %v to %v, converted template: %v", operation.Path, operation.Value, jobType, string(runSpecAfter))
		}
	}

	return nil
}

func (g *DefaultValidator) validateMetricsCollector(inst *experimentsv1beta1.Experiment) error {
	mcSpec := inst.Spec.MetricsCollectorSpec
	mcKind := mcSpec.Collector.Kind
	for _, mc := range mccommon.AutoInjectMetricsCollectorList {
		if mcKind != mc {
			continue
		}
		if _, err := g.GetMetricsCollectorConfigData(mcKind); err != nil {
			return fmt.Errorf("GetMetricsCollectorConfigData failed: %v", err)
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
			return fmt.Errorf("file path where metrics file exists is required by .spec.metricsCollectorSpec.source.fileSystemPath.path")
		}
		// Format
		fileFormat := mcSpec.Source.FileSystemPath.Format
		if fileFormat != commonapiv1beta1.TextFormat && fileFormat != commonapiv1beta1.JsonFormat {
			return fmt.Errorf("format of metrics file is required by .spec.metricsCollectorSpec.source.fileSystemPath.format")
		}
		if fileFormat == commonapiv1beta1.JsonFormat && mcSpec.Source.Filter != nil {
			return fmt.Errorf(".spec.metricsCollectorSpec.source.filter must be nil when format of metrics file is %v", commonapiv1beta1.JsonFormat)
		}
	case commonapiv1beta1.TfEventCollector:
		if mcSpec.Source == nil || mcSpec.Source.FileSystemPath == nil ||
			mcSpec.Source.FileSystemPath.Kind != commonapiv1beta1.DirectoryKind || !filepath.IsAbs(mcSpec.Source.FileSystemPath.Path) {
			return fmt.Errorf("directory path where tensorflow event files exist is required by .spec.metricsCollectorSpec.source.fileSystemPath.path")
		}
		if mcSpec.Source.FileSystemPath.Format != "" {
			return fmt.Errorf(".spec.metricsCollectorSpec.source.fileSystemPath.format must be empty")
		}
	case commonapiv1beta1.PrometheusMetricCollector:
		i, err := strconv.Atoi(mcSpec.Source.HttpGet.Port.String())
		if err != nil || i <= 0 {
			return fmt.Errorf(".spec.metricsCollectorSpec.source.httpGet.port must be a positive integer value for metrics collector kind: %v", mcKind)
		}
		if !strings.HasPrefix(mcSpec.Source.HttpGet.Path, "/") {
			return fmt.Errorf(".spec.metricsCollectorSpec.source.httpGet.path is invalid for metrics collector kind: %v", mcKind)
		}
	case commonapiv1beta1.CustomCollector:
		if mcSpec.Collector.CustomCollector == nil {
			return fmt.Errorf(".spec.metricsCollectorSpec.collector.customCollector is required for metrics collector kind: %v", mcKind)
		}
		if mcSpec.Source.FileSystemPath != nil {
			if !filepath.IsAbs(mcSpec.Source.FileSystemPath.Path) || (mcSpec.Source.FileSystemPath.Kind != commonapiv1beta1.DirectoryKind &&
				mcSpec.Source.FileSystemPath.Kind != commonapiv1beta1.FileKind) {
				return fmt.Errorf(".spec.metricsCollectorSpec.source is invalid")
			}
		}
	default:
		return fmt.Errorf("invalid metrics collector kind: %v", mcKind)
	}
	if mcSpec.Source != nil && mcSpec.Source.Filter != nil && len(mcSpec.Source.Filter.MetricsFormat) > 0 {
		// the filter regular expression must have two top subexpressions, the first matched one will be taken as metric name, the second one as metric value
		mustTwoBracket, _ := regexp.Compile(`.*\(.*\).*\(.*\).*`)
		for _, mFormat := range mcSpec.Source.Filter.MetricsFormat {
			if _, err := regexp.Compile(mFormat); err != nil {
				return fmt.Errorf("invalid %q in .spec.metricsCollectorSpec.source.filter: %v", mFormat, err)
			} else {
				if !mustTwoBracket.MatchString(mFormat) {
					return fmt.Errorf("invalid %q in .spec.metricsCollectorSpec.source.filter: two top subexpressions are required", mFormat)
				}
			}
		}
	}

	return nil
}

func isMetaKey(parameter string) bool {
	// Check if parameter is trial metadata reference as ${trailSpec.Name}, ${trialSpec.Labels[label]}, etc. used for substitution
	match := regexp.MustCompile(consts.TrialTemplateMetaReplaceFormatRegex).FindStringSubmatch(parameter)
	isMeta := false
	if len(match) > 0 {
		matchedKey := match[1]
		if contains(consts.TrialTemplateMetaKeys, matchedKey) {
			isMeta = true
		} else {
			// Check if it's Labels[label] or Annotations[annotation]
			subMatch := regexp.MustCompile(consts.TrialTemplateMetaParseFormatRegex).FindStringSubmatch(matchedKey)
			if len(subMatch) == 3 && contains(consts.TrialTemplateMetaKeys, subMatch[1]) {
				isMeta = true
			}
		}
	}
	return isMeta
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
