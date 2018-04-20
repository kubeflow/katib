package worker_interface

import (
	"errors"
	"fmt"
	"github.com/kubeflow/katib/api"
	"github.com/kubeflow/katib/db"
	"time"
)

type WorkerInterface interface {
	IsTrialComplete(studyId string, tID string) (bool, error)
	GetTrialObjValue(studyId string, tID string, objname string) (string, error)
	GetTrialEvLogs(studyId string, tID string, metrics []string, sinceTime string) ([]*api.EvaluationLog, error)
	CheckRunningTrials(studyId string, objname string) error
	SpawnWorkers(trials []*api.Trial, studyId string) error
	GetRunningTrials(studyId string) []*api.Trial
	GetCompletedTrials(studyId string) []*api.Trial
	CleanWorkers(studyId string) error
	CompleteTrial(studyId string, tID string, iscomplete bool) error
}

// Those functions can go after EvalLog transition to DB is complete.
func GetTrialObjValue(d db.VizierDBInterface, studyId string, tID string, objname string) (string, error) {

	log, err := d.GetTrialLogs(tID,
		&db.GetTrialLogOpts{Name: objname, Descending: true, Limit: 1})
	if err != nil {
		return "", err
	}
	if len(log) == 0 {
		return "", errors.New(fmt.Sprintf("No Objective Value Name %v  is found in log", objname))
	}
	return log[0].Value, nil
}

func GetTrialEvLogs(d db.VizierDBInterface, studyId string, tID string, metrics []string, sinceTime string) ([]*api.EvaluationLog, error) {
	var ret []*api.EvaluationLog
	logopts := db.GetTrialLogOpts{}
	if sinceTime != "" {
		t, err := time.Parse(time.RFC3339, sinceTime)
		if err != nil {
			return nil, err
		}
		logopts.SinceTime = &t
	}
	log, err := d.GetTrialLogs(tID, &logopts)
	if err != nil {
		return nil, err
	}
	for _, ls := range log {
		match := false
		for _, m := range metrics {
			if ls.Name == m {
				match = true
				break
			}
		}
		if !match {
			continue
		}
		time_str := ls.Time.Format(time.RFC3339Nano)
		metric := api.Metrics{Name: ls.Name, Value: ls.Value}
		if len(ret) > 0 && ret[len(ret)-1].Time == time_str {
			ret[len(ret)-1].Metrics = append(
				ret[len(ret)-1].Metrics, &metric)
			continue
		}
		ret = append(ret, &api.EvaluationLog{
			Time: time_str, Metrics: []*api.Metrics{&metric}})
	}
	return ret, nil
}
