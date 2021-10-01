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
		testDescription          string
		katibConfigMapSuggestion *SuggestionConfig
		expected                 *SuggestionConfig
		inputAlgorithmName       string
		err                      bool
		nonExistentSuggestion    bool
	}{
		{
			testDescription: "All parameters correctly are specified",
			katibConfigMapSuggestion: func() *SuggestionConfig {
				c := newFakeSuggestionConfig()
				c.ImagePullPolicy = corev1.PullAlways
				return c
			}(),
			expected: func() *SuggestionConfig {
				c := newFakeSuggestionConfig()
				c.ImagePullPolicy = corev1.PullAlways
				return c
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription:       fmt.Sprintf("There is not %s field in katib-config configMap", consts.LabelSuggestionTag),
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
				c := newFakeSuggestionConfig()
				c.Image = ""
				return c
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                true,
		},
		{
			testDescription: fmt.Sprintf("GetSuggestionConfigData sets %s to imagePullPolicy", consts.DefaultImagePullPolicy),
			katibConfigMapSuggestion: func() *SuggestionConfig {
				c := newFakeSuggestionConfig()
				c.ImagePullPolicy = ""
				return c
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: "GetSuggestionConfigData sets resource.requests and resource.limits for the suggestion service",
			katibConfigMapSuggestion: func() *SuggestionConfig {
				c := newFakeSuggestionConfig()
				c.Resource = corev1.ResourceRequirements{}
				return c
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: fmt.Sprintf("GetSuggestionConfigData sets %s to volumeMountPath", consts.DefaultContainerSuggestionVolumeMountPath),
			katibConfigMapSuggestion: func() *SuggestionConfig {
				c := newFakeSuggestionConfig()
				c.VolumeMountPath = ""
				return c
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: "GetSuggestionConfigData sets accessMode and resource.requests for PVC",
			katibConfigMapSuggestion: func() *SuggestionConfig {
				c := newFakeSuggestionConfig()
				c.PersistentVolumeClaimSpec = corev1.PersistentVolumeClaimSpec{}
				return c
			}(),
			expected:           newFakeSuggestionConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
		{
			testDescription: fmt.Sprintf("GetSuggestionConfigData does not set %s to persistentVolumeReclaimPolicy", corev1.PersistentVolumeReclaimDelete),
			katibConfigMapSuggestion: func() *SuggestionConfig {
				c := newFakeSuggestionConfig()
				c.PersistentVolumeSpec = corev1.PersistentVolumeSpec{}
				return c
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
		{
			testDescription: "All parameters correctly are specified",
			katibConfigMapEarlyStopping: func() *EarlyStoppingConfig {
				c := newFakeEarlyStoppingConfig()
				c.ImagePullPolicy = corev1.PullIfNotPresent
				return c
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
			testDescription:          fmt.Sprintf("There is not %s field in katib-config configMap", consts.LabelEarlyStoppingTag),
			err:                      true,
			nonExistentEarlyStopping: true,
		},
		{
			testDescription:             "There is not the AlgorithmName",
			katibConfigMapEarlyStopping: newFakeEarlyStoppingConfig(),
			inputAlgorithmName:          "invalid-algorithm-name",
			err:                         true,
		},
		{
			testDescription: "Image filed is empty in katib-config configMap",
			katibConfigMapEarlyStopping: func() *EarlyStoppingConfig {
				c := newFakeEarlyStoppingConfig()
				c.Image = ""
				return c
			}(),
			inputAlgorithmName: testAlgorithmName,
			err:                true,
		},
		{
			testDescription: fmt.Sprintf("GetEarlyStoppingConfigData sets %s to imagePullPolicy", consts.DefaultImagePullPolicy),
			katibConfigMapEarlyStopping: func() *EarlyStoppingConfig {
				c := newFakeEarlyStoppingConfig()
				c.ImagePullPolicy = ""
				return c
			}(),
			expected:           newFakeEarlyStoppingConfig(),
			inputAlgorithmName: testAlgorithmName,
			err:                false,
		},
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
	const (
		invalidCollectorKind commonv1beta1.CollectorKind = "invalid-collector-kind"
		testCollectorKind    commonv1beta1.CollectorKind = "testCollector"
	)

	nukeResource, _ := resource.ParseQuantity("-1")
	nukeResourceRequirements := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:              nukeResource,
		corev1.ResourceMemory:           nukeResource,
		corev1.ResourceEphemeralStorage: nukeResource,
	}

	tests := []struct {
		testDescription                string
		katibConfigMapMetricsCollector *MetricsCollectorConfig
		expected                       *MetricsCollectorConfig
		inputCollectorKind             commonv1beta1.CollectorKind
		err                            bool
		nonExistentMetricsCollector    bool
	}{
		{
			testDescription: "All parameters correctly are specified",
			katibConfigMapMetricsCollector: func() *MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig()
				c.ImagePullPolicy = corev1.PullNever
				return c
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
			testDescription:             fmt.Sprintf("There is not %s field in katib-config configMap", consts.LabelMetricsCollectorSidecar),
			err:                         true,
			nonExistentMetricsCollector: true,
		},
		{
			testDescription:                "There is not the cKind",
			katibConfigMapMetricsCollector: newFakeMetricsCollectorConfig(),
			inputCollectorKind:             invalidCollectorKind,
			err:                            true,
		},
		{
			testDescription: "Image filed is empty in katib-config configMap",
			katibConfigMapMetricsCollector: func() *MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig()
				c.Image = ""
				return c
			}(),
			inputCollectorKind: testCollectorKind,
			err:                true,
		},
		{
			testDescription: fmt.Sprintf("GetMetricsConfigData sets %s to imagePullPolicy", consts.DefaultImagePullPolicy),
			katibConfigMapMetricsCollector: func() *MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig()
				c.ImagePullPolicy = ""
				return c
			}(),
			expected:           newFakeMetricsCollectorConfig(),
			inputCollectorKind: testCollectorKind,
			err:                false,
		},
		{
			testDescription: "GetMetricsConfigData nukes resource.requests and resource.limits for the metrics collector",
			katibConfigMapMetricsCollector: func() *MetricsCollectorConfig {
				c := newFakeMetricsCollectorConfig()
				c.Resource = corev1.ResourceRequirements{
					Requests: nukeResourceRequirements,
					Limits:   nukeResourceRequirements,
				}
				return c
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
			// prepare test
			config := &katibConfig{}
			if !tt.nonExistentMetricsCollector {
				config.metricsCollector = map[commonv1beta1.CollectorKind]*MetricsCollectorConfig{
					testCollectorKind: tt.katibConfigMapMetricsCollector,
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
