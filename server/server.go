package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ListenAndServe() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := readTemplates()
		if err != nil {
			http.Error(w, fmt.Sprintf("read templates: %s", err), http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("execute template: %s", err), http.StatusInternalServerError)
			return
		}
	})

	// serve static files
	filesDir := http.Dir("server/static")
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := "/static/"
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx)))
	})

	return http.ListenAndServe(":8080", r)
}

func readTemplates() (*template.Template, error) {
	tmpl, err := template.ParseGlob("server/templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	return tmpl, nil
}
