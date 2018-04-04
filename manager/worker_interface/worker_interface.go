package worker_interface

import (
	"github.com/mlkube/katib/api"
)

type WorkerInterface interface {
	IsTrialComplete(studyId string, tID string) (bool, error)
	GetTrialObjValue(studyId string, tID string, objname string) (string, error)
	GetTrialEvLogs(studyId string, tID string, metrics []string, sinceTime string) ([]*api.EvaluationLog, error)
	CheckRunningTrials(studyId string, objname string, metrics []string) error
	SpawnWorkers(trials []*api.Trial, studyId string) error
	GetRunningTrials(studyId string) []*api.Trial
	GetCompletedTrials(studyId string) []*api.Trial
	CleanWorkers(studyId string) error
}
