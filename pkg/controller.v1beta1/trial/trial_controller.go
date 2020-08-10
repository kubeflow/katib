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

package trial

import (
	"context"
	"fmt"

	batchv1beta "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/trial/managerclient"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	jobv1beta1 "github.com/kubeflow/katib/pkg/job/v1beta1"
)

const (
	// ControllerName is the controller name.
	ControllerName = "trial-controller"
)

var (
	log = logf.Log.WithName(ControllerName)
)

// Add creates a new Trial Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	r := &ReconcileTrial{
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		ManagerClient: managerclient.New(),
		recorder:      mgr.GetRecorder(ControllerName),
		collector:     NewTrialsCollector(mgr.GetCache(), metrics.Registry),
	}
	r.updateStatusHandler = r.updateStatus
	return r
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("trial-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		log.Error(err, "Create trial controller error")
		return err
	}

	// Watch for changes to Trial
	err = c.Watch(&source.Kind{Type: &trialsv1beta1.Trial{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		log.Error(err, "Trial watch error")
		return err
	}

	// Watch for changes to Cronjob
	err = c.Watch(
		&source.Kind{Type: &batchv1beta.CronJob{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &trialsv1beta1.Trial{},
		})

	if err != nil {
		log.Error(err, "CronJob watch error")
		return err
	}

	for _, gvk := range jobv1beta1.SupportedJobList {
		unstructuredJob := &unstructured.Unstructured{}
		unstructuredJob.SetGroupVersionKind(gvk)
		err = c.Watch(
			&source.Kind{Type: unstructuredJob},
			&handler.EnqueueRequestForOwner{
				IsController: true,
				OwnerType:    &trialsv1beta1.Trial{},
			})
		if err != nil {
			if meta.IsNoMatchError(err) {
				log.Info("Job watch error. CRD might be missing. Please install CRD and restart katib-controller", "CRD Kind", gvk.Kind)
				continue
			}
			return err
		} else {
			log.Info("Job watch added successfully", "CRD Kind", gvk.Kind)
		}
	}
	log.Info("Trial controller created")
	return nil
}

var _ reconcile.Reconciler = &ReconcileTrial{}

// ReconcileTrial reconciles a Trial object
type ReconcileTrial struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder

	managerclient.ManagerClient
	// updateStatusHandler is defined for test purpose.
	updateStatusHandler updateStatusFunc
	// collector is a wrapper for experiment metrics.
	collector *TrialsCollector
}

// Reconcile reads that state of the cluster for a Trial object and makes changes based on the state read
// and what is in the Trial.Spec
// +kubebuilder:rbac:groups=trials.kubeflow.org,resources=trials,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=trials.kubeflow.org,resources=trials/status,verbs=get;update;patch
func (r *ReconcileTrial) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Trial instance
	logger := log.WithValues("Trial", request.NamespacedName)
	original := &trialsv1beta1.Trial{}
	err := r.Get(context.TODO(), request.NamespacedName, original)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Trial Get error")
		return reconcile.Result{}, err
	}

	instance := original.DeepCopy()

	if needUpdate, finalizers := needUpdateFinalizers(instance); needUpdate {
		return r.updateFinalizers(instance, finalizers)
	}

	if !instance.IsCreated() {
		if instance.Status.StartTime == nil {
			now := metav1.Now()
			instance.Status.StartTime = &now
		}
		if instance.Status.CompletionTime == nil {
			instance.Status.CompletionTime = &metav1.Time{}
		}
		msg := "Trial is created"
		instance.MarkTrialStatusCreated(TrialCreatedReason, msg)
	} else {
		err := r.reconcileTrial(instance)
		if err != nil {
			logger.Error(err, "Reconcile trial error")
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
			logger.Info("Update trial instance status failed, reconciler requeued", "err", err)
			return reconcile.Result{
				Requeue: true,
			}, nil
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileTrial) reconcileTrial(instance *trialsv1beta1.Trial) error {

	var err error
	logger := log.WithValues("Trial", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	desiredJob, err := r.getDesiredJobSpec(instance)
	if err != nil {
		logger.Error(err, "Job Spec Get error")
		return err
	}

	deployedJob, err := r.reconcileJob(instance, desiredJob)
	if err != nil {
		logger.Error(err, "Reconcile job error")
		return err
	}

	// Job already exists
	// TODO Can desired Spec differ from deployedSpec?
	if deployedJob != nil {
		kind := deployedJob.GetKind()
		jobProvider, err := jobv1beta1.New(kind)
		if err != nil {
			logger.Error(err, "Failed to create the provider")
			return err
		}
		// Currently jobCondition - part of commonv1 TF package for all jobs
		jobCondition, err := jobProvider.GetDeployedJobStatus(deployedJob)
		if err != nil {
			logger.Error(err, "Get deployed status error")
			return err
		}

		// Update trial observation when the job is succeeded.
		if isJobSucceeded(jobCondition) {
			if err = r.UpdateTrialStatusObservation(instance, deployedJob); err != nil {
				logger.Error(err, "Update trial status observation error")
				return err
			}
		}

		// Update Trial job status only
		//    if job has succeeded and if observation field is available.
		//    if job has failed
		// This will ensure that trial is set to be complete only if metric is collected at least once
		r.UpdateTrialStatusCondition(instance, deployedJob, jobCondition)

	}
	return nil
}

func (r *ReconcileTrial) reconcileJob(instance *trialsv1beta1.Trial, desiredJob *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var err error
	logger := log.WithValues("Trial", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	apiVersion := desiredJob.GetAPIVersion()
	kind := desiredJob.GetKind()
	gvk := schema.FromAPIVersionAndKind(apiVersion, kind)

	// Add annotation to desired Job to disable istio sidecar
	err = util.TrainingJobAnnotations(desiredJob)
	if err != nil {
		logger.Error(err, "TrainingJobAnnotations error")
		return nil, err
	}

	deployedJob := &unstructured.Unstructured{}
	deployedJob.SetGroupVersionKind(gvk)
	err = r.Get(context.TODO(), types.NamespacedName{Name: desiredJob.GetName(), Namespace: desiredJob.GetNamespace()}, deployedJob)
	if err != nil {
		if errors.IsNotFound(err) {
			if instance.IsCompleted() {
				return nil, nil
			}
			jobProvider, err := jobv1beta1.New(desiredJob.GetKind())
			if err != nil {
				return nil, err
			}
			// mutate desiredJob according to provider
			if err := jobProvider.MutateJob(instance, desiredJob); err != nil {
				logger.Error(err, "Mutating desiredSpec of km.Training error")
				return nil, err
			}
			logger.Info("Creating Job", "kind", kind,
				"name", desiredJob.GetName())
			err = r.Create(context.TODO(), desiredJob)
			if err != nil {
				logger.Error(err, "Create job error")
				return nil, err
			}
			eventMsg := fmt.Sprintf("Job %s has been created", desiredJob.GetName())
			r.recorder.Eventf(instance, corev1.EventTypeNormal, JobCreatedReason, eventMsg)
		} else {
			logger.Error(err, "Trial Get error")
			return nil, err
		}
	} else {
		if instance.IsCompleted() && !instance.Spec.RetainRun {
			if err = r.Delete(context.TODO(), desiredJob, client.PropagationPolicy(metav1.DeletePropagationForeground)); err != nil {
				logger.Error(err, "Delete job error")
				return nil, err
			} else {
				eventMsg := fmt.Sprintf("Job %s has been deleted", desiredJob.GetName())
				r.recorder.Eventf(instance, corev1.EventTypeNormal, JobDeletedReason, eventMsg)
				return nil, nil
			}
		}
	}

	return deployedJob, nil
}

func (r *ReconcileTrial) getDesiredJobSpec(instance *trialsv1beta1.Trial) (*unstructured.Unstructured, error) {

	logger := log.WithValues("Trial", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})

	desiredJobSpec := instance.Spec.RunSpec

	if err := controllerutil.SetControllerReference(instance, desiredJobSpec, r.scheme); err != nil {
		logger.Error(err, "Set controller reference error")
		return nil, err
	}

	return desiredJobSpec, nil
}
