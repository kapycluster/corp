package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kapycluster/corpy/panel/handlers/middleware"
	"github.com/kapycluster/corpy/panel/views"
)

func Setup(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(ctx))

	dashboard := Dashboard{}
	r.Get("/dashboard", dashboard.ShowDashboard)
	r.Get("/dashboard/*", dashboard.ShowDashboard)
	r.Get("/dashboard/controlplanes/create", dashboard.CreateControlPlane)

	r.Route("/static", func(r chi.Router) {
		r.Handle("/style.css", http.FileServerFS(views.Style()))
		prefix := "/static/js/"
		r.Handle(
			"/js/htmx.min.js",
			http.StripPrefix(prefix, http.FileServerFS(views.HTMX())),
		)
	})

	return r
}
