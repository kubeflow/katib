package genetic_algorithm

import (
	"testing"
)

func sampleEvaluateFunc0(offspring Offspring) float64 {
	// 30 is the best score
	score := 0
	for k, v := range offspring.GeneMap {
		switch k {
		case "cg0", "cg1", "cg2", "cg3", "cg4", "cg5", "cg6", "cg7", "cg8", "cg9":
			if v.(string) == "true" {
				score++
			}
		case "dg0", "dg1", "dg2", "dg3", "dg4", "dg5", "dg6", "dg7", "dg8", "dg9":
			if v.(int) == 0 {
				score++
			}
		case "cog0", "cog1", "cog2", "cog3", "cog4", "cog5", "cog6", "cog7", "cog8", "cog9":
			if v.(float64) < 0.2 {
				score++
			}
		}
	}
	return float64(score)
}

func sampleEvaluateFunc1(offspring Offspring) float64 {
	// 10 is the best score
	score := 0
	if offspring.GeneMap["cg0"].(string) == "true" &&
		offspring.GeneMap["dg0"].(int) == 0 &&
		offspring.GeneMap["cog0"].(float64) <= 0.2 {
		score++
	}
	if offspring.GeneMap["cg1"].(string) == "false" &&
		offspring.GeneMap["dg1"].(int) == 1 &&
		offspring.GeneMap["cog1"].(float64) >= 0.6 {
		score++
	}
	if offspring.GeneMap["cg2"].(string) == "true" &&
		offspring.GeneMap["dg2"].(int) == 0 &&
		offspring.GeneMap["cog2"].(float64) <= 0.5 {
		score++
	}
	if offspring.GeneMap["cg3"].(string) == "false" &&
		offspring.GeneMap["dg3"].(int) == 1 &&
		offspring.GeneMap["cog3"].(float64) >= 0.8 {
		score++
	}
	if offspring.GeneMap["cg4"].(string) == "true" &&
		offspring.GeneMap["dg4"].(int) == 0 &&
		offspring.GeneMap["cog4"].(float64) <= 0.2 {
		score++
	}
	if offspring.GeneMap["cg5"].(string) == "false" &&
		offspring.GeneMap["dg5"].(int) == 1 &&
		offspring.GeneMap["cog5"].(float64) >= 0.6 {
		score++
	}
	if offspring.GeneMap["cg6"].(string) == "false" &&
		offspring.GeneMap["dg6"].(int) == 0 &&
		offspring.GeneMap["cog6"].(float64) <= 0.4 {
		score++
	}
	if offspring.GeneMap["cg7"].(string) == "true" &&
		offspring.GeneMap["dg7"].(int) == 1 &&
		offspring.GeneMap["cog7"].(float64) >= 0.7 {
		score++
	}
	if offspring.GeneMap["cg8"].(string) == "tf" &&
		offspring.GeneMap["dg8"].(int) == 0 &&
		(offspring.GeneMap["cog8"].(float64) <= 0.2 ||
			offspring.GeneMap["cog8"].(float64) >= 0.8) {
		score++
	}
	if offspring.GeneMap["cg9"].(string) == "tf" &&
		offspring.GeneMap["dg9"].(int) == 1 &&
		(offspring.GeneMap["cog9"].(float64) >= 0.2 ||
			offspring.GeneMap["cog8"].(float64) <= 0.8) {
		score++
	}
	return float64(score)
}

func TestContinuousGeneType(t *testing.T) {
	g := NewContinuousGeneType(1.0, 2.0)
	got := g.generateRandom()
	if got >= 1.0 && got <= 2.0 {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestDiscreteGeneType(t *testing.T) {
	g := NewDiscreteGeneType(1, 3)
	got := int(g.generateRandom())
	if got >= 1 && got <= 3 {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestGetDiscreteRandomExcludeA(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := getDiscreteRandomExcludeA(1, 3, 2)
		if g == 2 {
			t.Fatal("NG")
		}
	}
	t.Log("OK")
}

func containStr(got string, arr []string) bool {
	for _, v := range arr {
		if v == got {
			return true
		}
	}
	return false
}

func TestCategoricalGeneType(t *testing.T) {
	arr := []string{"a", "b", "c"}
	g := NewCategoricalGeneType(arr)
	got := int(g.generateRandom())
	if containStr(arr[got], arr) {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestGetCategoricalRandomExcludeA(t *testing.T) {
	arr := []string{"a", "b", "c"}
	g := getCategoricalRandomExcludeA(arr, "b")
	for i := 0; i < 100; i++ {
		if arr[g] == "b" {
			t.Fatal("NG")
		}
	}
	t.Log("OK")
}

func TestGene(t *testing.T) {
	g := NewGene("gender")
	g.setGeneValue("male")
	if g.geneName == "gender" && g.geneValue == "male" {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestNewGenes(t *testing.T) {
	geneNames := []string{"g", "gender"}
	genes := NewGenes(geneNames)
	if genes[0].geneName == "g" && len(genes) == 2 {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestGeneRandomValue1(t *testing.T) {
	arr := []string{"male", "female", "none"}
	cg := NewCategoricalGeneType(arr)
	g := NewGene("gender")
	g.setRandomGeneValue(cg)
	got := g.geneValue.(string)
	if containStr(got, arr) {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestGeneRandomValue2(t *testing.T) {
	dg := NewDiscreteGeneType(1, 3)
	g := NewGene("gender")
	g.setRandomGeneValue(dg)
	got := g.geneValue.(int)
	if got >= 1 && got <= 3 {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestGeneRandomValue3(t *testing.T) {
	cg := NewContinuousGeneType(1.0, 2.0)
	g := NewGene("g")
	g.setRandomGeneValue(cg)
	got := g.geneValue.(float64)
	if got >= 1.0 && got <= 2.0 {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestOffspring(t *testing.T) {
	geneNames := []string{"gender", "somegene", "ag"}
	genes := NewGenes(geneNames)
	o := NewOffspring(genes)
	o.GeneMap["gender"] = "male"
	o.GeneMap["somegene"] = 1
	o.GeneMap["ag"] = 1.5
	if _, ok := o.GeneMap["gender"]; ok {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestGenerateInitialOffspring(t *testing.T) {
	arr := []string{"male", "female", "none"}
	cg := NewCategoricalGeneType(arr)
	gg := NewGene("gender")
	gg.setRandomGeneValue(cg)

	dg := NewDiscreteGeneType(1, 3)
	sg := NewGene("somegene")
	sg.setRandomGeneValue(dg)

	cog := NewContinuousGeneType(1.0, 2.0)
	ag := NewGene("ag")
	ag.setRandomGeneValue(cog)

	gs := []Gene{*gg, *sg, *ag}
	o := NewOffspring(gs)
	if _, ok := o.GeneMap["gender"]; ok {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestGenerateRandomOffspring(t *testing.T) {
	arr := []string{"male", "female", "none"}
	cg := NewCategoricalGeneType(arr)
	dg := NewDiscreteGeneType(1, 3)
	cog := NewContinuousGeneType(1.0, 2.0)
	geneNames := []string{"gender", "somegene", "ag"}
	geneTypes := []interface{}{*cg, *dg, *cog}
	genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	genes.setRandomGeneValues(geneTypeMap)

	o := NewOffspring(genes)
	if _, ok := o.GeneMap["gender"]; ok {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestGeneTypeMap(t *testing.T) {
	arr := []string{"male", "female", "none"}
	cg := NewCategoricalGeneType(arr)
	dg := NewDiscreteGeneType(1, 3)
	cog := NewContinuousGeneType(1.0, 2.0)
	geneNames := []string{"gender", "somegene", "ag"}
	geneTypes := []interface{}{*cg, *dg, *cog}
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	if _, ok := geneTypeMap["gender"]; ok {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestRandomGenes(t *testing.T) {
	arr := []string{"male", "female", "none"}
	cg := NewCategoricalGeneType(arr)
	dg := NewDiscreteGeneType(1, 3)
	cog := NewContinuousGeneType(1.0, 2.0)
	geneNames := []string{"gender", "somegene", "ag"}
	geneTypes := []interface{}{*cg, *dg, *cog}
	genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	genes.setRandomGeneValues(geneTypeMap)
	if genes[0].geneName == "gender" && genes[1].geneName == "somegene" && genes[2].geneName == "ag" {
		t.Log("OK")
		for _, v := range genes {
			t.Log(v.geneName, v.geneValue)
		}
	} else {
		t.Fatal("NG")
	}
}

func TestGAResult(t *testing.T) {
	gar := NewGAResult()
	n := 0
	score := 0.5
	gar.recordResult(n, score)
	if gar.result[n] == score {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}
func TestInitialRandomSuggest(t *testing.T) {
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 0.1
	maxGenerations := 5
	selection := "elite"
	selectNum := 2
	crossover := "uniform"
	geneMutation := "perturbation"
	generationChange := "discrete"
	evaluationFunction := "accuracy"
	evaluateHigh := true

	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	geneNames := []string{"cg0", "cg1",
		"dg0", "dg1",
		"cog0", "cog1"}
	geneTypes := []interface{}{*cg0, *cg1, *dg0, *dg1, *cog0, *cog1}
	genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)

	generation := ga.InitialRandomSuggest(geneTypeMap, genes)
	if len(generation.offsprings) == 10 {
		t.Log("OK")
	} else {
		t.Fatal("NG")
	}
}

func TestGetBestScores(t *testing.T) {
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 0.1
	maxGenerations := 5
	selection := "elite"
	selectNum := 2
	crossover := "uniform"
	geneMutation := "perturbation"
	generationChange := "discrete"
	evaluationFunction := "accuracy"
	evaluateHigh := true

	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	geneNames := []string{"cg0", "cg1",
		"dg0", "dg1",
		"cog0", "cog1"}
	geneTypes := []interface{}{*cg0, *cg1, *dg0, *dg1, *cog0, *cog1}
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)

	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)

	gar := NewGAResult()
	for i := 0; i < 100; i++ {
		score := float64(i%50) / 20
		gar.recordResult(i, score)
	}
	bestScoreMap := ga.GetBestScores(*gar, 1)
	if len(bestScoreMap) == 1 {
		t.Log("OK")
	} else {
		t.Log("NG")
	}
}

func TestRandomOffspringMutation(t *testing.T) {
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 1.0
	maxGenerations := 5
	selection := "elite"
	selectNum := 2
	crossover := "uniform"
	geneMutation := "perturbation"
	generationChange := "discrete"
	evaluationFunction := "accuracy"
	evaluateHigh := true
	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	geneNames := []string{"cg0", "cg1",
		"dg0", "dg1",
		"cog0", "cog1"}
	geneTypes := []interface{}{*cg0, *cg1, *dg0, *dg1, *cog0, *cog1}
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)

	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)
	genes := NewGenes(ga.geneNames)
	offspring := ga.randomOffspringMutation(genes)
	t.Log(offspring)
}

func TestSelectElite(t *testing.T) {
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 0.1
	maxGenerations := 10
	selection := "elite"
	selectNum := 2
	crossover := "uniform"
	geneMutation := "random"
	generationChange := "continuous"
	evaluationFunction := "accuracy"
	evaluateHigh := true
	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	geneNames := []string{"cg0", "cg1",
		"dg0", "dg1",
		"cog0", "cog1"}
	geneTypes := []interface{}{*cg0, *cg1, *dg0, *dg1, *cog0, *cog1}
	genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)

	generation := ga.InitialRandomSuggest(geneTypeMap, genes)
	gaResult := NewGAResult()
	for i, o := range generation.offsprings {
		gaResult.recordResult(i, sampleEvaluateFunc0(o))
	}
	bestOffsprings := ga.Select(*gaResult, generation)
	bestScoreMap := ga.GetBestScores(*gaResult, ga.selectNum)
	for _, v := range bestScoreMap {
		if sampleEvaluateFunc0(bestOffsprings[0]) != v && sampleEvaluateFunc0(bestOffsprings[1]) != v {
			t.Fatal("NG")
		}
	}
	t.Log("OK")
}

func TestSelectRoulette(t *testing.T) {
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 0.1
	maxGenerations := 10
	selection := "roulette"
	selectNum := 2
	crossover := "uniform"
	geneMutation := "random"
	generationChange := "continuous"
	evaluationFunction := "accuracy"
	evaluateHigh := true
	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	geneNames := []string{"cg0", "cg1",
		"dg0", "dg1",
		"cog0", "cog1"}
	geneTypes := []interface{}{*cg0, *cg1, *dg0, *dg1, *cog0, *cog1}
	genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)

	generation := ga.InitialRandomSuggest(geneTypeMap, genes)
	gaResult := NewGAResult()
	for i, o := range generation.offsprings {
		gaResult.recordResult(i, sampleEvaluateFunc0(o))
	}
	bestOffsprings := ga.Select(*gaResult, generation)
	t.Log(bestOffsprings)
	t.Log("OK")
}

func TestOptimizeSample0_0(t *testing.T) {
	// target score: 30
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 0.1
	maxGenerations := 200
	selection := "elite"
	selectNum := 2
	crossover := "uniform"
	geneMutation := "random"
	generationChange := "continuous"
	evaluationFunction := "accuracy"
	evaluateHigh := true
	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	cg2 := NewCategoricalGeneType(arr)
	cg3 := NewCategoricalGeneType(arr)
	cg4 := NewCategoricalGeneType(arr)
	cg5 := NewCategoricalGeneType(arr)
	cg6 := NewCategoricalGeneType(arr)
	cg7 := NewCategoricalGeneType(arr)
	cg8 := NewCategoricalGeneType(arr)
	cg9 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	dg2 := NewDiscreteGeneType(0, 4)
	dg3 := NewDiscreteGeneType(0, 4)
	dg4 := NewDiscreteGeneType(0, 3)
	dg5 := NewDiscreteGeneType(0, 3)
	dg6 := NewDiscreteGeneType(0, 2)
	dg7 := NewDiscreteGeneType(0, 2)
	dg8 := NewDiscreteGeneType(0, 1)
	dg9 := NewDiscreteGeneType(0, 1)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	cog2 := NewContinuousGeneType(0, 3.0)
	cog3 := NewContinuousGeneType(0, 3.0)
	cog4 := NewContinuousGeneType(0, 2.0)
	cog5 := NewContinuousGeneType(0, 2.0)
	cog6 := NewContinuousGeneType(0, 1.0)
	cog7 := NewContinuousGeneType(0, 1.0)
	cog8 := NewContinuousGeneType(0, 1.0)
	cog9 := NewContinuousGeneType(0, 1.0)
	geneNames := []string{"cg0", "cg1", "cg2", "cg3", "cg4", "cg5", "cg6", "cg7", "cg8", "cg9",
		"dg0", "dg1", "dg2", "dg3", "dg4", "dg5", "dg6", "dg7", "dg8", "dg9",
		"cog0", "cog1", "cog2", "cog3", "cog4", "cog5", "cog6", "cog7", "cog8", "cog9"}
	geneTypes := []interface{}{*cg0, *cg1, *cg2, *cg3, *cg4, *cg5, *cg6, *cg7, *cg8, *cg9, *dg0, *dg1, *dg2, *dg3, *dg4, *dg5, *dg6, *dg7, *dg8, *dg9, *cog0, *cog1, *cog2, *cog3, *cog4, *cog5, *cog6, *cog7, *cog8, *cog9}
	// genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)

	bestOffspring, bestScore := ga.Optimize(sampleEvaluateFunc0)
	t.Log(bestOffspring, bestScore)
}

func TestOptimizeSample0_1(t *testing.T) {
	// target score: 30
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 0.1
	maxGenerations := 200
	selection := "roulette" // not good
	selectNum := 2
	crossover := "uniform"
	geneMutation := "random"
	generationChange := "discrete"
	evaluationFunction := "accuracy"
	evaluateHigh := true
	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	cg2 := NewCategoricalGeneType(arr)
	cg3 := NewCategoricalGeneType(arr)
	cg4 := NewCategoricalGeneType(arr)
	cg5 := NewCategoricalGeneType(arr)
	cg6 := NewCategoricalGeneType(arr)
	cg7 := NewCategoricalGeneType(arr)
	cg8 := NewCategoricalGeneType(arr)
	cg9 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	dg2 := NewDiscreteGeneType(0, 4)
	dg3 := NewDiscreteGeneType(0, 4)
	dg4 := NewDiscreteGeneType(0, 3)
	dg5 := NewDiscreteGeneType(0, 3)
	dg6 := NewDiscreteGeneType(0, 2)
	dg7 := NewDiscreteGeneType(0, 2)
	dg8 := NewDiscreteGeneType(0, 1)
	dg9 := NewDiscreteGeneType(0, 1)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	cog2 := NewContinuousGeneType(0, 3.0)
	cog3 := NewContinuousGeneType(0, 3.0)
	cog4 := NewContinuousGeneType(0, 2.0)
	cog5 := NewContinuousGeneType(0, 2.0)
	cog6 := NewContinuousGeneType(0, 1.0)
	cog7 := NewContinuousGeneType(0, 1.0)
	cog8 := NewContinuousGeneType(0, 1.0)
	cog9 := NewContinuousGeneType(0, 1.0)
	geneNames := []string{"cg0", "cg1", "cg2", "cg3", "cg4", "cg5", "cg6", "cg7", "cg8", "cg9",
		"dg0", "dg1", "dg2", "dg3", "dg4", "dg5", "dg6", "dg7", "dg8", "dg9",
		"cog0", "cog1", "cog2", "cog3", "cog4", "cog5", "cog6", "cog7", "cog8", "cog9"}
	geneTypes := []interface{}{*cg0, *cg1, *cg2, *cg3, *cg4, *cg5, *cg6, *cg7, *cg8, *cg9, *dg0, *dg1, *dg2, *dg3, *dg4, *dg5, *dg6, *dg7, *dg8, *dg9, *cog0, *cog1, *cog2, *cog3, *cog4, *cog5, *cog6, *cog7, *cog8, *cog9}
	// genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)

	bestOffspring, bestScore := ga.Optimize(sampleEvaluateFunc0)
	t.Log(bestOffspring, bestScore)
}

func TestOptimizeSample1_0(t *testing.T) {
	// target score: 10
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 0.05
	maxGenerations := 300
	selection := "elite"
	selectNum := 2
	crossover := "uniform"
	geneMutation := "random"
	generationChange := "continuous"
	evaluationFunction := "accuracy"
	evaluateHigh := true
	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	cg2 := NewCategoricalGeneType(arr)
	cg3 := NewCategoricalGeneType(arr)
	cg4 := NewCategoricalGeneType(arr)
	cg5 := NewCategoricalGeneType(arr)
	cg6 := NewCategoricalGeneType(arr)
	cg7 := NewCategoricalGeneType(arr)
	cg8 := NewCategoricalGeneType(arr)
	cg9 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	dg2 := NewDiscreteGeneType(0, 4)
	dg3 := NewDiscreteGeneType(0, 4)
	dg4 := NewDiscreteGeneType(0, 3)
	dg5 := NewDiscreteGeneType(0, 3)
	dg6 := NewDiscreteGeneType(0, 2)
	dg7 := NewDiscreteGeneType(0, 2)
	dg8 := NewDiscreteGeneType(0, 1)
	dg9 := NewDiscreteGeneType(0, 1)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	cog2 := NewContinuousGeneType(0, 3.0)
	cog3 := NewContinuousGeneType(0, 3.0)
	cog4 := NewContinuousGeneType(0, 2.0)
	cog5 := NewContinuousGeneType(0, 2.0)
	cog6 := NewContinuousGeneType(0, 1.0)
	cog7 := NewContinuousGeneType(0, 1.0)
	cog8 := NewContinuousGeneType(0, 1.0)
	cog9 := NewContinuousGeneType(0, 1.0)
	geneNames := []string{"cg0", "cg1", "cg2", "cg3", "cg4", "cg5", "cg6", "cg7", "cg8", "cg9",
		"dg0", "dg1", "dg2", "dg3", "dg4", "dg5", "dg6", "dg7", "dg8", "dg9",
		"cog0", "cog1", "cog2", "cog3", "cog4", "cog5", "cog6", "cog7", "cog8", "cog9"}
	geneTypes := []interface{}{*cg0, *cg1, *cg2, *cg3, *cg4, *cg5, *cg6, *cg7, *cg8, *cg9, *dg0, *dg1, *dg2, *dg3, *dg4, *dg5, *dg6, *dg7, *dg8, *dg9, *cog0, *cog1, *cog2, *cog3, *cog4, *cog5, *cog6, *cog7, *cog8, *cog9}
	// genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)

	bestOffspring, bestScore := ga.Optimize(sampleEvaluateFunc1)
	t.Log(bestOffspring, bestScore)
}

func TestOptimizeSample1_1(t *testing.T) {
	// target score: 10
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 0.1
	maxGenerations := 100
	selection := "elite"
	selectNum := 2
	crossover := "uniform"
	geneMutation := "perturbation" // not good
	generationChange := "discrete"
	evaluationFunction := "accuracy"
	evaluateHigh := true
	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	cg2 := NewCategoricalGeneType(arr)
	cg3 := NewCategoricalGeneType(arr)
	cg4 := NewCategoricalGeneType(arr)
	cg5 := NewCategoricalGeneType(arr)
	cg6 := NewCategoricalGeneType(arr)
	cg7 := NewCategoricalGeneType(arr)
	cg8 := NewCategoricalGeneType(arr)
	cg9 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	dg2 := NewDiscreteGeneType(0, 4)
	dg3 := NewDiscreteGeneType(0, 4)
	dg4 := NewDiscreteGeneType(0, 3)
	dg5 := NewDiscreteGeneType(0, 3)
	dg6 := NewDiscreteGeneType(0, 2)
	dg7 := NewDiscreteGeneType(0, 2)
	dg8 := NewDiscreteGeneType(0, 1)
	dg9 := NewDiscreteGeneType(0, 1)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	cog2 := NewContinuousGeneType(0, 3.0)
	cog3 := NewContinuousGeneType(0, 3.0)
	cog4 := NewContinuousGeneType(0, 2.0)
	cog5 := NewContinuousGeneType(0, 2.0)
	cog6 := NewContinuousGeneType(0, 1.0)
	cog7 := NewContinuousGeneType(0, 1.0)
	cog8 := NewContinuousGeneType(0, 1.0)
	cog9 := NewContinuousGeneType(0, 1.0)
	geneNames := []string{"cg0", "cg1", "cg2", "cg3", "cg4", "cg5", "cg6", "cg7", "cg8", "cg9",
		"dg0", "dg1", "dg2", "dg3", "dg4", "dg5", "dg6", "dg7", "dg8", "dg9",
		"cog0", "cog1", "cog2", "cog3", "cog4", "cog5", "cog6", "cog7", "cog8", "cog9"}
	geneTypes := []interface{}{*cg0, *cg1, *cg2, *cg3, *cg4, *cg5, *cg6, *cg7, *cg8, *cg9, *dg0, *dg1, *dg2, *dg3, *dg4, *dg5, *dg6, *dg7, *dg8, *dg9, *cog0, *cog1, *cog2, *cog3, *cog4, *cog5, *cog6, *cog7, *cog8, *cog9}
	// genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)
	bestOffspring, bestScore := ga.Optimize(sampleEvaluateFunc1)
	t.Log(bestOffspring, bestScore)
}

func TestOptimizeSample1_2(t *testing.T) {
	// target score: 0
	numGenes := 5
	numOffsprings := 10
	geneMutationProbability := 0.05
	offspringMutationProbability := 0.1
	maxGenerations := 100
	selection := "elite"
	selectNum := 2
	crossover := "uniform"
	geneMutation := "random"
	generationChange := "discrete"
	evaluationFunction := "accuracy"
	evaluateHigh := false
	arr := []string{"true", "false", "tf"}
	cg0 := NewCategoricalGeneType(arr)
	cg1 := NewCategoricalGeneType(arr)
	cg2 := NewCategoricalGeneType(arr)
	cg3 := NewCategoricalGeneType(arr)
	cg4 := NewCategoricalGeneType(arr)
	cg5 := NewCategoricalGeneType(arr)
	cg6 := NewCategoricalGeneType(arr)
	cg7 := NewCategoricalGeneType(arr)
	cg8 := NewCategoricalGeneType(arr)
	cg9 := NewCategoricalGeneType(arr)
	dg0 := NewDiscreteGeneType(0, 5)
	dg1 := NewDiscreteGeneType(0, 5)
	dg2 := NewDiscreteGeneType(0, 4)
	dg3 := NewDiscreteGeneType(0, 4)
	dg4 := NewDiscreteGeneType(0, 3)
	dg5 := NewDiscreteGeneType(0, 3)
	dg6 := NewDiscreteGeneType(0, 2)
	dg7 := NewDiscreteGeneType(0, 2)
	dg8 := NewDiscreteGeneType(0, 1)
	dg9 := NewDiscreteGeneType(0, 1)
	cog0 := NewContinuousGeneType(0, 4.0)
	cog1 := NewContinuousGeneType(0, 4.0)
	cog2 := NewContinuousGeneType(0, 3.0)
	cog3 := NewContinuousGeneType(0, 3.0)
	cog4 := NewContinuousGeneType(0, 2.0)
	cog5 := NewContinuousGeneType(0, 2.0)
	cog6 := NewContinuousGeneType(0, 1.0)
	cog7 := NewContinuousGeneType(0, 1.0)
	cog8 := NewContinuousGeneType(0, 1.0)
	cog9 := NewContinuousGeneType(0, 1.0)
	geneNames := []string{"cg0", "cg1", "cg2", "cg3", "cg4", "cg5", "cg6", "cg7", "cg8", "cg9",
		"dg0", "dg1", "dg2", "dg3", "dg4", "dg5", "dg6", "dg7", "dg8", "dg9",
		"cog0", "cog1", "cog2", "cog3", "cog4", "cog5", "cog6", "cog7", "cog8", "cog9"}
	geneTypes := []interface{}{*cg0, *cg1, *cg2, *cg3, *cg4, *cg5, *cg6, *cg7, *cg8, *cg9, *dg0, *dg1, *dg2, *dg3, *dg4, *dg5, *dg6, *dg7, *dg8, *dg9, *cog0, *cog1, *cog2, *cog3, *cog4, *cog5, *cog6, *cog7, *cog8, *cog9}
	// genes := NewGenes(geneNames)
	geneTypeMap := NewGeneTypeMap(geneNames, geneTypes)
	ga := NewGA(numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap)
	bestOffspring, bestScore := ga.Optimize(sampleEvaluateFunc1)
	t.Log(bestOffspring, bestScore)
}
