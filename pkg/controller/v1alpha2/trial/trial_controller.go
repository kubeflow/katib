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
	"bytes"
	"context"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/trial/util"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Trial Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTrial{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("trial-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Trial
	err = c.Watch(&source.Kind{Type: &trialsv1alpha2.Trial{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileTrial{}

// ReconcileTrial reconciles a Trial object
type ReconcileTrial struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Trial object and makes changes based on the state read
// and what is in the Trial.Spec
// +kubebuilder:rbac:groups=trials.kubeflow.org,resources=trials,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=trials.kubeflow.org,resources=trials/status,verbs=get;update;patch
func (r *ReconcileTrial) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Trial instance
	instance := &trialsv1alpha2.Trial{}
	requeue := false
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	original := instance.DeepCopy()

	if instance.IsCompleted() {

		return reconcile.Result{}, nil

	}
	if !instance.IsCreated() {
		//Trial not created in DB
		err = util.CreateTrialinDB(instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		if instance.Status.StartTime == nil {
			now := metav1.Now()
			instance.Status.StartTime = &now
		}
		msg := "Trial is created"
		instance.MarkTrialStatusCreated(util.TrialCreatedReason, msg)
		requeue = true

	} else {
		// Trial already created in DB
		err := r.reconcileTrial(instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(original.Status, instance.Status) {
		//assuming that only status change
		err = util.UpdateTrialStatusinDB(instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{Requeue: requeue}, nil
}

func (r *ReconcileTrial) reconcileTrial(instance *trialsv1alpha2.Trial) error {

	var err error
	desiredJob, err := r.getDesiredJobSpec(instance)
	if err != nil {
		log.Info("Error in getting Job Spec from instance")
		return err
	}

	deployedJob, err := r.reconcileJob(instance, desiredJob)
	if err != nil {
		return err
	}

	//Job already exists
	//TODO Can desired Spec differ from deployedSpec?
	if deployedJob != nil {
		if err = util.UpdateTrialStatusCondition(instance, deployedJob); err != nil {
			return err
		}
		if err = util.UpdateTrialStatusObservation(instance, deployedJob); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReconcileTrial) reconcileJob(instance *trialsv1alpha2.Trial, desiredJob *unstructured.Unstructured) (*unstructured.Unstructured, error) {

	var err error
	apiVersion := desiredJob.GetAPIVersion()
	kind := desiredJob.GetKind()
	gvk := schema.FromAPIVersionAndKind(apiVersion, kind)

	deployedJob := &unstructured.Unstructured{}
	deployedJob.SetGroupVersionKind(gvk)
	err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, deployedJob)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Job", "namespace", instance.Namespace, "name", instance.Name, "kind", kind)
			err = r.Create(context.TODO(), desiredJob)
			if err != nil {
				log.Info("Error in creating job: %v ", err)
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	//TODO create Metric colletor

	msg := "Trial is running"
	instance.MarkTrialStatusRunning(util.TrialRunningReason, msg)
	return deployedJob, nil
}

func (r *ReconcileTrial) getDesiredJobSpec(instance *trialsv1alpha2.Trial) (*unstructured.Unstructured, error) {
	buf := bytes.NewBufferString(instance.Spec.RunSpec)
	bufSize := 1024

	desiredJobSpec := &unstructured.Unstructured{}
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, bufSize).Decode(desiredJobSpec); err != nil {
		log.Info("Yaml decode error %v", err)
		return nil, err
	}
	if err := controllerutil.SetControllerReference(instance, desiredJobSpec, r.scheme); err != nil {
		log.Info("SetControllerReference error %v", err)
		return nil, err
	}

	return desiredJobSpec, nil
}
