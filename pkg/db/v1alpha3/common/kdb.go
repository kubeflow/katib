package common

import (
	v1alpha3 "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
)

type KatibDBInterface interface {
	DBInit()
	SelectOne() error

	RegisterObservationLog(trialName string, observationLog *v1alpha3.ObservationLog) error
	GetObservationLog(trialName string, metricName string, startTime string, endTime string) (*v1alpha3.ObservationLog, error)
	DeleteObservationLog(trialName string) error
}
