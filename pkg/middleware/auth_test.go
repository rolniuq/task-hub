package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextKey(t *testing.T) {
	assert.Equal(t, contextKey("user_id"), UserIDKey)
	assert.Equal(t, contextKey("email"), EmailKey)
}

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "with user id",
			ctx:      context.WithValue(context.Background(), UserIDKey, "user-123"),
			expected: "user-123",
		},
		{
			name:     "without user id",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "with wrong type",
			ctx:      context.WithValue(context.Background(), UserIDKey, 123),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUserIDFromContext(tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEmailFromContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "with email",
			ctx:      context.WithValue(context.Background(), EmailKey, "test@example.com"),
			expected: "test@example.com",
		},
		{
			name:     "without email",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "with wrong type",
			ctx:      context.WithValue(context.Background(), EmailKey, 123),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEmailFromContext(tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthMiddleware_NoAuthHeader(t *testing.T) {
	middleware := &AuthMiddleware{}

	handler := middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Authorization header required")
}

func TestAuthMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	middleware := &AuthMiddleware{}

	tests := []struct {
		name       string
		authHeader string
	}{
		{
			name:       "no bearer prefix",
			authHeader: "token-only",
		},
		{
			name:       "wrong prefix",
			authHeader: "Basic token123",
		},
		{
			name:       "too many parts",
			authHeader: "Bearer token extra",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})
	}
}

func TestContextValues(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, "test-user-id")
	ctx = context.WithValue(ctx, EmailKey, "test@example.com")

	userID := GetUserIDFromContext(ctx)
	email := GetEmailFromContext(ctx)

	assert.Equal(t, "test-user-id", userID)
	assert.Equal(t, "test@example.com", email)
}

func TestNewAuthMiddleware(t *testing.T) {
	middleware := NewAuthMiddleware(nil)
	assert.NotNil(t, middleware)
}
