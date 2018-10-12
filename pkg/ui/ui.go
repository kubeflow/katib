package ui

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/kubeflow/katib/pkg"
	"github.com/kubeflow/katib/pkg/api"

	"github.com/pressly/chi"
	"google.golang.org/grpc"
)

const maxMsgSize = 1<<31 - 1

var colors = [...]string{
	"rgba(255, 99,  132, 0.7)",
	"rgba(54,  162, 235, 0.7)",
	"rgba(255, 206, 86,  0.7)",
	"rgba(75,  192, 192, 0.7)",
	"rgba(153, 102, 255, 0.7)",
	"rgba(255, 159, 64,  0.7)",
}

type IDList struct {
	StudyId  string
	WorkerId string
	TrialId  string
}
type KatibUIHandler struct {
}

func NewKatibUIHandler() *KatibUIHandler {
	return &KatibUIHandler{}
}

func (k *KatibUIHandler) connectManager() (*grpc.ClientConn, api.ManagerClient, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
	}
	conn, err := grpc.Dial(pkg.ManagerAddr, opts...)
	if err != nil {
		log.Printf("Connect katib manager error %v", err)
		return nil, nil, err
	}
	c := api.NewManagerClient(conn)
	return conn, c, nil
}

func (k *KatibUIHandler) Index(w http.ResponseWriter, r *http.Request) {
	conn, c, err := k.connectManager()
	if err != nil {
		return
	}
	defer conn.Close()
	gslrep, err := c.GetStudyList(
		context.Background(),
		&api.GetStudyListRequest{},
	)
	if err != nil {
		log.Printf("Get Study list failed %v", err)
		return
	}
	type StudyListView struct {
		IDList         *IDList
		StudyOverviews []*api.StudyOverview
	}
	slv := StudyListView{
		IDList:         &IDList{},
		StudyOverviews: gslrep.StudyOverviews,
	}
	t, err := template.ParseFiles("/template/layout.html", "/template/index.html", "/template/breadcrumb.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.ExecuteTemplate(w, "layout", slv); err != nil {
		log.Fatal(err)
	}
}

func (k *KatibUIHandler) Study(w http.ResponseWriter, r *http.Request) {
	studyID := chi.URLParam(r, "studyid")
	conn, c, err := k.connectManager()
	if err != nil {
		return
	}
	defer conn.Close()
	type HParam struct {
		Type string
		Name string
	}
	type StudyView struct {
		IDList    *IDList
		StudyConf *api.StudyConfig
		HParams   []*HParam
	}
	gsrep, err := c.GetStudy(
		context.Background(),
		&api.GetStudyRequest{
			StudyId: studyID,
		},
	)
	if err != nil {
		log.Printf("Get Study %s failed %v", studyID, err)
		return
	}
	sv := StudyView{
		IDList: &IDList{
			StudyId: studyID,
		},
		StudyConf: gsrep.StudyConfig,
		HParams:   make([]*HParam, len(gsrep.StudyConfig.ParameterConfigs.Configs)),
	}
	for i, p := range gsrep.StudyConfig.ParameterConfigs.Configs {
		sv.HParams[i] = &HParam{
			Type: "Number",
			Name: p.Name,
		}
		if p.ParameterType == api.ParameterType_CATEGORICAL {
			sv.HParams[i].Type = "String"
		}
	}
	t, err := template.ParseFiles("/template/layout.html", "/template/study.html", "/template/parallelcood.js", "/template/breadcrumb.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.ExecuteTemplate(w, "layout", sv); err != nil {
		log.Fatal(err)
	}
}

func (k *KatibUIHandler) StudyInfoCsv(w http.ResponseWriter, r *http.Request) {
	studyID := chi.URLParam(r, "studyid")
	conn, c, err := k.connectManager()
	if err != nil {
		return
	}
	defer conn.Close()
	retText := "WorkerID,TrialID"
	gsrep, err := c.GetStudy(
		context.Background(),
		&api.GetStudyRequest{
			StudyId: studyID,
		},
	)
	for _, m := range gsrep.StudyConfig.Metrics {
		retText += "," + m
	}
	for _, p := range gsrep.StudyConfig.ParameterConfigs.Configs {
		retText += "," + p.Name
	}
	retText += "\n"
	gwrep, err := c.GetWorkers(
		context.Background(),
		&api.GetWorkersRequest{
			StudyId: studyID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	trialParams := map[string]map[string]string{}
	gtrep, err := c.GetTrials(
		context.Background(),
		&api.GetTrialsRequest{
			StudyId: studyID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range gtrep.Trials {
		trialParams[t.TrialId] = map[string]string{}
		for _, p := range t.ParameterSet {
			trialParams[t.TrialId][p.Name] = p.Value
		}
	}

	workerMetrics := map[string]map[string]string{}
	for _, w := range gwrep.Workers {
		workerMetrics[w.WorkerId] = map[string]string{}
		gmrep, err := c.GetMetrics(
			context.Background(),
			&api.GetMetricsRequest{
				StudyId:   studyID,
				WorkerIds: []string{w.WorkerId},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		if len(gmrep.MetricsLogSets) > 0 {
			for _, m := range gmrep.MetricsLogSets[0].MetricsLogs {
				if len(m.Values) > 0 {
					workerMetrics[w.WorkerId][m.Name] = m.Values[len(m.Values)-1].Value
				}
			}
			retText += w.WorkerId + "," + w.TrialId
			for _, m := range gsrep.StudyConfig.Metrics {
				retText += "," + workerMetrics[w.WorkerId][m]
			}
			for _, p := range gsrep.StudyConfig.ParameterConfigs.Configs {
				retText += "," + trialParams[w.TrialId][p.Name]
			}
		}
		retText += "\n"
	}
	fmt.Fprint(w, retText)
}

func (k *KatibUIHandler) Trial(w http.ResponseWriter, r *http.Request) {
	studyID := chi.URLParam(r, "studyid")
	trialID := chi.URLParam(r, "trialid")
	conn, c, err := k.connectManager()
	if err != nil {
		return
	}
	defer conn.Close()
	type TrialView struct {
		IDList  *IDList
		Trial   *api.Trial
		Workers []*api.Worker
	}
	gtrep, err := c.GetTrials(
		context.Background(),
		&api.GetTrialsRequest{
			StudyId: studyID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	view := TrialView{
		IDList: &IDList{
			StudyId: studyID,
			TrialId: trialID,
		},
	}
	for _, t := range gtrep.Trials {
		if t.TrialId == trialID {
			view.Trial = t
		}
	}
	gwrep, err := c.GetWorkers(
		context.Background(),
		&api.GetWorkersRequest{
			StudyId: studyID,
			TrialId: trialID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	view.Workers = gwrep.Workers
	t, err := template.ParseFiles("/template/layout.html", "/template/trial.html", "/template/breadcrumb.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.ExecuteTemplate(w, "layout", view); err != nil {
		log.Fatal(err)
	}
}

func (k *KatibUIHandler) WorkerInfoCsv(w http.ResponseWriter, r *http.Request) {
	studyID := chi.URLParam(r, "studyid")
	workerID := chi.URLParam(r, "workerid")
	conn, c, err := k.connectManager()
	if err != nil {
		return
	}
	defer conn.Close()
	retText := "symbol,time,value\n"
	gmrep, err := c.GetMetrics(
		context.Background(),
		&api.GetMetricsRequest{
			StudyId:   studyID,
			WorkerIds: []string{workerID},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(gmrep.MetricsLogSets) > 0 {
		for _, m := range gmrep.MetricsLogSets[0].MetricsLogs {
			pvtime := ""
			for _, v := range m.Values {
				mvtime, _ := time.Parse(time.RFC3339Nano, v.Time)
				ctime := mvtime.Format("2006-01-02T15:4:5")
				if pvtime != ctime {
					retText += m.Name + "," + ctime + "," + v.Value + "\n"
					pvtime = ctime
				}
			}
		}
	}
	fmt.Fprint(w, retText)
}

func (k *KatibUIHandler) Worker(w http.ResponseWriter, r *http.Request) {
	studyID := chi.URLParam(r, "studyid")
	workerID := chi.URLParam(r, "workerid")
	conn, c, err := k.connectManager()
	if err != nil {
		return
	}
	defer conn.Close()
	type TimeValue struct {
		Time  float64
		Value string
	}
	type MetricsLog struct {
		Name      string
		Color     string
		LogValues []TimeValue
	}
	type WorkerView struct {
		IDList      *IDList
		Parameters  []*api.Parameter
		MetricsLogs []MetricsLog
	}
	gwrep, err := c.GetWorkers(
		context.Background(),
		&api.GetWorkersRequest{
			StudyId:  studyID,
			WorkerId: workerID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(gwrep.Workers) != 1 {
		fmt.Fprint(w, "Worker ID is wrong.")
		return
	}
	worker := gwrep.Workers[0]
	wv := WorkerView{
		IDList: &IDList{
			StudyId:  studyID,
			WorkerId: workerID,
			TrialId:  worker.TrialId,
		},
	}
	gtrep, err := c.GetTrials(
		context.Background(),
		&api.GetTrialsRequest{
			StudyId: studyID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range gtrep.Trials {
		if t.TrialId == worker.TrialId {
			wv.Parameters = t.ParameterSet
		}
	}

	gmrep, err := c.GetMetrics(
		context.Background(),
		&api.GetMetricsRequest{
			StudyId:   studyID,
			WorkerIds: []string{workerID},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(gmrep.MetricsLogSets) > 0 {
		wv.MetricsLogs = make([]MetricsLog, len(gmrep.MetricsLogSets[0].MetricsLogs))
		for i, m := range gmrep.MetricsLogSets[0].MetricsLogs {
			wv.MetricsLogs[i].Name = m.Name
			wv.MetricsLogs[i].Color = colors[i%len(colors)]
			wv.MetricsLogs[i].LogValues = []TimeValue{}
			var pvtime float64
			var baseTime time.Time
			if len(m.Values) > 0 {
				baseTime, _ = time.Parse(time.RFC3339Nano, m.Values[0].Time)
			}
			for _, v := range m.Values {
				mvtime, _ := time.Parse(time.RFC3339Nano, v.Time)
				tdiff := mvtime.Sub(baseTime)
				ctime := tdiff.Seconds()
				if pvtime != ctime {
					wv.MetricsLogs[i].LogValues = append(
						wv.MetricsLogs[i].LogValues,
						TimeValue{
							Time:  ctime,
							Value: v.Value,
						},
					)
					pvtime = ctime
				}
			}
			fmt.Printf("Log %s %v\n", wv.MetricsLogs[i].Name, wv.MetricsLogs[i].LogValues)
		}
	}
	t, err := template.ParseFiles("/template/layout.html", "/template/worker.html", "/template/linegraph.js", "/template/breadcrumb.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.ExecuteTemplate(w, "layout", wv); err != nil {
		log.Fatal(err)
	}
}
