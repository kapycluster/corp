package auth

import "github.com/gorilla/sessions"

type SessionOptions struct {
	CookiesKey string
	MaxAge     int
	HttpOnly   bool
	Secure     bool
}

func NewCookieStore(opts SessionOptions) *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(opts.CookiesKey))

	store.Options.MaxAge = opts.MaxAge
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure
	store.Options.Path = "/"

	return store
}
