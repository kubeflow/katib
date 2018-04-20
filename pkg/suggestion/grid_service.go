package suggestion

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/kubeflow/katib/pkg/api"
)

type GridSuggestParameters struct {
	defaultGridNum int
	gridConfig     map[string]int
	MaxParallel    int
}

type GridSuggestService struct {
	parameters  map[string]*GridSuggestParameters
	grids       map[string][][]*api.Parameter
	gridPointer map[string]int
}

func NewGridSuggestService() *GridSuggestService {
	return &GridSuggestService{parameters: make(map[string]*GridSuggestParameters), grids: make(map[string][][]*api.Parameter), gridPointer: make(map[string]int)}
}

func (s *GridSuggestService) SetSuggestionParameters(ctx context.Context, in *api.SetSuggestionParametersRequest) (*api.SetSuggestionParametersReply, error) {
	p := &GridSuggestParameters{gridConfig: make(map[string]int)}
	for _, sp := range in.SuggestionParameters {
		switch sp.Name {
		case "DefaultGrid":
			p.defaultGridNum, _ = strconv.Atoi(sp.Value)
		case "MaxParallel":
			p.MaxParallel, _ = strconv.Atoi(sp.Value)
		default:
			p.gridConfig[sp.Name], _ = strconv.Atoi(sp.Value)
		}
	}
	s.parameters[in.StudyId] = p
	return &api.SetSuggestionParametersReply{}, nil
}

func (s *GridSuggestService) allocInt(min int, max int, reqnum int) []string {
	ret := make([]string, reqnum)
	if reqnum == 1 {
		ret[0] = strconv.Itoa(min)
	} else {
		for i := 0; i < reqnum; i++ {
			ret[i] = strconv.Itoa(min + ((max - min) * i / (reqnum - 1)))
		}
	}
	return ret
}

func (s *GridSuggestService) allocFloat(min float64, max float64, reqnum int) []string {
	ret := make([]string, reqnum)
	if reqnum == 1 {
		ret[0] = strconv.FormatFloat(min, 'f', 4, 64)
	} else {
		for i := 0; i < reqnum; i++ {
			ret[i] = strconv.FormatFloat(min+(((max-min)/float64(reqnum-1))*float64(i)), 'f', 4, 64)
		}
	}
	return ret
}

func (s *GridSuggestService) allocCat(list []string, reqnum int) []string {
	ret := make([]string, reqnum)
	if reqnum == 1 {
		ret[0] = list[0]
	} else {
		for i := 0; i < reqnum; i++ {
			ret[i] = list[(((len(list) - 1) * i) / (reqnum - 1))]
			fmt.Printf("ret %v %v\n", i, ret[i])
		}
	}
	return ret
}

func (s *GridSuggestService) setP(gci int, p [][]*api.Parameter, pg [][]string, pcs []*api.ParameterConfig) {
	if gci == len(pg)-1 {
		for i := range pg[gci] {
			p[i] = append(p[i], &api.Parameter{
				Name:          pcs[gci].Name,
				ParameterType: pcs[gci].ParameterType,
				Value:         pg[gci][i],
			})

		}
		return
	} else {
		d := len(p) / len(pg[gci])
		for i := range pg[gci] {
			for j := d * i; j < d*(i+1); j++ {
				p[j] = append(p[j], &api.Parameter{
					Name:          pcs[gci].Name,
					ParameterType: pcs[gci].ParameterType,
					Value:         pg[gci][i],
				})
			}
			s.setP(gci+1, p[d*i:d*(i+1)], pg, pcs)
		}
	}
}

func (s *GridSuggestService) genGrids(studyId string, pcs []*api.ParameterConfig) [][]*api.Parameter {
	var pg [][]string
	var holenum = 1
	gcl := make([]int, len(pcs))
	for i, pc := range pcs {
		gc, ok := s.parameters[studyId].gridConfig[pc.Name]
		if !ok {
			gc = s.parameters[studyId].defaultGridNum
		}
		holenum *= gc
		gcl[i] = gc
		switch pc.ParameterType {
		case api.ParameterType_INT:
			imin, _ := strconv.Atoi(pc.Feasible.Min)
			imax, _ := strconv.Atoi(pc.Feasible.Max)
			pg = append(pg, s.allocInt(imin, imax, gc))
		case api.ParameterType_DOUBLE:
			dmin, _ := strconv.ParseFloat(pc.Feasible.Min, 64)
			dmax, _ := strconv.ParseFloat(pc.Feasible.Max, 64)
			pg = append(pg, s.allocFloat(dmin, dmax, gc))
		case api.ParameterType_CATEGORICAL:
			pg = append(pg, s.allocCat(pc.Feasible.List, gc))
		}
	}
	ret := make([][]*api.Parameter, holenum)
	s.setP(0, ret, pg, pcs)
	log.Printf("Study %v : %v parameters generated", studyId, holenum)
	return ret
}

func (s *GridSuggestService) GenerateTrials(ctx context.Context, in *api.GenerateTrialsRequest) (*api.GenerateTrialsReply, error) {
	if _, ok := s.grids[in.StudyId]; !ok {
		s.grids[in.StudyId] = s.genGrids(in.StudyId, in.Configs.ParameterConfigs.Configs)
		s.gridPointer[in.StudyId] = 0
	}
	if s.gridPointer[in.StudyId] >= len(s.grids[in.StudyId]) {
		if len(in.RunningTrials) == 0 {
			s.StopSuggestion(ctx, &api.StopSuggestionRequest{StudyId: in.StudyId})
			return &api.GenerateTrialsReply{Completed: true}, nil
		} else {
			return &api.GenerateTrialsReply{Completed: false}, nil
		}
	}
	var reqnum = 0
	if s.parameters[in.StudyId].MaxParallel <= 0 {
		reqnum = len(s.grids[in.StudyId])
	} else if len(s.grids[in.StudyId])-s.gridPointer[in.StudyId] < s.parameters[in.StudyId].MaxParallel-len(in.RunningTrials) {
		reqnum = len(s.grids[in.StudyId]) - s.gridPointer[in.StudyId]
	} else if len(in.RunningTrials) < s.parameters[in.StudyId].MaxParallel {
		reqnum = s.parameters[in.StudyId].MaxParallel - len(in.RunningTrials)
	}
	trials := make([]*api.Trial, reqnum)
	for i := 0; i < int(reqnum); i++ {
		trials[i] = &api.Trial{}
		trials[i].Status = api.TrialState_PENDING
		trials[i].EvalLogs = make([]*api.EvaluationLog, 0)
		trials[i].ParameterSet = s.grids[in.StudyId][s.gridPointer[in.StudyId]+i]
	}
	s.gridPointer[in.StudyId] += reqnum
	return &api.GenerateTrialsReply{Trials: trials, Completed: false}, nil
}

func (s *GridSuggestService) StopSuggestion(ctx context.Context, in *api.StopSuggestionRequest) (*api.StopSuggestionReply, error) {
	delete(s.gridPointer, in.StudyId)
	delete(s.grids, in.StudyId)
	delete(s.parameters, in.StudyId)
	return &api.StopSuggestionReply{}, nil
}
