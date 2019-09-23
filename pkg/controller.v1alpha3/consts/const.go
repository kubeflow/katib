package consts

import (
	"os"
)

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
	DefaultKatibNamespace = getEnvOrDefault(DefaultKatibNamespaceEnvName, "kubeflow")
)

func getEnvOrDefault(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
