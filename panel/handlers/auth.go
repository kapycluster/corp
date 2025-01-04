package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	authview "kapycluster.com/corp/panel/views/auth"
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
		http.Redirect(w, r, "/controlplanes", http.StatusSeeOther)
	} else {
		h.auth.BeginAuthHandler(w, r)
	}
}

func (h Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := h.auth.ClearUserSession(w, r)
	if err != nil {
		h.Error(r.Context(), w, "failed to clear user session", err)
		return
	}
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func (h Handler) HandleProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	u, err := h.auth.CompleteUserAuth(w, r)
	if err != nil {
		h.Error(r.Context(), w, "failed to complete user auth", err)
		return
	}

	if err := h.auth.StoreUserSession(w, r, u); err != nil {
		h.Error(r.Context(), w, "failed to store user session", err)
		return
	}

	// do something with the user (e.g. register or sign in)
	h.log.Info("user", "name", u.Name, "email", u.Email, "uid", u.UserID, "avatar", u.AvatarURL)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
