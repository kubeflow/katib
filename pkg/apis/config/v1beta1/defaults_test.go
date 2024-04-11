package v1beta1

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
)

func TestSetSuggestionConfigs(t *testing.T) {
	const testAlgorithmName = "test-suggestion"

	cases := map[string]struct {
		config     []SuggestionConfig
		wantConfig []SuggestionConfig
	}{
		"All parameters correctly are specified": {
			config: func() []SuggestionConfig {
				suggestionConfig := newFakeSuggestionConfig(testAlgorithmName)
				suggestionConfig.ImagePullPolicy = corev1.PullAlways
				suggestionConfig.Resources = *newFakeCustomResourceRequirements()
				return []SuggestionConfig{*suggestionConfig}
			}(),
			wantConfig: func() []SuggestionConfig {
				c := newFakeSuggestionConfig(testAlgorithmName)
				c.ImagePullPolicy = corev1.PullAlways
				c.Resources = *newFakeCustomResourceRequirements()
				return []SuggestionConfig{*c}
			}(),
		},
		fmt.Sprintf("GetSuggestionConfigData sets %s to imagePullPolicy", DefaultImagePullPolicy): {
			config: func() []SuggestionConfig {
				suggestion := newFakeSuggestionConfig(testAlgorithmName)
				suggestion.ImagePullPolicy = ""
				return []SuggestionConfig{*suggestion}
			}(),
			wantConfig: []SuggestionConfig{*newFakeSuggestionConfig(testAlgorithmName)},
		},
		"GetSuggestionConfigData sets resource.requests and resource.limits for the suggestion service": {
			config: func() []SuggestionConfig {
				suggestion := newFakeSuggestionConfig(testAlgorithmName)
				suggestion.Resources = corev1.ResourceRequirements{}
				return []SuggestionConfig{*suggestion}
			}(),
			wantConfig: []SuggestionConfig{*newFakeSuggestionConfig(testAlgorithmName)},
		},
		fmt.Sprintf("GetSuggestionConfigData sets %s to volumeMountPath", DefaultContainerSuggestionVolumeMountPath): {
			config: func() []SuggestionConfig {
				suggestion := newFakeSuggestionConfig(testAlgorithmName)
				suggestion.VolumeMountPath = ""
				return []SuggestionConfig{*suggestion}
			}(),
			wantConfig: []SuggestionConfig{*newFakeSuggestionConfig(testAlgorithmName)},
		},
		"GetSuggestionConfigData sets accessMode and resource.requests for PVC": {
			config: func() []SuggestionConfig {
				suggestion := newFakeSuggestionConfig(testAlgorithmName)
				suggestion.PersistentVolumeClaimSpec = corev1.PersistentVolumeClaimSpec{}
				return []SuggestionConfig{*suggestion}
			}(),
			wantConfig: []SuggestionConfig{*newFakeSuggestionConfig(testAlgorithmName)},
		},
		fmt.Sprintf("GetSuggestionConfigData does not set %s to persistentVolumeReclaimPolicy", corev1.PersistentVolumeReclaimDelete): {
			config: func() []SuggestionConfig {
				suggestion := newFakeSuggestionConfig(testAlgorithmName)
				suggestion.PersistentVolumeSpec = corev1.PersistentVolumeSpec{}
				return []SuggestionConfig{*suggestion}
			}(),
			wantConfig: func() []SuggestionConfig {
				c := newFakeSuggestionConfig(testAlgorithmName)
				c.PersistentVolumeSpec = corev1.PersistentVolumeSpec{}
				return []SuggestionConfig{*c}
			}(),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			kc := &KatibConfig{
				RuntimeConfig: RuntimeConfig{
					SuggestionConfigs: tc.config,
				},
			}
			SetDefaults_KatibConfig(kc)
			if diff := cmp.Diff(tc.wantConfig, kc.RuntimeConfig.SuggestionConfigs); len(diff) != 0 {
				t.Errorf("Unexpected SuggestionConfigs (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestSetEarlyStoppingConfigs(t *testing.T) {
	const testAlgorithmName = "test-early-stopping"

	cases := map[string]struct {
		config     []EarlyStoppingConfig
		wantConfig []EarlyStoppingConfig
	}{
		"All parameters correctly are specified": {
			config: func() []EarlyStoppingConfig {
				config := newFakeEarlyStoppingConfig(testAlgorithmName)
				config.ImagePullPolicy = corev1.PullIfNotPresent
				return []EarlyStoppingConfig{*config}
			}(),
			wantConfig: func() []EarlyStoppingConfig {
				c := newFakeEarlyStoppingConfig(testAlgorithmName)
				c.ImagePullPolicy = corev1.PullIfNotPresent
				return []EarlyStoppingConfig{*c}
			}(),
		},
		fmt.Sprintf("GetEarlyStoppingConfigData sets %s to imagePullPolicy", DefaultImagePullPolicy): {
			config: func() []EarlyStoppingConfig {
				c := newFakeEarlyStoppingConfig(testAlgorithmName)
				c.ImagePullPolicy = ""
				return []EarlyStoppingConfig{*c}
			}(),
			wantConfig: []EarlyStoppingConfig{*newFakeEarlyStoppingConfig(testAlgorithmName)},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			kc := &KatibConfig{
				RuntimeConfig: RuntimeConfig{
					EarlyStoppingConfigs: tc.config,
				},
			}
			SetDefaults_KatibConfig(kc)
			if diff := cmp.Diff(tc.wantConfig, kc.RuntimeConfig.EarlyStoppingConfigs); len(diff) != 0 {
				t.Errorf("Unexpected EarlyStoppingConfigs (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestSetMetricsCollectorConfigs(t *testing.T) {
	const testCollectorKind commonv1beta1.CollectorKind = "testCollector"
	nukeResource, _ := resource.ParseQuantity("-1")
	nukeResourceRequirements := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:              nukeResource,
		corev1.ResourceMemory:           nukeResource,
		corev1.ResourceEphemeralStorage: nukeResource,
	}

	cases := map[string]struct {
		config, wantConfig []MetricsCollectorConfig
	}{
		"All parameters correctly are specified": {
			config: func() []MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig(testCollectorKind)
				c.ImagePullPolicy = corev1.PullNever
				return []MetricsCollectorConfig{*c}
			}(),
			wantConfig: func() []MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig(testCollectorKind)
				c.ImagePullPolicy = corev1.PullNever
				return []MetricsCollectorConfig{*c}
			}(),
		},
		fmt.Sprintf("GetMetricsConfigData sets %s to imagePullPolicy", DefaultImagePullPolicy): {
			config: func() []MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig(testCollectorKind)
				c.ImagePullPolicy = ""
				return []MetricsCollectorConfig{*c}
			}(),
			wantConfig: []MetricsCollectorConfig{*newFakeMetricsCollectorConfig(testCollectorKind)},
		},
		"GetMetricsConfigData nukes resource.requests and resource.limits for the metrics collector": {
			config: func() []MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig(testCollectorKind)
				c.Resource = corev1.ResourceRequirements{
					Requests: nukeResourceRequirements,
					Limits:   nukeResourceRequirements,
				}
				return []MetricsCollectorConfig{*c}
			}(),
			wantConfig: func() []MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig(testCollectorKind)
				c.Resource = corev1.ResourceRequirements{
					Requests: map[corev1.ResourceName]resource.Quantity{},
					Limits:   map[corev1.ResourceName]resource.Quantity{},
				}
				return []MetricsCollectorConfig{*c}
			}(),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			kc := &KatibConfig{
				RuntimeConfig: RuntimeConfig{
					MetricsCollectorConfigs: tc.config,
				},
			}
			SetDefaults_KatibConfig(kc)
			if diff := cmp.Diff(tc.wantConfig, kc.RuntimeConfig.MetricsCollectorConfigs); len(diff) != 0 {
				t.Errorf("Unexpected MetricsCollectorConfigs (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestSetControllerConfig(t *testing.T) {
	disableGRPCProbeInSuggestion := false
	customizedWebhookPort := 18443

	cases := map[string]struct {
		config     ControllerConfig
		wantConfig ControllerConfig
	}{
		"All parameters correctly are specified": {
			config: ControllerConfig{
				ExperimentSuggestionName:    "test",
				MetricsAddr:                 ":8081",
				HealthzAddr:                 ":18081",
				InjectSecurityContext:       true,
				EnableGRPCProbeInSuggestion: &disableGRPCProbeInSuggestion,
				TrialResources: []string{
					"Job.v1.batch",
					"TFJob.v1.kubeflow.org",
				},
				WebhookPort:          &customizedWebhookPort,
				EnableLeaderElection: true,
				LeaderElectionID:     "xyz0123",
			},
			wantConfig: ControllerConfig{
				ExperimentSuggestionName:    "test",
				MetricsAddr:                 ":8081",
				HealthzAddr:                 ":18081",
				InjectSecurityContext:       true,
				EnableGRPCProbeInSuggestion: &disableGRPCProbeInSuggestion,
				TrialResources: []string{
					"Job.v1.batch",
					"TFJob.v1.kubeflow.org",
				},
				WebhookPort:          &customizedWebhookPort,
				EnableLeaderElection: true,
				LeaderElectionID:     "xyz0123",
			},
		},
		"ControllerConfig is empty": {
			config: ControllerConfig{},
			wantConfig: ControllerConfig{
				ExperimentSuggestionName:    DefaultExperimentSuggestionName,
				MetricsAddr:                 DefaultMetricsAddr,
				HealthzAddr:                 DefaultHealthzAddr,
				EnableGRPCProbeInSuggestion: &DefaultEnableGRPCProbeInSuggestion,
				TrialResources:              DefaultTrialResources,
				WebhookPort:                 &DefaultWebhookPort,
				LeaderElectionID:            DefaultLeaderElectionID,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			kc := &KatibConfig{
				InitConfig: InitConfig{
					ControllerConfig: tc.config,
				},
			}
			SetDefaults_KatibConfig(kc)
			if diff := cmp.Diff(tc.wantConfig, kc.InitConfig.ControllerConfig); len(diff) != 0 {
				t.Errorf("Unexpected ControllerConfig (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestSetCertGeneratorConfig(t *testing.T) {
	cases := map[string]struct {
		config     CertGeneratorConfig
		wantConfig CertGeneratorConfig
	}{
		"All parameters correctly are specified": {
			config: CertGeneratorConfig{
				Enable:             true,
				WebhookServiceName: "test",
				WebhookSecretName:  "katib-test",
			},
			wantConfig: CertGeneratorConfig{
				Enable:             true,
				WebhookServiceName: "test",
				WebhookSecretName:  "katib-test",
			},
		},
		"CertGeneratorConfig is empty": {
			config:     CertGeneratorConfig{},
			wantConfig: CertGeneratorConfig{},
		},
		"Enable is true and serviceName is empty": {
			config: CertGeneratorConfig{
				Enable: true,
			},
			wantConfig: CertGeneratorConfig{
				Enable:             true,
				WebhookServiceName: DefaultWebhookServiceName,
				WebhookSecretName:  DefaultWebhookSecretName,
			},
		},
		"cert-generator is forcefully enabled due to set webhookSecretName": {
			config: CertGeneratorConfig{
				WebhookSecretName: "katib-test",
			},
			wantConfig: CertGeneratorConfig{
				Enable:             true,
				WebhookServiceName: DefaultWebhookServiceName,
				WebhookSecretName:  "katib-test",
			},
		},
		"cert-generator is forcefully enabled due to set webhookServiceName": {
			config: CertGeneratorConfig{
				WebhookServiceName: "katib-test",
			},
			wantConfig: CertGeneratorConfig{
				Enable:             true,
				WebhookServiceName: "katib-test",
				WebhookSecretName:  DefaultWebhookSecretName,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			kc := &KatibConfig{
				InitConfig: InitConfig{
					CertGeneratorConfig: tc.config,
				},
			}
			SetDefaults_KatibConfig(kc)
			if diff := cmp.Diff(tc.wantConfig, kc.InitConfig.CertGeneratorConfig); len(diff) != 0 {
				t.Errorf("Unexpected CertGeneratorConfig (-want,+got):\n%s", diff)
			}
		})
	}
}

func newFakeSuggestionConfig(algorithmName string) *SuggestionConfig {
	defaultVolumeStorage, _ := resource.ParseQuantity(DefaultSuggestionVolumeStorage)

	return &SuggestionConfig{
		AlgorithmName: algorithmName,
		Container: corev1.Container{
			Image:           "suggestion-image",
			ImagePullPolicy: DefaultImagePullPolicy,
			Resources:       *setFakeResourceRequirements(),
		},
		VolumeMountPath: DefaultContainerSuggestionVolumeMountPath,
		PersistentVolumeClaimSpec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				DefaultSuggestionVolumeAccessMode,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: defaultVolumeStorage,
				},
			},
		},
		PersistentVolumeSpec: corev1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
		},
	}
}

func newFakeEarlyStoppingConfig(algorithmName string) *EarlyStoppingConfig {
	return &EarlyStoppingConfig{
		AlgorithmName:   algorithmName,
		Image:           "early-stopping-image",
		ImagePullPolicy: DefaultImagePullPolicy,
		Resource:        *setFakeResourceRequirements(),
	}
}

func newFakeMetricsCollectorConfig(collectorKind commonv1beta1.CollectorKind) *MetricsCollectorConfig {
	return &MetricsCollectorConfig{
		CollectorKind:   string(collectorKind),
		Image:           "metrics-collector-image",
		ImagePullPolicy: DefaultImagePullPolicy,
		Resource:        *setFakeResourceRequirements(),
	}
}

func setFakeResourceRequirements() *corev1.ResourceRequirements {
	defaultCPURequest, _ := resource.ParseQuantity(DefaultCPURequest)
	defaultMemoryRequest, _ := resource.ParseQuantity(DefaultMemRequest)
	defaultEphemeralStorageRequest, _ := resource.ParseQuantity(DefaultDiskRequest)

	defaultCPULimit, _ := resource.ParseQuantity(DefaultCPULimit)
	defaultMemoryLimit, _ := resource.ParseQuantity(DefaultMemLimit)
	defaultEphemeralStorageLimit, _ := resource.ParseQuantity(DefaultDiskLimit)

	return &corev1.ResourceRequirements{
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:              defaultCPURequest,
			corev1.ResourceMemory:           defaultMemoryRequest,
			corev1.ResourceEphemeralStorage: defaultEphemeralStorageRequest,
		},
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:              defaultCPULimit,
			corev1.ResourceMemory:           defaultMemoryLimit,
			corev1.ResourceEphemeralStorage: defaultEphemeralStorageLimit,
		},
	}
}

func newFakeCustomResourceRequirements() *corev1.ResourceRequirements {
	customCPURequest, _ := resource.ParseQuantity("25m")
	customMemoryRequest, _ := resource.ParseQuantity("200Mi")
	customEphemeralStorageRequest, _ := resource.ParseQuantity("550Mi")

	customCPULimit, _ := resource.ParseQuantity("250m")
	customMemoryLimit, _ := resource.ParseQuantity("2Gi")
	customEphemeralStorageLimit, _ := resource.ParseQuantity("15Gi")

	return &corev1.ResourceRequirements{
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:              customCPURequest,
			corev1.ResourceMemory:           customMemoryRequest,
			corev1.ResourceEphemeralStorage: customEphemeralStorageRequest,
		},
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:              customCPULimit,
			corev1.ResourceMemory:           customMemoryLimit,
			corev1.ResourceEphemeralStorage: customEphemeralStorageLimit,
		},
	}
}
