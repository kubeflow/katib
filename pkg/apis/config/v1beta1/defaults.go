package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

const (
	DefaultExperimentSuggestionName = "default"
	DefaultMetricsAddr              = ":8080"
	DefaultHealthzAddr              = ":18080"
	DefaultLeaderElectionID         = "3fbc96e9.katib.kubeflow.org"
)

var (
	DefaultInjectSecurityContext       = false
	DefaultEnableGRPCProbeInSuggestion = true
	DefaultWebhookPort                 = 8443
	DefaultEnableLeaderElection        = false
	DefaultTrialResources              = []string{"Job.v1.batch"}
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
	// Set ExperimentSuggestionName.
	if initConfig.ControllerConfig.ExperimentSuggestionName == "" {
		initConfig.ControllerConfig.ExperimentSuggestionName = DefaultExperimentSuggestionName
	}
	// Set MetricsAddr.
	if initConfig.ControllerConfig.MetricsAddr == "" {
		initConfig.ControllerConfig.MetricsAddr = DefaultMetricsAddr
	}
	// Set HealthzAddr.
	if initConfig.ControllerConfig.HealthzAddr == "" {
		initConfig.ControllerConfig.HealthzAddr = DefaultHealthzAddr
	}
	// Set InjectSecurityContext.
	if initConfig.ControllerConfig.InjectSecurityContext == nil {
		initConfig.ControllerConfig.InjectSecurityContext = &DefaultInjectSecurityContext
	}
	// Set EnableGRPCProbeInSuggestion.
	if initConfig.ControllerConfig.EnableGRPCProbeInSuggestion == nil {
		initConfig.ControllerConfig.EnableGRPCProbeInSuggestion = &DefaultEnableGRPCProbeInSuggestion
	}
	// Set TrialResources.
	if len(initConfig.ControllerConfig.TrialResources) == 0 {
		initConfig.ControllerConfig.TrialResources = DefaultTrialResources
	}
	// Set WebhookPort.
	if initConfig.ControllerConfig.WebhookPort == nil {
		initConfig.ControllerConfig.WebhookPort = &DefaultWebhookPort
	}
	// Set EnableLeaderElection.
	if initConfig.ControllerConfig.EnableLeaderElection == nil {
		initConfig.ControllerConfig.EnableLeaderElection = &DefaultEnableLeaderElection
	}
	// Set LeaderElectionID.
	if initConfig.ControllerConfig.LeaderElectionID == "" {
		initConfig.ControllerConfig.LeaderElectionID = DefaultLeaderElectionID
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
			suggestionConfigs[i].VolumeMountPath = consts.DefaultContainerSuggestionVolumeMountPath
		}

		// Get persistent volume claim spec from config
		pvcSpec := suggestionConfigs[i].PersistentVolumeClaimSpec

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
