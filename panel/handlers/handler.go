package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/markbates/goth"
	"kapycluster.com/corp/panel/auth"
	"kapycluster.com/corp/panel/config"
	"kapycluster.com/corp/panel/views"
)

type Handler struct {
	kc   KubeClient
	db   DBStore
	log  *slog.Logger
	c    *config.Config
	auth *auth.Auth
	dns  DNSClient
}

func (h Handler) MustGetUser(w http.ResponseWriter, r *http.Request) goth.User {
	u, err := h.auth.GetSessionUser(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return goth.User{}
	}
	return u
}

func (h Handler) RenderOrRedirect(w http.ResponseWriter, r *http.Request, c templ.Component, path string) {
	if r.Header.Get("hx-request") != "" {
		c.Render(r.Context(), w)
	} else {
		http.Redirect(w, r, path, http.StatusSeeOther)
	}
}

func (h Handler) Error(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	h.log.Error(msg, "error", err)
	views.Error(msg).Render(ctx, w)
}
