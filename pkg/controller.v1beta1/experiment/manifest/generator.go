package manifest

import (
	"errors"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	util "github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibclient"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibconfig"
)

const (
	defaultMetricsCollectorTemplateName = "defaultMetricsCollectorTemplate.yaml"
)

// Generator is the type for manifests Generator.
type Generator interface {
	InjectClient(c client.Client)
	// TODO (andreyvelich): Add this after changing validation for new Trial Template
	// GetRunSpec(e *experimentsv1beta1.Experiment, experiment, trial, namespace string) (string, error)
	GetRunSpecWithHyperParameters(experiment *experimentsv1beta1.Experiment, trialName, trialNamespace string, assignments []commonapiv1beta1.ParameterAssignment) (*unstructured.Unstructured, error)
	GetSuggestionConfigData(algorithmName string) (map[string]string, error)
	GetMetricsCollectorImage(cKind commonapiv1beta1.CollectorKind) (string, error)
}

// DefaultGenerator is the default implementation of Generator.
type DefaultGenerator struct {
	client katibclient.Client
}

// New creates a new Generator.
func New(c client.Client) Generator {
	katibClient := katibclient.NewWithGivenClient(c)
	return &DefaultGenerator{
		client: katibClient,
	}
}

func (g *DefaultGenerator) InjectClient(c client.Client) {
	g.client.InjectClient(c)
}

func (g *DefaultGenerator) GetMetricsCollectorImage(cKind commonapiv1beta1.CollectorKind) (string, error) {
	configData, err := katibconfig.GetMetricsCollectorConfigData(cKind, g.client.GetClient())
	if err != nil {
		return "", nil
	}
	return configData[consts.LabelMetricsCollectorSidecarImage], nil
}

func (g *DefaultGenerator) GetSuggestionConfigData(algorithmName string) (map[string]string, error) {
	return katibconfig.GetSuggestionConfigData(algorithmName, g.client.GetClient())
}

// GetRunSpecWithHyperParameters get the specification for trial with hyperparameters.
func (g *DefaultGenerator) GetRunSpecWithHyperParameters(experiment *experimentsv1beta1.Experiment, trialName, trialNamespace string, assignments []commonapiv1beta1.ParameterAssignment) (*unstructured.Unstructured, error) {

	// Get string Trial template from Experiment spec
	trialTemplate, err := g.getTrialTemplate(experiment)
	if err != nil {
		return nil, err
	}

	// Apply parameters to Trial Template from assignment
	replacedTemplate, err := g.applyParameters(trialTemplate, experiment.Spec.TrialTemplate.TrialParameters, assignments)
	if err != nil {
		return nil, err
	}
	// Convert Trial template to unstructured
	runSpec, err := util.ConvertStringToUnstructured(replacedTemplate)
	if err != nil {
		return nil, fmt.Errorf("ConvertStringToUnstructured failed: %v", err)
	}

	// Set name and namespace for Run Spec
	runSpec.SetName(trialName)
	runSpec.SetNamespace(trialNamespace)

	return runSpec, nil
}

func (g *DefaultGenerator) applyParameters(trialTemplate string, trialParams []experimentsv1beta1.TrialParameterSpec, assignments []commonapiv1beta1.ParameterAssignment) (string, error) {
	// Number of parameters must be equal
	if len(assignments) != len(trialParams) {
		return "", fmt.Errorf("Number of Trial assignment from Suggestion: %v not equal to number Trial parameters from Experiment: %v", len(assignments), len(trialParams))
	}
	// Convert parameter assignment to map key = parameter name, value = parameter value
	assignmentsMap := make(map[string]string)
	for _, assignment := range assignments {
		assignmentsMap[assignment.Name] = assignment.Value
	}

	// Replacing parameters from Trial parameters
	for _, parameter := range trialParams {

		if parameterValue, ok := assignmentsMap[parameter.Reference]; ok {
			trialTemplate = strings.Replace(trialTemplate, fmt.Sprintf("${trialParameters.%v}", parameter.Name), parameterValue, -1)
		} else {
			return "", fmt.Errorf("Unable to find parameter: %v in parameter assignment %v", parameter.Reference, assignmentsMap)
		}
	}
	return trialTemplate, nil
}

func (g *DefaultGenerator) getTrialTemplate(instance *experimentsv1beta1.Experiment) (string, error) {
	var trialTemplateString string
	var err error

	trialSource := instance.Spec.TrialTemplate.TrialSource
	if trialSource.TrialSpec != nil {
		trialTemplateString, err = util.ConvertUnstructuredToString(trialSource.TrialSpec)
		if err != nil {
			return "", fmt.Errorf("ConvertUnstructuredToString failed: %v", err)
		}
	} else {
		configMapNS := trialSource.ConfigMap.ConfigMapNamespace
		configMapName := trialSource.ConfigMap.ConfigMapName
		templatePath := trialSource.ConfigMap.TemplatePath
		configMap, err := g.client.GetConfigMap(configMapName, configMapNS)
		if err != nil {
			return "", err
		}
		var ok bool
		trialTemplateString, ok = configMap[templatePath]
		if !ok {
			err = errors.New(string(metav1.StatusReasonNotFound))
			return "", err
		}
	}

	return trialTemplateString, nil
}
