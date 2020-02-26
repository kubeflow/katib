package v1alpha3

import (
	"fmt"

	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
)

var (
	ProviderRegistry = make(map[string]Provider)
)

// Provider provides utilities for different jobs.
type Provider interface {
	// GetDeployedJobStatus get the deployed job status.
	GetDeployedJobStatus(
		deployedJob *unstructured.Unstructured) (*commonv1.JobCondition, error)
	// IsTrainingContainer returns if the c is the actual training container.
	IsTrainingContainer(index int, c corev1.Container) bool
	// Mutate jobSpec before creation if necessary
	MutateJob(*v1alpha3.Trial, *unstructured.Unstructured) error
	// Recreate Provider from kind
	Create(kind string) Provider
}

// New creates a new Provider.
func New(kind string) (Provider, error) {
	if ptr, ok := ProviderRegistry[kind]; ok {
		return ptr.Create(kind), nil
	} else {
		return nil, fmt.Errorf(
			"failed to create the provider: Unknown kind %s", kind)
	}
}
