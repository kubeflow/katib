package genetic_algorithm

// for gene with continuous variables in float64s
type ContinuousGeneType struct {
	min float64
	max float64
}

func NewContinuousGeneType(min, max float64) *ContinuousGeneType {
	return &ContinuousGeneType{min, max}
}

func (cg *ContinuousGeneType) generateRandom() float64 {
	return getContinuousRandom(cg.min, cg.max)
}

func (cg *ContinuousGeneType) setGeneRandomValue(g *Gene) {
	g.setGeneValue(cg.generateRandom())
}

// for gene with discrete variables in int
type DiscreteGeneType struct {
	min int
	max int
}

func NewDiscreteGeneType(min, max int) *DiscreteGeneType {
	return &DiscreteGeneType{min, max}
}

func (dg *DiscreteGeneType) generateRandom() int {
	return getDiscreteRandom(dg.min, dg.max)
}

func (dg *DiscreteGeneType) setGeneRandomValue(g *Gene) {
	g.setGeneValue(dg.generateRandom())
}

// for gene with categorical variables in string
type CategoricalGeneType struct {
	category []string
}

func NewCategoricalGeneType(category []string) *CategoricalGeneType {
	return &CategoricalGeneType{category}
}

func (cg *CategoricalGeneType) generateRandom() int {
	return getCategoricalRandom(cg.category)
}

func (cg *CategoricalGeneType) setGeneRandomValue(g *Gene) {
	g.setGeneValue(cg.category[cg.generateRandom()])
}

type setGeneRandomValueInterface interface {
	setGeneRandomValue(g *Gene)
}

// map for gene name and gene type
type GeneTypeMap map[string]interface{} // map[geneName]geneType

func NewGeneTypeMap(geneNames []string, geneTypes []interface{}) GeneTypeMap {
	geneTypeMap := make(map[string]interface{})
	for i := 0; i < len(geneNames); i++ {
		geneTypeMap[geneNames[i]] = geneTypes[i]
	}
	return geneTypeMap
}

// gene struct; one name and one value per one gene
type Gene struct {
	geneName  string      // name of the gene
	geneValue interface{} // value of the gene; must fit with the gene type
}

func NewGene(geneName string) *Gene {
	return &Gene{geneName, ""}
}

func (g *Gene) setGeneValue(v interface{}) {
	g.geneValue = v
}

func (g *Gene) setRandomGeneValue(s setGeneRandomValueInterface) {
	s.setGeneRandomValue(g)
}

// set random value depending on the gene type
func setRandomGeneValueWithType(gene *Gene, geneTypeMap GeneTypeMap) {
	switch s := geneTypeMap[gene.geneName].(type) {
	case CategoricalGeneType:
		gene.setRandomGeneValue(&s)
	case DiscreteGeneType:
		gene.setRandomGeneValue(&s)
	case ContinuousGeneType:
		gene.setRandomGeneValue(&s)
	}
}

// a slice of genes
type Genes []Gene

func NewGenes(geneNames []string) Genes {
	genes := []Gene{}
	for i := 0; i < len(geneNames); i++ {
		gene := NewGene(geneNames[i])
		genes = append(genes, *gene)
	}
	return genes
}

func (genes Genes) setRandomGeneValues(geneTypeMap GeneTypeMap) {
	for i := 0; i < len(genes); i++ {
		setRandomGeneValueWithType(&genes[i], geneTypeMap)
	}
}

// an offspring; a map of genes for the offspring; contains a map of genes with gene name and its value
type Offspring struct {
	GeneMap map[string]interface{} // map[geneName]geneValue
}

func NewOffspring(genes Genes) *Offspring {
	geneMap := make(map[string]interface{})
	for _, v := range genes {
		geneMap[v.geneName] = v.geneValue
	}
	return &Offspring{geneMap}
}

// a generation; a set of offsprings
type Generation struct {
	genNum     int // starts from 0
	offsprings []Offspring
}

func NewGeneration(genNum int) *Generation {
	return &Generation{genNum, nil}
}

func (generation *Generation) setGeneration(offsprings []Offspring) {
	generation.offsprings = offsprings
}
