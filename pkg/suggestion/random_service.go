package suggestion

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/kubeflow/katib/pkg/api"
	"google.golang.org/grpc"
)

type RandomSuggestService struct {
}

func NewRandomSuggestService() *RandomSuggestService {
	return &RandomSuggestService{}
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

func (s *RandomSuggestService) GetSuggestions(ctx context.Context, in *api.GetSuggestionsRequest) (*api.GetSuggestionsReply, error) {
	conn, err := grpc.Dial(manager, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return &api.GetSuggestionsReply{}, err
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)
	screq := &api.GetStudyRequest{
		StudyId: in.StudyId,
	}
	scr, err := c.GetStudy(ctx, screq)
	if err != nil {
		log.Fatalf("GetStudyConf failed: %v", err)
		return &api.GetSuggestionsReply{}, err
	}
	reqnum := int(in.RequestNumber)
	s_t := make([]*api.Trial, reqnum)
	for i := 0; i < reqnum; i++ {
		s_t[i] = &api.Trial{}
		s_t[i].StudyId = in.StudyId
		s_t[i].ParameterSet = make([]*api.Parameter, len(scr.StudyConfig.ParameterConfigs.Configs))
		s_t[i].Status = api.State_PENDING
		for j, pc := range scr.StudyConfig.ParameterConfigs.Configs {
			s_t[i].ParameterSet[j] = &api.Parameter{Name: pc.Name}
			s_t[i].ParameterSet[j].ParameterType = pc.ParameterType
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
		ctreq := &api.CreateTrialRequest{
			Trial: s_t[i],
		}
		ctret, err := c.CreateTrial(ctx, ctreq)
		if err != nil {
			return &api.GetSuggestionsReply{Trials: []*api.Trial{}}, err
		}
		s_t[i].TrialId = ctret.TrialId
	}
	return &api.GetSuggestionsReply{Trials: s_t}, nil
}
