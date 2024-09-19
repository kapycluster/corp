package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/kapycluster/corpy/panel/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

const SessionName = "kapy-panel"

type Auth struct {
}

func NewAuth(c *config.Config, store sessions.Store) *Auth {
	gothic.Store = store

	goth.UseProviders(github.New(
		c.OAuth.GitHub.Key,
		c.OAuth.GitHub.Secret,
		buildCallbackURL("github", c),
	))
	return &Auth{}
}

func (a *Auth) BeginAuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (a *Auth) CompleteUserAuth(w http.ResponseWriter, r *http.Request) (goth.User, error) {
	return gothic.CompleteUserAuth(w, r)
}

func (a *Auth) StoreUserSession(w http.ResponseWriter, r *http.Request, u goth.User) error {
	s, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	s.Values["user"] = u
	return s.Save(r, w)
}

func (a *Auth) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := a.GetSessionUser(r)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *Auth) GetSessionUser(r *http.Request) (goth.User, error) {
	s, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return goth.User{}, err
	}

	u, ok := s.Values["user"].(goth.User)
	if !ok {
		return goth.User{}, fmt.Errorf("user is not authenticated")
	}

	return u, nil
}

func buildCallbackURL(provider string, c *config.Config) string {
	return fmt.Sprintf("%s/auth/%s/callback", c.Server.BaseURL, provider)
}
