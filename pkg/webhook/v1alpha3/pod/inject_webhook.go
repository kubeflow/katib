/*
Copyright 2019 The Kubernetes Authors.

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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	katibmanagerv1alpha3 "github.com/kubeflow/katib/pkg/common/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	mccommon "github.com/kubeflow/katib/pkg/metricscollector/v1alpha3/common"
)

var log = logf.Log.WithName("injector-webhook")

// sidecarInjector that inject metrics collect sidecar into master pod
type sidecarInjector struct {
	client         client.Client
	decoder        types.Decoder
	managerService string
}

var _ admission.Handler = &sidecarInjector{}

func (s *sidecarInjector) Handle(ctx context.Context, req types.Request) types.Response {
	// Get the namespace from req since the namespace in the pod is empty.
	namespace := req.AdmissionRequest.Namespace
	pod := &v1.Pod{}
	err := s.decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	// Check whether the pod need to be mutated
	needMutate, err := s.MutationRequired(pod, namespace)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	} else {
		if !needMutate {
			return admission.ValidationResponse(true, "")
		}
	}

	// Do mutation
	mutatedPod, err := s.Mutate(pod, namespace)
	if err != nil {
		log.Error(err, "Failed to inject metrics collector")
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	return admission.PatchResponse(pod, mutatedPod)
}

var _ inject.Client = &sidecarInjector{}

func (s *sidecarInjector) InjectClient(c client.Client) error {
	s.client = c
	return nil
}

var _ inject.Decoder = &sidecarInjector{}

func (s *sidecarInjector) InjectDecoder(d types.Decoder) error {
	s.decoder = d
	return nil
}

func NewSidecarInjector(c client.Client, ms string) *sidecarInjector {
	return &sidecarInjector{
		client:         c,
		managerService: ms,
	}
}

func (s *sidecarInjector) MutationRequired(pod *v1.Pod, ns string) (bool, error) {
	jobKind, jobName, err := getKabitJob(pod)
	if err != nil {
		return false, nil
	}
	if !isMasterRole(pod, jobKind) {
		return false, nil
	}

	trialName := jobName
	trial := &trialsv1alpha3.Trial{}
	err = s.client.Get(context.TODO(), apitypes.NamespacedName{Name: trialName, Namespace: ns}, trial)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if trial.Spec.MetricsCollector.Collector.Kind == common.NoneCollector {
		return false, nil
	}
	return true, nil
}

func (s *sidecarInjector) Mutate(pod *v1.Pod, namespace string) (*v1.Pod, error) {
	mutatedPod := pod.DeepCopy()

	kind, trialName, _ := getKabitJob(pod)
	trial := &trialsv1alpha3.Trial{}
	if err := s.client.Get(context.TODO(), apitypes.NamespacedName{Name: trialName, Namespace: namespace}, trial); err != nil {
		return nil, err
	}

	metricName := trial.Spec.Objective.ObjectiveMetricName
	for _, v := range trial.Spec.Objective.AdditionalMetricNames {
		metricName += ";"
		metricName += v
	}

	image, err := s.getMetricsCollectorImage(trial.Spec.MetricsCollector.Collector.Kind)
	if err != nil {
		return nil, err
	}
	args := getMetricsCollectorArgs(trialName, metricName, trial.Spec.MetricsCollector)
	injectContainer := v1.Container{
		Name:            mccommon.MetricCollectorContainerName,
		Image:           image,
		Args:            args,
		ImagePullPolicy: v1.PullIfNotPresent,
	}
	mutatedPod.Spec.Containers = append(mutatedPod.Spec.Containers, injectContainer)
	mutatedPod.Spec.ServiceAccountName = pod.Spec.ServiceAccountName
	mutatedPod.Spec.ShareProcessNamespace = pointer.BoolPtr(true)

	if mountFile := getMountFile(trial.Spec.MetricsCollector); mountFile != "" {
		wrapWorkerContainer(mutatedPod, kind, mountFile, trial.Spec.MetricsCollector)
		if err = mutateVolume(mutatedPod, kind, mountFile); err != nil {
			return nil, err
		}
	}

	log.Info("Inject metrics collector sidecar container", "Pod", pod.Name, "Trial", trialName)

	return mutatedPod, nil
}

func (s *sidecarInjector) getMetricsCollectorImage(cKind common.CollectorKind) (string, error) {
	configMap := &v1.ConfigMap{}
	err := s.client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		log.Error(err, "Failed to find config map", "name", consts.KatibConfigMapName)
		// Error reading the object - requeue the request.
		return "", err
	}
	if mcs, ok := configMap.Data[MetricsCollectorSidecar]; ok {
		kind := string(cKind)
		mcsConfig := map[string]map[string]string{}
		if err := json.Unmarshal([]byte(mcs), &mcsConfig); err != nil {
			return "", err
		}
		if mc, ok := mcsConfig[kind]; ok {
			if image, yes := mc[MetricsCollectorSidecarImage]; yes {
				return image, nil
			} else {
				return "", errors.New("Failed to find " + MetricsCollectorSidecarImage + " configuration for metricsCollector kind " + kind)
			}
		} else {
			return "", errors.New("Cannot support metricsCollector injection for kind " + kind)
		}
	} else {
		return "", errors.New("Failed to find metrics collector configuration in configmap " + consts.KatibConfigMapName)
	}
}

func getMetricsCollectorArgs(trialName, metricName string, mc trialsv1alpha3.MetricsCollectorSpec) []string {
	args := []string{"-t", trialName, "-m", katibmanagerv1alpha3.GetManagerAddr(), "-mn", metricName}
	if mountFile := getMountFile(mc); mountFile != "" {
		args = append(args, "-f", mountFile)
	}
	return args
}

func getMountFile(mc trialsv1alpha3.MetricsCollectorSpec) string {
	if mc.Collector.Kind == common.StdOutCollector {
		return common.DefaultFilePath
	} else if mc.Collector.Kind == common.FileCollector {
		return mc.Source.FileSystemPath.Path
	} else {
		return ""
	}
}

func wrapWorkerContainer(pod *v1.Pod, jobKind, metricsFile string, mc trialsv1alpha3.MetricsCollectorSpec) {
	if mc.Collector.Kind != common.StdOutCollector {
		return
	}
	index := -1
	for i, c := range pod.Spec.Containers {
		if isWorkerContainer(jobKind, i, c) {
			index = i
			break
		}
	}
	if index >= 0 {
		// TODO(hougangliu): handle container.command is nil case
		c := &pod.Spec.Containers[index]
		command := []string{"sh", "-c"}
		args := []string{}
		if c.Command != nil {
			args = append(args, c.Command...)
		}
		if c.Args != nil {
			args = append(args, c.Args...)
		}
		redirectStr := fmt.Sprintf("1>%s 2>&1", metricsFile)
		args = append(args, redirectStr)
		argsStr := strings.Join(args, " ")
		c.Command = command
		c.Args = []string{argsStr}
	}
}

func isWorkerContainer(jobKind string, index int, c v1.Container) bool {
	switch jobKind {
	case BatchJob:
		if index == 0 {
			// for Job worker, the first container will be taken as worker container,
			// katib document should note it
			return true
		}
	case TFJob:
		if c.Name == TFJobWorkerContainerName {
			return true
		}
	case PyTorchJob:
		if c.Name == PyTorchJobWorkerContainerName {
			return true
		}
	default:
		log.Info("Invalid Katib worker kind", "JobKind", jobKind)
		return false
	}
	return false
}

func mutateVolume(pod *v1.Pod, jobKind, mountFile string) error {
	metricsVol := v1.Volume{
		Name: common.MetricsVolume,
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	}
	vm := v1.VolumeMount{
		Name:      metricsVol.Name,
		MountPath: filepath.Dir(mountFile),
	}
	index_list := []int{}
	for i, c := range pod.Spec.Containers {
		shouldMount := false
		if c.Name == mccommon.MetricCollectorContainerName {
			shouldMount = true
		} else {
			shouldMount = isWorkerContainer(jobKind, i, c)
		}
		if shouldMount {
			index_list = append(index_list, i)
		}
	}
	for _, i := range index_list {
		c := &pod.Spec.Containers[i]
		if c.VolumeMounts == nil {
			c.VolumeMounts = make([]v1.VolumeMount, 0)
		}
		c.VolumeMounts = append(c.VolumeMounts, vm)
		pod.Spec.Containers[i] = *c
	}
	pod.Spec.Volumes = append(pod.Spec.Volumes, metricsVol)

	return nil
}
