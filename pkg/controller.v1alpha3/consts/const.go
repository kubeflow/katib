package consts

import "github.com/kubeflow/katib/pkg/util/v1alpha3/env"

const (
	// ConfigExperimentSuggestionName is the config name of the
	// suggestion client implementation in experiment controller.
	ConfigExperimentSuggestionName = "experiment-suggestion-name"

	// LabelExperimentName is the label of experiment name.
	LabelExperimentName = "experiment"
	// LabelSuggestionName is the label of suggestion name.
	LabelSuggestionName = "suggestion"

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
	DefaultKatibNamespace = env.GetEnvOrDefault(DefaultKatibNamespaceEnvName, "kubeflow")
)
