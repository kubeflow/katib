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

type suggestionConfigJSON struct {
	Image              string                      `json:"image"`
	ImagePullPolicy    corev1.PullPolicy           `json:"imagePullPolicy"`
	Resource           corev1.ResourceRequirements `json:"resources"`
	ServiceAccountName string                      `json:"serviceAccountName"`
}

type metricsCollectorConfigJSON struct {
	Image           string                      `json:"image"`
	ImagePullPolicy corev1.PullPolicy           `json:"imagePullPolicy"`
	Resource        corev1.ResourceRequirements `json:"resources"`
}

// GetSuggestionConfigData gets the config data for the given algorithm name.
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

			// Get Image Pull Policy
			imagePullPolicy := suggestionConfig.ImagePullPolicy
			if imagePullPolicy == corev1.PullAlways || imagePullPolicy == corev1.PullIfNotPresent || imagePullPolicy == corev1.PullNever {
				suggestionConfigData[consts.LabelSuggestionImagePullPolicy] = string(imagePullPolicy)
			} else {
				suggestionConfigData[consts.LabelSuggestionImagePullPolicy] = consts.DefaultImagePullPolicy
			}

			// Get Service Account Name
			serviceAccountName := suggestionConfig.ServiceAccountName
			if strings.TrimSpace(serviceAccountName) != "" {
				suggestionConfigData[consts.LabelSuggestionServiceAccountName] = serviceAccountName
			}

			// Set default values for CPU, Memory and Disk
			suggestionConfigData[consts.LabelSuggestionCPURequestTag] = consts.DefaultCPURequest
			suggestionConfigData[consts.LabelSuggestionMemRequestTag] = consts.DefaultMemRequest
			suggestionConfigData[consts.LabelSuggestionDiskRequestTag] = consts.DefaultDiskRequest
			suggestionConfigData[consts.LabelSuggestionCPULimitTag] = consts.DefaultCPULimit
			suggestionConfigData[consts.LabelSuggestionMemLimitTag] = consts.DefaultMemLimit
			suggestionConfigData[consts.LabelSuggestionDiskLimitTag] = consts.DefaultDiskLimit

			// Get CPU, Memory and Disk Requests from config
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

			// Get CPU, Memory and Disk Limits from config
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
				suggestionConfigData[consts.LabelSuggestionDiskLimitTag] = diskLimit.String()
			}

		} else {
			return map[string]string{}, errors.New("Failed to find algorithm " + algorithmName + " config in configmap " + consts.KatibConfigMapName)
		}
	} else {
		return map[string]string{}, errors.New("Failed to find suggestions config in configmap " + consts.KatibConfigMapName)
	}
	return suggestionConfigData, nil
}

// GetMetricsCollectorConfigData gets the config data for the given kind.
func GetMetricsCollectorConfigData(cKind common.CollectorKind, client client.Client) (map[string]string, error) {
	configMap := &corev1.ConfigMap{}
	metricsCollectorConfigData := map[string]string{}
	err := client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		return metricsCollectorConfigData, err
	}
	// Get the config with name metrics-collector-sidecar.
	if config, ok := configMap.Data[consts.LabelMetricsCollectorSidecar]; ok {
		kind := string(cKind)
		mcsConfig := map[string]metricsCollectorConfigJSON{}
		if err := json.Unmarshal([]byte(config), &mcsConfig); err != nil {
			return metricsCollectorConfigData, err
		}
		// Get the config for the given cKind.
		if metricsCollectorConfig, ok := mcsConfig[kind]; ok {
			image := metricsCollectorConfig.Image
			// If the image is not empty, we set it into result.
			if strings.TrimSpace(image) != "" {
				metricsCollectorConfigData[consts.LabelMetricsCollectorSidecarImage] = image
			} else {
				return metricsCollectorConfigData, errors.New("Required value for " + consts.LabelMetricsCollectorSidecarImage + "configuration of metricsCollector kind " + kind)
			}

			// Get Image Pull Policy
			imagePullPolicy := metricsCollectorConfig.ImagePullPolicy
			if imagePullPolicy == corev1.PullAlways || imagePullPolicy == corev1.PullIfNotPresent || imagePullPolicy == corev1.PullNever {
				metricsCollectorConfigData[consts.LabelMetricsCollectorImagePullPolicy] = string(imagePullPolicy)
			} else {
				metricsCollectorConfigData[consts.LabelMetricsCollectorImagePullPolicy] = consts.DefaultImagePullPolicy
			}

			// Set default values for CPU, Memory and Disk
			metricsCollectorConfigData[consts.LabelMetricsCollectorCPURequestTag] = consts.DefaultCPURequest
			metricsCollectorConfigData[consts.LabelMetricsCollectorMemRequestTag] = consts.DefaultMemRequest
			metricsCollectorConfigData[consts.LabelMetricsCollectorDiskRequestTag] = consts.DefaultDiskRequest
			metricsCollectorConfigData[consts.LabelMetricsCollectorCPULimitTag] = consts.DefaultCPULimit
			metricsCollectorConfigData[consts.LabelMetricsCollectorMemLimitTag] = consts.DefaultMemLimit
			metricsCollectorConfigData[consts.LabelMetricsCollectorDiskLimitTag] = consts.DefaultDiskLimit

			// Get CPU, Memory and Disk Requests from config
			cpuRequest := metricsCollectorConfig.Resource.Requests[corev1.ResourceCPU]
			memRequest := metricsCollectorConfig.Resource.Requests[corev1.ResourceMemory]
			diskRequest := metricsCollectorConfig.Resource.Requests[corev1.ResourceEphemeralStorage]
			if !cpuRequest.IsZero() {
				metricsCollectorConfigData[consts.LabelSuggestionCPURequestTag] = cpuRequest.String()
			}
			if !memRequest.IsZero() {
				metricsCollectorConfigData[consts.LabelSuggestionMemRequestTag] = memRequest.String()
			}
			if !diskRequest.IsZero() {
				metricsCollectorConfigData[consts.LabelSuggestionDiskRequestTag] = diskRequest.String()
			}

			// Get CPU, Memory and Disk Limits from config
			cpuLimit := metricsCollectorConfig.Resource.Limits[corev1.ResourceCPU]
			memLimit := metricsCollectorConfig.Resource.Limits[corev1.ResourceMemory]
			diskLimit := metricsCollectorConfig.Resource.Limits[corev1.ResourceEphemeralStorage]
			if !cpuLimit.IsZero() {
				metricsCollectorConfigData[consts.LabelSuggestionCPULimitTag] = cpuLimit.String()
			}
			if !memLimit.IsZero() {
				metricsCollectorConfigData[consts.LabelSuggestionMemLimitTag] = memLimit.String()
			}
			if !diskLimit.IsZero() {
				metricsCollectorConfigData[consts.LabelSuggestionDiskLimitTag] = diskLimit.String()
			}

		} else {
			return metricsCollectorConfigData, errors.New("Cannot support metricsCollector injection for kind " + kind)
		}
	} else {
		return metricsCollectorConfigData, errors.New("Failed to find metrics collector configuration in configmap " + consts.KatibConfigMapName)
	}
	return metricsCollectorConfigData, nil
}
