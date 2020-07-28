package katibconfig

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

// SuggestionConfig is the JSON suggestion structure in Katib config
type SuggestionConfig struct {
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
func GetSuggestionConfigData(algorithmName string, client client.Client) (SuggestionConfig, error) {
	configMap := &corev1.ConfigMap{}
	suggestionConfigData := SuggestionConfig{}
	err := client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		return SuggestionConfig{}, err
	}

	// Try to find suggestion data in config map
	config, ok := configMap.Data[consts.LabelSuggestionTag]
	if !ok {
		return SuggestionConfig{}, errors.New("Failed to find suggestions config in configmap " + consts.KatibConfigMapName)
	}

	// Parse suggestion data to Suggestion Config
	suggestionsConfig := map[string]SuggestionConfig{}
	if err := json.Unmarshal([]byte(config), &suggestionsConfig); err != nil {
		return SuggestionConfig{}, err
	}

	// Try to find suggestion algorithm in suggestion data
	suggestionConfigData, ok = suggestionsConfig[algorithmName]
	if !ok {
		return SuggestionConfig{}, errors.New("Failed to find algorithm " + algorithmName + " config in configmap " + consts.KatibConfigMapName)
	}

	// Get image from config
	image := suggestionConfigData.Image
	if strings.TrimSpace(image) == "" {
		return SuggestionConfig{}, errors.New("Required value for image configuration of algorithm name " + algorithmName)
	}

	// Get Image Pull Policy
	imagePullPolicy := suggestionConfigData.ImagePullPolicy
	if imagePullPolicy != corev1.PullAlways && imagePullPolicy != corev1.PullIfNotPresent && imagePullPolicy != corev1.PullNever {
		// TODO (andreyvelich): Change it to consts once metrics collector config is refactored
		suggestionConfigData.ImagePullPolicy = corev1.PullIfNotPresent
	}

	// Set default values for CPU, Memory and Disk
	defaultCPURequest, _ := resource.ParseQuantity(consts.DefaultCPURequest)
	defaultMemRequest, _ := resource.ParseQuantity(consts.DefaultMemRequest)
	defaultDiskRequest, _ := resource.ParseQuantity(consts.DefaultDiskRequest)

	defaultCPULimit, _ := resource.ParseQuantity(consts.DefaultCPULimit)
	defaultMemLimit, _ := resource.ParseQuantity(consts.DefaultMemLimit)
	defaultDiskLimit, _ := resource.ParseQuantity(consts.DefaultDiskLimit)

	// Get CPU, Memory and Disk Requests from config
	cpuRequest := suggestionConfigData.Resource.Requests[corev1.ResourceCPU]
	memRequest := suggestionConfigData.Resource.Requests[corev1.ResourceMemory]
	diskRequest := suggestionConfigData.Resource.Requests[corev1.ResourceEphemeralStorage]

	// If resource is empty set default value
	if cpuRequest.IsZero() {
		suggestionConfigData.Resource.Requests[corev1.ResourceCPU] = defaultCPURequest
	}
	if memRequest.IsZero() {
		suggestionConfigData.Resource.Requests[corev1.ResourceMemory] = defaultMemRequest

	}
	if diskRequest.IsZero() {
		suggestionConfigData.Resource.Requests[corev1.ResourceEphemeralStorage] = defaultDiskRequest
	}

	// Get CPU, Memory and Disk Limits from config
	cpuLimit := suggestionConfigData.Resource.Limits[corev1.ResourceCPU]
	memLimit := suggestionConfigData.Resource.Limits[corev1.ResourceMemory]
	diskLimit := suggestionConfigData.Resource.Limits[corev1.ResourceEphemeralStorage]
	if cpuLimit.IsZero() {
		suggestionConfigData.Resource.Limits[corev1.ResourceCPU] = defaultCPULimit
	}
	if memLimit.IsZero() {
		suggestionConfigData.Resource.Limits[corev1.ResourceMemory] = defaultMemLimit
	}
	if diskLimit.IsZero() {
		suggestionConfigData.Resource.Limits[corev1.ResourceEphemeralStorage] = defaultDiskLimit
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
				metricsCollectorConfigData[consts.LabelMetricsCollectorCPURequestTag] = cpuRequest.String()
			}
			if !memRequest.IsZero() {
				metricsCollectorConfigData[consts.LabelMetricsCollectorMemRequestTag] = memRequest.String()
			}
			if !diskRequest.IsZero() {
				metricsCollectorConfigData[consts.LabelMetricsCollectorDiskRequestTag] = diskRequest.String()
			}

			// Get CPU, Memory and Disk Limits from config
			cpuLimit := metricsCollectorConfig.Resource.Limits[corev1.ResourceCPU]
			memLimit := metricsCollectorConfig.Resource.Limits[corev1.ResourceMemory]
			diskLimit := metricsCollectorConfig.Resource.Limits[corev1.ResourceEphemeralStorage]
			if !cpuLimit.IsZero() {
				metricsCollectorConfigData[consts.LabelMetricsCollectorCPULimitTag] = cpuLimit.String()
			}
			if !memLimit.IsZero() {
				metricsCollectorConfigData[consts.LabelMetricsCollectorMemLimitTag] = memLimit.String()
			}
			if !diskLimit.IsZero() {
				metricsCollectorConfigData[consts.LabelMetricsCollectorDiskLimitTag] = diskLimit.String()
			}

		} else {
			return metricsCollectorConfigData, errors.New("Cannot support metricsCollector injection for kind " + kind)
		}
	} else {
		return metricsCollectorConfigData, errors.New("Failed to find metrics collector configuration in configmap " + consts.KatibConfigMapName)
	}
	return metricsCollectorConfigData, nil
}
