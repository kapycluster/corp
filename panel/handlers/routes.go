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
	r.Get("/", dashboard.ShowDashboard)
	r.Route("/static", func(r chi.Router) {
		r.Handle("/style.css", http.FileServerFS(views.Style()))
		r.Handle(
			"/js/preline.js",
			http.StripPrefix("/static/js/", http.FileServerFS(views.Preline())),
		)
	})

	return r
}
