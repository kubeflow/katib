package manifest

import (
	"bytes"
	"errors"
	"text/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibclient"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibconfig"
)

const (
	defaultMetricsCollectorTemplateName = "defaultMetricsCollectorTemplate.yaml"
)

// Generator is the type for manifests Generator.
type Generator interface {
	InjectClient(c client.Client)
	GetRunSpec(e *experimentsv1beta1.Experiment, experiment, trial, namespace string) (string, error)
	GetRunSpecWithHyperParameters(e *experimentsv1beta1.Experiment, experiment, trial, namespace string, hps []commonapiv1beta1.ParameterAssignment) (string, error)
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

// GetRunSpec get the specification for trial.
func (g *DefaultGenerator) GetRunSpec(e *experimentsv1beta1.Experiment, experiment, trial, namespace string) (string, error) {
	params := trialTemplateParams{
		Experiment: experiment,
		Trial:      trial,
		NameSpace:  namespace,
	}
	return g.getRunSpec(e, params)
}

// GetRunSpecWithHyperParameters get the specification for trial with hyperparameters.
func (g *DefaultGenerator) GetRunSpecWithHyperParameters(e *experimentsv1beta1.Experiment, experiment, trial, namespace string, hps []commonapiv1beta1.ParameterAssignment) (string, error) {
	params := trialTemplateParams{
		Experiment:      experiment,
		Trial:           trial,
		NameSpace:       namespace,
		HyperParameters: hps,
	}
	return g.getRunSpec(e, params)
}

func (g *DefaultGenerator) getRunSpec(e *experimentsv1beta1.Experiment, params trialTemplateParams) (string, error) {
	var buf bytes.Buffer
	trialTemplate, err := g.getTrialTemplate(e)
	if err != nil {
		return "", err
	}
	if err := trialTemplate.Execute(&buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g *DefaultGenerator) getTrialTemplate(instance *experimentsv1beta1.Experiment) (*template.Template, error) {
	var err error
	var tpl *template.Template = nil

	trialTemplate := instance.Spec.TrialTemplate
	if trialTemplate.GoTemplate.RawTemplate != "" {
		tpl, err = template.New("Trial").Parse(trialTemplate.GoTemplate.RawTemplate)
		if err != nil {
			return nil, err
		}
	} else {
		templateSpec := trialTemplate.GoTemplate.TemplateSpec
		configMapNS := templateSpec.ConfigMapNamespace
		configMapName := templateSpec.ConfigMapName
		templatePath := templateSpec.TemplatePath

		configMap, err := g.client.GetConfigMap(configMapName, configMapNS)
		if err != nil {
			return nil, err
		}

		if configMapTemplate, ok := configMap[templatePath]; !ok {
			err = errors.New(string(metav1.StatusReasonNotFound))
			return nil, err
		} else {
			tpl, err = template.New("Trial").Parse(configMapTemplate)
			if err != nil {
				return nil, err
			}
		}
	}

	return tpl, nil
}

type trialTemplateParams struct {
	Experiment      string
	Trial           string
	NameSpace       string
	HyperParameters []commonapiv1beta1.ParameterAssignment
}
