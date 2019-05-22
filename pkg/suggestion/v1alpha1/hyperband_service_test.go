package suggestion

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"

	api "github.com/kubeflow/katib/pkg/api/v1alpha1"
	mockapi "github.com/kubeflow/katib/pkg/mock/v1alpha1/api"
)

func getSampleBracket(size int) Bracket {
	sample := []Evals{}
	for i := 0; i < size; i++ {
		sample = append(sample, Evals{id: "test" + strconv.Itoa(i+1), value: float64(i + 1)})
	}
	return sample
}

func getSampleHyperBandParameters() HyperBandParameters {
	sample := HyperBandParameters{
		eta:                4.0,
		sMax:               100,
		bL:                 60.0,
		rL:                 30.0,
		r:                  8.0,
		n:                  50,
		shloopitr:          2,
		currentS:           3,
		ResourceName:       "testResource",
		ObjectiveValueName: "testValue",
		evaluatingTrials:   []string{"trial1", "trial2", "trial3"},
	}
	return sample
}

func getSampleStudyReply() *api.GetStudyReply {
	fs := &api.FeasibleSpace{
		Max:  "10",
		Min:  "1",
		List: []string{"alpha", "beta"},
	}
	c1 := &api.ParameterConfig{
		Name:          "config1",
		ParameterType: 2,
		Feasible:      fs,
	}
	c2 := &api.ParameterConfig{
		Name:          "config2",
		ParameterType: 2,
		Feasible:      fs,
	}
	sp := []*api.ParameterConfig{c1, c2}
	scpc := &api.StudyConfig_ParameterConfigs{
		Configs: sp,
	}
	sc := &api.StudyConfig{
		OptimizationType: 1,
		ParameterConfigs: scpc,
	}
	srep := &api.GetStudyReply{
		StudyConfig: sc,
	}

	return srep
}

func getSampleTrialsReply() *api.GetTrialsReply {
	t1 := &api.Trial{
		TrialId: "trialId1",
	}
	t2 := &api.Trial{
		TrialId: "trialId2",
	}
	ts := []*api.Trial{t1, t2}
	trep := &api.GetTrialsReply{
		Trials: ts,
	}

	return trep
}

func getExpectedTrials(n int) ([]string, []*api.Trial) {
	exp_tids := []string{"trialId", "trialId"}
	exp_ts := make([]*api.Trial, n)
	parameter1 := &api.Parameter{
		Name:          "config1",
		ParameterType: 1,
		Value:         "1.0",
	}
	parameter2 := &api.Parameter{
		Name:          "config2",
		ParameterType: 2,
		Value:         "1",
	}
	parameter3 := &api.Parameter{
		Name:          "config3",
		ParameterType: 4,
		Value:         "test",
	}
	parameterSet := []*api.Parameter{parameter1, parameter2, parameter3}
	trial := &api.Trial{
		TrialId:      "trialId",
		StudyId:      "studyId",
		ParameterSet: parameterSet,
	}
	for i := range exp_ts {
		exp_ts[i] = trial
	}

	return exp_tids, exp_ts
}

func getSampleSuggestionParameters() []*api.SuggestionParameter {
	exp_sparam := make([]*api.SuggestionParameter, 0)

	eta := &api.SuggestionParameter{
		Name:  "eta",
		Value: "1.5",
	}
	exp_sparam = append(exp_sparam, eta)
	r_l := &api.SuggestionParameter{
		Name:  "r_l",
		Value: "10.5",
	}
	exp_sparam = append(exp_sparam, r_l)
	ResourceName := &api.SuggestionParameter{
		Name:  "ResourceName",
		Value: "testResourceName",
	}
	exp_sparam = append(exp_sparam, ResourceName)
	ObjectiveValueName := &api.SuggestionParameter{
		Name:  "ObjectiveValueName",
		Value: "testObjectiveValueName",
	}
	exp_sparam = append(exp_sparam, ObjectiveValueName)
	b_l := &api.SuggestionParameter{
		Name:  "b_l",
		Value: "10.5",
	}
	exp_sparam = append(exp_sparam, b_l)
	sMax := &api.SuggestionParameter{
		Name:  "sMax",
		Value: "100",
	}
	exp_sparam = append(exp_sparam, sMax)
	r := &api.SuggestionParameter{
		Name:  "r",
		Value: "10.5",
	}
	exp_sparam = append(exp_sparam, r)
	n := &api.SuggestionParameter{
		Name:  "n",
		Value: "10",
	}
	exp_sparam = append(exp_sparam, n)
	shloopitr := &api.SuggestionParameter{
		Name:  "shloopitr",
		Value: "10",
	}
	exp_sparam = append(exp_sparam, shloopitr)
	currentS := &api.SuggestionParameter{
		Name:  "currentS",
		Value: "10",
	}
	exp_sparam = append(exp_sparam, currentS)
	evaluatingTrials := &api.SuggestionParameter{
		Name:  "evaluatingTrials",
		Value: "trial1,trial2,trial3",
	}
	exp_sparam = append(exp_sparam, evaluatingTrials)

	return exp_sparam
}

func getSampleWorkersRequest() []*api.GetWorkersRequest {
	wreq1 := &api.GetWorkersRequest{
		StudyId: "studyId",
		TrialId: "trial1",
	}
	wreq2 := &api.GetWorkersRequest{
		StudyId: "studyId",
		TrialId: "trial2",
	}
	wreq3 := &api.GetWorkersRequest{
		StudyId: "studyId",
		TrialId: "trial3",
	}
	wreq := []*api.GetWorkersRequest{wreq1, wreq2, wreq3}

	return wreq
}

func getSampleWorkersReplyArray() []*api.GetWorkersReply {
	worker1 := &api.Worker{
		WorkerId: "worker1",
		TrialId:  "trial1",
	}
	worker2 := &api.Worker{
		WorkerId: "worker2",
		TrialId:  "trial2",
	}
	worker3 := &api.Worker{
		WorkerId: "worker3",
		TrialId:  "trial3",
	}
	wrep1 := &api.GetWorkersReply{
		Workers: []*api.Worker{worker1},
	}
	wrep2 := &api.GetWorkersReply{
		Workers: []*api.Worker{worker2},
	}
	wrep3 := &api.GetWorkersReply{
		Workers: []*api.Worker{worker3},
	}
	wrep := []*api.GetWorkersReply{wrep1, wrep2, wrep3}
	return wrep
}

func getSampleMetricsRequestArray(ovn string) []*api.GetMetricsRequest {
	metrics1 := &api.GetMetricsRequest{
		StudyId:      "studyId",
		WorkerIds:    []string{"worker1"},
		MetricsNames: []string{ovn},
	}
	metrics2 := &api.GetMetricsRequest{
		StudyId:      "studyId",
		WorkerIds:    []string{"worker2"},
		MetricsNames: []string{ovn},
	}
	metrics3 := &api.GetMetricsRequest{
		StudyId:      "studyId",
		WorkerIds:    []string{"worker3"},
		MetricsNames: []string{ovn},
	}
	metrics := []*api.GetMetricsRequest{metrics1, metrics2, metrics3}

	return metrics
}

func getSampleMetricsReplyArray() []*api.GetMetricsReply {
	mvt := make([]*api.MetricsValueTime, 0)
	for i := 0; i < 24; i++ {
		mvt = append(mvt, &api.MetricsValueTime{Time: "Jan 1 12:00:00 2018 UTC", Value: strconv.Itoa(i)})
	}

	mLog := make([]*api.MetricsLog, 0)
	for j := 0; j < 6; j++ {
		mLog = append(mLog, &api.MetricsLog{Name: "metricsLog" + strconv.Itoa(j), Values: mvt[j*4 : (j+1)*4]})
	}

	mLogSet := make([]*api.MetricsLogSet, 0)
	for k := 0; k < 3; k++ {
		mLogSet = append(mLogSet, &api.MetricsLogSet{WorkerId: "worker" + strconv.Itoa(k+1), MetricsLogs: mLog[k*2 : (k+1)*2], WorkerStatus: 2})
	}

	mrepArray := make([]*api.GetMetricsReply, 0)
	for l := 0; l < 3; l++ {
		mrepArray = append(mrepArray, &api.GetMetricsReply{MetricsLogSets: []*api.MetricsLogSet{mLogSet[l]}})
	}

	return mrepArray
}

func TestLen(t *testing.T) {
	size := 2
	b := getSampleBracket(size)
	exp := 2
	rtn := b.Len()

	if exp != rtn {
		t.Errorf("expected %v, but returned %v", exp, rtn)
	}
}

func TestSwap(t *testing.T) {
	size := 2
	b := getSampleBracket(size)
	exp := Bracket{
		Evals{
			id:    "test2",
			value: 2.0,
		},
		Evals{
			id:    "test1",
			value: 1.0,
		},
	}

	b.Swap(0, 1)

	if exp[0] != b[0] || exp[1] != b[1] {
		t.Errorf("expected %v, but returned %v", exp, b)
	}
}

func TestLess(t *testing.T) {
	size := 2
	b := getSampleBracket(size)
	exp := false
	rtn := b.Less(0, 1)

	if exp != rtn {
		t.Errorf("expected %v, but returned %v", exp, rtn)
	}
}

func TestMakeMasterBracket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := HyperBandSuggestService{}
	mockAPI := mockapi.NewMockManagerClient(ctrl)

	sreq := &api.GetStudyRequest{
		StudyId: "studyId",
	}
	srep := getSampleStudyReply()
	mockAPI.EXPECT().GetStudy(context.Background(), sreq).Return(srep, nil)

	trep := &api.CreateTrialReply{
		TrialId: "trialId",
	}
	mockAPI.EXPECT().CreateTrial(context.Background(), gomock.Any()).Return(trep, nil).AnyTimes()

	n := 3
	r := 2.0
	p := getSampleHyperBandParameters()

	exp_tids, exp_ts := getExpectedTrials(n)
	rtn_tids, rtn_ts, err := h.makeMasterBracket(context.Background(), mockAPI, "studyId", n, r, &p)
	if err != nil {
		t.Errorf("makeMasterBracket Error: %v", err)
	}

	for j := range exp_tids {
		if exp_tids[j] != rtn_tids[j] {
			t.Errorf("expected %v, but returned %v", exp_tids[j], rtn_tids[j])
		}
		if exp_ts[j].TrialId != rtn_ts[j].TrialId || exp_ts[j].StudyId != rtn_ts[j].StudyId {
			t.Errorf("expected %v, but returned %v", exp_ts[j], rtn_ts[j])
		}
	}
}

func TestMakeChildBracket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := HyperBandSuggestService{}
	mockAPI := mockapi.NewMockManagerClient(ctrl)

	sreq := &api.GetStudyRequest{
		StudyId: "studyId",
	}
	srep := getSampleStudyReply()
	mockAPI.EXPECT().GetStudy(context.Background(), sreq).Return(srep, nil)

	treq := &api.GetTrialsRequest{
		StudyId: "studyId",
	}
	strep := getSampleTrialsReply()
	mockAPI.EXPECT().GetTrials(context.Background(), treq).Return(strep, nil)

	trep := &api.CreateTrialReply{
		TrialId: "trialId",
	}
	mockAPI.EXPECT().CreateTrial(context.Background(), gomock.Any()).Return(trep, nil).AnyTimes()

	n := 3
	r := 2.0
	p := getSampleHyperBandParameters()
	size := 10
	b := getSampleBracket(size)
	exp_tids, exp_ts := getExpectedTrials(n)
	rtn_tids, rtn_ts, err := h.makeChildBracket(context.Background(), mockAPI, b, "studyId", n, r, &p)
	if err != nil {
		t.Errorf("makeChildBracket Error: %v", err)
	}

	for j := range exp_tids {
		if exp_tids[j] != rtn_tids[j] {
			t.Errorf("expected %v, but returned %v", exp_tids[j], rtn_tids[j])
		}
		if exp_ts[j].TrialId != rtn_ts[j].TrialId || exp_ts[j].StudyId != rtn_ts[j].StudyId {
			t.Errorf("expected %v, but returned %v", exp_ts[j], rtn_ts[j])
		}
	}

}

func TestParseSuggestionParameters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := HyperBandSuggestService{}
	p1 := getSampleSuggestionParameters()
	mockAPI := mockapi.NewMockManagerClient(ctrl)

	rtn_param1, _ := h.parseSuggestionParameters(context.Background(), mockAPI, "studyId", p1)

	exp_param1 := &HyperBandParameters{
		eta:                1.5,
		sMax:               100,
		bL:                 10.5,
		rL:                 10.5,
		r:                  10.5,
		n:                  10,
		shloopitr:          10,
		currentS:           10,
		ResourceName:       "testResourceName",
		ObjectiveValueName: "testObjectiveValueName",
		evaluatingTrials:   []string{"trial1", "trial2", "trial3"},
	}

	if !reflect.DeepEqual(exp_param1, rtn_param1) {
		t.Errorf("expected %v , but returned %v", exp_param1, rtn_param1)
	}

	p2 := make([]*api.SuggestionParameter, 0)
	rtn_param2, _ := h.parseSuggestionParameters(context.Background(), mockAPI, "studyId", p2)

	if rtn_param2 != nil {
		t.Errorf("expected nil, but returned %v", rtn_param2)
	}

	p3 := make([]*api.SuggestionParameter, 0)
	r_l := &api.SuggestionParameter{
		Name:  "r_l",
		Value: "27.0",
	}
	p3 = append(p3, r_l)
	ResourceName := &api.SuggestionParameter{
		Name:  "ResourceName",
		Value: "testResourceName",
	}
	p3 = append(p3, ResourceName)
	sreq := &api.GetStudyRequest{
		StudyId: "studyId",
	}
	srep := getSampleStudyReply()
	mockAPI.EXPECT().GetStudy(context.Background(), sreq).Return(srep, nil)
	rtn_param3, _ := h.parseSuggestionParameters(context.Background(), mockAPI, "studyId", p3)

	exp_param3 := &HyperBandParameters{
		eta:                3,
		sMax:               3,
		bL:                 108,
		rL:                 27,
		r:                  1,
		n:                  27,
		shloopitr:          0,
		currentS:           3,
		ResourceName:       "testResourceName",
		ObjectiveValueName: "",
		evaluatingTrials:   []string{},
	}
	if !reflect.DeepEqual(exp_param3, rtn_param3) {
		t.Errorf("exptected %v, but returned %v", exp_param3, rtn_param3)
	}
}

func TestSaveSuggestionParameters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := HyperBandSuggestService{}
	p := getSampleHyperBandParameters()
	mockAPI := mockapi.NewMockManagerClient(ctrl)

	mockAPI.EXPECT().SetSuggestionParameters(context.Background(), gomock.Any()).Return(nil, nil)

	rtn := h.saveSuggestionParameters(context.Background(), mockAPI, "studyID", "testAlgorithm", "paramID", &p)

	if rtn != nil {
		t.Errorf("returned error: %v", rtn)
	}
}

func TestEvalWorkers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := HyperBandSuggestService{}
	p := getSampleHyperBandParameters()
	mockAPI := mockapi.NewMockManagerClient(ctrl)

	wreq := getSampleWorkersRequest()
	wrepArray := getSampleWorkersReplyArray()

	mockAPI.EXPECT().GetWorkers(context.Background(), wreq[0]).Return(wrepArray[0], nil)
	mockAPI.EXPECT().GetWorkers(context.Background(), wreq[1]).Return(wrepArray[1], nil)
	mockAPI.EXPECT().GetWorkers(context.Background(), wreq[2]).Return(wrepArray[2], nil)

	mreqArray := getSampleMetricsRequestArray(p.ObjectiveValueName)
	mrep := getSampleMetricsReplyArray()
	mockAPI.EXPECT().GetMetrics(context.Background(), mreqArray[0]).Return(mrep[0], nil)
	mockAPI.EXPECT().GetMetrics(context.Background(), mreqArray[1]).Return(mrep[1], nil)
	mockAPI.EXPECT().GetMetrics(context.Background(), mreqArray[2]).Return(mrep[2], nil)

	_, rtn_bracket := h.evalWorkers(context.Background(), mockAPI, "studyId", &p)

	exp_bracket := []Evals{{"trial3", 19}, {"trial2", 11}, {"trial1", 3}}

	for i, ebkt := range exp_bracket {
		if ebkt != rtn_bracket[i] {
			t.Errorf("Evals[%v]: expected %v, but returned %v", i, ebkt, rtn_bracket[i])
		}
	}
}

func TestHbLoopParamUpdate(t *testing.T) {
	h := HyperBandSuggestService{}
	p := getSampleHyperBandParameters()

	studyId := "testStudyId"
	h.hbLoopParamUpdate(studyId, &p)

	exp := getSampleHyperBandParameters()
	exp.r = 0.46875
	exp.n = 32
	exp.shloopitr = 0

	if exp.r != p.r || exp.n != p.n || exp.shloopitr != p.shloopitr {
		t.Errorf("expected %v, but returned %v", exp, p)
	}
}

func TestGetLoopParam(t *testing.T) {
	h := HyperBandSuggestService{}
	p := getSampleHyperBandParameters()

	studyId := "testStudyId"
	n_i, r_i := h.getLoopParam(studyId, &p)

	exp_n_i, exp_r_i := 3, 128.0

	if exp_n_i != n_i || exp_r_i != r_i {
		t.Errorf("expected {%v %v}, but returned {%v %v}", exp_n_i, exp_r_i, n_i, r_i)
	}
}

func TestShLoopParamUpdate(t *testing.T) {
	h := HyperBandSuggestService{}
	p1 := getSampleHyperBandParameters()

	studyId := "testStudyId"
	exp1 := p1.currentS
	h.shLoopParamUpdate(studyId, &p1)
	rtn1 := p1.currentS

	if exp1 != rtn1 {
		t.Errorf("expected currentS = %v, but returned currentS = %v", exp1, rtn1)
	}

	p2 := getSampleHyperBandParameters()
	p2.shloopitr = 100
	exp2 := p2.currentS
	exp2--

	h.shLoopParamUpdate(studyId, &p2)
	rtn2 := p2.currentS

	if exp2 != rtn2 {
		t.Errorf("expected currentS = %v, but returned currentS = %v", exp2, rtn2)
	}
}
