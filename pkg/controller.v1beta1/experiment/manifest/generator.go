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

package manifest

import (
	"fmt"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibclient"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibconfig"
)

// Generator is the type for manifests Generator.
type Generator interface {
	InjectClient(c client.Client)
	GetTrialTemplate(instance *experimentsv1beta1.Experiment) (string, error)
	GetRunSpecWithHyperParameters(experiment *experimentsv1beta1.Experiment, trialName, trialNamespace string, assignments []commonapiv1beta1.ParameterAssignment) (*unstructured.Unstructured, error)
	GetSuggestionConfigData(algorithmName string) (configv1beta1.SuggestionConfig, error)
	GetEarlyStoppingConfigData(algorithmName string) (configv1beta1.EarlyStoppingConfig, error)
	GetMetricsCollectorConfigData(cKind commonapiv1beta1.CollectorKind) (configv1beta1.MetricsCollectorConfig, error)
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

// GetMetricsCollectorConfigData returns metrics collector configuration for a given collector kind.
func (g *DefaultGenerator) GetMetricsCollectorConfigData(cKind commonapiv1beta1.CollectorKind) (configv1beta1.MetricsCollectorConfig, error) {
	return katibconfig.GetMetricsCollectorConfigData(cKind, g.client.GetClient())
}

// GetSuggestionConfigData returns suggestion configuration for a given algorithm name.
func (g *DefaultGenerator) GetSuggestionConfigData(algorithmName string) (configv1beta1.SuggestionConfig, error) {
	return katibconfig.GetSuggestionConfigData(algorithmName, g.client.GetClient())
}

// GetEarlyStoppingConfigData returns early stopping configuration for a given algorithm.
func (g *DefaultGenerator) GetEarlyStoppingConfigData(algorithmName string) (configv1beta1.EarlyStoppingConfig, error) {
	return katibconfig.GetEarlyStoppingConfigData(algorithmName, g.client.GetClient())
}

// GetRunSpecWithHyperParameters returns the specification for trial with hyperparameters.
func (g *DefaultGenerator) GetRunSpecWithHyperParameters(experiment *experimentsv1beta1.Experiment, trialName, trialNamespace string, assignments []commonapiv1beta1.ParameterAssignment) (*unstructured.Unstructured, error) {

	// Apply parameters to Trial Template from assignment
	replacedTemplate, err := g.applyParameters(experiment, trialName, trialNamespace, assignments)
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

func (g *DefaultGenerator) applyParameters(experiment *experimentsv1beta1.Experiment, trialName, trialNamespace string, assignments []commonapiv1beta1.ParameterAssignment) (string, error) {
	// Get string Trial template from Experiment spec
	trialTemplate, err := g.GetTrialTemplate(experiment)
	if err != nil {
		return "", err
	}

	trialSpec := experiment.Spec.TrialTemplate.TrialSpec
	// If trialSpec is not defined in TrialTemplate, deserialize templateString to fetch it
	if trialSpec == nil {
		trialSpec, err = util.ConvertStringToUnstructured(trialTemplate)
		if err != nil {
			return "", fmt.Errorf("ConvertStringToUnstructured failed: %v", err)
		}
	}

	// Convert parameter assignment to map key = parameter name, value = parameter value
	assignmentsMap := make(map[string]string)
	for _, assignment := range assignments {
		assignmentsMap[assignment.Name] = assignment.Value
	}

	placeHolderToValueMap := make(map[string]string)
	var metaRefKey, metaRefIndex string
	nonMetaParamCount := 0
	for _, param := range experiment.Spec.TrialTemplate.TrialParameters {
		metaMatchRegex := regexp.MustCompile(consts.TrialTemplateMetaReplaceFormatRegex)
		sub := metaMatchRegex.FindStringSubmatch(param.Reference)
		// handle trial parameters which consume trial assignments
		if len(sub) == 0 {
			if value, ok := assignmentsMap[param.Reference]; ok {
				placeHolderToValueMap[param.Name] = value
				nonMetaParamCount += 1
				continue
			} else {
				return "", fmt.Errorf("Unable to find parameter: %v in parameter assignment %v", param.Reference, assignmentsMap)
			}
		}
		metaRefKey = sub[1]

		// handle trial parameters which consume trial meta data
		// extract index (key) of Labels and Annotations if exists
		if sub := regexp.MustCompile(consts.TrialTemplateMetaParseFormatRegex).FindStringSubmatch(metaRefKey); len(sub) > 0 {
			if len(sub) != 3 {
				return "", fmt.Errorf("illegal reference of trial metadata: %v", param.Reference)
			}
			metaRefKey = sub[1]
			metaRefIndex = sub[2]
		}
		// fetch metadata value
		switch metaRefKey {
		case consts.TrialTemplateMetaKeyOfName:
			placeHolderToValueMap[param.Name] = trialName
		case consts.TrialTemplateMetaKeyOfNamespace:
			placeHolderToValueMap[param.Name] = trialNamespace
		case consts.TrialTemplateMetaKeyOfKind:
			placeHolderToValueMap[param.Name] = trialSpec.GetKind()
		case consts.TrialTemplateMetaKeyOfAPIVersion:
			placeHolderToValueMap[param.Name] = trialSpec.GetAPIVersion()
		case consts.TrialTemplateMetaKeyOfAnnotations:
			if value, ok := trialSpec.GetAnnotations()[metaRefIndex]; !ok {
				return "", fmt.Errorf("illegal reference of trial metadata: %v; failed to fetch Annotation: %v", param.Reference, metaRefIndex)
			} else {
				placeHolderToValueMap[param.Name] = value
			}
		case consts.TrialTemplateMetaKeyOfLabels:
			if value, ok := trialSpec.GetLabels()[metaRefIndex]; !ok {
				return "", fmt.Errorf("illegal reference of trial metadata: %v; failed to fetch Label: %v", param.Reference, metaRefIndex)
			} else {
				placeHolderToValueMap[param.Name] = value
			}
		default:
			return "", fmt.Errorf("illegal reference of trial metadata: %v", param.Reference)
		}
	}

	// Number of parameters must be equal
	if len(assignments) != nonMetaParamCount {
		return "", fmt.Errorf("Number of TrialAssignment: %v != number of nonMetaTrialParameters in TrialSpec: %v", len(assignments), nonMetaParamCount)
	}

	// Replacing placeholders with parameter values
	for placeHolder, paramValue := range placeHolderToValueMap {
		trialTemplate = strings.Replace(trialTemplate, fmt.Sprintf(consts.TrialTemplateParamReplaceFormat, placeHolder), paramValue, -1)
	}

	return trialTemplate, nil
}

// GetTrialTemplate returns string Trial template from experiment
func (g *DefaultGenerator) GetTrialTemplate(instance *experimentsv1beta1.Experiment) (string, error) {
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
			return "", fmt.Errorf("GetConfigMap failed: %v", err)
		}
		var ok bool
		trialTemplateString, ok = configMap[templatePath]
		if !ok {
			return "", fmt.Errorf("TemplatePath: %v not found in configMap: %v", templatePath, configMap)
		}
	}

	return trialTemplateString, nil
}
