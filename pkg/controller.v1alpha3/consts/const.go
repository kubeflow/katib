package consts

import (
	"github.com/kubeflow/katib/pkg/util/v1alpha3/env"
	"time"
)

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
	// ConfigEnableGRPCProbeInSuggestion is the config name which indicates
	// if we should set GRPC probe in suggestion deployments.
	ConfigEnableGRPCProbeInSuggestion = "enable-grpc-probe-in-suggestion"

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
	// LabelSuggestionServiceAccountName is the name of suggestion service account in configmap.
	LabelSuggestionServiceAccountName = "serviceAccountName"
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

	// built-in JobRoles
	JobRole        = "job-role"
	JobRoleTF      = "tf-job-role"
	JobRolePyTorch = "pytorch-job-role"

	// AnnotationIstioSidecarInjectName is the annotation of Istio Sidecar
	AnnotationIstioSidecarInjectName = "sidecar.istio.io/inject"

	// AnnotationIstioSidecarInjectValue is the value of Istio Sidecar annotation
	AnnotationIstioSidecarInjectValue = "false"

	// LabelTrialTemplateConfigMapName is the label name for the Trial templates configMap
	LabelTrialTemplateConfigMapName = "app"
	// LabelTrialTemplateConfigMapValue is the label value for the Trial templates configMap
	LabelTrialTemplateConfigMapValue = "katib-trial-templates"
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
)
