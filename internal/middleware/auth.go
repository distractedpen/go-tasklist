package middleware

import (
	"go-tasklist/internal/auth"
	"net/http"
	"strings"
)

type EnsureAuth struct {
	handler http.Handler
}

func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if len(authHeader) == 0 {
				http.Error(w, "Not Authorized.", http.StatusUnauthorized)
				return
			}
			token, found := strings.CutPrefix(authHeader, "Bearer ")
			if !found {
				http.Error(w, "Not Authorized.", http.StatusUnauthorized)
			}

			valid, err := auth.VerifyToken(token)
			if nil != err || !valid {
				http.Error(w, "Not Authorized.", http.StatusUnauthorized)
			}

			next.ServeHTTP(w, r)
		})
}


