package web

import (
	"net/http"
)

type authenticationMiddleware struct {
	adminKeys map[string]struct{}
}

func NewAuthenticationMiddleware(adminKey string) *authenticationMiddleware {
	return &authenticationMiddleware{
		adminKeys: map[string]struct{}{
			adminKey: {},
		},
	}
}

func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Admin-Token")

		if _, found := amw.adminKeys[token]; found {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
