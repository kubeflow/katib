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

	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibconfig"
)

var log = logf.Log.WithName("injector-webhook")

// sidecarInjector that inject metrics collect sidecar into master pod
type sidecarInjector struct {
	client  client.Client
	decoder types.Decoder

	// injectSecurityContext indicates if we should inject the security
	// context into the metrics collector sidecar.
	injectSecurityContext bool
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

// NewSidecarInjector returns a new sidecar injector.
func NewSidecarInjector(c client.Client) *sidecarInjector {
	return &sidecarInjector{
		injectSecurityContext: viper.GetBool(consts.ConfigInjectSecurityContext),
		client:                c,
	}
}

func (s *sidecarInjector) MutationRequired(pod *v1.Pod, ns string) (bool, error) {
	object, err := util.ConvertObjectToUnstructured(pod)
	if err != nil {
		return false, err
	}

	// Try to get Katib job kind and job name from mutating pod
	jobKind, jobName, err := s.getKatibJob(object, ns)
	if err != nil {
		return false, nil
	}

	trial := &trialsv1beta1.Trial{}
	// jobName and Trial name is equal
	if err := s.client.Get(context.TODO(), apitypes.NamespacedName{Name: jobName, Namespace: ns}, trial); err != nil {
		return false, err
	}

	// If PrimaryPodLabel is not set we mutate all pods which are related to Trial job
	// Otherwise mutate pod only with appropriate labels
	if trial.Spec.PrimaryPodLabels != nil {
		if !isPrimaryPod(pod.Labels, trial.Spec.PrimaryPodLabels) {
			return false, nil
		}
	} else {
		// TODO (andreyvelich): This can be deleted after switch to custom CRD
		if !isMasterRole(pod, jobKind) {
			return false, nil
		}
	}

	if trial.Spec.MetricsCollector.Collector.Kind == common.NoneCollector {
		return false, nil
	}
	return true, nil
}

func (s *sidecarInjector) Mutate(pod *v1.Pod, namespace string) (*v1.Pod, error) {
	mutatedPod := pod.DeepCopy()

	object, err := util.ConvertObjectToUnstructured(pod)
	if err != nil {
		return nil, err
	}

	// Try to get Katib job kind and job name from mutating pod
	jobKind, jobName, _ := s.getKatibJob(object, namespace)

	trial := &trialsv1beta1.Trial{}
	// jobName and Trial name is equal
	if err := s.client.Get(context.TODO(), apitypes.NamespacedName{Name: jobName, Namespace: namespace}, trial); err != nil {
		return nil, err
	}

	injectContainer, err := s.getMetricsCollectorContainer(trial, pod)
	if err != nil {
		return nil, err
	}
	mutatedPod.Spec.Containers = append(mutatedPod.Spec.Containers, *injectContainer)

	mutatedPod.Spec.ShareProcessNamespace = pointer.BoolPtr(true)

	mountPath, pathKind := getMountPath(trial.Spec.MetricsCollector)
	if mountPath != "" {
		if err = mutateVolume(mutatedPod, jobKind, mountPath, injectContainer.Name, trial.Spec.PrimaryContainerName, pathKind); err != nil {
			return nil, err
		}
	}
	if needWrapWorkerContainer(trial.Spec.MetricsCollector) {
		if err = wrapWorkerContainer(mutatedPod, namespace, jobKind, mountPath, pathKind, trial); err != nil {
			return nil, err
		}
	}

	// For Job kind mutated pod has only generate name
	if mutatedPod.Name != "" {
		log.Info("Inject metrics collector sidecar container", "Pod Name", mutatedPod.Name, "Trial", jobName)
	} else {
		log.Info("Inject metrics collector sidecar container", "Pod Generate Name", mutatedPod.GenerateName, "Trial", jobName)
	}
	return mutatedPod, nil
}

func (s *sidecarInjector) getMetricsCollectorContainer(trial *trialsv1beta1.Trial, originalPod *v1.Pod) (*v1.Container, error) {
	mc := trial.Spec.MetricsCollector
	if mc.Collector.Kind == common.CustomCollector {
		return mc.Collector.CustomCollector, nil
	}
	metricName := trial.Spec.Objective.ObjectiveMetricName
	for _, v := range trial.Spec.Objective.AdditionalMetricNames {
		metricName += ";"
		metricName += v
	}
	metricsCollectorConfigData, err := katibconfig.GetMetricsCollectorConfigData(mc.Collector.Kind, s.client)
	if err != nil {
		return nil, err
	}
	args := getMetricsCollectorArgs(trial.Name, metricName, mc)
	sidecarContainerName := getSidecarContainerName(trial.Spec.MetricsCollector.Collector.Kind)

	injectContainer := v1.Container{
		Name:            sidecarContainerName,
		Image:           metricsCollectorConfigData.Image,
		Args:            args,
		ImagePullPolicy: metricsCollectorConfigData.ImagePullPolicy,
		Resources:       metricsCollectorConfigData.Resource,
	}

	// Inject the security context when the flag is enabled.
	if s.injectSecurityContext {
		if len(originalPod.Spec.Containers) != 0 &&
			originalPod.Spec.Containers[0].SecurityContext != nil {
			injectContainer.SecurityContext = originalPod.Spec.Containers[0].SecurityContext.DeepCopy()
		}
	}

	return &injectContainer, nil
}

func (s *sidecarInjector) getKatibJob(object *unstructured.Unstructured, namespace string) (string, string, error) {
	owners := object.GetOwnerReferences()
	// jobKind and jobName points to the object kind and name that Trial is created
	jobKind := ""
	jobName := ""
	// Search for Trial owner in object owner references
	// Trial is owned object if kind = Trial kind and API version = Trial API version
	for _, owner := range owners {
		if owner.Kind == TrialKind && owner.APIVersion == TrialAPIVersion {
			jobKind = object.GetKind()
			jobName = object.GetName()
		}
	}
	// If Trial is not found in object owners search for nested owners
	if jobKind == "" {
		i := 0
		// Search for Trial ownership unless jobKind is empty and owners is exists
		for jobKind == "" && i < len(owners) {
			nestedJob := &unstructured.Unstructured{}
			// Get group and version from owner API version
			gv, err := schema.ParseGroupVersion(owners[i].APIVersion)
			if err != nil {
				return "", "", err
			}
			gvk := schema.GroupVersionKind{
				Group:   gv.Group,
				Version: gv.Version,
				Kind:    owners[i].Kind,
			}
			// Set GVK for nested unstructured object
			nestedJob.SetGroupVersionKind(gvk)
			// Get nested object from cluster.
			// Nested object namespace must be equal to object namespace
			err = s.client.Get(context.TODO(), apitypes.NamespacedName{Name: owners[i].Name, Namespace: namespace}, nestedJob)
			if err != nil {
				return "", "", err
			}
			// Recursively search for Trial ownership in nested object
			jobKind, jobName, err = s.getKatibJob(nestedJob, namespace)
			i++
		}
	}

	// If jobKind is empty after the loop, Trial doesn't own the object
	if jobKind == "" {
		return "", "", errors.New("The Pod doesn't belong to Katib Job")
	}

	return jobKind, jobName, nil
}
