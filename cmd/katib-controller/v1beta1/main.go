/*
Copyright 2018 The Kubeflow Authors
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
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	apis "github.com/kubeflow/katib/pkg/apis/controller"
	controller "github.com/kubeflow/katib/pkg/controller.v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	trialutil "github.com/kubeflow/katib/pkg/controller.v1beta1/trial/util"
	webhook "github.com/kubeflow/katib/pkg/webhook/v1beta1"
)

func main() {
	logf.SetLogger(logf.ZapLogger(false))
	log := logf.Log.WithName("entrypoint")

	var experimentSuggestionName string
	var metricsAddr string
	var webhookPort int
	var certLocalFS bool
	var injectSecurityContext bool
	var serviceName string
	var enableGRPCProbeInSuggestion bool
	var trialResources trialutil.GvkListFlag

	flag.StringVar(&experimentSuggestionName, "experiment-suggestion-name",
		"default", "The implementation of suggestion interface in experiment controller (default)")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.IntVar(&webhookPort, "webhook-port", 8443, "The port number to be used for admission webhook server.")
	flag.BoolVar(&certLocalFS, "cert-localfs", false, "Store the webhook cert in local file system")
	flag.BoolVar(&injectSecurityContext, "webhook-inject-securitycontext", false, "Inject the securityContext of container[0] in the sidecar")
	flag.StringVar(&serviceName, "webhook-service-name", "katib-controller", "The service name which will be used in webhook")
	flag.BoolVar(&enableGRPCProbeInSuggestion, "enable-grpc-probe-in-suggestion", true, "enable grpc probe in suggestions")
	flag.Var(&trialResources, "trial-resources", "The list of resources that can be used as trial template, in the form: Kind.version.group (e.g. TFJob.v1.kubeflow.org)")

	flag.Parse()

	// Set the config in viper.
	viper.Set(consts.ConfigExperimentSuggestionName, experimentSuggestionName)
	viper.Set(consts.ConfigCertLocalFS, certLocalFS)
	viper.Set(consts.ConfigInjectSecurityContext, injectSecurityContext)
	viper.Set(consts.ConfigEnableGRPCProbeInSuggestion, enableGRPCProbeInSuggestion)
	viper.Set(consts.ConfigTrialResources, trialResources)

	log.Info("Config:",
		consts.ConfigExperimentSuggestionName,
		viper.GetString(consts.ConfigExperimentSuggestionName),
		consts.ConfigCertLocalFS,
		viper.GetBool(consts.ConfigCertLocalFS),
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
	})
	if err != nil {
		log.Error(err, "unable add APIs to scheme")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "Fail to create the manager")
		os.Exit(1)
	}

	// Setup all Controllers
	log.Info("Setting up controller")
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "unable to register controllers to the manager")
		os.Exit(1)
	}

	log.Info("Setting up webhooks")
	if err := webhook.AddToManager(mgr, int32(webhookPort), serviceName); err != nil {
		log.Error(err, "unable to register webhooks to the manager")
		os.Exit(1)
	}

	// Start the Cmd
	log.Info("Starting the Cmd.")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to run the manager")
		os.Exit(1)
	}
}
