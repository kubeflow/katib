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

package experiment

import (
	"context"
	"sort"

	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/manifest"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/suggestion"
	suggestionfake "github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/suggestion/fake"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/util"
)

const (
	// ControllerName is the controller name.
	ControllerName = "experiment-controller"
)

var (
	log = logf.Log.WithName(ControllerName)
)

// Add creates a new Experiment Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	r := &ReconcileExperiment{
		Client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetRecorder(ControllerName),
	}
	imp := viper.GetString(consts.ConfigExperimentSuggestionName)
	r.Suggestion = newSuggestion(imp, mgr.GetScheme(), mgr.GetClient())

	r.Generator = manifest.New(r.Client)
	r.updateStatusHandler = r.updateStatus
	r.collector = util.NewExpsCollector(mgr.GetCache(), metrics.Registry)
	return r
}

// newSuggestion returns the new Suggestion for the given config.
func newSuggestion(config string, scheme *runtime.Scheme, client client.Client) suggestion.Suggestion {
	// Use different implementation according to the configuration.
	switch config {
	case "fake":
		log.Info("Using the fake suggestion implementation")
		return suggestionfake.New()
	case "default":
		log.Info("Using the default suggestion implementation")
		return suggestion.New(scheme, client)
	default:
		log.Info("No valid name specified, using the default suggestion implementation",
			"implementation", config)
		return suggestion.New(scheme, client)
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("experiment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		log.Error(err, "Failed to create experiment controller")
		return err
	}

	if err = addWatch(mgr, c); err != nil {
		log.Error(err, "Trial watch failed")
		return err
	}

	log.Info("Experiment controller created")
	return nil
}

// addWatch adds a new Controller to mgr with r as the reconcile.Reconciler
func addWatch(mgr manager.Manager, c controller.Controller) error {
	// Watch for changes to Experiment
	err := c.Watch(&source.Kind{Type: &experimentsv1alpha3.Experiment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		log.Error(err, "Experiment watch failed")
		return err
	}

	// Watch for trials for the experiments
	err = c.Watch(
		&source.Kind{Type: &trialsv1alpha3.Trial{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &experimentsv1alpha3.Experiment{},
		})

	if err != nil {
		log.Error(err, "Trial watch failed")
		return err
	}

	err = c.Watch(
		&source.Kind{Type: &suggestionsv1alpha3.Suggestion{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &experimentsv1alpha3.Experiment{},
		})

	if err != nil {
		log.Error(err, "Suggestion watch failed")
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcileExperiment{}

// ReconcileExperiment reconciles a Experiment object
type ReconcileExperiment struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder

	suggestion.Suggestion
	manifest.Generator
	// updateStatusHandler is defined for test purpose.
	updateStatusHandler updateStatusFunc
	// collector is a wrapper for experiment metrics.
	collector *util.ExperimentsCollector
}

// Reconcile reads that state of the cluster for a Experiment object and makes changes based on the state read
// and what is in the Experiment.Spec
// +kubebuilder:rbac:groups=experiments.kubeflow.org,resources=experiments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=experiments.kubeflow.org,resources=experiments/status,verbs=get;update;patch
func (r *ReconcileExperiment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Experiment instance
	logger := log.WithValues("Experiment", request.NamespacedName)
	original := &experimentsv1alpha3.Experiment{}
	err := r.Get(context.TODO(), request.NamespacedName, original)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Experiment Get error")
		return reconcile.Result{}, err
	}
	instance := original.DeepCopy()

	if needUpdate, finalizers := needUpdateFinalizers(instance); needUpdate {
		return r.updateFinalizers(instance, finalizers)
	}

	if instance.IsCompleted() {
		// Check if completed instance is restartable
		// Experiment is restartable only if it is in succeeded state by reaching max trials
		if util.IsCompletedExperimentRestartable(instance) {
			// Check if max trials is reconfigured
			if (instance.Spec.MaxTrialCount != nil &&
				*instance.Spec.MaxTrialCount != instance.Status.Trials) ||
				(instance.Spec.MaxTrialCount == nil && instance.Status.Trials != 0) {
				msg := "Experiment is restarted"
				instance.MarkExperimentStatusRestarting(util.ExperimentRestartingReason, msg)
			}
		} else {
			// If experiment is completed with no running trials, stop reconcile
			if !instance.HasRunningTrials() {
				return reconcile.Result{}, nil
			}
		}
	}
	if !instance.IsCreated() {
		if instance.Status.StartTime == nil {
			now := metav1.Now()
			instance.Status.StartTime = &now
		}
		if instance.Status.CompletionTime == nil {
			instance.Status.CompletionTime = &metav1.Time{}
		}
		msg := "Experiment is created"
		instance.MarkExperimentStatusCreated(util.ExperimentCreatedReason, msg)
	} else {
		err := r.ReconcileExperiment(instance)
		if err != nil {
			logger.Error(err, "Reconcile experiment error")
			r.recorder.Eventf(instance,
				corev1.EventTypeWarning, ReconcileFailedReason,
				"Failed to reconcile: %v", err)
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(original.Status, instance.Status) {
		//assuming that only status change
		err = r.updateStatusHandler(instance)
		if err != nil {
			logger.Error(err, "Update experiment instance status error")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// ReconcileExperiment is the main reconcile loop.
func (r *ReconcileExperiment) ReconcileExperiment(instance *experimentsv1alpha3.Experiment) error {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	trials := &trialsv1alpha3.TrialList{}
	labels := map[string]string{consts.LabelExperimentName: instance.Name}
	lo := &client.ListOptions{}
	lo.MatchingLabels(labels).InNamespace(instance.Namespace)

	if err := r.List(context.TODO(), lo, trials); err != nil {
		logger.Error(err, "Trial List error")
		return err
	}
	if len(trials.Items) > 0 {
		if err := util.UpdateExperimentStatus(r.collector, instance, trials); err != nil {
			logger.Error(err, "Update experiment status error")
			return err
		}
	}
	reconcileRequired := !instance.IsCompleted()
	if reconcileRequired {
		r.ReconcileTrials(instance, trials.Items)
	} else if instance.Spec.ResumePolicy != experimentsv1alpha3.LongRunning {
		return r.terminateSuggestion(instance)
	}
	return nil
}

// ReconcileTrials syncs trials.
func (r *ReconcileExperiment) ReconcileTrials(instance *experimentsv1alpha3.Experiment, trials []trialsv1alpha3.Trial) error {

	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})

	parallelCount := *instance.Spec.ParallelTrialCount
	activeCount := instance.Status.TrialsPending + instance.Status.TrialsRunning
	completedCount := instance.Status.TrialsSucceeded + instance.Status.TrialsFailed + instance.Status.TrialsKilled

	if activeCount > parallelCount {
		deleteCount := activeCount - parallelCount
		if deleteCount > 0 {
			//delete 'deleteCount' number of trails. Sort them?
			logger.Info("DeleteTrials", "deleteCount", deleteCount)
			if err := r.deleteTrials(instance, trials, deleteCount); err != nil {
				logger.Error(err, "Delete trials error")
				return err
			}
		}

	} else if activeCount < parallelCount {
		var requiredActiveCount int32
		if instance.Spec.MaxTrialCount == nil {
			requiredActiveCount = parallelCount
		} else {
			requiredActiveCount = *instance.Spec.MaxTrialCount - completedCount
			if requiredActiveCount > parallelCount {
				requiredActiveCount = parallelCount
			}
		}

		addCount := requiredActiveCount - activeCount
		if addCount < 0 {
			logger.Info("Invalid setting", "requiredActiveCount", requiredActiveCount, "MaxTrialCount",
				*instance.Spec.MaxTrialCount, "CompletedCount", completedCount)
			addCount = 0
		}

		logger.Info("Statistics",
			"requiredActiveCount", requiredActiveCount,
			"parallelCount", parallelCount,
			"activeCount", activeCount,
			"completedCount", completedCount,
		)

		//skip if no trials need to be created
		if addCount > 0 {
			//create "addCount" number of trials
			logger.Info("CreateTrials", "addCount", addCount)
			if err := r.createTrials(instance, trials, addCount); err != nil {
				logger.Error(err, "Create trials error")
				return err
			}
		}
	}

	return nil

}

func (r *ReconcileExperiment) createTrials(instance *experimentsv1alpha3.Experiment, trialList []trialsv1alpha3.Trial, addCount int32) error {

	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	currentCount := int32(len(trialList))
	trials, err := r.ReconcileSuggestions(instance, currentCount, addCount)
	if err != nil {
		logger.Error(err, "Get suggestions error")
		return err
	}
	for _, trial := range trials {
		if err = r.createTrialInstance(instance, &trial); err != nil {
			logger.Error(err, "Create trial instance error", "trial", trial)
			continue
		}
	}
	return nil
}

func (r *ReconcileExperiment) deleteTrials(instance *experimentsv1alpha3.Experiment,
	trials []trialsv1alpha3.Trial,
	deleteCount int32) error {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})

	trialSlice := trials
	sort.Slice(trialSlice, func(i, j int) bool {
		return trialSlice[i].CreationTimestamp.Time.
			After(trialSlice[j].CreationTimestamp.Time)
	})

	for i := 0; i < int(deleteCount); i++ {
		if err := r.Delete(context.TODO(), &trialSlice[i]); err != nil {
			logger.Error(err, "Trial Delete error")
			return err
		}
	}
	return nil
}

// ReconcileSuggestions gets or creates the suggestion if needed.
func (r *ReconcileExperiment) ReconcileSuggestions(instance *experimentsv1alpha3.Experiment, currentCount, addCount int32) ([]suggestionsv1alpha3.TrialAssignment, error) {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	var assignments []suggestionsv1alpha3.TrialAssignment
	suggestionRequestsCount := currentCount + addCount

	original, err := r.GetOrCreateSuggestion(instance, suggestionRequestsCount)
	logger.Info("GetOrCreateSuggestion", "Instance name", instance.Name, "suggestionRequestsCount", suggestionRequestsCount)
	if err != nil {
		logger.Error(err, "GetOrCreateSuggestion failed", "instance", instance.Name, "suggestionRequestsCount", suggestionRequestsCount)
		return nil, err
	} else {
		if original != nil {
			if original.IsFailed() {
				msg := "Suggestion has failed"
				instance.MarkExperimentStatusFailed(util.ExperimentFailedReason, msg)
			} else {
				suggestion := original.DeepCopy()
				if len(suggestion.Status.Suggestions) > int(currentCount) {
					suggestions := suggestion.Status.Suggestions
					assignments = suggestions[currentCount:]
				}
				if suggestion.Spec.Requests != suggestionRequestsCount {
					suggestion.Spec.Requests = suggestionRequestsCount
					if err := r.UpdateSuggestion(suggestion); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return assignments, nil
}
