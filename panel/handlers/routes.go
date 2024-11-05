package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kapycluster/corpy/log"
	"github.com/kapycluster/corpy/panel/auth"
	"github.com/kapycluster/corpy/panel/config"
	"github.com/kapycluster/corpy/panel/handlers/middleware"
	"github.com/kapycluster/corpy/panel/kube"
	"github.com/kapycluster/corpy/panel/store"
	"github.com/kapycluster/corpy/panel/views"
)

func Setup(ctx context.Context, config *config.Config) (*chi.Mux, error) {
	r := chi.NewRouter()

	kubeClient, err := kube.NewKube(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube client: %w", err)
	}

	sessionStore := auth.NewCookieStore(
		auth.SessionOptions{
			CookiesKey: config.Session.Secret,
			HttpOnly:   config.Session.HttpOnly,
			Secure:     config.Session.Secure,
			MaxAge:     config.Session.MaxAge,
		},
	)
	auth := auth.NewAuth(config, sessionStore)

	dbStore, err := store.NewDB()
	if err != nil {
		return nil, fmt.Errorf("failed to create db store: %w", err)
	}

	handler := Handler{
		kc:   kubeClient,
		db:   dbStore,
		log:  log.FromContext(ctx),
		c:    config,
		auth: auth,
	}

	// Show* functions render templ templates.
	// Handle* functions handle form submissions/affect the state of the application.
	// Fetch* functions fetch data from the database or other sources.
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.RequestLogger(ctx))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/controlplanes", http.StatusMovedPermanently)
		})

		r.Route("/controlplanes", func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Get("/", handler.ShowDashboard)
			r.Get("/create", handler.ShowCreateControlPlaneForm)
			r.Post("/create", handler.HandleCreateControlPlaneForm)
			// r.Get("/controlplane/{id}", handler.FetchControlPlaneInfo)
		})

		r.Route("/auth/", func(r chi.Router) {
			r.Get("/{provider}", handler.HandleProviderLogin)
			r.Get("/{provider}/callback", handler.HandleProviderCallback)
			r.Get("/{provider}/logout", handler.HandleLogout)
			r.Get("/login", handler.ShowLogin)
		})

	})

	r.Route("/static", func(r chi.Router) {
		r.Handle("/style.css", http.FileServerFS(views.Style()))
		prefix := "/static/js/"
		r.Handle(
			"/js/htmx.min.js",
			http.StripPrefix(prefix, http.FileServerFS(views.HTMX())),
		)
		r.Handle(
			"/js/cdn.min.js",
			http.StripPrefix(prefix, http.FileServerFS(views.Alpine())),
		)
	})

	return r, nil
}
