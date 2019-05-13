package main

import (
	"net/http"

	ui "github.com/kubeflow/katib/pkg/ui/v1alpha2"
)

var (
	port = "80"
)

func main() {
	kuh := ui.NewKatibUIHandler()

	frontend := http.FileServer(http.Dir("/app/build/"))
	http.Handle("/katib/", http.StripPrefix("/katib/", frontend))

	http.HandleFunc("/katib/fetch_hp_jobs/", kuh.FetchHPJobs)
	http.HandleFunc("/katib/fetch_nas_jobs/", kuh.FetchNASJobs)
	http.HandleFunc("/katib/submit_yaml/", kuh.SubmitYamlJob)
	http.HandleFunc("/katib/submit_hp_job/", kuh.SubmitHPJob)
	http.HandleFunc("/katib/submit_nas_job/", kuh.SubmitNASJob)

	//TODO: Add it in Katib client
	http.HandleFunc("/katib/delete_job/", kuh.DeleteJob)

	http.HandleFunc("/katib/fetch_hp_job_info/", kuh.FetchHPJobInfo)
	http.HandleFunc("/katib/fetch_hp_job_trial_info/", kuh.FetchHPJobTrialInfo)
	http.HandleFunc("/katib/fetch_nas_job_info/", kuh.FetchNASJobInfo)

	http.HandleFunc("/katib/fetch_trial_templates/", kuh.FetchTrialTemplates)
	http.HandleFunc("/katib/fetch_collector_templates/", kuh.FetchMetricsCollectorTemplates)
	http.HandleFunc("/katib/update_template/", kuh.AddEditDeleteTemplate)

	http.ListenAndServe(":"+port, nil)
}
