package handlers

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/kapycluster/corpy/panel/handlers/middleware"
)

func Setup(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(ctx))

	dashboard := Dashboard{}
	r.Get("/", dashboard.ShowDashboard)

	return r
}
