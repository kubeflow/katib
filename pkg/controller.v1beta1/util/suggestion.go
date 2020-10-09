package util

import (
	"fmt"

	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

// GetSuggestionDeploymentName returns name for the suggestion's deployment
func GetSuggestionDeploymentName(s *suggestionsv1beta1.Suggestion) string {
	return s.Name + "-" + s.Spec.AlgorithmName
}

// GetSuggestionServiceName returns name for the suggestion's service
func GetSuggestionServiceName(s *suggestionsv1beta1.Suggestion) string {
	return s.Name + "-" + s.Spec.AlgorithmName
}

// GetSuggestionPersistentVolumeName returns name for the suggestion's PV
func GetSuggestionPersistentVolumeName(s *suggestionsv1beta1.Suggestion) string {
	return s.Name + "-" + s.Spec.AlgorithmName + "-" + s.Namespace
}

// GetSuggestionPersistentVolumeClaimName returns name for the suggestion's PVC
func GetSuggestionPersistentVolumeClaimName(s *suggestionsv1beta1.Suggestion) string {
	return s.Name + "-" + s.Spec.AlgorithmName
}

// GetSuggestionRBACName returns name for the suggestion's ServiceAccount, Role and RoleBinding
func GetSuggestionRBACName(s *suggestionsv1beta1.Suggestion) string {
	return s.Name + "-" + s.Spec.AlgorithmName
}

// GetAlgorithmEndpoint returns the endpoint of the algorithm service.
func GetAlgorithmEndpoint(s *suggestionsv1beta1.Suggestion) string {
	serviceName := GetSuggestionServiceName(s)
	return fmt.Sprintf("%s.%s:%d",
		serviceName,
		s.Namespace,
		consts.DefaultSuggestionPort)
}

// GetEarlyStoppingEndpoint returns the endpoint of the early stopping service.
func GetEarlyStoppingEndpoint(s *suggestionsv1beta1.Suggestion) string {
	serviceName := GetSuggestionServiceName(s)
	return fmt.Sprintf("%s.%s:%d",
		serviceName,
		s.Namespace,
		consts.DefaultEarlyStoppingPort)
}
