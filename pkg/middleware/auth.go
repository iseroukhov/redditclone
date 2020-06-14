package middleware

import (
	"context"
	"net/http"
	"redditclone/pkg/response"
	"redditclone/pkg/user"
	"strings"
)

func Auth(repo *user.Repository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		token := strings.Replace(bearerToken, "Bearer ", "", 1)
		if token != "" {
			usr, err := repo.GetByToken(token)
			if err != nil {
				response.Error(w, err, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", usr)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}
