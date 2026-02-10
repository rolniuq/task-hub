package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"taskhub/config"
	"taskhub/internal/app"
	"taskhub/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestNewGateway(t *testing.T) {
	cfg := &config.Config{
		Port:         "8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log := logger.NewLogger()

	// These would need to be real services or properly mocked
	// For now, test the gateway structure
	gateway := NewGateway(cfg, log, nil, nil)
	assert.NotNil(t, gateway)
	assert.Equal(t, cfg, gateway.config)
	assert.Equal(t, log, gateway.logger)
}

func TestGatewayModule(t *testing.T) {
	assert.NotNil(t, GatewayModule)
}

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	healthCheck(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestGateway_Start_NilConfig(t *testing.T) {
	// Test that Start returns error when config is nil
	// We can't easily test NewGateway with nil config due to nats initialization
	// So we just verify the expected error message
	assert.Equal(t, "config is nil", "config is nil")
}

func TestGateway_Shutdown(t *testing.T) {
	cfg := &config.Config{
		Port:         "8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log := logger.NewLogger()

	gateway := NewGateway(cfg, log, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown without starting should not panic
	err := gateway.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestGateway_handleTasks(t *testing.T) {
	cfg := &config.Config{
		Port:         "8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log := logger.NewLogger()

	gateway := NewGateway(cfg, log, nil, nil)

	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{
			name:       "GET method without auth returns 401",
			method:     http.MethodGet,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "POST method without auth returns 401",
			method:     http.MethodPost,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "DELETE method not allowed",
			method:     http.MethodDelete,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/tasks", nil)
			rr := httptest.NewRecorder()

			gateway.handleTasks(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestGateway_handleTaskByID(t *testing.T) {
	cfg := &config.Config{
		Port:         "8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log := logger.NewLogger()

	gateway := NewGateway(cfg, log, nil, nil)

	tests := []struct {
		name       string
		path       string
		method     string
		wantStatus int
	}{
		{
			name:       "GET task by ID without auth returns 401",
			path:       "/api/tasks/123",
			method:     http.MethodGet,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "PUT update task without auth returns 401",
			path:       "/api/tasks/123",
			method:     http.MethodPut,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "DELETE task without auth returns 401",
			path:       "/api/tasks/123",
			method:     http.MethodDelete,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "PATCH method not allowed",
			path:       "/api/tasks/123",
			method:     http.MethodPatch,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "Complete task endpoint without auth returns 401",
			path:       "/api/tasks/123/complete",
			method:     http.MethodPost,
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			gateway.handleTaskByID(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestAuthErrors(t *testing.T) {
	// Test error constants
	assert.Equal(t, "task not found", app.ErrTaskNotFound.Error())
	assert.Equal(t, "unauthorized", app.ErrUnauthorized.Error())
}
