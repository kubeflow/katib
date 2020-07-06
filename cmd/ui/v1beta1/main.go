package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	common_v1beta1 "github.com/kubeflow/katib/pkg/common/v1beta1"
	ui "github.com/kubeflow/katib/pkg/ui/v1beta1"
)

var (
	port, host, buildDir, dbManagerAddr *string
)

func init() {
	port = flag.String("port", "80", "The port to listen to for incoming HTTP connections")
	host = flag.String("host", "0.0.0.0", "The host to listen to for incoming HTTP connections")
	buildDir = flag.String("build-dir", "/app/build", "The dir of frontend")
	dbManagerAddr = flag.String("db-manager-address", common_v1beta1.GetDBManagerAddr(), "The address of Katib DB manager")
}

func main() {
	flag.Parse()
	kuh := ui.NewKatibUIHandler(*dbManagerAddr)

	log.Printf("Serving the frontend dir %s", *buildDir)
	frontend := http.FileServer(http.Dir(*buildDir))
	http.Handle("/katib/", http.StripPrefix("/katib/", frontend))

	http.HandleFunc("/katib/fetch_experiments/", kuh.FetchAllExperiments)

	http.HandleFunc("/katib/submit_yaml/", kuh.SubmitYamlJob)
	http.HandleFunc("/katib/submit_hp_job/", kuh.SubmitParamsJob)
	http.HandleFunc("/katib/submit_nas_job/", kuh.SubmitParamsJob)

	http.HandleFunc("/katib/delete_experiment/", kuh.DeleteExperiment)

	http.HandleFunc("/katib/fetch_experiment/", kuh.FetchExperiment)
	http.HandleFunc("/katib/fetch_suggestion/", kuh.FetchSuggestion)

	http.HandleFunc("/katib/fetch_hp_job_info/", kuh.FetchHPJobInfo)
	http.HandleFunc("/katib/fetch_hp_job_trial_info/", kuh.FetchHPJobTrialInfo)
	http.HandleFunc("/katib/fetch_nas_job_info/", kuh.FetchNASJobInfo)

	http.HandleFunc("/katib/fetch_trial_templates/", kuh.FetchTrialTemplates)
	http.HandleFunc("/katib/add_template/", kuh.AddTemplate)
	http.HandleFunc("/katib/edit_template/", kuh.EditTemplate)
	http.HandleFunc("/katib/delete_template/", kuh.DeleteTemplate)
	http.HandleFunc("/katib/fetch_namespaces", kuh.FetchNamespaces)

	log.Printf("Serving at %s:%s", *host, *port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", *host, *port), nil); err != nil {
		panic(err)
	}
}
