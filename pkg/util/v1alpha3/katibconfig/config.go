package katibconfig

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
)

func GetSuggestionConfigData(algorithmName string, client client.Client) (map[string]string, error) {
	configMap := &corev1.ConfigMap{}
	suggestionConfigData := map[string]string{}
	err := client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		return map[string]string{}, err
	}
	if config, ok := configMap.Data[consts.LabelSuggestionTag]; ok {
		suggestionsConfig := map[string]map[string]string{}
		if err := json.Unmarshal([]byte(config), &suggestionsConfig); err != nil {
			return map[string]string{}, err
		}
		if suggestionConfig, ok := suggestionsConfig[algorithmName]; ok {
			// Get image from config
			if image, yes := suggestionConfig[consts.LabelSuggestionImageTag]; yes {
				if strings.TrimSpace(image) != "" {
					suggestionConfigData[consts.LabelSuggestionImageTag] = image
				} else {
					return map[string]string{}, errors.New("Required value for " + consts.LabelSuggestionImageTag + " configuration of algorithm name " + algorithmName)
				}
			} else {
				return map[string]string{}, errors.New("Failed to find " + consts.LabelSuggestionImageTag + " configuration of algorithm name " + algorithmName)
			}
			// Get CPU Limit from config
			cpuLimit, yes := suggestionConfig[consts.LabelSuggestionCPULimitTag]
			if yes && strings.TrimSpace(cpuLimit) != "" {
				suggestionConfigData[consts.LabelSuggestionCPULimitTag] = cpuLimit
			} else {
				// Set default value
				suggestionConfigData[consts.LabelSuggestionCPULimitTag] = consts.DefaultCPULimit
			}
			// Get CPU Request from config
			cpuRequest, yes := suggestionConfig[consts.LabelSuggestionCPURequestTag]
			if yes && strings.TrimSpace(cpuRequest) != "" {
				suggestionConfigData[consts.LabelSuggestionCPURequestTag] = cpuRequest
			} else {
				// Set default value
				suggestionConfigData[consts.LabelSuggestionCPURequestTag] = consts.DefaultCPURequest
			}
			// Get Mem Limit from config
			memLimit, yes := suggestionConfig[consts.LabelSuggestionMemLimitTag]
			if yes && strings.TrimSpace(memLimit) != "" {
				suggestionConfigData[consts.LabelSuggestionMemLimitTag] = memLimit
			} else {
				// Set default value
				suggestionConfigData[consts.LabelSuggestionMemLimitTag] = consts.DefaultMemLimit
			}
			// Get Mem Request from config
			memRequest, yes := suggestionConfig[consts.LabelSuggestionMemRequestTag]
			if yes && strings.TrimSpace(memRequest) != "" {
				suggestionConfigData[consts.LabelSuggestionMemRequestTag] = memRequest
			} else {
				// Set default value
				suggestionConfigData[consts.LabelSuggestionMemRequestTag] = consts.DefaultMemRequest
			}
		} else {
			return map[string]string{}, errors.New("Failed to find algorithm " + algorithmName + " config in configmap " + consts.KatibConfigMapName)
		}
	} else {
		return map[string]string{}, errors.New("Failed to find suggestions config in configmap " + consts.KatibConfigMapName)
	}
	return suggestionConfigData, nil
}

func GetMetricsCollectorImage(cKind common.CollectorKind, client client.Client) (string, error) {
	configMap := &corev1.ConfigMap{}
	err := client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		return "", err
	}
	if mcs, ok := configMap.Data[consts.LabelMetricsCollectorSidecar]; ok {
		kind := string(cKind)
		mcsConfig := map[string]map[string]string{}
		if err := json.Unmarshal([]byte(mcs), &mcsConfig); err != nil {
			return "", err
		}
		if mc, ok := mcsConfig[kind]; ok {
			if image, yes := mc[consts.LabelMetricsCollectorSidecarImage]; yes {
				if strings.TrimSpace(image) != "" {
					return image, nil
				} else {
					return "", errors.New("Required value for " + consts.LabelMetricsCollectorSidecarImage + "configuration of metricsCollector kind " + kind)
				}
			} else {
				return "", errors.New("Failed to find " + consts.LabelMetricsCollectorSidecarImage + " configuration of metricsCollector kind " + kind)
			}
		} else {
			return "", errors.New("Cannot support metricsCollector injection for kind " + kind)
		}
	} else {
		return "", errors.New("Failed to find metrics collector configuration in configmap " + consts.KatibConfigMapName)
	}
}
