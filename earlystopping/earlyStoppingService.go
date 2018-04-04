package earlystopping

import (
	//	"fmt"
	"github.com/mlkube/katib/api"
	"sort"
	//	"strconv"
)

type EarlyStoppingService interface {
	ShouldStoppingTrial(runningTrials []*api.Trial, completedTrials []*api.Trial, leastStep int) []api.Trial
}

type MedianStoppingRule struct{}

func NewMedianStoppingRule() *MedianStoppingRule {
	m := &MedianStoppingRule{}
	return m
}

func (m *MedianStoppingRule) getMedianRunningAverage(completedTrials []api.Trial, step int) float64 {
	r := []float64{}
	for _, ct := range completedTrials {
		if ct.Status == api.TrialState_COMPLETED {
			//			var ra float64
			//			for s := 0; s < step; s++ {
			//				p, _ := strconv.ParseFloat(ct.EvalLogs[s].Value, 64)
			//				ra += p
			//			}
			//			ra = ra / float64(len(ct.EvalLogs))
			//			r = append(r, ra)
			//				var p float64
			//				if len(ct.EvalLogs) < step {
			//					p, _ = strconv.ParseFloat(ct.EvalLogs[len(ct.EvalLogs)-1].Value, 64)
			//				} else {
			//					p, _ = strconv.ParseFloat(ct.EvalLogs[step-1].Value, 64)
			//				}
			//				r = append(r, p)
		}
	}
	if len(r) == 0 {
		return 0
	} else {
		sort.Float64s(r)
		return r[len(r)/2]
	}
}

func (m *MedianStoppingRule) ShouldStoppingTrial(runningTrials []api.Trial, completedTrials []api.Trial, leastStep int) []api.Trial {
	s_t := []api.Trial{}
	for _, t := range runningTrials {
		if t.Status != api.TrialState_RUNNING {
			continue
		}
		s := len(t.EvalLogs)
		if s < leastStep {
			continue
		}
		//om := m.getMedianRunningAverage(completedTrials, s)
		//		v, _ := strconv.ParseFloat(t.EvalLogs[s-1].Value, 64)
		//		fmt.Printf("Trial %v, Current value %v Median value in step %v %v\n", t.TrialId, v, s, om)
		//		if v < (om - 0.03) {
		//			s_t = append(s_t, t)
		//		}
	}
	return s_t
}
