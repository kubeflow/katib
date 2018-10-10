package main

import (
	"github.com/pressly/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	http.ListenAndServe(":80", r)
}
