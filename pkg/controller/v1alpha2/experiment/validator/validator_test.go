package validator

import (
	"testing"

	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	mockmanifest "github.com/kubeflow/katib/pkg/mock/v1alpha2/experiment/manifest"
)

func TestValidateTFJobTrialTemplate(t *testing.T) {
	trialTFJobTemplate := `apiVersion: "kubeflow.org/v1beta1"
kind: "TFJob"
metadata:
    name: "dist-mnist-for-e2e-test"
spec:
    tfReplicaSpecs:
        Worker:
            template:
                spec:
                    containers:
                      - name: tensorflow
                        image: gaocegege/mnist:1`

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	p := mockmanifest.NewMockProducer(mockCtrl)
	g := New(p)

	p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(trialTFJobTemplate, nil)

	instance := newFakeInstance()
	if err := g.(*General).validateTrialTemplate(instance); err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestValidateJobTrialTemplate(t *testing.T) {
	trialTFJobTemplate := `apiVersion: "batc1/v1"
kind: "Job"
metadata:
    name: "dist-mnist-for-e2e-test"
spec:
    tfReplicaSpecs:
        Worker:
            template:
                spec:
                    containers:
                      - name: tensorflow
                        image: gaocegege/mnist:1`

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	p := mockmanifest.NewMockProducer(mockCtrl)
	g := New(p)

	p.EXPECT().GetRunSpec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(trialTFJobTemplate, nil)

	instance := newFakeInstance()
	if err := g.(*General).validateTrialTemplate(instance); err != nil {
		t.Errorf("Expected nil, got err %v", err)
	}
}

func newFakeInstance() *experimentsv1alpha2.Experiment {
	return &experimentsv1alpha2.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "fakens",
		},
	}
}
