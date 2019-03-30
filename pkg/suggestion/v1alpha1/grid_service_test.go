package suggestion

import (
	"reflect"
	"testing"

	api "github.com/kubeflow/katib/pkg/api/v1alpha1"
)

func getSampleParameterConfigs() []*api.ParameterConfig {
	parameterConfigs := make([]*api.ParameterConfig, 0)

	config1 := &api.ParameterConfig{
		Name:          "config1",
		ParameterType: api.ParameterType_INT,
		Feasible: &api.FeasibleSpace{
			Max: "2",
			Min: "1",
		},
	}
	config2 := &api.ParameterConfig{
		Name:          "config2",
		ParameterType: api.ParameterType_DOUBLE,
		Feasible: &api.FeasibleSpace{
			Max: "5.5",
			Min: "3.5",
		},
	}
	config3 := &api.ParameterConfig{
		Name:          "config3",
		ParameterType: api.ParameterType_CATEGORICAL,
		Feasible: &api.FeasibleSpace{
			List: []string{"alpha", "beta", "gamma"},
		},
	}

	parameterConfigs = append(parameterConfigs, config1)
	parameterConfigs = append(parameterConfigs, config2)
	parameterConfigs = append(parameterConfigs, config3)

	return parameterConfigs
}

func getSampleSuggestionParameter(pattern int) []*api.SuggestionParameter {
	suggestParam := make([]*api.SuggestionParameter, 0)

	switch pattern {
	case 1:
		param1 := &api.SuggestionParameter{
			Name:  "DefaultGrid",
			Value: "1",
		}
		param2 := &api.SuggestionParameter{
			Name:  "Iteration",
			Value: "2",
		}
		param3 := &api.SuggestionParameter{
			Name:  "learning-rate",
			Value: "3",
		}
		suggestParam = append(suggestParam, param1)
		suggestParam = append(suggestParam, param2)
		suggestParam = append(suggestParam, param3)

	case 2:
		param := &api.SuggestionParameter{
			Name:  "DefaultGrid",
			Value: "-1",
		}
		suggestParam = append(suggestParam, param)
	}
	return suggestParam
}

func TestAllocInt(t *testing.T) {
	s := &GridSuggestService{}

	min1 := 1
	max1 := 2
	reqnum1 := 1

	exp1 := "1"
	rtn1 := s.allocInt(min1, max1, reqnum1)

	if rtn1[0] != exp1 {
		t.Errorf("expected [%v], but %v is returned", exp1, rtn1)
	}

	min2 := 1
	max2 := 9
	reqnum2 := 5

	exp2 := []string{"1", "3", "5", "7", "9"}
	rtn2 := s.allocInt(min2, max2, reqnum2)

	for i := 0; i < 5; i++ {
		if rtn2[i] != exp2[i] {
			t.Errorf("expected Array[%v] = %v, but %v is returned", i, exp2[i], rtn2[i])
		}
	}
}

func TestAllocFloat(t *testing.T) {
	s := &GridSuggestService{}

	min1 := 1.0
	max1 := 2.0
	reqnum1 := 1

	exp1 := "1.0000"
	rtn1 := s.allocFloat(min1, max1, reqnum1)

	if rtn1[0] != exp1 {
		t.Errorf("expected [%v], but %v is returned", exp1, rtn1)
	}

	min2 := 1.0
	max2 := 9.0
	reqnum2 := 5

	exp2 := []string{"1.0000", "3.0000", "5.0000", "7.0000", "9.0000"}
	rtn2 := s.allocFloat(min2, max2, reqnum2)

	for i := 0; i < 5; i++ {
		if rtn2[i] != exp2[i] {
			t.Errorf("expected Array[%v] = %v, but %v is returned", i, exp2[i], rtn2[i])
		}
	}
}

func TestAllocCat(t *testing.T) {
	s := &GridSuggestService{}

	list := []string{"alpha", "beta", "gamma"}
	reqnum1 := 1

	exp1 := []string{"alpha"}
	rtn1 := s.allocCat(list, reqnum1)

	if rtn1[0] != exp1[0] {
		t.Errorf("exptected %v, but %v", rtn1, exp1)
	}

	reqnum2 := 5
	exp2 := []string{"alpha", "alpha", "beta", "beta", "gamma"}
	rtn2 := s.allocCat(list, reqnum2)

	for i := 0; i < 5; i++ {
		if rtn2[i] != exp2[i] {
			t.Errorf("expected Array[%v] = %v, but %v is returned", i, exp2[i], rtn2[i])
		}
	}
}

func TestSetP(t *testing.T) {
	s := &GridSuggestService{}

	gci := 0
	p := make([][]*api.Parameter, 4)
	pcs := getSampleParameterConfigs()
	pg := make([][]string, 0)
	pg = append(pg, []string{"1", "2"})
	pg = append(pg, []string{"3.5000", "5.5000"})

	s.setP(gci, p, pg, pcs)

	exp := make([][]*api.Parameter, 0)

	p1 := &api.Parameter{
		Name:          "config1",
		ParameterType: api.ParameterType_INT,
		Value:         "1",
	}
	p2 := &api.Parameter{
		Name:          "config1",
		ParameterType: api.ParameterType_INT,
		Value:         "2",
	}
	p3 := &api.Parameter{
		Name:          "config2",
		ParameterType: api.ParameterType_DOUBLE,
		Value:         "3.5000",
	}
	p4 := &api.Parameter{
		Name:          "config2",
		ParameterType: api.ParameterType_DOUBLE,
		Value:         "5.5000",
	}
	exp = append(exp, []*api.Parameter{p1, p3})
	exp = append(exp, []*api.Parameter{p1, p4})
	exp = append(exp, []*api.Parameter{p2, p3})
	exp = append(exp, []*api.Parameter{p2, p4})

	for i, rtn := range p {
		if !reflect.DeepEqual(rtn, exp[i]) {
			t.Errorf("expected %v, but %v is returned", exp[i], rtn)
		}
	}
}

func TestParseSuggestParam(t *testing.T) {
	s := &GridSuggestService{}

	sp1 := getSampleSuggestionParameter(1)
	rtn_defaultGrid1, rtn_i, rtn_ret := s.parseSuggestParam(sp1)
	exp_defaultGrid1, exp_i := 1, 2
	exp_ret := make(map[string]int)
	exp_ret["learning-rate"] = 3

	if rtn_defaultGrid1 != exp_defaultGrid1 {
		t.Errorf("expected DefaultGrid is %v, but %v is returned", exp_defaultGrid1, rtn_defaultGrid1)
	}
	if rtn_i != exp_i {
		t.Errorf("expected Iteration is %v, but %v is returned", exp_i, rtn_i)
	}
	if !reflect.DeepEqual(rtn_ret, exp_ret) {
		t.Errorf("expected parameter is %v, but %v is returned", exp_ret, rtn_ret)
	}

	sp2 := getSampleSuggestionParameter(2)
	rtn_defaultGrid2, _, _ := s.parseSuggestParam(sp2)
	exp_defaultGrid2 := 1

	if rtn_defaultGrid2 != exp_defaultGrid2 {
		t.Errorf("expected DefaultGrid is %v, but %v is returned", exp_defaultGrid2, rtn_defaultGrid2)
	}
}

func TestGenGrids(t *testing.T) {
	s := &GridSuggestService{}

	studyID := "testStudy"
	pcs := getSampleParameterConfigs()
	df := 1
	glist := make(map[string]int)
	glist["config1"] = 2
	glist["config2"] = 3
	glist["config3"] = 3

	rtn := s.genGrids(studyID, pcs, df, glist)

	exp_len := 1
	for _, v := range glist {
		exp_len *= v
	}
	if len(rtn) != exp_len {
		t.Errorf("expected %v parameters generated, but %v parameters generated", exp_len, len(rtn))
	}

	config1 := []string{"1", "2"}
	config2 := []string{"3.5000", "4.5000", "5.5000"}
	config3 := []string{"alpha", "beta", "gamma"}
	iter := 0
	for _, c1 := range config1 {
		for _, c2 := range config2 {
			for _, c3 := range config3 {
				exp_struct1 := &api.Parameter{
					Name:          "config1",
					ParameterType: api.ParameterType_INT,
					Value:         c1,
				}
				exp_struct2 := &api.Parameter{
					Name:          "config2",
					ParameterType: api.ParameterType_DOUBLE,
					Value:         c2,
				}
				exp_struct3 := &api.Parameter{
					Name:          "config3",
					ParameterType: api.ParameterType_CATEGORICAL,
					Value:         c3,
				}
				exp := []*api.Parameter{exp_struct1, exp_struct2, exp_struct3}
				if !reflect.DeepEqual(rtn[iter], exp) {
					t.Errorf("expected %v, but %v is returned", exp, rtn[iter])
				}
				iter++
			}
		}
	}
}
