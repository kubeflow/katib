package suggestion

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	api "github.com/kubeflow/katib/pkg/api/v1alpha1"
	common "github.com/kubeflow/katib/pkg/common/v1alpha1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

type Evals struct {
	id    string
	value float64
}
type Bracket []Evals

func (b Bracket) Len() int {
	return len(b)
}

func (b Bracket) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Bracket) Less(i, j int) bool {
	return b[i].value < b[j].value
}

type HyperBandParameters struct {
	eta                float64
	sMax               int
	bL                 float64
	rL                 float64
	r                  float64
	n                  int
	shloopitr          int
	currentS           int
	ResourceName       string
	ObjectiveValueName string
	evaluatingTrials   []string
}

type HyperBandSuggestService struct {
	RandomSuggestService
}

func NewHyperBandSuggestService() *HyperBandSuggestService {
	return &HyperBandSuggestService{}
}

func (h *HyperBandSuggestService) makeBracket(ctx context.Context, c api.ManagerClient, studyID string, hbparam *HyperBandParameters) ([]string, []*api.Trial, error) {
	var tids []string
	var trials []*api.Trial
	var err error
	n, r := h.getLoopParam(hbparam)

	if len(hbparam.evaluatingTrials) == 0 {
		tids, trials, err = h.makeMasterBracket(ctx, c, studyID, n, r, hbparam)
	} else {
		err, b := h.evalWorkers(ctx, c, studyID, hbparam)
		if err != nil {
			return nil, nil, err
		}
		if b == nil {
			return nil, nil, nil
		}
		newTrialNum := int(math.Ceil(float64(n) / hbparam.eta))
		h.shLoopParamUpdate(hbparam)
		_, ri := h.getLoopParam(hbparam)
		tids, trials, err = h.makeChildBracket(ctx, c, b, studyID, newTrialNum, ri, hbparam)
	}
	// evaluatingTrials is no need for last inner loop since makeMasterBracket will work for next loop
	if hbparam.shloopitr < hbparam.currentS {
		hbparam.evaluatingTrials = tids
	} else {
		hbparam.evaluatingTrials = []string{}
	}
	klog.Infof("Hyb Param sMax %v", hbparam.sMax)
	klog.Infof("Hyb Param B %v", hbparam.bL)
	klog.Infof("Hyb Param n %v", hbparam.n)
	klog.Infof("Hyb Param r %v", hbparam.r)
	klog.Infof("Hyb Param currentS %v", hbparam.currentS)
	klog.Infof("Hyb Param shloopitr %v", hbparam.shloopitr)
	klog.Infof("Hyb Param evaluatingTrials %v", hbparam.evaluatingTrials)
	if err == nil && len(hbparam.evaluatingTrials) == 0 {
		hbparam.currentS--
		if hbparam.currentS >= 0 {
			h.hbLoopParamUpdate(hbparam)
		}
	}
	return tids, trials, err
}

func (h *HyperBandSuggestService) makeMasterBracket(ctx context.Context, c api.ManagerClient, studyID string, n int, r float64, hbparam *HyperBandParameters) ([]string, []*api.Trial, error) {
	klog.Infof("Make MasterBracket %v Trials", n)
	gsreq := &api.GetStudyRequest{
		StudyId: studyID,
	}
	gsrep, err := c.GetStudy(ctx, gsreq)
	if err != nil {
		klog.Info("GetStudy Error")
		return nil, nil, err
	}
	sconf := gsrep.StudyConfig
	tids := make([]string, n)
	ts := make([]*api.Trial, n)
	for i := 0; i < n; i++ {
		t := &api.Trial{
			StudyId: studyID,
		}
		t.ParameterSet = make([]*api.Parameter, len(sconf.ParameterConfigs.Configs))
		for j, pc := range sconf.ParameterConfigs.Configs {
			t.ParameterSet[j] = &api.Parameter{Name: pc.Name}
			t.ParameterSet[j].ParameterType = pc.ParameterType
			if pc.Name == hbparam.ResourceName {
				if pc.ParameterType == api.ParameterType_INT {
					t.ParameterSet[j].Value = strconv.Itoa(int(r))
				} else {
					t.ParameterSet[j].Value = strconv.FormatFloat(r, 'f', 4, 64)
				}
			} else {
				switch pc.ParameterType {
				case api.ParameterType_INT:
					imin, _ := strconv.Atoi(pc.Feasible.Min)
					imax, _ := strconv.Atoi(pc.Feasible.Max)
					t.ParameterSet[j].Value = strconv.Itoa(h.IntRandom(imin, imax))
				case api.ParameterType_DOUBLE:
					dmin, _ := strconv.ParseFloat(pc.Feasible.Min, 64)
					dmax, _ := strconv.ParseFloat(pc.Feasible.Max, 64)
					t.ParameterSet[j].Value = strconv.FormatFloat(h.DoubleRandom(dmin, dmax), 'f', 4, 64)
				case api.ParameterType_CATEGORICAL:
					t.ParameterSet[j].Value = pc.Feasible.List[h.IntRandom(0, len(pc.Feasible.List)-1)]
				}
			}
		}
		req := &api.CreateTrialRequest{
			Trial: t,
		}
		ret, err := c.CreateTrial(ctx, req)
		if err != nil {
			klog.Info("CreateTrial Error")
			return nil, nil, err
		}
		tids[i] = ret.TrialId
		t.TrialId = ret.TrialId
		ts[i] = t
	}
	return tids, ts, nil
}

func (h *HyperBandSuggestService) makeChildBracket(ctx context.Context, c api.ManagerClient, parent Bracket, studyID string, n int, rI float64, hbparam *HyperBandParameters) ([]string, []*api.Trial, error) {
	klog.Infof("Make ChildBracket %v Trials", n)
	gsreq := &api.GetStudyRequest{
		StudyId: studyID,
	}
	gsrep, err := c.GetStudy(ctx, gsreq)
	if err != nil {
		klog.Info("GetStudy Error")
		return nil, nil, err
	}
	sconf := gsrep.StudyConfig
	child := Bracket{}

	if sconf.OptimizationType == api.OptimizationType_MINIMIZE {
		child = parent[:n]
	} else if sconf.OptimizationType == api.OptimizationType_MAXIMIZE {
		child = parent[len(parent)-n:]
	}
	gtreq := &api.GetTrialsRequest{
		StudyId: studyID,
	}
	gtrep, err := c.GetTrials(ctx, gtreq)
	if err != nil {
		klog.Info("GetTrials Error")
		return nil, nil, err
	}
	tids := make([]string, n)
	ts := make([]*api.Trial, n)
	var rtype api.ParameterType
	for _, pc := range sconf.ParameterConfigs.Configs {
		if pc.Name == hbparam.ResourceName {
			rtype = pc.ParameterType
		}
	}
	for i, tid := range child {
		t := &api.Trial{
			StudyId: studyID,
		}
		for _, pt := range gtrep.Trials {
			if pt.TrialId == tid.id {
				t.ParameterSet = pt.ParameterSet
			}
		}
		for i, p := range t.ParameterSet {
			if p.Name == hbparam.ResourceName {
				if rtype == api.ParameterType_INT {
					t.ParameterSet[i].Value = strconv.Itoa(int(rI))
				} else {
					t.ParameterSet[i].Value = strconv.FormatFloat(rI, 'f', 4, 64)
				}
			}
		}
		req := &api.CreateTrialRequest{
			Trial: t,
		}
		ret, err := c.CreateTrial(ctx, req)
		if err != nil {
			klog.Info("CreateTrial Error")
			return nil, nil, err
		}
		tids[i] = ret.TrialId
		t.TrialId = ret.TrialId
		ts[i] = t
	}
	return tids, ts, nil
}

func (h *HyperBandSuggestService) parseSuggestionParameters(ctx context.Context, c api.ManagerClient, studyID string, sparam []*api.SuggestionParameter) (*HyperBandParameters, error) {
	p := &HyperBandParameters{
		eta:                -1,
		sMax:               -1,
		bL:                 -1,
		rL:                 -1,
		r:                  -1,
		n:                  -1,
		shloopitr:          -1,
		currentS:           -2,
		ResourceName:       "",
		ObjectiveValueName: "",
		evaluatingTrials:   []string{},
	}
	for _, sp := range sparam {
		switch sp.Name {
		case "eta":
			p.eta, _ = strconv.ParseFloat(sp.Value, 64)
		case "r_l":
			p.rL, _ = strconv.ParseFloat(sp.Value, 64)
		case "ResourceName":
			p.ResourceName = sp.Value
		case "ObjectiveValueName":
			p.ObjectiveValueName = sp.Value
		case "b_l":
			p.bL, _ = strconv.ParseFloat(sp.Value, 64)
		case "sMax":
			p.sMax, _ = strconv.Atoi(sp.Value)
		case "r":
			p.r, _ = strconv.ParseFloat(sp.Value, 64)
		case "n":
			p.n, _ = strconv.Atoi(sp.Value)
		case "shloopitr":
			p.shloopitr, _ = strconv.Atoi(sp.Value)
		case "currentS":
			p.currentS, _ = strconv.Atoi(sp.Value)
		case "evaluatingTrials":
			if sp.Value != "" {
				p.evaluatingTrials = strings.Split(sp.Value, ",")
			}
		default:
			klog.Infof("Unknown Suggestion Parameter %v", sp.Name)
		}
	}
	if p.currentS == -1 {
		klog.Infof("HyperBand outlerloop has finished")
		return p, nil
	}
	if p.rL <= 0 || p.ResourceName == "" {
		klog.Info("Failed to parse Suggestion Parameter. r_l and ResourceName must be set.")
		return nil, fmt.Errorf("Suggestion Parameter set Error")
	}
	if p.eta <= 0 {
		p.eta = 3
	}
	if p.ObjectiveValueName == "" {
		gsreq := &api.GetStudyRequest{
			StudyId: studyID,
		}
		gsrep, err := c.GetStudy(ctx, gsreq)
		if err != nil {
			klog.Info("GetStudy Error")
			return nil, err
		}
		p.ObjectiveValueName = gsrep.StudyConfig.ObjectiveValueName
	}
	if p.sMax == -1 {
		p.sMax = int(math.Trunc(math.Log(p.rL) / math.Log(p.eta)))
	}
	if p.bL == -1 {
		p.bL = float64((p.sMax + 1.0)) * p.rL
	}
	if p.n == -1 {
		p.n = int(math.Ceil((p.bL / p.rL) * (math.Pow(p.eta, float64(p.sMax)) / float64(p.sMax+1))))
	}
	if p.currentS == -2 {
		p.currentS = p.sMax
	}
	if p.shloopitr == -1 {
		p.shloopitr = 0
	}
	if p.r == -1 {
		p.r = p.rL * math.Pow(p.eta, float64(-p.sMax))
	}
	return p, nil
}

func (h *HyperBandSuggestService) saveSuggestionParameters(ctx context.Context, c api.ManagerClient, studyID string, algorithm string, paramID string, hbparam *HyperBandParameters) error {
	req := &api.SetSuggestionParametersRequest{
		StudyId:             studyID,
		SuggestionAlgorithm: algorithm,
		ParamId:             paramID,
	}
	sp := []*api.SuggestionParameter{}
	sp = append(sp, &api.SuggestionParameter{
		Name:  "eta",
		Value: strconv.FormatFloat(hbparam.eta, 'f', 4, 64),
	})
	sp = append(sp, &api.SuggestionParameter{
		Name:  "sMax",
		Value: strconv.Itoa(hbparam.sMax),
	})
	sp = append(sp, &api.SuggestionParameter{
		Name:  "b_l",
		Value: strconv.FormatFloat(hbparam.bL, 'f', 4, 64),
	})
	sp = append(sp, &api.SuggestionParameter{
		Name:  "r_l",
		Value: strconv.FormatFloat(hbparam.rL, 'f', 4, 64),
	})
	sp = append(sp, &api.SuggestionParameter{
		Name:  "r",
		Value: strconv.FormatFloat(hbparam.r, 'f', 4, 64),
	})
	sp = append(sp, &api.SuggestionParameter{
		Name:  "shloopitr",
		Value: strconv.Itoa(hbparam.shloopitr),
	})
	sp = append(sp, &api.SuggestionParameter{
		Name:  "n",
		Value: strconv.Itoa(hbparam.n),
	})
	sp = append(sp, &api.SuggestionParameter{
		Name:  "currentS",
		Value: strconv.Itoa(hbparam.currentS),
	})
	sp = append(sp, &api.SuggestionParameter{
		Name:  "ResourceName",
		Value: hbparam.ResourceName,
	})
	sp = append(sp, &api.SuggestionParameter{
		Name:  "evaluatingTrials",
		Value: strings.Join(hbparam.evaluatingTrials, ","),
	})
	req.SuggestionParameters = sp
	_, err := c.SetSuggestionParameters(ctx, req)
	return err
}

func (h *HyperBandSuggestService) evalWorkers(ctx context.Context, c api.ManagerClient, studyID string, hbparam *HyperBandParameters) (error, Bracket) {
	bracket := Bracket{}
	for _, tid := range hbparam.evaluatingTrials {
		gwreq := &api.GetWorkersRequest{
			StudyId: studyID,
			TrialId: tid,
		}
		gwrep, err := c.GetWorkers(ctx, gwreq)
		if err != nil {
			klog.Errorf("GetWorkers error %v", err)
			return err, nil
		}
		wl := make([]string, len(gwrep.Workers))
		for i, w := range gwrep.Workers {
			wl[i] = w.WorkerId
		}
		gmreq := &api.GetMetricsRequest{
			StudyId:      studyID,
			WorkerIds:    wl,
			MetricsNames: []string{hbparam.ObjectiveValueName},
		}
		gmrep, err := c.GetMetrics(ctx, gmreq)
		if err != nil {
			klog.Errorf("GetMetrics error %v", err)
			return err, nil
		}
		vs := 0.0
		for _, ml := range gmrep.MetricsLogSets {
			if ml.WorkerStatus != api.State_COMPLETED {
				return nil, nil
			}
			if len(ml.MetricsLogs) == 0 {
				return nil, nil
			}
			v, _ := strconv.ParseFloat(ml.MetricsLogs[0].Values[len(ml.MetricsLogs[0].Values)-1].Value, 64)
			vs += v
		}
		if len(gwrep.Workers) > 0 {
			bracket = append(bracket, Evals{
				id:    gwrep.Workers[0].TrialId,
				value: vs / float64(len(gwrep.Workers)),
			})
		} else {
			return nil, nil
		}

	}
	sort.Sort(bracket)
	return nil, bracket
}

func (h *HyperBandSuggestService) hbLoopParamUpdate(hbparam *HyperBandParameters) {
	hbparam.shloopitr = 0
	hbparam.n = int(math.Ceil((hbparam.bL / hbparam.rL) * (math.Pow(hbparam.eta, float64(hbparam.currentS)) / float64(hbparam.currentS+1))))
	hbparam.r = hbparam.rL * math.Pow(hbparam.eta, float64(-hbparam.currentS))
}

func (h *HyperBandSuggestService) getLoopParam(hbparam *HyperBandParameters) (int, float64) {
	nI := int(math.Ceil(float64(hbparam.n) * math.Pow(hbparam.eta, float64(-hbparam.shloopitr))))
	rI := hbparam.r * math.Pow(hbparam.eta, float64(hbparam.shloopitr))
	return nI, rI
}

func (h *HyperBandSuggestService) shLoopParamUpdate(hbparam *HyperBandParameters) {
	hbparam.shloopitr++
	if hbparam.shloopitr > hbparam.currentS {
		hbparam.currentS--
		if hbparam.currentS >= 0 {
			h.hbLoopParamUpdate(hbparam)
		}
	}
}

func (h *HyperBandSuggestService) GetSuggestions(ctx context.Context, in *api.GetSuggestionsRequest) (*api.GetSuggestionsReply, error) {
	conn, err := grpc.Dial(common.ManagerAddr, grpc.WithInsecure())
	if err != nil {
		klog.Fatalf("could not connect: %v", err)
		return &api.GetSuggestionsReply{}, err
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)
	spreq := &api.GetSuggestionParametersRequest{
		ParamId: in.ParamId,
	}
	spr, err := c.GetSuggestionParameters(ctx, spreq)
	if err != nil {
		klog.Fatalf("GetParameter failed: %v", err)
		return &api.GetSuggestionsReply{}, err
	}
	hbparam, err := h.parseSuggestionParameters(ctx, c, in.StudyId, spr.SuggestionParameters)
	if err != nil {
		return &api.GetSuggestionsReply{}, err
	}

	if hbparam.currentS < 0 {
		return &api.GetSuggestionsReply{}, nil
	}

	tids, ts, err := h.makeBracket(ctx, c, in.StudyId, hbparam)
	if err != nil {
		return &api.GetSuggestionsReply{}, err
	}
	if tids == nil {
		return &api.GetSuggestionsReply{}, status.Errorf(codes.FailedPrecondition, "Previous workers are not completed.")
	}

	err = h.saveSuggestionParameters(ctx, c, in.StudyId, in.SuggestionAlgorithm, in.ParamId, hbparam)
	if err != nil {
		klog.Fatalf("saveSuggestionParameters failed: %v", err)
		return &api.GetSuggestionsReply{}, err
	}
	return &api.GetSuggestionsReply{
		Trials: ts,
	}, nil
}

func (s *HyperBandSuggestService) ValidateSuggestionParameters(ctx context.Context, in *api.ValidateSuggestionParametersRequest) (*api.ValidateSuggestionParametersReply, error) {

	return &api.ValidateSuggestionParametersReply{}, nil
}
