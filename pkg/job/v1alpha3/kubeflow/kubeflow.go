package kubeflow

import (
	pytorchv1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	job "github.com/kubeflow/katib/pkg/job/v1alpha3"
)

var (
	log = logf.Log.WithName("provider-kubeflow")
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
			log.Error(unerr, "NestedFieldCopy unstructured to status error")
			return nil, unerr
		}
		log.Info("NestedFieldCopy unstructured to status error",
			"err", "Status is not found in job")
		return nil, nil
	}

	statusMap := status.(map[string]interface{})
	jobStatus := commonv1.JobStatus{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(statusMap, &jobStatus)
	if err != nil {
		log.Error(err, "Convert unstructured to status error")
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
		log.Info("Invalid Katib worker kind", "JobKind", k.Kind)
		return false
	}
	return false
}

func (k Kubeflow) MutateJob(*v1alpha3.Trial, *unstructured.Unstructured) error {
	return nil
}

func (k *Kubeflow) Create(kind string) job.Provider {
	return &Kubeflow{Kind: kind}
}

func Register() {
	job.ProviderRegistry[consts.JobKindTF] = &Kubeflow{}
	job.SupportedJobList[consts.JobKindTF] = schema.GroupVersionKind{
		Group:   "kubeflow.org",
		Version: "v1",
		Kind:    "TFJob",
	}
	job.JobRoleMap[consts.JobKindTF] = []string{"job-role", "tf-job-role"}
	job.ProviderRegistry[consts.JobKindPyTorch] = &Kubeflow{}
	job.SupportedJobList[consts.JobKindPyTorch] = schema.GroupVersionKind{
		Group:   "kubeflow.org",
		Version: "v1",
		Kind:    "PyTorchJob",
	}
	job.JobRoleMap[consts.JobKindPyTorch] = []string{"job-role", "pytorch-job-role"}
}

func init() {
	Register()
}
