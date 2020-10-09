package composer

import (
	"fmt"

	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"

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
	DesiredRBAC(s *suggestionsv1beta1.Suggestion) (*corev1.ServiceAccount, *rbacv1.Role, *rbacv1.RoleBinding, error)
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
			Name:        util.GetSuggestionDeploymentName(s),
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
					Containers: g.desiredContainers(s, suggestionConfigData),
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
						ClaimName: util.GetSuggestionPersistentVolumeClaimName(s),
					},
				},
			},
		}
	}

	// Attach service account if early stopping is used
	if s.Spec.EarlyStoppingAlgorithmName != "" {
		d.Spec.Template.Spec.ServiceAccountName = util.GetSuggestionRBACName(s)
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
	if s.Spec.EarlyStoppingAlgorithmName != "" {
		earlyStoppingPort := corev1.ServicePort{
			Name: consts.DefaultEarlyStoppingPortName,
			Port: consts.DefaultEarlyStoppingPort,
		}
		ports = append(ports, earlyStoppingPort)
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GetSuggestionServiceName(s),
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

func (g *General) desiredContainers(s *suggestionsv1beta1.Suggestion, suggestionConfigData katibconfig.SuggestionConfig) []corev1.Container {

	containers := []corev1.Container{}
	suggestionContainer := corev1.Container{
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
		suggestionContainer.ReadinessProbe = &corev1.Probe{
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
		suggestionContainer.LivenessProbe = &corev1.Probe{
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
		suggestionContainer.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      consts.ContainerSuggestionVolumeName,
				MountPath: suggestionConfigData.VolumeMountPath,
			},
		}
	}
	containers = append(containers, suggestionContainer)

	// TODO (andreyvelich): Take parameters from the config
	if s.Spec.EarlyStoppingAlgorithmName != "" {
		earlyStoppingContainer := corev1.Container{
			Name:            consts.ContainerEarlyStopping,
			Image:           "docker.io/andreyvelichkevich/earlystopping-median",
			ImagePullPolicy: "Always",
			Ports: []corev1.ContainerPort{
				{
					Name:          consts.DefaultEarlyStoppingPortName,
					ContainerPort: consts.DefaultEarlyStoppingPort,
				},
			},
			// TODO (andreyvelich): Change to Early Stopping
			Resources: suggestionConfigData.Resource,
		}

		containers = append(containers, earlyStoppingContainer)
	}
	return containers
}

// DesiredVolume returns desired PVC and PV for suggestion.
// If StorageClassName != DefaultSuggestionStorageClassName returns only PVC.
func (g *General) DesiredVolume(s *suggestionsv1beta1.Suggestion) (*corev1.PersistentVolumeClaim, *corev1.PersistentVolume, error) {

	suggestionConfigData, err := katibconfig.GetSuggestionConfigData(s.Spec.AlgorithmName, g.Client)
	if err != nil {
		return nil, nil, err
	}

	persistentVolumeName := util.GetSuggestionPersistentVolumeName(s)

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GetSuggestionPersistentVolumeClaimName(s),
			Namespace: s.Namespace,
		},
		Spec: suggestionConfigData.PersistentVolumeClaimSpec,
	}

	// Add owner reference to the pvc so that it could be GC after the suggestion is deleted
	if err := controllerutil.SetControllerReference(s, pvc, g.scheme); err != nil {
		return nil, nil, err
	}

	var pv *corev1.PersistentVolume
	// Create PV with local hostPath by default
	if *pvc.Spec.StorageClassName == consts.DefaultSuggestionStorageClassName {
		localLabel := map[string]string{"type": "local"}

		pv = &corev1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:   persistentVolumeName,
				Labels: localLabel,
			},
			Spec: suggestionConfigData.PersistentVolumeSpec,
		}

		// If default host path is specified attach pv name to the path.
		// Full default local path = DefaultSuggestionVolumeLocalPathPrefix<suggestion-name>-<suggestion-algorithm>-<suggestion-namespace>
		if pv.Spec.PersistentVolumeSource.HostPath != nil &&
			pv.Spec.PersistentVolumeSource.HostPath.Path == consts.DefaultSuggestionVolumeLocalPathPrefix {
			pv.Spec.PersistentVolumeSource.HostPath.Path = pv.Spec.PersistentVolumeSource.HostPath.Path + persistentVolumeName
		}

		// Add owner reference to the pv so that it could be GC after the suggestion is deleted
		if err := controllerutil.SetControllerReference(s, pv, g.scheme); err != nil {
			return nil, nil, err
		}

	}

	return pvc, pv, nil
}

// DesiredRBAC returns desired ServiceAccount, Role and RoleBinding for the Suggestion
func (g *General) DesiredRBAC(s *suggestionsv1beta1.Suggestion) (*corev1.ServiceAccount, *rbacv1.Role, *rbacv1.RoleBinding, error) {

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GetSuggestionRBACName(s),
			Namespace: s.Namespace,
		},
	}

	// Add owner reference to the ServiceAccount so that it could be GC after the suggestion is deleted
	if err := controllerutil.SetControllerReference(s, serviceAccount, g.scheme); err != nil {
		return nil, nil, nil, err
	}

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GetSuggestionRBACName(s),
			Namespace: s.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					trialsv1beta1.Group,
				},
				Resources: []string{
					consts.PluralTrial,
					fmt.Sprintf("%v/status", consts.PluralTrial),
				},
				Verbs: []string{
					rbacv1.VerbAll,
				},
			},
		},
	}

	// Add owner reference to the Role so that it could be GC after the suggestion is deleted
	if err := controllerutil.SetControllerReference(s, role, g.scheme); err != nil {
		return nil, nil, nil, err
	}

	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GetSuggestionRBACName(s),
			Namespace: s.Namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      util.GetSuggestionRBACName(s),
				Namespace: s.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     util.GetSuggestionRBACName(s),
		},
	}

	// Add owner reference to the RoleBinding so that it could be GC after the suggestion is deleted
	if err := controllerutil.SetControllerReference(s, roleBinding, g.scheme); err != nil {
		return nil, nil, nil, err
	}

	return serviceAccount, role, roleBinding, nil
}

// CreateComposer create instance of composer interface with given manager
func (g *General) CreateComposer(mgr manager.Manager) Composer {
	return &General{mgr.GetScheme(), mgr.GetClient()}
}

func init() {
	ComposerRegistry[consts.DefaultComposer] = &General{}
}
