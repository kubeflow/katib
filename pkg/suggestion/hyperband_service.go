package suggestion

import (
	"context"
	//	"crypto/rand"
	//	"fmt"
	//	"github.com/kubeflow/katib/pkg/db"
	//	"log"
	//	"math"
	//	"sort"
	//	"strconv"

	"github.com/kubeflow/katib/pkg/api"
)

//type Evals struct {
//	id    string
//	value float64
//}
//type Bracket []Evals
//
//func (b Bracket) Len() int {
//	return len(b)
//}
//
//func (b Bracket) Swap(i, j int) {
//	b[i], b[j] = b[j], b[i]
//}
//
//func (b Bracket) Less(i, j int) bool {
//	return b[i].value > b[j].value
//}
//
type HyperBandParameters struct {
	eta          float64
	sMax         int
	b_l          float64
	r_l          float64
	r            float64
	n            int
	shloopitr    int
	currentS     int
	ResourceName string
}

type HyperBandSuggestService struct {
	RandomSuggestService
	parameters HyperBandParameters
}

func NewHyperBandSuggestService() *HyperBandSuggestService {
	return &HyperBandSuggestService{}
}

//
//func (h *HyperBandSuggestService) generate_randid() string {
//	// UUID isn't quite handy in the Go world
//	id_ := make([]byte, 8)
//	_, err := rand.Read(id_)
//	if err != nil {
//		log.Fatalf("Error reading random: %v", err)
//	}
//	return fmt.Sprintf("%016x", id_)
//}
//
//func (h *HyperBandSuggestService) makeMasterBracket(sconf *api.StudyConfig, n int) Bracket {
//	log.Printf("Make MasterBracket %v Trials", n)
//	s_t := make([]*api.Trial, n)
//	for i := 0; i < n; i++ {
//		s_t[i] = &api.Trial{}
//		s_t[i].ParameterSet = make([]*api.Parameter, len(sconf.ParameterConfigs.Configs))
//		for j, pc := range sconf.ParameterConfigs.Configs {
//			s_t[i].ParameterSet[j] = &api.Parameter{Name: pc.Name}
//			s_t[i].ParameterSet[j].ParameterType = pc.ParameterType
//			switch pc.ParameterType {
//			case api.ParameterType_INT:
//				imin, _ := strconv.Atoi(pc.Feasible.Min)
//				imax, _ := strconv.Atoi(pc.Feasible.Max)
//				s_t[i].ParameterSet[j].Value = strconv.Itoa(h.IntRandom(imin, imax))
//			case api.ParameterType_DOUBLE:
//				dmin, _ := strconv.ParseFloat(pc.Feasible.Min, 64)
//				dmax, _ := strconv.ParseFloat(pc.Feasible.Max, 64)
//				s_t[i].ParameterSet[j].Value = strconv.FormatFloat(h.DoubelRandom(dmin, dmax), 'f', 4, 64)
//			case api.ParameterType_CATEGORICAL:
//				s_t[i].ParameterSet[j].Value = pc.Feasible.List[h.IntRandom(0, len(pc.Feasible.List)-1)]
//			}
//		}
//		s_t[i].Tags = append(s_t[i].Tags, &api.Tag{Name: "HyperBand_BracketID", Value: h.generate_randid()})
//	}
//	return Bracket(s_t)
//}
//
//func (h *HyperBandSuggestService) purseSuggestionParameters(sparam []*api.SuggestionParameters) (HyperBandParameters, error) {
//	p := &HyperBandParameters{
//		eta:          -1,
//		sMax:         -1,
//		b_l:          -1,
//		r_l:          -1,
//		r:            -1,
//		n:            -1,
//		shloopitr:    -1,
//		currentS:     -1,
//		ResourceName: -1,
//	}
//	for _, sp := range sparam {
//		switch sp.Name {
//		case "Eta":
//			p.eta, _ = strconv.ParseFloat(sp.Value, 64)
//		case "R":
//			p.r_l, _ = strconv.ParseFloat(sp.Value, 64)
//		case "ResourceName":
//			p.ResourceName = sp.Value
//		case "b_l":
//			p.b_l, _ = strconv.ParseFloat(sp.Value, 64)
//		case "sMax":
//			p.sMax, _ = strconv.AtoI(sp.Value)
//		case "r_s":
//			p.r, _ = strconv.ParseFloat(sp.Value, 64)
//		case "n_s":
//			p.n, _ = strconv.AtoI(sp.Value)
//		case "shloopitr":
//			p.shloopitr, _ = strconv.AtoI(sp.Value)
//		case "currentS":
//			p.currentS, _ = strconv.AtoI(sp.Value)
//		default:
//			log.Printf("Unknown Suggestion Parameter %v", sp.Name)
//		}
//	}
//	if p.eta == 0 || p.r_l == 0 || p.ResourceName == "" {
//		log.Printf("Failed to Suggestion Parameter set.")
//		return &api.SetSuggestionParametersReply{}, fmt.Errorf("Suggestion Parameter set Error")
//	}
//	if p.sMax == -1 {
//		p.sMax = int(math.Log(p.r_l) / math.Log(p.eta))
//	}
//	if p.b_l == -1 {
//		p.b_l = float64((p.sMax + 1.0)) * p.r_l
//	}
//	if p.n == -1 {
//		p.n = int((p.b_l/p.r_l)*(math.Pow(p.eta, float64(p.sMax))/float64(p.sMax+1))) + 1
//	}
//	if p.currentS == -1 {
//		p.currentS = p.sMax + 1
//	}
//	if p.shloopit == -1 {
//		p.shloopitr = p.currentS + 1
//	}
//	if p.r == -1 {
//		p.r = p.r_l * math.Pow(p.eta, float64(-p.sMax))
//	}
//	p.MasterBracket = h.makeMasterBracket(in.Configs, p.n)
//	h.parameters[in.StudyId] = p
//	log.Printf("Smax = %v", p.sMax)
//	return &api.SetSuggestionParametersReply{}, nil
//}
//
//func (h *HyperBandSuggestService) getHyperParameter(studyId string, sconf *api.StudyConfig, n int) Bracket {
//	s_t := make([]*api.Trial, n)
//	for i := 0; i < n; i++ {
//		s_t[i] = &api.Trial{}
//		s_t[i].ParameterSet = make([]*api.Parameter, len(sconf.ParameterConfigs.Configs))
//		s_t[i].Status = api.TrialState_PENDING
//		s_t[i].EvalLogs = make([]*api.EvaluationLog, 0)
//		var j int
//		if sconf.OptimizationType == api.OptimizationType_MAXIMIZE {
//			j = i
//		} else if sconf.OptimizationType == api.OptimizationType_MINIMIZE {
//			j = len(h.parameters[studyId].MasterBracket) - 1 - i
//		}
//		for k, v := range h.parameters[studyId].MasterBracket[j].ParameterSet {
//			s_t[i].ParameterSet[k] = v
//		}
//		for _, t := range h.parameters[studyId].MasterBracket[j].Tags {
//			s_t[i].Tags = append(s_t[i].Tags, t)
//		}
//	}
//	return Bracket(s_t)
//}
//
//func (h *HyperBandSuggestService) hbLoopParamUpdate(studyId string) {
//	log.Printf("HB loop s = %v", h.parameters[studyId].currentS)
//	h.parameters[studyId].shloopitr = 0
//	h.parameters[studyId].n = int((h.parameters[studyId].b_l/h.parameters[studyId].r_l)*(math.Pow(h.parameters[studyId].eta, float64(h.parameters[studyId].currentS))/float64(h.parameters[studyId].currentS+1))) + 1
//	h.parameters[studyId].r = h.parameters[studyId].r_l * math.Pow(h.parameters[studyId].eta, float64(-h.parameters[studyId].currentS))
//}
//
//func (h *HyperBandSuggestService) shLoopParamUpdate(studyId string) (int, int) {
//	log.Printf("SH loop i = %v", h.parameters[studyId].shloopitr)
//	pn_i := int(float64(h.parameters[studyId].n) * math.Pow(h.parameters[studyId].eta, float64(-h.parameters[studyId].shloopitr+1)) / h.parameters[studyId].eta)
//	r_i := int(h.parameters[studyId].r * math.Pow(h.parameters[studyId].eta, float64(h.parameters[studyId].shloopitr)))
//	return pn_i, r_i
//}
//
//func (h *HyperBandSuggestService) makeBracker(ctx context.Context, cl *api.ManagerClient, sid string, logs []string) (Bracket, error) {
//	req := &api.GetWorkersRequest{StudyId: sid}
//	r, err := api.GetWorkers(ctx, req)
//	if err != nil {
//		return nil, err
//	}
//	mreq := &api.GetMetrics{
//		StudyId:      sid,
//		WorkerId:     logs,
//		MetricsNames: []string{h.ResourceName},
//	}
//	mr, err := api.GetMetrics(ctx, mreq)
//	if err != nil {
//		return nil, err
//	}
//	e_l := make([]Evals, len(logs))
//	for i, l := range logs {
//		for _, m := range mr.MetricsLogs {
//			if m.WorkerId == l {
//				if len(m.MetricsLogs) == 0 {
//					break
//				}
//				e_l[i].Value = strconv.ParseFloat(m.MetricsLogs[len(m.MetricsLogs)-1].Value, 64)
//			}
//		}
//		for _, w := range r.Workers {
//			if w.WorkerId == l {
//				e_l[i].Id = w.TrialId
//			}
//		}
//	}
//	return Bracket(e_l)
//}

func (h *HyperBandSuggestService) GetSuggestions(ctx context.Context, in *api.GetSuggestionsRequest) (*api.GetSuggestionsReply, error) {
	return &api.GetSuggestionsReply{}, nil
	//	conn, err := grpc.Dial(manager, grpc.WithInsecure())
	//	if err != nil {
	//		log.Fatalf("could not connect: %v", err)
	//		return
	//	}
	//	defer conn.Close()
	//	c := api.NewManagerClient(conn)
	//	screq := &api.GetStudyRequest{
	//		StudyId: in.StudyId,
	//	}
	//	scr, err := GetStudy(ctx, screq)
	//	if err != nil {
	//		log.Fatalf("GetStudyConf failed: %v", err)
	//		return &api.GetSuggestionsReply{}, err
	//	}
	//	spreq := &api.GetSuggestionParametersRequest{
	//		StudyId:             in.StudyId,
	//		SuggestionAlgorithm: in.SuggestionAlgorithm,
	//	}
	//	spr, err := c.GetSuggestionParameters(ctx, spreq)
	//	if err != nil {
	//		log.Fatalf("GetParameter failed: %v", err)
	//		return &api.GetSuggestionsReply{}, err
	//	}
	//	hp, err := h.purseSuggestionParameters(spr.SuggestionParameters)
	//	if err != nil {
	//		return &api.GetSuggestionsReply{}, err
	//	}
	//
	//	if hp.currentS <= 0 {
	//		return &api.GetSuggestionsReply{}, nil
	//	}
	//
	//	if len(in.LogWorkerIds) > 0 {
	//
	//		var schec int
	//		var bid string
	//		for _, c := range in.CompletedTrials {
	//			schec = 0
	//			value, _ := h.dbIf.GetTrialLogs(c.TrialId,
	//				&db.GetTrialLogOpts{Objective: true, Descending: true, Limit: 1})
	//			if len(value) != 1 {
	//				log.Printf("objective value for %s not found",
	//					c.TrialId)
	//				continue
	//			}
	//			c.ObjectiveValue = value[0].Value
	//			for _, t := range c.Tags {
	//				if t.Name == "HyperBand_shi" && t.Value == strconv.Itoa(h.parameters[in.StudyId].shloopitr) {
	//					schec++
	//				}
	//				if t.Name == "HyperBand_s" && t.Value == strconv.Itoa(h.parameters[in.StudyId].currentS) {
	//					schec++
	//				}
	//				if t.Name == "HyperBand_BracketID" {
	//					bid = t.Value
	//				}
	//			}
	//			if schec == 2 {
	//				for _, b := range h.parameters[in.StudyId].MasterBracket {
	//					for _, t := range b.Tags {
	//						if t.Name == "HyperBand_BracketID" && t.Value == bid {
	//							b.ObjectiveValue = c.ObjectiveValue
	//						}
	//					}
	//				}
	//			}
	//		}
	//		sort.Sort(h.parameters[in.StudyId].MasterBracket)
	//	}
	//	var evalT []*api.Trial
	//	var r_i int
	//	if h.parameters.shloopitr > h.parameters.currentS {
	//		h.parameters.currentS--
	//		h.hbLoopParamUpdate(in.StudyId)
	//		_, r_i = h.shLoopParamUpdate(in.StudyId)
	//		h.parameters.MasterBracket = h.makeMasterBracket(in.Configs, h.parameters.n)
	//		evalT = h.getHyperParameter(in.StudyId, in.Configs, h.parameters.n)
	//		h.parameters.shloopitr++
	//	} else {
	//		var pn_i int
	//		pn_i, r_i = h.shLoopParamUpdate(in.StudyId)
	//		evalT = h.getHyperParameter(in.StudyId, in.Configs, pn_i)
	//		h.parameters.shloopitr++
	//	}
	//	for i := range evalT {
	//		for j := range evalT[i].ParameterSet {
	//			if evalT[i].ParameterSet[j].Name == h.parameters.ResourceName {
	//				evalT[i].ParameterSet[j].Value = strconv.Itoa(r_i)
	//			}
	//		}
	//		evalT[i].Tags = append(evalT[i].Tags, &api.Tag{Name: "HyperBand_s", Value: strconv.Itoa(h.parameters.currentS)})
	//		evalT[i].Tags = append(evalT[i].Tags, &api.Tag{Name: "HyperBand_r", Value: strconv.Itoa(r_i)})
	//		evalT[i].Tags = append(evalT[i].Tags, &api.Tag{Name: "HyperBand_shi", Value: strconv.Itoa(h.parameters.shloopitr)})
	//		log.Printf("Gen Trial %v", evalT[i].Tags)
	//	}
	//	return &api.GenerateTrialsReply{Trials: evalT, Completed: false}, nil
}
