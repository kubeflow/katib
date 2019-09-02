package composer

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/consts"
)

type Composer interface {
	DesiredDeployment(s *suggestionsv1alpha2.Suggestion) (*appsv1.Deployment, error)
	DesiredService(s *suggestionsv1alpha2.Suggestion) (*corev1.Service, error)
}

type General struct {
	scheme *runtime.Scheme
}

func New(scheme *runtime.Scheme) Composer {
	return &General{
		scheme: scheme,
	}
}

func (g *General) DesiredDeployment(s *suggestionsv1alpha2.Suggestion) (*appsv1.Deployment, error) {
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"deployment":               s.Name,
					consts.LabelExperimentName: s.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"deployment":               s.Name,
						consts.LabelExperimentName: s.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*g.desiredContainer(s),
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

func (g *General) DesiredService(s *suggestionsv1alpha2.Suggestion) (*corev1.Service, error) {
	ports := []corev1.ServicePort{
		corev1.ServicePort{
			Name: "katib-api",
			Port: 6789,
		},
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"deployment":               s.Name,
				consts.LabelExperimentName: s.Name,
			},
			Ports: ports,
			Type:  corev1.ServiceTypeClusterIP,
		},
	}

	// Add owner reference to the service so that it could be GC after the suggestion is deleted
	if err := controllerutil.SetControllerReference(s, service, g.scheme); err != nil {
		return nil, err
	}

	return service, nil
}

// TODO(gaocegege): Implement switch logic.
func (g *General) desiredContainer(s *suggestionsv1alpha2.Suggestion) *corev1.Container {
	c := &corev1.Container{
		Name: consts.ContainerSuggestion,
	}
	c.Image = "katib/v1alpha2/suggestion-random"
	return c
}
