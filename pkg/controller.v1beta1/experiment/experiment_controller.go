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

package experiment

import (
	"context"
	"fmt"
	"sort"
	"time"

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
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/manifest"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/suggestion"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/util"
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
		recorder: mgr.GetEventRecorderFor(ControllerName),
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
		log.Error(err, "addWatch failed")
		return err
	}

	log.Info("Experiment controller created")
	return nil
}

// addWatch adds a new Controller to mgr with r as the reconcile.Reconciler
func addWatch(mgr manager.Manager, c controller.Controller) error {
	// Watch for changes to Experiment
	err := c.Watch(source.Kind(mgr.GetCache(), &experimentsv1beta1.Experiment{}), &handler.EnqueueRequestForObject{})
	if err != nil {
		log.Error(err, "Experiment watch failed")
		return err
	}

	// Watch for trials for the experiments
	eventHandler := handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &experimentsv1beta1.Experiment{}, handler.OnlyControllerOwner())
	if err = c.Watch(source.Kind(mgr.GetCache(), &trialsv1beta1.Trial{}), eventHandler); err != nil {
		log.Error(err, "Trial watch failed")
		return err
	}
	if err = c.Watch(source.Kind(mgr.GetCache(), &suggestionsv1beta1.Suggestion{}), eventHandler); err != nil {
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
func (r *ReconcileExperiment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Experiment instance
	logger := log.WithValues("Experiment", request.NamespacedName)
	original := &experimentsv1beta1.Experiment{}
	err := r.Get(ctx, request.NamespacedName, original)
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
		// Cleanup suggestion after Experiment is finished.
		// If ResumePolicy = Never or ResumePolicy = FromVolume, delete suggestion deployment, service and mark suggestion status succeeded.
		if instance.Spec.ResumePolicy == experimentsv1beta1.NeverResume || instance.Spec.ResumePolicy == experimentsv1beta1.FromVolume {
			err := r.cleanupSuggestionResources(instance)
			if err != nil {
				logger.Error(err, "cleanupSuggestionResources error")
				return reconcile.Result{}, err
			}
		}
		// Check if experiment is restartable and max trials is reconfigured
		// That means experiment is restarting
		if (util.IsCompletedExperimentRestartable(instance) &&
			instance.Spec.MaxTrialCount != nil &&
			*instance.Spec.MaxTrialCount > instance.Status.Trials) ||
			(instance.Spec.MaxTrialCount == nil && instance.Status.Trials != 0) {
			logger.Info("Experiment is restarting",
				"MaxTrialCount", instance.Spec.MaxTrialCount,
				"ParallelTrialCount", instance.Spec.ParallelTrialCount,
				"MaxFailedTrialCount", instance.Spec.MaxFailedTrialCount)
			msg := "Experiment is restarted"
			instance.MarkExperimentStatusRestarting(util.ExperimentRestartingReason, msg)
			// If ResumePolicy = FromVolume, suggestion must remove succeeded status and add running status when restarting
			if instance.Spec.ResumePolicy == experimentsv1beta1.FromVolume {
				err := r.restartSuggestion(instance)
				if err != nil {
					logger.Error(err, "restartSuggestion error")
					return reconcile.Result{}, err
				}
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
				corev1.EventTypeWarning, consts.ReconcileErrorReason,
				"Failed to reconcile: %v", err)
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(original.Status, instance.Status) {
		//assuming that only status change
		err = r.updateStatusHandler(instance)
		if err != nil {
			logger.Info("Update experiment instance status failed, reconciler requeued", "err", err)
			return reconcile.Result{
				Requeue: true,
			}, nil
		}
	}

	return reconcile.Result{}, nil
}

// ReconcileExperiment is the main reconcile loop.
func (r *ReconcileExperiment) ReconcileExperiment(instance *experimentsv1beta1.Experiment) error {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	trials := &trialsv1beta1.TrialList{}
	labels := map[string]string{consts.LabelExperimentName: instance.Name}

	if err := r.List(context.TODO(), trials, client.InNamespace(instance.Namespace), client.MatchingLabels(labels)); err != nil {
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
		return r.ReconcileTrials(instance, trials.Items)
	}

	return nil
}

// ReconcileTrials syncs trials.
func (r *ReconcileExperiment) ReconcileTrials(instance *experimentsv1beta1.Experiment, trials []trialsv1beta1.Trial) error {

	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})

	parallelCount := *instance.Spec.ParallelTrialCount
	activeCount := instance.Status.TrialsPending + instance.Status.TrialsRunning
	completedCount := instance.Status.TrialsSucceeded + instance.Status.TrialsFailed + instance.Status.TrialsKilled + instance.Status.TrialsEarlyStopped

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
			if err := r.createTrials(instance, trials, addCount); err != nil {
				logger.Error(err, "Create trials error")
				return err
			}
		}
	}

	return nil

}

func (r *ReconcileExperiment) createTrials(instance *experimentsv1beta1.Experiment, trialList []trialsv1beta1.Trial, addCount int32) error {

	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	logger.Info("Reconcile Suggestion", "addCount", addCount)
	trials, err := r.ReconcileSuggestions(instance, trialList, addCount)
	if err != nil {
		logger.Error(err, "Get suggestions error")
		return err
	}
	var trialNames []string
	for _, trial := range trials {
		trialInstance, err := r.getTrialInstance(instance, &trial)
		if err != nil {
			logger.Error(err, "Get trial instance error", "trial", trial)
			return err
		}
		// Due to unsynchronised policy of Kubernetes controllers, trial creation can fail.
		if err = r.Create(context.TODO(), trialInstance); err != nil {
			logger.Error(err, "Trial create error", "Trial name", trial.Name)
			continue
		}
		trialNames = append(trialNames, trial.Name)
	}
	// Print created Trial names
	if len(trialNames) != 0 {
		logger.Info("Created Trials", "trialNames", trialNames)
	}
	return nil
}

func (r *ReconcileExperiment) deleteTrials(instance *experimentsv1beta1.Experiment,
	trials []trialsv1beta1.Trial,
	expectedDeletions int32) error {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})

	trialSlice := trials
	sort.Slice(trialSlice, func(i, j int) bool {
		return trialSlice[i].CreationTimestamp.Time.
			After(trialSlice[j].CreationTimestamp.Time)
	})

	expected := int(expectedDeletions)
	actual := len(trialSlice)
	// If the number of trials < expected, we delete all we have.
	if actual < expected {
		logger.Info("deleteTrials does not find enough trials, we will delete all trials instead",
			"expectedDeletions", expected, "trials", actual)
		expected = actual
	}
	deletedNames := []string{}
	for i := 0; i < expected; i++ {
		if err := r.Delete(context.TODO(), &trialSlice[i]); err != nil {
			logger.Error(err, "Trial Delete error")
			return err
		}
		deletedNames = append(deletedNames, trialSlice[i].Name)
	}

	// Check if trials were deleted
	timeout := 60 * time.Second
	endTime := time.Now().Add(timeout)

	for _, name := range deletedNames {
		var err error
		for !errors.IsNotFound(err) && time.Now().Before(endTime) {
			err = r.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: instance.GetNamespace()}, &trialsv1beta1.Trial{})
		}
		// If trials were deleted, err == IsNotFound, return error if timeout is out
		if !errors.IsNotFound(err) {
			return fmt.Errorf("Unable to delete trials %v, error: %v", deletedNames, err)
		}
	}

	// We have to delete trials from suggestion status and update SuggestionCount
	suggestion := &suggestionsv1beta1.Suggestion{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()}, suggestion)
	if err != nil {
		logger.Error(err, "Suggestion Get error")
		return err
	}

	deletedNamesMap := make(map[string]bool)
	for _, name := range deletedNames {
		deletedNamesMap[name] = true
	}
	// Create new Trial Assignment without deleted trials
	newTrialAssignment := []suggestionsv1beta1.TrialAssignment{}
	for _, ta := range suggestion.Status.Suggestions {
		if _, ok := deletedNamesMap[ta.Name]; !ok {
			newTrialAssignment = append(newTrialAssignment, ta)
		}
	}

	// Update suggestion spec first.
	// If Requests <= SuggestionCount suggestion controller returns nil.
	suggestion.Spec.Requests = int32(len(newTrialAssignment))
	if err := r.UpdateSuggestion(suggestion); err != nil {
		return err
	}

	// Update suggestion status
	suggestion.Status.Suggestions = newTrialAssignment
	suggestion.Status.SuggestionCount = int32(len(newTrialAssignment))
	if err := r.UpdateSuggestionStatus(suggestion); err != nil {
		return err
	}

	logger.Info("Trials were successfully deleted", "trialNames", deletedNames)

	return nil
}

// ReconcileSuggestions gets or creates the suggestion if needed.
func (r *ReconcileExperiment) ReconcileSuggestions(instance *experimentsv1beta1.Experiment, trialList []trialsv1beta1.Trial, addCount int32) ([]suggestionsv1beta1.TrialAssignment, error) {
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	var assignments []suggestionsv1beta1.TrialAssignment
	currentCount := int32(len(trialList))
	incompleteEarlyStoppingCount := int32(0)
	trialNames := map[string]bool{}
	for _, trial := range trialList {
		trialNames[trial.Name] = true

		if !trial.IsObservationAvailable() && trial.IsEarlyStopped() {
			incompleteEarlyStoppingCount += 1
		}
	}

	// Exclude the number of incomplete early stopping trials from total suggestion requests
	// This means that no new trials will be requested from suggestion service until observations are ready for early stopped trials
	suggestionRequestsCount := currentCount + addCount - incompleteEarlyStoppingCount

	logger.Info("GetOrCreateSuggestion", "name", instance.Name, "Suggestion Requests", suggestionRequestsCount)
	original, err := r.GetOrCreateSuggestion(instance, suggestionRequestsCount)
	if err != nil {
		logger.Error(err, "GetOrCreateSuggestion failed", "name", instance.Name, "Suggestion Requests", suggestionRequestsCount)
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
					for _, suggestion := range suggestions {
						if !trialNames[suggestion.Name] {
							assignments = append(assignments, suggestion)
						}
					}
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
