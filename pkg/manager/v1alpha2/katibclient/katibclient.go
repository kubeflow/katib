package katibclient

import (
	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type KatibClient struct {
	client client.Client
}

func newKatibClient() (*KatibClient, error) {
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		return nil, err
	}
	return &KatibClient{
			client: cl,
		},
		nil

}

func (k *KatibClient) GetExperimentList() (*experimentv1alpha2.ExperimentList, error) {

	// Using a unstructured object.
	u := &unstructured.UnstructuredList{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apps",
		Kind:    "Experiment",
		Version: "v1",
	})
	_ = k.client.List()

	return u.Items, nil
}
