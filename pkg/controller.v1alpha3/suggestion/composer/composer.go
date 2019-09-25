package composer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/util"
)

const (
	defaultInitialDelaySeconds = 10
	defaultPeriod              = 10
	// Ref https://github.com/grpc-ecosystem/grpc-health-probe/
	defaultGRPCHealthCheckProbe = "/bin/grpc_health_probe"
)

var log = logf.Log.WithName("suggestion-composer")

type Composer interface {
	DesiredDeployment(s *suggestionsv1alpha3.Suggestion) (*appsv1.Deployment, error)
	DesiredService(s *suggestionsv1alpha3.Suggestion) (*corev1.Service, error)
}

type General struct {
	scheme *runtime.Scheme
	client.Client
}

func New(scheme *runtime.Scheme, client client.Client) Composer {
	return &General{
		scheme: scheme,
		Client: client,
	}
}

func (g *General) DesiredDeployment(s *suggestionsv1alpha3.Suggestion) (*appsv1.Deployment, error) {
	container, err := g.desiredContainer(s)
	if err != nil {
		log.Error(err, "Error in constructing container")
		return nil, err
	}
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: util.SuggestionLabels(s),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: util.SuggestionLabels(s),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*container,
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(s, d, g.scheme); err != nil {
		return nil, err
	}
	return d, nil
}

func (g *General) DesiredService(s *suggestionsv1alpha3.Suggestion) (*corev1.Service, error) {
	ports := []corev1.ServicePort{
		{
			Name: consts.DefaultSuggestionPortName,
			Port: consts.DefaultSuggestionPort,
		},
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: util.SuggestionLabels(s),
			Ports:    ports,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	// Add owner reference to the service so that it could be GC after the suggestion is deleted
	if err := controllerutil.SetControllerReference(s, service, g.scheme); err != nil {
		return nil, err
	}

	return service, nil
}

func (g *General) desiredContainer(s *suggestionsv1alpha3.Suggestion) (*corev1.Container, error) {
	suggestionContainerImage, err := g.getSuggestionContainerImage(s.Spec.AlgorithmName)
	if err != nil {
		return nil, err
	}
	c := &corev1.Container{
		Name: consts.ContainerSuggestion,
	}
	c.Image = suggestionContainerImage
	c.Ports = []corev1.ContainerPort{
		{
			Name:          consts.DefaultSuggestionPortName,
			ContainerPort: consts.DefaultSuggestionPort,
		},
	}
	c.ReadinessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{
					defaultGRPCHealthCheckProbe,
					fmt.Sprintf("-addr=:%d", consts.DefaultSuggestionPort),
					fmt.Sprintf("-service=%s", consts.DefaultGRPCService),
				},
			},
		},
		InitialDelaySeconds: defaultInitialDelaySeconds,
		PeriodSeconds:       defaultPeriod,
	}
	c.LivenessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{
					defaultGRPCHealthCheckProbe,
					fmt.Sprintf("-addr=:%d", consts.DefaultSuggestionPort),
					fmt.Sprintf("-service=%s", consts.DefaultGRPCService),
				},
			},
		},
		InitialDelaySeconds: defaultInitialDelaySeconds,
		PeriodSeconds:       defaultPeriod,
	}
	return c, nil
}

func (g *General) getSuggestionContainerImage(algorithmName string) (string, error) {
	configMap := &corev1.ConfigMap{}
	err := g.Client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		log.Error(err, "Failed to find config map", "name", consts.KatibConfigMapName)
		// Error reading the object - requeue the request.
		return "", err
	}
	if config, ok := configMap.Data[consts.LabelSuggestionTag]; ok {
		suggestionConfig := map[string]map[string]string{}
		if err := json.Unmarshal([]byte(config), &suggestionConfig); err != nil {
			log.Error(err, "Json Unmarshal error", "Config", config)
			return "", err
		}
		if imageConfig, ok := suggestionConfig[algorithmName]; ok {
			if image, yes := imageConfig[consts.LabelSuggestionImageTag]; yes {
				return image, nil
			} else {
				return "", errors.New("Failed to find " + consts.LabelSuggestionImageTag + " configuration for algorithm name " + algorithmName)
			}
		} else {
			return "", errors.New("Failed to find algorithm image mapping " + algorithmName)
		}
	} else {
		return "", errors.New("Failed to find algorithm image mapping in configmap " + consts.KatibConfigMapName)
	}
}
