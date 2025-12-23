package handler

import (
	"encoding/json"
	"fmt"
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

	isHTMX := isHTMXRequest(r)
	var req app.CreateTaskRequest

	if isHTMX {
		req.Title = r.FormValue("title")
		req.Description = r.FormValue("description")
		req.Priority = task.TaskPriority(r.FormValue("priority"))
		if deadlineStr := r.FormValue("deadline"); deadlineStr != "" {
			deadline, err := time.Parse("2006-01-02T15:04", deadlineStr)
			if err == nil {
				req.Deadline = &deadline
			}
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
	}

	if req.Title == "" {
		if isHTMX {
			writeHTMXError(w, "Title is required")
			return
		}
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	resp, err := h.taskService.CreateTask(r.Context(), &req, userID)
	if err != nil {
		if isHTMX {
			writeHTMXError(w, "Failed to create task")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	if isHTMX {
		w.Header().Set("HX-Trigger", "taskCreated")
		h.renderTaskCard(w, resp.Task)
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

	isHTMX := isHTMXRequest(r)
	var req app.UpdateTaskRequest

	if isHTMX {
		req.Title = r.FormValue("title")
		req.Description = r.FormValue("description")
		req.Status = task.TaskStatus(r.FormValue("status"))
		req.Priority = task.TaskPriority(r.FormValue("priority"))
		if deadlineStr := r.FormValue("deadline"); deadlineStr != "" {
			deadline, err := time.Parse("2006-01-02T15:04", deadlineStr)
			if err == nil {
				req.Deadline = &deadline
			}
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
	}

	resp, err := h.taskService.UpdateTask(r.Context(), taskID, &req, userID)
	if err != nil {
		if err == app.ErrTaskNotFound {
			if isHTMX {
				writeHTMXError(w, "Task not found")
				return
			}
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		if err == app.ErrUnauthorized {
			if isHTMX {
				writeHTMXError(w, "Unauthorized")
				return
			}
			writeError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if isHTMX {
			writeHTMXError(w, "Failed to update task")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	if isHTMX {
		w.Header().Set("HX-Trigger", "taskUpdated")
		h.renderTaskCard(w, resp.Task)
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

	isHTMX := isHTMXRequest(r)
	err = h.taskService.DeleteTask(r.Context(), taskID, userID)
	if err != nil {
		if err == app.ErrTaskNotFound {
			if isHTMX {
				writeHTMXError(w, "Task not found")
				return
			}
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		if err == app.ErrUnauthorized {
			if isHTMX {
				writeHTMXError(w, "Unauthorized")
				return
			}
			writeError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if isHTMX {
			writeHTMXError(w, "Failed to delete task")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	if isHTMX {
		w.Header().Set("HX-Trigger", "taskDeleted")
		w.WriteHeader(http.StatusNoContent)
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

	if searchStr := query.Get("search"); searchStr != "" {
		req.Search = searchStr
	}

	resp, err := h.taskService.ListTasks(r.Context(), req, userID)
	if err != nil {
		if isHTMXRequest(r) {
			writeHTMXError(w, "Failed to load tasks")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}

	if isHTMXRequest(r) {
		if len(resp.Tasks) == 0 {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<div class="empty-state">
				<h3>No tasks found</h3>
				<p>Create your first task to get started!</p>
			</div>`)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		for _, task := range resp.Tasks {
			h.renderTaskCard(w, task)
		}
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

	isHTMX := isHTMXRequest(r)
	resp, err := h.taskService.CompleteTask(r.Context(), taskID, userID)
	if err != nil {
		if err == app.ErrTaskNotFound {
			if isHTMX {
				writeHTMXError(w, "Task not found")
				return
			}
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		if err == app.ErrUnauthorized {
			if isHTMX {
				writeHTMXError(w, "Unauthorized")
				return
			}
			writeError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if isHTMX {
			writeHTMXError(w, "Failed to complete task")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to complete task")
		return
	}

	if isHTMX {
		w.Header().Set("HX-Trigger", "taskCompleted")
		h.renderTaskCard(w, resp.Task)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) renderTaskCard(w http.ResponseWriter, t *task.Task) {
	w.Header().Set("Content-Type", "text/html")

	statusClass := "badge-" + string(t.Status)
	priorityClass := "badge-" + string(t.Priority)

	deadlineText := ""
	if t.Deadline != nil {
		deadlineText = t.Deadline.Format("Jan 2, 2006 3:04 PM")
	}

	fmt.Fprintf(w, `
	<div class="task-card" id="task-%s">
		<div class="task-info">
			<h4>%s</h4>
			<p>%s</p>
			<div class="task-meta">
				<span class="badge %s">%s</span>
				<span class="badge %s">%s</span>
				%s
			</div>
		</div>
		<div class="task-actions">
			%s
			<button class="btn btn-sm btn-outline" onclick="editTask('%s')">Edit</button>
			<button class="btn btn-sm btn-danger" hx-delete="/api/tasks/%s" hx-target="#task-%s" hx-swap="outerHTML">Delete</button>
		</div>
	</div>`,
		t.Id.String(),
		t.Title,
		t.Description,
		statusClass, string(t.Status),
		priorityClass, string(t.Priority),
		func() string {
			if deadlineText != "" {
				return fmt.Sprintf(`<span class="badge">📅 %s</span>`, deadlineText)
			}
			return ""
		}(),
		func() string {
			if t.Status != task.StatusDone {
				return fmt.Sprintf(`<button class="btn btn-sm btn-success" hx-post="/api/tasks/%s/complete" hx-target="#task-%s" hx-swap="outerHTML">✓</button>`, t.Id.String(), t.Id.String())
			}
			return ""
		}(),
		t.Id.String(),
		t.Id.String(),
		t.Id.String(),
	)
}
