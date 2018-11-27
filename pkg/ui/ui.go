package ui

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kubeflow/katib/pkg"
	"github.com/kubeflow/katib/pkg/api"
	studyjobv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	"github.com/kubeflow/katib/pkg/manager/studyjobclient"

	"github.com/pressly/chi"
	"google.golang.org/grpc"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

const maxMsgSize = 1<<31 - 1

var colors = [...]string{
	"rgba(255, 99,  132, 0.6)",
	"rgba(54,  162, 235, 0.6)",
	"rgba(255, 206, 86,  0.6)",
	"rgba(75,  192, 192, 0.6)",
	"rgba(153, 102, 255, 0.6)",
	"rgba(255, 159, 64,  0.6)",
}

type IDList struct {
	StudyId  string
	WorkerId string
	TrialId  string
}
type KatibUIHandler struct {
	studyjobClient *studyjobclient.StudyjobClient
}

func NewKatibUIHandler() *KatibUIHandler {
	sjc, err := studyjobclient.NewStudyjobClient(nil)
	if err != nil {
		panic(err)
	}
	return &KatibUIHandler{
		studyjobClient: sjc,
	}
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
	type StudyNameStack struct {
		StudyId           string
		Owner             string
		StudyJobName      string
		StudyJobCondition string
	}
	type StudySummary struct {
		StudyNameStacks    []*StudyNameStack
		LatestJobCondition string
		LastJobUpdateTime  *time.Time
	}
	type StudyListView struct {
		IDList        *IDList
		StudySummarys map[string]*StudySummary
	}
	slv := &StudyListView{
		IDList:        &IDList{},
		StudySummarys: make(map[string]*StudySummary),
	}
	for _, so := range gslrep.StudyOverviews {
		ss, ok := slv.StudySummarys[so.Name]
		if !ok {
			ss = &StudySummary{
				StudyNameStacks:    []*StudyNameStack{},
				LatestJobCondition: "Unknown",
			}
			slv.StudySummarys[so.Name] = ss
		}
		ss.StudyNameStacks = append(ss.StudyNameStacks, &StudyNameStack{
			StudyId:           so.Id,
			Owner:             so.Owner,
			StudyJobName:      "None",
			StudyJobCondition: "None",
		})
	}
	sl, err := k.studyjobClient.GetStudyJobList()
	if err != nil {
		log.Printf("StudyjobClient List err %v", err)
	} else {
		for _, sj := range sl.Items {
			ss, ok := slv.StudySummarys[sj.Spec.StudyName]
			if !ok {
				continue
			}
			for i := range ss.StudyNameStacks {
				if ss.StudyNameStacks[i].StudyId == sj.Status.StudyID {
					ss.StudyNameStacks[i].StudyJobName = sj.Name
					ss.StudyNameStacks[i].StudyJobCondition = string(sj.Status.Condition)
					if sj.Status.LastReconcileTime != nil {
						if ss.LastJobUpdateTime != nil {
							if ss.LastJobUpdateTime.Before(sj.Status.LastReconcileTime.Time) {
								continue
							}
						}
						ss.LatestJobCondition = string(sj.Status.Condition)
						ss.LastJobUpdateTime = &sj.Status.LastReconcileTime.Time
					}
				}
			}
		}
	}
	t, err := template.ParseFiles("/app/template/layout.html", "/app/template/index.html", "/app/template/breadcrumb.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.ExecuteTemplate(w, "layout", slv); err != nil {
		log.Fatal(err)
	}
}

func (k *KatibUIHandler) CreateStudyJob(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.PostForm
	sjm := form["StudyJobManifest"][0]
	var job *studyjobv1alpha1.StudyJob
	BUFSIZE := 1024
	k8syaml.NewYAMLOrJSONDecoder(strings.NewReader(sjm), BUFSIZE).Decode(&job)
	_, err := k.studyjobClient.CreateStudyJob(job)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fail to create Studyjob: " + err.Error()))
		return
	}
}

func (k *KatibUIHandler) StudyJobGen(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	sid := query.Get("studyid")
	sname := query.Get("studyname")
	conn, c, err := k.connectManager()
	defer conn.Close()
	if err != nil {
		log.Printf("Failed to connect Katib manager %v", err)
	}
	if sid == "" && sname != "" {
		sl, err := k.studyjobClient.GetStudyJobList()
		if err != nil {
			log.Printf("StudyjobClient List err %v", err)
		} else {
			var ltime *time.Time
			for _, sj := range sl.Items {
				if sj.Spec.StudyName == sname {
					if ltime != nil {
						if sj.Status.LastReconcileTime != nil {
							if ltime.Before(sj.Status.LastReconcileTime.Time) {
								continue
							}
						}
					}
					ltime = &sj.Status.LastReconcileTime.Time
					sid = sj.Status.StudyID
				}
			}
		}
		if sid == "" {
			if c != nil {
				gslrep, err := c.GetStudyList(
					context.Background(),
					&api.GetStudyListRequest{},
				)
				if err != nil {
					log.Printf("Get Study List err %v", err)
				} else {
					for _, so := range gslrep.StudyOverviews {
						if so.Name == sname {
							sid = so.Id
							break
						}
					}
				}
			}
		}
	}
	wt, err := k.studyjobClient.GetWorkerTemplates()
	if err != nil {
		log.Printf("GetWorkerTemplates err %v", err)
	}
	type Param struct {
		Name string
		Type string
		Min  string
		Max  string
		List string
	}
	type StudyJobDefault struct {
		IDList             *IDList
		StudyName          string
		Owner              string
		OptimizationType   string
		OptimizationGoal   float64
		ObjectiveValueName string
		Metrics            string
		ParamConf          []*Param
		WorkerTemplates    map[string]string
	}
	sjd := &StudyJobDefault{
		IDList:          &IDList{},
		WorkerTemplates: wt,
	}
	if sid != "" {
		if c != nil {
			gsrep, err := c.GetStudy(
				context.Background(),
				&api.GetStudyRequest{
					StudyId: sid,
				},
			)
			if err != nil {
				log.Printf("Get Study %s failed %v", sid, err)
			} else {
				sjd.StudyName = gsrep.StudyConfig.Name
				sjd.Owner = gsrep.StudyConfig.Owner
				sjd.OptimizationType = strings.ToLower(
					string(
						api.OptimizationType_name[int32(gsrep.StudyConfig.OptimizationType)],
					),
				)
				sjd.OptimizationGoal = gsrep.StudyConfig.OptimizationGoal
				sjd.ObjectiveValueName = gsrep.StudyConfig.ObjectiveValueName
				sjd.Metrics = strings.Join(gsrep.StudyConfig.Metrics, " ")
				sjd.ParamConf = make([]*Param, len(gsrep.StudyConfig.ParameterConfigs.Configs))
				for i, p := range gsrep.StudyConfig.ParameterConfigs.Configs {
					sjd.ParamConf[i] = &Param{}
					sjd.ParamConf[i].Name = p.Name
					sjd.ParamConf[i].Type = strings.ToLower(
						string(
							api.ParameterType_name[int32(p.ParameterType)],
						),
					)
					sjd.ParamConf[i].Min = p.Feasible.Min
					sjd.ParamConf[i].Max = p.Feasible.Max
					sjd.ParamConf[i].List = strings.Join(p.Feasible.List, " ")
				}
			}
		}
	}
	t, err := template.ParseFiles("/app/template/layout.html", "/app/template/studyjobgen.html", "/app/template/studyjobgen.js", "/app/template/breadcrumb.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.ExecuteTemplate(w, "layout", sjd); err != nil {
		log.Fatal(err)
	}
}

func (k *KatibUIHandler) WorkerTemplate(w http.ResponseWriter, r *http.Request) {
	wt, err := k.studyjobClient.GetWorkerTemplates()
	if err != nil {
		log.Printf("GetWorkerTemplates err %v", err)
	}
	t, err := template.ParseFiles("/app/template/layout.html", "/app/template/workertemplate.html", "/app/template/workertemplate.js", "/app/template/breadcrumb.html")
	type WorkerTemplateView struct {
		IDList         *IDList
		WorkerTemplate map[string]string
	}
	wtv := WorkerTemplateView{
		IDList:         &IDList{},
		WorkerTemplate: wt,
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := t.ExecuteTemplate(w, "layout", wtv); err != nil {
		log.Fatal(err)
	}
}

func (k *KatibUIHandler) UpdateWorkerTemplate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.PostForm
	wt := make(map[string]string, len(form))
	for k, v := range form {
		if len(v) > 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Template Name duplicated"))
			return
		}
		wt[k] = v[0]
	}
	err := k.studyjobClient.UpdateWorkerTemplates(wt)
	if err != nil {
		log.Print("fail to UpdateWorkerTemplate %v", err)
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
	t, err := template.ParseFiles("/app/template/layout.html", "/app/template/study.html", "/app/template/parallelcood.js", "/app/template/breadcrumb.html")
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
		log.Println(err)
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
	if err != nil {
		log.Println(err)
		return
	}
	metricsList := map[string]int{}
	for i, m := range gsrep.StudyConfig.Metrics {
		retText += "," + m
		metricsList[m] = i
	}
	paramList := map[string]int{}
	for i, p := range gsrep.StudyConfig.ParameterConfigs.Configs {
		retText += "," + p.Name
		paramList[p.Name] = i + len(metricsList)
	}
	gwfirep, err := c.GetWorkerFullInfo(
		context.Background(),
		&api.GetWorkerFullInfoRequest{
			StudyId:       studyID,
			OnlyLatestLog: true,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	retText += "\n"
	for _, wfi := range gwfirep.WorkerFullInfos {
		restext := make([]string, len(metricsList)+len(paramList))
		for _, m := range wfi.MetricsLogs {
			if len(m.Values) > 0 {
				restext[metricsList[m.Name]] = m.Values[len(m.Values)-1].Value
			}
		}
		for _, p := range wfi.ParameterSet {
			restext[paramList[p.Name]] = p.Value
		}
		retText += wfi.Worker.WorkerId + "," + wfi.Worker.TrialId + "," + strings.Join(restext, ",") + "\n"
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
	t, err := template.ParseFiles("/app/template/layout.html", "/app/template/trial.html", "/app/template/breadcrumb.html")
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
	gwfirep, err := c.GetWorkerFullInfo(
		context.Background(),
		&api.GetWorkerFullInfoRequest{
			StudyId:  studyID,
			WorkerId: workerID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(gwfirep.WorkerFullInfos) > 0 {
		for _, m := range gwfirep.WorkerFullInfos[0].MetricsLogs {
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
	gwfirep, err := c.GetWorkerFullInfo(
		context.Background(),
		&api.GetWorkerFullInfoRequest{
			StudyId:  studyID,
			WorkerId: workerID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(gwfirep.WorkerFullInfos) != 1 {
		fmt.Fprint(w, "Worker ID is wrong.")
		return
	}
	worker := gwfirep.WorkerFullInfos[0].Worker
	wv := WorkerView{
		IDList: &IDList{
			StudyId:  studyID,
			WorkerId: workerID,
			TrialId:  worker.TrialId,
		},
		Parameters: gwfirep.WorkerFullInfos[0].ParameterSet,
	}
	wv.MetricsLogs = make([]MetricsLog, len(gwfirep.WorkerFullInfos[0].MetricsLogs))
	for i, m := range gwfirep.WorkerFullInfos[0].MetricsLogs {
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
	t, err := template.ParseFiles("/app/template/layout.html", "/app/template/worker.html", "/app/template/linegraph.js", "/app/template/breadcrumb.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.ExecuteTemplate(w, "layout", wv); err != nil {
		log.Fatal(err)
	}
}
