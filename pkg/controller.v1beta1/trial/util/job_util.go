package util

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	"github.com/tidwall/gjson"
)

// TrialJobStatus is the internal representation for deployed Job status
type TrialJobStatus struct {
	// Condition describes the state of the Job at a certain point.
	// Condition can be either Running, Succeeded or Failed
	Condition ConditionType `json:"condition,omitempty"`

	// The reason received from Job's status, if it is possible
	Reason string `json:"reason,omitempty"`

	// The message received from Job's status, if it is possible
	Message string `json:"message,omitempty"`
}

// ConditionType describes the various conditions a Job can be in.
type ConditionType string

const (
	// JobRunning means that Job was deployed by trial.
	// Job doesn't have succeeded or failed condition.
	JobRunning ConditionType = "Running"
	// JobSucceeded means that Job status satisfies Trial success condition
	JobSucceeded ConditionType = "Succeeded"
	// JobFailed means that Job status satisfies Trial failure condition
	JobFailed ConditionType = "Failed"
)

var (
	log = logf.Log.WithName("job-util")
)

// GetDeployedJobStatus returns internal representation for deployed Job status.
func GetDeployedJobStatus(trial *trialsv1beta1.Trial, deployedJob *unstructured.Unstructured) (*TrialJobStatus, error) {
	logger := log.WithValues("Trial", types.NamespacedName{Name: trial.GetName(), Namespace: trial.GetNamespace()})

	trialJobStatus := &TrialJobStatus{}

	// Marshal unstructured Job to JSON
	// Deployed Job is valid JSON
	deployedJobJSON, _ := util.ConvertUnstructuredToString(deployedJob)

	// Try to get failure condition using spec.failureCondition expression
	failureJobCondition := gjson.Get(deployedJobJSON, trial.Spec.FailureCondition)

	// Condition exists if failureJobCondition is object or failureJobCondition is array with len > 0
	if failureJobCondition.IsObject() || (failureJobCondition.IsArray() && len(failureJobCondition.Array()) > 0) {
		strCondition := failureJobCondition.String()

		// If failureJobCondition is array we take first element
		if failureJobCondition.IsArray() {
			strCondition = failureJobCondition.Array()[0].String()
		}

		// Unmarshal condition to Trial Job representation to get message and reason if it exists
		err := json.Unmarshal([]byte(strCondition), &trialJobStatus)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal failure condition to Trial Job status failed %v", err)
		}

		// Job condition is failed
		trialJobStatus.Condition = JobFailed
		logger.Info("Deployed Job status is failed", "Job", deployedJob.GetName())
		return trialJobStatus, nil
	}

	// Try to get success condition using spec.successCondition expression
	successJobCondition := gjson.Get(deployedJobJSON, trial.Spec.SuccessCondition)

	// The same logic as for failure condition
	if successJobCondition.IsObject() || (successJobCondition.IsArray() && len(successJobCondition.Array()) > 0) {
		strCondition := successJobCondition.String()

		if successJobCondition.IsArray() {
			strCondition = successJobCondition.Array()[0].String()
		}

		err := json.Unmarshal([]byte(strCondition), &trialJobStatus)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal success condition to Trial Job status failed %v", err)
		}

		// Job condition is succeeded failed
		trialJobStatus.Condition = JobSucceeded
		logger.Info("Deployed Job status is succeeded", "Job", deployedJob.GetName())
		return trialJobStatus, nil
	}

	// Set default Job condition is running when name is generated.
	// Check if trial is not running and is not completed
	if !trial.IsRunning() && deployedJob.GetName() != "" && !trial.IsCompleted() {
		trialJobStatus.Condition = JobRunning
		logger.Info("Deployed Job status is running", "Job", deployedJob.GetName())
		return trialJobStatus, nil
	}

	return nil, nil
}
