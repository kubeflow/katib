package goptuna

import (
	"errors"
	"fmt"
	"time"
)

// StudySummary holds basic attributes and aggregated results of Study.
type StudySummary struct {
	ID            int               `json:"study_id"`
	Name          string            `json:"study_name"`
	Direction     StudyDirection    `json:"direction"`
	BestTrial     FrozenTrial       `json:"best_trial"`
	UserAttrs     map[string]string `json:"user_attrs"`
	SystemAttrs   map[string]string `json:"system_attrs"`
	DatetimeStart time.Time         `json:"datetime_start"`
}

// FrozenTrial holds the status and results of a Trial.
type FrozenTrial struct {
	ID                 int                    `json:"trial_id"`
	StudyID            int                    `json:"study_id"`
	Number             int                    `json:"number"`
	State              TrialState             `json:"state"`
	Value              float64                `json:"value"`
	IntermediateValues map[int]float64        `json:"intermediate_values"`
	DatetimeStart      time.Time              `json:"datetime_start"`
	DatetimeComplete   time.Time              `json:"datetime_complete"`
	InternalParams     map[string]float64     `json:"internal_params"`
	Params             map[string]interface{} `json:"params"`
	Distributions      map[string]interface{} `json:"distributions"`
	UserAttrs          map[string]string      `json:"user_attrs"`
	SystemAttrs        map[string]string      `json:"system_attrs"`
}

// GetLatestStep returns the latest step in intermediate values.
func (t FrozenTrial) GetLatestStep() (step int, exist bool) {
	if len(t.IntermediateValues) == 0 {
		return -1, false
	}
	var maxStep int
	for k := range t.IntermediateValues {
		if k > maxStep {
			maxStep = k
		}
	}
	return maxStep, true
}

// Validate returns error if invalid.
func (t FrozenTrial) validate() error {
	if t.DatetimeStart.IsZero() {
		return errors.New("`DatetimeStart` is supposed to be set")
	}

	if t.State.IsFinished() {
		if t.DatetimeComplete.IsZero() {
			return errors.New("`DatetimeComplete` is supposed to be set for a finished trial")
		}
	} else {
		if !t.DatetimeComplete.IsZero() {
			return errors.New("`DatetimeComplete` is supposed to not be set for a finished trial")
		}
	}

	if len(t.InternalParams) != len(t.Distributions) {
		return errors.New("`Params` and `Distributions` should be the same length")
	}
	for name := range t.InternalParams {
		if _, ok := t.Distributions[name]; !ok {
			return fmt.Errorf("distribution '%s' is not found", name)
		}
	}

	if len(t.Params) != len(t.InternalParams) {
		return errors.New("`Params` and `InternalParams` should be the same length")
	}
	for name := range t.InternalParams {
		ir := t.InternalParams[name]
		d := t.Distributions[name]

		switch typedDistribution := d.(type) {
		case UniformDistribution:
			if !typedDistribution.Contains(ir) {
				return fmt.Errorf("internal param is out of the distribution range")
			}
		case LogUniformDistribution:
			if !typedDistribution.Contains(ir) {
				return fmt.Errorf("internal param is out of the distribution range")
			}
		case IntUniformDistribution:
			if !typedDistribution.Contains(ir) {
				return fmt.Errorf("internal param is out of the distribution range")
			}
		case DiscreteUniformDistribution:
			if !typedDistribution.Contains(ir) {
				return fmt.Errorf("internal param is out of the distribution range")
			}
		case CategoricalDistribution:
			if !typedDistribution.Contains(ir) {
				return fmt.Errorf("internal param is out of the distribution range")
			}
		default:
			return errors.New("unsupported distribution")
		}

		expectedXr, err := ToExternalRepresentation(d, ir)
		if err != nil {
			return err
		}

		actualXr, ok := t.Params[name]
		if !ok {
			return fmt.Errorf("params '%s' is not found", name)
		}

		if expectedXr != actualXr {
			return fmt.Errorf("internal params and external param does not match: %s", name)
		}
	}
	return nil
}
