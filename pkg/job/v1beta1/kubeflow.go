package v1beta1

import (
	pytorchv1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

var (
	kfLogger = logf.Log.WithName("provider-kubeflow")
)

// Kubeflow is the provider of Kubeflow kinds.
type Kubeflow struct {
	Kind string
}

// GetDeployedJobStatus get the deployed job status.
func (k Kubeflow) GetDeployedJobStatus(
	deployedJob *unstructured.Unstructured) (*commonv1.JobCondition, error) {
	jobCondition := commonv1.JobCondition{}
	// Set default type to running.
	jobCondition.Type = commonv1.JobRunning
	status, ok, unerr := unstructured.NestedFieldCopy(deployedJob.Object, "status")
	if !ok {
		if unerr != nil {
			kfLogger.Error(unerr, "NestedFieldCopy unstructured to status error")
			return nil, unerr
		}
		kfLogger.Info("NestedFieldCopy unstructured to status error",
			"Info", "Kubeflow Job doesn't have status yet")
		return nil, nil
	}

	statusMap := status.(map[string]interface{})
	jobStatus := commonv1.JobStatus{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(statusMap, &jobStatus)
	if err != nil {
		kfLogger.Error(err, "Convert unstructured to status error")
		return nil, err
	}
	// Get the latest condition and set it to jobCondition.
	if len(jobStatus.Conditions) > 0 {
		lc := jobStatus.Conditions[len(jobStatus.Conditions)-1]
		jobCondition.Type = lc.Type
		jobCondition.Message = lc.Message
		jobCondition.Status = lc.Status
		jobCondition.Reason = lc.Reason
	}

	return &jobCondition, nil
}

// IsTrainingContainer returns if the c is the actual training container.
func (k Kubeflow) IsTrainingContainer(index int, c corev1.Container) bool {
	switch k.Kind {
	case consts.JobKindTF:
		if c.Name == tfv1.DefaultContainerName {
			return true
		}
	case consts.JobKindPyTorch:
		if c.Name == pytorchv1.DefaultContainerName {
			return true
		}
	default:
		kfLogger.Info("Invalid Katib worker kind", "JobKind", k.Kind)
		return false
	}
	return false
}

func (k Kubeflow) MutateJob(*v1beta1.Trial, *unstructured.Unstructured) error {
	return nil
}

func (k *Kubeflow) Create(kind string) Provider {
	return &Kubeflow{Kind: kind}
}

func init() {
	ProviderRegistry[consts.JobKindTF] = &Kubeflow{}
	SupportedJobList[consts.JobKindTF] = schema.GroupVersionKind{
		Group:   "kubeflow.org",
		Version: "v1",
		Kind:    consts.JobKindTF,
	}
	JobRoleMap[consts.JobKindTF] = []string{consts.JobRole, consts.JobRoleTF}
	ProviderRegistry[consts.JobKindPyTorch] = &Kubeflow{}
	SupportedJobList[consts.JobKindPyTorch] = schema.GroupVersionKind{
		Group:   "kubeflow.org",
		Version: "v1",
		Kind:    consts.JobKindPyTorch,
	}
	JobRoleMap[consts.JobKindPyTorch] = []string{consts.JobRole, consts.JobRolePyTorch}
}
