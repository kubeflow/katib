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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		SuggestionClient: suggestionclient.New(),
		scheme:           mgr.GetScheme(),
		Composer:         composer.New(mgr.GetScheme(), mgr.GetClient()),
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
	log.Info("Suggestion controller created")
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
	oldS := &suggestionsv1alpha3.Suggestion{}
	err := r.Get(context.TODO(), request.NamespacedName, oldS)
	if err != nil {
		if errors.IsNotFound(err) {
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	instance := oldS.DeepCopy()
	if instance.IsCompleted() {
		return reconcile.Result{}, nil
	}
	if !instance.IsCreated() {
		if instance.Status.StartTime == nil {
			now := metav1.Now()
			instance.Status.StartTime = &now
		}
		msg := "Suggestion is created"
		instance.MarkSuggestionStatusCreated(SuggestionCreatedReason, msg)
	} else {
		err := r.ReconcileSuggestion(instance)
		if err != nil {
			logger.Error(err, "Reconcile Suggestion error")
			return reconcile.Result{}, err
		}
	}

	if err := r.updateStatus(instance, oldS); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileSuggestion) ReconcileSuggestion(instance *suggestionsv1alpha3.Suggestion) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	service, err := r.DesiredService(instance)
	if err != nil {
		return err
	}
	_, err = r.reconcileService(service)
	if err != nil {
		return err
	}

	deploy, err := r.DesiredDeployment(instance)
	if err != nil {
		return err
	}
	if foundDeploy, err := r.reconcileDeployment(deploy); err != nil {
		return err
	} else {
		if isReady := r.checkDeploymentReady(foundDeploy); isReady != true {
			// deployment is not ready yet
			return nil
		}
	}
	experiment := &experimentsv1alpha3.Experiment{}
	trials := &trialsv1alpha3.TrialList{}

	if err := r.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace,
	}, experiment); err != nil {
		return err
	}

	if err := r.List(context.TODO(),
		client.MatchingLabels(util.TrialLabels(experiment)), trials); err != nil {
		return err
	}
	if !instance.IsRunning() {
		if err = r.ValidateAlgorithmSettings(instance, experiment); err != nil {
			logger.Error(err, "Marking suggestion failed as algorithm settings validation failed")
			msg := fmt.Sprintf("Validation failed: %v", err)
			instance.MarkSuggestionStatusFailed(SuggestionFailedReason, msg)
			// return nil since it is a terminal condition
			return nil
		}
		msg := "Suggestion is running"
		instance.MarkSuggestionStatusRunning(SuggestionRunningReason, msg)
	}
	logger.Info("Sync assignments", "suggestions", instance.Spec.Requests)
	if err = r.SyncAssignments(instance, experiment, trials.Items); err != nil {
		return err
	}

	return nil
}

func (r *ReconcileSuggestion) checkDeploymentReady(deploy *appsv1.Deployment) bool {
	if deploy == nil {
		return false
	} else {
		for _, cond := range deploy.Status.Conditions {
			if cond.Type == appsv1.DeploymentAvailable && cond.Status == corev1.ConditionTrue {
				return true
			}
		}
	}
	return false
}
