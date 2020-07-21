package suggestion

import (
	"context"

	"github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// TODO (andreyvelich): Setup logger with suggestion name/namespace here

func (r *ReconcileSuggestion) reconcileDeployment(deploy *appsv1.Deployment) (*appsv1.Deployment, error) {
	foundDeploy := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, foundDeploy)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
		err = r.Create(context.TODO(), deploy)
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return foundDeploy, nil
}

func (r *ReconcileSuggestion) reconcileService(service *corev1.Service) (*corev1.Service, error) {
	foundService := &corev1.Service{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Service", "namespace", service.Namespace, "name", service.Name)
		err = r.Create(context.TODO(), service)
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return foundService, nil
}

func (r *ReconcileSuggestion) reconcileVolume(pvc *corev1.PersistentVolumeClaim, pv *corev1.PersistentVolume) (*corev1.PersistentVolumeClaim, *corev1.PersistentVolume, error) {
	foundPVC := &corev1.PersistentVolumeClaim{}
	foundPV := &corev1.PersistentVolume{}

	// Try to find/create PV, if PV has to be created
	if pv != nil {
		err := r.Get(context.TODO(), types.NamespacedName{Name: pv.Name}, foundPV)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating Persistent Volume", "name", pv.Name)
			err = r.Create(context.TODO(), pv)
			// Return only if Create was failed, otherwise try to find/create PVC
			if err != nil {
				return nil, nil, err
			}
		} else if err != nil {
			return nil, nil, err
		}
	}

	// Try to find/create PVC
	err := r.Get(context.TODO(), types.NamespacedName{Name: pvc.Name, Namespace: pvc.Namespace}, foundPVC)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Persistent Volume Claim", "namespace", pvc.Namespace, "name", pvc.Name)
		err = r.Create(context.TODO(), pvc)
		return nil, nil, err
	} else if err != nil {
		return nil, nil, err
	}
	return foundPVC, foundPV, nil
}

func (r *ReconcileSuggestion) deleteDeployment(instance *v1beta1.Suggestion) error {
	deploy, err := r.DesiredDeployment(instance)
	if err != nil {
		return err
	}
	realDeploy := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, realDeploy)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	log.Info("Deleting Suggestion Deployment", "namespace", realDeploy.Namespace, "name", realDeploy.Name)

	err = r.Delete(context.TODO(), realDeploy)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileSuggestion) deleteService(instance *v1beta1.Suggestion) error {
	service, err := r.DesiredService(instance)
	if err != nil {
		return err
	}
	realService := &corev1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, realService)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	log.Info("Deleting Suggestion Service", "namespace", realService.Namespace, "name", realService.Name)

	err = r.Delete(context.TODO(), realService)
	if err != nil {
		return err
	}

	return nil
}
