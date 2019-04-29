/*

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

package util

import (
	"bytes"
	"context"
	"errors"
	"text/template"

	apiv1alpha2 "github.com/kubeflow/katib/pkg/api/v1alpha2"
	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TrialTemplateParams struct {
	Experiment      string
	Trial           string
	NameSpace       string
	HyperParameters []*apiv1alpha2.ParameterAssignment
}

func GetRunSpec(e *experimentsv1alpha2.Experiment, params TrialTemplateParams) (string, error) {
	var buf bytes.Buffer
	trialTemplate, err := getTrialTemplate(e)
	if err != nil {
		return "", err
	}
	if err := trialTemplate.Execute(&buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getTrialTemplate(instance *experimentsv1alpha2.Experiment) (*template.Template, error) {
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

		configMap, err := getConfigMap(configMapName, configMapNS)
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

func getConfigMap(name, namespace string) (map[string]string, error) {
	client, err := NewClient(client.Options{})
	if err != nil {
		return map[string]string{}, err
	}
	configMap := &apiv1.ConfigMap{}
	if err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, configMap); err != nil {
		return map[string]string{}, err
	}
	return configMap.Data, nil
}
