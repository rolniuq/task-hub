package middleware

import (
	"context"
	"net/http"
	"strings"
	"taskhub/internal/app"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

type AuthMiddleware struct {
	authService *app.AuthService
}

func NewAuthMiddleware(authService *app.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := m.authService.ValidateAccessToken(tokenString)
		if err != nil {
			if err == app.ErrTokenExpired {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

func GetEmailFromContext(ctx context.Context) string {
	if email, ok := ctx.Value(EmailKey).(string); ok {
		return email
	}
	return ""
}
