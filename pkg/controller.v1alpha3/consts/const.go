package consts

import "github.com/kubeflow/katib/pkg/util/v1alpha3/env"

const (
	ConfigExperimentSuggestionName = "experiment-suggestion-name"

	LabelExperimentName = "experiment"
	LabelSuggestionName = "suggestion"

	ContainerSuggestion = "suggestion"

	DefaultSuggestionPort     = 6789
	DefaultSuggestionPortName = "katib-api"
	DefaultGRPCService        = "manager.v1alpha3.Suggestion"

	// Default env name of katib namespace
	DefaultKatibNamespaceEnvName = "KATIB_CORE_NAMESPACE"

	// Katib config map constants
	// Configmap name which includes Katib's configuration
	KatibConfigMapName = "katib-config"

	LabelSuggestionTag = "suggestion"

	LabelSuggestionImageTag = "image"

	ReconcileErrorReason = "ReconcileError"
)

var (
	DefaultKatibNamespace = env.GetEnvOrDefault(DefaultKatibNamespaceEnvName, "kubeflow")
)
