package earlystopping

import (
	"context"
	"errors"
	"github.com/kubeflow/katib/pkg/api"
	vdb "github.com/kubeflow/katib/pkg/db"
	"log"
	"sort"
	"strconv"
)

const (
	defaultLeastStep = 20
	defaultMargin    = 0
	defaultBurnIn    = 0
)

type MedianStoppingParam struct {
	LeastStep  int
	Margin     float64
	EvalMetric string
	BurnIn     int
}

type MedianStoppingRule struct {
	dbIf vdb.VizierDBInterface
}

func NewMedianStoppingRule() *MedianStoppingRule {
	m := &MedianStoppingRule{}
	m.dbIf = vdb.New()
	return m
}

func (m *MedianStoppingRule) purseEarlyStoppingParameters(sc *api.StudyConfig, eps []*api.EarlyStoppingParameter) (*MedianStoppingParam, error) {
	p := &MedianStoppingParam{LeastStep: defaultLeastStep, Margin: defaultMargin, EvalMetric: sc.ObjectiveValueName, BurnIn: defaultBurnIn}
	for _, ep := range eps {
		switch ep.Name {
		case "LeastStep":
			l, err := strconv.Atoi(ep.Value)
			if err != nil {
				log.Printf("Fail to puerse parameter %s : %s", ep.Name, ep.Value)
			} else {
				p.LeastStep = l
			}
		case "EvalMargin":
			mar, err := strconv.ParseFloat(ep.Value, 64)
			if err != nil {
				log.Printf("Fail to puerse parameter %s : %s", ep.Name, ep.Value)
			} else {
				p.Margin = mar
			}
		case "EvalMetrics":
			p.EvalMetric = ep.Value
		case "BurnInPeriod":
			b, err := strconv.Atoi(ep.Value)
			if err != nil {
				log.Printf("Fail to puerse parameter %s : %s", ep.Name, ep.Value)
			} else {
				p.BurnIn = b
			}
		default:
			log.Printf("Unknown EarlyStopping Parameter %v", ep.Name)
		}
	}
	log.Printf("Parameter: LeastStep %d, Margin %v, EvalMetric %s, BurnInPeriod %d", p.LeastStep, p.Margin, p.EvalMetric, p.BurnIn)
	return p, nil
}

func (m *MedianStoppingRule) getMedianRunningAverage(completedWorkerslogs [][]*vdb.WorkerLog, step int, burnin int) float64 {
	r := []float64{}
	var ra float64
	for _, cwl := range completedWorkerslogs {
		ra = 0
		var st int
		var errParce bool = false
		if step > len(cwl) {
			st = len(cwl)
		} else {
			st = step
		}
		for s := burnin; s < st; s++ {
			v, err := strconv.ParseFloat(cwl[s].Value, 64)
			if err != nil {
				log.Printf("Fail to Parse %s : %s", cwl[s].Name, cwl[s].Value)
				errParce = true
				break
			}
			ra += v
		}
		if errParce {
			continue
		}
		ra = ra / float64(st)
		r = append(r, ra)
	}
	if len(r) == 0 {
		return 0
	} else {
		sort.Float64s(r)
		return r[len(r)/2]
	}
}

func (m *MedianStoppingRule) getBestValue(sid string, sc *api.StudyConfig, logs []*vdb.WorkerLog) (float64, error) {
	if len(logs) == 0 {
		return 0, errors.New("Evaluation Log is missing")
	}
	ot := sc.OptimizationType
	if ot != api.OptimizationType_MAXIMIZE && ot != api.OptimizationType_MINIMIZE {
		return 0, errors.New("OptimizationType Unknown.")
	}
	var ret float64
	var target_objlog []float64
	for _, l := range logs {
		v, err := strconv.ParseFloat(l.Value, 64)
		if err != nil {
			log.Printf("Fail to Parse %s : %s", l.Name, l.Value)
			continue
		}
		target_objlog = append(target_objlog, v)
	}
	if len(target_objlog) == 0 {
		return 0, errors.New("No Objective value log in Logs")
	}
	sort.Float64s(target_objlog)
	if ot == api.OptimizationType_MAXIMIZE {
		ret = target_objlog[len(target_objlog)-1]
	} else if ot == api.OptimizationType_MINIMIZE {
		ret = target_objlog[0]
	}
	return ret, nil
}
func (m *MedianStoppingRule) GetShouldStopWorkers(ctx context.Context, in *api.GetShouldStopWorkersRequest) (*api.GetShouldStopWorkersReply, error) {
	wl, err := m.dbIf.GetWorkerList(in.StudyId, "")
	if err != nil {
		return &api.GetShouldStopWorkersReply{}, err
	}
	sc, err := m.dbIf.GetStudyConfig(in.StudyId)
	if err != nil {
		return &api.GetShouldStopWorkersReply{}, err
	}
	eparam, err := m.dbIf.GetEarlyStopParam(in.ParamId)
	if err != nil {
		return &api.GetShouldStopWorkersReply{}, err
	}
	p, err := m.purseEarlyStoppingParameters(sc, eparam)
	if err != nil {
		return &api.GetShouldStopWorkersReply{}, err
	}

	rwids := []string{}
	cwl := make([][]*vdb.WorkerLog, 0, len(wl))
	s_w := []string{}
	for _, w := range wl {
		switch w.Status {
		case api.State_RUNNING:
			rwids = append(rwids, w.WorkerId)
		case api.State_COMPLETED:
			wl, err := m.dbIf.GetWorkerLogs(w.WorkerId, &vdb.GetWorkerLogOpts{Name: p.EvalMetric})
			if err != nil {
				log.Printf("Fail to get worker %v logs", w.WorkerId)
				continue
			}
			if len(wl) > p.BurnIn {
				cwl = append(cwl, wl)
			}
		default:
		}
	}
	if len(cwl) == 0 {
		return &api.GetShouldStopWorkersReply{}, err
	}
	for _, w := range rwids {
		wl, err := m.dbIf.GetWorkerLogs(w, &vdb.GetWorkerLogOpts{Name: p.EvalMetric})
		if err != nil {
			log.Printf("Fail to get worker %v logs", w)
			continue
		}
		if len(wl) < p.LeastStep || len(wl) <= p.BurnIn {
			continue
		}
		v, err := m.getBestValue(in.StudyId, sc, wl)
		if err != nil {
			log.Printf("Fail to Get Best Value at %s: %v Log:%v", w, err, wl)
			continue
		}
		om := m.getMedianRunningAverage(cwl, len(wl), p.BurnIn)
		log.Printf("Worker %s, In step %d Current value: %v Median value: %v\n", w, len(wl), v, om)
		if (v < (om-p.Margin) && sc.OptimizationType == api.OptimizationType_MAXIMIZE) || v > (om+p.Margin) && sc.OptimizationType == api.OptimizationType_MINIMIZE {
			log.Printf("Worker %s shuold be stopped", w)
			s_w = append(s_w, w)
		}
	}
	return &api.GetShouldStopWorkersReply{ShouldStopWorkerIds: s_w}, nil
}
