package tpe

import (
	"math"
	"math/rand"
	"sort"
	"sync"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/internal/random"
	"gonum.org/v1/gonum/floats"
)

const eps = 1e-12

// FuncGamma is a type of gamma function.
type FuncGamma func(int) int

// FuncWeights is a type of weights function.
type FuncWeights func(int) []float64

// DefaultGamma is a default gamma function.
func DefaultGamma(x int) int {
	a := int(math.Ceil(0.1 * float64(x)))
	if a > 25 {
		return 25
	}
	return a
}

// HyperoptDefaultGamma is a default gamma function of Hyperopt.
func HyperoptDefaultGamma(x int) int {
	a := int(math.Ceil(0.25 * float64(x)))
	if a > 25 {
		return a
	}
	return 25
}

// DefaultWeights is a default weights function.
func DefaultWeights(x int) []float64 {
	if x == 0 {
		return []float64{}
	} else if x < 25 {
		return ones1d(x)
	} else {
		ramp := linspace(1.0/float64(x), 1.0, x-25, true)
		flat := ones1d(25)
		return append(ramp, flat...)
	}
}

var _ goptuna.Sampler = &Sampler{}

// Sampler returns the next search points by using TPE.
type Sampler struct {
	numberOfStartupTrials int
	numberOfEICandidates  int
	gamma                 FuncGamma
	params                ParzenEstimatorParams
	rng                   *rand.Rand
	randomSampler         *goptuna.RandomSearchSampler
	mu                    sync.Mutex
}

// NewSampler returns the TPE sampler.
func NewSampler(opts ...SamplerOption) *Sampler {
	sampler := &Sampler{
		numberOfStartupTrials: 10,
		numberOfEICandidates:  24,
		gamma:                 DefaultGamma,
		params: ParzenEstimatorParams{
			ConsiderPrior:     true,
			PriorWeight:       1.0,
			ConsiderMagicClip: true,
			ConsiderEndpoints: false,
			Weights:           DefaultWeights,
		},
		rng:           rand.New(rand.NewSource(0)),
		randomSampler: goptuna.NewRandomSearchSampler(),
	}

	for _, opt := range opts {
		opt(sampler)
	}
	return sampler
}

func (s *Sampler) splitObservationPairs(
	configVals []float64,
	lossVals [][2]float64,
) ([]float64, []float64) {
	nbelow := s.gamma(len(configVals))
	lossAscending := argSort2d(lossVals)

	sort.Ints(lossAscending[:nbelow])
	below := choice(configVals, lossAscending[:nbelow])

	sort.Ints(lossAscending[nbelow:])
	above := choice(configVals, lossAscending[nbelow:])
	return below, above
}

func (s *Sampler) sampleFromGMM(parzenEstimator *ParzenEstimator, low, high float64, size int, q float64, isLog bool) []float64 {
	weights := parzenEstimator.Weights
	mus := parzenEstimator.Mus
	sigmas := parzenEstimator.Sigmas
	nsamples := size

	if low > high {
		panic("the low should be lower than the high")
	}

	samples := make([]float64, 0, nsamples)
	for {
		if len(samples) == nsamples {
			break
		}
		active, err := random.ArgMaxMultinomial(weights)
		if err != nil {
			panic(err)
		}
		x := s.rng.NormFloat64()
		draw := x*sigmas[active] + mus[active]
		if low <= draw && draw < high {
			samples = append(samples, draw)
		}
	}

	if isLog {
		for i := range samples {
			samples[i] = math.Exp(samples[i])
		}
	}

	if q > 0 {
		for i := range samples {
			samples[i] = math.Round(samples[i]/q) * q
		}
	}
	return samples
}

func (s *Sampler) normalCDF(x float64, mu []float64, sigma []float64) []float64 {
	l := len(mu)
	results := make([]float64, l)
	for i := 0; i < l; i++ {
		denominator := x - mu[i]
		numerator := math.Max(math.Sqrt(2)*sigma[i], eps)
		z := denominator / numerator
		results[i] = 0.5 * (1 + math.Erf(z))
	}
	return results
}

func (s *Sampler) logNormalCDF(x float64, mu []float64, sigma []float64) []float64 {
	if x < 0 {
		panic("negative argument is given to logNormalCDF")
	}
	l := len(mu)
	results := make([]float64, l)
	for i := 0; i < l; i++ {
		denominator := math.Log(math.Max(x, eps)) - mu[i]
		numerator := math.Max(math.Sqrt(2)*sigma[i], eps)
		z := denominator / numerator
		results[i] = 0.5 + (0.5 * math.Erf(z))
	}
	return results
}

func (s *Sampler) logsumRows(x [][]float64) []float64 {
	y := make([]float64, len(x))
	for i := range x {
		m := floats.Max(x[i])

		sum := 0.0
		for j := range x[i] {
			sum += math.Log(math.Exp(x[i][j] - m))
		}
		y[i] = sum + m
	}
	return y
}

func (s *Sampler) gmmLogPDF(samples []float64, parzenEstimator *ParzenEstimator, low, high float64, q float64, isLog bool) []float64 {
	weights := parzenEstimator.Weights
	mus := parzenEstimator.Mus
	sigmas := parzenEstimator.Sigmas

	if len(samples) == 0 {
		return []float64{}
	}

	highNormalCdf := s.normalCDF(high, mus, sigmas)
	lowNormalCdf := s.normalCDF(low, mus, sigmas)
	if len(weights) != len(highNormalCdf) {
		panic("the length should be the same with weights")
	}

	paccept := 0.0
	for i := 0; i < len(highNormalCdf); i++ {
		paccept += highNormalCdf[i]*weights[i] - lowNormalCdf[i]
	}

	if q > 0 {
		probabilities := make([]float64, len(samples))
		if len(weights) != len(mus) || len(weights) != len(sigmas) {
			panic("should be the same length of weights, mus and sigmas")
		}
		for i := range weights {
			w := weights[i]
			mu := mus[i]
			sigma := sigmas[i]
			upperBound := make([]float64, len(samples))
			lowerBound := make([]float64, len(samples))
			for i := range upperBound {
				if isLog {
					upperBound[i] = math.Min(samples[i]+q/2.0, math.Exp(high))
					lowerBound[i] = math.Max(samples[i]-q/2.0, math.Exp(low))
					lowerBound[i] = math.Max(0, lowerBound[i])
				} else {
					upperBound[i] = math.Min(samples[i]+q/2.0, high)
					lowerBound[i] = math.Max(samples[i]-q/2.0, low)
				}
			}

			incAmt := make([]float64, len(samples))
			for j := range upperBound {
				if isLog {
					incAmt[j] = w * s.logNormalCDF(upperBound[j], []float64{mu}, []float64{sigma})[0]
					incAmt[j] -= w * s.logNormalCDF(lowerBound[j], []float64{mu}, []float64{sigma})[0]
				} else {
					incAmt[j] = w * s.normalCDF(upperBound[j], []float64{mu}, []float64{sigma})[0]
					incAmt[j] -= w * s.normalCDF(lowerBound[j], []float64{mu}, []float64{sigma})[0]
				}
			}
			for j := range probabilities {
				probabilities[j] += incAmt[j]
			}
		}
		returnValue := make([]float64, len(samples))
		for i := range probabilities {
			returnValue[i] = math.Log(probabilities[i]+eps) + math.Log(paccept+eps)
		}
		return returnValue
	}

	var (
		jacobian []float64
		distance [][]float64
	)
	if isLog {
		jacobian = samples
	} else {
		jacobian = ones1d(len(samples))
	}
	distance = make([][]float64, len(samples))
	for i := range samples {
		distance[i] = make([]float64, len(mus))
		for j := range mus {
			if isLog {
				distance[i][j] = math.Log(samples[i]) - mus[j]
			} else {
				distance[i][j] = samples[i] - mus[j]
			}
		}
	}
	mahalanobis := make([][]float64, len(distance))
	for i := range distance {
		mahalanobis[i] = make([]float64, len(distance[i]))
		for j := range distance[i] {
			mahalanobis[i][j] = distance[i][j] / math.Pow(math.Max(sigmas[j], eps), 2)
		}
	}
	z := make([][]float64, len(distance))
	for i := range distance {
		z[i] = make([]float64, len(distance[i]))
		for j := range distance[i] {
			z[i][j] = math.Sqrt(2*math.Pi) * sigmas[j] * jacobian[i]
		}
	}
	coefficient := make([][]float64, len(distance))
	for i := range distance {
		coefficient[i] = make([]float64, len(distance[i]))
		for j := range distance[i] {
			coefficient[i][j] = weights[j] / z[i][j] / paccept
		}
	}

	y := make([][]float64, len(distance))
	for i := range distance {
		y[i] = make([]float64, len(distance[i]))
		for j := range distance[i] {
			y[i][j] = -0.5*mahalanobis[i][j] + math.Log(coefficient[i][j])
		}
	}
	return s.logsumRows(y)
}

func (s *Sampler) sampleFromCategoricalDist(probabilities []float64, size int) []int {
	if size == 0 {
		return []int{}
	}
	sample := random.Multinomial(1, probabilities, size)

	returnVals := make([]int, size)
	for i := 0; i < size; i++ {
		for j := range sample[i] {
			returnVals[i] += sample[i][j] * j
		}
	}
	return returnVals
}

func (s *Sampler) categoricalLogPDF(sample []int, p []float64) []float64 {
	if len(sample) == 0 {
		return []float64{}
	}

	result := make([]float64, len(sample))
	for i := 0; i < len(sample); i++ {
		result[i] = math.Log(p[sample[i]])
	}
	return result
}

func (s *Sampler) compare(samples []float64, logL []float64, logG []float64) []float64 {
	if len(samples) == 0 {
		return []float64{}
	}
	if len(logL) != len(logG) {
		panic("the size of the log_l and log_g should be same")
	}
	score := make([]float64, len(logL))
	for i := range score {
		score[i] = logL[i] - logG[i]
	}
	if len(samples) != len(score) {
		panic("the size of the samples and score should be same")
	}

	argMax := func(s []float64) int {
		max := s[0]
		maxIdx := 0
		for i := range s {
			if i == 0 {
				continue
			}
			if s[i] > max {
				max = s[i]
				maxIdx = i
			}
		}
		return maxIdx
	}
	best := argMax(score)
	results := make([]float64, len(samples))
	for i := range results {
		results[i] = samples[best]
	}
	return results
}

func (s *Sampler) sampleNumerical(low, high float64, below, above []float64, q float64, isLog bool) float64 {
	if isLog {
		low = math.Log(low)
		high = math.Log(high)
		for i := range below {
			below[i] = math.Log(below[i])
		}
		for i := range above {
			above[i] = math.Log(above[i])
		}
	}
	size := s.numberOfEICandidates
	parzenEstimatorBelow := NewParzenEstimator(below, low, high, s.params)
	sampleBelow := s.sampleFromGMM(parzenEstimatorBelow, low, high, size, q, isLog)
	logLikelihoodsBelow := s.gmmLogPDF(sampleBelow, parzenEstimatorBelow, low, high, q, isLog)

	parzenEstimatorAbove := NewParzenEstimator(above, low, high, s.params)
	logLikelihoodsAbove := s.gmmLogPDF(sampleBelow, parzenEstimatorAbove, low, high, q, isLog)

	return s.compare(sampleBelow, logLikelihoodsBelow, logLikelihoodsAbove)[0]
}

func (s *Sampler) sampleUniform(distribution goptuna.UniformDistribution, below, above []float64) float64 {
	low := distribution.Low
	high := distribution.High
	return s.sampleNumerical(low, high, below, above, 0, false)
}

func (s *Sampler) sampleLogUniform(distribution goptuna.LogUniformDistribution, below, above []float64) float64 {
	low := distribution.Low
	high := distribution.High
	return s.sampleNumerical(low, high, below, above, 0, true)
}

func (s *Sampler) sampleInt(distribution goptuna.IntUniformDistribution, below, above []float64) float64 {
	q := 1.0
	low := float64(distribution.Low) - 0.5*q
	high := float64(distribution.High) + 0.5*q
	return s.sampleNumerical(low, high, below, above, q, false)
}

func (s *Sampler) sampleStepInt(distribution goptuna.StepIntUniformDistribution, below, above []float64) float64 {
	q := 1.0
	low := float64(distribution.Low) - 0.5*q
	high := float64(distribution.High) + 0.5*q
	return s.sampleNumerical(low, high, below, above, q, false)
}

func (s *Sampler) sampleDiscreteUniform(distribution goptuna.DiscreteUniformDistribution, below, above []float64) float64 {
	q := distribution.Q
	r := distribution.High - distribution.Low

	// [low, high] is shifted to [0, r] to align sampled values at regular intervals.
	// See https://github.com/optuna/optuna/pull/917#issuecomment-586114630 for details.
	low := 0 - 0.5*q
	high := r + 0.5*q

	// Shift below and above to [0, r]
	for i := range below {
		below[i] -= distribution.Low
	}
	for i := range above {
		above[i] -= distribution.Low
	}

	best := s.sampleNumerical(low, high, below, above, q, false) + distribution.Low
	return math.Min(math.Max(best, distribution.Low), distribution.High)
}

func (s *Sampler) sampleCategorical(distribution goptuna.CategoricalDistribution, below, above []float64) float64 {
	belowInt := make([]int, len(below))
	for i := range below {
		belowInt[i] = int(below[i])
	}
	aboveInt := make([]int, len(above))
	for i := range above {
		aboveInt[i] = int(above[i])
	}
	upper := len(distribution.Choices)
	size := s.numberOfEICandidates

	// below
	weightsBelow := s.params.Weights(len(below))
	countsBelow := bincount(belowInt, weightsBelow, upper)
	weightedBelowSum := 0.0
	weightedBelow := make([]float64, len(countsBelow))
	for i := range countsBelow {
		weightedBelow[i] = countsBelow[i] + s.params.PriorWeight
		weightedBelowSum += weightedBelow[i]
	}
	for i := range weightedBelow {
		weightedBelow[i] /= weightedBelowSum
	}
	samplesBelow := s.sampleFromCategoricalDist(weightedBelow, size)
	logLikelihoodsBelow := s.categoricalLogPDF(samplesBelow, weightedBelow)

	// above
	weightsAbove := s.params.Weights(len(above))
	countsAbove := bincount(aboveInt, weightsAbove, upper)
	weightedAboveSum := 0.0
	weightedAbove := make([]float64, len(countsAbove))
	for i := range countsAbove {
		weightedAbove[i] = countsAbove[i] + s.params.PriorWeight
		weightedAboveSum += weightedAbove[i]
	}
	for i := range weightedAbove {
		weightedAbove[i] /= weightedAboveSum
	}
	samplesAbove := s.sampleFromCategoricalDist(weightedAbove, size)
	logLikelihoodsAbove := s.categoricalLogPDF(samplesAbove, weightedAbove)

	floatSamplesBelow := make([]float64, len(samplesBelow))
	for i := range samplesBelow {
		floatSamplesBelow[i] = float64(samplesBelow[i])
	}
	return s.compare(floatSamplesBelow, logLikelihoodsBelow, logLikelihoodsAbove)[0]
}

// Sample a parameter for a given distribution.
func (s *Sampler) Sample(
	study *goptuna.Study,
	trial goptuna.FrozenTrial,
	paramName string,
	paramDistribution interface{},
) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	values, scores, err := getObservationPairs(study, paramName)
	if err != nil {
		return 0, err
	}
	n := len(values)

	if n < s.numberOfStartupTrials {
		return s.randomSampler.Sample(study, trial, paramName, paramDistribution)
	}

	belowParamValues, aboveParamValues := s.splitObservationPairs(values, scores)

	switch d := paramDistribution.(type) {
	case goptuna.UniformDistribution:
		return s.sampleUniform(d, belowParamValues, aboveParamValues), nil
	case goptuna.LogUniformDistribution:
		return s.sampleLogUniform(d, belowParamValues, aboveParamValues), nil
	case goptuna.IntUniformDistribution:
		return s.sampleInt(d, belowParamValues, aboveParamValues), nil
	case goptuna.StepIntUniformDistribution:
		return s.sampleStepInt(d, belowParamValues, aboveParamValues), nil
	case goptuna.DiscreteUniformDistribution:
		return s.sampleDiscreteUniform(d, belowParamValues, aboveParamValues), nil
	case goptuna.CategoricalDistribution:
		return s.sampleCategorical(d, belowParamValues, aboveParamValues), nil
	}
	return 0, goptuna.ErrUnknownDistribution
}

func getObservationPairs(study *goptuna.Study, paramName string) ([]float64, [][2]float64, error) {
	var sign float64 = 1
	if study.Direction() == goptuna.StudyDirectionMaximize {
		sign = -1
	}

	trials, err := study.GetTrials()
	if err != nil {
		return nil, nil, err
	}

	values := make([]float64, 0, len(trials))
	scores := make([][2]float64, 0, len(trials))
	for _, trial := range trials {
		ir, ok := trial.InternalParams[paramName]
		if !ok {
			continue
		}

		var paramValue, score0, score1 float64
		paramValue = ir
		if trial.State == goptuna.TrialStateComplete {
			score0 = math.Inf(-1)
			score1 = sign * trial.Value
		} else if trial.State == goptuna.TrialStatePruned {
			if len(trial.IntermediateValues) > 0 {
				var step int
				var intermediateValue float64

				for key := range trial.IntermediateValues {
					if key > step {
						step = key
						intermediateValue = trial.IntermediateValues[key]
					}
				}
				score0 = float64(-step)
				score1 = sign * intermediateValue
			} else {
				score0 = math.Inf(1)
				score1 = 0.0
			}
		} else {
			continue
		}
		values = append(values, paramValue)
		scores = append(scores, [2]float64{score0, score1})
	}
	return values, scores, nil
}
