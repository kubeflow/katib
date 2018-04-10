package earlystopping

import (
	"context"
	"errors"
	"github.com/kubeflow/hp-tuning/api"
	vdb "github.com/kubeflow/hp-tuning/db"
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
	confList map[string]*MedianStoppingParam
	dbIf     vdb.VizierDBInterface
}

func NewMedianStoppingRule() *MedianStoppingRule {
	m := &MedianStoppingRule{}
	m.confList = make(map[string]*MedianStoppingParam)
	m.dbIf = vdb.New()
	return m
}

func (m *MedianStoppingRule) SetEarlyStoppingParameter(ctx context.Context, in *api.SetEarlyStoppingParameterRequest) (*api.SetEarlyStoppingParameterReply, error) {
	sc, err := m.dbIf.GetStudyConfig(in.StudyId)
	if err != nil {
		return &api.SetEarlyStoppingParameterReply{}, err
	}
	p := &MedianStoppingParam{LeastStep: defaultLeastStep, Margin: defaultMargin, EvalMetric: sc.ObjectiveValueName, BurnIn: defaultBurnIn}
	for _, ep := range in.EarlyStoppingParameters {
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
	m.confList[in.StudyId] = p
	log.Printf("Parameter for Study %s : LeastStep %d, Margin %v, EvalMetric %s, BurnInPeriod %d", in.StudyId, p.LeastStep, p.Margin, p.EvalMetric, p.BurnIn)
	return &api.SetEarlyStoppingParameterReply{}, nil
}

func (m *MedianStoppingRule) getMedianRunningAverage(completedTrialslogs [][]*vdb.TrialLog, step int, burnin int) float64 {
	r := []float64{}
	var ra float64
	for _, ctl := range completedTrialslogs {
		ra = 0
		var st int
		var errParce bool = false
		if step > len(ctl) {
			st = len(ctl)
		} else {
			st = step
		}
		for s := burnin; s < st; s++ {
			v, err := strconv.ParseFloat(ctl[s].Value, 64)
			if err != nil {
				log.Printf("Fail to Parse %s : %s", ctl[s].Name, ctl[s].Value)
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
		log.Printf("running avg list %v", r)
		sort.Float64s(r)
		return r[len(r)/2]
	}
}

func (m *MedianStoppingRule) getBestValue(sid string, sc *api.StudyConfig, logs []*vdb.TrialLog) (float64, error) {
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
func (m *MedianStoppingRule) ShouldTrialStop(ctx context.Context, in *api.ShouldTrialStopRequest) (*api.ShouldTrialStopReply, error) {
	if _, ok := m.confList[in.StudyId]; !ok {
		return &api.ShouldTrialStopReply{}, errors.New("EarlyStopping config is not set.")
	}
	tl, err := m.dbIf.GetTrialList(in.StudyId)
	if err != nil {
		return &api.ShouldTrialStopReply{}, err
	}
	sc, err := m.dbIf.GetStudyConfig(in.StudyId)
	if err != nil {
		return &api.ShouldTrialStopReply{}, err
	}
	rtl := []*api.Trial{}
	ctl := make([][]*vdb.TrialLog, 0, len(tl))
	s_t := []*api.Trial{}
	for _, t := range tl {
		switch t.Status {
		case api.TrialState_RUNNING:
			rtl = append(rtl, t)
		case api.TrialState_COMPLETED:
			tl, err := m.dbIf.GetTrialLogs(t.TrialId, &vdb.GetTrialLogOpts{Name: m.confList[in.StudyId].EvalMetric})
			if err != nil {
				log.Printf("Fail to get trial %v logs", t.TrialId)
				continue
			}
			if len(tl) > m.confList[in.StudyId].BurnIn {
				ctl = append(ctl, tl)
			}
		default:
		}
	}
	if len(ctl) == 0 {
		return &api.ShouldTrialStopReply{}, nil
	}
	for _, t := range rtl {
		tl, err := m.dbIf.GetTrialLogs(t.TrialId, &vdb.GetTrialLogOpts{Name: m.confList[in.StudyId].EvalMetric})
		if err != nil {
			log.Printf("Fail to get trial %v logs", t.TrialId)
			continue
		}
		if len(tl) < m.confList[in.StudyId].LeastStep || len(tl) <= m.confList[in.StudyId].BurnIn {
			continue
		}
		v, err := m.getBestValue(in.StudyId, sc, tl)
		if err != nil {
			log.Printf("Fail to Get Best Value at %s: %v Log:%v", t.TrialId, err, tl)
			continue
		}
		om := m.getMedianRunningAverage(ctl, len(tl), m.confList[in.StudyId].BurnIn)
		log.Printf("Trial %s, In step %d Current value: %v Median value: %v\n", t.TrialId, len(tl), v, om)
		if (v < (om-m.confList[in.StudyId].Margin) && sc.OptimizationType == api.OptimizationType_MAXIMIZE) || v > (om+m.confList[in.StudyId].Margin) && sc.OptimizationType == api.OptimizationType_MINIMIZE {
			log.Print("Trial %s shuold be stopped", t.TrialId)
			s_t = append(s_t, t)
		}
	}
	return &api.ShouldTrialStopReply{Trials: s_t}, nil
}
