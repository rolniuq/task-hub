package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"message": "hello"}

	writeJSON(rec, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var result map[string]string
	err := json.NewDecoder(rec.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "hello", result["message"])
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	writeError(rec, http.StatusBadRequest, "something went wrong")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var result ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "something went wrong", result.Error)
}

func TestErrorResponse(t *testing.T) {
	errResp := ErrorResponse{Error: "test error"}

	data, err := json.Marshal(errResp)
	assert.NoError(t, err)

	var decoded ErrorResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "test error", decoded.Error)
}

func TestNewAuthHandler(t *testing.T) {
	handler := NewAuthHandler(nil)
	assert.NotNil(t, handler)
}

func TestAuthHandler_Register_MethodNotAllowed(t *testing.T) {
	handler := NewAuthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/register", nil)
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	handler := NewAuthHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	handler := NewAuthHandler(nil)

	body := `{"email": "test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_Login_MethodNotAllowed(t *testing.T) {
	handler := NewAuthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	handler := NewAuthHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	handler := NewAuthHandler(nil)

	body := `{"email": "test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_RefreshToken_MethodNotAllowed(t *testing.T) {
	handler := NewAuthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/refresh", nil)
	rec := httptest.NewRecorder()

	handler.RefreshToken(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestAuthHandler_RefreshToken_InvalidBody(t *testing.T) {
	handler := NewAuthHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.RefreshToken(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_RefreshToken_MissingToken(t *testing.T) {
	handler := NewAuthHandler(nil)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.RefreshToken(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
