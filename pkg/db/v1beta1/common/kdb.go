package common

import (
	v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

type KatibDBInterface interface {
	DBInit()
	SelectOne() error

	RegisterObservationLog(trialName string, observationLog *v1beta1.ObservationLog) error
	GetObservationLog(trialName string, metricName string, startTime string, endTime string) (*v1beta1.ObservationLog, error)
	DeleteObservationLog(trialName string) error
}
