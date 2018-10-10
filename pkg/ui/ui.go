package ui

import (
	"net/http"
)

type KatibUIHandler struct {
}

func (k *KatibUIHandler) Index(w http.ResponseWriter, r *http.Request) {
}

func (k *KatibUIHandler) Study(w http.ResponseWriter, r *http.Request) {
	studyID := chi.URLParam(r, "studyid")
}
