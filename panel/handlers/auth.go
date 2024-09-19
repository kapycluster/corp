package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kapycluster/corpy/panel/views"
	authview "github.com/kapycluster/corpy/panel/views/auth"
	"github.com/kapycluster/corpy/panel/views/dashboard"
)

func (h Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	_, err := h.auth.GetSessionUser(r)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	authview.Login().Render(r.Context(), w)
}

func (h Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	if _, err := h.auth.CompleteUserAuth(w, r); err == nil {
		// user logged in
		dashboard.ControlPlanes().Render(r.Context(), w)
	} else {
		h.auth.BeginAuthHandler(w, r)
	}
}

func (h Handler) HandleProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	u, err := h.auth.CompleteUserAuth(w, r)
	if err != nil {
		views.Error(err.Error()).Render(r.Context(), w)
		return
	}

	if err := h.auth.StoreUserSession(w, r, u); err != nil {
		views.Error("failed to store user session").Render(r.Context(), w)
		return
	}

	// do something with the user (e.g. register or sign in)
	h.log.Info("user", "name", u.Name, "email", u.Email, "uid", u.UserID, "avatar", u.AvatarURL)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
