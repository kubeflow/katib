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

package composer

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/yaml"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	apis "github.com/kubeflow/katib/pkg/apis/controller"
	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

var (
	cfg     *rest.Config
	timeout = time.Second * 40
	ctx     context.Context
	cancel  context.CancelFunc

	suggestionName         = "test-suggestion"
	suggestionAlgorithm    = "random"
	earlyStoppingAlgorithm = "median-stop"
	suggestionLabels       = map[string]string{
		"custom-label": "test",
	}
	suggestionAnnotations = map[string]string{
		"custom-annotation": "test",
	}

	deploymentLabels = map[string]string{
		"custom-label":             "test",
		consts.LabelDeploymentName: suggestionName + "-" + suggestionAlgorithm,
		consts.LabelExperimentName: suggestionName,
		consts.LabelSuggestionName: suggestionName,
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

	cpu    = "2m"
	memory = "3Mi"
	disk   = "4Gi"

	refFlag bool = true

	storageClassName = "test-storage-class"
	pvLabels         = map[string]string{"type": "local"}
)

func TestMain(m *testing.M) {
	// To avoid the `timeout waiting for process kube-apiserver to stop` error,
	// we must use the `context.WithCancel`.
	// Ref: https://github.com/kubernetes-sigs/controller-runtime/issues/1571#issuecomment-945535598
	ctx, cancel = context.WithCancel(context.TODO())

	// Start test k8s server
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "manifests", "v1beta1", "components", "crd"),
		},
	}
	var err error
	if err = apis.AddToScheme(scheme.Scheme); err != nil {
		stdlog.Fatal(err)
	}

	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}

	code := m.Run()
	cancel()
	if err = t.Stop(); err != nil {
		stdlog.Fatal(err)
	}
	os.Exit(code)
}

func TestDesiredDeployment(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mgr, err := manager.New(cfg, manager.Options{Metrics: metricsserver.Options{BindAddress: "0"}})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(configv1beta1.AddToScheme(mgr.GetScheme())).NotTo(gomega.HaveOccurred())

	// Start test manager.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(ctx)).NotTo(gomega.HaveOccurred())
	}()

	c := mgr.GetClient()
	composer := New(mgr)
	// Create kubeflow namespace.
	kubeflowNS := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	g.Expect(c.Create(ctx, kubeflowNS)).NotTo(gomega.HaveOccurred())

	tcs := []struct {
		suggestion         *suggestionsv1beta1.Suggestion
		configMap          *corev1.ConfigMap
		expectedDeployment *appsv1.Deployment
		err                bool
		testDescription    string
	}{
		{
			suggestion:      newFakeSuggestion(),
			configMap:       newFakeKatibConfig(newFakeSuggestionConfig(), newFakeEarlyStoppingConfig()),
			err:             true,
			testDescription: "Set controller reference error",
		},
		{
			suggestion:         newFakeSuggestion(),
			configMap:          newFakeKatibConfig(newFakeSuggestionConfig(), newFakeEarlyStoppingConfig()),
			expectedDeployment: newFakeDeployment(),
			err:                false,
			testDescription:    "Desired Deployment valid run",
		},
		{
			suggestion: newFakeSuggestion(),
			configMap: func() *corev1.ConfigMap {
				cm := newFakeKatibConfig(newFakeSuggestionConfig(), newFakeEarlyStoppingConfig())
				cm.Data[consts.LabelKatibConfigTag] = strings.ReplaceAll(cm.Data[consts.LabelKatibConfigTag], string(imagePullPolicy), "invalid")
				return cm
			}(),
			expectedDeployment: func() *appsv1.Deployment {
				deploy := newFakeDeployment()
				deploy.Spec.Template.Spec.Containers[0].ImagePullPolicy = configv1beta1.DefaultImagePullPolicy
				deploy.Spec.Template.Spec.Containers[1].ImagePullPolicy = configv1beta1.DefaultImagePullPolicy
				return deploy
			}(),
			err:             false,
			testDescription: "Image Pull Policy set to default",
		},
		{
			suggestion: newFakeSuggestion(),
			configMap: func() *corev1.ConfigMap {
				cm := newFakeKatibConfig(newFakeSuggestionConfig(), newFakeEarlyStoppingConfig())
				cm.Data[consts.LabelKatibConfigTag] = strings.ReplaceAll(cm.Data[consts.LabelKatibConfigTag], cpu, "invalid")
				return cm
			}(),
			err:             true,
			testDescription: "Get suggestion config error, invalid CPU limit",
		},
		{
			suggestion: newFakeSuggestion(),
			configMap: func() *corev1.ConfigMap {
				sc := newFakeSuggestionConfig()
				sc.VolumeMountPath = "/custom/container/path"
				cm := newFakeKatibConfig(sc, newFakeEarlyStoppingConfig())
				return cm
			}(),
			expectedDeployment: func() *appsv1.Deployment {
				deploy := newFakeDeployment()
				deploy.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath = "/custom/container/path"
				return deploy
			}(),
			err:             false,
			testDescription: "Suggestion container with custom volume mount path",
		},
		{
			suggestion: newFakeSuggestion(),
			configMap: func() *corev1.ConfigMap {
				sc := newFakeSuggestionConfig()
				sc.ServiceAccountName = ""
				cm := newFakeKatibConfig(sc, newFakeEarlyStoppingConfig())
				return cm
			}(),
			expectedDeployment: func() *appsv1.Deployment {
				deploy := newFakeDeployment()
				deploy.Spec.Template.Spec.ServiceAccountName = suggestionName + "-" + suggestionAlgorithm
				return deploy
			}(),
			err:             false,
			testDescription: "Desired Deployment valid run with default serviceAccount",
		},
		{
			suggestion: newFakeSuggestion(),
			configMap: func() *corev1.ConfigMap {
				esC := newFakeEarlyStoppingConfig()
				esC.Image = ""
				cm := newFakeKatibConfig(newFakeSuggestionConfig(), esC)
				return cm
			}(),
			err:             true,
			testDescription: "Get early stopping config error, image is missed",
		},
	}

	viper.Set(consts.ConfigEnableGRPCProbeInSuggestion, true)

	for idx, tc := range tcs {
		// Create configMap with Katib config
		g.Expect(c.Create(ctx, tc.configMap)).NotTo(gomega.HaveOccurred())

		// Wait that Config Map is created
		g.Eventually(func() error {
			return c.Get(ctx, types.NamespacedName{Namespace: namespace, Name: configMap}, &corev1.ConfigMap{})
		}, timeout).ShouldNot(gomega.HaveOccurred())

		// Get deployment
		var actualDeployment *appsv1.Deployment
		var err error
		// For the first Test we run DesiredDeployment with empty Scheme to fail Set Controller Reference
		if idx == 0 {
			c := General{
				scheme: &runtime.Scheme{},
				Client: mgr.GetClient(),
			}
			actualDeployment, err = c.DesiredDeployment(tc.suggestion)
		} else {
			actualDeployment, err = composer.DesiredDeployment(tc.suggestion)
		}

		if !tc.err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
		} else if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		} else if !tc.err && !metaEqual(tc.expectedDeployment.ObjectMeta, actualDeployment.ObjectMeta) {
			t.Errorf("Case: %v failed. \nExpected deploy metadata %v\n Got %v", tc.testDescription, tc.expectedDeployment.ObjectMeta, actualDeployment.ObjectMeta)
		} else if !tc.err && !equality.Semantic.DeepEqual(tc.expectedDeployment.Spec, actualDeployment.Spec) {
			t.Errorf("Case: %v failed. \nExpected deploy spec %v\n Got %v", tc.testDescription, tc.expectedDeployment.Spec, actualDeployment.Spec)
		}

		// Delete configMap with Katib config
		g.Expect(c.Delete(ctx, tc.configMap)).NotTo(gomega.HaveOccurred())

		// Wait that Config Map is deleted
		g.Eventually(func() bool {
			return errors.IsNotFound(
				c.Get(ctx, types.NamespacedName{Namespace: namespace, Name: configMap}, &corev1.ConfigMap{}))
		}, timeout).Should(gomega.BeTrue())

	}
}

func TestDesiredService(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mgr, err := manager.New(cfg, manager.Options{Metrics: metricsserver.Options{BindAddress: "0"}})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	composer := New(mgr)

	expectedService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      suggestionName + "-" + suggestionAlgorithm,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "kubeflow.org/v1beta1",
					Kind:               "Suggestion",
					Name:               suggestionName,
					Controller:         &refFlag,
					BlockOwnerDeletion: &refFlag,
				},
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: deploymentLabels,
			Ports: []corev1.ServicePort{
				{
					Name: consts.DefaultSuggestionPortName,
					Port: consts.DefaultSuggestionPort,
				},
				{
					Name: consts.DefaultEarlyStoppingPortName,
					Port: consts.DefaultEarlyStoppingPort,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	tcs := []struct {
		suggestion      *suggestionsv1beta1.Suggestion
		expectedService *corev1.Service
		err             bool
		testDescription string
	}{
		{
			suggestion:      newFakeSuggestion(),
			err:             true,
			testDescription: "Set controller reference error",
		},
		{
			suggestion:      newFakeSuggestion(),
			expectedService: expectedService,
			err:             false,
			testDescription: "Desired Service valid run",
		},
	}

	for idx, tc := range tcs {

		// Get service
		var actualService *corev1.Service
		var err error
		// For the first Test we run DesiredService with empty Scheme to fail Set Controller Reference
		if idx == 0 {
			c := General{
				scheme: &runtime.Scheme{},
				Client: mgr.GetClient(),
			}
			actualService, err = c.DesiredService(tc.suggestion)
		} else {
			actualService, err = composer.DesiredService(tc.suggestion)
		}

		if !tc.err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
		} else if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		} else if !tc.err && !metaEqual(tc.expectedService.ObjectMeta, actualService.ObjectMeta) {
			t.Errorf("Case: %v failed. \nExpected service metadata %v\n Got %v", tc.testDescription, tc.expectedService.ObjectMeta, actualService.ObjectMeta)
		} else if !tc.err && !equality.Semantic.DeepEqual(tc.expectedService.Spec, actualService.Spec) {
			t.Errorf("Case: %v failed. \nExpected service spec %v\n Got %v", tc.testDescription, tc.expectedService.Spec, actualService.Spec)
		}
	}
}

func TestDesiredVolume(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mgr, err := manager.New(cfg, manager.Options{Metrics: metricsserver.Options{BindAddress: "0"}})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(configv1beta1.AddToScheme(mgr.GetScheme())).NotTo(gomega.HaveOccurred())

	// Start test manager.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(ctx)).NotTo(gomega.HaveOccurred())
	}()

	c := mgr.GetClient()
	composer := New(mgr)

	tcs := []struct {
		suggestion      *suggestionsv1beta1.Suggestion
		configMap       *corev1.ConfigMap
		expectedPVC     *corev1.PersistentVolumeClaim
		expectedPV      *corev1.PersistentVolume
		err             bool
		testDescription string
	}{
		{
			suggestion:      newFakeSuggestion(),
			configMap:       newFakeKatibConfig(newFakeSuggestionConfig(), newFakeEarlyStoppingConfig()),
			err:             true,
			testDescription: "Set controller reference error",
		},
		{
			suggestion:      newFakeSuggestion(),
			err:             true,
			testDescription: "Get suggestion config error, not found Katib config",
		},
		{
			suggestion:      newFakeSuggestion(),
			configMap:       newFakeKatibConfig(newFakeSuggestionConfig(), newFakeEarlyStoppingConfig()),
			expectedPVC:     newFakePVC(),
			expectedPV:      nil,
			err:             false,
			testDescription: "Desired Volume valid run with default PVC",
		},
		{
			suggestion: newFakeSuggestion(),
			configMap: func() *corev1.ConfigMap {
				sc := newFakeSuggestionConfig()

				sc.PersistentVolumeClaimSpec = newFakePVC().Spec

				// Change StorageClass and volume storage.
				sc.PersistentVolumeClaimSpec.StorageClassName = &storageClassName

				sc.PersistentVolumeSpec = newFakePV().Spec
				// This policy will be changed to "Delete".
				sc.PersistentVolumeSpec.PersistentVolumeReclaimPolicy = corev1.PersistentVolumeReclaimRetain
				sc.PersistentVolumeLabels = pvLabels

				cm := newFakeKatibConfig(sc, newFakeEarlyStoppingConfig())
				return cm
			}(),
			expectedPVC: func() *corev1.PersistentVolumeClaim {
				pvc := newFakePVC()
				pvc.Spec.StorageClassName = &storageClassName
				return pvc
			}(),
			expectedPV:      newFakePV(),
			err:             false,
			testDescription: "Custom PVC and PV",
		},
	}

	for idx, tc := range tcs {

		if tc.configMap != nil {
			// Create ConfigMap with Katib config
			g.Expect(c.Create(ctx, tc.configMap)).NotTo(gomega.HaveOccurred())

			// Expect that ConfigMap is created
			g.Eventually(func() error {
				return c.Get(ctx, types.NamespacedName{Namespace: namespace, Name: configMap}, &corev1.ConfigMap{})
			}, timeout).ShouldNot(gomega.HaveOccurred())
		}

		// Get PVC and PV
		var actualPVC *corev1.PersistentVolumeClaim
		var actualPV *corev1.PersistentVolume
		var err error
		// For the first Test we run DesiredVolume with empty Scheme to fail Set Controller Reference
		if idx == 0 {
			c := General{
				scheme: &runtime.Scheme{},
				Client: mgr.GetClient(),
			}
			actualPVC, actualPV, err = c.DesiredVolume(tc.suggestion)
		} else {
			actualPVC, actualPV, err = composer.DesiredVolume(tc.suggestion)
		}

		if !tc.err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)

		} else if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)

		} else if !tc.err && ((tc.expectedPV == nil && actualPV != nil) || (tc.expectedPV != nil && actualPV == nil)) {
			t.Errorf("Case: %v failed. \nExpected PV: %v\n Got %v", tc.testDescription, tc.expectedPV, actualPV)

		} else if !tc.err && (!metaEqual(tc.expectedPVC.ObjectMeta, actualPVC.ObjectMeta) ||
			(tc.expectedPV != nil && !metaEqual(tc.expectedPV.ObjectMeta, actualPV.ObjectMeta))) {
			t.Errorf("Case: %v failed. \nExpected PVC metadata %v\n Got %v.\nExpected PV metadata %v\n Got %v",
				tc.testDescription, tc.expectedPVC.ObjectMeta, actualPVC.ObjectMeta, tc.expectedPV.ObjectMeta, actualPV.ObjectMeta)

		} else if !tc.err && (!equality.Semantic.DeepEqual(tc.expectedPVC.Spec, actualPVC.Spec) ||
			(tc.expectedPV != nil && !equality.Semantic.DeepEqual(tc.expectedPV.Spec, actualPV.Spec))) {
			t.Errorf("Case: %v failed. \nExpected PVC spec %v\n Got %v.\nExpected PV spec %v\n Got %v",
				tc.testDescription, tc.expectedPVC.Spec, actualPVC.Spec, tc.expectedPV, actualPV)
		}

		if tc.configMap != nil {
			// Delete ConfigMap with Katib config
			g.Expect(c.Delete(ctx, tc.configMap)).NotTo(gomega.HaveOccurred())
			// Expect that ConfigMap is deleted
			g.Eventually(func() bool {
				return errors.IsNotFound(
					c.Get(ctx, types.NamespacedName{Namespace: namespace, Name: configMap}, &corev1.ConfigMap{}))
			}, timeout).Should(gomega.BeTrue())
		}
	}
}

func TestDesiredRBAC(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mgr, err := manager.New(cfg, manager.Options{Metrics: metricsserver.Options{BindAddress: "0"}})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	composer := New(mgr)

	expectedServiceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      suggestionName + "-" + suggestionAlgorithm,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "kubeflow.org/v1beta1",
					Kind:               "Suggestion",
					Name:               suggestionName,
					Controller:         &refFlag,
					BlockOwnerDeletion: &refFlag,
				},
			},
		},
	}

	expectedRole := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      suggestionName + "-" + suggestionAlgorithm,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "kubeflow.org/v1beta1",
					Kind:               "Suggestion",
					Name:               suggestionName,
					Controller:         &refFlag,
					BlockOwnerDeletion: &refFlag,
				},
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					trialsv1beta1.Group,
				},
				Resources: []string{
					consts.PluralTrial,
					fmt.Sprintf("%v/status", consts.PluralTrial),
				},
				Verbs: []string{
					rbacv1.VerbAll,
				},
			},
		},
	}

	expectedRoleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      suggestionName + "-" + suggestionAlgorithm,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "kubeflow.org/v1beta1",
					Kind:               "Suggestion",
					Name:               suggestionName,
					Controller:         &refFlag,
					BlockOwnerDeletion: &refFlag,
				},
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      suggestionName + "-" + suggestionAlgorithm,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     suggestionName + "-" + suggestionAlgorithm,
		},
	}

	tcs := []struct {
		suggestion             *suggestionsv1beta1.Suggestion
		expectedServiceAccount *corev1.ServiceAccount
		expectedRole           *rbacv1.Role
		expectedRoleBinding    *rbacv1.RoleBinding
		err                    bool
		testDescription        string
	}{
		{
			suggestion:             newFakeSuggestion(),
			expectedServiceAccount: expectedServiceAccount,
			expectedRole:           expectedRole,
			expectedRoleBinding:    expectedRoleBinding,
			err:                    false,
			testDescription:        "Desired RBAC valid run",
		},
	}

	for _, tc := range tcs {

		actualServiceAccount, actualRole, actualRoleBinding, err := composer.DesiredRBAC(tc.suggestion)

		if !tc.err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
		} else if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		} else if !tc.err && (!equality.Semantic.DeepEqual(tc.expectedServiceAccount, actualServiceAccount) ||
			!equality.Semantic.DeepEqual(tc.expectedRole, actualRole) ||
			!equality.Semantic.DeepEqual(tc.expectedRoleBinding, actualRoleBinding)) {
			t.Errorf("Case: %v failed. \nExpected SA %v\n Got %v.\nExpected Role %v\n Got %v.\nExpected RoleBinding %v\n Got %v",
				tc.testDescription,
				tc.expectedServiceAccount, actualServiceAccount,
				tc.expectedRole, actualRole,
				tc.expectedRoleBinding, actualRoleBinding)
		}
	}
}

func metaEqual(expected, actual metav1.ObjectMeta) bool {
	return expected.Name == actual.Name &&
		expected.Namespace == actual.Namespace &&
		reflect.DeepEqual(expected.Labels, actual.Labels) &&
		reflect.DeepEqual(expected.Annotations, actual.Annotations) &&
		(len(actual.OwnerReferences) > 0 &&
			expected.OwnerReferences[0].APIVersion == actual.OwnerReferences[0].APIVersion &&
			expected.OwnerReferences[0].Kind == actual.OwnerReferences[0].Kind &&
			expected.OwnerReferences[0].Name == actual.OwnerReferences[0].Name &&
			*expected.OwnerReferences[0].Controller == *actual.OwnerReferences[0].Controller &&
			*expected.OwnerReferences[0].BlockOwnerDeletion == *actual.OwnerReferences[0].BlockOwnerDeletion ||
			len(actual.OwnerReferences) == 0)
}

func newFakeSuggestionConfig() configv1beta1.SuggestionConfig {
	cpuQ, _ := resource.ParseQuantity(cpu)
	memoryQ, _ := resource.ParseQuantity(memory)
	diskQ, _ := resource.ParseQuantity(disk)

	return configv1beta1.SuggestionConfig{
		AlgorithmName: suggestionAlgorithm,
		Container: corev1.Container{
			Image:           image,
			ImagePullPolicy: imagePullPolicy,
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
		},
		ServiceAccountName: serviceAccount,
	}
}

func newFakeEarlyStoppingConfig() configv1beta1.EarlyStoppingConfig {
	cpuQ, _ := resource.ParseQuantity(cpu)
	memoryQ, _ := resource.ParseQuantity(memory)
	diskQ, _ := resource.ParseQuantity(disk)

	return configv1beta1.EarlyStoppingConfig{
		AlgorithmName:   earlyStoppingAlgorithm,
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
	}
}

func newFakeKatibConfig(suggestionConfig configv1beta1.SuggestionConfig, earlyStoppingConfig configv1beta1.EarlyStoppingConfig) *corev1.ConfigMap {
	katibConfig := configv1beta1.KatibConfig{
		RuntimeConfig: configv1beta1.RuntimeConfig{
			SuggestionConfigs: []configv1beta1.SuggestionConfig{
				suggestionConfig,
			},
			EarlyStoppingConfigs: []configv1beta1.EarlyStoppingConfig{
				earlyStoppingConfig,
			},
		},
	}

	bKatibConfig, err := yaml.Marshal(katibConfig)
	if err != nil {
		stdlog.Fatal(err)
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMap,
			Namespace: namespace,
		},
		Data: map[string]string{
			consts.LabelKatibConfigTag: string(bKatibConfig),
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
			Requests: 1,
			Algorithm: &commonv1beta1.AlgorithmSpec{
				AlgorithmName: suggestionAlgorithm,
			},
			EarlyStopping: &commonv1beta1.EarlyStoppingSpec{
				AlgorithmName: earlyStoppingAlgorithm,
			},
			ResumePolicy: experimentsv1beta1.FromVolume,
		},
	}
}

func newFakeDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        suggestionName + "-" + suggestionAlgorithm,
			Namespace:   namespace,
			Labels:      suggestionLabels,
			Annotations: suggestionAnnotations,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "kubeflow.org/v1beta1",
					Kind:               "Suggestion",
					Name:               suggestionName,
					Controller:         &refFlag,
					BlockOwnerDeletion: &refFlag,
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
				ProbeHandler: corev1.ProbeHandler{
					GRPC: &corev1.GRPCAction{
						Port:    consts.DefaultSuggestionPort,
						Service: &consts.DefaultGRPCService,
					},
				},
				InitialDelaySeconds: defaultInitialDelaySeconds,
				PeriodSeconds:       defaultPeriodForReady,
			},
			LivenessProbe: &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					GRPC: &corev1.GRPCAction{
						Port:    consts.DefaultSuggestionPort,
						Service: &consts.DefaultGRPCService,
					},
				},
				InitialDelaySeconds: defaultInitialDelaySeconds,
				PeriodSeconds:       defaultPeriodForLive,
				FailureThreshold:    defaultFailureThreshold,
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      consts.ContainerSuggestionVolumeName,
					MountPath: configv1beta1.DefaultContainerSuggestionVolumeMountPath,
				},
			},
		},
		{
			Name:            consts.ContainerEarlyStopping,
			Image:           image,
			ImagePullPolicy: imagePullPolicy,
			Ports: []corev1.ContainerPort{
				{
					Name:          consts.DefaultEarlyStoppingPortName,
					ContainerPort: consts.DefaultEarlyStoppingPort,
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
		},
	}
}

func newFakePVC() *corev1.PersistentVolumeClaim {

	volumeStorage, _ := resource.ParseQuantity(configv1beta1.DefaultSuggestionVolumeStorage)

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      suggestionName + "-" + suggestionAlgorithm,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         "kubeflow.org/v1beta1",
					Kind:               "Suggestion",
					Name:               suggestionName,
					Controller:         &refFlag,
					BlockOwnerDeletion: &refFlag,
				},
			},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				configv1beta1.DefaultSuggestionVolumeAccessMode,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: volumeStorage,
				},
			},
		},
	}
}

func newFakePV() *corev1.PersistentVolume {
	pvName := suggestionName + "-" + suggestionAlgorithm + "-" + namespace
	volumeStorage, _ := resource.ParseQuantity(configv1beta1.DefaultSuggestionVolumeStorage)

	return &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pvName,
			Labels: pvLabels,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName:              storageClassName,
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				configv1beta1.DefaultSuggestionVolumeAccessMode,
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "tmp/katib/suggestion" + pvName,
				},
			},
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: volumeStorage,
			},
		},
	}
}
