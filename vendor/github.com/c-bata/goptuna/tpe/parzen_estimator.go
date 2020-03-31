package tpe

import (
	"math"

	"gonum.org/v1/gonum/floats"
)

// ParzenEstimatorParams holds the parameters of ParzenEstimator
type ParzenEstimatorParams struct {
	ConsiderPrior     bool
	ConsiderMagicClip bool
	ConsiderEndpoints bool
	Weights           FuncWeights
	PriorWeight       float64 // optional
}

// ParzenEstimator is a surrogate model for TPE>
type ParzenEstimator struct {
	Weights []float64
	Mus     []float64
	Sigmas  []float64
}

func buildEstimator(
	mus []float64,
	low float64,
	high float64,
	params ParzenEstimatorParams,
) ([]float64, []float64, []float64) {
	considerPrior := params.ConsiderPrior
	priorWeight := params.PriorWeight
	considerMagicClip := params.ConsiderMagicClip
	considerEndpoints := params.ConsiderEndpoints
	weightsFunc := params.Weights

	var sortedWeights []float64
	var sortedMus []float64
	var sigma []float64

	var order []int
	var priorPos int
	var priorSigma float64
	if considerPrior {
		priorMu := 0.5 * (low + high)
		priorSigma = 1.0 * (high - low)
		if len(mus) == 0 {
			sortedMus = []float64{priorMu}
			sigma = []float64{priorSigma}
			priorPos = 0
			order = make([]int, 0)
		} else {
			order = make([]int, len(mus))
			floats.Argsort(mus, order)
			priorPos = location(choice(mus, order), priorMu)
			sortedMus = make([]float64, 0, len(mus)+1)
			sortedMus = append(sortedMus, choice(mus, order[:priorPos])...)
			sortedMus = append(sortedMus, priorMu)
			sortedMus = append(sortedMus, choice(mus, order[priorPos:])...)
		}
	} else {
		order = make([]int, len(mus))
		floats.Argsort(mus, order)
		sortedMus = choice(mus, order)
	}

	// we decide the sigma.
	if len(mus) > 0 {
		lowSortedMusHigh := append(sortedMus, high)
		lowSortedMusHigh = append([]float64{low}, lowSortedMusHigh...)

		l := len(lowSortedMusHigh)
		sigma = make([]float64, l)
		for i := 0; i < l-2; i++ {
			sigma[i+1] = math.Max(lowSortedMusHigh[i+1]-lowSortedMusHigh[i], lowSortedMusHigh[i+2]-lowSortedMusHigh[i+1])
		}
		if !considerEndpoints && len(lowSortedMusHigh) > 2 {
			sigma[1] = lowSortedMusHigh[2] - lowSortedMusHigh[1]
			sigma[l-2] = lowSortedMusHigh[l-2] - lowSortedMusHigh[l-3]
		}
		sigma = sigma[1 : l-1]
	}

	// we decide the weights.
	unsortedWeights := weightsFunc(len(mus))
	if considerPrior {
		sortedWeights = make([]float64, 0, len(sortedMus))
		sortedWeights = append(sortedWeights, choice(unsortedWeights, order[:priorPos])...)
		sortedWeights = append(sortedWeights, priorWeight)
		sortedWeights = append(sortedWeights, choice(unsortedWeights, order[priorPos:])...)
	} else {
		sortedWeights = choice(unsortedWeights, order)
	}
	sumSortedWeights := floats.Sum(sortedWeights)
	for i := range sortedWeights {
		sortedWeights[i] /= sumSortedWeights
	}

	// We adjust the range of the 'sigma' according to the 'consider_magic_clip' flag.
	maxSigma := 1.0 * (high - low)
	var minSigma float64
	if considerMagicClip {
		minSigma = 1.0 * (high - low) / math.Min(100.0, 1.0+float64(len(sortedMus)))
	} else {
		minSigma = eps
	}
	clip(sigma, minSigma, maxSigma)
	if considerPrior {
		sigma[priorPos] = priorSigma
	}
	return sortedWeights, sortedMus, sigma
}

// NewParzenEstimator returns the parzen estimator object.
func NewParzenEstimator(mus []float64, low, high float64, params ParzenEstimatorParams) *ParzenEstimator {
	estimator := &ParzenEstimator{
		Weights: nil,
		Mus:     nil,
		Sigmas:  nil,
	}

	sWeights, sMus, sigma := buildEstimator(mus, low, high, params)
	estimator.Weights = sWeights
	estimator.Mus = sMus
	estimator.Sigmas = sigma
	return estimator
}
