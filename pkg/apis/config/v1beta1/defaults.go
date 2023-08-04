package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// DefaultExperimentSuggestionName is the default name for the suggestions.
	DefaultExperimentSuggestionName = "default"
	// DefaultMetricsAddr is the default address for the prometheus metrics.
	DefaultMetricsAddr = ":8080"
	// DefaultHealthzAddr is the default address for the health probe.
	DefaultHealthzAddr = ":18080"
	// DefaultLeaderElectionID is the default LeaderElectionID for the controller.
	DefaultLeaderElectionID = "3fbc96e9.katib.kubeflow.org"
	// DefaultContainerSuggestionVolumeMountPath is the default mount path in suggestion container.
	DefaultContainerSuggestionVolumeMountPath = "/opt/katib/data"
	// DefaultSuggestionVolumeAccessMode is the default value for suggestion's volume access mode.
	DefaultSuggestionVolumeAccessMode = corev1.ReadWriteOnce
	// DefaultSuggestionVolumeStorage is the default value for suggestion's volume storage.
	DefaultSuggestionVolumeStorage = "1Gi"
	// DefaultImagePullPolicy is the default value for image pull policy.
	DefaultImagePullPolicy = corev1.PullIfNotPresent
	// DefaultCPULimit is the default value for CPU limit.
	DefaultCPULimit = "500m"
	// DefaultCPURequest is the default value for CPU request.
	DefaultCPURequest = "50m"
	// DefaultMemLimit is the default value for memory limit.
	DefaultMemLimit = "100Mi"
	// DefaultMemRequest is the default value for memory request.
	DefaultMemRequest = "10Mi"
	// DefaultDiskLimit is the default value for disk limit.
	DefaultDiskLimit = "5Gi"
	// DefaultDiskRequest is the default value for disk request.
	DefaultDiskRequest = "500Mi"
	// DefaultWebhookServiceName is the default service name for the admission webhooks.
	DefaultWebhookServiceName = "katib-controller"
	// DefaultWebhookSecretName is the default secret name to save the certs for the admission webhooks.
	DefaultWebhookSecretName = "katib-webhook-cert"
)

var (
	// DefaultEnableGRPCProbeInSuggestion is the default value whether enable to gRPC probe in suggestions.
	DefaultEnableGRPCProbeInSuggestion = true
	// DefaultWebhookPort is the default port for the admission webhook.
	DefaultWebhookPort = 8443
	// DefaultTrialResources is the default resource which can be used as a trial template.
	DefaultTrialResources = []string{"Job.v1.batch"}
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&KatibConfig{}, func(obj interface{}) {
		SetDefaults_KatibConfig(obj.(*KatibConfig))
	})
	return nil
}

func SetDefaults_KatibConfig(cfg *KatibConfig) {
	if cfg == nil {
		return
	}
	setInitConfig(&cfg.InitConfig)
	setRuntimeConfig(&cfg.RuntimeConfig)
}

func setInitConfig(initConfig *InitConfig) {
	setControllerConfig(&initConfig.ControllerConfig)
	setCertGeneratorConfig(&initConfig.CertGeneratorConfig)
}

func setControllerConfig(controllerConfig *ControllerConfig) {
	// Set ExperimentSuggestionName.
	if controllerConfig.ExperimentSuggestionName == "" {
		controllerConfig.ExperimentSuggestionName = DefaultExperimentSuggestionName
	}
	// Set MetricsAddr.
	if controllerConfig.MetricsAddr == "" {
		controllerConfig.MetricsAddr = DefaultMetricsAddr
	}
	// Set HealthzAddr.
	if controllerConfig.HealthzAddr == "" {
		controllerConfig.HealthzAddr = DefaultHealthzAddr
	}
	// Set EnableGRPCProbeInSuggestion.
	if controllerConfig.EnableGRPCProbeInSuggestion == nil {
		controllerConfig.EnableGRPCProbeInSuggestion = &DefaultEnableGRPCProbeInSuggestion
	}
	// Set TrialResources.
	if len(controllerConfig.TrialResources) == 0 {
		controllerConfig.TrialResources = DefaultTrialResources
	}
	// Set WebhookPort.
	if controllerConfig.WebhookPort == nil {
		controllerConfig.WebhookPort = &DefaultWebhookPort
	}
	// Set LeaderElectionID.
	if controllerConfig.LeaderElectionID == "" {
		controllerConfig.LeaderElectionID = DefaultLeaderElectionID
	}
}

func setCertGeneratorConfig(certGeneratorConfig *CertGeneratorConfig) {
	if len(certGeneratorConfig.WebhookServiceName) != 0 || len(certGeneratorConfig.WebhookSecretName) != 0 {
		certGeneratorConfig.Enable = true
	}
	if certGeneratorConfig.Enable && len(certGeneratorConfig.WebhookServiceName) == 0 {
		certGeneratorConfig.WebhookServiceName = DefaultWebhookServiceName
	}
	if certGeneratorConfig.Enable && len(certGeneratorConfig.WebhookSecretName) == 0 {
		certGeneratorConfig.WebhookSecretName = DefaultWebhookSecretName
	}
}

func setRuntimeConfig(runtimeConfig *RuntimeConfig) {
	setSuggestionConfigs(runtimeConfig.SuggestionConfigs)
	setMetricsCollectorConfigs(runtimeConfig.MetricsCollectorConfigs)
	setEarlyStoppingConfigs(runtimeConfig.EarlyStoppingConfigs)
}

func setSuggestionConfigs(suggestionConfigs []SuggestionConfig) {
	for i := range suggestionConfigs {
		// Set Image Pull Policy
		suggestionConfigs[i].ImagePullPolicy = setImagePullPolicy(suggestionConfigs[i].ImagePullPolicy)

		// Set resource requirements for suggestion
		suggestionConfigs[i].Resources = setResourceRequirements(suggestionConfigs[i].Resources)

		// Set default suggestion container volume mount path
		if suggestionConfigs[i].VolumeMountPath == "" {
			suggestionConfigs[i].VolumeMountPath = DefaultContainerSuggestionVolumeMountPath
		}

		// Get persistent volume claim spec from config
		pvcSpec := suggestionConfigs[i].PersistentVolumeClaimSpec

		// Set default access modes
		if len(pvcSpec.AccessModes) == 0 {
			pvcSpec.AccessModes = []corev1.PersistentVolumeAccessMode{
				DefaultSuggestionVolumeAccessMode,
			}
		}

		// Set default resources
		if len(pvcSpec.Resources.Requests) == 0 {
			defaultVolumeStorage, _ := resource.ParseQuantity(DefaultSuggestionVolumeStorage)
			pvcSpec.Resources.Requests = make(map[corev1.ResourceName]resource.Quantity)
			pvcSpec.Resources.Requests[corev1.ResourceStorage] = defaultVolumeStorage
		}

		// Set PVC back for suggestion config.
		suggestionConfigs[i].PersistentVolumeClaimSpec = pvcSpec

		// Get PV from config only if it exists.
		if !equality.Semantic.DeepEqual(suggestionConfigs[i].PersistentVolumeSpec, corev1.PersistentVolumeSpec{}) {

			// Set PersistentVolumeReclaimPolicy to "Delete" to automatically delete PV once PVC is deleted.
			// Kubernetes doesn't allow to specify ownerReferences for the cluster-scoped
			// resources (which PV is) with namespace-scoped owner (which Suggestion is).
			suggestionConfigs[i].PersistentVolumeSpec.PersistentVolumeReclaimPolicy = corev1.PersistentVolumeReclaimDelete
		}
	}
}

func setMetricsCollectorConfigs(metricsCollectorConfigs []MetricsCollectorConfig) {
	for i := range metricsCollectorConfigs {
		// Set Image Pull Policy
		metricsCollectorConfigs[i].ImagePullPolicy = setImagePullPolicy(metricsCollectorConfigs[i].ImagePullPolicy)

		// Set resource requirements for metrics collector
		metricsCollectorConfigs[i].Resource = setResourceRequirements(metricsCollectorConfigs[i].Resource)
	}
}

func setEarlyStoppingConfigs(earlyStoppingConfigs []EarlyStoppingConfig) {
	for i := range earlyStoppingConfigs {
		// Set Image Pull Policy.
		earlyStoppingConfigs[i].ImagePullPolicy = setImagePullPolicy(earlyStoppingConfigs[i].ImagePullPolicy)

		// Set resource requirements
		earlyStoppingConfigs[i].Resource = setResourceRequirements(earlyStoppingConfigs[i].Resource)
	}
}

func setImagePullPolicy(imagePullPolicy corev1.PullPolicy) corev1.PullPolicy {
	if imagePullPolicy != corev1.PullAlways && imagePullPolicy != corev1.PullIfNotPresent && imagePullPolicy != corev1.PullNever {
		return DefaultImagePullPolicy
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
		defaultCPURequest, _ := resource.ParseQuantity(DefaultCPURequest)
		configResource.Requests[corev1.ResourceCPU] = defaultCPURequest
	}
	if memRequest.IsZero() {
		defaultMemRequest, _ := resource.ParseQuantity(DefaultMemRequest)
		configResource.Requests[corev1.ResourceMemory] = defaultMemRequest
	}
	if diskRequest.IsZero() {
		defaultDiskRequest, _ := resource.ParseQuantity(DefaultDiskRequest)
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
		defaultCPULimit, _ := resource.ParseQuantity(DefaultCPULimit)
		configResource.Limits[corev1.ResourceCPU] = defaultCPULimit
	}
	if memLimit.IsZero() {
		defaultMemLimit, _ := resource.ParseQuantity(DefaultMemLimit)
		configResource.Limits[corev1.ResourceMemory] = defaultMemLimit
	}
	if diskLimit.IsZero() {
		defaultDiskLimit, _ := resource.ParseQuantity(DefaultDiskLimit)
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
