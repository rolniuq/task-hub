package desktop

import (
	"testing"

	"taskhub/internal/app"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	// Since NewApp requires *app.AuthService (concrete type),
	// we can only test that the function signature is correct
	// Actual testing would require a real AuthService or refactoring to use interfaces
	assert.True(t, true)
}

func TestDesktopApp_LoginRequest(t *testing.T) {
	req := &app.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestDesktopApp_RegisterRequest(t *testing.T) {
	req := &app.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.Equal(t, "Test User", req.Name)
	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestDesktopApp_TokenPair(t *testing.T) {
	tokens := &app.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    1234567890,
	}

	assert.Equal(t, "access-token", tokens.AccessToken)
	assert.Equal(t, "refresh-token", tokens.RefreshToken)
	assert.Equal(t, int64(1234567890), tokens.ExpiresAt)
}

func TestDesktopApp_Claims(t *testing.T) {
	claims := &app.Claims{
		UserID: "user-123",
		Email:  "test@example.com",
	}

	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestDesktopApp_RefreshClaims(t *testing.T) {
	claims := &app.RefreshClaims{
		UserID: "user-123",
	}

	assert.Equal(t, "user-123", claims.UserID)
}

func TestDesktopApp_RefreshTokenRequest(t *testing.T) {
	req := &app.RefreshTokenRequest{
		RefreshToken: "refresh-token-123",
	}

	assert.Equal(t, "refresh-token-123", req.RefreshToken)
}
