package sidecarmetricscollector

import (
	"fmt"
	"math"
	"strconv"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
)

type RuleSet struct {
	spec   []commonv1beta1.EarlyStoppingRule
	status []struct {
		pruner earlyStoppingPruner
		reach  bool
	}
}

func NewRuleSet(
	objMetric string,
	objType commonv1beta1.ObjectiveType,
	spec []commonv1beta1.EarlyStoppingRule,
) (*RuleSet, error) {
	s := &RuleSet{
		spec: spec,
		status: make([]struct {
			pruner earlyStoppingPruner
			reach  bool
		}, len(spec)),
	}

	for i, rule := range spec {
		pruner, err := defaultFactory(rule)
		if err != nil {
			return nil, err
		}
		if objMetric == rule.Name {
			pruner = &objPruner{
				objType:         objType,
				optimalObjValue: math.NaN(),
				sub:             pruner,
			}
		}
		s.status[i].pruner = pruner
	}

	return s, nil
}

func (s *RuleSet) LiveMetrics() []string {
	ls := make([]string, 0, len(s.spec))
	for i, rule := range s.spec {
		if !s.status[i].reach {
			ls = append(ls, rule.Name)
		}
	}
	return ls
}

func (s *RuleSet) UpdateMetric(name string, metricValue float64) error {
	for i := range s.spec {
		rule := &s.spec[i]
		status := &s.status[i]
		if rule.Name != name || status.reach {
			continue
		}

		reach, err := status.pruner.Pruner(metricValue)
		if err != nil {
			return err
		}
		if reach {
			status.reach = true
		}
	}
	return nil
}

type earlyStoppingPruner interface {
	Pruner(metricValue float64) (bool, error)
}

func defaultFactory(rule commonv1beta1.EarlyStoppingRule) (earlyStoppingPruner, error) {
	r := rule
	switch rule.Comparison {
	case commonv1beta1.ComparisonTypeGreater, commonv1beta1.ComparisonTypeLess, commonv1beta1.ComparisonTypeEqual:
		value, err := strconv.ParseFloat(r.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse value to float for rule metric %s: %w", r.Name, err)
		}
		return &basicPruner{
			target:    value,
			startStep: r.StartStep,
			cmp:       r.Comparison,
		}, nil
	default:
		return nil, fmt.Errorf("unknown rule comparison: %s", r.Comparison)
	}
}

type basicPruner struct {
	target    float64
	step      int
	startStep int
	cmp       commonv1beta1.ComparisonType
}

func (p *basicPruner) Pruner(metricValue float64) (bool, error) {
	p.step++
	if p.startStep > 0 && p.step < p.startStep {
		return false, nil
	}
	switch p.cmp {
	case commonv1beta1.ComparisonTypeLess:
		return metricValue < p.target, nil
	case commonv1beta1.ComparisonTypeGreater:
		return metricValue > p.target, nil
	case commonv1beta1.ComparisonTypeEqual:
		return metricValue == p.target, nil
	default:
		return false, fmt.Errorf("unknown rule comparison: %s", p.cmp)
	}
}

type objPruner struct {
	objType         commonv1beta1.ObjectiveType
	optimalObjValue float64
	sub             earlyStoppingPruner
}

func (p *objPruner) Pruner(metricValue float64) (bool, error) {
	// For objective metric we calculate best optimal value from the recorded metrics.
	// This is workaround for Median Stop algorithm.
	// TODO (andreyvelich): Think about it, maybe define latest, max or min strategy type in stop-rule as well ?

	if math.IsNaN(p.optimalObjValue) {
		p.optimalObjValue = metricValue
	} else if p.objType == commonv1beta1.ObjectiveTypeMaximize && metricValue > p.optimalObjValue {
		p.optimalObjValue = metricValue
	} else if p.objType == commonv1beta1.ObjectiveTypeMinimize && metricValue < p.optimalObjValue {
		p.optimalObjValue = metricValue
	}
	// Assign best optimal value to metric value.
	metricValue = p.optimalObjValue

	return p.sub.Pruner(metricValue)
}
