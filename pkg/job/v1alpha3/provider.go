package v1alpha3

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	"github.com/kubeflow/katib/pkg/job/v1alpha3/job"
	"github.com/kubeflow/katib/pkg/job/v1alpha3/kubeflow"
	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
)

// Provider provides utilities for different jobs.
type Provider interface {
	// GetDeployedJobStatus get the deployed job status.
	GetDeployedJobStatus(
		deployedJob *unstructured.Unstructured) (*commonv1.JobCondition, error)
	// IsTrainingContainer returns if the c is the actual training container.
	IsTrainingContainer(index int, c corev1.Container) bool
}

// New creates a new Provider.
func New(kind string) (Provider, error) {
	switch kind {
	case consts.JobKindJob:
		return &job.Job{}, nil
	case consts.JobKindPyTorch, consts.JobKindTF:
		return &kubeflow.Kubeflow{
			Kind: kind,
		}, nil
	default:
		return nil, fmt.Errorf(
			"Failed to create the provider: Unknown kind %s", kind)
	}
}
