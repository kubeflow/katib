/*
Copyright 2022 The Kubeflow Authors.

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

package composer

import (
	"fmt"

	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
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
)

var (
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
	return ComposerRegistry[consts.DefaultComposer].CreateComposer(mgr)
}

// DesiredDeployment returns desired deployment for suggestion
func (g *General) DesiredDeployment(s *suggestionsv1beta1.Suggestion) (*appsv1.Deployment, error) {

	suggestionConfigData, err := katibconfig.GetSuggestionConfigData(s.Spec.Algorithm.AlgorithmName, g.Client)
	if err != nil {
		return nil, err
	}
	if containsContainerPortWithName(suggestionConfigData.Ports, consts.DefaultSuggestionPortName) ||
		containsContainerPort(suggestionConfigData.Ports, consts.DefaultSuggestionPort) {
		return nil, fmt.Errorf("invalid suggestion config: a port with name %q or number %d must not be specified",
			consts.DefaultSuggestionPortName, consts.DefaultSuggestionPort)
	}

	// If early stopping is used, get the config data.
	earlyStoppingConfigData := configv1beta1.EarlyStoppingConfig{}
	if s.Spec.EarlyStopping != nil && s.Spec.EarlyStopping.AlgorithmName != "" {
		earlyStoppingConfigData, err = katibconfig.GetEarlyStoppingConfigData(s.Spec.EarlyStopping.AlgorithmName, g.Client)
		if err != nil {
			return nil, err
		}
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
					Containers: g.desiredContainers(s, suggestionConfigData, earlyStoppingConfigData),
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

	// Attach ServiceAccount if early stopping is used.
	// For custom service account user should manually add appropriate Role to change Trial status.
	if s.Spec.EarlyStopping != nil && s.Spec.EarlyStopping.AlgorithmName != "" && suggestionConfigData.ServiceAccountName == "" {
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
	if s.Spec.EarlyStopping != nil && s.Spec.EarlyStopping.AlgorithmName != "" {
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

func (g *General) desiredContainers(s *suggestionsv1beta1.Suggestion,
	suggestionConfigData configv1beta1.SuggestionConfig,
	earlyStoppingConfigData configv1beta1.EarlyStoppingConfig) []corev1.Container {

	var (
		containers          []corev1.Container
		suggestionContainer corev1.Container
	)

	suggestionConfigData.Container.DeepCopyInto(&suggestionContainer)

	// Assign default values for suggestionContainer fields that are not set via
	// the suggestion config.

	if suggestionContainer.Name == "" {
		suggestionContainer.Name = consts.ContainerSuggestion
	}

	suggestionPort := corev1.ContainerPort{
		Name:          consts.DefaultSuggestionPortName,
		ContainerPort: consts.DefaultSuggestionPort,
	}
	suggestionContainer.Ports = append(suggestionContainer.Ports, suggestionPort)

	if viper.GetBool(consts.ConfigEnableGRPCProbeInSuggestion) && suggestionContainer.ReadinessProbe == nil {
		suggestionContainer.ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				GRPC: &corev1.GRPCAction{
					Port:    consts.DefaultSuggestionPort,
					Service: &consts.DefaultGRPCService,
				},
			},
			InitialDelaySeconds: defaultInitialDelaySeconds,
			PeriodSeconds:       defaultPeriodForReady,
		}
	}
	if viper.GetBool(consts.ConfigEnableGRPCProbeInSuggestion) && suggestionContainer.LivenessProbe == nil {
		suggestionContainer.LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				GRPC: &corev1.GRPCAction{
					Port:    consts.DefaultSuggestionPort,
					Service: &consts.DefaultGRPCService,
				},
			},
			// Ref https://srcco.de/posts/kubernetes-liveness-probes-are-dangerous.html
			InitialDelaySeconds: defaultInitialDelaySeconds,
			PeriodSeconds:       defaultPeriodForLive,
			FailureThreshold:    defaultFailureThreshold,
		}
	}

	if s.Spec.ResumePolicy == experimentsv1beta1.FromVolume && !containsVolumeMountWithName(suggestionContainer.VolumeMounts, consts.ContainerSuggestionVolumeName) {
		suggestionVolume := corev1.VolumeMount{
			Name:      consts.ContainerSuggestionVolumeName,
			MountPath: suggestionConfigData.VolumeMountPath,
		}
		suggestionContainer.VolumeMounts = append(suggestionContainer.VolumeMounts, suggestionVolume)
	}

	containers = append(containers, suggestionContainer)

	if s.Spec.EarlyStopping != nil && s.Spec.EarlyStopping.AlgorithmName != "" {
		earlyStoppingContainer := corev1.Container{
			Name:            consts.ContainerEarlyStopping,
			Image:           earlyStoppingConfigData.Image,
			ImagePullPolicy: earlyStoppingConfigData.ImagePullPolicy,
			Ports: []corev1.ContainerPort{
				{
					Name:          consts.DefaultEarlyStoppingPortName,
					ContainerPort: consts.DefaultEarlyStoppingPort,
				},
			},
			Resources: earlyStoppingConfigData.Resource,
		}

		containers = append(containers, earlyStoppingContainer)
	}
	return containers
}

func containsVolumeMountWithName(volumeMounts []corev1.VolumeMount, name string) bool {
	for i := range volumeMounts {
		if volumeMounts[i].Name == name {
			return true
		}
	}

	return false
}

func containsContainerPortWithName(ports []corev1.ContainerPort, name string) bool {
	for i := range ports {
		if ports[i].Name == name {
			return true
		}
	}

	return false
}

func containsContainerPort(ports []corev1.ContainerPort, port int32) bool {
	for i := range ports {
		if ports[i].ContainerPort == port {
			return true
		}
	}

	return false
}

// DesiredVolume returns desired PVC and PV for Suggestion.
// If PV doesn't exist in Katib config return nil for PV.
func (g *General) DesiredVolume(s *suggestionsv1beta1.Suggestion) (*corev1.PersistentVolumeClaim, *corev1.PersistentVolume, error) {

	suggestionConfigData, err := katibconfig.GetSuggestionConfigData(s.Spec.Algorithm.AlgorithmName, g.Client)
	if err != nil {
		return nil, nil, err
	}

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
	// Create PV if Katib config contains it.
	if !equality.Semantic.DeepEqual(suggestionConfigData.PersistentVolumeSpec, corev1.PersistentVolumeSpec{}) {

		persistentVolumeName := util.GetSuggestionPersistentVolumeName(s)

		pv = &corev1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:   persistentVolumeName,
				Labels: suggestionConfigData.PersistentVolumeLabels,
			},
			Spec: suggestionConfigData.PersistentVolumeSpec,
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
