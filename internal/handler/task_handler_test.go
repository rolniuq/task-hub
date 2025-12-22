package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"taskhub/pkg/middleware"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTaskHandler(t *testing.T) {
	handler := NewTaskHandler(nil)
	assert.NotNil(t, handler)
}

func TestTaskHandler_Create_MethodNotAllowed(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestTaskHandler_Create_InvalidUser(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks", nil)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTaskHandler_Create_InvalidBody(t *testing.T) {
	handler := NewTaskHandler(nil)

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, uuid.New().String())
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString("invalid"))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTaskHandler_Create_MissingTitle(t *testing.T) {
	handler := NewTaskHandler(nil)

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, uuid.New().String())
	body := `{"description": "no title"}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(body))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTaskHandler_Get_MethodNotAllowed(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks/123", nil)
	rec := httptest.NewRecorder()

	handler.Get(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestTaskHandler_Get_InvalidUser(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/123", nil)
	rec := httptest.NewRecorder()

	handler.Get(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTaskHandler_Get_InvalidTaskID(t *testing.T) {
	handler := NewTaskHandler(nil)

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, uuid.New().String())
	req := httptest.NewRequest(http.MethodGet, "/api/tasks/invalid-uuid", nil)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.Get(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTaskHandler_Update_MethodNotAllowed(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks/123", nil)
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestTaskHandler_Update_InvalidUser(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodPut, "/api/tasks/123", nil)
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTaskHandler_Delete_MethodNotAllowed(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks/123", nil)
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestTaskHandler_Delete_InvalidUser(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/123", nil)
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTaskHandler_List_MethodNotAllowed(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks", nil)
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestTaskHandler_List_InvalidUser(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTaskHandler_Complete_MethodNotAllowed(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/123/complete", nil)
	rec := httptest.NewRecorder()

	handler.Complete(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestTaskHandler_Complete_InvalidUser(t *testing.T) {
	handler := NewTaskHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks/123/complete", nil)
	rec := httptest.NewRecorder()

	handler.Complete(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTaskHandler_Complete_InvalidTaskID(t *testing.T) {
	handler := NewTaskHandler(nil)

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, uuid.New().String())
	req := httptest.NewRequest(http.MethodPost, "/api/tasks/invalid-uuid/complete", nil)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.Complete(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
