package app

import (
	"context"
	"taskhub/internal/domains/task"
	"taskhub/internal/domains/task/repo"
	"taskhub/pkg/logger"

	"github.com/google/uuid"
)

type TaskService struct {
	logger         *logger.Logger
	taskRepository *repo.TaskRepository
}

func NewTaskService(
	taskRepository *repo.TaskRepository,
) *TaskService {
	return &TaskService{
		logger:         logger.NewLogger(),
		taskRepository: taskRepository,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *task.Task) (*task.Task, error) {
	res, err := s.taskRepository.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type RequestUpdateTask struct {
	id   uuid.UUID
	task *task.Task
}

func (s *TaskService) UpdateTask(ctx context.Context, req *RequestUpdateTask) (*task.Task, error) {
	res, err := s.taskRepository.UpdateById(ctx, req.id, req.task)
	if err != nil {
		return nil, err
	}

	return res, nil
}
