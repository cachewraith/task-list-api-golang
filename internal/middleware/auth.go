package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/example/todo-api/internal/service"
	"github.com/example/todo-api/pkg/response"
)

// contextKey is a type for context keys
type contextKey string

const (
	// UserIDContextKey is the key for user ID in context
	UserIDContextKey contextKey = "user_id"
	// UserEmailContextKey is the key for user email in context
	UserEmailContextKey contextKey = "user_email"
)

// Auth creates an authentication middleware
func Auth(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.NewError("UNAUTHORIZED", "Authorization header required").Write(w, http.StatusUnauthorized)
				return
			}

			// Extract Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.NewError("UNAUTHORIZED", "Invalid authorization format. Use: Bearer <token>").Write(w, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				response.NewError("UNAUTHORIZED", "Invalid or expired token").Write(w, http.StatusUnauthorized)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailContextKey, claims.Email)

			// Continue with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	if ok {
		return userID, true
	}
	// Try uuid.UUID type
	userUUID, ok := ctx.Value(UserIDContextKey).(fmt.Stringer)
	if ok {
		return userUUID.String(), true
	}
	return "", false
}

// GetUserEmailFromContext extracts user email from context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailContextKey).(string)
	return email, ok
}
