package task

import (
	"context"
	"taskhub/pkg/base/entity"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string
type TaskPriority string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	entity.BaseEntity
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	Deadline    *time.Time   `json:"deadline,omitempty"`
	UserID      uuid.UUID    `json:"user_id"`
}

func NewTask(ctx context.Context, t *Task, userID uuid.UUID) *Task {
	now := time.Now()
	return &Task{
		BaseEntity: entity.BaseEntity{
			Id:        uuid.New(),
			CreatedAt: now,
			CreatedBy: userID,
		},
		Title:       t.Title,
		Description: t.Description,
		Status:      StatusTodo,
		Priority:    t.Priority,
		Deadline:    t.Deadline,
		UserID:      userID,
	}
}

type TaskFilter struct {
	Status   *TaskStatus
	Priority *TaskPriority
	UserID   *uuid.UUID
	Deadline *time.Time
}

func (t *Task) MarkAsCompleted(userID uuid.UUID) {
	now := time.Now()
	t.Status = StatusDone
	t.UpdateAt = &now
	t.UpdateBy = &userID
}

func (t *Task) MarkAsInProgress(userID uuid.UUID) {
	now := time.Now()
	t.Status = StatusInProgress
	t.UpdateAt = &now
	t.UpdateBy = &userID
}
