package pod

import (
	"reflect"
	"testing"

	"path/filepath"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	mccommon "github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWrapWorkerContainer(t *testing.T) {
	testCases := []struct {
		Pod           *v1.Pod
		Namespace     string
		JobKind       string
		MetricsFile   string
		PathKind      common.FileSystemKind
		MC            common.MetricsCollectorSpec
		Expected      *v1.Pod
		ExpectedError error
		Name          string
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
			MC: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
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
			ExpectedError: nil,
			Name:          "tensorflow container without sh -c",
		},
		{
			Pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "test",
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
			MC: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
				},
			},
			Expected: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "test",
							Command: []string{
								"python main.py",
							},
						},
					},
				},
			},
			ExpectedError: nil,
			Name:          "test container without sh -c",
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
			MC: common.MetricsCollectorSpec{
				Collector: &common.CollectorSpec{
					Kind: common.StdOutCollector,
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
			ExpectedError: nil,
			Name:          "Tensorflow container with sh -c",
		},
	}

	for _, c := range testCases {
		err := wrapWorkerContainer(c.Pod, c.Namespace, c.JobKind, c.MetricsFile, c.PathKind, c.MC)
		if err != c.ExpectedError {
			t.Errorf("Expected error %v, got %v", c.ExpectedError, err)
		}
		if err == nil {
			if !equality.Semantic.DeepEqual(c.Pod.Spec.Containers, c.Expected.Spec.Containers) {
				t.Errorf("Case %s: Expected pod %v, got %v",
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

func TestGetKatibJob(t *testing.T) {
	testCases := []struct {
		Pod             v1.Pod
		ExpectedJobKind string
		ExpectedJobName string
		Err             bool
		Name            string
	}{
		{
			Pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "kubeflow.org/v1",
							Kind:       "PyTorchJob",
							Name:       "OwnerName",
						},
					},
				},
			},
			ExpectedJobKind: "PyTorchJob",
			ExpectedJobName: "OwnerName",
			Err:             false,
			Name:            "Valid Pod",
		},
		{
			Pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "notkubeflow.org/v1",
							Kind:       "PyTorchJob",
							Name:       "OwnerName",
						},
					},
				},
			},
			Err:  true,
			Name: "Invalid APIVersion",
		},
		{
			Pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "kubeflow.org/v1",
							Kind:       "MXJob",
							Name:       "OwnerName",
						},
					},
				},
			},
			Err:  true,
			Name: "Invalid Kind",
		},
	}

	for _, tc := range testCases {
		jobKind, jobName, err := getKatibJob(&tc.Pod)
		if !tc.Err && err != nil {
			t.Errorf("Case %v failed. Error %v", tc.Name, err)
		} else if !tc.Err && (tc.ExpectedJobKind != jobKind || tc.ExpectedJobName != jobName) {
			t.Errorf("Case %v failed. Expected jobKind %v, got %v, Expected jobName %v, got %v",
				tc.Name, tc.ExpectedJobKind, jobKind, tc.ExpectedJobName, jobName)
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
