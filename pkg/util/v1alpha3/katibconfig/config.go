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
	type suggestionConfigJSON struct {
		Image    string                      `json:"image"`
		Resource corev1.ResourceRequirements `json:"resources"`
	}
	if config, ok := configMap.Data[consts.LabelSuggestionTag]; ok {
		suggestionsConfig := map[string]suggestionConfigJSON{}
		if err := json.Unmarshal([]byte(config), &suggestionsConfig); err != nil {
			return map[string]string{}, err
		}
		if suggestionConfig, ok := suggestionsConfig[algorithmName]; ok {
			// Get image from config
			image := suggestionConfig.Image
			if strings.TrimSpace(image) != "" {
				suggestionConfigData[consts.LabelSuggestionImageTag] = image
			} else {
				return map[string]string{}, errors.New("Required value for " + consts.LabelSuggestionImageTag + " configuration of algorithm name " + algorithmName)
			}

			// Set default values for CPU and Memory
			suggestionConfigData[consts.LabelSuggestionCPURequestTag] = consts.DefaultCPURequest
			suggestionConfigData[consts.LabelSuggestionMemRequestTag] = consts.DefaultMemRequest
			suggestionConfigData[consts.LabelSuggestionDiskRequestTag] = consts.DefaultDiskRequest
			suggestionConfigData[consts.LabelSuggestionCPULimitTag] = consts.DefaultCPULimit
			suggestionConfigData[consts.LabelSuggestionMemLimitTag] = consts.DefaultMemLimit
			suggestionConfigData[consts.LabelSuggestionDiskLimitTag] = consts.DefaultDiskLimit

			// Get CPU and Memory Requests from config
			cpuRequest := suggestionConfig.Resource.Requests[corev1.ResourceCPU]
			memRequest := suggestionConfig.Resource.Requests[corev1.ResourceMemory]
			diskRequest := suggestionConfig.Resource.Requests[corev1.ResourceEphemeralStorage]
			if !cpuRequest.IsZero() {
				suggestionConfigData[consts.LabelSuggestionCPURequestTag] = cpuRequest.String()
			}
			if !memRequest.IsZero() {
				suggestionConfigData[consts.LabelSuggestionMemRequestTag] = memRequest.String()
			}
			if !diskRequest.IsZero() {
				suggestionConfigData[consts.LabelSuggestionDiskRequestTag] = diskRequest.String()
			}

			// Get CPU and Memory Limits from config
			cpuLimit := suggestionConfig.Resource.Limits[corev1.ResourceCPU]
			memLimit := suggestionConfig.Resource.Limits[corev1.ResourceMemory]
			diskLimit := suggestionConfig.Resource.Limits[corev1.ResourceEphemeralStorage]
			if !cpuLimit.IsZero() {
				suggestionConfigData[consts.LabelSuggestionCPULimitTag] = cpuLimit.String()
			}
			if !memLimit.IsZero() {
				suggestionConfigData[consts.LabelSuggestionMemLimitTag] = memLimit.String()
			}
			if !diskLimit.IsZero() {
				suggestionConfigData[consts.LabelSuggestionDiskRequestTag] = diskLimit.String()
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
