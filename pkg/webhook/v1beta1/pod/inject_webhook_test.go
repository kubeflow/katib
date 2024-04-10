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

package pod

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	apis "github.com/kubeflow/katib/pkg/apis/controller"
	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	mccommon "github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
)

var (
	timeout = time.Second * 5
)

func TestWrapWorkerContainer(t *testing.T) {
	primaryContainer := "tensorflow"
	trial := &trialsv1beta1.Trial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trial-name",
			Namespace: "trial-namespace",
		},
		Spec: trialsv1beta1.TrialSpec{
			MetricsCollector: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			PrimaryContainerName: primaryContainer,
			SuccessCondition:     experimentsv1beta1.DefaultJobSuccessCondition,
			FailureCondition:     experimentsv1beta1.DefaultJobFailureCondition,
		},
	}

	metricsFile := "metric.log"

	testCases := []struct {
		Trial           *trialsv1beta1.Trial
		Pod             *v1.Pod
		MetricsFile     string
		PathKind        common.FileSystemKind
		ExpectedPod     *v1.Pod
		Err             bool
		TestDescription string
	}{
		{
			Trial: trial,
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: primaryContainer,
							Command: []string{
								"python main.py",
							},
						},
					},
				},
			},
			MetricsFile: metricsFile,
			PathKind:    common.FileKind,
			ExpectedPod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: primaryContainer,
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								fmt.Sprintf("python main.py 1>%v 2>&1 && echo completed > $$$$.pid", metricsFile),
							},
						},
					},
				},
			},
			Err:             false,
			TestDescription: "Tensorflow container without sh -c",
		},
		{
			Trial: trial,
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: primaryContainer,
							Command: []string{
								"sh", "-c",
								"python main.py",
							},
						},
					},
				},
			},
			MetricsFile: metricsFile,
			PathKind:    common.FileKind,
			ExpectedPod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: primaryContainer,
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								fmt.Sprintf("python main.py 1>%v 2>&1 && echo completed > $$$$.pid", metricsFile),
							},
						},
					},
				},
			},
			Err:             false,
			TestDescription: "Tensorflow container with sh -c",
		},
		{
			Trial: trial,
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "not-primary-container",
						},
					},
				},
			},
			PathKind:        common.FileKind,
			Err:             true,
			TestDescription: "Training pod doesn't have primary container",
		},
		{
			Trial: func() *trialsv1beta1.Trial {
				t := trial.DeepCopy()
				t.Spec.EarlyStoppingRules = []common.EarlyStoppingRule{
					{
						Name:       "accuracy",
						Value:      "0.6",
						Comparison: common.ComparisonTypeLess,
					},
				}
				return t
			}(),
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: primaryContainer,
							Command: []string{
								"python main.py",
							},
						},
					},
				},
			},
			MetricsFile: metricsFile,
			PathKind:    common.FileKind,
			ExpectedPod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: primaryContainer,
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								fmt.Sprintf("python main.py 1>%v 2>&1 || "+
									"if test -f $$$$.pid && [ $(head -n 1 $$.pid) = early-stopped ]; then "+
									"echo Training Container was Early Stopped; "+
									"else echo Training Container was Failed; exit 1; fi "+
									"&& echo completed > $$$$.pid", metricsFile),
							},
						},
					},
				},
			},
			Err:             false,
			TestDescription: "Container with early stopping command",
		},
	}

	for _, c := range testCases {
		err := wrapWorkerContainer(c.Trial, c.Pod, c.Trial.Namespace, c.MetricsFile, c.PathKind)
		if c.Err && err == nil {
			t.Errorf("Case %s failed. Expected error, got nil", c.TestDescription)
		} else if !c.Err {
			if err != nil {
				t.Errorf("Case %s failed. Expected nil, got error: %v", c.TestDescription, err)
			} else if !equality.Semantic.DeepEqual(c.Pod.Spec.Containers, c.ExpectedPod.Spec.Containers) {
				t.Errorf("Case %s failed. Expected pod: %v, got: %v",
					c.TestDescription, c.ExpectedPod.Spec.Containers, c.Pod.Spec.Containers)
			}
		}
	}
}

func TestGetMetricsCollectorArgs(t *testing.T) {

	// Start test k8s server
	envTest := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "manifests", "v1beta1", "components", "crd"),
		},
	}
	if err := apis.AddToScheme(scheme.Scheme); err != nil {
		t.Error(err)
	}

	cfg, err := envTest.Start()
	if err != nil {
		t.Error(err)
	}

	g := gomega.NewGomegaWithT(t)

	mgr, err := manager.New(cfg, manager.Options{Metrics: metricsserver.Options{BindAddress: "0"}})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// Start test manager.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(context.TODO())).NotTo(gomega.HaveOccurred())
	}()

	c := mgr.GetClient()
	si := NewSidecarInjector(c, admission.NewDecoder(mgr.GetScheme()))

	testTrialName := "test-trial"
	testSuggestionName := "test-suggestion"
	testNamespace := "kubeflow"
	testAlgorithm := "random"
	testObjective := common.ObjectiveTypeMaximize
	testMetricName := "accuracy"
	katibDBAddress := fmt.Sprintf("katib-db-manager.%v:%v", testNamespace, consts.DefaultSuggestionPort)
	katibEarlyStopAddress := fmt.Sprintf("%v-%v.%v:%v", testSuggestionName, testAlgorithm, testNamespace, consts.DefaultEarlyStoppingPort)
	waitAllProcessesValue := false
	testPath := "/test/path"

	// Create kubeflow namespace.
	kubeflowNS := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
		},
	}
	g.Expect(c.Create(context.TODO(), kubeflowNS)).NotTo(gomega.HaveOccurred())

	earlyStoppingRules := []string{
		"accuracy;0.6;less;5",
		"loss;2;greater",
	}

	testSuggestion := &suggestionsv1beta1.Suggestion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSuggestionName,
			Namespace: testNamespace,
		},
		Spec: suggestionsv1beta1.SuggestionSpec{
			Algorithm: &common.AlgorithmSpec{
				AlgorithmName: testAlgorithm,
			},
		},
	}

	testTrial := &trialsv1beta1.Trial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testTrialName,
			Namespace: testNamespace,
			Labels: map[string]string{
				consts.LabelExperimentName: testSuggestionName,
			},
		},
		Spec: trialsv1beta1.TrialSpec{
			Objective: &common.ObjectiveSpec{
				Type: testObjective,
			},
		},
	}

	testCases := []struct {
		Trial              *trialsv1beta1.Trial
		MetricNames        string
		MCSpec             common.MetricsCollectorSpec
		EarlyStoppingRules []string
		KatibConfig        configv1beta1.MetricsCollectorConfig
		ExpectedArgs       []string
		Name               string
		Err                bool
	}{
		{
			Trial:       testTrial,
			MetricNames: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			KatibConfig: configv1beta1.MetricsCollectorConfig{
				WaitAllProcesses: &waitAllProcessesValue,
			},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-o-type", string(testObjective),
				"-s-db", katibDBAddress,
				"-path", common.DefaultFilePath,
				"-format", string(common.TextFormat),
				"-w", "false",
			},
			Name: "StdOut MC",
		},
		{
			Trial:       testTrial,
			MetricNames: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.FileCollector,
				},
				Source: &common.SourceSpec{
					FileSystemPath: &common.FileSystemPath{
						Path:   testPath,
						Format: common.TextFormat,
					},
					Filter: &common.FilterSpec{
						MetricsFormat: []string{
							"{mn1: ([a-b]), mv1: [0-9]}",
							"{mn2: ([a-b]), mv2: ([0-9])}",
						},
					},
				},
			},
			KatibConfig: configv1beta1.MetricsCollectorConfig{},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-o-type", string(testObjective),
				"-s-db", katibDBAddress,
				"-path", testPath,
				"-f", "{mn1: ([a-b]), mv1: [0-9]};{mn2: ([a-b]), mv2: ([0-9])}",
				"-format", string(common.TextFormat),
			},
			Name: "File MC with Filter",
		},
		{
			Trial:       testTrial,
			MetricNames: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.FileCollector,
				},
				Source: &common.SourceSpec{
					FileSystemPath: &common.FileSystemPath{
						Path:   testPath,
						Format: common.JsonFormat,
					},
				},
			},
			KatibConfig: configv1beta1.MetricsCollectorConfig{},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-o-type", string(testObjective),
				"-s-db", katibDBAddress,
				"-path", testPath,
				"-format", string(common.JsonFormat),
			},
			Name: "File MC with Json Format",
		},
		{
			Trial:       testTrial,
			MetricNames: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.TfEventCollector,
				},
				Source: &common.SourceSpec{
					FileSystemPath: &common.FileSystemPath{
						Path: testPath,
					},
				},
			},
			KatibConfig: configv1beta1.MetricsCollectorConfig{},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-o-type", string(testObjective),
				"-s-db", katibDBAddress,
				"-path", testPath,
			},
			Name: "Tf Event MC",
		},
		{
			Trial:       testTrial,
			MetricNames: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.CustomCollector,
				},
			},
			KatibConfig: configv1beta1.MetricsCollectorConfig{},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-o-type", string(testObjective),
				"-s-db", katibDBAddress,
			},
			Name: "Custom MC without Path",
		},
		{
			Trial:       testTrial,
			MetricNames: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.CustomCollector,
				},
				Source: &common.SourceSpec{
					FileSystemPath: &common.FileSystemPath{
						Path: testPath,
					},
				},
			},
			KatibConfig: configv1beta1.MetricsCollectorConfig{},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-o-type", string(testObjective),
				"-s-db", katibDBAddress,
				"-path", testPath,
			},
			Name: "Custom MC with Path",
		},
		{
			Trial:       testTrial,
			MetricNames: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.PrometheusMetricCollector,
				},
			},
			KatibConfig: configv1beta1.MetricsCollectorConfig{},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-o-type", string(testObjective),
				"-s-db", katibDBAddress,
			},
			Name: "Prometheus MC without Path",
		},
		{
			Trial:       testTrial,
			MetricNames: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			EarlyStoppingRules: earlyStoppingRules,
			KatibConfig:        configv1beta1.MetricsCollectorConfig{},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-o-type", string(testObjective),
				"-s-db", katibDBAddress,
				"-path", common.DefaultFilePath,
				"-format", string(common.TextFormat),
				"-stop-rule", earlyStoppingRules[0],
				"-stop-rule", earlyStoppingRules[1],
				"-s-earlystop", katibEarlyStopAddress,
			},
			Name: "Trial with EarlyStopping rules",
		},
		{
			Trial: func() *trialsv1beta1.Trial {
				trial := testTrial.DeepCopy()
				trial.ObjectMeta.Labels[consts.LabelExperimentName] = "invalid-name"
				return trial
			}(),
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			EarlyStoppingRules: earlyStoppingRules,
			KatibConfig:        configv1beta1.MetricsCollectorConfig{},
			Name:               "Trial with invalid Experiment label name. Suggestion is not created",
			Err:                true,
		},
	}

	g.Expect(c.Create(context.TODO(), testSuggestion)).NotTo(gomega.HaveOccurred())

	// Wait that Suggestion is created
	g.Eventually(func() error {
		return c.Get(context.TODO(), types.NamespacedName{Namespace: testNamespace, Name: testSuggestionName}, testSuggestion)
	}, timeout).ShouldNot(gomega.HaveOccurred())

	for _, tc := range testCases {
		args, err := si.getMetricsCollectorArgs(tc.Trial, tc.MetricNames, tc.MCSpec, tc.KatibConfig, tc.EarlyStoppingRules)

		if !tc.Err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.Name, err)
		} else if tc.Err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.Name)
		} else if !tc.Err && !reflect.DeepEqual(tc.ExpectedArgs, args) {
			t.Errorf("Case %v failed. ExpectedArgs: %v, got %v", tc.Name, tc.ExpectedArgs, args)
		}
	}
}

func TestNeedWrapWorkerContainer(t *testing.T) {
	testCases := []struct {
		MCSpec   common.MetricsCollectorSpec
		needWrap bool
	}{
		{
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			needWrap: true,
		},
		{
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.CustomCollector,
				},
			},
			needWrap: false,
		},
	}

	for _, tc := range testCases {
		needWrap := needWrapWorkerContainer(tc.MCSpec)
		if needWrap != tc.needWrap {
			t.Errorf("Expected needWrap %v, got %v", tc.needWrap, needWrap)
		}
	}
}

func TestMutateMetricsCollectorVolume(t *testing.T) {
	tc := struct {
		Pod                  v1.Pod
		ExpectedPod          v1.Pod
		JobKind              string
		MountPath            string
		SidecarContainerName string
		PrimaryContainerName string
		PathKind             common.FileSystemKind
		Err                  bool
	}{
		Pod: v1.Pod{
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name: "train-job",
					},
					{
						Name: "init-container",
					},
					{
						Name: "metrics-collector",
					},
				},
			},
		},
		ExpectedPod: v1.Pod{
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name: "train-job",
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      common.MetricsVolume,
								MountPath: filepath.Dir(common.DefaultFilePath),
							},
						},
					},
					{
						Name: "init-container",
					},
					{
						Name: "metrics-collector",
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      common.MetricsVolume,
								MountPath: filepath.Dir(common.DefaultFilePath),
							},
						},
					},
				},
				Volumes: []v1.Volume{
					{
						Name: common.MetricsVolume,
						VolumeSource: v1.VolumeSource{
							EmptyDir: &v1.EmptyDirVolumeSource{},
						},
					},
				},
			},
		},
		MountPath:            common.DefaultFilePath,
		SidecarContainerName: "metrics-collector",
		PrimaryContainerName: "train-job",
		PathKind:             common.FileKind,
	}

	err := mutateMetricsCollectorVolume(
		&tc.Pod,
		tc.MountPath,
		tc.SidecarContainerName,
		tc.PrimaryContainerName,
		tc.PathKind)
	if err != nil {
		t.Errorf("mutateMetricsCollectorVolume failed: %v", err)
	} else if !equality.Semantic.DeepEqual(tc.Pod, tc.ExpectedPod) {
		t.Errorf("Expected pod %v, got %v", tc.ExpectedPod, tc.Pod)
	}
}

func TestGetSidecarContainerName(t *testing.T) {
	testCases := []struct {
		CollectorKind         common.CollectorKind
		ExpectedCollectorKind string
	}{
		{
			CollectorKind:         common.StdOutCollector,
			ExpectedCollectorKind: mccommon.MetricLoggerCollectorContainerName,
		},
		{
			CollectorKind:         common.TfEventCollector,
			ExpectedCollectorKind: mccommon.MetricCollectorContainerName,
		},
	}

	for _, tc := range testCases {
		collectorKind := getSidecarContainerName(tc.CollectorKind)
		if collectorKind != tc.ExpectedCollectorKind {
			t.Errorf("Expected Collector Kind: %v, got %v", tc.ExpectedCollectorKind, collectorKind)
		}
	}
}

func TestGetKatibJob(t *testing.T) {
	// Start test k8s server
	envTest := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "manifests", "v1beta1", "components", "crd"),
		},
	}
	if err := apis.AddToScheme(scheme.Scheme); err != nil {
		t.Error(err)
	}

	cfg, err := envTest.Start()
	if err != nil {
		t.Error(err)
	}

	g := gomega.NewGomegaWithT(t)

	mgr, err := manager.New(cfg, manager.Options{Metrics: metricsserver.Options{BindAddress: "0"}})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// Start test manager.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Expect(mgr.Start(context.TODO())).NotTo(gomega.HaveOccurred())
	}()

	c := mgr.GetClient()
	si := NewSidecarInjector(c, admission.NewDecoder(mgr.GetScheme()))

	namespace := "default"
	trialName := "trial-name"
	podName := "pod-name"
	deployName := "deploy-name"
	jobName := "job-name"

	testCases := []struct {
		Pod             *v1.Pod
		Job             *batchv1.Job
		Deployment      *appsv1.Deployment
		ExpectedJobKind string
		ExpectedJobName string
		Err             bool
		TestDescription string
	}{
		{
			Pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "batch/v1",
							Kind:       "Job",
							Name:       jobName + "-1",
						},
					},
				},
			},
			Job: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      jobName + "-1",
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "kubeflow.org/v1beta1",
							Kind:       "Trial",
							Name:       trialName + "-1",
							UID:        "test-uid",
						},
					},
				},
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							RestartPolicy: v1.RestartPolicyNever,
							Containers: []v1.Container{
								{
									Name:  "test",
									Image: "test",
								},
							},
						},
					},
				},
			},
			ExpectedJobKind: "Job",
			ExpectedJobName: jobName + "-1",
			Err:             false,
			TestDescription: "Valid run with ownership sequence: Trial -> Job -> Pod",
		},
		{
			Pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "batch/v1",
							Kind:       "Job",
							Name:       jobName + "-2",
						},
						{
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Name:       deployName + "-2",
						},
					},
				},
			},
			Job: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      jobName + "-2",
					Namespace: namespace,
				},
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							RestartPolicy: v1.RestartPolicyNever,
							Containers: []v1.Container{
								{
									Name:  "test",
									Image: "test",
								},
							},
						},
					},
				},
			},
			Deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deployName + "-2",
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "kubeflow.org/v1beta1",
							Kind:       "Trial",
							Name:       trialName + "-2",
							UID:        "test-uid",
						},
					},
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"test-key": "test-value",
						},
					},
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"test-key": "test-value",
							},
						},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  "test",
									Image: "test",
								},
							},
						},
					},
				},
			},
			ExpectedJobKind: "Deployment",
			ExpectedJobName: deployName + "-2",
			Err:             false,
			TestDescription: "Valid run with ownership sequence: Trial -> Deployment -> Pod, Job -> Pod",
		},
		{
			Pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "batch/v1",
							Kind:       "Job",
							Name:       jobName + "-3",
						},
					},
				},
			},
			Job: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      jobName + "-3",
					Namespace: namespace,
				},
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							RestartPolicy: v1.RestartPolicyNever,
							Containers: []v1.Container{
								{
									Name:  "test",
									Image: "test",
								},
							},
						},
					},
				},
			},
			Err:             true,
			TestDescription: "Run for not Trial's pod with ownership sequence: Job -> Pod",
		},
		{
			Pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "batch/v1",
							Kind:       "Job",
							Name:       jobName + "-4",
						},
					},
				},
			},
			Err:             true,
			TestDescription: "Run when Pod owns Job that doesn't exists",
		},
		{
			Pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "invalid/api/version",
							Kind:       "Job",
							Name:       jobName + "-4",
						},
					},
				},
			},
			Err:             true,
			TestDescription: "Run when Pod owns Job with invalid API version",
		},
	}

	for _, tc := range testCases {
		// Create Job if it is needed
		if tc.Job != nil {
			jobUnstr, err := util.ConvertObjectToUnstructured(tc.Job)
			gvk := schema.GroupVersionKind{
				Group:   "batch",
				Version: "v1",
				Kind:    "Job",
			}
			jobUnstr.SetGroupVersionKind(gvk)
			if err != nil {
				t.Errorf("ConvertObjectToUnstructured error %v", err)
			}

			g.Expect(c.Create(context.TODO(), jobUnstr)).NotTo(gomega.HaveOccurred())

			// Wait that Job is created
			g.Eventually(func() error {
				return c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: tc.Job.Name}, jobUnstr)
			}, timeout).ShouldNot(gomega.HaveOccurred())
		}

		// Create Deployment if it is needed
		if tc.Deployment != nil {
			g.Expect(c.Create(context.TODO(), tc.Deployment)).NotTo(gomega.HaveOccurred())

			// Wait that Deployment is created
			g.Eventually(func() error {
				return c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: tc.Deployment.Name}, tc.Deployment)
			}, timeout).ShouldNot(gomega.HaveOccurred())
		}

		object, _ := util.ConvertObjectToUnstructured(tc.Pod)
		jobKind, jobName, err := si.getKatibJob(object, namespace)
		if !tc.Err && err != nil {
			t.Errorf("Case %v failed. Error %v", tc.TestDescription, err)
		} else if !tc.Err && (tc.ExpectedJobKind != jobKind || tc.ExpectedJobName != jobName) {
			t.Errorf("Case %v failed. Expected jobKind %v, got %v, Expected jobName %v, got %v",
				tc.TestDescription, tc.ExpectedJobKind, jobKind, tc.ExpectedJobName, jobName)
		} else if tc.Err && err == nil {
			t.Errorf("Expected error got nil")
		}
	}
}

func TestIsPrimaryPod(t *testing.T) {
	testCases := []struct {
		podLabels        map[string]string
		primaryPodLabels map[string]string
		isPrimary        bool
		testDescription  string
	}{
		{
			podLabels: map[string]string{
				"test-key-1": "test-value-1",
				"test-key-2": "test-value-2",
				"test-key-3": "test-value-3",
			},
			primaryPodLabels: map[string]string{
				"test-key-1": "test-value-1",
				"test-key-2": "test-value-2",
			},
			isPrimary:       true,
			testDescription: "Pod contains all labels from primary pod labels",
		},
		{
			podLabels: map[string]string{
				"test-key-1": "test-value-1",
			},
			primaryPodLabels: map[string]string{
				"test-key-1": "test-value-1",
				"test-key-2": "test-value-2",
			},
			isPrimary:       false,
			testDescription: "Pod doesn't contain primary label",
		},
		{
			podLabels: map[string]string{
				"test-key-1": "invalid",
			},
			primaryPodLabels: map[string]string{
				"test-key-1": "test-value-1",
			},
			isPrimary:       false,
			testDescription: "Pod contains label with incorrect value",
		},
	}

	for _, tc := range testCases {
		isPrimary := isPrimaryPod(tc.podLabels, tc.primaryPodLabels)
		if isPrimary != tc.isPrimary {
			t.Errorf("Case %v. Expected isPrimary %v, got %v", tc.testDescription, tc.isPrimary, isPrimary)
		}
	}
}

func TestMutatePodMetadata(t *testing.T) {
	mutatedPodLabels := map[string]string{
		"custom-pod-label":    "custom-value",
		"katib-experiment":    "katib-value",
		consts.LabelTrialName: "test-trial",
	}

	testCases := []struct {
		pod             *v1.Pod
		trial           *trialsv1beta1.Trial
		mutatedPod      *v1.Pod
		testDescription string
	}{
		{
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"custom-pod-label": "custom-value",
					},
				},
			},
			trial: &trialsv1beta1.Trial{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-trial",
					Labels: map[string]string{
						"katib-experiment": "katib-value",
					},
				},
			},
			mutatedPod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: mutatedPodLabels,
				},
			},
			testDescription: "Mutated Pod should contain label from the origin Pod and Trial",
		},
	}

	for _, tc := range testCases {
		mutatePodMetadata(tc.pod, tc.trial)
		if !reflect.DeepEqual(tc.mutatedPod, tc.pod) {
			t.Errorf("Case %v. Expected Pod %v, got %v", tc.testDescription, tc.mutatedPod, tc.pod)
		}
	}
}
