package suggestion

import (
	"context"
	"github.com/mlkube/katib/api"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type SuggestService interface {
	SetSuggestionParameters(ctx context.Context, in *api.SetSuggestionParametersRequest) (*api.SetSuggestionParametersReply, error)
	GenerateTrials(ctx context.Context, in *api.GenerateTrialsRequest) (*api.GenerateTrialsReply, error)
}

type RandomSuggestParameters struct {
	SuggestionNum int
	MaxParallel   int
}
type RandomSuggestService struct {
	parameters map[string]*RandomSuggestParameters
}

func NewRandomSuggestService() *RandomSuggestService {
	return &RandomSuggestService{parameters: make(map[string]*RandomSuggestParameters)}
}

func (s *RandomSuggestService) DoubelRandom(min, max float64) float64 {
	if min == max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

func (s *RandomSuggestService) IntRandom(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func (s *RandomSuggestService) SetSuggestionParameters(ctx context.Context, in *api.SetSuggestionParametersRequest) (*api.SetSuggestionParametersReply, error) {
	p := &RandomSuggestParameters{}
	for _, sp := range in.SuggestionParameters {
		switch sp.Name {
		case "SuggestionNum":
			p.SuggestionNum, _ = strconv.Atoi(sp.Value)
		case "MaxParallel":
			p.MaxParallel, _ = strconv.Atoi(sp.Value)
		default:
			log.Printf("Unknown Suggestion Parameter %v", sp.Name)
		}
	}
	s.parameters[in.StudyId] = p
	return &api.SetSuggestionParametersReply{}, nil
}

func (s *RandomSuggestService) GenerateTrials(ctx context.Context, in *api.GenerateTrialsRequest) (*api.GenerateTrialsReply, error) {
	if len(in.CompletedTrials) >= s.parameters[in.StudyId].SuggestionNum {
		s.StopSuggestion(ctx, &api.StopSuggestionRequest{StudyId: in.StudyId})
		return &api.GenerateTrialsReply{Completed: true}, nil
	}
	if s.parameters[in.StudyId].MaxParallel < 1 && len(in.RunningTrials) > 0 {
		return &api.GenerateTrialsReply{Completed: false}, nil
	} else {
		if len(in.RunningTrials) >= s.parameters[in.StudyId].MaxParallel {
			return &api.GenerateTrialsReply{Completed: false}, nil
		}
		if s.parameters[in.StudyId].SuggestionNum-len(in.CompletedTrials)-len(in.RunningTrials) <= 0 {
			return &api.GenerateTrialsReply{Completed: false}, nil
		}
	}
	var reqnum int = 0
	if s.parameters[in.StudyId].MaxParallel < 1 {
		reqnum = s.parameters[in.StudyId].SuggestionNum
	} else if s.parameters[in.StudyId].SuggestionNum-len(in.CompletedTrials) <= s.parameters[in.StudyId].MaxParallel {
		reqnum = s.parameters[in.StudyId].SuggestionNum - len(in.CompletedTrials) - len(in.RunningTrials)
	} else {
		reqnum = s.parameters[in.StudyId].MaxParallel - len(in.RunningTrials)
	}
	s_t := make([]*api.Trial, reqnum)
	for i := 0; i < reqnum; i++ {
		s_t[i] = &api.Trial{}
		s_t[i].ParameterSet = make([]*api.Parameter, len(in.Configs.ParameterConfigs.Configs))
		s_t[i].Status = api.TrialState_PENDING
		s_t[i].EvalLogs = make([]*api.EvaluationLog, 0)
		for j, pc := range in.Configs.ParameterConfigs.Configs {
			s_t[i].ParameterSet[j] = &api.Parameter{Name: pc.Name}
			switch pc.ParameterType {
			case api.ParameterType_INT:
				imin, _ := strconv.Atoi(pc.Feasible.Min)
				imax, _ := strconv.Atoi(pc.Feasible.Max)
				s_t[i].ParameterSet[j].Value = strconv.Itoa(s.IntRandom(imin, imax))
			case api.ParameterType_DOUBLE:
				dmin, _ := strconv.ParseFloat(pc.Feasible.Min, 64)
				dmax, _ := strconv.ParseFloat(pc.Feasible.Max, 64)
				s_t[i].ParameterSet[j].Value = strconv.FormatFloat(s.DoubelRandom(dmin, dmax), 'f', 4, 64)
			case api.ParameterType_CATEGORICAL:
				s_t[i].ParameterSet[j].Value = pc.Feasible.List[s.IntRandom(0, len(pc.Feasible.List)-1)]
			}
		}
	}
	return &api.GenerateTrialsReply{Trials: s_t, Completed: false}, nil
}

func (s *RandomSuggestService) StopSuggestion(ctx context.Context, in *api.StopSuggestionRequest) (*api.StopSuggestionReply, error) {
	delete(s.parameters, in.StudyId)
	return &api.StopSuggestionReply{}, nil
}
