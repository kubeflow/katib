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

package katibconfig

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	configapi "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrKatibConfigNil         = fmt.Errorf("failed to parse katib-config.yaml in ConfigMap: %s", consts.KatibConfigMapName)
	ErrInvalidGVKFormat       = errors.New("invalid GroupVersionKinds")
	ErrTrialResourcesAreEmpty = errors.New("trialResources are empty")
)

func TrialResourcesToGVKs(trialResources []string) ([]schema.GroupVersionKind, error) {
	if len(trialResources) == 0 {
		return nil, ErrTrialResourcesAreEmpty
	}
	gvks := make([]schema.GroupVersionKind, 0, len(trialResources))
	for i := range trialResources {
		gvk, _ := schema.ParseKindArg(trialResources[i])
		if gvk == nil {
			return nil, ErrInvalidGVKFormat
		}
		gvks = append(gvks, *gvk)
	}
	return gvks, nil
}

// GetSuggestionConfigData gets the config data for the given suggestion algorithm name.
func GetSuggestionConfigData(algorithmName string, client client.Client) (configapi.SuggestionConfig, error) {
	katibCfg := &configapi.KatibConfig{}
	if err := fromConfigMap(katibCfg, client); err != nil {
		return configapi.SuggestionConfig{}, err
	}

	// Try to find SuggestionConfig for the algorithm
	var suggestionConfigData *configapi.SuggestionConfig
	for i := range katibCfg.RuntimeConfig.SuggestionConfigs {
		if katibCfg.RuntimeConfig.SuggestionConfigs[i].AlgorithmName == algorithmName {
			suggestionConfigData = &katibCfg.RuntimeConfig.SuggestionConfigs[i]
		}
	}
	if suggestionConfigData == nil {
		return configapi.SuggestionConfig{}, fmt.Errorf("failed to find suggestion config for algorithm: %s in ConfigMap: %s", algorithmName, consts.KatibConfigMapName)
	}

	// Get image from config
	image := suggestionConfigData.Image
	if strings.TrimSpace(image) == "" {
		return configapi.SuggestionConfig{}, fmt.Errorf("required value for image configuration of algorithm name: %s", algorithmName)
	}
	return *suggestionConfigData, nil
}

// GetEarlyStoppingConfigData gets the config data for the given early stopping algorithm name.
func GetEarlyStoppingConfigData(algorithmName string, client client.Client) (configapi.EarlyStoppingConfig, error) {
	katibCfg := &configapi.KatibConfig{}
	if err := fromConfigMap(katibCfg, client); err != nil {
		return configapi.EarlyStoppingConfig{}, err
	}

	// Try to find EarlyStoppingConfig for the algorithm
	var earlyStoppingConfigData *configapi.EarlyStoppingConfig
	for i := range katibCfg.RuntimeConfig.EarlyStoppingConfigs {
		if katibCfg.RuntimeConfig.EarlyStoppingConfigs[i].AlgorithmName == algorithmName {
			earlyStoppingConfigData = &katibCfg.RuntimeConfig.EarlyStoppingConfigs[i]
		}
	}
	if earlyStoppingConfigData == nil {
		return configapi.EarlyStoppingConfig{}, fmt.Errorf("failed to find early stopping config for algorithm: %s in ConfigMap: %s", algorithmName, consts.KatibConfigMapName)
	}

	// Get image from config.
	image := earlyStoppingConfigData.Image
	if strings.TrimSpace(image) == "" {
		return configapi.EarlyStoppingConfig{}, fmt.Errorf("required value for image configuration of algorithm name: %s", algorithmName)
	}

	return *earlyStoppingConfigData, nil
}

// GetMetricsCollectorConfigData gets the config data for the given collector kind.
func GetMetricsCollectorConfigData(cKind common.CollectorKind, client client.Client) (configapi.MetricsCollectorConfig, error) {
	katibCfg := &configapi.KatibConfig{}
	if err := fromConfigMap(katibCfg, client); err != nil {
		return configapi.MetricsCollectorConfig{}, err
	}

	// Try to find MetricsCollectorConfig for the collector kind
	var metricsCollectorConfigData *configapi.MetricsCollectorConfig
	kind := string(cKind)
	for i := range katibCfg.RuntimeConfig.MetricsCollectorConfigs {
		if katibCfg.RuntimeConfig.MetricsCollectorConfigs[i].CollectorKind == kind {
			metricsCollectorConfigData = &katibCfg.RuntimeConfig.MetricsCollectorConfigs[i]
		}
	}
	if metricsCollectorConfigData == nil {
		return configapi.MetricsCollectorConfig{}, fmt.Errorf("failed to find metrics collector config for kind: %s in ConfigMap: %s", kind, consts.KatibConfigMapName)
	}

	// Get image from config
	image := metricsCollectorConfigData.Image
	if strings.TrimSpace(image) == "" {
		return configapi.MetricsCollectorConfig{}, fmt.Errorf("required value for image configuration of metrics collector kind: %s", kind)
	}

	return *metricsCollectorConfigData, nil
}

// GetInitConfigData gets the init config data.
func GetInitConfigData(scheme *runtime.Scheme, katibCfgPath string) (configapi.InitConfig, error) {
	var katibCfg configapi.KatibConfig
	if err := fromFile(scheme, &katibCfg, katibCfgPath); err != nil {
		return configapi.InitConfig{}, fmt.Errorf("%w: %s", ErrKatibConfigNil, err.Error())
	}
	return katibCfg.InitConfig, nil
}

func fromFile(scheme *runtime.Scheme, katibConfig *configapi.KatibConfig, katibConfigPath string) error {
	if len(katibConfigPath) == 0 {
		scheme.Default(katibConfig)
		return nil
	}
	config, err := os.ReadFile(katibConfigPath)
	if err != nil {
		return err
	}
	codecs := serializer.NewCodecFactory(scheme)
	return runtime.DecodeInto(codecs.UniversalDecoder(), config, katibConfig)
}

func fromConfigMap(katibConfig *configapi.KatibConfig, client client.Client) error {
	configMap := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace}, configMap)
	if err != nil {
		return err
	}
	// Try to find katib-config.yaml data in configMap.
	config, ok := configMap.Data[consts.LabelKatibConfigTag]
	if !ok {
		return fmt.Errorf("failed to find katib-config.yaml in ConfigMap: %s", consts.KatibConfigMapName)
	}
	codecs := serializer.NewCodecFactory(client.Scheme())
	return runtime.DecodeInto(codecs.UniversalDecoder(), []byte(config), katibConfig)
}
