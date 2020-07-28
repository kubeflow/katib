package composer

import (
	"context"
	"encoding/json"
	"fmt"
	stdlog "log"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"sigs.k8s.io/controller-runtime/pkg/manager"

	apis "github.com/kubeflow/katib/pkg/apis/controller"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

var (
	cfg     *rest.Config
	timeout = time.Second * 40

	suggestionName      = "test-suggestion"
	suggestionAlgorithm = "random"
	suggestionLabels    = map[string]string{
		"custom-label": "test",
	}
	suggestionAnnotations = map[string]string{
		"custom-annotation": "test",
	}

	deploymentLabels = map[string]string{
		"custom-label": "test",
		"deployment":   suggestionName + "-" + suggestionAlgorithm,
		"experiment":   suggestionName,
		"suggestion":   suggestionName,
	}

	podAnnotations = map[string]string{
		"custom-annotation":       "test",
		"sidecar.istio.io/inject": "false",
	}

	namespace       = "kubeflow"
	configMap       = "katib-config"
	serviceAccount  = "test-serviceaccount"
	image           = "test-image"
	imagePullPolicy = corev1.PullAlways

	cpu    = "1m"
	memory = "2Mi"
	disk   = "3Gi"
)

func TestMain(m *testing.M) {
	// Start test k8s server
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "manifests", "v1beta1", "katib-controller"),
		},
	}
	apis.AddToScheme(scheme.Scheme)
	var err error
	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}

	code := m.Run()
	t.Stop()
	os.Exit(code)
}

func StartTestManager(mgr manager.Manager, g *gomega.GomegaWithT) (chan struct{}, *sync.WaitGroup) {
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(stop)).NotTo(gomega.HaveOccurred())
	}()
	return stop, wg
}

func TestDesiredDeployment(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	apis.AddToScheme(scheme.Scheme)
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)
	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	c := mgr.GetClient()

	composer := New(mgr)

	tcs := []struct {
		suggestion         *suggestionsv1beta1.Suggestion
		configMap          *corev1.ConfigMap
		expectedDeployment *appsv1.Deployment
		err                bool
		testDescription    string
	}{
		{
			suggestion:         newFakeSuggestion(),
			configMap:          newFakeKatibConfig(),
			expectedDeployment: newFakeDeployment(),
			err:                false,
			testDescription:    "Desired Deployment valid run",
		},
	}

	viper.Set(consts.ConfigEnableGRPCProbeInSuggestion, true)

	for _, tc := range tcs {
		// Create configMap with Katib config
		g.Expect(c.Create(context.TODO(), tc.configMap)).NotTo(gomega.HaveOccurred())

		// Wait that Config Map is created
		g.Eventually(func() error {
			return c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: configMap}, &corev1.ConfigMap{})
		}, timeout).ShouldNot(gomega.HaveOccurred())

		// Get deployment
		actualDeployment, err := composer.DesiredDeployment(tc.suggestion)

		if !tc.err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
		} else if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		} else if metaEqual(tc.expectedDeployment.ObjectMeta, actualDeployment.ObjectMeta) {
			t.Errorf("Case: %v failed. \nExpected deploy metadata %v\n Got %v", tc.testDescription, tc.expectedDeployment.ObjectMeta, actualDeployment.ObjectMeta)
		} else if !equality.Semantic.DeepEqual(tc.expectedDeployment.Spec, actualDeployment.Spec) {
			t.Errorf("Case: %v failed. \nExpected deploy spec %v\n Got %v", tc.testDescription, tc.expectedDeployment.Spec, actualDeployment.Spec)
		}

		// Delete configMap with Katib config
		g.Expect(c.Delete(context.TODO(), tc.configMap)).NotTo(gomega.HaveOccurred())

		// Wait that Config Map is deleted
		g.Eventually(func() bool {
			return errors.IsNotFound(
				c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: configMap}, &corev1.ConfigMap{}))
		}, timeout).Should(gomega.BeTrue())

	}
}

func metaEqual(expected, actual metav1.ObjectMeta) bool {
	return expected.Name == actual.Name &&
		expected.Namespace == actual.Namespace &&
		reflect.DeepEqual(expected.Labels, actual.Labels) &&
		reflect.DeepEqual(expected.Annotations, actual.Annotations) &&
		len(actual.OwnerReferences) > 0 &&
		expected.OwnerReferences[0].APIVersion == expected.OwnerReferences[0].APIVersion &&
		expected.OwnerReferences[0].Kind == expected.OwnerReferences[0].Kind &&
		expected.OwnerReferences[0].Name == expected.OwnerReferences[0].Name &&
		expected.OwnerReferences[0].Controller == expected.OwnerReferences[0].Controller &&
		expected.OwnerReferences[0].BlockOwnerDeletion == expected.OwnerReferences[0].BlockOwnerDeletion
}

func newFakeKatibConfig() *corev1.ConfigMap {
	cpuQ, _ := resource.ParseQuantity(cpu)
	memoryQ, _ := resource.ParseQuantity(memory)
	diskQ, _ := resource.ParseQuantity(disk)

	type suggestionConfigJSON struct {
		Image              string                      `json:"image"`
		ImagePullPolicy    corev1.PullPolicy           `json:"imagePullPolicy"`
		Resource           corev1.ResourceRequirements `json:"resources"`
		ServiceAccountName string                      `json:"serviceAccountName"`
	}
	jsonConfig := map[string]suggestionConfigJSON{
		"random": suggestionConfigJSON{
			Image:           image,
			ImagePullPolicy: imagePullPolicy,
			Resource: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:              cpuQ,
					corev1.ResourceMemory:           memoryQ,
					corev1.ResourceEphemeralStorage: diskQ,
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:              cpuQ,
					corev1.ResourceMemory:           memoryQ,
					corev1.ResourceEphemeralStorage: diskQ,
				},
			},
			ServiceAccountName: serviceAccount,
		},
	}

	b, _ := json.Marshal(jsonConfig)

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMap,
			Namespace: namespace,
		},
		Data: map[string]string{
			"suggestion": string(b),
		},
	}
}

func newFakeSuggestion() *suggestionsv1beta1.Suggestion {
	return &suggestionsv1beta1.Suggestion{
		ObjectMeta: metav1.ObjectMeta{
			Name:        suggestionName,
			Namespace:   namespace,
			Labels:      suggestionLabels,
			Annotations: suggestionAnnotations,
		},
		Spec: suggestionsv1beta1.SuggestionSpec{
			Requests:      1,
			AlgorithmName: suggestionAlgorithm,
			ResumePolicy:  experimentsv1beta1.FromVolume,
		},
	}
}

func newFakeDeployment() *appsv1.Deployment {
	var flag bool = true
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        suggestionName + "-" + suggestionAlgorithm,
			Namespace:   namespace,
			Labels:      suggestionLabels,
			Annotations: suggestionLabels,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "kubeflow.org/v1beta1",
					Kind:               "Suggestion",
					Name:               suggestionName,
					Controller:         &flag,
					BlockOwnerDeletion: &flag,
				},
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: deploymentLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      deploymentLabels,
					Annotations: podAnnotations,
				},
				Spec: corev1.PodSpec{
					Containers:         newFakeContainers(),
					ServiceAccountName: serviceAccount,
					Volumes: []corev1.Volume{
						{
							Name: consts.ContainerSuggestionVolumeName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: suggestionName + "-" + suggestionAlgorithm,
								},
							},
						},
					},
				},
			},
		},
	}
}

func newFakeContainers() []corev1.Container {

	cpuQ, _ := resource.ParseQuantity(cpu)
	memoryQ, _ := resource.ParseQuantity(memory)
	diskQ, _ := resource.ParseQuantity(disk)

	return []corev1.Container{
		{
			Name:            consts.ContainerSuggestion,
			Image:           image,
			ImagePullPolicy: corev1.PullAlways,
			Ports: []corev1.ContainerPort{
				{
					Name:          consts.DefaultSuggestionPortName,
					ContainerPort: consts.DefaultSuggestionPort,
				},
			},
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:              cpuQ,
					corev1.ResourceMemory:           memoryQ,
					corev1.ResourceEphemeralStorage: diskQ,
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:              cpuQ,
					corev1.ResourceMemory:           memoryQ,
					corev1.ResourceEphemeralStorage: diskQ,
				},
			},
			ReadinessProbe: &corev1.Probe{
				Handler: corev1.Handler{
					Exec: &corev1.ExecAction{
						Command: []string{
							defaultGRPCHealthCheckProbe,
							fmt.Sprintf("-addr=:%d", consts.DefaultSuggestionPort),
							fmt.Sprintf("-service=%s", consts.DefaultGRPCService),
						},
					},
				},
				InitialDelaySeconds: defaultInitialDelaySeconds,
				PeriodSeconds:       defaultPeriodForReady,
			},
			LivenessProbe: &corev1.Probe{
				Handler: corev1.Handler{
					Exec: &corev1.ExecAction{
						Command: []string{
							defaultGRPCHealthCheckProbe,
							fmt.Sprintf("-addr=:%d", consts.DefaultSuggestionPort),
							fmt.Sprintf("-service=%s", consts.DefaultGRPCService),
						},
					},
				},
				InitialDelaySeconds: defaultInitialDelaySeconds,
				PeriodSeconds:       defaultPeriodForLive,
				FailureThreshold:    defaultFailureThreshold,
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      consts.ContainerSuggestionVolumeName,
					MountPath: consts.DefaultContainerSuggestionVolumeMountPath,
				},
			},
		},
	}
}
