package middleware

import (
	"app/finance/configs"
	"app/finance/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	ContextIDKey key = "ContextIDKey"
)

func ISAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)
		if !isValid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextIDKey, data.UserID)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
