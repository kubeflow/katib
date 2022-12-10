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
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"github.com/google/go-containerregistry/pkg/name"
	crv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v1 "k8s.io/api/core/v1"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	mccommon "github.com/kubeflow/katib/pkg/metricscollector/v1beta1/common"
)

func isPrimaryPod(podLabels, primaryLabels map[string]string) bool {

	for primaryKey, primaryValue := range primaryLabels {
		if podValue, ok := podLabels[primaryKey]; ok {
			if podValue != primaryValue {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func getPrimaryContainerIndex(containers []v1.Container, primaryContainerName string) int {
	primaryContainerIndex := -1
	for i, c := range containers {
		if c.Name == primaryContainerName {
			primaryContainerIndex = i
			break
		}
	}
	return primaryContainerIndex
}

func getRemoteImage(pod *v1.Pod, namespace string, containerIndex int) (crv1.Image, error) {
	// verify the image name, then download the remote config file
	c := pod.Spec.Containers[containerIndex]
	ref, err := name.ParseReference(c.Image, name.WeakValidation)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse image %q: %v", c.Image, err)
	}
	imagePullSecrets := []string{}
	for _, s := range pod.Spec.ImagePullSecrets {
		imagePullSecrets = append(imagePullSecrets, s.Name)
	}
	kc, err := k8schain.NewInCluster(context.TODO(),
		k8schain.Options{
			Namespace:          namespace,
			ServiceAccountName: pod.Spec.ServiceAccountName,
			ImagePullSecrets:   imagePullSecrets,
		})
	if err != nil {
		return nil, fmt.Errorf("Failed to create k8schain: %v", err)
	}

	mkc := authn.NewMultiKeychain(kc)
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(mkc))
	if err != nil {
		return nil, fmt.Errorf("Failed to get container image %q info from registry: %v", c.Image, err)
	}

	return img, nil
}

func getContainerCommand(pod *v1.Pod, namespace string, containerIndex int) ([]string, error) {
	// https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#notes
	var err error
	var img crv1.Image
	var cfg *crv1.ConfigFile
	args := []string{}
	c := pod.Spec.Containers[containerIndex]
	if len(c.Command) != 0 {
		args = append(args, c.Command...)
	} else {
		img, err = getRemoteImage(pod, namespace, containerIndex)
		if err != nil {
			return nil, err
		}
		cfg, err = img.ConfigFile()
		if err != nil {
			return nil, fmt.Errorf("Failed to get config for image %q: %v", c.Image, err)
		}
		if len(cfg.Config.Entrypoint) != 0 {
			args = append(args, cfg.Config.Entrypoint...)
		}
	}
	if len(c.Args) != 0 {
		args = append(args, c.Args...)
	} else {
		if cfg != nil && len(cfg.Config.Cmd) != 0 {
			args = append(args, cfg.Config.Cmd...)
		}
	}
	return args, nil
}

func getMountPath(mc common.MetricsCollectorSpec) (string, common.FileSystemKind) {
	if mc.Collector.Kind == common.StdOutCollector {
		return common.DefaultFilePath, common.FileKind
	} else if mc.Collector.Kind == common.FileCollector {
		return mc.Source.FileSystemPath.Path, common.FileKind
	} else if mc.Collector.Kind == common.TfEventCollector {
		return mc.Source.FileSystemPath.Path, common.DirectoryKind
	} else if mc.Collector.Kind == common.CustomCollector {
		if mc.Source == nil || mc.Source.FileSystemPath == nil {
			return "", common.InvalidKind
		}
		return mc.Source.FileSystemPath.Path, mc.Source.FileSystemPath.Kind
	} else {
		return "", common.InvalidKind
	}
}

func needWrapWorkerContainer(mc common.MetricsCollectorSpec) bool {
	mcKind := mc.Collector.Kind
	for _, kind := range NeedWrapWorkerMetricsCollectorList {
		if mcKind == kind {
			return true
		}
	}
	return false
}

func wrapWorkerContainer(trial *trialsv1beta1.Trial, pod *v1.Pod, namespace,
	metricsFile string, pathKind common.FileSystemKind) error {
	// Search for primary container.
	index := getPrimaryContainerIndex(pod.Spec.Containers, trial.Spec.PrimaryContainerName)
	if index >= 0 {
		command := []string{"sh", "-c"}
		args, err := getContainerCommand(pod, namespace, index)
		if err != nil {
			return err
		}
		// If the first two commands are sh -c, we do not inject command.
		if args[0] == "sh" || args[0] == "bash" {
			if args[1] == "-c" {
				command = args[0:2]
				args = args[2:]
			}
		}
		mc := trial.Spec.MetricsCollector
		if mc.Collector.Kind == common.StdOutCollector {
			redirectStr := fmt.Sprintf("1>%s 2>&1", metricsFile)
			args = append(args, redirectStr)
		}

		// Get metrics file directory
		metricsFileDir := metricsFile
		if pathKind == common.FileKind {
			metricsFileDir = filepath.Dir(metricsFile)
		}

		// If early stopping is set add appropriate command
		if trial.Spec.EarlyStoppingRules != nil {
			args = append(args, "||", getEarlyStoppingCommand(metricsFileDir, pathKind))
		}
		// Add completed command to run without early stopping
		args = append(args, "&&", getMarkCompletedCommand(metricsFileDir, pathKind))

		argsStr := strings.Join(args, " ")
		c := &pod.Spec.Containers[index]
		c.Command = command
		c.Args = []string{argsStr}
	} else {
		return fmt.Errorf("Unable to find primary container %v in mutated pod containers %v",
			trial.Spec.PrimaryContainerName, pod.Spec.Containers)
	}
	return nil
}

func getEarlyStoppingCommand(metricsFileDir string, pathKind common.FileSystemKind) string {

	// $$$$ is process id in shell
	// In condition: [ $(head -n $$.pid) ], process id can be received with $$
	pidFile := filepath.Join(metricsFileDir, "$$$$.pid")
	pidFileCondition := filepath.Join(metricsFileDir, "$$.pid")

	containerStopped := "Training Container was Early Stopped"
	containerFailed := "Training Container was Failed"

	// If main training process failed, it checks that $$.pid (e.g. 6.pid) file exists and contains "early-stopped" line.
	// Otherwise, training was failed and Job must be failed. In that case, we execute "exit 1" at the end
	return fmt.Sprintf(`if test -f %v && [ $(head -n 1 %v) = %v ]; then echo %v; else echo %v; exit 1; fi`,
		pidFile, pidFileCondition, mccommon.TrainingEarlyStopped, containerStopped, containerFailed)
}

func getMarkCompletedCommand(metricsFileDir string, pathKind common.FileSystemKind) string {
	// $$$$ is process id in shell
	pidFile := filepath.Join(metricsFileDir, "$$$$.pid")
	return fmt.Sprintf("echo %s > %s", mccommon.TrainingCompleted, pidFile)
}

func addContainerVolumeMount(c *v1.Container, vm *v1.VolumeMount) {
	if c.VolumeMounts == nil {
		c.VolumeMounts = make([]v1.VolumeMount, 0)
	}
	c.VolumeMounts = append(c.VolumeMounts, *vm)
}

func mutateMetricsCollectorVolume(pod *v1.Pod, mountPath, sidecarContainerName, primaryContainerName string, pathKind common.FileSystemKind) error {
	metricsVol := v1.Volume{
		Name: common.MetricsVolume,
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	}
	dir := mountPath
	if pathKind == common.FileKind {
		dir = filepath.Dir(mountPath)
	}
	vm := v1.VolumeMount{
		Name:      metricsVol.Name,
		MountPath: dir,
	}

	indexList := []int{}
	for i, c := range pod.Spec.Containers {
		shouldMount := false
		// We should mount volume only on sidecar and primary containers
		if c.Name == sidecarContainerName || c.Name == primaryContainerName {
			shouldMount = true
		}
		if shouldMount {
			indexList = append(indexList, i)
		}
	}

	for _, i := range indexList {
		addContainerVolumeMount(&pod.Spec.Containers[i], &vm)
	}
	pod.Spec.Volumes = append(pod.Spec.Volumes, metricsVol)

	return nil
}

func mutatePodMetadata(pod *v1.Pod, trial *trialsv1beta1.Trial) {
	podLabels := map[string]string{}

	// Get labels from the created pod.
	if pod.Labels != nil {
		podLabels = pod.Labels
	}

	// Get labels from Trial.
	for k, v := range trial.Labels {
		podLabels[k] = v
	}

	// Add Trial name label.
	podLabels[consts.LabelTrialName] = trial.GetName()

	// Append label to the Pod metadata.
	pod.Labels = podLabels
}

func getSidecarContainerName(cKind common.CollectorKind) string {
	if cKind == common.StdOutCollector || cKind == common.FileCollector {
		return mccommon.MetricLoggerCollectorContainerName
	} else {
		return mccommon.MetricCollectorContainerName
	}
}
