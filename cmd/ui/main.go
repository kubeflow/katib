package main

import (
	"github.com/kubeflow/katib/pkg/ui"
	"github.com/pressly/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	kuh := ui.NewKatibUIHandler()
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("/static"))))
	r.Route("/katib", func(r chi.Router) {
		r.Get("/", kuh.Index)
		r.Route("/{studyid}", func(r chi.Router) {
			r.Get("/", kuh.Study)
			r.Get("/csv", kuh.StudyInfoCsv)
			r.Route("/TrialID/{trialid}", func(r chi.Router) {
				r.Get("/", kuh.Trial)
			})
			r.Route("/WorkerID/{workerid}", func(r chi.Router) {
				r.Get("/", kuh.Worker)
				r.Get("/csv", kuh.WorkerInfoCsv)
			})
		})
	})
	http.ListenAndServe(":80", r)
}
