package v1alpha2

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetSupportedJobList() []schema.GroupVersionKind {
	supportedJobList := []schema.GroupVersionKind{
		schema.GroupVersionKind{
			Group:   "batch",
			Version: "v1",
			Kind:    "Job",
		},
		schema.GroupVersionKind{
			Group:   "kubeflow.org",
			Version: "v1beta2",
			Kind:    "TFJob",
		},
		schema.GroupVersionKind{
			Group:   "kubeflow.org",
			Version: "v1beta2",
			Kind:    "PyTorchJob",
		},
	}
	return supportedJobList
}
