package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func htmxRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/holder", func(w http.ResponseWriter, r *http.Request) {

	})
	return r
}
