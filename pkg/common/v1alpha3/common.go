package v1alpha3

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertTime2RFC3339(t *metav1.Time) string {
	if t != nil {
		return t.UTC().Format(time.RFC3339)
	} else {
		zero := &metav1.Time{}
		return zero.UTC().Format(time.RFC3339)
	}
}

func GetJobLabelMap(jobKind string, trialName string) map[string]string {
	labelMap := make(map[string]string)

	if jobKind == "TFJob" {
		labelMap["job-name"] = trialName
		labelMap["job-role"] = "master"
	} else if jobKind == "PyTorchJob" {
		labelMap["job-name"] = trialName
		labelMap["job-role"] = "master"
	} else {
		labelMap["job-name"] = trialName
	}

	return labelMap
}
