package v1beta1

import (
	"fmt"

	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
)

var (
	ProviderRegistry = make(map[string]Provider)
	// JobRoleMap is the map which is used to determin if the replica is master.
	// Katib will inject metrics collector into master replica.
	JobRoleMap = make(map[string][]string)
	// SupportedJobList returns the list of the supported jobs' GVK.
	SupportedJobList = make(map[string]schema.GroupVersionKind)
)

// Provider provides utilities for different jobs.
type Provider interface {
	// GetDeployedJobStatus get the deployed job status.
	GetDeployedJobStatus(
		deployedJob *unstructured.Unstructured) (*commonv1.JobCondition, error)
	// IsTrainingContainer returns if the c is the actual training container.
	IsTrainingContainer(index int, c corev1.Container) bool
	// Mutate jobSpec before creation if necessary
	MutateJob(*v1beta1.Trial, *unstructured.Unstructured) error
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
