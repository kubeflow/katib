package suggestion

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/kubeflow/katib/pkg"
	"github.com/kubeflow/katib/pkg/api"
	"google.golang.org/grpc"
)

type RandomSuggestService struct {
}

func NewRandomSuggestService() *RandomSuggestService {
	return &RandomSuggestService{}
}

func (s *RandomSuggestService) DoubleRandom(min, max float64) float64 {
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
	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
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
	sT := make([]*api.Trial, reqnum)
	for i := 0; i < reqnum; i++ {
		sT[i] = &api.Trial{}
		sT[i].StudyId = in.StudyId
		sT[i].ParameterSet = make([]*api.Parameter, len(scr.StudyConfig.ParameterConfigs.Configs))
		for j, pc := range scr.StudyConfig.ParameterConfigs.Configs {
			sT[i].ParameterSet[j] = &api.Parameter{Name: pc.Name}
			sT[i].ParameterSet[j].ParameterType = pc.ParameterType
			switch pc.ParameterType {
			case api.ParameterType_INT:
				imin, _ := strconv.Atoi(pc.Feasible.Min)
				imax, _ := strconv.Atoi(pc.Feasible.Max)
				sT[i].ParameterSet[j].Value = strconv.Itoa(s.IntRandom(imin, imax))
			case api.ParameterType_DOUBLE:
				dmin, _ := strconv.ParseFloat(pc.Feasible.Min, 64)
				dmax, _ := strconv.ParseFloat(pc.Feasible.Max, 64)
				sT[i].ParameterSet[j].Value = strconv.FormatFloat(s.DoubleRandom(dmin, dmax), 'f', 4, 64)
			case api.ParameterType_CATEGORICAL:
				sT[i].ParameterSet[j].Value = pc.Feasible.List[s.IntRandom(0, len(pc.Feasible.List)-1)]
			}
		}
		ctreq := &api.CreateTrialRequest{
			Trial: sT[i],
		}
		ctret, err := c.CreateTrial(ctx, ctreq)
		if err != nil {
			return &api.GetSuggestionsReply{Trials: []*api.Trial{}}, err
		}
		sT[i].TrialId = ctret.TrialId
	}
	return &api.GetSuggestionsReply{Trials: sT}, nil
}

func (s *RandomSuggestService) ValidateSuggestionParameters(ctx context.Context, in *api.ValidateSuggestionParametersRequest) (*api.ValidateSuggestionParametersReply, error) {

	return &api.ValidateSuggestionParametersReply{}, nil
}