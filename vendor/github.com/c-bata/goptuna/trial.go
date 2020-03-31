package goptuna

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
)

//go:generate stringer -trimprefix TrialState -output stringer_trial_state.go -type=TrialState

// TrialState is a state of Trial
type TrialState int

const (
	// TrialStateRunning means Trial is running.
	TrialStateRunning TrialState = iota
	// TrialStateComplete means Trial has been finished without any error.
	TrialStateComplete
	// TrialStatePruned means Trial has been pruned.
	TrialStatePruned
	// TrialStateFail means Trial has failed due to an uncaught error.
	TrialStateFail
	// TrialStateWaiting means Trial has been stopped, but may be resuming.
	TrialStateWaiting
)

// IsFinished returns true if trial is not running.
func (i TrialState) IsFinished() bool {
	return i != TrialStateRunning && i != TrialStateWaiting
}

// Trial is a process of evaluating an objective function.
//
// This object is passed to an objective function and provides interfaces to get parameter
// suggestion, manage the trial's state of the trial.
// Note that this object is seamlessly instantiated and passed to the objective function behind;
// hence, in typical use cases, library users do not care about instantiation of this object.
type Trial struct {
	Study               *Study
	ID                  int
	state               TrialState
	value               float64
	relativeParams      map[string]float64
	relativeSearchSpace map[string]interface{}
}

func (t *Trial) isFixedParam(name string, distribution interface{}) (float64, bool, error) {
	systemAttrs, err := t.GetSystemAttrs()
	if err != nil {
		return 0, false, err
	}
	fixedParamsJSON, ok := systemAttrs["fixed_params"]
	if !ok {
		return 0, false, nil
	}

	var fixedParams map[string]float64
	err = json.Unmarshal([]byte(fixedParamsJSON), &fixedParams)
	if err != nil {
		return 0, false, err
	}

	internalParam, ok := fixedParams[name]
	if !ok {
		return 0, false, nil
	}

	switch typedDistribution := distribution.(type) {
	case UniformDistribution:
		if !typedDistribution.Contains(internalParam) {
			return 0, false, nil
		}
	case LogUniformDistribution:
		if !typedDistribution.Contains(internalParam) {
			return 0, false, nil
		}
	case DiscreteUniformDistribution:
		if !typedDistribution.Contains(internalParam) {
			return 0, false, nil
		}
	case IntUniformDistribution:
		if !typedDistribution.Contains(internalParam) {
			return 0, false, nil
		}
	case CategoricalDistribution:
		if !typedDistribution.Contains(internalParam) {
			return 0, false, nil
		}
	default:
		return 0, false, errors.New("unsupported distribution")
	}
	return internalParam, true, nil
}

func (t *Trial) isRelativeParam(name string, distribution interface{}) bool {
	expected, ok := t.relativeSearchSpace[name]
	if !ok {
		return false
	}
	return reflect.DeepEqual(expected, distribution)
}

func (t *Trial) suggest(name string, distribution interface{}) (float64, error) {
	trial, err := t.Study.Storage.GetTrial(t.ID)
	if err != nil {
		return 0.0, err
	}

	if value, ok, err := t.isFixedParam(name, distribution); err != nil {
		return 0.0, err
	} else if ok {
		err = t.Study.Storage.SetTrialParam(t.ID, name, value, distribution)
		return value, err
	}

	if t.isRelativeParam(name, distribution) {
		// isRelativeParam ensure that 'distribution' is same
		// with the one's in relativeSearchSpace.
		value, ok := t.relativeParams[name]
		if ok {
			err = t.Study.Storage.SetTrialParam(trial.ID, name, value, distribution)
			return value, err
		}
	}

	v, err := t.Study.Sampler.Sample(t.Study, trial, name, distribution)
	if err != nil {
		return 0.0, err
	}

	err = t.Study.Storage.SetTrialParam(trial.ID, name, v, distribution)
	return v, err
}

// Report an intermediate value of an objective function
func (t *Trial) Report(value float64, step int) error {
	if step < 0 {
		return errors.New("step should be larger equal than 0")
	}
	return t.Study.Storage.SetTrialIntermediateValue(t.ID, step, value)
}

// ShouldPrune judges whether the trial should be pruned.
// This method calls prune method of the pruner, which judges whether
// the trial should be pruned at the given step.
func (t *Trial) ShouldPrune(value float64) (bool, error) {
	if t.Study.Pruner == nil {
		t.Study.logger.Warn("Although it's not registered pruner, but you calls ShouldPrune method")
		return false, nil
	}

	trial, err := t.Study.Storage.GetTrial(t.ID)
	if err != nil {
		return false, err
	}
	return t.Study.Pruner.Prune(t.Study, trial)
}

// Number return trial's number which is consecutive and unique in a study.
func (t *Trial) Number() (int, error) {
	return t.Study.Storage.GetTrialNumberFromID(t.ID)
}

// SuggestUniform suggests a value from a uniform distribution.
func (t *Trial) SuggestUniform(name string, low, high float64) (float64, error) {
	if low > high {
		return 0, errors.New("'low' must be smaller than or equal to the 'high'")
	}
	return t.suggest(name, UniformDistribution{
		High: high, Low: low,
	})
}

// SuggestLogUniform suggests a value from a uniform distribution in the log domain.
func (t *Trial) SuggestLogUniform(name string, low, high float64) (float64, error) {
	if low > high {
		return 0, errors.New("'low' must be smaller than or equal to the 'high'")
	}
	v, err := t.suggest(name, LogUniformDistribution{
		High: high, Low: low,
	})
	return v, err
}

// SuggestInt suggests an integer parameter.
func (t *Trial) SuggestInt(name string, low, high int) (int, error) {
	if low > high {
		return 0, errors.New("'low' must be smaller than or equal to the 'high'")
	}
	v, err := t.suggest(name, IntUniformDistribution{
		High: high, Low: low,
	})
	return int(v), err
}

// SuggestDiscreteUniform suggests a value from a discrete uniform distribution.
func (t *Trial) SuggestDiscreteUniform(name string, low, high, q float64) (float64, error) {
	if low > high {
		return 0, errors.New("'low' must be smaller than or equal to the 'high'")
	}
	d := DiscreteUniformDistribution{
		High: high, Low: low, Q: q,
	}
	ir, err := t.suggest(name, d)
	if err != nil {
		return 0, err
	}
	return d.ToExternalRepr(ir).(float64), err
}

// SuggestCategorical suggests an categorical parameter.
func (t *Trial) SuggestCategorical(name string, choices []string) (string, error) {
	if len(choices) == 0 {
		return "", errors.New("'choices' must contains one or more elements")
	}
	v, err := t.suggest(name, CategoricalDistribution{
		Choices: choices,
	})
	return choices[int(v)], err
}

// SetUserAttr to store the value for the user.
func (t *Trial) SetUserAttr(key, value string) error {
	return t.Study.Storage.SetTrialUserAttr(t.ID, key, value)
}

// SetSystemAttr to store the value for the system.
func (t *Trial) SetSystemAttr(key, value string) error {
	return t.Study.Storage.SetTrialSystemAttr(t.ID, key, value)
}

// GetUserAttrs to store the value for the user.
func (t *Trial) GetUserAttrs() (map[string]string, error) {
	return t.Study.Storage.GetTrialUserAttrs(t.ID)
}

// GetSystemAttrs to store the value for the system.
func (t *Trial) GetSystemAttrs() (map[string]string, error) {
	return t.Study.Storage.GetTrialSystemAttrs(t.ID)
}

// GetContext returns a context which is registered at 'study.WithContext()'.
func (t *Trial) GetContext() context.Context {
	return t.Study.ctx
}
