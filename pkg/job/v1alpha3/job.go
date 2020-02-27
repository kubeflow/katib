package v1alpha3

import (
	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
)

var (
	jobLogger = logf.Log.WithName("provider-job")
)

// Job is the provider of Job kind.
type Job struct{}

// GetDeployedJobStatus get the deployed job status.
func (j Job) GetDeployedJobStatus(
	deployedJob *unstructured.Unstructured) (*commonv1.JobCondition, error) {
	jobCondition := commonv1.JobCondition{}
	// Set default type to running.
	jobCondition.Type = commonv1.JobRunning
	status, ok, unerr := unstructured.NestedFieldCopy(deployedJob.Object, "status")
	if !ok {
		if unerr != nil {
			jobLogger.Error(unerr, "NestedFieldCopy unstructured to status error")
			return nil, unerr
		}
		// Job does not have the running condition in status, thus we think
		// the job is running when it is created.
		jobLogger.Info("NestedFieldCopy", "err", "status cannot be found in job")
		return nil, nil
	}

	statusMap := status.(map[string]interface{})
	jobStatus := batchv1.JobStatus{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(statusMap, &jobStatus)
	if err != nil {
		jobLogger.Error(err, "Convert unstructured to status error")
		return nil, err
	}
	for _, cond := range jobStatus.Conditions {
		if cond.Type == batchv1.JobComplete && cond.Status == corev1.ConditionTrue {
			jobCondition.Type = commonv1.JobSucceeded
			//  JobConditions message not populated when succeeded for batchv1 Job
			break
		}
		if cond.Type == batchv1.JobFailed && cond.Status == corev1.ConditionTrue {
			jobCondition.Type = commonv1.JobFailed
			jobCondition.Message = cond.Message
			break
		}
	}
	return &jobCondition, nil
}

// IsTrainingContainer returns if the c is the actual training container.
func (j Job) IsTrainingContainer(index int, c corev1.Container) bool {
	if index == 0 {
		// for Job worker, the first container will be taken as worker container,
		// katib document should note it
		return true
	}
	return false
}
func (j Job) MutateJob(*v1alpha3.Trial, *unstructured.Unstructured) error {
	return nil
}

func (j *Job) Create(kind string) Provider {
	return &Job{}
}

func init() {
	ProviderRegistry[consts.JobKindJob] = &Job{}
	SupportedJobList[consts.JobKindJob] = schema.GroupVersionKind{
		Group:   "batch",
		Version: "v1",
		Kind:    "Job",
	}
	JobRoleMap[consts.JobKindJob] = []string{}
}
