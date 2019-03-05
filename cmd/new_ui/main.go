package main

import (
	"net/http"

	"github.com/kubeflow/katib/pkg/ui"
)

func main() {
	kuh := ui.NewKatibUIHandler()

	http.HandleFunc("/katib/fetch_hp_jobs/", kuh.FetchHPJobs)
	http.HandleFunc("/katib/fetch_nas_jobs/", kuh.FetchNASJobs)
	http.HandleFunc("/katib/submit_yaml/", kuh.SubmitYamlJob)
	http.HandleFunc("/katib/fetch_worker_templates/", kuh.FetchWorkerTemplates)
	http.HandleFunc("/katib/fetch_collector_templates/", kuh.FetchCollectorTemplates)
	http.HandleFunc("/katib/update_template/", kuh.AddEditTemplate)
	http.HandleFunc("/katib/delete_template/", kuh.DeleteTemplate)

	http.ListenAndServe(":9303", nil)
}
