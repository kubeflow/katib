/*
Copyright 2022 The Kubeflow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package katibclient

import (
	"context"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

type Client interface {
	InjectClient(c client.Client)
	GetClient() client.Client
	GetExperimentList(namespace ...string) (*experimentsv1beta1.ExperimentList, error)
	GetExperiment(name string, namespace ...string) (*experimentsv1beta1.Experiment, error)
	GetConfigMap(name string, namespace ...string) (map[string]string, error)
	GetTrial(name string, namespace ...string) (*trialsv1beta1.Trial, error)
	GetTrialList(name string, namespace ...string) (*trialsv1beta1.TrialList, error)
	GetTrialTemplates(namespace ...string) (*apiv1.ConfigMapList, error)
	GetSuggestion(name string, namespace ...string) (*suggestionsv1beta1.Suggestion, error)
	GetNamespaceList() (*apiv1.NamespaceList, error)
	CreateRuntimeObject(object client.Object) error
	DeleteRuntimeObject(object client.Object) error
	UpdateRuntimeObject(object client.Object) error
}

type KatibClient struct {
	client client.Client
}

func NewWithGivenClient(c client.Client) Client {
	return &KatibClient{
		client: c,
	}
}

func NewClient(options client.Options) (Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	if err = experimentsv1beta1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}
	if err = trialsv1beta1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}
	if err = suggestionsv1beta1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}
	if err = configv1beta1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}
	cl, err := client.New(cfg, options)
	if err != nil {
		return nil, err
	}
	return &KatibClient{
		client: cl,
	}, nil
}

func (k *KatibClient) InjectClient(c client.Client) {
	k.client = c
}

func (k *KatibClient) GetClient() client.Client {
	return k.client
}

func (k *KatibClient) GetExperimentList(namespace ...string) (*experimentsv1beta1.ExperimentList, error) {
	ns := getNamespace(namespace...)
	expList := &experimentsv1beta1.ExperimentList{}
	listOpt := client.InNamespace(ns)

	if err := k.client.List(context.Background(), expList, listOpt); err != nil {
		return expList, err
	}
	return expList, nil

}

// GetSuggestion returns the Suggestion CR for the given name and namespace
func (k *KatibClient) GetSuggestion(name string, namespace ...string) (
	*suggestionsv1beta1.Suggestion, error) {
	ns := getNamespace(namespace...)
	suggestion := &suggestionsv1beta1.Suggestion{}

	if err := k.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: ns}, suggestion); err != nil {
		return nil, err
	}
	return suggestion, nil

}

// GetTrial returns the Trial for the given name and namespace
func (k *KatibClient) GetTrial(name string, namespace ...string) (*trialsv1beta1.Trial, error) {
	ns := getNamespace(namespace...)
	trial := &trialsv1beta1.Trial{}

	if err := k.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: ns}, trial); err != nil {
		return nil, err
	}
	return trial, nil

}

func (k *KatibClient) GetTrialList(name string, namespace ...string) (*trialsv1beta1.TrialList, error) {
	ns := getNamespace(namespace...)
	trialList := &trialsv1beta1.TrialList{}
	labels := map[string]string{consts.LabelExperimentName: name}

	if err := k.client.List(context.Background(), trialList, client.InNamespace(ns), client.MatchingLabels(labels)); err != nil {
		return trialList, err
	}
	return trialList, nil

}

func (k *KatibClient) GetExperiment(name string, namespace ...string) (*experimentsv1beta1.Experiment, error) {
	ns := getNamespace(namespace...)
	exp := &experimentsv1beta1.Experiment{}
	if err := k.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: ns}, exp); err != nil {
		return nil, err
	}
	return exp, nil
}

// GetConfigMap returns the configmap for the given name and namespace.
func (k *KatibClient) GetConfigMap(name string, namespace ...string) (map[string]string, error) {
	ns := getNamespace(namespace...)
	configMap := &apiv1.ConfigMap{}
	if err := k.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: ns}, configMap); err != nil {
		return map[string]string{}, err
	}
	return configMap.Data, nil
}

// GetTrialTemplates returns all trial templates from the given namespace
func (k *KatibClient) GetTrialTemplates(namespace ...string) (*apiv1.ConfigMapList, error) {
	ns := getNamespace(namespace...)

	templatesConfigMapList := &apiv1.ConfigMapList{}

	templateLabel := map[string]string{consts.LabelTrialTemplateConfigMapName: consts.LabelTrialTemplateConfigMapValue}

	err := k.client.List(context.TODO(), templatesConfigMapList, client.InNamespace(ns), client.MatchingLabels(templateLabel))

	if err != nil {
		return nil, err
	}

	return templatesConfigMapList, nil

}

func getNamespace(namespace ...string) string {
	if len(namespace) == 0 {
		return consts.DefaultKatibNamespace
	}
	return namespace[0]
}

func (k *KatibClient) GetNamespaceList() (*apiv1.NamespaceList, error) {

	namespaceList := &apiv1.NamespaceList{}

	if err := k.client.List(context.TODO(), namespaceList); err != nil {
		return namespaceList, err
	}
	return namespaceList, nil
}

// CreateRuntimeObject creates the given runtime object in Kubernetes cluster.
func (k *KatibClient) CreateRuntimeObject(object client.Object) error {

	if err := k.client.Create(context.Background(), object); err != nil {
		return err
	}
	return nil
}

// DeleteRuntimeObject deletes the given runtime object in Kubernetes cluster.
func (k *KatibClient) DeleteRuntimeObject(object client.Object) error {

	if err := k.client.Delete(context.Background(), object); err != nil {
		return err
	}
	return nil
}

// UpdateRuntimeObject updates the given runtime object in Kubernetes cluster.
func (k *KatibClient) UpdateRuntimeObject(object client.Object) error {

	if err := k.client.Update(context.Background(), object); err != nil {
		return err
	}
	return nil
}
