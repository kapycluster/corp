package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"kapycluster.com/corp/log"
	"kapycluster.com/corp/panel/auth"
	"kapycluster.com/corp/panel/config"
	"kapycluster.com/corp/panel/dns"
	"kapycluster.com/corp/panel/handlers/middleware"
	"kapycluster.com/corp/panel/kube"
	"kapycluster.com/corp/panel/store"
	"kapycluster.com/corp/panel/views"
)

func Setup(ctx context.Context, cfg *config.Config) (*chi.Mux, error) {
	r := chi.NewRouter()

	kubeClient, err := kube.NewKube(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube client: %w", err)
	}

	sessionStore := auth.NewCookieStore(
		auth.SessionOptions{
			CookiesKey: cfg.Session.Secret,
			HttpOnly:   cfg.Session.HttpOnly,
			Secure:     cfg.Session.Secure,
			MaxAge:     cfg.Session.MaxAge,
		},
	)
	auth := auth.NewAuth(cfg, sessionStore)

	dbStore, err := store.New(cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create db store: %w", err)
	}

	err = dbStore.Setup(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup db store: %w", err)
	}

	cloudflare, err := dns.NewCloudflare(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudflare client: %w", err)
	}

	handler := Handler{
		kc:   kubeClient,
		db:   dbStore,
		log:  log.FromContext(ctx),
		c:    cfg,
		auth: auth,
		dns:  cloudflare,
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
			r.Get("/{id}/kubeconfig", handler.DownloadKubeconfigStub)
			r.Get("/{id}/kubeconfig/download", handler.DownloadKubeconfig)
			// r.Get("/controlplane/{id}/more", handler.)
		})

		r.Route("/auth/", func(r chi.Router) {
			r.Get("/{provider}", handler.HandleProviderLogin)
			r.Get("/{provider}/callback", handler.HandleProviderCallback)
			r.Get("/{provider}/logout", handler.HandleLogout)
			r.With(middleware.ValidateInvite(dbStore)).Get("/login", handler.ShowLogin)
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
