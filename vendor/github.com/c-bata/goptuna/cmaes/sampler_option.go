package cmaes

import (
	"math/rand"
)

// SamplerOption is a type of the function to customizing CMA-ES sampler.
type SamplerOption func(sampler *Sampler)

// SamplerOptionSeed sets seed number.
func SamplerOptionSeed(seed int64) SamplerOption {
	return func(sampler *Sampler) {
		sampler.rng = rand.New(rand.NewSource(seed))
	}
}

// SamplerOptionInitialMean sets the initial mean vectors.
func SamplerOptionInitialMean(mean map[string]float64) SamplerOption {
	return func(sampler *Sampler) {
		sampler.x0 = mean
	}
}

// SamplerOptionInitialSigma sets the initial sigma.
func SamplerOptionInitialSigma(sigma float64) SamplerOption {
	return func(sampler *Sampler) {
		sampler.sigma0 = sigma
	}
}

// SamplerOptionOptimizerOptions sets the options for Optimizer.
func SamplerOptionOptimizerOptions(opts ...OptimizerOption) SamplerOption {
	return func(sampler *Sampler) {
		sampler.optimizerOptions = opts
	}
}

// SamplerOptionNStartupTrials sets the number of startup trials.
func SamplerOptionNStartupTrials(nStartupTrials int) SamplerOption {
	return func(sampler *Sampler) {
		sampler.nStartUpTrials = nStartupTrials
	}
}
