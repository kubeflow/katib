/*
Copyright 2019 The Kubernetes Authors.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/suggestion/composer"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/suggestion/suggestionclient"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/suggestion/suggestionclient/fake"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/util"
)

const (
	ControllerName = "suggestion-controller"
)

var log = logf.Log.WithName(ControllerName)

// Add creates a new Suggestion Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSuggestion{
		Client:           mgr.GetClient(),
		SuggestionClient: fake.New(),
		scheme:           mgr.GetScheme(),
		Composer:         composer.New(mgr.GetScheme()),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("suggestion-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &suggestionsv1alpha3.Suggestion{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &suggestionsv1alpha3.Suggestion{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &suggestionsv1alpha3.Suggestion{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileSuggestion{}

// ReconcileSuggestion reconciles a Suggestion object
type ReconcileSuggestion struct {
	client.Client
	composer.Composer
	suggestionclient.SuggestionClient

	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Suggestion object and makes changes based on the state read
// and what is in the Suggestion.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=katib.kubeflow.org,resources=suggestions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=katib.kubeflow.org,resources=suggestions/status,verbs=get;update;patch
func (r *ReconcileSuggestion) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("Suggestion", request.NamespacedName)
	// Fetch the Suggestion instance
	instance := &suggestionsv1alpha3.Suggestion{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	oldS := instance.DeepCopy()

	deploy, err := r.DesiredDeployment(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	_, err = r.createOrUpdateDeployment(deploy)
	if err != nil {
		return reconcile.Result{}, err
	}

	service, err := r.DesiredService(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	_, err = r.createOrUpdateService(service)
	if err != nil {
		return reconcile.Result{}, err
	}

	experiment := &experimentsv1alpha3.Experiment{}
	trials := &trialsv1alpha3.TrialList{}

	if err := r.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace,
	}, experiment); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.List(context.TODO(),
		client.MatchingLabels(util.TrialLabels(experiment)), trials); err != nil {
		return reconcile.Result{}, err
	}

	logger.Info("Sync assignments", "suggestions", instance.Spec.Suggestions)
	if err = r.SyncAssignments(instance, experiment, trials.Items); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(gaocegege): Check deployment status and update suggestion's status.
	if err := r.updateStatus(instance, oldS); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
