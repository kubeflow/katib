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
	Image                     string                           `json:"image"`
	ImagePullPolicy           corev1.PullPolicy                `json:"imagePullPolicy"`
	Resource                  corev1.ResourceRequirements      `json:"resources"`
	ServiceAccountName        string                           `json:"serviceAccountName"`
	VolumeMountPath           string                           `json:"volumeMountPath"`
	PersistentVolumeClaimSpec corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimSpec"`
	PersistentVolumeSpec      corev1.PersistentVolumeSpec      `json:"persistentVolumeSpec"`
}

// MetricsCollectorConfig is the JSON metrics collector structure in Katib config
type MetricsCollectorConfig struct {
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
		return SuggestionConfig{}, errors.New("Failed to find suggestions config in ConfigMap: " + consts.KatibConfigMapName)
	}

	// Parse suggestion data to map where key = algorithm name, value = SuggestionConfig
	suggestionsConfig := map[string]SuggestionConfig{}
	if err := json.Unmarshal([]byte(config), &suggestionsConfig); err != nil {
		return SuggestionConfig{}, err
	}

	// Try to find SuggestionConfig for the algorithm
	suggestionConfigData, ok = suggestionsConfig[algorithmName]
	if !ok {
		return SuggestionConfig{}, errors.New("Failed to find suggestion config for algorithm: " + algorithmName + " in ConfigMap: " + consts.KatibConfigMapName)
	}

	// Get image from config
	image := suggestionConfigData.Image
	if strings.TrimSpace(image) == "" {
		return SuggestionConfig{}, errors.New("Required value for image configuration of algorithm name: " + algorithmName)
	}

	// Get Image Pull Policy
	imagePullPolicy := suggestionConfigData.ImagePullPolicy
	if imagePullPolicy != corev1.PullAlways && imagePullPolicy != corev1.PullIfNotPresent && imagePullPolicy != corev1.PullNever {
		suggestionConfigData.ImagePullPolicy = consts.DefaultImagePullPolicy
	}

	// Set resource requirements for suggestion
	suggestionConfigData.Resource = setResourceRequirements(suggestionConfigData.Resource)

	// Set default suggestion container volume mount path
	if suggestionConfigData.VolumeMountPath == "" {
		suggestionConfigData.VolumeMountPath = consts.DefaultContainerSuggestionVolumeMountPath
	}

	// Get persistent volume claim spec from config
	pvcSpec := suggestionConfigData.PersistentVolumeClaimSpec

	// Set default storage class
	defaultStorageClassName := consts.DefaultSuggestionStorageClassName
	if pvcSpec.StorageClassName == nil {
		pvcSpec.StorageClassName = &defaultStorageClassName
	}

	// Set default access modes
	if len(pvcSpec.AccessModes) == 0 {
		pvcSpec.AccessModes = []corev1.PersistentVolumeAccessMode{
			consts.DefaultSuggestionVolumeAccessMode,
		}
	}

	// Set default resources
	defaultVolumeStorage, _ := resource.ParseQuantity(consts.DefaultSuggestionVolumeStorage)
	if len(pvcSpec.Resources.Requests) == 0 {

		pvcSpec.Resources.Requests = make(map[corev1.ResourceName]resource.Quantity)
		pvcSpec.Resources.Requests[corev1.ResourceStorage] = defaultVolumeStorage
	}

	// Set pvc back for suggestion config
	suggestionConfigData.PersistentVolumeClaimSpec = pvcSpec

	// Get pv from config only if pvc storage class name = DefaultSuggestionStorageClassName
	if *pvcSpec.StorageClassName == consts.DefaultSuggestionStorageClassName {
		pvSpec := suggestionConfigData.PersistentVolumeSpec

		// Set default storage class
		pvSpec.StorageClassName = defaultStorageClassName

		// Set default access modes
		if len(pvSpec.AccessModes) == 0 {
			pvSpec.AccessModes = []corev1.PersistentVolumeAccessMode{
				consts.DefaultSuggestionVolumeAccessMode,
			}
		}

		// Set default pv source.
		// In composer we add name, algorithm and namespace to host path.
		if pvSpec.PersistentVolumeSource == (corev1.PersistentVolumeSource{}) {
			pvSpec.PersistentVolumeSource = corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: consts.DefaultSuggestionVolumeLocalPathPrefix,
				},
			}
		}

		// Set default local path if it is empty
		if pvSpec.PersistentVolumeSource.HostPath != nil && pvSpec.PersistentVolumeSource.HostPath.Path == "" {
			pvSpec.PersistentVolumeSource.HostPath.Path = consts.DefaultSuggestionVolumeLocalPathPrefix
		}

		// Set default capacity
		if len(pvSpec.Capacity) == 0 {
			pvSpec.Capacity = make(map[corev1.ResourceName]resource.Quantity)
			pvSpec.Capacity[corev1.ResourceStorage] = defaultVolumeStorage
		}

		// Set pv back for suggestion config
		suggestionConfigData.PersistentVolumeSpec = pvSpec

	}

	return suggestionConfigData, nil
}

// GetMetricsCollectorConfigData gets the config data for the given collector kind.
func GetMetricsCollectorConfigData(cKind common.CollectorKind, client client.Client) (MetricsCollectorConfig, error) {
	configMap := &corev1.ConfigMap{}
	metricsCollectorConfigData := MetricsCollectorConfig{}
	err := client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		return MetricsCollectorConfig{}, err
	}

	// Try to find metrics collector data in config map
	config, ok := configMap.Data[consts.LabelMetricsCollectorSidecar]
	if !ok {
		return MetricsCollectorConfig{}, errors.New("Failed to find metrics collector config in ConfigMap: " + consts.KatibConfigMapName)
	}
	// Parse metrics collector data to map where key = collector kind, value = MetricsCollectorConfig
	kind := string(cKind)
	mcsConfig := map[string]MetricsCollectorConfig{}
	if err := json.Unmarshal([]byte(config), &mcsConfig); err != nil {
		return MetricsCollectorConfig{}, err
	}

	// Try to find MetricsCollectorConfig for the collector kind
	metricsCollectorConfigData, ok = mcsConfig[kind]
	if !ok {
		return MetricsCollectorConfig{}, errors.New("Failed to find metrics collector config for kind: " + kind + " in ConfigMap: " + consts.KatibConfigMapName)
	}

	// Get image from config
	image := metricsCollectorConfigData.Image
	if strings.TrimSpace(image) == "" {
		return MetricsCollectorConfig{}, errors.New("Required value for image configuration of metrics collector kind: " + kind)
	}

	// Get Image Pull Policy
	imagePullPolicy := metricsCollectorConfigData.ImagePullPolicy
	if imagePullPolicy != corev1.PullAlways && imagePullPolicy != corev1.PullIfNotPresent && imagePullPolicy != corev1.PullNever {
		metricsCollectorConfigData.ImagePullPolicy = consts.DefaultImagePullPolicy
	}

	// Set resource requirements for metrics collector
	metricsCollectorConfigData.Resource = setResourceRequirements(metricsCollectorConfigData.Resource)

	return metricsCollectorConfigData, nil
}

func setResourceRequirements(configResource corev1.ResourceRequirements) corev1.ResourceRequirements {

	// If requests are empty create new map
	if len(configResource.Requests) == 0 {
		configResource.Requests = make(map[corev1.ResourceName]resource.Quantity)
	}

	// Get CPU, Memory and Disk Requests from config
	cpuRequest := configResource.Requests[corev1.ResourceCPU]
	memRequest := configResource.Requests[corev1.ResourceMemory]
	diskRequest := configResource.Requests[corev1.ResourceEphemeralStorage]

	// If resource is empty set default value for CPU, Memory, Disk
	if cpuRequest.IsZero() {
		defaultCPURequest, _ := resource.ParseQuantity(consts.DefaultCPURequest)
		configResource.Requests[corev1.ResourceCPU] = defaultCPURequest
	}
	if memRequest.IsZero() {
		defaultMemRequest, _ := resource.ParseQuantity(consts.DefaultMemRequest)
		configResource.Requests[corev1.ResourceMemory] = defaultMemRequest
	}
	if diskRequest.IsZero() {
		defaultDiskRequest, _ := resource.ParseQuantity(consts.DefaultDiskRequest)
		configResource.Requests[corev1.ResourceEphemeralStorage] = defaultDiskRequest
	}

	// If limits are empty create new map
	if len(configResource.Limits) == 0 {
		configResource.Limits = make(map[corev1.ResourceName]resource.Quantity)
	}

	// Get CPU, Memory and Disk Limits from config
	cpuLimit := configResource.Limits[corev1.ResourceCPU]
	memLimit := configResource.Limits[corev1.ResourceMemory]
	diskLimit := configResource.Limits[corev1.ResourceEphemeralStorage]

	// If limit is empty set default value for CPU, Memory, Disk
	if cpuLimit.IsZero() {
		defaultCPULimit, _ := resource.ParseQuantity(consts.DefaultCPULimit)
		configResource.Limits[corev1.ResourceCPU] = defaultCPULimit
	}
	if memLimit.IsZero() {
		defaultMemLimit, _ := resource.ParseQuantity(consts.DefaultMemLimit)
		configResource.Limits[corev1.ResourceMemory] = defaultMemLimit
	}
	if diskLimit.IsZero() {
		defaultDiskLimit, _ := resource.ParseQuantity(consts.DefaultDiskLimit)
		configResource.Limits[corev1.ResourceEphemeralStorage] = defaultDiskLimit
	}

	// If user explicitly sets ephemeral-storage value to something negative, nuke it.
	// This enables compability with the GKE nodepool autoscalers, which cannot scale
	// pods which define ephemeral-storage resource constraints.
	if diskLimit.Sign() == -1 && diskRequest.Sign() == -1 {
		delete(configResource.Limits, corev1.ResourceEphemeralStorage)
		delete(configResource.Requests, corev1.ResourceEphemeralStorage)
	}
	return configResource
}
