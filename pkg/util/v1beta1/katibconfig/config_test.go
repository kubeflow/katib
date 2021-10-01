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
	metricsCollector map[string]*MetricsCollectorConfig
}

func TestGetSuggestionConfigData(t *testing.T) {
	const testAlgorithmName = "test-suggestion"

	tests := []struct {
		testDescription          string
		katibConfigMapSuggestion *SuggestionConfig
		expected                 *SuggestionConfig
		inputAlgorithmName       string
		err                      bool
		nonExistentSuggestion    bool
	}{
		{
			testDescription: "All parameters are specified",
			katibConfigMapSuggestion: func() *SuggestionConfig {
				s := newFakeSuggestionConfig()
				s.ImagePullPolicy = corev1.PullAlways
				return s
			}(),
			expected: func() *SuggestionConfig {
				s := newFakeSuggestionConfig()
				s.ImagePullPolicy = corev1.PullAlways
				return s
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription:       "There is not the suggestion field in katib-config configMap",
			err:                   true,
			nonExistentSuggestion: true,
		},
		{
			testDescription:          "There is not the AlgorithmName",
			katibConfigMapSuggestion: newFakeSuggestionConfig(),
			inputAlgorithmName:       "invalid-algorithm-name",
			err:                      true,
		},
		{
			testDescription: "Image filed is empty in katib-config configMap",
			katibConfigMapSuggestion: func() *SuggestionConfig {
				s := newFakeSuggestionConfig()
				s.Image = ""
				return s
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                true,
		},
		{
			testDescription: fmt.Sprintf("GetSuggestionConfigData sets %s to imagePullPolicy", consts.DefaultImagePullPolicy),
			katibConfigMapSuggestion: func() *SuggestionConfig {
				s := newFakeSuggestionConfig()
				s.ImagePullPolicy = ""
				return s
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: "GetSuggestionConfigData sets resource.requests and resource.limits for the suggestion service",
			katibConfigMapSuggestion: func() *SuggestionConfig {
				s := newFakeSuggestionConfig()
				s.Resource = corev1.ResourceRequirements{}
				return s
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: fmt.Sprintf("GetSuggestionConfigData sets %s to volumeMountPath", consts.DefaultContainerSuggestionVolumeMountPath),
			katibConfigMapSuggestion: func() *SuggestionConfig {
				s := newFakeSuggestionConfig()
				s.VolumeMountPath = ""
				return s
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: "GetSuggestionConfigData sets accessMode and resource.requests for PVC",
			katibConfigMapSuggestion: func() *SuggestionConfig {
				s := newFakeSuggestionConfig()
				s.PersistentVolumeClaimSpec = corev1.PersistentVolumeClaimSpec{}
				return s
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: fmt.Sprintf("GetSuggestionConfigData does not set %s to persistentVolumeReclaimPolicy", corev1.PersistentVolumeReclaimDelete),
			katibConfigMapSuggestion: func() *SuggestionConfig {
				s := newFakeSuggestionConfig()
				s.PersistentVolumeSpec = corev1.PersistentVolumeSpec{}
				return s
			}(),
			expected: func() *SuggestionConfig {
				s := newFakeSuggestionConfig()
				s.PersistentVolumeSpec = corev1.PersistentVolumeSpec{}
				return s
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			// prepare test
			config := &katibConfig{}
			if !tt.nonExistentSuggestion {
				config.suggestion = map[string]*SuggestionConfig{
					testAlgorithmName: tt.katibConfigMapSuggestion,
				}
			}
			c := newFakeKubeClient(newFakeKatibConfigMap(config))

			// start test
			actual, err := GetSuggestionConfigData(tt.inputAlgorithmName, c)
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
		testDescription             string
		katibConfigMapEarlyStopping *EarlyStoppingConfig
		expected                    *EarlyStoppingConfig
		inputAlgorithmName          string
		err                         bool
		nonExistentEarlyStopping    bool
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			// prepare test
			config := &katibConfig{}
			if !tt.nonExistentEarlyStopping {
				config.earlyStopping = map[string]*EarlyStoppingConfig{
					testAlgorithmName: tt.katibConfigMapEarlyStopping,
				}
			}
			c := newFakeKubeClient(newFakeKatibConfigMap(config))

			// start test
			actual, err := GetEarlyStoppingConfigData(tt.inputAlgorithmName, c)
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
	const testMetricsCollector = "test-metrics-collector"

	tests := []struct {
		testDescription                string
		katibConfigMapMetricsCollector *MetricsCollectorConfig
		expected           *MetricsCollectorConfig
		inputCollectorKind commonv1beta1.CollectorKind
		err                bool
		nonExistentMetricsCollector    bool
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			// prepare test
			config := &katibConfig{}
			if !tt.nonExistentMetricsCollector {
				config.metricsCollector = map[string]*MetricsCollectorConfig{
					testMetricsCollector: tt.katibConfigMapMetricsCollector,
				}
			}
			c := newFakeKubeClient(newFakeKatibConfigMap(config))

			// start test
			actual, err := GetMetricsCollectorConfigData(tt.inputCollectorKind, c)
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

	suggestionConfig := ""
	if config.suggestion != nil {
		bSuggestionConfig, _ := json.Marshal(config.suggestion)
		suggestionConfig = string(bSuggestionConfig)
	}

	earlyStoppingConfig := ""
	if config.earlyStopping != nil {
		bEarlyStoppingConfig, _ := json.Marshal(config.earlyStopping)
		earlyStoppingConfig = string(bEarlyStoppingConfig)
	}

	metricsCollector := ""
	if config.metricsCollector != nil {
		bMetricsCollector, _ := json.Marshal(config.metricsCollector)
		metricsCollector = string(bMetricsCollector)
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
		Data: map[string]string{
			consts.LabelSuggestionTag:           suggestionConfig,
			consts.LabelEarlyStoppingTag:        earlyStoppingConfig,
			consts.LabelMetricsCollectorSidecar: metricsCollector,
		},
	}
}

func newFakeSuggestionConfig() *SuggestionConfig {
	defaultCPURequest, _ := resource.ParseQuantity(consts.DefaultCPURequest)
	defaultMemoryRequest, _ := resource.ParseQuantity(consts.DefaultMemRequest)
	defaultEphemeralStorageRequest, _ := resource.ParseQuantity(consts.DefaultDiskRequest)

	defaultCPULimit, _ := resource.ParseQuantity(consts.DefaultCPULimit)
	defaultMemoryLimit, _ := resource.ParseQuantity(consts.DefaultMemLimit)
	defaultEphemeralStorageLimit, _ := resource.ParseQuantity(consts.DefaultDiskLimit)

	defaultVolumeStorage, _ := resource.ParseQuantity(consts.DefaultSuggestionVolumeStorage)

	return &SuggestionConfig{
		Image:           "suggestion-image",
		ImagePullPolicy: consts.DefaultImagePullPolicy,
		Resource: corev1.ResourceRequirements{
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
		},
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
		Image: "early-stopping-image",
	}
}

func newFakeMetricsCollectorConfig() *MetricsCollectorConfig {
	return &MetricsCollectorConfig{
		Image: "metrics-collector-image",
	}
}
