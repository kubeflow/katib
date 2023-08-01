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
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	katibmanagerv1beta1 "github.com/kubeflow/katib/pkg/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibconfig"
)

var log = logf.Log.WithName("injector-webhook")

// SidecarInjector injects metrics collect sidecar to the primary pod.
type SidecarInjector struct {
	client  client.Client
	decoder *admission.Decoder

	// injectSecurityContext indicates if we should inject the security
	// context into the metrics collector sidecar.
	injectSecurityContext bool
}

// NewSidecarInjector returns a new sidecar injector with the given client.
func NewSidecarInjector(c client.Client, d *admission.Decoder) *SidecarInjector {
	return &SidecarInjector{
		injectSecurityContext: viper.GetBool(consts.ConfigInjectSecurityContext),
		client:                c,
		decoder:               d,
	}
}

func (s *SidecarInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	// Get the namespace from req since the namespace in the pod is empty.
	namespace := req.AdmissionRequest.Namespace
	pod := &v1.Pod{}
	err := s.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// Check whether the pod need to be mutated
	needMutate, err := s.MutationRequired(pod, namespace)
	if err != nil {
		log.Info("Unable to run MutationRequired", "Error", err)
		return admission.Errored(http.StatusInternalServerError, err)
	} else if !needMutate {
		return admission.ValidationResponse(true, "")
	}

	// Do mutation
	mutatedPod, err := s.Mutate(pod, namespace)
	if err != nil {
		log.Error(err, "Failed to mutate Trial's pod")
		return admission.Errored(http.StatusBadRequest, err)
	}

	marshaledPod, err := json.Marshal(mutatedPod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.AdmissionRequest.Object.Raw, marshaledPod)
}

func (s *SidecarInjector) MutationRequired(pod *v1.Pod, ns string) (bool, error) {
	object, err := util.ConvertObjectToUnstructured(pod)
	if err != nil {
		return false, err
	}

	// Try to get Katib Job name from mutating pod
	_, jobName, err := s.getKatibJob(object, ns)
	if err != nil {
		return false, nil
	}

	trial := &trialsv1beta1.Trial{}
	// Job name and Trial name is equal
	if err := s.client.Get(context.TODO(), apitypes.NamespacedName{Name: jobName, Namespace: ns}, trial); err != nil {
		return false, err
	}

	return true, nil
}

func (s *SidecarInjector) Mutate(pod *v1.Pod, namespace string) (*v1.Pod, error) {
	mutatedPod := pod.DeepCopy()

	object, err := util.ConvertObjectToUnstructured(pod)
	if err != nil {
		return nil, err
	}

	// Try to get Katib job kind and job name from mutating pod
	_, jobName, _ := s.getKatibJob(object, namespace)

	trial := &trialsv1beta1.Trial{}
	// jobName and Trial name is equal
	if err := s.client.Get(context.TODO(), apitypes.NamespacedName{Name: jobName, Namespace: namespace}, trial); err != nil {
		return nil, err
	}

	// Add Katib Trial labels to the Pod metadata.
	mutatePodMetadata(mutatedPod, trial)

	// Do the following mutation only for the Primary pod.
	// If PrimaryPodLabel is not set we mutate all pods which are related to Trial job.
	// Otherwise, mutate pod only with the appropriate labels.
	if trial.Spec.PrimaryPodLabels != nil && !isPrimaryPod(pod.Labels, trial.Spec.PrimaryPodLabels) {
		return mutatedPod, nil
	}

	// If Metrics Collector in None, skip the mutation.
	if trial.Spec.MetricsCollector.Collector.Kind == common.NoneCollector {
		return mutatedPod, nil
	}

	// Create metrics sidecar container spec
	injectContainer, err := s.getMetricsCollectorContainer(trial, pod)
	if err != nil {
		return nil, err
	}
	mutatedPod.Spec.Containers = append(mutatedPod.Spec.Containers, *injectContainer)

	// Enable shared volume between suggestion <> trial
	if err = s.mutateSuggestionVolume(mutatedPod, injectContainer.Name, trial); err != nil {
		return nil, err
	}

	isShareProcessNamespace := true
	mutatedPod.Spec.ShareProcessNamespace = &isShareProcessNamespace

	mountPath, pathKind := getMountPath(trial.Spec.MetricsCollector)
	if mountPath != "" {
		if err = mutateMetricsCollectorVolume(mutatedPod, mountPath, injectContainer.Name, trial.Spec.PrimaryContainerName, pathKind); err != nil {
			return nil, err
		}
	}
	if needWrapWorkerContainer(trial.Spec.MetricsCollector) {
		if err = wrapWorkerContainer(trial, mutatedPod, namespace, mountPath, pathKind); err != nil {
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

func (s *SidecarInjector) getMetricsCollectorContainer(trial *trialsv1beta1.Trial, originalPod *v1.Pod) (*v1.Container, error) {
	mc := trial.Spec.MetricsCollector
	if mc.Collector.Kind == common.CustomCollector {
		return mc.Collector.CustomCollector, nil
	}
	metricNames := trial.Spec.Objective.ObjectiveMetricName
	for _, v := range trial.Spec.Objective.AdditionalMetricNames {
		metricNames += ";"
		metricNames += v
	}

	// Convert rules to flag value with name;value;comparison;startStep order, e.g. accuracy;0.8;less;4.
	// If start step is empty, we apply rule from the first recorded metrics and flag is equal to accuracy;0.8;less;0.
	earlyStoppingRules := []string{}
	for _, rule := range trial.Spec.EarlyStoppingRules {
		newRule := rule.Name + ";" + rule.Value + ";" + string(rule.Comparison) + ";" + strconv.Itoa(rule.StartStep)
		earlyStoppingRules = append(earlyStoppingRules, newRule)
	}
	metricsCollectorConfigData, err := katibconfig.GetMetricsCollectorConfigData(mc.Collector.Kind, s.client)
	if err != nil {
		return nil, err
	}

	args, err := s.getMetricsCollectorArgs(trial, metricNames, mc, metricsCollectorConfigData, earlyStoppingRules)
	if err != nil {
		return nil, err
	}

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

func (s *SidecarInjector) getKatibJob(object *unstructured.Unstructured, namespace string) (string, string, error) {
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
			if err != nil && jobKind != "" {
				return "", "", err
			}
			i++
		}
	}

	// If jobKind is empty after the loop, Trial doesn't own the object
	if jobKind == "" {
		return "", "", errors.New("The Pod doesn't belong to Katib Job")
	}

	return jobKind, jobName, nil
}

func (s *SidecarInjector) getMetricsCollectorArgs(trial *trialsv1beta1.Trial, metricNames string, mc common.MetricsCollectorSpec, metricsCollectorConfigData configv1beta1.MetricsCollectorConfig, esRules []string) ([]string, error) {
	args := []string{"-t", trial.Name, "-m", metricNames, "-o-type", string(trial.Spec.Objective.Type), "-s-db", katibmanagerv1beta1.GetDBManagerAddr()}
	if mountPath, _ := getMountPath(mc); mountPath != "" {
		args = append(args, "-path", mountPath)
	}
	if mc.Source != nil && mc.Source.Filter != nil && len(mc.Source.Filter.MetricsFormat) > 0 {
		args = append(args, "-f", strings.Join(mc.Source.Filter.MetricsFormat, ";"))
	}
	if mc.Collector.Kind == common.FileCollector && mc.Source != nil {
		if mc.Source.FileSystemPath != nil {
			args = append(args, "-format", string(mc.Source.FileSystemPath.Format))
		}
	}
	if mc.Collector.Kind == common.StdOutCollector {
		args = append(args, "-format", string(common.TextFormat))
	}
	if metricsCollectorConfigData.WaitAllProcesses != nil {
		args = append(args, "-w", strconv.FormatBool(*metricsCollectorConfigData.WaitAllProcesses))
	}
	// Add stop rules and service endpoint for Early Stopping
	if len(esRules) > 0 {
		for _, rule := range esRules {
			args = append(args, "-stop-rule", rule)
		}
		// Suggestion name == Experiment name
		// Suggestion namespace == Trial namespace
		// Get suggestion to set early stopping service endpoint
		suggestionName := trial.ObjectMeta.Labels[consts.LabelExperimentName]
		suggestion := &suggestionsv1beta1.Suggestion{}
		err := s.client.Get(context.TODO(), apitypes.NamespacedName{Name: suggestionName, Namespace: trial.Namespace}, suggestion)
		if err != nil {
			return nil, err
		}
		args = append(args, "-s-earlystop", util.GetEarlyStoppingEndpoint(suggestion))
	}

	return args, nil
}

// Mutate trial container with shared Suggestions PVC when algorithm settings contains suggestion_trial_dir
func (s *SidecarInjector) mutateSuggestionVolume(pod *v1.Pod, primaryContainerName string, trial *trialsv1beta1.Trial) error {
	// Suggestion name == Experiment name
	// Suggestion namespace == Trial namespace
	experimentName := trial.ObjectMeta.Labels[consts.LabelExperimentName]
	experiment := &experimentsv1beta1.Experiment{}
	if err := s.client.Get(context.TODO(), apitypes.NamespacedName{Name: experimentName, Namespace: trial.Namespace}, experiment); err != nil {
		return err
	}
	suggestion := &suggestionsv1beta1.Suggestion{}
	if err := s.client.Get(context.TODO(), apitypes.NamespacedName{Name: experimentName, Namespace: trial.Namespace}, suggestion); err != nil {
		return err
	}

	// Check if mutation is needed
	checkpointPath := ""
	for _, s := range suggestion.Spec.Algorithm.AlgorithmSettings {
		if s.Name == consts.SuggestionVolumeMountKey && s.Value != "" {
			checkpointPath = s.Value
			break
		}
	}
	if checkpointPath == "" {
		return nil
	}

	// Generate folder name in format: <ExperimentName>/<TrialName>
	checkpointFolder := filepath.Join(experimentName, trial.Name)

	// Suggestion volume for the trial to the MetricsCollector
	suggestionVolume := v1.Volume{
		Name: consts.ContainerSuggestionVolumeName,
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
				ClaimName: util.GetSuggestionPersistentVolumeClaimName(suggestion),
			},
		},
	}

	vm := v1.VolumeMount{
		Name:      suggestionVolume.Name,
		MountPath: checkpointPath,
		SubPath:   checkpointFolder,
	}

	primaryContainerIndex := getPrimaryContainerIndex(pod.Spec.Containers, trial.Spec.PrimaryContainerName)
	addContainerVolumeMount(&pod.Spec.Containers[primaryContainerIndex], &vm)

	pod.Spec.Volumes = append(pod.Spec.Volumes, suggestionVolume)

	return nil
}
