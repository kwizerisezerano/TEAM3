package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware is a placeholder for JWT authentication
// TODO: Implement proper JWT authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// TODO: Validate JWT token here
		// For now, we'll just check if token is not empty
		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// TODO: Extract user info from token and add to context
		// ctx := context.WithValue(r.Context(), "userID", userID)
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// RoleMiddleware checks if user has required role
// TODO: Implement role-based access control
func RoleMiddleware(allowedRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Get user role from context and check against allowedRoles
			// For now, just pass through
			next.ServeHTTP(w, r)
		})
	}
}
