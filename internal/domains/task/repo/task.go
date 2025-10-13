package repo

import (
	"context"
	"database/sql"
	"taskhub/config"
	"taskhub/internal/domains/task"
	"taskhub/pkg/db"

	"github.com/google/uuid"
)

type TaskRepository struct {
	conn *sql.DB
}

func NewTaskRepository(config *config.Config) *TaskRepository {
	conn := db.NewDB(config).GetConnection()
	return &TaskRepository{conn}
}

func (r *TaskRepository) Create(ctx context.Context, task *task.Task) (*task.Task, error) {
	return nil, nil
}

func (r *TaskRepository) UpdateById(ctx context.Context, id uuid.UUID, task *task.Task) (*task.Task, error) {
	return nil, nil
}

func (r *TaskRepository) FindById(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	return nil, nil
}

func (r *TaskRepository) FindAll(ctx context.Context) ([]*task.Task, error) {
	return nil, nil
}

func (r *TaskRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	return nil
}
