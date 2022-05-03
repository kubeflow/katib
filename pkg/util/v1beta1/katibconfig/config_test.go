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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

type katibConfig struct {
	suggestion       map[string]*SuggestionConfig
	earlyStopping    map[string]*EarlyStoppingConfig
	metricsCollector map[commonv1beta1.CollectorKind]*MetricsCollectorConfig
}

func TestGetSuggestionConfigData(t *testing.T) {
	const testAlgorithmName = "test-suggestion"

	tests := []struct {
		testDescription    string
		katibConfig        *katibConfig
		expected           *SuggestionConfig
		inputAlgorithmName string
		err                bool
	}{
		{
			testDescription: "All parameters correctly are specified",
			katibConfig: func() *katibConfig {
				kc := &katibConfig{suggestion: map[string]*SuggestionConfig{testAlgorithmName: newFakeSuggestionConfig()}}
				kc.suggestion[testAlgorithmName].ImagePullPolicy = corev1.PullAlways
				kc.suggestion[testAlgorithmName].Resource = *newFakeCustomResourceRequirements()
				return kc
			}(),
			expected: func() *SuggestionConfig {
				c := newFakeSuggestionConfig()
				c.ImagePullPolicy = corev1.PullAlways
				c.Resource = *newFakeCustomResourceRequirements()
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
			testDescription: fmt.Sprintf("There is not %s field in katib-config configMap", consts.LabelSuggestionTag),
			katibConfig:     &katibConfig{},
			err:             true,
		},
		{
			testDescription:    "There is not the AlgorithmName",
			katibConfig:        &katibConfig{suggestion: map[string]*SuggestionConfig{testAlgorithmName: newFakeSuggestionConfig()}},
			inputAlgorithmName: "invalid-algorithm-name",
			err:                true,
		},
		{
			testDescription: "Image filed is empty in katib-config configMap",
			katibConfig: func() *katibConfig {
				kc := &katibConfig{suggestion: map[string]*SuggestionConfig{testAlgorithmName: newFakeSuggestionConfig()}}
				kc.suggestion[testAlgorithmName].Image = ""
				return kc
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                true,
		},
		{
			testDescription: fmt.Sprintf("GetSuggestionConfigData sets %s to imagePullPolicy", consts.DefaultImagePullPolicy),
			katibConfig: func() *katibConfig {
				kc := &katibConfig{suggestion: map[string]*SuggestionConfig{testAlgorithmName: newFakeSuggestionConfig()}}
				kc.suggestion[testAlgorithmName].ImagePullPolicy = ""
				return kc
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: "GetSuggestionConfigData sets resource.requests and resource.limits for the suggestion service",
			katibConfig: func() *katibConfig {
				kc := &katibConfig{suggestion: map[string]*SuggestionConfig{testAlgorithmName: newFakeSuggestionConfig()}}
				kc.suggestion[testAlgorithmName].Resource = corev1.ResourceRequirements{}
				return kc
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: fmt.Sprintf("GetSuggestionConfigData sets %s to volumeMountPath", consts.DefaultContainerSuggestionVolumeMountPath),
			katibConfig: func() *katibConfig {
				kc := &katibConfig{suggestion: map[string]*SuggestionConfig{testAlgorithmName: newFakeSuggestionConfig()}}
				kc.suggestion[testAlgorithmName].VolumeMountPath = ""
				return kc
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: "GetSuggestionConfigData sets accessMode and resource.requests for PVC",
			katibConfig: func() *katibConfig {
				kc := &katibConfig{suggestion: map[string]*SuggestionConfig{testAlgorithmName: newFakeSuggestionConfig()}}
				kc.suggestion[testAlgorithmName].PersistentVolumeClaimSpec = corev1.PersistentVolumeClaimSpec{}
				return kc
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: fmt.Sprintf("GetSuggestionConfigData does not set %s to persistentVolumeReclaimPolicy", corev1.PersistentVolumeReclaimDelete),
			katibConfig: func() *katibConfig {
				kc := &katibConfig{suggestion: map[string]*SuggestionConfig{testAlgorithmName: newFakeSuggestionConfig()}}
				kc.suggestion[testAlgorithmName].PersistentVolumeSpec = corev1.PersistentVolumeSpec{}
				return kc
			}(),
			expected: func() *SuggestionConfig {
				c := newFakeSuggestionConfig()
				c.PersistentVolumeSpec = corev1.PersistentVolumeSpec{}
				return c
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			fakeKubeClient := newFakeKubeClient(newFakeKatibConfigMap(tt.katibConfig))
			actual, err := GetSuggestionConfigData(tt.inputAlgorithmName, fakeKubeClient)
			if (err != nil) != tt.err {
				t.Errorf("want error: %v, actual: %v", tt.err, err)
			} else if tt.expected != nil {
				if !reflect.DeepEqual(actual, *tt.expected) {
					t.Errorf("Generated SuggestionConfig is invalid.\n\nactual:\n%v\n\nexpected:\n%v\n\n", actual, *tt.expected)
				}
			}
		})
	}

}

func TestGetEarlyStoppingConfigData(t *testing.T) {
	const testAlgorithmName = "test-early-stopping"

	tests := []struct {
		testDescription    string
		katibConfig        *katibConfig
		expected           *EarlyStoppingConfig
		inputAlgorithmName string
		err                bool
	}{
		{
			testDescription: "All parameters correctly are specified",
			katibConfig: func() *katibConfig {
				kc := &katibConfig{earlyStopping: map[string]*EarlyStoppingConfig{testAlgorithmName: newFakeEarlyStoppingConfig()}}
				kc.earlyStopping[testAlgorithmName].ImagePullPolicy = corev1.PullIfNotPresent
				return kc
			}(),
			expected: func() *EarlyStoppingConfig {
				c := newFakeEarlyStoppingConfig()
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
			testDescription: fmt.Sprintf("There is not %s field in katib-config configMap", consts.LabelEarlyStoppingTag),
			katibConfig:     &katibConfig{},
			err:             true,
		},
		{
			testDescription:    "There is not the AlgorithmName",
			katibConfig:        &katibConfig{earlyStopping: map[string]*EarlyStoppingConfig{testAlgorithmName: newFakeEarlyStoppingConfig()}},
			inputAlgorithmName: "invalid-algorithm-name",
			err:                true,
		},
		{
			testDescription: "Image filed is empty in katib-config configMap",
			katibConfig: func() *katibConfig {
				kc := &katibConfig{earlyStopping: map[string]*EarlyStoppingConfig{testAlgorithmName: newFakeEarlyStoppingConfig()}}
				kc.earlyStopping[testAlgorithmName].Image = ""
				return kc
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                true,
		},
		{
			testDescription: fmt.Sprintf("GetEarlyStoppingConfigData sets %s to imagePullPolicy", consts.DefaultImagePullPolicy),
			katibConfig: func() *katibConfig {
				kc := &katibConfig{earlyStopping: map[string]*EarlyStoppingConfig{testAlgorithmName: newFakeEarlyStoppingConfig()}}
				kc.earlyStopping[testAlgorithmName].ImagePullPolicy = ""
				return kc
			}(),
			expected:           newFakeEarlyStoppingConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			fakeKubeClient := newFakeKubeClient(newFakeKatibConfigMap(tt.katibConfig))
			actual, err := GetEarlyStoppingConfigData(tt.inputAlgorithmName, fakeKubeClient)
			if (err != nil) != tt.err {
				t.Errorf("want error: %v, actual: %v", tt.err, err)
			} else if tt.expected != nil {
				if !reflect.DeepEqual(actual, *tt.expected) {
					t.Errorf("Generated EarlyStoppingConfig is invalid.\n\nactual:\n%v\n\nexpected:\n%v\n\n", actual, *tt.expected)
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

	nukeResource, _ := resource.ParseQuantity("-1")
	nukeResourceRequirements := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:              nukeResource,
		corev1.ResourceMemory:           nukeResource,
		corev1.ResourceEphemeralStorage: nukeResource,
	}

	tests := []struct {
		testDescription    string
		katibConfig        *katibConfig
		expected           *MetricsCollectorConfig
		inputCollectorKind commonv1beta1.CollectorKind
		err                bool
	}{
		{
			testDescription: "All parameters correctly are specified",
			katibConfig: func() *katibConfig {
				kc := &katibConfig{metricsCollector: map[commonv1beta1.CollectorKind]*MetricsCollectorConfig{testCollectorKind: newFakeMetricsCollectorConfig()}}
				kc.metricsCollector[testCollectorKind].ImagePullPolicy = corev1.PullNever
				return kc
			}(),
			expected: func() *MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig()
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
			testDescription: fmt.Sprintf("There is not %s field in katib-config configMap", consts.LabelMetricsCollectorSidecar),
			katibConfig:     &katibConfig{},
			err:             true,
		},
		{
			testDescription:    "There is not the cKind",
			katibConfig:        &katibConfig{metricsCollector: map[commonv1beta1.CollectorKind]*MetricsCollectorConfig{testCollectorKind: newFakeMetricsCollectorConfig()}},
			inputCollectorKind: invalidCollectorKind,
			err:                true,
		},
		{
			testDescription: "Image filed is empty in katib-config configMap",
			katibConfig: func() *katibConfig {
				kc := &katibConfig{metricsCollector: map[commonv1beta1.CollectorKind]*MetricsCollectorConfig{testCollectorKind: newFakeMetricsCollectorConfig()}}
				kc.metricsCollector[testCollectorKind].Image = ""
				return kc
			}(),
			inputCollectorKind: testCollectorKind,
			err:                true,
		},
		{
			testDescription: fmt.Sprintf("GetMetricsConfigData sets %s to imagePullPolicy", consts.DefaultImagePullPolicy),
			katibConfig: func() *katibConfig {
				kc := &katibConfig{metricsCollector: map[commonv1beta1.CollectorKind]*MetricsCollectorConfig{testCollectorKind: newFakeMetricsCollectorConfig()}}
				kc.metricsCollector[testCollectorKind].ImagePullPolicy = ""
				return kc
			}(),
			expected:           newFakeMetricsCollectorConfig(),
			inputCollectorKind: testCollectorKind,
			err:                false,
		},
		{
			testDescription: "GetMetricsConfigData nukes resource.requests and resource.limits for the metrics collector",
			katibConfig: func() *katibConfig {
				kc := &katibConfig{metricsCollector: map[commonv1beta1.CollectorKind]*MetricsCollectorConfig{testCollectorKind: newFakeMetricsCollectorConfig()}}
				kc.metricsCollector[testCollectorKind].Resource = corev1.ResourceRequirements{
					Requests: nukeResourceRequirements,
					Limits:   nukeResourceRequirements,
				}
				return kc
			}(),
			expected: func() *MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig()
				c.Resource = corev1.ResourceRequirements{
					Requests: map[corev1.ResourceName]resource.Quantity{},
					Limits:   map[corev1.ResourceName]resource.Quantity{},
				}
				return c
			}(),
			inputCollectorKind: testCollectorKind,
			err:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			fakeKubeClient := newFakeKubeClient(newFakeKatibConfigMap(tt.katibConfig))
			actual, err := GetMetricsCollectorConfigData(tt.inputCollectorKind, fakeKubeClient)
			if (err != nil) != tt.err {
				t.Errorf("want error: %v, actual: %v", tt.err, err)
			} else if tt.expected != nil {
				if !reflect.DeepEqual(actual, *tt.expected) {
					t.Errorf("Generated MetricsCollectorConfig is invalid.\n\nactual:\n%v\n\nexpected:\n%v\n\n", actual, *tt.expected)
				}
			}
		})
	}
}

func newFakeKubeClient(katibConfigMap *corev1.ConfigMap) client.Client {
	fakeClientBuilder := fake.NewClientBuilder().WithScheme(scheme.Scheme)
	if katibConfigMap != nil {
		fakeClientBuilder.WithObjects(katibConfigMap)
	}
	return fakeClientBuilder.Build()
}

func newFakeKatibConfigMap(config *katibConfig) *corev1.ConfigMap {
	if config == nil {
		return nil
	}

	data := map[string]string{}
	if config.suggestion != nil {
		bSuggestionConfig, _ := json.Marshal(config.suggestion)
		data[consts.LabelSuggestionTag] = string(bSuggestionConfig)
	}
	if config.earlyStopping != nil {
		bEarlyStoppingConfig, _ := json.Marshal(config.earlyStopping)
		data[consts.LabelEarlyStoppingTag] = string(bEarlyStoppingConfig)
	}
	if config.metricsCollector != nil {
		bMetricsCollector, _ := json.Marshal(config.metricsCollector)
		data[consts.LabelMetricsCollectorSidecar] = string(bMetricsCollector)
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

func newFakeSuggestionConfig() *SuggestionConfig {
	defaultVolumeStorage, _ := resource.ParseQuantity(consts.DefaultSuggestionVolumeStorage)

	return &SuggestionConfig{
		Image:           "suggestion-image",
		ImagePullPolicy: consts.DefaultImagePullPolicy,
		Resource:        *setFakeResourceRequirements(),
		VolumeMountPath: consts.DefaultContainerSuggestionVolumeMountPath,
		PersistentVolumeClaimSpec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				consts.DefaultSuggestionVolumeAccessMode,
			},
			Resources: corev1.ResourceRequirements{
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

func newFakeEarlyStoppingConfig() *EarlyStoppingConfig {
	return &EarlyStoppingConfig{
		Image:           "early-stopping-image",
		ImagePullPolicy: consts.DefaultImagePullPolicy,
	}
}

func newFakeMetricsCollectorConfig() *MetricsCollectorConfig {
	return &MetricsCollectorConfig{
		Image:           "metrics-collector-image",
		ImagePullPolicy: consts.DefaultImagePullPolicy,
		Resource:        *setFakeResourceRequirements(),
	}
}

func setFakeResourceRequirements() *corev1.ResourceRequirements {
	defaultCPURequest, _ := resource.ParseQuantity(consts.DefaultCPURequest)
	defaultMemoryRequest, _ := resource.ParseQuantity(consts.DefaultMemRequest)
	defaultEphemeralStorageRequest, _ := resource.ParseQuantity(consts.DefaultDiskRequest)

	defaultCPULimit, _ := resource.ParseQuantity(consts.DefaultCPULimit)
	defaultMemoryLimit, _ := resource.ParseQuantity(consts.DefaultMemLimit)
	defaultEphemeralStorageLimit, _ := resource.ParseQuantity(consts.DefaultDiskLimit)

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
