package app

import (
	"context"
	"testing"
	"time"

	"taskhub/internal/domains/task"
	"taskhub/pkg/base/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTaskRequest(t *testing.T) {
	deadline := time.Now().Add(24 * time.Hour)
	req := &CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    task.PriorityHigh,
		Deadline:    &deadline,
	}

	assert.Equal(t, "Test Task", req.Title)
	assert.Equal(t, "Test Description", req.Description)
	assert.Equal(t, task.PriorityHigh, req.Priority)
	assert.NotNil(t, req.Deadline)
}

func TestUpdateTaskRequest(t *testing.T) {
	deadline := time.Now().Add(48 * time.Hour)
	req := &UpdateTaskRequest{
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      task.StatusInProgress,
		Priority:    task.PriorityMedium,
		Deadline:    &deadline,
	}

	assert.Equal(t, "Updated Task", req.Title)
	assert.Equal(t, task.StatusInProgress, req.Status)
	assert.Equal(t, task.PriorityMedium, req.Priority)
}

func TestListTasksRequest(t *testing.T) {
	status := task.StatusTodo
	priority := task.PriorityHigh
	deadline := time.Now()

	req := &ListTasksRequest{
		Status:   &status,
		Priority: &priority,
		Deadline: &deadline,
	}

	assert.NotNil(t, req.Status)
	assert.Equal(t, task.StatusTodo, *req.Status)
	assert.NotNil(t, req.Priority)
	assert.Equal(t, task.PriorityHigh, *req.Priority)
}

func TestTaskResponse(t *testing.T) {
	testTask := &task.Task{
		BaseEntity: entity.BaseEntity{
			Id:        uuid.New(),
			CreatedAt: time.Now(),
		},
		Title:       "Test Task",
		Description: "Description",
		Status:      task.StatusTodo,
		Priority:    task.PriorityMedium,
	}

	resp := &TaskResponse{Task: testTask}

	assert.NotNil(t, resp.Task)
	assert.Equal(t, "Test Task", resp.Task.Title)
}

func TestListTasksResponse(t *testing.T) {
	tasks := []*task.Task{
		{
			BaseEntity: entity.BaseEntity{Id: uuid.New()},
			Title:      "Task 1",
		},
		{
			BaseEntity: entity.BaseEntity{Id: uuid.New()},
			Title:      "Task 2",
		},
	}

	resp := &ListTasksResponse{Tasks: tasks}

	assert.Len(t, resp.Tasks, 2)
	assert.Equal(t, "Task 1", resp.Tasks[0].Title)
	assert.Equal(t, "Task 2", resp.Tasks[1].Title)
}

func TestTaskServiceErrors(t *testing.T) {
	assert.Equal(t, "task not found", ErrTaskNotFound.Error())
	assert.Equal(t, "unauthorized", ErrUnauthorized.Error())
}

func TestRequestUpdateTask(t *testing.T) {
	deadline := time.Now().Add(24 * time.Hour)
	req := &UpdateTaskRequest{
		Title:       "Task Title",
		Description: "Task Description",
		Status:      task.StatusDone,
		Priority:    task.PriorityLow,
		Deadline:    &deadline,
	}

	assert.Equal(t, "Task Title", req.Title)
	assert.Equal(t, "Task Description", req.Description)
	assert.Equal(t, task.StatusDone, req.Status)
	assert.Equal(t, task.PriorityLow, req.Priority)
	assert.NotNil(t, req.Deadline)
}

type MockTaskRepository struct {
	tasks map[uuid.UUID]*task.Task
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		tasks: make(map[uuid.UUID]*task.Task),
	}
}

func (m *MockTaskRepository) Create(ctx context.Context, t *task.Task) (*task.Task, error) {
	m.tasks[t.Id] = t
	return t, nil
}

func (m *MockTaskRepository) UpdateById(ctx context.Context, id uuid.UUID, t *task.Task) (*task.Task, error) {
	if _, ok := m.tasks[id]; !ok {
		return nil, nil
	}
	m.tasks[id] = t
	return t, nil
}

func (m *MockTaskRepository) FindById(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	if t, ok := m.tasks[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *MockTaskRepository) FindAll(ctx context.Context, filter *task.TaskFilter) ([]*task.Task, error) {
	result := make([]*task.Task, 0)
	for _, t := range m.tasks {
		if filter != nil {
			if filter.UserID != nil && t.UserID != *filter.UserID {
				continue
			}
			if filter.Status != nil && t.Status != *filter.Status {
				continue
			}
			if filter.Priority != nil && t.Priority != *filter.Priority {
				continue
			}
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *MockTaskRepository) FindByUserId(ctx context.Context, userID uuid.UUID, filter *task.TaskFilter) ([]*task.Task, error) {
	if filter == nil {
		filter = &task.TaskFilter{}
	}
	filter.UserID = &userID
	return m.FindAll(ctx, filter)
}

func (m *MockTaskRepository) DeleteById(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	delete(m.tasks, id)
	return nil
}

func (m *MockTaskRepository) MarkAsCompleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if t, ok := m.tasks[id]; ok {
		t.Status = task.StatusDone
		return nil
	}
	return nil
}

func TestTaskFilter(t *testing.T) {
	status := task.StatusTodo
	priority := task.PriorityHigh
	userID := uuid.New()
	deadline := time.Now()

	filter := &task.TaskFilter{
		Status:   &status,
		Priority: &priority,
		UserID:   &userID,
		Deadline: &deadline,
	}

	assert.Equal(t, task.StatusTodo, *filter.Status)
	assert.Equal(t, task.PriorityHigh, *filter.Priority)
	assert.Equal(t, userID, *filter.UserID)
	assert.NotNil(t, filter.Deadline)
}

func TestTaskStatus(t *testing.T) {
	tests := []struct {
		status   task.TaskStatus
		expected string
	}{
		{task.StatusTodo, "todo"},
		{task.StatusInProgress, "in_progress"},
		{task.StatusDone, "done"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestTaskPriority(t *testing.T) {
	tests := []struct {
		priority task.TaskPriority
		expected string
	}{
		{task.PriorityLow, "low"},
		{task.PriorityMedium, "medium"},
		{task.PriorityHigh, "high"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.priority))
		})
	}
}
