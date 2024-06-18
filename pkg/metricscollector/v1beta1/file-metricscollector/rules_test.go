package sidecarmetricscollector

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
)

func TestRuleSet(t *testing.T) {
	testCases := []struct {
		name      string
		objMetric string
		objType   commonv1beta1.ObjectiveType
		spec      []commonv1beta1.EarlyStoppingRule
		action    func(t *testing.T, s *RuleSet)
	}{
		{
			name:      "simple",
			objMetric: "obj",
			objType:   commonv1beta1.ObjectiveTypeMinimize,
			spec: []commonv1beta1.EarlyStoppingRule{
				{
					Name:       "a",
					Value:      "0.2",
					Comparison: commonv1beta1.ComparisonTypeGreater,
					StartStep:  2,
				},
				{
					Name:       "b",
					Value:      "0.5",
					Comparison: commonv1beta1.ComparisonTypeLess,
					StartStep:  3,
				},
				{
					Name:       "c",
					Value:      "1",
					Comparison: commonv1beta1.ComparisonTypeEqual,
					StartStep:  0,
				},
			},
			action: func(t *testing.T, s *RuleSet) {
				diff(t, []string{"a", "b", "c"}, s.LiveMetrics())
				err := s.UpdateMetric("c", 1)
				if err != nil {
					t.Error(err)
				}
				diff(t, []string{"a", "b"}, s.LiveMetrics())
				err = s.UpdateMetric("a", 1)
				if err != nil {
					t.Error(err)
				}
				err = s.UpdateMetric("b", 0)
				if err != nil {
					t.Error(err)
				}
				err = s.UpdateMetric("b", 0)
				if err != nil {
					t.Error(err)
				}
				diff(t, []string{"a", "b"}, s.LiveMetrics())
				err = s.UpdateMetric("a", 0.1)
				if err != nil {
					t.Error(err)
				}
				diff(t, []string{"a", "b"}, s.LiveMetrics())
				err = s.UpdateMetric("a", 0.21)
				if err != nil {
					t.Error(err)
				}
				diff(t, []string{"b"}, s.LiveMetrics())
				err = s.UpdateMetric("b", 0.2)
				if err != nil {
					t.Error(err)
				}
				diff(t, []string{}, s.LiveMetrics())
			},
		},
		{
			name:      "obj",
			objMetric: "obj",
			objType:   commonv1beta1.ObjectiveTypeMaximize,
			spec: []commonv1beta1.EarlyStoppingRule{
				{
					Name:       "obj",
					Value:      "0.8",
					Comparison: commonv1beta1.ComparisonTypeGreater,
					StartStep:  2,
				},
				{
					Name:       "a",
					Value:      "0.5",
					Comparison: commonv1beta1.ComparisonTypeLess,
					StartStep:  2,
				},
			},
			action: func(t *testing.T, s *RuleSet) {
				diff(t, []string{"obj", "a"}, s.LiveMetrics())
				err := s.UpdateMetric("obj", 1)
				if err != nil {
					t.Error(err)
				}
				err = s.UpdateMetric("a", 0.6)
				if err != nil {
					t.Error(err)
				}
				diff(t, []string{"obj", "a"}, s.LiveMetrics())
				err = s.UpdateMetric("obj", 0.7)
				if err != nil {
					t.Error(err)
				}
				err = s.UpdateMetric("a", 0.6)
				if err != nil {
					t.Error(err)
				}
				diff(t, []string{"a"}, s.LiveMetrics())
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewRuleSet(tt.objMetric, tt.objType, tt.spec)
			if err != nil {
				t.Fatalf("failed to NewRuleSet: %v", err)
			}

			tt.action(t, s)
		})
	}
}

func diff(t *testing.T, want, got any) {
	t.Helper()
	if diff := cmp.Diff(want, got); len(diff) != 0 {
		t.Errorf("Unexpected error (-want,+got):\n%s", diff)
	}
}
