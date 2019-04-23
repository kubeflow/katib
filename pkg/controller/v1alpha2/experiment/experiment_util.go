package experiment

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"text/template"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	apiv1alpha2 "github.com/kubeflow/katib/pkg/api/v1alpha2"
)

type TrialTemplateParams struct {
	Experiment      string
	Trial           string
	NameSpace       string
	HyperParameters []*apiv1alpha2.ParameterAssignment
}

func (r *ReconcileExperiment) createTrialInstance(expInstance *experimentsv1alpha2.Experiment, trialInstance *apiv1alpha2.Trial, trialTemplate *template.Template) error {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: expInstance.GetName(), Namespace: expInstance.GetNamespace()})

	trial := &trialsv1alpha2.Trial{}
	trial.Name = fmt.Sprintf("%s-%s", expInstance.GetName(), utilrand.String(8))
	trial.Namespace = expInstance.GetNamespace()
	trial.Labels = map[string]string{"experiment": expInstance.GetName()}

	if err := controllerutil.SetControllerReference(expInstance, trial, r.scheme); err != nil {
		logger.Error(err, "Set controller reference error")
		return err
	}

	trialParams := TrialTemplateParams{
		Experiment: expInstance.GetName(),
		Trial:      trial.Name,
		NameSpace:  trial.Namespace,
	}

	var buf bytes.Buffer
	if trialInstance.Spec != nil && trialInstance.Spec.ParameterAssignments != nil {
		for _, p := range trialInstance.Spec.ParameterAssignments.Assignments {
			trialParams.HyperParameters = append(trialParams.HyperParameters, p)
		}
	}
	if err := trialTemplate.Execute(&buf, trialParams); err != nil {
		logger.Error(err, "Template execute error")
		return err
	}

	trial.Spec.RunSpec = buf.String()

	if err := r.Create(context.TODO(), trial); err != nil {
		logger.Error(err, "Trial create error", "Trial name", trial.Name)
		return err
	}
	return nil

}

func (r *ReconcileExperiment) getTrialTemplate(instance *experimentsv1alpha2.Experiment) (*template.Template, error) {

	var err error
	var tpl *template.Template = nil
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	trialTemplate := instance.Spec.TrialTemplate
	if trialTemplate != nil && trialTemplate.GoTemplate.RawTemplate != "" {
		tpl, err = template.New("Trial").Parse(trialTemplate.GoTemplate.RawTemplate)
	} else {
		//default values if user hasn't set
		configMapNS := os.Getenv(experimentsv1alpha2.DefaultKatibNamespaceEnvName)
		configMapName := experimentsv1alpha2.DefaultTrialConfigMapName
		templatePath := experimentsv1alpha2.DefaultTrialTemplatePath

		if trialTemplate != nil && trialTemplate.GoTemplate.TemplateSpec != nil {
			templateSpec := trialTemplate.GoTemplate.TemplateSpec
			if templateSpec.ConfigMapName != "" {
				configMapName = templateSpec.ConfigMapName
			}
			if templateSpec.ConfigMapNamespace != "" {
				configMapNS = templateSpec.ConfigMapNamespace
			}
			if templateSpec.TemplatePath != "" {
				templatePath = templateSpec.TemplatePath
			}
		}
		configMap, err := r.getConfigMap(configMapName, configMapNS)
		if err != nil {
			logger.Error(err, "Get config map error", "configMapName", configMapName, "configMapNS", configMapNS)
			return nil, err
		}
		if configMapTemplate, ok := configMap[templatePath]; !ok {
			err = errors.New(string(metav1.StatusReasonNotFound))
			logger.Error(err, "Config map template not found", "templatePath", templatePath)
			return nil, err
		} else {
			tpl, err = template.New("Trial").Parse(configMapTemplate)
		}
	}
	if err != nil {
		logger.Error(err, "Template parse error")
		return nil, err
	}

	return tpl, nil
}

func (r *ReconcileExperiment) getConfigMap(name, namespace string) (map[string]string, error) {

	configMap := &apiv1.ConfigMap{}
	if err := r.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, configMap); err != nil {
		return map[string]string{}, err
	}
	return configMap.Data, nil
}
