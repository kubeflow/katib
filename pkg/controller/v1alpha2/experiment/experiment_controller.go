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
	"bytes"
	"context"
	"os"
	"text/template"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	apiv1alpha2 "github.com/kubeflow/katib/pkg/api/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/experiment/util"
)

var log = logf.Log.WithName("experiment-controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Experiment Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileExperiment{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("experiment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		log.Error(err, "Failed to create experiment controller")
		return err
	}

	// Watch for changes to Experiment
	err = c.Watch(&source.Kind{Type: &experimentsv1alpha2.Experiment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		log.Error(err, "Experiment watch failed")
		return err
	}

	// Watch for trials for the experiments
	err = c.Watch(
		&source.Kind{Type: &trialsv1alpha2.Trial{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &experimentsv1alpha2.Experiment{},
		})

	if err != nil {
		log.Error(err, "Trial watch failed")
		return err
	}

	log.Info("Experiment controller created")
	return nil
}

var _ reconcile.Reconciler = &ReconcileExperiment{}

// ReconcileExperiment reconciles a Experiment object
type ReconcileExperiment struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Experiment object and makes changes based on the state read
// and what is in the Experiment.Spec
// +kubebuilder:rbac:groups=experiments.kubeflow.org,resources=experiments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=experiments.kubeflow.org,resources=experiments/status,verbs=get;update;patch
func (r *ReconcileExperiment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Experiment instance
	logger := log.WithValues("Experiment", request.NamespacedName)
	original := &experimentsv1alpha2.Experiment{}
	requeue := false
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

	if instance.IsCompleted() {

		return reconcile.Result{}, nil

	}
	if !instance.IsCreated() {
		//Experiment not created in DB
		err = util.CreateExperimentInDB(instance)
		if err != nil {
			logger.Error(err, "Create experiment in DB error")
			return reconcile.Result{}, err
		}

		if instance.Status.StartTime == nil {
			now := metav1.Now()
			instance.Status.StartTime = &now
		}
		msg := "Experiment is created"
		instance.MarkExperimentStatusCreated(util.ExperimentCreatedReason, msg)
		requeue = true
	} else {
		// Experiment already created in DB
		err := r.ReconcileExperiment(instance)
		if err != nil {
			logger.Error(err, "Reconcile experiment error")
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(original.Status, instance.Status) {
		//assuming that only status change
		err = util.UpdateExperimentStatusInDB(instance)
		if err != nil {
			logger.Error(err, "Update experiment status in DB error")
			return reconcile.Result{}, err
		}
		err = r.Status().Update(context.TODO(), instance)
		if err != nil {
			logger.Error(err, "Update experiment instance status error")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{Requeue: requeue}, nil
}

func (r *ReconcileExperiment) ReconcileExperiment(instance *experimentsv1alpha2.Experiment) error {

	var err error
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	trials := &trialsv1alpha2.TrialList{}
	labels := map[string]string{"experiment": instance.Name}
	lo := &client.ListOptions{}
	lo.MatchingLabels(labels).InNamespace(instance.Namespace)

	err = r.List(context.TODO(), lo, trials)
	if err != nil {
		logger.Error(err, "Trial List error")
		return err
	}
	if len(trials.Items) > 0 {
		err := util.UpdateExperimentStatus(instance, trials)
		if err != nil {
			logger.Error(err, "Update experiment status error")
			return err
		}
	}
	reconcileRequired := !instance.IsCompleted()
	if reconcileRequired {
		r.ReconcileTrials(instance)
	}
	return err
}

func (r *ReconcileExperiment) ReconcileTrials(instance *experimentsv1alpha2.Experiment) error {

	var err error
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	parallelCount := 0

	if instance.Spec.ParallelTrialCount != nil {
		parallelCount = *instance.Spec.ParallelTrialCount
	} else {
		parallelCount = 3
	}
	activeCount := instance.Status.TrialsPending + instance.Status.TrialsRunning
	completedCount := instance.Status.TrialsSucceeded + instance.Status.TrialsFailed + instance.Status.TrialsKilled

	if activeCount > parallelCount {
		deleteCount := activeCount - parallelCount
		if deleteCount > 0 {
			//delete 'deleteCount' number of trails. Sort them?
			logger.Info("DeleteTrials", "deleteCount", deleteCount)
			err = r.deleteTrials(instance, deleteCount)
		}

	} else if activeCount < parallelCount {
		requiredActiveCount := 0
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

		//create "addCount" number of trials
		logger.Info("CreateTrials", "addCount", addCount)
		err = r.createTrials(instance, addCount)

	}

	return err

}

type TrialTemplateParams struct {
	Experiment      string
	Trial           string
	NameSpace       string
	HyperParameters []*apiv1alpha2.ParameterAssignment
}

func (r *ReconcileExperiment) createTrials(instance *experimentsv1alpha2.Experiment, addCount int) error {

	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	trials, err := util.GetSuggestions(instance, addCount)
	/*trials := []apiv1alpha2.Trial{
		apiv1alpha2.Trial{Spec: &apiv1alpha2.TrialSpec{}}, apiv1alpha2.Trial{Spec: &apiv1alpha2.TrialSpec{}},
	}*/

	trialTemplate, err := r.getTrialTemplate(instance)
	if err != nil {
		logger.Error(err, "Get trial template error")
	}
	for _, elem := range trials {

		trial := &trialsv1alpha2.Trial{}
		trial.Name = string(uuid.NewUUID())
		trial.Namespace = instance.GetNamespace()
		trial.Labels = map[string]string{"experiment": instance.GetName()}

		if err := controllerutil.SetControllerReference(instance, trial, r.scheme); err != nil {
			logger.Error(err, "Set controller reference error")
		}

		trialParams := TrialTemplateParams{
			Experiment: instance.GetName(),
			Trial:      trial.Name,
			NameSpace:  trial.Namespace,
		}

		var buf bytes.Buffer
		if elem.Spec != nil && elem.Spec.ParameterAssignments != nil {
			for _, p := range elem.Spec.ParameterAssignments.Assignments {
				trialParams.HyperParameters = append(trialParams.HyperParameters, p)
			}
		}
		err = trialTemplate.Execute(&buf, trialParams)
		if err != nil {
			logger.Error(err, "Template execute error")
		}

		trial.Spec.RunSpec = buf.String()

		err := r.Create(context.TODO(), trial)
		if err != nil {
			logger.Error(err, "Trial create error", "Trial name", trial.Name)
		}

	}
	return err
}

func (r *ReconcileExperiment) getTrialTemplate(instance *experimentsv1alpha2.Experiment) (*template.Template, error) {

	var err error
	var tpl *template.Template = nil
	logger := log.WithValues("Experiment", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	trialTemplate := instance.Spec.TrialTemplate
	if trialTemplate != nil && trialTemplate.GoTemplate.RawTemplate != "" {
		tpl, err = template.New("Trial").Parse(trialTemplate.GoTemplate.RawTemplate)
	} else {
		//default values if user hasn't set
		configMapNS := os.Getenv("KATIB_CORE_NAMESPACE")
		configMapName := "trial-template"
		templatePath := "defaultTrialTemplate.yaml"
		if trialTemplate != nil && trialTemplate.GoTemplate.TemplateSpec != nil {
			templateSpec := trialTemplate.GoTemplate.TemplateSpec
			if templateSpec.ConfigMapName != "" {
				configMapName = templateSpec.ConfigMapName
			}
			if templateSpec.ConfigMapNamespace != "" {
				configMapNS = templateSpec.ConfigMapNamespace
			}
			if templateSpec.TemplatePath != "" {
				templatePath = templateSpec.TemplatePath
			}
		}
		configMap, err := r.getConfigMap(configMapName, configMapNS)
		if err != nil {
			logger.Error(err, "Get config map error", "configMapName", configMapName, "configMapNS", configMapNS)
		}
		if configMapTemplate, ok := configMap[templatePath]; !ok {
		} else {
			tpl, err = template.New("Trial").Parse(configMapTemplate)
		}
	}
	if err != nil {
		logger.Error(err, "Template parse error")
	}

	return tpl, err
}

func (r *ReconcileExperiment) deleteTrials(instance *experimentsv1alpha2.Experiment, deleteCount int) error {

	return nil
}

func (r *ReconcileExperiment) getConfigMap(name, namespace string) (map[string]string, error) {

	return nil, nil
}
