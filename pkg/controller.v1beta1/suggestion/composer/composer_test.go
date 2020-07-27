package composer

import (
	"context"
	"encoding/json"
	stdlog "log"
	"os"
	"path/filepath"
	"testing"
	"time"

	apis "github.com/kubeflow/katib/pkg/apis/controller"
	"github.com/kubeflow/katib/pkg/client/controller/clientset/versioned/scheme"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
)

var (
	suggestionName = "test-suggestion"
	namespace      = "kubeflow"
	configMap      = "katib-config"
	timeout        = time.Second * 5
	cfg            *rest.Config
)

func TestMain(m *testing.M) {
	// Start test k8s server
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "manifests", "v1beta1", "katib-controller"),
		},
	}
	apis.AddToScheme(scheme.Scheme)

	var err error
	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}

	code := m.Run()
	t.Stop()
	os.Exit(code)
}
func TestDesiredDeployment(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c := mgr.GetClient()

	composer := New(mgr)

	tcs := []struct {
		suggestion         *suggestionsv1beta1.Suggestion
		configMap          *corev1.ConfigMap
		expectedDeployment *appsv1.Deployment
		err                bool
		testDescription    string
	}{
		{
			suggestion:         nil,
			expectedDeployment: nil,
			err:                true,
			testDescription:    "",
		},
	}

	for _, tc := range tcs {
		actualDeployment, err := composer.DesiredDeployment(tc.suggestion)
		// Create configMap with Katib config
		g.Expect(c.Create(context.TODO(), tc.configMap)).NotTo(gomega.HaveOccurred())
		g.Eventually(func() error {
			return c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: configMap}, &corev1.ConfigMap{})
		}, timeout).ShouldNot(gomega.HaveOccurred())

		if !tc.err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
		} else if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		} else if !equality.Semantic.DeepEqual(tc.expectedDeployment, actualDeployment) {
			t.Errorf("Case: %v failed. Expected deploy %v, got %v", tc.testDescription, tc.expectedDeployment, actualDeployment)
		}
	}
}

func newFakeKatibConfig() *corev1.ConfigMap {
	suggestionConfig := map[string]map[string]string{
		"random": {
			"image":           "test",
			"imagePullPolicy": "Always",
		},
	}
	b, _ := json.Marshal(suggestionConfig)

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "katib-config",
			Namespace: "kubeflow",
		},
		Data: map[string]string{
			"suggestion": string(b),
		},
	}
}

func newFakeSuggestion() *suggestionsv1beta1.Suggestion {
	return &suggestionsv1beta1.Suggestion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      suggestionName,
			Namespace: namespace,
		},
		Spec: suggestionsv1beta1.SuggestionSpec{
			Requests:      1,
			AlgorithmName: "random",
		},
	}
}
