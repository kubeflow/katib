package consts

import "github.com/kubeflow/katib/pkg/util/v1alpha3/env"

const (
	// ConfigExperimentSuggestionName is the config name of the
	// suggestion client implementation in experiment controller.
	ConfigExperimentSuggestionName = "experiment-suggestion-name"
	// ConfigCertLocalFS is the config name which indicates if we
	// should store the cert in file system.
	ConfigCertLocalFS = "cert-local-filesystem"
	// ConfigInjectSecurityContext is the config name which indicates
	// if we should inject the security context into the metrics collector
	// sidecar.
	ConfigInjectSecurityContext = "inject-security-context"

	// LabelExperimentName is the label of experiment name.
	LabelExperimentName = "experiment"
	// LabelSuggestionName is the label of suggestion name.
	LabelSuggestionName = "suggestion"
	// LabelDeploymentName is the label of deployment name.
	LabelDeploymentName = "deployment"

	// ContainerSuggestion is the container name in Suggestion.
	ContainerSuggestion = "suggestion"

	// DefaultSuggestionPort is the default port of suggestion service.
	DefaultSuggestionPort = 6789
	// DefaultSuggestionPortName is the default port name of suggestion service.
	DefaultSuggestionPortName = "katib-api"
	// DefaultGRPCService is the default service name in Suggestion,
	// which is used to run healthz check using grpc probe.
	DefaultGRPCService = "manager.v1alpha3.Suggestion"

	// DefaultKatibNamespaceEnvName is the default env name of katib namespace
	DefaultKatibNamespaceEnvName = "KATIB_CORE_NAMESPACE"

	// KatibConfigMapName is the config map constants
	// Configmap name which includes Katib's configuration
	KatibConfigMapName = "katib-config"
	// LabelSuggestionTag is the name of suggestion config in configmap.
	LabelSuggestionTag = "suggestion"
	// LabelSuggestionImageTag is the name of suggestion image config in configmap.
	LabelSuggestionImageTag = "image"
	// LabelSuggestionCPULimitTag is the name of suggestion CPU Limit config in configmap.
	LabelSuggestionCPULimitTag = "cpuLimit"
	// DefaultCPULimit is the default value for CPU Limit
	DefaultCPULimit = "500m"
	// LabelSuggestionCPURequestTag is the name of suggestion CPU Request config in configmap.
	LabelSuggestionCPURequestTag = "cpuRequest"
	// DefaultCPURequest is the default value for CPU Request
	DefaultCPURequest = "50m"
	// LabelSuggestionMemLimitTag is the name of suggestion Mem Limit config in configmap.
	LabelSuggestionMemLimitTag = "memLimit"
	// DefaultMemLimit is the default value for mem Limit
	DefaultMemLimit = "100Mi"
	// LabelSuggestionMemRequestTag is the name of suggestion Mem Request config in configmap.
	LabelSuggestionMemRequestTag = "memRequest"
	// DefaultMemRequest is the default value for mem Request
	DefaultMemRequest = "10Mi"
	// LabelSuggestionDiskLimitTag is the name of suggestion Disk Limit config in configmap.
	LabelSuggestionDiskLimitTag = "diskLimit"
	// DefaultDiskLimit is the default value for disk limit.
	DefaultDiskLimit = "5Gi"
	// LabelSuggestionDiskRequestTag is the name of suggestion Disk Request config in configmap.
	LabelSuggestionDiskRequestTag = "diskRequest"
	// DefaultDiskRequest is the default value for disk request.
	DefaultDiskRequest = "500Mi"
	// LabelSuggestionImagePullPolicy is the name of suggestion image pull policy in configmap.
	LabelSuggestionImagePullPolicy = "imagePullPolicy"
	// DefaultImagePullPolicy is the default value for image pull policy.
	DefaultImagePullPolicy = "IfNotPresent"
	// LabelMetricsCollectorSidecar is the name of metrics collector config in configmap.
	LabelMetricsCollectorSidecar = "metrics-collector-sidecar"
	// LabelMetricsCollectorSidecarImage is the name of metrics collector image config in configmap.
	LabelMetricsCollectorSidecarImage = "image"
	// LabelMetricsCollectorCPULimitTag is the name of metrics collector CPU Limit config in configmap.
	LabelMetricsCollectorCPULimitTag = "cpuLimit"
	// LabelMetricsCollectorCPURequestTag is the name of metrics collector CPU Request config in configmap.
	LabelMetricsCollectorCPURequestTag = "cpuRequest"
	// LabelMetricsCollectorMemLimitTag is the name of metrics collector Mem Limit config in configmap.
	LabelMetricsCollectorMemLimitTag = "memLimit"
	// LabelMetricsCollectorMemRequestTag is the name of metrics collector Mem Request config in configmap.
	LabelMetricsCollectorMemRequestTag = "memRequest"
	// LabelMetricsCollectorDiskLimitTag is the name of metrics collector Disk Limit config in configmap.
	LabelMetricsCollectorDiskLimitTag = "diskLimit"
	// LabelMetricsCollectorDiskRequestTag is the name of metrics collector Disk Request config in configmap.
	LabelMetricsCollectorDiskRequestTag = "diskRequest"
	// LabelMetricsCollectorImagePullPolicy is the name of metrics collector image pull policy in configmap.
	LabelMetricsCollectorImagePullPolicy = "imagePullPolicy"

	// ReconcileErrorReason is the reason when there is a reconcile error.
	ReconcileErrorReason = "ReconcileError"

	// JobKindJob is the kind of the Kubernetes Job.
	JobKindJob = "Job"
	// JobKindTF is the kind of TFJob.
	JobKindTF = "TFJob"
	// JobKindPyTorch is the kind of PyTorchJob.
	JobKindPyTorch = "PyTorchJob"

	// JobVersionJob is the api version of Kubernetes Job.
	JobVersionJob = "v1"
	// JobVersionTF is the api version of TFJob.
	JobVersionTF = "v1"
	// JobVersionPyTorch is the api version of PyTorchJob.
	JobVersionPyTorch = "v1"

	// JobGroupJob is the group name of Kubernetes Job.
	JobGroupJob = "batch"
	// JobGroupKubeflow is the group name of Kubeflow.
	JobGroupKubeflow = "kubeflow.org"
)

var (
	// DefaultKatibNamespace is the default namespace of katib deployment.
	DefaultKatibNamespace = env.GetEnvOrDefault(DefaultKatibNamespaceEnvName, "kubeflow")
)
