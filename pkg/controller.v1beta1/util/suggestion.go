package util

import (
	"fmt"

	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

// GetAlgorithmDeploymentName returns name for the suggestion's deployment
func GetAlgorithmDeploymentName(s *suggestionsv1beta1.Suggestion) string {
	return s.Name + "-" + s.Spec.AlgorithmName
}

// GetAlgorithmServiceName returns name for the suggestion's service
func GetAlgorithmServiceName(s *suggestionsv1beta1.Suggestion) string {
	return s.Name + "-" + s.Spec.AlgorithmName
}

// GetAlgorithmPersistentVolumeName returns name for the suggestion's PV
func GetAlgorithmPersistentVolumeName(s *suggestionsv1beta1.Suggestion) string {
	return s.Name + "-" + s.Spec.AlgorithmName + "-" + s.Namespace
}

// GetAlgorithmPersistentVolumeClaimName returns name for the suggestion's PVC
func GetAlgorithmPersistentVolumeClaimName(s *suggestionsv1beta1.Suggestion) string {
	return s.Name + "-" + s.Spec.AlgorithmName
}

// GetAlgorithmEndpoint returns the endpoint of the algorithm service.
func GetAlgorithmEndpoint(s *suggestionsv1beta1.Suggestion) string {
	serviceName := GetAlgorithmServiceName(s)
	return fmt.Sprintf("%s.%s:%d",
		serviceName,
		s.Namespace,
		consts.DefaultSuggestionPort)
}
