package common

import (
	v1alpha3 "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
)

type KatibDBInterface interface {
	DBInit()
	SelectOne() error

	RegisterExperiment(experiment *v1alpha3.Experiment) error
	PreCheckRegisterExperiment(experiment *v1alpha3.Experiment) (bool, error)
	DeleteExperiment(experimentName string) error
	GetExperiment(experimentName string) (*v1alpha3.Experiment, error)
	GetExperimentList() ([]*v1alpha3.ExperimentSummary, error)
	UpdateExperimentStatus(experimentName string, newStatus *v1alpha3.ExperimentStatus) error

	UpdateAlgorithmExtraSettings(experimentName string, extraAlgorithmSetting []*v1alpha3.AlgorithmSetting) error
	GetAlgorithmExtraSettings(experimentName string) ([]*v1alpha3.AlgorithmSetting, error)

	RegisterTrial(trial *v1alpha3.Trial) error
	GetTrialList(experimentName string, filter string) ([]*v1alpha3.Trial, error)
	GetTrial(trialName string) (*v1alpha3.Trial, error)
	UpdateTrialStatus(trialName string, newStatus *v1alpha3.TrialStatus) error
	DeleteTrial(trialName string) error

	RegisterObservationLog(trialName string, observationLog *v1alpha3.ObservationLog) error
	GetObservationLog(trialName string, metricName string, startTime string, endTime string) (*v1alpha3.ObservationLog, error)
	DeleteObservationLog(trialName string) error
}
