package composer

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/util"
	"github.com/kubeflow/katib/pkg/util/v1alpha3/katibconfig"
)

const (
	defaultInitialDelaySeconds = 10
	defaultPeriodForReady      = 10
	defaultPeriodForLive       = 120
	defaultFailureThreshold    = 12
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
			Name:      util.GetAlgorithmDeploymentName(s),
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
			Name:      util.GetAlgorithmServiceName(s),
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
	suggestionConfigData, err := katibconfig.GetSuggestionConfigData(s.Spec.AlgorithmName, g.Client)
	if err != nil {
		return nil, err
	}
	// Get Suggestion data from config
	suggestionContainerImage := suggestionConfigData[consts.LabelSuggestionImageTag]
	suggestionCPULimit := suggestionConfigData[consts.LabelSuggestionCPULimitTag]
	suggestionCPURequest := suggestionConfigData[consts.LabelSuggestionCPURequestTag]
	suggestionMemLimit := suggestionConfigData[consts.LabelSuggestionMemLimitTag]
	suggestionMemRequest := suggestionConfigData[consts.LabelSuggestionMemRequestTag]
	suggestionDiskLimit := suggestionConfigData[consts.LabelSuggestionDiskLimitTag]
	suggestionDiskRequest := suggestionConfigData[consts.LabelSuggestionDiskRequestTag]
	c := &corev1.Container{
		Name: consts.ContainerSuggestion,
	}
	c.Image = suggestionContainerImage
	c.ImagePullPolicy = corev1.PullIfNotPresent
	c.Ports = []corev1.ContainerPort{
		{
			Name:          consts.DefaultSuggestionPortName,
			ContainerPort: consts.DefaultSuggestionPort,
		},
	}

	cpuLimitQuantity, err := resource.ParseQuantity(suggestionCPULimit)
	if err != nil {
		return nil, err
	}
	cpuRequestQuantity, err := resource.ParseQuantity(suggestionCPURequest)
	if err != nil {
		return nil, err
	}
	memLimitQuantity, err := resource.ParseQuantity(suggestionMemLimit)
	if err != nil {
		return nil, err
	}
	memRequestQuantity, err := resource.ParseQuantity(suggestionMemRequest)
	if err != nil {
		return nil, err
	}
	diskLimitQuantity, err := resource.ParseQuantity(suggestionDiskLimit)
	if err != nil {
		return nil, err
	}
	diskRequestQuantity, err := resource.ParseQuantity(suggestionDiskRequest)
	if err != nil {
		return nil, err
	}

	c.Resources = corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:              cpuLimitQuantity,
			corev1.ResourceMemory:           memLimitQuantity,
			corev1.ResourceEphemeralStorage: diskLimitQuantity,
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:              cpuRequestQuantity,
			corev1.ResourceMemory:           memRequestQuantity,
			corev1.ResourceEphemeralStorage: diskRequestQuantity,
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
		PeriodSeconds:       defaultPeriodForReady,
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
		// Ref https://srcco.de/posts/kubernetes-liveness-probes-are-dangerous.html
		InitialDelaySeconds: defaultInitialDelaySeconds,
		PeriodSeconds:       defaultPeriodForLive,
		FailureThreshold:    defaultFailureThreshold,
	}
	return c, nil
}
