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

package suggestion

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	//utilrand "k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
)

var log = logf.Log.WithName("experiment-suggestion-client")

type Suggestion interface {
	GetOrCreateSuggestion(instance *experimentsv1beta1.Experiment, suggestionRequests int32) (*suggestionsv1beta1.Suggestion, error)
	UpdateSuggestion(suggestion *suggestionsv1beta1.Suggestion) error
	UpdateSuggestionStatus(suggestion *suggestionsv1beta1.Suggestion) error
}

type General struct {
	scheme *runtime.Scheme
	client.Client
}

func New(scheme *runtime.Scheme, client client.Client) Suggestion {
	return &General{scheme: scheme, Client: client}
}

func (g *General) GetOrCreateSuggestion(instance *experimentsv1beta1.Experiment, suggestionRequests int32) (*suggestionsv1beta1.Suggestion, error) {
	logger := log.WithValues("experiment", types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace})
	suggestion := &suggestionsv1beta1.Suggestion{}
	err := g.Get(context.TODO(),
		types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()}, suggestion)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating Suggestion", "namespace", instance.Namespace, "name", instance.Name, "Suggestion Requests", suggestionRequests)
		if err := g.createSuggestion(instance, suggestionRequests); err != nil {
			logger.Error(err, "CreateSuggestion failed", "instance", instance.Name)
			return nil, err
		}
	} else if err != nil {
		logger.Error(err, "Suggestion get failed", "instance", instance.Name)
		return nil, err
	} else {
		return suggestion, nil
	}
	return nil, nil
}

func (g *General) createSuggestion(instance *experimentsv1beta1.Experiment, suggestionRequests int32) error {
	logger := log.WithValues("experiment", types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace})
	suggestion := &suggestionsv1beta1.Suggestion{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.Name,
			Namespace:   instance.Namespace,
			Labels:      instance.Labels,
			Annotations: instance.Annotations,
		},
		Spec: suggestionsv1beta1.SuggestionSpec{
			Algorithm:    instance.Spec.Algorithm.DeepCopy(),
			Requests:     suggestionRequests,
			ResumePolicy: instance.Spec.ResumePolicy,
		},
	}

	if instance.Spec.EarlyStopping != nil {
		suggestion.Spec.EarlyStopping = instance.Spec.EarlyStopping.DeepCopy()
	}

	if err := controllerutil.SetControllerReference(instance, suggestion, g.scheme); err != nil {
		logger.Error(err, "Error in setting controller reference")
		return err
	}

	if err := g.Create(context.TODO(), suggestion); err != nil {
		return err
	}
	logger.Info("Suggestion created", "namespace", instance.Namespace, "name", instance.Name)
	return nil
}

func (g *General) UpdateSuggestion(suggestion *suggestionsv1beta1.Suggestion) error {
	if err := g.Update(context.TODO(), suggestion); err != nil {
		return err
	}
	return nil
}

func (g *General) UpdateSuggestionStatus(suggestion *suggestionsv1beta1.Suggestion) error {
	if err := g.Status().Update(context.TODO(), suggestion); err != nil {
		return err
	}

	return nil
}
