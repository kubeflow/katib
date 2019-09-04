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
	"net/http"
	"os"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	apitypes "k8s.io/apimachinery/pkg/types"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	katibmanagerv1alpha3 "github.com/kubeflow/katib/pkg/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
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

	return true, nil
}

func (s *sidecarInjector) Mutate(pod *v1.Pod, namespace string) (*v1.Pod, error) {
	mutatedPod := pod.DeepCopy()

	kind, trialName, _ := getKabitJob(pod)
	trial := &trialsv1alpha3.Trial{}
	if err := s.client.Get(context.TODO(), apitypes.NamespacedName{Name: trialName, Namespace: namespace}, trial); err != nil {
		return nil, err
	}

	experimentName := "TODO_Experiment"
	for _, v := range trial.OwnerReferences {
		if v.Kind == "Experiment" {
			experimentName = v.Name
		}
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
	injectContainer := v1.Container{
		Name:            "sidecar-metrics-collector",
		Image:           image,
		Args:            []string{"-e", experimentName, "-t", trialName, "-k", kind, "-n", namespace, "-m", katibmanagerv1alpha3.GetManagerAddr(), "-mn", metricName},
		ImagePullPolicy: v1.PullIfNotPresent,
		VolumeMounts:    pod.Spec.Containers[0].VolumeMounts,
	}
	mutatedPod.Spec.Containers = append(mutatedPod.Spec.Containers, injectContainer)
	mutatedPod.Spec.ServiceAccountName = pod.Spec.ServiceAccountName

	log.Info("Inject metrics collector sidecar container", "Pod", pod.Name, "Trial", trialName, "Experiment", experimentName)
	return mutatedPod, nil
}

func (s *sidecarInjector) getMetricsCollectorImage(cKind common.CollectorKind) (string, error) {
	configMap := &v1.ConfigMap{}
	err := s.client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: experimentsv1alpha3.KatibConfigMapName, Namespace: os.Getenv(experimentsv1alpha3.DefaultKatibNamespaceEnvName)},
		configMap)
	if err != nil {
		log.Error(err, "Failed to find config map", "name", experimentsv1alpha3.KatibConfigMapName)
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
			return "",  errors.New("Cannot support metricsCollector injection for kind " + kind)
		}
	} else {
		return "",  errors.New("Failed to find metrics collector configuration in configmap " + experimentsv1alpha3.KatibConfigMapName)
	}
}
