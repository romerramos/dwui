package auth

import (
	"net/http"

	dwuiHttp "github.com/dwui/cmd/http"
	"github.com/dwui/cmd/session"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if publicPaths(r) {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(session.SessionCookieName)
		if err != nil || !session.Validate(cookie.Value) {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func publicPaths(r *http.Request) bool {
	return r.URL.Path == "/signin" ||
		r.URL.Path == "/auth/signin" ||
		r.URL.Path == "/auth/signout" ||
		dwuiHttp.IsStaticFile(r.URL.Path)
}
