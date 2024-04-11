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
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/yaml"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

func TestTrialResourcesToGVKs(t *testing.T) {
	cases := map[string]struct {
		trialResources []string
		wantGVKs       []schema.GroupVersionKind
		wantError      error
	}{
		"All GVKs are appropriate": {
			trialResources: []string{
				"Job.v1.batch",
				"TFJob.v1.kubeflow.org",
			},
			wantGVKs: []schema.GroupVersionKind{
				{Group: "batch", Version: "v1", Kind: "Job"},
				{Group: "kubeflow.org", Version: "v1", Kind: "TFJob"},
			},
		},
		"TrialResources are empty": {
			trialResources: []string{},
			wantError:      ErrTrialResourcesAreEmpty,
		},
		"GVK with invalid schema": {
			trialResources: []string{
				"invalid;;invalid",
			},
			wantError: ErrInvalidGVKFormat,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := TrialResourcesToGVKs(tc.trialResources)
			if diff := cmp.Diff(tc.wantError, err, cmpopts.EquateErrors()); len(diff) != 0 {
				t.Errorf("Unexpected error from TrialResourcesToGVKs (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantGVKs, got); len(diff) != 0 {
				t.Errorf("Unexpected gvks from TrialResourcesToGVKs (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestGetSuggestionConfigData(t *testing.T) {
	const testAlgorithmName = "test-suggestion"
	scm := runtime.NewScheme()
	if err := configv1beta1.AddToScheme(scm); err != nil {
		t.Fatal(err)
	}
	if err := clientgoscheme.AddToScheme(scm); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		testDescription    string
		katibConfig        *configv1beta1.KatibConfig
		expected           *configv1beta1.SuggestionConfig
		inputAlgorithmName string
		err                bool
	}{
		{
			testDescription: "All parameters correctly are specified",
			katibConfig: func() *configv1beta1.KatibConfig {
				kc := &configv1beta1.KatibConfig{
					TypeMeta: metav1.TypeMeta{
						Kind:       "KatibConfig",
						APIVersion: "config.kubeflow.org/v1beta1",
					},
					RuntimeConfig: configv1beta1.RuntimeConfig{
						SuggestionConfigs: []configv1beta1.SuggestionConfig{
							*newFakeSuggestionConfig(testAlgorithmName),
						},
					},
				}
				kc.RuntimeConfig.SuggestionConfigs[0].ImagePullPolicy = corev1.PullAlways
				kc.RuntimeConfig.SuggestionConfigs[0].Resources = *newFakeCustomResourceRequirements()
				return kc
			}(),
			expected: func() *configv1beta1.SuggestionConfig {
				c := newFakeSuggestionConfig(testAlgorithmName)
				c.ImagePullPolicy = corev1.PullAlways
				c.Resources = *newFakeCustomResourceRequirements()
				return c
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: "There is not katib-config.",
			katibConfig:     nil,
			err:             true,
		},
		{
			testDescription: "There is not runtime.suggestions field in katib-config configMap",
			katibConfig:     &configv1beta1.KatibConfig{},
			err:             true,
		},
		{
			testDescription: "There is not the AlgorithmName",
			katibConfig: &configv1beta1.KatibConfig{
				RuntimeConfig: configv1beta1.RuntimeConfig{
					SuggestionConfigs: []configv1beta1.SuggestionConfig{
						*newFakeSuggestionConfig(testAlgorithmName),
					},
				},
			},
			inputAlgorithmName: "invalid-algorithm-name",
			err:                true,
		},
		{
			testDescription: "Image filed is empty in katib-config configMap",
			katibConfig: func() *configv1beta1.KatibConfig {
				kc := &configv1beta1.KatibConfig{
					RuntimeConfig: configv1beta1.RuntimeConfig{
						SuggestionConfigs: []configv1beta1.SuggestionConfig{
							*newFakeSuggestionConfig(testAlgorithmName),
						},
					},
				}
				kc.RuntimeConfig.SuggestionConfigs[0].Image = ""
				return kc
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			fakeKubeClient := newFakeKubeClient(scm, newFakeKatibConfigMap(tt.katibConfig))
			actual, err := GetSuggestionConfigData(tt.inputAlgorithmName, fakeKubeClient)
			if (err != nil) != tt.err {
				t.Errorf("want error: %v, actual: %v", tt.err, err)
			} else if tt.expected != nil {
				if diff := cmp.Diff(*tt.expected, actual); len(diff) != 0 {
					t.Logf("katibConfig: %v", tt.katibConfig)
					t.Errorf("Generated SuggestionConfig is invalid. (-want,+got):\n%s", diff)
				}
			}
		})
	}
}

func TestGetEarlyStoppingConfigData(t *testing.T) {
	const testAlgorithmName = "test-early-stopping"
	scm := runtime.NewScheme()
	if err := configv1beta1.AddToScheme(scm); err != nil {
		t.Fatal(err)
	}
	if err := clientgoscheme.AddToScheme(scm); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		testDescription    string
		katibConfig        *configv1beta1.KatibConfig
		expected           *configv1beta1.EarlyStoppingConfig
		inputAlgorithmName string
		err                bool
	}{
		{
			testDescription: "All parameters correctly are specified",
			katibConfig: func() *configv1beta1.KatibConfig {
				kc := &configv1beta1.KatibConfig{
					RuntimeConfig: configv1beta1.RuntimeConfig{
						EarlyStoppingConfigs: []configv1beta1.EarlyStoppingConfig{
							*newFakeEarlyStoppingConfig(testAlgorithmName),
						},
					},
				}
				kc.RuntimeConfig.EarlyStoppingConfigs[0].ImagePullPolicy = corev1.PullIfNotPresent
				return kc
			}(),
			expected: func() *configv1beta1.EarlyStoppingConfig {
				c := newFakeEarlyStoppingConfig(testAlgorithmName)
				c.ImagePullPolicy = corev1.PullIfNotPresent
				return c
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: "There is not katib-config.",
			katibConfig:     nil,
			err:             true,
		},
		{
			testDescription: "There is not runtime.earlyStoppings field in katib-config configMap",
			katibConfig:     &configv1beta1.KatibConfig{},
			err:             true,
		},
		{
			testDescription: "There is not the AlgorithmName",
			katibConfig: &configv1beta1.KatibConfig{
				RuntimeConfig: configv1beta1.RuntimeConfig{
					EarlyStoppingConfigs: []configv1beta1.EarlyStoppingConfig{
						*newFakeEarlyStoppingConfig(testAlgorithmName),
					},
				},
			},
			inputAlgorithmName: "invalid-algorithm-name",
			err:                true,
		},
		{
			testDescription: "Image filed is empty in katib-config configMap",
			katibConfig: func() *configv1beta1.KatibConfig {
				kc := &configv1beta1.KatibConfig{
					RuntimeConfig: configv1beta1.RuntimeConfig{
						EarlyStoppingConfigs: []configv1beta1.EarlyStoppingConfig{
							*newFakeEarlyStoppingConfig(testAlgorithmName),
						},
					},
				}
				kc.RuntimeConfig.EarlyStoppingConfigs[0].Image = ""
				return kc
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			fakeKubeClient := newFakeKubeClient(scm, newFakeKatibConfigMap(tt.katibConfig))
			actual, err := GetEarlyStoppingConfigData(tt.inputAlgorithmName, fakeKubeClient)
			if (err != nil) != tt.err {
				t.Errorf("want error: %v, actual: %v", tt.err, err)
			} else if tt.expected != nil {
				if diff := cmp.Diff(*tt.expected, actual); len(diff) != 0 {
					t.Errorf("Generated EarlyStoppingConfig is invalid. (-want,+got):\n%s", diff)
				}
			}
		})
	}
}

func TestGetMetricsCollectorConfigData(t *testing.T) {
	const (
		invalidCollectorKind commonv1beta1.CollectorKind = "invalidCollector"
		testCollectorKind    commonv1beta1.CollectorKind = "testCollector"
	)
	scm := runtime.NewScheme()
	if err := configv1beta1.AddToScheme(scm); err != nil {
		t.Fatal(err)
	}
	if err := clientgoscheme.AddToScheme(scm); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		testDescription    string
		katibConfig        *configv1beta1.KatibConfig
		expected           *configv1beta1.MetricsCollectorConfig
		inputCollectorKind commonv1beta1.CollectorKind
		err                bool
	}{
		{
			testDescription: "All parameters correctly are specified",
			katibConfig: func() *configv1beta1.KatibConfig {
				kc := &configv1beta1.KatibConfig{
					RuntimeConfig: configv1beta1.RuntimeConfig{
						MetricsCollectorConfigs: []configv1beta1.MetricsCollectorConfig{
							*newFakeMetricsCollectorConfig(testCollectorKind),
						},
					},
				}
				kc.RuntimeConfig.MetricsCollectorConfigs[0].ImagePullPolicy = corev1.PullNever
				return kc
			}(),
			expected: func() *configv1beta1.MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig(testCollectorKind)
				c.ImagePullPolicy = corev1.PullNever
				return c
			}(),
			inputCollectorKind: testCollectorKind,
			err:                false,
		},
		{
			testDescription: "There is not katib-config.",
			katibConfig:     nil,
			err:             true,
		},
		{
			testDescription: "There is not runtime.metricsCollectors field in katib-config configMap",
			katibConfig:     &configv1beta1.KatibConfig{},
			err:             true,
		},
		{
			testDescription: "There is not the cKind",
			katibConfig: &configv1beta1.KatibConfig{
				RuntimeConfig: configv1beta1.RuntimeConfig{
					MetricsCollectorConfigs: []configv1beta1.MetricsCollectorConfig{
						*newFakeMetricsCollectorConfig(testCollectorKind),
					},
				},
			},
			inputCollectorKind: invalidCollectorKind,
			err:                true,
		},
		{
			testDescription: "Image filed is empty in katib-config configMap",
			katibConfig: func() *configv1beta1.KatibConfig {
				kc := &configv1beta1.KatibConfig{
					RuntimeConfig: configv1beta1.RuntimeConfig{
						MetricsCollectorConfigs: []configv1beta1.MetricsCollectorConfig{
							*newFakeMetricsCollectorConfig(testCollectorKind),
						},
					},
				}
				kc.RuntimeConfig.MetricsCollectorConfigs[0].Image = ""
				return kc
			}(),
			inputCollectorKind: testCollectorKind,
			err:                true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			fakeKubeClient := newFakeKubeClient(scm, newFakeKatibConfigMap(tt.katibConfig))
			actual, err := GetMetricsCollectorConfigData(tt.inputCollectorKind, fakeKubeClient)
			if (err != nil) != tt.err {
				t.Errorf("want error: %v, actual: %v", tt.err, err)
			} else if tt.expected != nil {
				if diff := cmp.Diff(*tt.expected, actual); len(diff) != 0 {
					t.Errorf("Generated MetricsCollectorConfig is invalid. (-want,+got):\n%s", diff)
				}
			}
		})
	}
}

func TestGetInitConfigData(t *testing.T) {
	diableGRPCProbeInSuggestion := false
	customizedWebhookPort := 18443

	tmpDir, err := os.MkdirTemp("", "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	scm := runtime.NewScheme()
	if err = configv1beta1.AddToScheme(scm); err != nil {
		t.Fatal(err)
	}

	fullInitConfig := filepath.Join(tmpDir, "fullInitConfig.yaml")
	if err = os.WriteFile(fullInitConfig, []byte(`
apiVersion: config.kubeflow.org/v1beta1
kind: KatibConfig
init:
  certGenerator:
    enable: true
    webhookServiceName: katib-test
    webhookSecretName: katib-test-secret
  controller:
    experimentSuggestionName: test
    metricsAddr: :8081
    healthzAddr: :18081
    injectSecurityContext: true
    enableGRPCProbeInSuggestion: false
    trialResources:
    - Job.v1.batch
    - TFJob.v1.kubeflow.org
    - PyTorchJob.v1.kubeflow.org
    - MPIJob.v1.kubeflow.org
    - XGBoostJob.v1.kubeflow.org
    - MXJob.v1.kubeflow.org
    webhookPort: 18443
    enableLeaderElection: true
    leaderElectionID: xyz0123
runtime:
  suggestions:
  - algorithmName: random
    image: docker.io/kubeflowkatib/suggestion-hyperopt:latest
`), os.FileMode(0600)); err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		katibConfigFile    string
		wantInitConfigData configv1beta1.InitConfig
		wantError          error
	}{
		"KatibConfigFile is empty": {
			katibConfigFile: "",
			wantInitConfigData: configv1beta1.InitConfig{
				ControllerConfig: configv1beta1.ControllerConfig{
					ExperimentSuggestionName:    configv1beta1.DefaultExperimentSuggestionName,
					MetricsAddr:                 configv1beta1.DefaultMetricsAddr,
					HealthzAddr:                 configv1beta1.DefaultHealthzAddr,
					EnableGRPCProbeInSuggestion: &configv1beta1.DefaultEnableGRPCProbeInSuggestion,
					TrialResources:              configv1beta1.DefaultTrialResources,
					WebhookPort:                 &configv1beta1.DefaultWebhookPort,
					LeaderElectionID:            configv1beta1.DefaultLeaderElectionID,
				},
			},
		},
		"invalid katibConfigFile": {
			katibConfigFile: "invalid",
			wantError:       ErrKatibConfigNil,
		},
		"full init config": {
			katibConfigFile: fullInitConfig,
			wantInitConfigData: configv1beta1.InitConfig{
				CertGeneratorConfig: configv1beta1.CertGeneratorConfig{
					Enable:             true,
					WebhookServiceName: "katib-test",
					WebhookSecretName:  "katib-test-secret",
				},
				ControllerConfig: configv1beta1.ControllerConfig{
					ExperimentSuggestionName:    "test",
					MetricsAddr:                 ":8081",
					HealthzAddr:                 ":18081",
					InjectSecurityContext:       true,
					EnableGRPCProbeInSuggestion: &diableGRPCProbeInSuggestion,
					TrialResources: []string{
						"Job.v1.batch",
						"TFJob.v1.kubeflow.org",
						"PyTorchJob.v1.kubeflow.org",
						"MPIJob.v1.kubeflow.org",
						"XGBoostJob.v1.kubeflow.org",
						"MXJob.v1.kubeflow.org",
					},
					WebhookPort:          &customizedWebhookPort,
					EnableLeaderElection: true,
					LeaderElectionID:     "xyz0123",
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := GetInitConfigData(scm, tc.katibConfigFile)
			if diff := cmp.Diff(tc.wantError, err, cmpopts.EquateErrors()); len(diff) != 0 {
				t.Errorf("Unexpected error from GetInitConfigData() (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantInitConfigData, got); len(diff) != 0 {
				t.Errorf("Unexpected InitConfig from GetInitConfigData() (-want,+got):\n%s", diff)
			}
		})
	}
}

func newFakeKubeClient(scm *runtime.Scheme, katibConfigMap *corev1.ConfigMap) client.Client {
	fakeClientBuilder := fake.NewClientBuilder().WithScheme(scm)
	if katibConfigMap != nil {
		fakeClientBuilder.WithObjects(katibConfigMap)
	}
	return fakeClientBuilder.Build()
}

func newFakeKatibConfigMap(config *configv1beta1.KatibConfig) *corev1.ConfigMap {
	if config == nil {
		return nil
	}

	data := map[string]string{}
	if config != nil {
		bKatibConfig, err := yaml.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		data[consts.LabelKatibConfigTag] = string(bKatibConfig)
	}
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      consts.KatibConfigMapName,
			Namespace: consts.DefaultKatibNamespace,
		},
		Data: data,
	}
}

func newFakeSuggestionConfig(algorithmName string) *configv1beta1.SuggestionConfig {
	defaultVolumeStorage, _ := resource.ParseQuantity(configv1beta1.DefaultSuggestionVolumeStorage)

	return &configv1beta1.SuggestionConfig{
		AlgorithmName: algorithmName,
		Container: corev1.Container{
			Image:           "suggestion-image",
			ImagePullPolicy: configv1beta1.DefaultImagePullPolicy,
			Resources:       *setFakeResourceRequirements(),
		},
		VolumeMountPath: configv1beta1.DefaultContainerSuggestionVolumeMountPath,
		PersistentVolumeClaimSpec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				configv1beta1.DefaultSuggestionVolumeAccessMode,
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

func newFakeEarlyStoppingConfig(algorithmName string) *configv1beta1.EarlyStoppingConfig {
	return &configv1beta1.EarlyStoppingConfig{
		AlgorithmName:   algorithmName,
		Image:           "early-stopping-image",
		ImagePullPolicy: configv1beta1.DefaultImagePullPolicy,
		Resource:        *setFakeResourceRequirements(),
	}
}

func newFakeMetricsCollectorConfig(collectorKind commonv1beta1.CollectorKind) *configv1beta1.MetricsCollectorConfig {
	return &configv1beta1.MetricsCollectorConfig{
		CollectorKind:   string(collectorKind),
		Image:           "metrics-collector-image",
		ImagePullPolicy: configv1beta1.DefaultImagePullPolicy,
		Resource:        *setFakeResourceRequirements(),
	}
}

func setFakeResourceRequirements() *corev1.ResourceRequirements {
	defaultCPURequest, _ := resource.ParseQuantity(configv1beta1.DefaultCPURequest)
	defaultMemoryRequest, _ := resource.ParseQuantity(configv1beta1.DefaultMemRequest)
	defaultEphemeralStorageRequest, _ := resource.ParseQuantity(configv1beta1.DefaultDiskRequest)

	defaultCPULimit, _ := resource.ParseQuantity(configv1beta1.DefaultCPULimit)
	defaultMemoryLimit, _ := resource.ParseQuantity(configv1beta1.DefaultMemLimit)
	defaultEphemeralStorageLimit, _ := resource.ParseQuantity(configv1beta1.DefaultDiskLimit)

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
