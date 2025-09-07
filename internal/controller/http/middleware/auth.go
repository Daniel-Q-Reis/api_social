package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type contextKey string

const UserContextKey contextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In a real implementation, you would validate the JWT token here
		// For now, we'll just check for an Authorization header and extract a user ID
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "bearer token required", http.StatusUnauthorized)
			return
		}

		// In a real implementation, you would parse and validate the JWT token
		// For now, we'll just use a placeholder user ID
		// In a real app, this would come from the validated token
		userID, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
		if err != nil {
			http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
