package handlers

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/kapycluster/corpy/panel/auth"
	"github.com/kapycluster/corpy/panel/config"
	"github.com/markbates/goth"
)

type Handler struct {
	kc   KubeClient
	db   DBStore
	log  *slog.Logger
	c    *config.Config
	auth *auth.Auth
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
