package cmaes

import (
	"errors"
	"math"
	"math/rand"
	"sort"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

type Solution struct {
	// Params is a parameter transformed to N(m, σ^2 C) from Z.
	Params []float64
	// Value represents an evaluation value.
	Value float64
}

// Optimizer is CMA-ES stochastic optimizer class with ask-and-tell interface.
type Optimizer struct {
	mean  *mat.VecDense
	sigma float64
	c     *mat.SymDense

	dim     int
	mu      int
	muEff   float64
	popsize int
	cc      float64
	c1      float64
	cmu     float64
	cSigma  float64
	dSigma  float64
	cm      float64
	chiN    float64
	pSigma  *mat.VecDense
	pc      *mat.VecDense
	weights *mat.VecDense

	bounds        mat.Matrix
	maxReSampling int

	rng *rand.Rand
	g   int
}

// NewOptimizer returns an optimizer object based on CMA-ES.
func NewOptimizer(mean []float64, sigma float64, opts ...OptimizerOption) (*Optimizer, error) {
	if sigma <= 0 {
		return nil, errors.New("sigma should be non-zero positive number")
	}
	dim := len(mean)
	popsize := 4 + int(math.Floor(3*math.Log(float64(dim))))
	mu := popsize / 2

	sumWeightsPrimeBeforeMu := 0.
	sumWeightsPrimeSquareBeforeMu := 0.
	sumWeightsPrimeAfterMu := 0.
	sumWeightsPrimeSquareAfterMu := 0.
	weightsPrime := make([]float64, popsize)
	weightsPrimePositiveSum := 0.0
	weightsPrimeNegativeSum := 0.0
	for i := 0; i < popsize; i++ {
		wp := math.Log((float64(popsize)+1)/2) - math.Log(float64(i+1))
		weightsPrime[i] = wp

		if i < mu {
			sumWeightsPrimeBeforeMu += wp
			sumWeightsPrimeSquareBeforeMu += math.Pow(wp, 2)
		} else {
			sumWeightsPrimeAfterMu += weightsPrime[i]
			sumWeightsPrimeSquareAfterMu += math.Pow(wp, 2)
		}

		if wp > 0 {
			weightsPrimePositiveSum += wp
		} else {
			weightsPrimeNegativeSum -= wp
		}
	}
	muEff := math.Pow(sumWeightsPrimeBeforeMu, 2) / sumWeightsPrimeSquareBeforeMu
	muEffMinus := math.Pow(sumWeightsPrimeAfterMu, 2) / sumWeightsPrimeSquareAfterMu

	alphaCov := 2.0
	// learning rate for the rank-one update
	c1 := alphaCov / (math.Pow(float64(dim)+1.3, 2) + muEff)
	// learning rate for the rank-μ update
	cmu := math.Min(
		1-c1,
		alphaCov*(muEff-2+1/muEff)/(math.Pow(float64(dim+2), 2)+alphaCov*muEff/2),
	)
	if c1+cmu > 1 {
		return nil, errors.New("invalid learning rate for the rank-one and rank-μ update")
	}

	alphaMin := math.Min(
		1+c1/cmu,                   // α_μ-
		1+(2*muEffMinus)/(muEff+2), // α_μ_eff-
	)
	alphaMin = math.Min(alphaMin, (1-c1-cmu)/(float64(dim)*cmu)) // α_{pos_def}^{minus}

	weights := make([]float64, popsize)
	for i := 0; i < popsize; i++ {
		if weightsPrime[i] > 0 {
			weights[i] = 1 / weightsPrimePositiveSum * weightsPrime[i]
		} else {
			weights[i] = alphaMin / weightsPrimeNegativeSum * weightsPrime[i]
		}
	}
	cm := 1.0

	// learning rate for the cumulation for the step-size control (eq.55)
	cSigma := (muEff + 2) / (float64(dim) + muEff + 5)
	dSigma := 1 + 2*math.Max(0, math.Sqrt((muEff-1)/(float64(dim)+1))-1) + cSigma
	if cSigma >= 1 {
		return nil, errors.New("invalid learning rate for cumulation for the ste-size control")
	}

	// learning rate for cumulation for the rank-one update (eq.56)
	cc := (4 + muEff/float64(dim)) / (float64(dim) + 4 + 2*muEff/float64(dim))
	if cc > 1 {
		return nil, errors.New("invalid learning rate for cumulation for the rank-one update")
	}

	chiN := math.Sqrt(float64(dim)) * (1.0 - (1.0 / (4.0 * float64(dim))) + 1.0/(21.0*(math.Pow(float64(dim), 2))))

	cma := &Optimizer{
		mean:          mat.NewVecDense(dim, mean),
		sigma:         sigma,
		c:             initC(dim),
		dim:           dim,
		popsize:       popsize,
		mu:            mu,
		muEff:         muEff,
		cc:            cc,
		c1:            c1,
		cmu:           cmu,
		cSigma:        cSigma,
		dSigma:        dSigma,
		cm:            cm,
		chiN:          chiN,
		pSigma:        mat.NewVecDense(dim, make([]float64, dim)),
		pc:            mat.NewVecDense(dim, make([]float64, dim)),
		weights:       mat.NewVecDense(popsize, weights),
		bounds:        nil,
		maxReSampling: 100,
		rng:           rand.New(rand.NewSource(0)),
		g:             0,
	}

	for _, opt := range opts {
		opt(cma)
	}
	return cma, nil
}

// Generation is incremented when a multivariate normal distribution is updated.
func (o *Optimizer) Generation() int {
	return o.g
}

// PopulationSize returns the population size.
func (o *Optimizer) PopulationSize() int {
	return o.popsize
}

// Ask a next parameter.
func (o *Optimizer) Ask() ([]float64, error) {
	x, err := o.sampleSolution()
	if err != nil {
		return nil, err
	}
	for i := 0; i < o.maxReSampling; i++ {
		if o.isFeasible(x) {
			return x.RawVector().Data, nil
		}
		x, err = o.sampleSolution()
		if err != nil {
			return nil, err
		}
	}
	err = o.repairInfeasibleParams(x)
	if err != nil {
		return nil, err
	}
	return x.RawVector().Data, nil
}

func (o *Optimizer) isFeasible(values *mat.VecDense) bool {
	if o.bounds == nil {
		return true
	}
	if values.Len() != o.dim {
		return false
	}
	for i := 0; i < o.dim; i++ {
		v := values.AtVec(i)
		if !(o.bounds.At(i, 0) < v && o.bounds.At(i, 1) > v) {
			return false
		}
	}
	return true
}

func (o *Optimizer) repairInfeasibleParams(values *mat.VecDense) error {
	if o.bounds == nil {
		return nil
	}
	if values.Len() != o.dim {
		return errors.New("invalid matrix size")
	}

	for i := 0; i < o.dim; i++ {
		v := values.AtVec(i)
		if o.bounds.At(i, 0) > v {
			values.SetVec(i, o.bounds.At(i, 0))
		}
		if o.bounds.At(i, 1) < v {
			values.SetVec(i, o.bounds.At(i, 1))
		}
	}
	return nil
}

func (o *Optimizer) sampleSolution() (*mat.VecDense, error) {
	// TODO(o-bata): Cache B and D
	var eigsym mat.EigenSym
	ok := eigsym.Factorize(o.c, true)
	if !ok {
		return nil, errors.New("symmetric eigendecomposition failed")
	}

	var b mat.Dense
	eigsym.VectorsTo(&b)
	d := make([]float64, o.dim)
	eigsym.Values(d) // d^2
	floatsSqrtTo(d)  // d

	z := make([]float64, o.dim)
	for i := 0; i < o.dim; i++ {
		z[i] = o.rng.NormFloat64()
	}

	var bd mat.Dense
	bd.Mul(&b, mat.NewDiagDense(o.dim, d))

	values := mat.NewVecDense(o.dim, z) // ~ N(0, I)
	values.MulVec(&bd, values)          // ~ N(0, C)
	values.ScaleVec(o.sigma, values)    // ~ N(0, σ^2 C)
	values.AddVec(values, o.mean)       // ~ N(m, σ^2 C)
	return values, nil
}

// Tell evaluation values.
func (o *Optimizer) Tell(solutions []*Solution) error {
	if len(solutions) != o.popsize {
		return errors.New("must tell popsize-length solutions")
	}

	o.g++
	sort.Slice(solutions, func(i, j int) bool {
		return solutions[i].Value < solutions[j].Value
	})

	var eigsym mat.EigenSym
	ok := eigsym.Factorize(o.c, true)
	if !ok {
		return errors.New("symmetric eigendecomposition failed")
	}

	var b mat.Dense
	eigsym.VectorsTo(&b)
	d := make([]float64, o.dim)
	eigsym.Values(d) // d^2
	floatsSqrtTo(d)  // d

	yk := mat.NewDense(o.popsize, o.dim, nil)
	for i := 0; i < o.popsize; i++ {
		xi := solutions[i].Params           // ~ N(m, σ^2 C)
		xiSubMean := make([]float64, o.dim) // ~ N(0, σ^2 C)
		floats.SubTo(xiSubMean, xi, o.mean.RawVector().Data)
		yk.SetRow(i, xiSubMean)
	}
	yk.Scale(1/o.sigma, yk) // ~ N(0, C)

	// Selection and recombination
	ydotw := mat.NewDense(o.mu, o.dim, nil)
	ydotw.Copy(yk.Slice(0, o.mu, 0, o.dim))
	weightsmu := stackvec(o.dim, o.mu, o.weights)
	ydotw.MulElem(ydotw, weightsmu.T())

	yw := sumColumns(ydotw.T())
	meandiff := mat.NewVecDense(o.dim, nil)
	meandiff.CopyVec(yw)
	meandiff.ScaleVec(o.cm*o.sigma, meandiff)
	o.mean.AddVec(o.mean, meandiff)

	// Step-size control
	dinv := mat.NewDiagDense(o.dim, arrinv(d))
	c2 := mat.NewDense(o.dim, o.dim, nil)
	c2.Product(&b, dinv, b.T()) // C^(-1/2) = B D^(-1) B^T

	c2yw := mat.NewDense(o.dim, 1, nil)
	c2yw.Product(c2, yw)
	c2yw.Scale(math.Sqrt(o.cSigma*(2-o.cSigma)*o.muEff), c2yw)
	o.pSigma.ScaleVec(1-o.cSigma, o.pSigma)
	o.pSigma.AddVec(o.pSigma, mat.NewVecDense(o.dim, c2yw.RawMatrix().Data))

	normPSigma := mat.Norm(o.pSigma, 2)
	o.sigma *= math.Exp((o.cSigma / o.dSigma) * (normPSigma/o.chiN - 1))

	hSigmaCondLeft := normPSigma / math.Sqrt(
		1-math.Pow(1-o.cSigma, float64(2*(o.g+1))))
	hSigmaCondRight := (1.4 + 2/float64(o.dim+1)) * o.chiN
	hSigma := 0.0
	if hSigmaCondLeft < hSigmaCondRight {
		hSigma = 1.0
	}

	// eq.45
	o.pc.ScaleVec(1-o.cc, o.pc)
	o.pc.AddScaledVec(o.pc, hSigma*math.Sqrt(o.cc*(2-o.cc)*o.muEff), yw)

	// eq.46
	wio := mat.NewVecDense(o.weights.Len(), nil)
	wio.CopyVec(o.weights)
	c2yk := mat.NewDense(o.dim, o.popsize, nil)
	c2yk.Product(c2, yk.T())
	wio.MulElemVec(wio, vecapply(o.weights, func(i int, a float64) float64 {
		if a > 0 {
			return 1.0
		}
		c2xinorm := mat.Norm(c2yk.ColView(i), 2)
		return float64(o.dim) / math.Pow(c2xinorm, 2)
	}))

	deltaHSigma := (1 - hSigma) * o.cc * (2 - o.cc)
	if deltaHSigma > 1 {
		panic("invalid delta_h_sigma")
	}

	// eq.47
	rankOne := mat.NewSymDense(o.dim, nil)
	rankOne.SymOuterK(1.0, o.pc)

	rankMu := mat.NewSymDense(o.dim, nil)
	for i := 0; i < o.popsize; i++ {
		wi := wio.AtVec(i)
		yi := yk.RowView(i)
		s := mat.NewSymDense(o.dim, nil)
		s.SymOuterK(wi, yi)
		rankMu.AddSym(rankMu, s)
	}

	o.c.ScaleSym(1+o.c1*deltaHSigma-o.c1-o.cmu*mat.Sum(o.weights), o.c)
	rankOne.ScaleSym(o.c1, rankOne)
	rankMu.ScaleSym(o.cmu, rankMu)
	o.c.AddSym(o.c, rankOne)
	o.c.AddSym(o.c, rankMu)

	// Avoid eigendecomposition error by arithmetic overflow
	o.c.AddSym(o.c, initMinC(o.dim))
	return nil
}

// OptimizerOption is a type of the function to customizing CMA-ES.
type OptimizerOption func(*Optimizer)

// OptimizerOptionSeed sets seed number.
func OptimizerOptionSeed(seed int64) OptimizerOption {
	return func(cma *Optimizer) {
		cma.rng = rand.New(rand.NewSource(seed))
	}
}

// OptimizerOptionMaxReSampling sets a number of max re-sampling.
func OptimizerOptionMaxReSampling(n int) OptimizerOption {
	return func(cma *Optimizer) {
		cma.maxReSampling = n
	}
}

// OptimizerOptionBounds sets the range of parameters.
func OptimizerOptionBounds(bounds *mat.Dense) OptimizerOption {
	row, column := bounds.Dims()
	if column != 2 {
		panic("invalid matrix size")
	}

	return func(cma *Optimizer) {
		if row != cma.dim {
			panic("invalid dimensions")
		}
		cma.bounds = bounds
	}
}
