package composer

import (
	"fmt"

	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"

	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibconfig"
)

const (
	defaultInitialDelaySeconds = 10
	defaultPeriodForReady      = 10
	defaultPeriodForLive       = 120
	defaultFailureThreshold    = 12
	// Ref https://github.com/grpc-ecosystem/grpc-health-probe/
	defaultGRPCHealthCheckProbe = "/bin/grpc_health_probe"
)

var (
	log              = logf.Log.WithName("suggestion-composer")
	ComposerRegistry = make(map[string]Composer)
)

type Composer interface {
	DesiredDeployment(s *suggestionsv1beta1.Suggestion) (*appsv1.Deployment, error)
	DesiredService(s *suggestionsv1beta1.Suggestion) (*corev1.Service, error)
	DesiredVolume(s *suggestionsv1beta1.Suggestion) (*corev1.PersistentVolumeClaim, *corev1.PersistentVolume, error)
	CreateComposer(mgr manager.Manager) Composer
}

type General struct {
	scheme *runtime.Scheme
	client.Client
}

func New(mgr manager.Manager) Composer {
	// We assume DefaultComposer always exists in ComposerRegistry.
	ptr, _ := ComposerRegistry[consts.DefaultComposer]
	return ptr.CreateComposer(mgr)
}

// DesiredDeployment returns desired deployment for suggestion
func (g *General) DesiredDeployment(s *suggestionsv1beta1.Suggestion) (*appsv1.Deployment, error) {

	suggestionConfigData, err := katibconfig.GetSuggestionConfigData(s.Spec.AlgorithmName, g.Client)
	if err != nil {
		return nil, err
	}

	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        util.GetAlgorithmDeploymentName(s),
			Namespace:   s.Namespace,
			Labels:      s.Labels,
			Annotations: s.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: util.SuggestionLabels(s),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      util.SuggestionLabels(s),
					Annotations: util.SuggestionAnnotations(s),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*g.desiredContainer(s, suggestionConfigData),
					},
				},
			},
		},
	}

	// Get Suggestion Service Account Name from config
	if suggestionConfigData.ServiceAccountName != "" {
		d.Spec.Template.Spec.ServiceAccountName = suggestionConfigData.ServiceAccountName
	}

	// Attach volume to the suggestion pod spec if ResumePolicy = FromVolume
	if s.Spec.ResumePolicy == experimentsv1beta1.FromVolume {
		d.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: consts.ContainerSuggestionVolumeName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: util.GetAlgorithmPersistentVolumeClaimName(s),
					},
				},
			},
		}
	}

	if err := controllerutil.SetControllerReference(s, d, g.scheme); err != nil {
		return nil, err
	}

	return d, nil
}

// DesiredService returns desired service for suggestion
func (g *General) DesiredService(s *suggestionsv1beta1.Suggestion) (*corev1.Service, error) {
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

func (g *General) desiredContainer(s *suggestionsv1beta1.Suggestion, suggestionConfigData katibconfig.SuggestionConfig) *corev1.Container {

	c := &corev1.Container{
		Name:            consts.ContainerSuggestion,
		Image:           suggestionConfigData.Image,
		ImagePullPolicy: suggestionConfigData.ImagePullPolicy,
		Ports: []corev1.ContainerPort{
			{
				Name:          consts.DefaultSuggestionPortName,
				ContainerPort: consts.DefaultSuggestionPort,
			},
		},
		Resources: suggestionConfigData.Resource,
	}

	if viper.GetBool(consts.ConfigEnableGRPCProbeInSuggestion) {
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
	}

	// Attach volume mounts to the suggestion container if ResumePolicy = FromVolume
	if s.Spec.ResumePolicy == experimentsv1beta1.FromVolume {
		c.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      consts.ContainerSuggestionVolumeName,
				MountPath: consts.DefaultContainerSuggestionVolumeMountPath,
			},
		}
	}
	return c
}

// DesiredVolume returns desired PVC and PV for suggestion.
// If StorageClassName != DefaultSuggestionStorageClassName returns only PVC.
func (g *General) DesiredVolume(s *suggestionsv1beta1.Suggestion) (*corev1.PersistentVolumeClaim, *corev1.PersistentVolume, error) {
	persistentVolumeName := util.GetAlgorithmPersistentVolumeName(s)

	// TODO (andreyvelich): Enable to specify these values from Katib config
	storageClassName := consts.DefaultSuggestionStorageClassName
	persistentVolumePath := consts.DefaultSuggestionVolumeLocalPathPrefix + persistentVolumeName
	volumeAccessModes := consts.DefaultSuggestionVolumeAccessMode

	volumeStorage, _ := resource.ParseQuantity(consts.DefaultSuggestionVolumeStorage)

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GetAlgorithmPersistentVolumeClaimName(s),
			Namespace: s.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				volumeAccessModes,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: volumeStorage,
				},
			},
		},
	}

	// Add owner reference to the pvc so that it could be GC after the suggestion is deleted
	if err := controllerutil.SetControllerReference(s, pvc, g.scheme); err != nil {
		return nil, nil, err
	}

	var pv *corev1.PersistentVolume
	// Create PV with local hostPath by default
	if storageClassName == consts.DefaultSuggestionStorageClassName {
		localLabel := map[string]string{"type": "local"}

		pv = &corev1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:   persistentVolumeName,
				Labels: localLabel,
			},
			Spec: corev1.PersistentVolumeSpec{
				StorageClassName: consts.DefaultSuggestionStorageClassName,
				AccessModes: []corev1.PersistentVolumeAccessMode{
					volumeAccessModes,
				},
				PersistentVolumeSource: corev1.PersistentVolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: persistentVolumePath,
					},
				},
				Capacity: corev1.ResourceList{
					corev1.ResourceStorage: volumeStorage,
				},
			},
		}

		// Add owner reference to the pv so that it could be GC after the suggestion is deleted
		if err := controllerutil.SetControllerReference(s, pv, g.scheme); err != nil {
			return nil, nil, err
		}

	}

	return pvc, pv, nil
}

func (g *General) CreateComposer(mgr manager.Manager) Composer {
	return &General{mgr.GetScheme(), mgr.GetClient()}
}

func init() {
	ComposerRegistry[consts.DefaultComposer] = &General{}
}
