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
	"errors"
	"net/http"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	v1 "k8s.io/api/core/v1"
	apitypes "k8s.io/apimachinery/pkg/types"

	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	katibmanagerv1alpha2 "github.com/kubeflow/katib/pkg/common/v1alpha2"
)

const (
	// JobNameLabel represents the label key for the job name, the value is job name
	JobNameLabel = "job-name"
	// JobRoleLabel represents the label key for the job role, e.g. the value is master
	JobRoleLabel = "job-role"
)

// For debug
var log = logf.Log.WithName("injector-webhook")

// sidecarInjector that inject metrics collect sidecar into master pod
type sidecarInjector struct {
	client         client.Client
	decoder        types.Decoder
	managerService string
}

var _ admission.Handler = &sidecarInjector{}

func (s *sidecarInjector) Handle(ctx context.Context, req types.Request) types.Response {
	pod := &v1.Pod{}
	err := s.decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	// Check whether the pod need to be mutated
	if !s.MutationRequired(pod) {
		return admission.ValidationResponse(true, "")
	}

	// Get the namespace from req since the namespace in the pod is empty.
	namespace := req.AdmissionRequest.Namespace
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

func (s *sidecarInjector) MutationRequired(pod *v1.Pod) bool {
	value, err := s.GetLabel(pod, JobRoleLabel)
	if err != nil || value != "master" {
		return false
	}
	return true
}

func (s *sidecarInjector) GetLabel(pod *v1.Pod, targetLabel string) (string, error) {
	labels := pod.Labels
	for k, v := range labels {
		if k == targetLabel {
			return v, nil
		}
	}
	return "", errors.New("Label " + targetLabel + " not found.")
}

func (s *sidecarInjector) Mutate(pod *v1.Pod, namespace string) (*v1.Pod, error) {
	mutatedPod := pod.DeepCopy()

	// Get the trial info from client
	trialName, err := s.GetLabel(pod, JobNameLabel)
	log.Info("FOR TEST: trialName: " + trialName + ", NameSpace: " + namespace)
	if err != nil {
		log.Info(err.Error())
		return nil, err
	}
	trial := &trialsv1alpha2.Trial{}
	err = s.client.Get(context.TODO(), apitypes.NamespacedName{Name: trialName, Namespace: namespace}, trial)
	if err != nil {
		log.Info(err.Error())
		return nil, err
	}

	metricName := trial.Spec.Objective.ObjectiveMetricName
	for _, v := range trial.Spec.Objective.AdditionalMetricNames {
		metricName += ";"
		metricName += v
	}

	// Hard code container, inject metrics collector
	injectContainer := v1.Container{
		Name:            "metrics-collector",
		Image:           "gcr.io/kubeflow-images-public/katib/v1alpha2/metrics-collector",
		Command:         []string{"./metricscollector"},
		Args:            []string{"-e", "TODO_Experiment", "-t", trialName, "-k", "TFJob", "-n", namespace, "-m", katibmanagerv1alpha2.GetManagerAddr(), "-mn", metricName},
		ImagePullPolicy: v1.PullIfNotPresent,
		VolumeMounts:    pod.Spec.Containers[0].VolumeMounts,
	}
	mutatedPod.Spec.Containers = append(mutatedPod.Spec.Containers, injectContainer)

	log.Info("-t " + trialName + " -k " + "TFJob" + " -n " + namespace + " -m " + katibmanagerv1alpha2.GetManagerAddr() + " -mn " + metricName)

	mutatedPod.Spec.ServiceAccountName = pod.Spec.ServiceAccountName

	return mutatedPod, nil
}
