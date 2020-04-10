package goptuna

import (
	"errors"
	"math"
	"math/rand"
	"sync"
)

var (
	// ErrUnsupportedSearchSpace represents sampler does not support a given search space.
	ErrUnsupportedSearchSpace = errors.New("unsupported search space")
)

// Sampler is the interface for sampling algorithms that do not use
// relationship between parameters such as random sampling and TPE.
//
// Note that if study object has RelativeSampler, this interface is used
// only for parameters that are not sampled by RelativeSampler.
type Sampler interface {
	// Sample a single parameter for a given distribution.
	Sample(*Study, FrozenTrial, string, interface{}) (float64, error)
}

// RelativeSampler is the interface for sampling algorithms that use
// relationship between parameters such as Gaussian Process and CMA-ES.
//
// This interface is called once at the beginning of each trial,
// i.e., right before the evaluation of the objective function.
type RelativeSampler interface {
	// SampleRelative samples multiple dimensional parameters in a given search space.
	SampleRelative(*Study, FrozenTrial, map[string]interface{}) (map[string]float64, error)
}

// IntersectionSearchSpace return return the intersection search space of the Study.
//
// Intersection search space contains the intersection of parameter distributions that have been
// suggested in the completed trials of the study so far.
// If there are multiple parameters that have the same name but different distributions,
// neither is included in the resulting search space
// (i.e., the parameters with dynamic value ranges are excluded).
func IntersectionSearchSpace(study *Study) (map[string]interface{}, error) {
	var searchSpace map[string]interface{}

	trials, err := study.GetTrials()
	if err != nil {
		return nil, err
	}

	for i := range trials {
		if trials[i].State != TrialStateComplete {
			continue
		}

		if searchSpace == nil {
			searchSpace = trials[i].Distributions
			continue
		}

		exists := func(name string) bool {
			for name2 := range trials[i].Distributions {
				if name == name2 {
					return true
				}
			}
			return false
		}

		deleteParams := make([]string, 0, len(searchSpace))
		for name := range searchSpace {
			if !exists(name) {
				deleteParams = append(deleteParams, name)
			} else if trials[i].Distributions[name] != searchSpace[name] {
				deleteParams = append(deleteParams, name)
			}
		}

		for j := range deleteParams {
			delete(searchSpace, deleteParams[j])
		}
	}
	return searchSpace, nil
}

var _ Sampler = &RandomSearchSampler{}

// RandomSearchSampler for random search
type RandomSearchSampler struct {
	rng *rand.Rand
	mu  sync.Mutex
}

// RandomSearchSamplerOption is a type of function to set change the option.
type RandomSearchSamplerOption func(sampler *RandomSearchSampler)

// RandomSearchSamplerOptionSeed sets seed number.
func RandomSearchSamplerOptionSeed(seed int64) RandomSearchSamplerOption {
	return func(sampler *RandomSearchSampler) {
		sampler.rng = rand.New(rand.NewSource(seed))
	}
}

// NewRandomSearchSampler implements random search algorithm.
func NewRandomSearchSampler(opts ...RandomSearchSamplerOption) *RandomSearchSampler {
	s := &RandomSearchSampler{
		rng: rand.New(rand.NewSource(0)),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Sample a parameter for a given distribution.
func (s *RandomSearchSampler) Sample(
	study *Study,
	trial FrozenTrial,
	paramName string,
	paramDistribution interface{},
) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch d := paramDistribution.(type) {
	case UniformDistribution:
		if d.Single() {
			return d.Low, nil
		}
		return s.rng.Float64()*(d.High-d.Low) + d.Low, nil
	case LogUniformDistribution:
		if d.Single() {
			return d.Low, nil
		}
		logLow := math.Log(d.Low)
		logHigh := math.Log(d.High)
		return math.Exp(s.rng.Float64()*(logHigh-logLow) + logLow), nil
	case IntUniformDistribution:
		if d.Single() {
			return float64(d.Low), nil
		}
		return float64(s.rng.Intn(d.High-d.Low) + d.Low), nil
	case StepIntUniformDistribution:
		if d.Single() {
			return float64(d.Low), nil
		}
		r := (d.High - d.Low) / d.Step
		v := (s.rng.Intn(r) * d.Step) + d.Low
		return float64(v), nil
	case DiscreteUniformDistribution:
		if d.Single() {
			return d.Low, nil
		}
		q := d.Q
		r := d.High - d.Low
		// [low, high] is shifted to [0, r] to align sampled values at regular intervals.
		low := 0 - 0.5*q
		high := r + 0.5*q
		x := s.rng.Float64()*(high-low) + low
		v := math.Round(x/q)*q + d.Low
		return math.Min(math.Max(v, d.Low), d.High), nil
	case CategoricalDistribution:
		if d.Single() {
			return float64(0), nil
		}
		return float64(rand.Intn(len(d.Choices))), nil
	default:
		return 0.0, errors.New("undefined distribution")
	}
}
