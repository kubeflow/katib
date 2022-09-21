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

package consts

import (
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/kubeflow/katib/pkg/util/v1beta1/env"
)

const (

	// PluralTrial is the plural for Trial object
	PluralTrial = "trials"

	// ConfigExperimentSuggestionName is the config name of the
	// suggestion client implementation in experiment controller.
	ConfigExperimentSuggestionName = "experiment-suggestion-name"

	// ConfigCertLocalFS is the config name which indicates if we
	// should store the cert in file system.
	// TODO (andreyvelich): Currently is is not possible to store webhook cert in the local file system
	// ConfigCertLocalFS = "cert-local-filesystem"

	// ConfigInjectSecurityContext is the config name which indicates
	// if we should inject the security context into the metrics collector
	// sidecar.
	ConfigInjectSecurityContext = "inject-security-context"
	// ConfigEnableGRPCProbeInSuggestion is the config name which indicates
	// if we should set GRPC probe in suggestion deployments.
	ConfigEnableGRPCProbeInSuggestion = "enable-grpc-probe-in-suggestion"
	// ConfigTrialResources is the config name which indicates
	// resources list which can be used as trial template
	ConfigTrialResources = "trial-resources"

	// LabelExperimentName is the label of experiment name.
	LabelExperimentName = "katib.kubeflow.org/experiment"
	// LabelSuggestionName is the label of suggestion name.
	LabelSuggestionName = "katib.kubeflow.org/suggestion"
	// LabelDeploymentName is the label of deployment name.
	LabelDeploymentName = "katib.kubeflow.org/deployment"

	// ContainerSuggestion is the container name to run Suggestion service.
	ContainerSuggestion = "suggestion"
	// ContainerEarlyStopping is the container name to run EarlyStopping service.
	ContainerEarlyStopping = "early-stopping"
	// ContainerSuggestionVolumeName is the volume name that mounted on suggestion container
	ContainerSuggestionVolumeName = "suggestion-volume"

	// DefaultSuggestionPortName is the default port name of Suggestion service.
	DefaultSuggestionPortName = "suggestion-api"
	// DefaultSuggestionPort is the default port of Suggestion service.
	DefaultSuggestionPort = 6789
	// DefaultEarlyStoppingPortName is the default port name of EarlyStopping service.
	DefaultEarlyStoppingPortName = "earlystop-api"
	// DefaultEarlyStoppingPort is the default port of EarlyStopping service.
	DefaultEarlyStoppingPort = 6788

	// DefaultGRPCService is the default suggestion service name,
	// which is used to run healthz check using grpc probe.
	DefaultGRPCService = "manager.v1beta1.Suggestion"

	// DefaultGRPCRetryAttempts is the the maximum number of retries for gRPC calls
	DefaultGRPCRetryAttempts = 10
	// DefaultGRPCRetryPeriod is a fixed period of time between gRPC call retries
	DefaultGRPCRetryPeriod = 3 * time.Second

	// DefaultKatibNamespaceEnvName is the default env name of katib namespace
	DefaultKatibNamespaceEnvName = "KATIB_CORE_NAMESPACE"
	// DefaultKatibComposerEnvName is the default env name of katib suggestion composer
	DefaultKatibComposerEnvName = "KATIB_SUGGESTION_COMPOSER"

	// DefaultKatibDBManagerServiceNamespaceEnvName is the env name of Katib DB Manager namespace
	DefaultKatibDBManagerServiceNamespaceEnvName = "KATIB_DB_MANAGER_SERVICE_NAMESPACE"
	// DefaultKatibDBManagerServiceIPEnvName is the env name of Katib DB Manager IP
	DefaultKatibDBManagerServiceIPEnvName = "KATIB_DB_MANAGER_SERVICE_IP"
	// DefaultKatibDBManagerServicePortEnvName is the env name of Katib DB Manager Port
	DefaultKatibDBManagerServicePortEnvName = "KATIB_DB_MANAGER_SERVICE_PORT"

	// KatibConfigMapName is the configmap name which includes Katib's configuration.
	KatibConfigMapName = "katib-config"
	// LabelSuggestionTag is the name of suggestion config in Katib configmap.
	LabelSuggestionTag = "suggestion"
	// LabelMetricsCollectorSidecar is the name of metrics collector config in Katib configmap.
	LabelMetricsCollectorSidecar = "metrics-collector-sidecar"
	// LabelEarlyStoppingTag is the name of early stopping config in Katib configmap.
	LabelEarlyStoppingTag = "early-stopping"
	// DefaultImagePullPolicy is the default value for image pull policy.
	DefaultImagePullPolicy = corev1.PullIfNotPresent
	// DefaultCPULimit is the default value for CPU limit.
	DefaultCPULimit = "500m"
	// DefaultCPURequest is the default value for CPU request.
	DefaultCPURequest = "50m"
	// DefaultMemLimit is the default value for memory limit.
	DefaultMemLimit = "100Mi"
	// DefaultMemRequest is the default value for memory request.
	DefaultMemRequest = "10Mi"
	// DefaultDiskLimit is the default value for disk limit.
	DefaultDiskLimit = "5Gi"
	// DefaultDiskRequest is the default value for disk request.
	DefaultDiskRequest = "500Mi"

	// DefaultContainerSuggestionVolumeMountPath is the default mount path in suggestion container
	DefaultContainerSuggestionVolumeMountPath = "/opt/katib/data"

	// DefaultSuggestionVolumeStorage is the default value for suggestion's volume storage
	DefaultSuggestionVolumeStorage = "1Gi"

	// DefaultSuggestionVolumeAccessMode is the default value for suggestion's volume access mode
	DefaultSuggestionVolumeAccessMode = corev1.ReadWriteOnce

	// SuggestionVolumeMountKey specifies the AlgorithmSettings key used to toggle Suggestion managed trial storage
	SuggestionVolumeMountKey = "suggestion_trial_dir"

	// ReconcileErrorReason is the reason when there is a reconcile error.
	ReconcileErrorReason = "ReconcileError"

	// JobKindJob is the kind of the Kubernetes Job.
	JobKindJob = "Job"

	// AnnotationIstioSidecarInjectName is the annotation of Istio Sidecar
	AnnotationIstioSidecarInjectName = "sidecar.istio.io/inject"

	// AnnotationIstioSidecarInjectValue is the value of Istio Sidecar annotation
	AnnotationIstioSidecarInjectValue = "false"

	// LabelTrialTemplateConfigMapName is the label name for the Trial templates configMap
	LabelTrialTemplateConfigMapName = "katib.kubeflow.org/component"
	// LabelTrialTemplateConfigMapValue is the label value for the Trial templates configMap
	LabelTrialTemplateConfigMapValue = "trial-templates"

	// TrialTemplateParamReplaceFormat is the format to make substitution in Trial template from Names in TrialParameters
	// E.g if Name = learningRate, according value in Trial template must be ${trialParameters.learningRate}
	TrialTemplateParamReplaceFormat = "${trialParameters.%v}"

	// TrialTemplateParamReplaceFormatRegex is the regex for TrialParameters format in Trial template
	TrialTemplateParamReplaceFormatRegex = "\\$\\{trialParameters\\..+?\\}"

	// TrialTemplateMetaReplaceFormatRegex is the regex for TrialMetadata format in Trial template
	TrialTemplateMetaReplaceFormatRegex = "\\$\\{trialSpec\\.(.+?)\\}"
	// TrialTemplateMetaParseFormatRegex is the regex to parse the index of Annotations and Labels from meta key
	TrialTemplateMetaParseFormatRegex = "(.+)\\[(.+)]"

	// valid keys of trial metadata which are used to make substitution in Trial template
	TrialTemplateMetaKeyOfName        = "Name"
	TrialTemplateMetaKeyOfNamespace   = "Namespace"
	TrialTemplateMetaKeyOfKind        = "Kind"
	TrialTemplateMetaKeyOfAPIVersion  = "APIVersion"
	TrialTemplateMetaKeyOfAnnotations = "Annotations"
	TrialTemplateMetaKeyOfLabels      = "Labels"

	// UnavailableMetricValue is the value when metric was not reported or metric value can't be converted to float64
	// This value is recorded in to DB when metrics collector can't parse objective metric from the training logs.
	UnavailableMetricValue = "unavailable"
)

var (
	// DefaultKatibNamespace is the default namespace of katib deployment.
	DefaultKatibNamespace = env.GetEnvOrDefault(DefaultKatibNamespaceEnvName, "kubeflow")
	// DefaultComposer is the default composer of katib suggestion.
	DefaultComposer = env.GetEnvOrDefault(DefaultKatibComposerEnvName, "General")

	// DefaultKatibDBManagerServiceNamespace is the default namespace of Katib DB Manager
	DefaultKatibDBManagerServiceNamespace = env.GetEnvOrDefault(DefaultKatibDBManagerServiceNamespaceEnvName, DefaultKatibNamespace)
	// DefaultKatibDBManagerServiceIP is the default IP of Katib DB Manager
	DefaultKatibDBManagerServiceIP = env.GetEnvOrDefault(DefaultKatibDBManagerServiceIPEnvName, "katib-db-manager")
	// DefaultKatibDBManagerServicePort is the default Port of Katib DB Manager
	DefaultKatibDBManagerServicePort = env.GetEnvOrDefault(DefaultKatibDBManagerServicePortEnvName, "6789")

	// List of all valid keys of trial metadata for substitution in Trial template
	TrialTemplateMetaKeys = []string{
		TrialTemplateMetaKeyOfName,
		TrialTemplateMetaKeyOfNamespace,
		TrialTemplateMetaKeyOfKind,
		TrialTemplateMetaKeyOfAPIVersion,
		TrialTemplateMetaKeyOfAnnotations,
		TrialTemplateMetaKeyOfLabels,
	}
)
