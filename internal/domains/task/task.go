package task

import (
	"context"
	"taskhub/pkg/base/entity"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusTodo    TaskStatus = "todo"
	StatusPending TaskStatus = "pending"
	StatusDone    TaskStatus = "done"
)

type Task struct {
	entity.BaseEntity
	Title    string
	Content  string
	Status   TaskStatus
	Deadline *time.Time
}

func NewTask(ctx context.Context, task *Task) *Task {
	return &Task{
		BaseEntity: entity.BaseEntity{
			Id:        uuid.New(),
			CreatedAt: time.Now(),
			CreatedBy: uuid.New(), //TODO: get current user
		},
		Title:    task.Title,
		Content:  task.Content,
		Status:   StatusTodo,
		Deadline: task.Deadline,
	}
}
