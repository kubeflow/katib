package genetic_algorithm

import (
	"sort"
)

// struct for genetic algorithm
type GA struct {
	numGenes                     int     // number of genes per offspring
	numOffsprings                int     // number of offspring per generation
	geneMutationProbability      float64 // below 1.0; probability of gene mutation
	offspringMutationProbability float64 // below 1.0; probability of random offspring mutation
	maxGenerations               int     // max iterations
	selection                    string  // either roulette or elite; recommend elite
	selectNum                    int     // must be >=2; recommend 2 or 3
	crossover                    string  // only uniform
	geneMutation                 string  // either random or perturbation; recommend random
	generationChange             string  // either discrete or continuous
	evaluationFunction           string  // evaluation function
	evaluateHigh                 bool    // true for higher is better; false for lower is better
	geneNames                    []string
	geneTypeMap                  GeneTypeMap
}

func NewGA(numGenes int, numOffsprings int, geneMutationProbability float64, offspringMutationProbability float64, maxGenerations int, selection string, selectNum int, crossover string, geneMutation string, generationChange string, evaluationFunction string, evaluateHigh bool, geneNames []string, geneTypeMap GeneTypeMap) *GA {
	return &GA{numGenes, numOffsprings, geneMutationProbability, offspringMutationProbability, maxGenerations, selection, selectNum, crossover, geneMutation, generationChange, evaluationFunction, evaluateHigh, geneNames, geneTypeMap}
}

// randomly generate a generation; for the initial generation to start with
func (ga *GA) InitialRandomSuggest(geneTypeMap GeneTypeMap, genes Genes) Generation {
	offsprings := []Offspring{}
	for i := 0; i < ga.numOffsprings; i++ {
		genes.setRandomGeneValues(geneTypeMap)
		offspring := NewOffspring(genes)
		offsprings = append(offsprings, *offspring)
	}
	generation := NewGeneration(0)
	generation.setGeneration(offsprings)
	return *generation
}

// genetic algorithm optimizer; optimizes against evaluateFunc with GA variables
func (ga *GA) Optimize(evaluateFunc func(offspring Offspring) float64) (Offspring, float64) {
	genes := NewGenes(ga.geneNames)
	bestOffspring := *NewOffspring(genes)
	var bestScore float64

	// initial generation and the scores of its offsprings
	generation := ga.InitialRandomSuggest(ga.geneTypeMap, genes)
	gaResult := NewGAResult()
	for i, o := range generation.offsprings {
		gaResult.recordResult(i, evaluateFunc(o))
	}
	bestOffsprings := ga.Select(gaResult, generation)

	// run genetic algorithm for maxGeneration iterations to optimize genes
	for g := 0; g < ga.maxGenerations; g++ {
		generation = ga.Crossover(bestOffsprings, generation)
		gaResult = NewGAResult()
		for i, o := range generation.offsprings {
			gaResult.recordResult(i, evaluateFunc(o))
		}
		bestOffsprings = ga.Select(gaResult, generation)
	}

	// best genes after optimization
	bestScoreMap := ga.GetBestScores(gaResult, 1)
	for k, v := range bestScoreMap {
		bestScore = v
		bestOffspring = generation.offsprings[k]
	}

	// returns the best optimized offspring and its score
	return bestOffspring, bestScore
}

// score of a gene
type GAScore struct {
	ScoreId int
	Score   float64
}

// a slice of evaluation result for each offspring
type GAResult []GAScore

func NewGAResult() GAResult {
	gar := GAResult{}
	return gar
}

func (gar *GAResult) recordResult(n int, score float64) {
	gas := GAScore{n, score}
	*gar = append(*gar, gas)
}

func (gar GAResult) Len() int {
	return len(gar)
}

func (gar GAResult) Swap(i, j int) {
	gar[i], gar[j] = gar[j], gar[i]
}

func (gar GAResult) Less(i, j int) bool {
	return gar[i].Score < gar[j].Score
}

func (gar GAResult) SortScore(more bool) GAResult {
	if more {
		sort.Sort(sort.Reverse(gar))
	} else {
		sort.Sort(gar)
	}
	return gar
}

// get the best score from ga result
func (ga *GA) GetBestScores(gaResult GAResult, numScores int) map[int]float64 {
	bestScoreMap := make(map[int]float64)
	gaSorted := gaResult.SortScore(ga.evaluateHigh)
	for _, v := range gaSorted {
		bestScoreMap[v.ScoreId] = v.Score
		if len(bestScoreMap) >= numScores {
			break
		}
	}
	return bestScoreMap
}

// selection; either elite or roulette is available
func (ga *GA) Select(gaResult GAResult, generation Generation) []Offspring {
	if ga.selection == "elite" {
		return ga.selectElite(gaResult, generation)
	} else if ga.selection == "roulette" {
		return ga.selectRoulette(gaResult, generation)
	} else {
		return nil
	}
}

// select by randomly choosing with probabilities depending on ga result
func (ga *GA) selectRoulette(gaResult GAResult, generation Generation) []Offspring {
	bestOffsprings := []Offspring{}
	roulette := []float64{}
	positions := []int{}
	var sum float64
	var choice float64
	chosen := make(map[int]bool)
	if ga.evaluateHigh {
		// if higher is better, higher result gets higher probability
		for k, v := range gaResult {
			if v.Score == 0 {
				continue
			}
			sum = sum + v.Score
			positions = append(positions, k)
			roulette = append(roulette, sum)
		}
	} else {
		// if lower is better, lower result gets higher probability
		var rsum float64
		for _, v := range gaResult {
			rsum = rsum + v.Score
		}
		for k, v := range gaResult {
			positions = append(positions, k)
			sum = sum + rsum - v.Score
			roulette = append(roulette, sum)
		}
	}
	// randomly choose depending on the probability
	for i := 0; i < 1; {
		choice = getContinuousRandom(0, sum)
		for k, v := range roulette {
			if choice < v {
				if _, ok := chosen[k]; ok {
					break
				}
				chosen[k] = true
				bestOffsprings = append(bestOffsprings, generation.offsprings[positions[k]])
				break
			}
		}
		if len(bestOffsprings) == ga.selectNum {
			break
		}
	}
	return bestOffsprings
}

// select depending on the highest results
func (ga *GA) selectElite(gaResult GAResult, generation Generation) []Offspring {
	bestScoreMap := ga.GetBestScores(gaResult, ga.selectNum)
	bestOffsprings := []Offspring{}
	for k, _ := range bestScoreMap {
		bestOffsprings = append(bestOffsprings, generation.offsprings[k])
	}
	return bestOffsprings
}

// offspring mutation; only random mutation is available
func (ga *GA) OffspringMutation(genes Genes) *Offspring {
	return ga.randomOffspringMutation(genes)
}

// generates an offspring with random genes
func (ga *GA) randomOffspringMutation(genes Genes) *Offspring {
	genes.setRandomGeneValues(ga.geneTypeMap)
	return NewOffspring(genes)
}

// gene mutation; random or perturbation
func (ga *GA) GeneMutation(gene *Gene) {
	if ga.geneMutation == "random" {
		ga.randomGeneMutation(gene)
	} else if ga.geneMutation == "perturbation" {
		ga.perturbationGeneMutation(gene)
	}
}

// generates a random gene
func (ga *GA) randomGeneMutation(gene *Gene) {
	setRandomGeneValueWithType(gene, ga.geneTypeMap)
}

// gives a little change on a discrete or continuous gene; random for categorical gene
func (ga *GA) perturbationGeneMutation(gene *Gene) {
	switch s := ga.geneTypeMap[gene.geneName].(type) {
	case CategoricalGeneType:
		// random for categorical
		gene.setRandomGeneValue(&s)
	case DiscreteGeneType:
		// add 1 or subtract 1 from the original gene for discrete gene
		// should be in between min and max of the gene
		var r int
		if gene.geneValue == s.min {
			r = 1
		} else if gene.geneValue == s.max {
			r = -1
		} else {
			r = getDiscreteRandomExcludeA(-1, 1, 0)
		}
		gene.setGeneValue(gene.geneValue.(int) + r)
	case ContinuousGeneType:
		// multiply the original gene by 0.9 to 1.1 for continuous gene
		// should be in between min and max of the gene
		for i := 0; i < 1; {
			r := gene.geneValue.(float64) * getContinuousRandom(0.9, 1.1)
			if r >= s.min && r <= s.max {
				gene.setGeneValue(r)
				break
			}
		}
	}
}

// cross over; generates next generation from the selected offsprings as parents
func (ga *GA) Crossover(bestOffsprings []Offspring, generation Generation) Generation {
	newOffspringNum := 0
	if ga.generationChange == "discrete" {
		// all the offsprings in the next generation are new; no former offsprings survive
		newOffspringNum = ga.numOffsprings
	} else if ga.generationChange == "continuous" {
		// best offsprings survive to the next generation
		newOffspringNum = ga.numOffsprings - len(bestOffsprings)
	}
	// only uniform crossover for now
	if ga.crossover == "uniform" {
		return ga.crossoverUniform(bestOffsprings, generation, newOffspringNum)
	} else {
		return ga.crossoverUniform(bestOffsprings, generation, newOffspringNum)
	}
}

// uniform crossover; just randomly generates next offsprings using parents' genes
func (ga *GA) crossoverUniform(bestOffsprings []Offspring, generation Generation, newOffspringNum int) Generation {
	offsprings := []Offspring{}
	nextGeneration := NewGeneration(generation.genNum + 1)
	offsprings = append(offsprings, bestOffsprings...)
	for i := 0; i < newOffspringNum; i++ {
		// chooses a parent pair
		choice0 := getDiscreteRandom(0, len(bestOffsprings)-1)
		choice1 := getDiscreteRandomExcludeA(0, len(bestOffsprings)-1, choice0)
		parents := []Offspring{bestOffsprings[choice0], bestOffsprings[choice1]}
		genes0 := NewGenes(ga.geneNames)
		genes1 := NewGenes(ga.geneNames)
		if getContinuousRandom(0.0, 1.0) < ga.offspringMutationProbability {
			// offspring mutation; generates offspring randomly, without dependency on the parents
			offspring := ga.OffspringMutation(genes0)
			offsprings = append(offsprings, *offspring)
		} else {
			// generate two offsprings from a pair of parents
			for _, g := range ga.geneNames {
				gene0 := NewGene(g)
				gene1 := NewGene(g)
				// randomly chooses which parent to inherit for each gene
				choice0, choice1 = shuffleTwo(0, 1)
				gene0.setGeneValue(parents[choice0].GeneMap[g])
				gene1.setGeneValue(parents[choice1].GeneMap[g])
				if getContinuousRandom(0.0, 1.0) < ga.geneMutationProbability {
					// gene mutation; mutates the inherited genes
					ga.GeneMutation(gene0)
					ga.GeneMutation(gene1)
				}
				genes0 = append(genes0, *gene0)
				genes1 = append(genes1, *gene1)
			}
			offspring0 := NewOffspring(genes0)
			offsprings = append(offsprings, *offspring0)
			offspring1 := NewOffspring(genes1)
			offsprings = append(offsprings, *offspring1)
		}
		if len(offsprings) >= ga.numOffsprings {
			break
		}
	}
	// generates the next generation
	nextGeneration.setGeneration(offsprings)
	if len(nextGeneration.offsprings) > ga.numOffsprings {
		nextGeneration.offsprings = nextGeneration.offsprings[:ga.numOffsprings]
	}
	return *nextGeneration
}
