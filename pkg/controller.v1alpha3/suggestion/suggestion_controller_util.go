package suggestion

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileSuggestion) createOrUpdateDeployment(deploy *appsv1.Deployment) (*appsv1.Deployment, error) {
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

func (r *ReconcileSuggestion) createOrUpdateService(service *corev1.Service) (*corev1.Service, error) {
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
