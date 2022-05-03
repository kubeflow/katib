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
	"encoding/json"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

// SuggestionConfig is the JSON suggestion structure in Katib config.
type SuggestionConfig struct {
	Image                     string                           `json:"image"`
	ImagePullPolicy           corev1.PullPolicy                `json:"imagePullPolicy,omitempty"`
	Resource                  corev1.ResourceRequirements      `json:"resources,omitempty"`
	ServiceAccountName        string                           `json:"serviceAccountName,omitempty"`
	VolumeMountPath           string                           `json:"volumeMountPath,omitempty"`
	PersistentVolumeClaimSpec corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimSpec,omitempty"`
	PersistentVolumeSpec      corev1.PersistentVolumeSpec      `json:"persistentVolumeSpec,omitempty"`
	PersistentVolumeLabels    map[string]string                `json:"persistentVolumeLabels,omitempty"`
}

// EarlyStoppingConfig is the JSON early stopping structure in Katib config.
type EarlyStoppingConfig struct {
	Image           string            `json:"image"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

// MetricsCollectorConfig is the JSON metrics collector structure in Katib config.
type MetricsCollectorConfig struct {
	Image            string                      `json:"image"`
	ImagePullPolicy  corev1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	Resource         corev1.ResourceRequirements `json:"resources,omitempty"`
	WaitAllProcesses *bool                       `json:"waitAllProcesses,omitempty"`
}

// GetSuggestionConfigData gets the config data for the given suggestion algorithm name.
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
		return SuggestionConfig{}, fmt.Errorf("failed to find suggestions config in ConfigMap: %s", consts.KatibConfigMapName)
	}

	// Parse suggestion data to map where key = algorithm name, value = SuggestionConfig
	suggestionsConfig := map[string]SuggestionConfig{}
	if err := json.Unmarshal([]byte(config), &suggestionsConfig); err != nil {
		return SuggestionConfig{}, err
	}

	// Try to find SuggestionConfig for the algorithm
	suggestionConfigData, ok = suggestionsConfig[algorithmName]
	if !ok {
		return SuggestionConfig{}, fmt.Errorf("failed to find suggestion config for algorithm: %s in ConfigMap: %s", algorithmName, consts.KatibConfigMapName)
	}

	// Get image from config
	image := suggestionConfigData.Image
	if strings.TrimSpace(image) == "" {
		return SuggestionConfig{}, fmt.Errorf("required value for image configuration of algorithm name: %s", algorithmName)
	}

	// Set Image Pull Policy
	suggestionConfigData.ImagePullPolicy = setImagePullPolicy(suggestionConfigData.ImagePullPolicy)

	// Set resource requirements for suggestion
	suggestionConfigData.Resource = setResourceRequirements(suggestionConfigData.Resource)

	// Set default suggestion container volume mount path
	if suggestionConfigData.VolumeMountPath == "" {
		suggestionConfigData.VolumeMountPath = consts.DefaultContainerSuggestionVolumeMountPath
	}

	// Get persistent volume claim spec from config
	pvcSpec := suggestionConfigData.PersistentVolumeClaimSpec

	// Set default access modes
	if len(pvcSpec.AccessModes) == 0 {
		pvcSpec.AccessModes = []corev1.PersistentVolumeAccessMode{
			consts.DefaultSuggestionVolumeAccessMode,
		}
	}

	// Set default resources
	if len(pvcSpec.Resources.Requests) == 0 {
		defaultVolumeStorage, _ := resource.ParseQuantity(consts.DefaultSuggestionVolumeStorage)
		pvcSpec.Resources.Requests = make(map[corev1.ResourceName]resource.Quantity)
		pvcSpec.Resources.Requests[corev1.ResourceStorage] = defaultVolumeStorage
	}

	// Set PVC back for suggestion config.
	suggestionConfigData.PersistentVolumeClaimSpec = pvcSpec

	// Get PV from config only if it exists.
	if !equality.Semantic.DeepEqual(suggestionConfigData.PersistentVolumeSpec, corev1.PersistentVolumeSpec{}) {

		// Set PersistentVolumeReclaimPolicy to "Delete" to automatically delete PV once PVC is deleted.
		// Kubernetes doesn't allow to specify ownerReferences for the cluster-scoped
		// resources (which PV is) with namespace-scoped owner (which Suggestion is).
		suggestionConfigData.PersistentVolumeSpec.PersistentVolumeReclaimPolicy = corev1.PersistentVolumeReclaimDelete

	}

	return suggestionConfigData, nil
}

// GetEarlyStoppingConfigData gets the config data for the given early stopping algorithm name.
func GetEarlyStoppingConfigData(algorithmName string, client client.Client) (EarlyStoppingConfig, error) {
	configMap := &corev1.ConfigMap{}
	earlyStoppingConfigData := EarlyStoppingConfig{}
	err := client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		return EarlyStoppingConfig{}, err
	}

	// Try to find early stopping data in config map.
	config, ok := configMap.Data[consts.LabelEarlyStoppingTag]
	if !ok {
		return EarlyStoppingConfig{}, fmt.Errorf("failed to find early stopping config in ConfigMap: %s", consts.KatibConfigMapName)
	}

	// Parse early stopping data to map where key = algorithm name, value = EarlyStoppingConfig.
	earlyStoppingsConfig := map[string]EarlyStoppingConfig{}
	if err := json.Unmarshal([]byte(config), &earlyStoppingsConfig); err != nil {
		return EarlyStoppingConfig{}, err
	}

	// Try to find EarlyStoppingConfig for the algorithm.
	earlyStoppingConfigData, ok = earlyStoppingsConfig[algorithmName]
	if !ok {
		return EarlyStoppingConfig{}, fmt.Errorf("failed to find early stopping config for algorithm: %s in ConfigMap: %s", algorithmName, consts.KatibConfigMapName)
	}

	// Get image from config.
	image := earlyStoppingConfigData.Image
	if strings.TrimSpace(image) == "" {
		return EarlyStoppingConfig{}, fmt.Errorf("required value for image configuration of algorithm name: %s", algorithmName)
	}

	// Set Image Pull Policy.
	earlyStoppingConfigData.ImagePullPolicy = setImagePullPolicy(earlyStoppingConfigData.ImagePullPolicy)

	return earlyStoppingConfigData, nil
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
		return MetricsCollectorConfig{}, fmt.Errorf("failed to find metrics collector config in ConfigMap: %s", consts.KatibConfigMapName)
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
		return MetricsCollectorConfig{}, fmt.Errorf("failed to find metrics collector config for kind: %s in ConfigMap: %s", kind, consts.KatibConfigMapName)
	}

	// Get image from config
	image := metricsCollectorConfigData.Image
	if strings.TrimSpace(image) == "" {
		return MetricsCollectorConfig{}, fmt.Errorf("required value for image configuration of metrics collector kind: %s", kind)
	}

	// Set Image Pull Policy
	metricsCollectorConfigData.ImagePullPolicy = setImagePullPolicy(metricsCollectorConfigData.ImagePullPolicy)

	// Set resource requirements for metrics collector
	metricsCollectorConfigData.Resource = setResourceRequirements(metricsCollectorConfigData.Resource)

	return metricsCollectorConfigData, nil
}

func setImagePullPolicy(imagePullPolicy corev1.PullPolicy) corev1.PullPolicy {
	if imagePullPolicy != corev1.PullAlways && imagePullPolicy != corev1.PullIfNotPresent && imagePullPolicy != corev1.PullNever {
		return consts.DefaultImagePullPolicy
	}
	return imagePullPolicy
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

	// If user explicitly sets CPU value to -1, nuke it.
	if cpuLimit.Sign() == -1 {
		delete(configResource.Limits, corev1.ResourceCPU)
	}
	if cpuRequest.Sign() == -1 {
		delete(configResource.Requests, corev1.ResourceCPU)
	}

	// If user explicitly sets Memory value to -1, nuke it.
	if memLimit.Sign() == -1 {
		delete(configResource.Limits, corev1.ResourceMemory)
	}
	if memRequest.Sign() == -1 {
		delete(configResource.Requests, corev1.ResourceMemory)
	}

	// If user explicitly sets ephemeral-storage value to something negative, nuke it.
	// This enables compatibility with the GKE nodepool autoscalers, which cannot scale
	// pods which define ephemeral-storage resource constraints.
	if diskLimit.Sign() == -1 {
		delete(configResource.Limits, corev1.ResourceEphemeralStorage)
	}
	if diskRequest.Sign() == -1 {
		delete(configResource.Requests, corev1.ResourceEphemeralStorage)
	}

	return configResource
}
