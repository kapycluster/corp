package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kapycluster/corpy/panel/handlers/middleware"
	"github.com/kapycluster/corpy/panel/kube"
	"github.com/kapycluster/corpy/panel/store"
	"github.com/kapycluster/corpy/panel/views"
)

func Setup(ctx context.Context) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(ctx))

	kubeClient, err := kube.NewKube()
	if err != nil {
		return nil, fmt.Errorf("failed to create kube client: %w", err)
	}

	dbStore, err := store.NewDB()
	if err != nil {
		return nil, fmt.Errorf("failed to create db store: %w", err)
	}

	dashboard := Dashboard{
		kc: kubeClient,
		db: dbStore,
	}

	r.Route("/dashboard", func(r chi.Router) {
		r.Get("/", dashboard.ShowDashboard)
		r.Get("/*", dashboard.ShowDashboard)
		r.Get("/controlplanes/create", dashboard.ShowCreateControlPlaneForm)

		r.Post("/controlplanes/create", dashboard.HandleCreateControlPlaneForm)
	})

	r.Route("/static", func(r chi.Router) {
		r.Handle("/style.css", http.FileServerFS(views.Style()))
		prefix := "/static/js/"
		r.Handle(
			"/js/htmx.min.js",
			http.StripPrefix(prefix, http.FileServerFS(views.HTMX())),
		)
	})

	return r, nil
}
