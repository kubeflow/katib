package tpe

import (
	"math/rand"

	"github.com/c-bata/goptuna"
)

// SamplerOption is a type of the function to customizing TPE sampler.
type SamplerOption func(sampler *Sampler)

// SamplerOptionSeed sets seed number.
func SamplerOptionSeed(seed int64) SamplerOption {
	randomSampler := goptuna.NewRandomSearchSampler(
		goptuna.RandomSearchSamplerOptionSeed(seed))

	return func(sampler *Sampler) {
		sampler.rng = rand.New(rand.NewSource(seed))
		sampler.randomSampler = randomSampler
	}
}

// SamplerOptionGammaFunc sets the gamma function.
func SamplerOptionGammaFunc(gamma FuncGamma) SamplerOption {
	return func(sampler *Sampler) {
		sampler.gamma = gamma
	}
}

// SamplerOptionNumberOfEICandidates sets the number of EI candidates (default 24).
func SamplerOptionNumberOfEICandidates(n int) SamplerOption {
	return func(sampler *Sampler) {
		sampler.numberOfEICandidates = n
	}
}

// SamplerOptionNumberOfStartupTrials sets the number of start up trials (default 10).
func SamplerOptionNumberOfStartupTrials(n int) SamplerOption {
	return func(sampler *Sampler) {
		sampler.numberOfStartupTrials = n
	}
}

// SamplerOptionParzenEstimatorParams sets the parameter of ParzenEstimator.
func SamplerOptionParzenEstimatorParams(params ParzenEstimatorParams) SamplerOption {
	return func(sampler *Sampler) {
		sampler.params = params
	}
}
