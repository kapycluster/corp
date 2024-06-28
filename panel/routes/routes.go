package routes

import (
	"context"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/decantor/panel/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Init sets up middlewares and routes for the application
func Init(app *config.App) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(AppMiddleware(app))

	// Routes
	r.Get("/login", loginPageHandler)

	return r
}

// AppMiddleware is a custom middleware that sets the app context on the request
func AppMiddleware(app *config.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "app", app)
			(next).ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// loginPageHandler renders the login page
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Get the app context
	a := r.Context().Value("app").(*config.App)

	// Render the login page
	tmpl, err := template.ParseFiles(filepath.Join(a.Config.Dirs.Templates, "index.html"))
	if err != nil {
		log.Error().Msgf("failed to parse template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Error().Msgf("failed to execute template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
