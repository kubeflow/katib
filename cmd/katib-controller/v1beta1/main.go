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

/*
 Katib-controller is a controller (operator) for Experiments and Trials
*/
package main

import (
	"flag"
	"os"

	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	apis "github.com/kubeflow/katib/pkg/apis/controller"
	controller "github.com/kubeflow/katib/pkg/controller.v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	trialutil "github.com/kubeflow/katib/pkg/controller.v1beta1/trial/util"
	webhook "github.com/kubeflow/katib/pkg/webhook/v1beta1"
)

func main() {
	logf.SetLogger(zap.New())
	log := logf.Log.WithName("entrypoint")

	var experimentSuggestionName string
	var metricsAddr string
	var webhookPort int
	var injectSecurityContext bool
	var enableGRPCProbeInSuggestion bool
	var trialResources trialutil.GvkListFlag
	var enableLeaderElection bool
	var leaderElectionID string

	flag.StringVar(&experimentSuggestionName, "experiment-suggestion-name",
		"default", "The implementation of suggestion interface in experiment controller (default)")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&injectSecurityContext, "webhook-inject-securitycontext", false, "Inject the securityContext of container[0] in the sidecar")
	flag.BoolVar(&enableGRPCProbeInSuggestion, "enable-grpc-probe-in-suggestion", true, "enable grpc probe in suggestions")
	flag.Var(&trialResources, "trial-resources", "The list of resources that can be used as trial template, in the form: Kind.version.group (e.g. TFJob.v1.kubeflow.org)")
	flag.IntVar(&webhookPort, "webhook-port", 8443, "The port number to be used for admission webhook server.")
	// For leader election
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false, "Enable leader election for katib-controller. Enabling this will ensure there is only one active katib-controller.")
	flag.StringVar(&leaderElectionID, "leader-election-id", "3fbc96e9.katib.kubeflow.org", "The ID for leader election.")

	// TODO (andreyvelich): Currently it is not possible to set different webhook service name.
	// flag.StringVar(&serviceName, "webhook-service-name", "katib-controller", "The service name which will be used in webhook")
	// TODO (andreyvelich): Currently is is not possible to store webhook cert in the local file system.
	// flag.BoolVar(&certLocalFS, "cert-localfs", false, "Store the webhook cert in local file system")

	flag.Parse()

	// Set the config in viper.
	viper.Set(consts.ConfigExperimentSuggestionName, experimentSuggestionName)
	viper.Set(consts.ConfigInjectSecurityContext, injectSecurityContext)
	viper.Set(consts.ConfigEnableGRPCProbeInSuggestion, enableGRPCProbeInSuggestion)
	viper.Set(consts.ConfigTrialResources, trialResources)

	log.Info("Config:",
		consts.ConfigExperimentSuggestionName,
		viper.GetString(consts.ConfigExperimentSuggestionName),
		"webhook-port",
		webhookPort,
		"metrics-addr",
		metricsAddr,
		consts.ConfigInjectSecurityContext,
		viper.GetBool(consts.ConfigInjectSecurityContext),
		consts.ConfigEnableGRPCProbeInSuggestion,
		viper.GetBool(consts.ConfigEnableGRPCProbeInSuggestion),
		"trial-resources",
		viper.Get(consts.ConfigTrialResources),
	)

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Fail to get the config")
		os.Exit(1)
	}

	// Create a new katib controller to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   leaderElectionID,
	})
	if err != nil {
		log.Error(err, "Failed to create the manager")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "Unable to add APIs to scheme")
		os.Exit(1)
	}

	// Setup all Controllers
	log.Info("Setting up controller.")
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "Unable to register controllers to the manager")
		os.Exit(1)
	}

	log.Info("Setting up webhooks.")
	if err := webhook.AddToManager(mgr, webhookPort); err != nil {
		log.Error(err, "Unable to register webhooks to the manager")
		os.Exit(1)
	}

	// Start the Cmd
	log.Info("Starting the Cmd.")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "Unable to run the manager")
		os.Exit(1)
	}
}
