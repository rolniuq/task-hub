package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"taskhub/internal/app"
	"taskhub/internal/domains/task"
	"taskhub/pkg/middleware"
	"time"

	"github.com/google/uuid"
)

type TaskHandler struct {
	taskService *app.TaskService
}

func NewTaskHandler(taskService *app.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := middleware.GetUserIDFromContext(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	var req app.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	resp, err := h.taskService.CreateTask(r.Context(), &req, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := middleware.GetUserIDFromContext(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	taskIDStr := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	resp, err := h.taskService.GetTask(r.Context(), taskID, userID)
	if err != nil {
		if err == app.ErrTaskNotFound {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		if err == app.ErrUnauthorized {
			writeError(w, http.StatusForbidden, "unauthorized")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get task")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := middleware.GetUserIDFromContext(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	taskIDStr := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req app.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.taskService.UpdateTask(r.Context(), taskID, &req, userID)
	if err != nil {
		if err == app.ErrTaskNotFound {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		if err == app.ErrUnauthorized {
			writeError(w, http.StatusForbidden, "unauthorized")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := middleware.GetUserIDFromContext(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	taskIDStr := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	err = h.taskService.DeleteTask(r.Context(), taskID, userID)
	if err != nil {
		if err == app.ErrTaskNotFound {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		if err == app.ErrUnauthorized {
			writeError(w, http.StatusForbidden, "unauthorized")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := middleware.GetUserIDFromContext(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	query := r.URL.Query()
	req := &app.ListTasksRequest{}

	if statusStr := query.Get("status"); statusStr != "" {
		status := task.TaskStatus(statusStr)
		req.Status = &status
	}

	if priorityStr := query.Get("priority"); priorityStr != "" {
		priority := task.TaskPriority(priorityStr)
		req.Priority = &priority
	}

	if deadlineStr := query.Get("deadline"); deadlineStr != "" {
		deadline, err := time.Parse(time.RFC3339, deadlineStr)
		if err == nil {
			req.Deadline = &deadline
		}
	}

	resp, err := h.taskService.ListTasks(r.Context(), req, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) Complete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := middleware.GetUserIDFromContext(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/api/tasks/")
	path = strings.TrimSuffix(path, "/complete")
	taskID, err := uuid.Parse(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	resp, err := h.taskService.CompleteTask(r.Context(), taskID, userID)
	if err != nil {
		if err == app.ErrTaskNotFound {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		if err == app.ErrUnauthorized {
			writeError(w, http.StatusForbidden, "unauthorized")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to complete task")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
