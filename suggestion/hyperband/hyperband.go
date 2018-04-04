package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/mlkube/katib/api"
	"github.com/mlkube/katib/suggestion"
	"log"
	"math"
	"sort"
	"strconv"
)

type Bracket []*api.Trial

func (b Bracket) Len() int {
	return len(b)
}

func (b Bracket) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Bracket) Less(i, j int) bool {
	vi, _ := strconv.ParseFloat(b[i].ObjectiveValue, 64)
	vj, _ := strconv.ParseFloat(b[j].ObjectiveValue, 64)
	return vi > vj
}

type HyperBandParameters struct {
	eta           float64
	sMax          int
	b_l           float64
	r_l           float64
	r             float64
	n             int
	shloopitr     int
	currentS      int
	MasterBracket Bracket
	ResourceName  string
}
type HyperBandSuggestService struct {
	suggestion.RandomSuggestService
	parameters map[string]*HyperBandParameters
}

func NewHyperBandSuggestService() *HyperBandSuggestService {
	return &HyperBandSuggestService{parameters: make(map[string]*HyperBandParameters)}
}

func (h *HyperBandSuggestService) generate_randid() string {
	// UUID isn't quite handy in the Go world
	id_ := make([]byte, 8)
	_, err := rand.Read(id_)
	if err != nil {
		log.Fatalf("Error reading random: %v", err)
	}
	return fmt.Sprintf("%016x", id_)
}

func (h *HyperBandSuggestService) makeMasterBracket(sconf *api.StudyConfig, n int) Bracket {
	log.Printf("Make MasterBracket %v Trials", n)
	s_t := make([]*api.Trial, n)
	for i := 0; i < n; i++ {
		s_t[i] = &api.Trial{}
		s_t[i].ParameterSet = make([]*api.Parameter, len(sconf.ParameterConfigs.Configs))
		for j, pc := range sconf.ParameterConfigs.Configs {
			s_t[i].ParameterSet[j] = &api.Parameter{Name: pc.Name}
			switch pc.ParameterType {
			case api.ParameterType_INT:
				imin, _ := strconv.Atoi(pc.Feasible.Min)
				imax, _ := strconv.Atoi(pc.Feasible.Max)
				s_t[i].ParameterSet[j].Value = strconv.Itoa(h.IntRandom(imin, imax))
			case api.ParameterType_DOUBLE:
				dmin, _ := strconv.ParseFloat(pc.Feasible.Min, 64)
				dmax, _ := strconv.ParseFloat(pc.Feasible.Max, 64)
				s_t[i].ParameterSet[j].Value = strconv.FormatFloat(h.DoubelRandom(dmin, dmax), 'f', 4, 64)
			case api.ParameterType_CATEGORICAL:
				s_t[i].ParameterSet[j].Value = pc.Feasible.List[h.IntRandom(0, len(pc.Feasible.List)-1)]
			}
		}
		s_t[i].Tags = append(s_t[i].Tags, &api.Tag{Name: "HyperBand_BracketID", Value: h.generate_randid()})
	}
	return Bracket(s_t)
}

func (h *HyperBandSuggestService) SetSuggestionParameters(ctx context.Context, in *api.SetSuggestionParametersRequest) (*api.SetSuggestionParametersReply, error) {
	p := &HyperBandParameters{}
	for _, sp := range in.SuggestionParameters {
		switch sp.Name {
		case "Eta":
			p.eta, _ = strconv.ParseFloat(sp.Value, 64)
		case "R":
			p.r_l, _ = strconv.ParseFloat(sp.Value, 64)
		case "ResourceName":
			p.ResourceName = sp.Value
		default:
			log.Printf("Unknown Suggestion Parameter %v", sp.Name)
		}
	}
	if p.eta == 0 || p.r_l == 0 || p.ResourceName == "" {
		log.Printf("Failed to Suggestion Parameter set.")
		return &api.SetSuggestionParametersReply{}, fmt.Errorf("Suggestion Parameter set Error")
	}
	p.sMax = int(math.Log(p.r_l) / math.Log(p.eta))
	p.b_l = float64((p.sMax + 1.0)) * p.r_l
	p.n = int((p.b_l/p.r_l)*(math.Pow(p.eta, float64(p.sMax))/float64(p.sMax+1))) + 1
	p.currentS = p.sMax + 1
	p.shloopitr = p.currentS + 1
	p.r = p.r_l * math.Pow(p.eta, float64(-p.sMax))
	p.MasterBracket = h.makeMasterBracket(in.Configs, p.n)
	h.parameters[in.StudyId] = p
	log.Printf("Smax = %v", p.sMax)
	return &api.SetSuggestionParametersReply{}, nil
}

func (h *HyperBandSuggestService) getHyperParameter(studyId string, sconf *api.StudyConfig, n int) Bracket {
	s_t := make([]*api.Trial, n)
	for i := 0; i < n; i++ {
		s_t[i] = &api.Trial{}
		s_t[i].ParameterSet = make([]*api.Parameter, len(sconf.ParameterConfigs.Configs))
		s_t[i].Status = api.TrialState_PENDING
		s_t[i].EvalLogs = make([]*api.EvaluationLog, 0)
		var j int
		if sconf.OptimizationType == api.OptimizationType_MAXIMIZE {
			j = i
		} else if sconf.OptimizationType == api.OptimizationType_MINIMIZE {
			j = len(h.parameters[studyId].MasterBracket) - 1 - i
		}
		for k, v := range h.parameters[studyId].MasterBracket[j].ParameterSet {
			s_t[i].ParameterSet[k] = v
		}
		for _, t := range h.parameters[studyId].MasterBracket[j].Tags {
			s_t[i].Tags = append(s_t[i].Tags, t)
		}
	}

	return Bracket(s_t)
}
func (h *HyperBandSuggestService) hbLoopParamUpdate(studyId string) {
	log.Printf("HB loop s = %v", h.parameters[studyId].currentS)
	h.parameters[studyId].shloopitr = 0
	h.parameters[studyId].n = int((h.parameters[studyId].b_l/h.parameters[studyId].r_l)*(math.Pow(h.parameters[studyId].eta, float64(h.parameters[studyId].currentS))/float64(h.parameters[studyId].currentS+1))) + 1
	h.parameters[studyId].r = h.parameters[studyId].r_l * math.Pow(h.parameters[studyId].eta, float64(-h.parameters[studyId].currentS))
}

func (h *HyperBandSuggestService) shLoopParamUpdate(studyId string) (int, int) {
	log.Printf("SH loop i = %v", h.parameters[studyId].shloopitr)
	pn_i := int(float64(h.parameters[studyId].n) * math.Pow(h.parameters[studyId].eta, float64(-h.parameters[studyId].shloopitr+1)) / h.parameters[studyId].eta)
	r_i := int(h.parameters[studyId].r * math.Pow(h.parameters[studyId].eta, float64(h.parameters[studyId].shloopitr)))
	return pn_i, r_i
}

func (h *HyperBandSuggestService) GenerateTrials(ctx context.Context, in *api.GenerateTrialsRequest) (*api.GenerateTrialsReply, error) {
	if h.parameters[in.StudyId].currentS <= 0 {
		h.StopSuggestion(ctx, &api.StopSuggestionRequest{StudyId: in.StudyId})
		return &api.GenerateTrialsReply{Completed: true}, nil
	}
	if len(in.RunningTrials) > 0 {
		return &api.GenerateTrialsReply{Completed: false}, nil
	}
	if len(in.CompletedTrials) > 0 {
		var schec int
		var bid string
		for _, c := range in.CompletedTrials {
			schec = 0
			for _, t := range c.Tags {
				if t.Name == "HyperBand_shi" && t.Value == strconv.Itoa(h.parameters[in.StudyId].shloopitr) {
					schec++
				}
				if t.Name == "HyperBand_s" && t.Value == strconv.Itoa(h.parameters[in.StudyId].currentS) {
					schec++
				}
				if t.Name == "HyperBand_BracketID" {
					bid = t.Value
				}
			}
			if schec == 2 {
				for _, b := range h.parameters[in.StudyId].MasterBracket {
					for _, t := range b.Tags {
						if t.Name == "HyperBand_BracketID" && t.Value == bid {
							b.ObjectiveValue = c.ObjectiveValue
						}
					}
				}
			}
		}
		sort.Sort(h.parameters[in.StudyId].MasterBracket)
	}
	var evalT []*api.Trial
	var r_i int
	if h.parameters[in.StudyId].shloopitr > h.parameters[in.StudyId].currentS {
		h.parameters[in.StudyId].currentS--
		h.hbLoopParamUpdate(in.StudyId)
		_, r_i = h.shLoopParamUpdate(in.StudyId)
		h.parameters[in.StudyId].MasterBracket = h.makeMasterBracket(in.Configs, h.parameters[in.StudyId].n)
		evalT = h.getHyperParameter(in.StudyId, in.Configs, h.parameters[in.StudyId].n)
		h.parameters[in.StudyId].shloopitr++
	} else {
		var pn_i int
		pn_i, r_i = h.shLoopParamUpdate(in.StudyId)
		evalT = h.getHyperParameter(in.StudyId, in.Configs, pn_i)
		h.parameters[in.StudyId].shloopitr++
	}
	for i := range evalT {
		for j := range evalT[i].ParameterSet {
			if evalT[i].ParameterSet[j].Name == h.parameters[in.StudyId].ResourceName {
				evalT[i].ParameterSet[j].Value = strconv.Itoa(r_i)
			}
		}
		evalT[i].Tags = append(evalT[i].Tags, &api.Tag{Name: "HyperBand_s", Value: strconv.Itoa(h.parameters[in.StudyId].currentS)})
		evalT[i].Tags = append(evalT[i].Tags, &api.Tag{Name: "HyperBand_r", Value: strconv.Itoa(r_i)})
		evalT[i].Tags = append(evalT[i].Tags, &api.Tag{Name: "HyperBand_shi", Value: strconv.Itoa(h.parameters[in.StudyId].shloopitr)})
		log.Printf("Gen Trial %v", evalT[i].Tags)
	}
	return &api.GenerateTrialsReply{Trials: evalT, Completed: false}, nil
}

func (h *HyperBandSuggestService) StopSuggestion(ctx context.Context, in *api.StopSuggestionRequest) (*api.StopSuggestionReply, error) {
	delete(h.parameters, in.StudyId)
	return &api.StopSuggestionReply{}, nil
}
