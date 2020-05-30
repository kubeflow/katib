package suggestionclient

import (
	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	"github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"strconv"
	"testing"
)

func init() {
	logf.SetLogger(logf.ZapLogger(true))
}

func TestConvertTrialObservation(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var strategies = map[string]commonv1beta1.MetricStrategy{
		"error":    commonv1beta1.ExtractByMin,
		"auc":      commonv1beta1.ExtractByMax,
		"accuracy": commonv1beta1.ExtractByLatest,
	}
	var observation = &commonv1beta1.Observation{
		Metrics: []commonv1beta1.Metric{
			{Name: "error", Min: 0.01, Max: 0.08, Latest: "0.05"},
			{Name: "auc", Min: 0.70, Max: 0.95, Latest: "0.90"},
			{Name: "accuracy", Min: 0.8, Max: 0.94, Latest: "0.93"},
		},
	}
	obsPb := convertTrialObservation(strategies, observation)
	g.Expect(obsPb.Metrics[0].Name).To(gomega.Equal("error"))
	value, _ := strconv.ParseFloat(obsPb.Metrics[0].Value, 64)
	g.Expect(value).To(gomega.Equal(0.01))
	g.Expect(obsPb.Metrics[1].Name).To(gomega.Equal("auc"))
	value, _ = strconv.ParseFloat(obsPb.Metrics[1].Value, 64)
	g.Expect(value).To(gomega.Equal(0.95))
	g.Expect(obsPb.Metrics[2].Name).To(gomega.Equal("accuracy"))
	value, _ = strconv.ParseFloat(obsPb.Metrics[2].Value, 64)
	g.Expect(value).To(gomega.Equal(0.93))
}
