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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/suggestion/composer"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/suggestion/suggestionclient"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
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
		Composer:         composer.New(mgr),
		recorder:         mgr.GetEventRecorderFor(ControllerName),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("suggestion-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	if err = c.Watch(source.Kind(mgr.GetCache(), &suggestionsv1beta1.Suggestion{}), &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	eventHandler := handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &suggestionsv1beta1.Suggestion{}, handler.OnlyControllerOwner())
	if err = c.Watch(source.Kind(mgr.GetCache(), &appsv1.Deployment{}), eventHandler); err != nil {
		return err
	}
	if err = c.Watch(source.Kind(mgr.GetCache(), &corev1.Service{}), eventHandler); err != nil {
		return err
	}
	if err = c.Watch(source.Kind(mgr.GetCache(), &corev1.PersistentVolumeClaim{}), eventHandler); err != nil {
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

	scheme   *runtime.Scheme
	recorder record.EventRecorder
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
func (r *ReconcileSuggestion) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("Suggestion", request.NamespacedName)
	// Fetch the Suggestion instance
	oldS := &suggestionsv1beta1.Suggestion{}
	err := r.Get(ctx, request.NamespacedName, oldS)
	if err != nil {
		if errors.IsNotFound(err) {
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	instance := oldS.DeepCopy()
	// Suggestion will be succeeded if ResumePolicy = Never or ResumePolicy = FromVolume
	if instance.IsSucceeded() {
		err = r.deleteDeployment(instance, request.NamespacedName)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.deleteService(instance, request.NamespacedName)
		if err != nil {
			return reconcile.Result{}, err
		}
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
			r.recorder.Eventf(instance, corev1.EventTypeWarning,
				consts.ReconcileErrorReason, err.Error())

			// Try updating just the status condition when possible
			// Status conditions might need to be updated even in error
			// Ignore all other status fields else it will be inconsistent during retry
			_ = r.updateStatusCondition(instance, oldS)
			logger.Error(err, "Reconcile Suggestion error")
			return reconcile.Result{}, err
		}
	}

	if err := r.updateStatus(instance, oldS); err != nil {
		logger.Info("Update suggestion instance status failed, reconciler requeued", "err", err)
		return reconcile.Result{
			Requeue: true,
		}, nil
	}
	return reconcile.Result{}, nil
}

// ReconcileSuggestion is the main reconcile loop for suggestion CR.
func (r *ReconcileSuggestion) ReconcileSuggestion(instance *suggestionsv1beta1.Suggestion) error {
	suggestionNsName := types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()}
	logger := log.WithValues("Suggestion", suggestionNsName)

	// If ResumePolicy = FromVolume volume is reconciled for suggestion
	if instance.Spec.ResumePolicy == experimentsv1beta1.FromVolume {
		pvc, pv, err := r.DesiredVolume(instance)
		if err != nil {
			return err
		}

		// Reconcile PVC and PV
		_, _, err = r.reconcileVolume(pvc, pv, suggestionNsName)
		if err != nil {
			return err
		}

	}

	service, err := r.DesiredService(instance)
	if err != nil {
		return err
	}
	_, err = r.reconcileService(service, suggestionNsName)
	if err != nil {
		return err
	}

	deploy, err := r.DesiredDeployment(instance)
	if err != nil {
		return err
	}

	// If early stopping is used, create RBAC.
	// If controller should reconcile RBAC,
	// ServiceAccount name must be equal to <suggestion-name>-<suggestion-algorithm>
	if instance.Spec.EarlyStopping != nil && deploy.Spec.Template.Spec.ServiceAccountName == util.GetSuggestionRBACName(instance) {

		serviceAccount, role, roleBinding, err := r.DesiredRBAC(instance)
		if err != nil {
			return err
		}

		// Reconcile ServiceAccount, Role and RoleBinding
		err = r.reconcileRBAC(serviceAccount, role, roleBinding, suggestionNsName)
		if err != nil {
			return err
		}
	}

	if foundDeploy, err := r.reconcileDeployment(deploy, suggestionNsName); err != nil {
		return err
	} else {
		if !r.checkDeploymentReady(foundDeploy) {
			// deployment is not ready yet
			msg := "Deployment is not ready"
			instance.MarkSuggestionStatusDeploymentReady(corev1.ConditionFalse, SuggestionDeploymentNotReady, msg)
			return nil
		} else {
			msg := "Deployment is ready"
			instance.MarkSuggestionStatusDeploymentReady(corev1.ConditionTrue, SuggestionDeploymentReady, msg)
		}

	}
	experiment := &experimentsv1beta1.Experiment{}
	trials := &trialsv1beta1.TrialList{}

	if err := r.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace,
	}, experiment); err != nil {
		return err
	}

	if err := r.List(context.TODO(), trials, client.MatchingLabels(util.TrialLabels(experiment))); err != nil {
		return err
	}
	// TODO (andreyvelich): Do we want to run ValidateAlgorithmSettings when Experiment is restarting?
	// Currently it is running.
	if !instance.IsRunning() {
		if err = r.ValidateAlgorithmSettings(instance, experiment); err != nil {
			logger.Error(err, "Marking suggestion failed as algorithm settings validation failed")
			msg := fmt.Sprintf("Validation failed: %v", err)
			instance.MarkSuggestionStatusFailed(SuggestionFailedReason, msg)
			// return nil since it is a terminal condition
			return nil
		}
		if instance.Spec.EarlyStopping != nil {
			if err = r.ValidateEarlyStoppingSettings(instance, experiment); err != nil {
				logger.Error(err, "Marking suggestion failed as early stopping settings validation failed")
				msg := fmt.Sprintf("Validation failed: %v", err)
				instance.MarkSuggestionStatusFailed(SuggestionFailedReason, msg)
				// return nil since it is a terminal condition
				return nil
			}
		}
		msg := "Suggestion is running"
		instance.MarkSuggestionStatusRunning(corev1.ConditionTrue, SuggestionRunningReason, msg)
	}
	logger.Info("Sync assignments", "Suggestion Requests", instance.Spec.Requests,
		"Suggestion Count", instance.Status.SuggestionCount)
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
