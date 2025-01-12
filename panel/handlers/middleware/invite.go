package middleware

import (
	"context"
	"net/http"

	"kapycluster.com/corp/panel/store"
)

func getInviteFromQuery(r *http.Request) (string, bool) {
	inviteCode := r.URL.Query().Get("invite")
	return inviteCode, inviteCode != ""
}

func getInviteFromCookie(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("kapy-invite")
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

func setInviteContext(r *http.Request, inviteCode string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "invite", inviteCode)
	return r.WithContext(ctx)
}

func writeUnauthorizedResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(message))
}

// TODO: fix dependency on store.DB
func ValidateInvite(db *store.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get invite from query params first
			inviteCode, ok := getInviteFromQuery(r)
			if ok {
				r = setInviteContext(r, inviteCode)
			} else {
				// Try to get invite from cookie
				inviteCode, ok = getInviteFromCookie(r)
				if !ok {
					writeUnauthorizedResponse(w, "you need an invite")
					return
				}
				// Set invite in context and continue; we don't need to validate it
				r = setInviteContext(r, inviteCode)
				next.ServeHTTP(w, r)
				return
			}

			// Validate invite code
			storedInvite, _ := db.GetInvite(r.Context(), inviteCode)
			if storedInvite == nil || storedInvite.Used {
				writeUnauthorizedResponse(w, "invalid invite code")
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "kapy-invite",
				Value: inviteCode,
			})
			next.ServeHTTP(w, r)
		})
	}
}
