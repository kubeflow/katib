package pod

import (
	"context"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	mccommon "github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
)

func TestWrapWorkerContainer(t *testing.T) {
	testCases := []struct {
		Pod         *v1.Pod
		Namespace   string
		JobKind     string
		MetricsFile string
		PathKind    common.FileSystemKind
		Trial       *trialsv1beta1.Trial
		Expected    *v1.Pod
		Err         bool
		Name        string
	}{
		{
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "tensorflow",
							Command: []string{
								"python main.py",
							},
						},
					},
				},
			},
			Namespace:   "nohere",
			JobKind:     "TFJob",
			MetricsFile: "testfile",
			PathKind:    common.FileKind,
			Trial: &trialsv1beta1.Trial{
				Spec: trialsv1beta1.TrialSpec{
					MetricsCollector: common.MetricsCollectorSpec{
						Collector: &common.CollectorSpec{
							Kind: common.StdOutCollector,
						},
					},
				},
			},
			Expected: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "tensorflow",
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								"python main.py 1>testfile 2>&1 && echo completed > $$$$.pid",
							},
						},
					},
				},
			},
			Err:  false,
			Name: "tensorflow container without sh -c",
		},
		{
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "tensorflow",
							Command: []string{
								"sh", "-c",
								"python main.py",
							},
						},
					},
				},
			},
			Namespace:   "nohere",
			JobKind:     "TFJob",
			MetricsFile: "testfile",
			PathKind:    common.FileKind,
			Trial: &trialsv1beta1.Trial{
				Spec: trialsv1beta1.TrialSpec{
					MetricsCollector: common.MetricsCollectorSpec{
						Collector: &common.CollectorSpec{
							Kind: common.StdOutCollector,
						},
					},
				},
			},
			Expected: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "tensorflow",
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								"python main.py 1>testfile 2>&1 && echo completed > $$$$.pid",
							},
						},
					},
				},
			},
			Err:  false,
			Name: "tensorflow container with sh -c",
		},
		{
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "primary-container",
							Command: []string{
								"sh", "-c",
								"python main.py",
							},
						},
						{
							Name: "not-primary-container",
							Command: []string{
								"sh", "-c",
								"python main.py",
							},
						},
					},
				},
			},
			MetricsFile: "testfile",
			PathKind:    common.FileKind,
			Trial: &trialsv1beta1.Trial{
				Spec: trialsv1beta1.TrialSpec{
					PrimaryContainerName: "primary-container",
					MetricsCollector: common.MetricsCollectorSpec{
						Collector: &common.CollectorSpec{
							Kind: common.StdOutCollector,
						},
					},
				},
			},
			Expected: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "primary-container",
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								"python main.py 1>testfile 2>&1 && echo completed > $$$$.pid",
							},
						},
						{
							Name: "not-primary-container",
							Command: []string{
								"sh", "-c",
								"python main.py",
							},
						},
					},
				},
			},
			Err:  false,
			Name: "Primary container name is set for training pod",
		},
		{
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "not-primary-container",
						},
					},
				},
			},
			PathKind: common.FileKind,
			Trial: &trialsv1beta1.Trial{
				Spec: trialsv1beta1.TrialSpec{
					PrimaryContainerName: "primary-container",
				},
			},
			Err:  true,
			Name: "Training pod doesn't have primary container name",
		},
	}

	for _, c := range testCases {
		err := wrapWorkerContainer(c.Pod, c.Namespace, c.JobKind, c.MetricsFile, c.PathKind, c.Trial)
		if c.Err && err == nil {
			t.Errorf("Case %s failed. Expected error, got nil", c.Name)
		} else if !c.Err {
			if err != nil {
				t.Errorf("Case %s failed. Expected nil, got error: %v", c.Name, err)
			} else if !equality.Semantic.DeepEqual(c.Pod.Spec.Containers, c.Expected.Spec.Containers) {
				t.Errorf("Case %s failed. Expected pod: %v, got: %v",
					c.Name, c.Expected.Spec.Containers, c.Pod.Spec.Containers)
			}
		}

	}
}

func TestGetMetricsCollectorArgs(t *testing.T) {
	testTrialName := "test-trial"
	testMetricName := "accuracy"
	katibDBAddress := "katib-db-manager.kubeflow:6789"
	testPath := "/test/path"
	testCases := []struct {
		TrialName    string
		MetricName   string
		MCSpec       common.MetricsCollectorSpec
		ExpectedArgs []string
		Name         string
	}{
		{
			TrialName:  testTrialName,
			MetricName: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-s", katibDBAddress,
				"-path", common.DefaultFilePath,
			},
			Name: "StdOut MC",
		},
		{
			TrialName:  testTrialName,
			MetricName: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.FileCollector,
				},
				Source: &common.SourceSpec{
					FileSystemPath: &common.FileSystemPath{
						Path: testPath,
					},
					Filter: &common.FilterSpec{
						MetricsFormat: []string{
							"{mn1: ([a-b]), mv1: [0-9]}",
							"{mn2: ([a-b]), mv2: ([0-9])}",
						},
					},
				},
			},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-s", katibDBAddress,
				"-path", testPath,
				"-f", "{mn1: ([a-b]), mv1: [0-9]};{mn2: ([a-b]), mv2: ([0-9])}",
			},
			Name: "File MC with Filter",
		},
		{
			TrialName:  testTrialName,
			MetricName: testMetricName,
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
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-s", katibDBAddress,
				"-path", testPath,
			},
			Name: "Tf Event MC",
		},
		{
			TrialName:  testTrialName,
			MetricName: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.CustomCollector,
				},
			},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-s", katibDBAddress,
			},
			Name: "Custom MC without Path",
		},
		{
			TrialName:  testTrialName,
			MetricName: testMetricName,
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
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-s", katibDBAddress,
				"-path", testPath,
			},
			Name: "Custom MC with Path",
		},
		{
			TrialName:  testTrialName,
			MetricName: testMetricName,
			MCSpec: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.PrometheusMetricCollector,
				},
			},
			ExpectedArgs: []string{
				"-t", testTrialName,
				"-m", testMetricName,
				"-s", katibDBAddress,
			},
			Name: "Prometheus MC without Path",
		},
	}

	for _, tc := range testCases {
		args := getMetricsCollectorArgs(tc.TrialName, tc.MetricName, tc.MCSpec)
		if !reflect.DeepEqual(tc.ExpectedArgs, args) {
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

func TestMutateVolume(t *testing.T) {
	tc := struct {
		Pod                  v1.Pod
		ExpectedPod          v1.Pod
		JobKind              string
		MountPath            string
		SidecarContainerName string
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
		JobKind:              "Job",
		MountPath:            common.DefaultFilePath,
		SidecarContainerName: "train-job",
		PathKind:             common.FileKind,
	}

	err := mutateVolume(
		&tc.Pod,
		tc.JobKind,
		tc.MountPath,
		tc.SidecarContainerName,
		tc.PathKind)
	if err != nil {
		t.Errorf("mutateVolume failed: %v", err)
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

func TestGetKatibJob(t *testing.T) {
	// Start test k8s server
	envTest := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "manifests", "v1beta1", "katib-controller"),
			filepath.Join("..", "..", "..", "..", "test", "unit", "v1beta1", "crds"),
		},
	}

	cfg, err := envTest.Start()
	if err != nil {
		t.Error(err)
	}

	g := gomega.NewGomegaWithT(t)

	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)
	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	c := mgr.GetClient()
	si := NewSidecarInjector(c)

	namespace := "default"
	trialName := "trial-name"
	podName := "pod-name"
	deployName := "deploy-name"
	tfJobName := "tfjob-name"
	timeout := time.Second * 5

	testCases := []struct {
		Pod             *v1.Pod
		TFJob           *tfv1.TFJob
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
							APIVersion: "kubeflow.org/v1",
							Kind:       "TFJob",
							Name:       tfJobName + "-1",
						},
					},
				},
			},
			TFJob: &tfv1.TFJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tfJobName + "-1",
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
			},
			ExpectedJobKind: "TFJob",
			ExpectedJobName: tfJobName + "-1",
			Err:             false,
			TestDescription: "Valid run with ownership sequence: Trial -> TFJob -> Pod",
		},
		{
			Pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "kubeflow.org/v1",
							Kind:       "TFJob",
							Name:       tfJobName + "-2",
						},
						{
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Name:       deployName + "-2",
						},
					},
				},
			},
			TFJob: &tfv1.TFJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tfJobName + "-2",
					Namespace: namespace,
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
			TestDescription: "Valid run with ownership sequence: Trial -> Deployment -> Pod, TFJob -> Pod",
		},
		{
			Pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "kubeflow.org/v1",
							Kind:       "TFJob",
							Name:       tfJobName + "-3",
						},
					},
				},
			},
			TFJob: &tfv1.TFJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tfJobName + "-3",
					Namespace: namespace,
				},
			},
			Err:             true,
			TestDescription: "Run for not Trial's pod with ownership sequence: TFJob -> Pod",
		},
		{
			Pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "kubeflow.org/v1",
							Kind:       "TFJob",
							Name:       tfJobName + "-4",
						},
					},
				},
			},
			Err:             true,
			TestDescription: "Run when Pod owns TFJob that doesn't exists",
		},
		{
			Pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "invalid/api/version",
							Kind:       "TFJob",
							Name:       tfJobName + "-4",
						},
					},
				},
			},
			Err:             true,
			TestDescription: "Run when Pod owns TFJob with invalid API version",
		},
	}

	for _, tc := range testCases {
		// Create TFJob if it is needed
		if tc.TFJob != nil {
			tfJobUnstr, err := util.ConvertObjectToUnstructured(tc.TFJob)
			gvk := schema.GroupVersionKind{
				Group:   "kubeflow.org",
				Version: "v1",
				Kind:    "TFJob",
			}
			tfJobUnstr.SetGroupVersionKind(gvk)
			if err != nil {
				t.Errorf("ConvertObjectToUnstructured error %v", err)
			}

			g.Expect(c.Create(context.TODO(), tfJobUnstr)).NotTo(gomega.HaveOccurred())

			// Wait that TFJob is created
			g.Eventually(func() error {
				return c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: tc.TFJob.Name}, tfJobUnstr)
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

func TestIsMasterRole(t *testing.T) {
	masterRoleLabel := make(map[string]string)
	masterRoleLabel[consts.JobRole] = MasterRole
	invalidLabel := make(map[string]string)
	invalidLabel["invalid-label"] = "invalid"
	testCases := []struct {
		Pod      v1.Pod
		JobKind  string
		IsMaster bool
		Name     string
	}{
		{
			JobKind:  "Job",
			IsMaster: true,
			Name:     "Kubernetes Batch Job Pod",
		},
		{
			Pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: masterRoleLabel,
				},
			},
			JobKind:  "PyTorchJob",
			IsMaster: true,
			Name:     "Pytorch Master Pod",
		},
		{
			Pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: invalidLabel,
				},
			},
			JobKind:  "PyTorchJob",
			IsMaster: false,
			Name:     "Pytorch Pod with invalid label",
		},
	}

	for _, tc := range testCases {
		isMaster := isMasterRole(&tc.Pod, tc.JobKind)
		if isMaster != tc.IsMaster {
			t.Errorf("Case %v. Expected isMaster %v, got %v", tc.Name, tc.IsMaster, isMaster)
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
