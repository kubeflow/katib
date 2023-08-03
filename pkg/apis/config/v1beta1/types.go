/*
Copyright 2023 The Kubeflow Authors.

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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// KatibConfig is the katib-config.yaml structure in Katib config.
type KatibConfig struct {
	metav1.TypeMeta `json:",inline"`

	RuntimeConfig RuntimeConfig `json:"runtime,omitempty"`
	InitConfig    InitConfig    `json:"init,omitempty"`
}

// RuntimeConfig is the runtime structure in Katib config.
type RuntimeConfig struct {
	SuggestionConfigs       []SuggestionConfig       `json:"suggestions,omitempty"`
	EarlyStoppingConfigs    []EarlyStoppingConfig    `json:"earlyStoppings,omitempty"`
	MetricsCollectorConfigs []MetricsCollectorConfig `json:"metricsCollectors,omitempty"`
}

// InitConfig is the YAML init structure in Katib config.
type InitConfig struct {
	ControllerConfig    ControllerConfig    `json:"controller,omitempty"`
	CertGeneratorConfig CertGeneratorConfig `json:"certGenerator,omitempty"`

	// TODO: Adding a config for the following components would be nice.
	// - Katib DB
	// - Katib DB Manager
	// - Katib UI
}

// ControllerConfig is the controller structure in Katib config.
type ControllerConfig struct {
	// ExperimentSuggestionName is the implementation of suggestion interface in experiment controller.
	// Defaults to 'default'.
	ExperimentSuggestionName string `json:"experimentSuggestionName,omitempty"`
	// MetricsAddr is the address the metric endpoint binds to.
	// Defaults to ':8080'.
	MetricsAddr string `json:"metricsAddr,omitempty"`
	// HealthzAddr is the address the healthz endpoint binds to.
	// Defaults to ':18080'.
	HealthzAddr string `json:"healthzAddr,omitempty"`
	// InjectSecurityContext indicates whether inject the securityContext of container[0] in the sidecar.
	// Defaults to 'false'.
	InjectSecurityContext bool `json:"injectSecurityContext,omitempty"`
	// EnableGRPCProbeInSuggestion indicates whether enable grpc probe in suggestions.
	// Defaults to 'true'.
	EnableGRPCProbeInSuggestion *bool `json:"enableGRPCProbeInSuggestion,omitempty"`
	// TrialResources is the list of resources that can be used as trial template,
	// in the form: Kind.Version.Group (e.g. TFJob.v1.kubeflow.org)
	// Defaults to 'Job.v1.batch'
	TrialResources []string `json:"trialResources,omitempty"`
	// WebhookPort is the port number to be used for admission webhook server.
	// Defaults to '8443'.
	WebhookPort *int `json:"webhookPort,omitempty"`
	// EnableLeaderElection indicates whether enable leader election for katib-controller.
	// Enabling this will ensure there is only one active katib-controller.
	// Defaults to 'false'.
	EnableLeaderElection bool `json:"enableLeaderElection,omitempty"`
	// LeaderElectionID is the ID for leader election.
	// Defaults to '3fbc96e9.katib.kubeflow.org'.
	LeaderElectionID string `json:"leaderElectionID,omitempty"`
}

// CertGeneratorConfig is the certGenerator structure in Katib config.
type CertGeneratorConfig struct {
	// Enable indicates the internal cert-generator is enabled.
	// Defaults to 'false'.
	Enable bool `json:"enable,omitempty"`
	// WebhookServiceName indicates which service is used for the admission webhook.
	// If it is set, the cert-generator forcefully is enabled even if the '.init.certGenerator.enable' is false.
	// Defaults to 'katib-controller'.
	WebhookServiceName string `json:"webhookServiceName,omitempty"`
	// WebhookSecretName indicates which secrets is used to save the certs for the admission webhook.
	// If it is set, the cert-generator forcefully is enabled even if the '.init.certGenerator.enable' is false.
	// Defaults to 'katib-webhook-cert'.
	WebhookSecretName string `json:"webhookSecretName,omitempty"`
}

// SuggestionConfig is the suggestion structure in Katib config.
type SuggestionConfig struct {
	AlgorithmName             string `json:"algorithmName"`
	corev1.Container          `json:",inline"`
	ServiceAccountName        string                           `json:"serviceAccountName,omitempty"`
	VolumeMountPath           string                           `json:"volumeMountPath,omitempty"`
	PersistentVolumeClaimSpec corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimSpec,omitempty"`
	PersistentVolumeSpec      corev1.PersistentVolumeSpec      `json:"persistentVolumeSpec,omitempty"`
	PersistentVolumeLabels    map[string]string                `json:"persistentVolumeLabels,omitempty"`
}

// EarlyStoppingConfig is the early stopping structure in Katib config.
type EarlyStoppingConfig struct {
	AlgorithmName   string                      `json:"algorithmName"`
	Image           string                      `json:"image"`
	ImagePullPolicy corev1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	Resource        corev1.ResourceRequirements `json:"resources,omitempty"`
}

// MetricsCollectorConfig is the metrics collector structure in Katib config.
type MetricsCollectorConfig struct {
	CollectorKind    string                      `json:"kind"`
	Image            string                      `json:"image"`
	ImagePullPolicy  corev1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	Resource         corev1.ResourceRequirements `json:"resources,omitempty"`
	WaitAllProcesses *bool                       `json:"waitAllProcesses,omitempty"`
}
