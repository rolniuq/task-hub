package app

import (
	"context"
	"errors"
	"strings"
	"taskhub/internal/domains/task"
	"taskhub/internal/domains/task/repo"
	"taskhub/pkg/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

var TaskServiceModule = fx.Module(
	"task-service",
	fx.Provide(NewTaskService),
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrUnauthorized = errors.New("unauthorized")
)

type TaskService struct {
	logger   *logger.Logger
	taskRepo *repo.TaskRepository
}

func NewTaskService(logger *logger.Logger, taskRepo *repo.TaskRepository) *TaskService {
	return &TaskService{
		logger:   logger,
		taskRepo: taskRepo,
	}
}

type CreateTaskRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Priority    task.TaskPriority `json:"priority"`
	Deadline    *time.Time        `json:"deadline,omitempty"`
}

type TaskResponse struct {
	Task *task.Task `json:"task"`
}

func (s *TaskService) CreateTask(ctx context.Context, req *CreateTaskRequest, userID uuid.UUID) (*TaskResponse, error) {
	newTask := task.NewTask(ctx, &task.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Deadline:    req.Deadline,
	}, userID)

	createdTask, err := s.taskRepo.Create(ctx, newTask)
	if err != nil {
		return nil, err
	}

	return &TaskResponse{Task: createdTask}, nil
}

type UpdateTaskRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      task.TaskStatus   `json:"status"`
	Priority    task.TaskPriority `json:"priority"`
	Deadline    *time.Time        `json:"deadline,omitempty"`
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID uuid.UUID, req *UpdateTaskRequest, userID uuid.UUID) (*TaskResponse, error) {
	existingTask, err := s.taskRepo.FindById(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if existingTask == nil {
		return nil, ErrTaskNotFound
	}

	if existingTask.UserID != userID {
		return nil, ErrUnauthorized
	}

	now := time.Now()
	existingTask.Title = req.Title
	existingTask.Description = req.Description
	existingTask.Status = req.Status
	existingTask.Priority = req.Priority
	existingTask.Deadline = req.Deadline
	existingTask.UpdateAt = &now
	existingTask.UpdateBy = &userID

	updatedTask, err := s.taskRepo.UpdateById(ctx, taskID, existingTask)
	if err != nil {
		return nil, err
	}

	return &TaskResponse{Task: updatedTask}, nil
}

func (s *TaskService) GetTask(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) (*TaskResponse, error) {
	t, err := s.taskRepo.FindById(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrTaskNotFound
	}

	if t.UserID != userID {
		return nil, ErrUnauthorized
	}

	return &TaskResponse{Task: t}, nil
}

type ListTasksRequest struct {
	Status   *task.TaskStatus   `json:"status,omitempty"`
	Priority *task.TaskPriority `json:"priority,omitempty"`
	Deadline *time.Time         `json:"deadline,omitempty"`
	Search   string             `json:"search,omitempty"`
}

type ListTasksResponse struct {
	Tasks []*task.Task `json:"tasks"`
}

func (s *TaskService) ListTasks(ctx context.Context, req *ListTasksRequest, userID uuid.UUID) (*ListTasksResponse, error) {
	filter := &task.TaskFilter{
		Status:   req.Status,
		Priority: req.Priority,
		Deadline: req.Deadline,
		UserID:   &userID,
	}

	tasks, err := s.taskRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Apply search filter client-side for now
	if req.Search != "" {
		var filteredTasks []*task.Task
		searchLower := strings.ToLower(req.Search)
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), searchLower) ||
				strings.Contains(strings.ToLower(t.Description), searchLower) {
				filteredTasks = append(filteredTasks, t)
			}
		}
		tasks = filteredTasks
	}

	return &ListTasksResponse{Tasks: tasks}, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) error {
	existingTask, err := s.taskRepo.FindById(ctx, taskID)
	if err != nil {
		return err
	}
	if existingTask == nil {
		return ErrTaskNotFound
	}

	if existingTask.UserID != userID {
		return ErrUnauthorized
	}

	return s.taskRepo.DeleteById(ctx, taskID, userID)
}

func (s *TaskService) CompleteTask(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) (*TaskResponse, error) {
	existingTask, err := s.taskRepo.FindById(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if existingTask == nil {
		return nil, ErrTaskNotFound
	}

	if existingTask.UserID != userID {
		return nil, ErrUnauthorized
	}

	if err := s.taskRepo.MarkAsCompleted(ctx, taskID, userID); err != nil {
		return nil, err
	}

	existingTask.Status = task.StatusDone
	return &TaskResponse{Task: existingTask}, nil
}
