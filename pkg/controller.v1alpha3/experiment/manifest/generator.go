package manifest

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	commonapiv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	commonv1alpha3 "github.com/kubeflow/katib/pkg/common/v1alpha3"
	"github.com/kubeflow/katib/pkg/util/v1alpha3/katibclient"
)

const (
	defaultMetricsCollectorTemplateName = "defaultMetricsCollectorTemplate.yaml"
)

// Generator is the type for manifests Generator.
type Generator interface {
	InjectClient(c client.Client)
	GetRunSpec(e *experimentsv1alpha3.Experiment, experiment, trial, namespace string) (string, error)
	GetRunSpecWithHyperParameters(e *experimentsv1alpha3.Experiment, experiment, trial, namespace string, hps []commonapiv1alpha3.ParameterAssignment) (string, error)
	GetMetricsCollectorManifest(experimentName string, trialName string, jobKind string, namespace string, metricNames []string, mcs *experimentsv1alpha3.MetricsCollectorSpec) (string, error)
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

func (g *DefaultGenerator) GetMetricsCollectorManifest(experimentName string, trialName string, jobKind string, namespace string, metricNames []string, mcs *experimentsv1alpha3.MetricsCollectorSpec) (string, error) {
	var mtp *template.Template = nil
	var err error
	tmpValues := map[string]string{
		"Experiment":     experimentName,
		"Trial":          trialName,
		"JobKind":        jobKind,
		"NameSpace":      namespace,
		"ManagerService": commonv1alpha3.GetManagerAddr(),
		"MetricNames":    strings.Join(metricNames, ";"),
	}
	if mcs != nil && mcs.GoTemplate.RawTemplate != "" {
		mtp, err = template.New("MetricsCollector").Parse(mcs.GoTemplate.RawTemplate)
	} else {
		mctp := defaultMetricsCollectorTemplateName
		if mcs != nil && mcs.GoTemplate.TemplateSpec != nil {
			mctp = mcs.GoTemplate.TemplateSpec.TemplatePath
		}
		mtl, err := g.client.GetMetricsCollectorTemplates()
		if err != nil {
			return "", err
		}
		if mt, ok := mtl[mctp]; !ok {
			return "", fmt.Errorf("No MetricsCollector template name %s", mctp)
		} else {
			mtp, err = template.New("MetricsCollector").Parse(mt)
		}
	}
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = mtp.Execute(&b, tmpValues)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// GetRunSpec get the specification for trial.
func (g *DefaultGenerator) GetRunSpec(e *experimentsv1alpha3.Experiment, experiment, trial, namespace string) (string, error) {
	params := trialTemplateParams{
		Experiment: experiment,
		Trial:      trial,
		NameSpace:  namespace,
	}
	return g.getRunSpec(e, params)
}

// GetRunSpecWithHyperParameters get the specification for trial with hyperparameters.
func (g *DefaultGenerator) GetRunSpecWithHyperParameters(e *experimentsv1alpha3.Experiment, experiment, trial, namespace string, hps []commonapiv1alpha3.ParameterAssignment) (string, error) {
	params := trialTemplateParams{
		Experiment:      experiment,
		Trial:           trial,
		NameSpace:       namespace,
		HyperParameters: hps,
	}
	return g.getRunSpec(e, params)
}

func (g *DefaultGenerator) getRunSpec(e *experimentsv1alpha3.Experiment, params trialTemplateParams) (string, error) {
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

func (g *DefaultGenerator) getTrialTemplate(instance *experimentsv1alpha3.Experiment) (*template.Template, error) {
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
	HyperParameters []commonapiv1alpha3.ParameterAssignment
}
