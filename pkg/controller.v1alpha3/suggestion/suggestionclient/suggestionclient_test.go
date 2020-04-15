package suggestionclient

import (
	commonv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	"github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"testing"
)

func init() {
	logf.SetLogger(logf.ZapLogger(true))
}

func TestConvertTrialObservation(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var strategies = map[string]commonv1alpha3.MetricStrategy{
		"error":    commonv1alpha3.ExtractByMin,
		"auc":      commonv1alpha3.ExtractByMax,
		"accuracy": commonv1alpha3.ExtractByLatest,
	}
	var observation = &commonv1alpha3.Observation{
		Metrics: []commonv1alpha3.Metric{
			{Name: "error", Min: 0.01, Max: 0.08, Latest: 0.05},
			{Name: "auc", Min: 0.70, Max: 0.95, Latest: 0.90},
			{Name: "accuracy", Min: 0.8, Max: 0.94, Latest: 0.93},
		},
	}
	obsPb := convertTrialObservation(strategies, observation)
	g.Expect(obsPb.Metrics[0].Name).To(gomega.Equal("error"))
	g.Expect(obsPb.Metrics[0].Value).To(gomega.Equal(0.01))
	g.Expect(obsPb.Metrics[1].Name).To(gomega.Equal("auc"))
	g.Expect(obsPb.Metrics[1].Value).To(gomega.Equal(0.95))
	g.Expect(obsPb.Metrics[2].Name).To(gomega.Equal("accuracy"))
	g.Expect(obsPb.Metrics[2].Value).To(gomega.Equal(0.93))
}
