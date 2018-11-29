package suggestion

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kubeflow/katib/pkg"
	"github.com/kubeflow/katib/pkg/api"
	"google.golang.org/grpc"
)

type ManualSuggestService struct {
}

func NewManualSuggestService() *ManualSuggestService {
	return &ManualSuggestService{}
}

func (s *ManualSuggestService) ValidateParameter(pvalue string, pc *api.ParameterConfig) error {
	switch pc.ParameterType {
	case api.ParameterType_INT:
		imin, _ := strconv.Atoi(pc.Feasible.Min)
		imax, _ := strconv.Atoi(pc.Feasible.Max)
		ivalue, err := strconv.Atoi(pvalue)
		if err != nil {
			return err
		}
		if ivalue < imin || ivalue > imax {
			return fmt.Errorf("%s :Parameter outof range", pc.Name)
		}
	case api.ParameterType_DOUBLE:
		dmin, _ := strconv.ParseFloat(pc.Feasible.Min, 64)
		dmax, _ := strconv.ParseFloat(pc.Feasible.Max, 64)
		dvalue, err := strconv.ParseFloat(pvalue, 64)
		if err != nil {
			return err
		}
		if dvalue < dmin || dvalue > dmax {
			return fmt.Errorf("%s :Parameter outof range", pc.Name)
		}
	case api.ParameterType_CATEGORICAL:
		ok := false
		for _, l := range pc.Feasible.List {
			if l == pvalue {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("%s :Parameter outof range", pc.Name)
		}
	case api.ParameterType_DISCRETE:
		ok := false
		for _, l := range pc.Feasible.List {
			if l == pvalue {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("%s :Parameter outof range", pc.Name)
		}
	}
	return nil
}

func (s *ManualSuggestService) ParseParameters(ctx context.Context, c api.ManagerClient, studyId string, sconf *api.StudyConfig, params []*api.SuggestionParameter) (*api.Trial, error) {
	mparam := map[string]string{}
	var tid string
	for _, param := range params {
		if strings.ToLower(param.Name) == "trialid" {
			tid = param.Value
		} else if param.Name == "SuggestionCount" {
			scount, _ := strconv.Atoi(param.Value)
			if scount > 0 {
				return nil, nil
			}
		} else {
			mparam[param.Name] = param.Value
		}
	}
	if tid != "" {
		gtreq := &api.GetTrialRequest{
			TrialId: tid,
		}
		gtrep, err := c.GetTrial(ctx, gtreq)
		if err != nil {
			return nil, err
		}
		return gtrep.Trial, nil
	}
	trial := &api.Trial{
		StudyId:      studyId,
		ParameterSet: make([]*api.Parameter, len(sconf.ParameterConfigs.Configs)),
	}
	for i, pc := range sconf.ParameterConfigs.Configs {
		if pvalue, ok := mparam[pc.Name]; ok {
			err := s.ValidateParameter(pvalue, pc)
			if err != nil {
				return nil, err
			}
			trial.ParameterSet[i] = &api.Parameter{}
			trial.ParameterSet[i].Name = pc.Name
			trial.ParameterSet[i].Value = pvalue
		} else {
			return nil, fmt.Errorf("%s :Parameter is not includes", pc.Name)
		}
	}
	ctreq := &api.CreateTrialRequest{
		Trial: trial,
	}
	ctret, err := c.CreateTrial(ctx, ctreq)
	if err != nil {
		return nil, err
	}
	trial.TrialId = ctret.TrialId
	return trial, nil
}

func (s *ManualSuggestService) GetSuggestions(ctx context.Context, in *api.GetSuggestionsRequest) (*api.GetSuggestionsReply, error) {
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
	spreq := &api.GetSuggestionParametersRequest{
		ParamId: in.ParamId,
	}
	spr, err := c.GetSuggestionParameters(ctx, spreq)
	if err != nil {
		log.Printf("GetParameter failed: %v", err)
		return &api.GetSuggestionsReply{}, err
	}
	t, err := s.ParseParameters(ctx, c, in.StudyId, scr.StudyConfig, spr.SuggestionParameters)
	if err != nil {
		log.Printf("Parse Parameter failed: %v", err)
		return &api.GetSuggestionsReply{}, err
	}
	sT := []*api.Trial{}
	if t != nil {
		sT = append(sT, t)
	}
	return &api.GetSuggestionsReply{Trials: sT}, nil
}
