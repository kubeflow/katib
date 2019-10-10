package v1alpha3

import (
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// JobNameLabel represents the label key for the job name, the value is job name
	JobNameLabel = "job-name"
	// JobRoleLabel represents the label key for the job role, e.g. the value is master
	JobRoleLabel = "job-role"
	// TFJobRoleLabel is deprecated in kubeflow 0.7, but we need to be compatible.
	TFJobRoleLabel = "tf-job-role"
	// PyTorchJobRoleLabel is deprecated in kubeflow 0.7, but we need to be compatible.
	PyTorchJobRoleLabel = "pytorch-job-role"
)

// JobRoleMap is the map which is used to determin if the replica is master.
// Katib will inject metrics collector into master replica.
var JobRoleMap = map[string][]string{
	// Job kind does not support distributed training, thus no master.
	consts.JobKindJob:     {},
	consts.JobKindTF:      {JobRoleLabel, TFJobRoleLabel},
	consts.JobKindPyTorch: {JobRoleLabel, PyTorchJobRoleLabel},
}

// GetSupportedJobList returns the list of the supported jobs' GVK.
func GetSupportedJobList() []schema.GroupVersionKind {
	supportedJobList := []schema.GroupVersionKind{
		{
			Group:   consts.JobGroupJob,
			Version: consts.JobVersionJob,
			Kind:    consts.JobKindJob,
		},
		{
			Group:   consts.JobGroupKubeflow,
			Version: consts.JobVersionTF,
			Kind:    consts.JobKindTF,
		},
		{
			Group:   consts.JobGroupKubeflow,
			Version: consts.JobVersionPyTorch,
			Kind:    consts.JobKindPyTorch,
		},
	}
	return supportedJobList
}
