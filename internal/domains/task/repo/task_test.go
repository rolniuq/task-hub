package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"taskhub/internal/domains/task"
	"taskhub/pkg/base/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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
		return nil, errors.New("task not found")
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
	return errors.New("task not found")
}

func (m *MockTaskRepository) FindTasksNearDeadline(ctx context.Context, hoursAhead int) ([]*task.Task, error) {
	result := make([]*task.Task, 0)
	now := time.Now()
	for _, t := range m.tasks {
		if t.Status != task.StatusDone && t.Deadline != nil && t.Deadline.Before(now.Add(time.Duration(hoursAhead)*time.Hour)) {
			result = append(result, t)
		}
	}
	return result, nil
}

func TestTaskRepositoryModule(t *testing.T) {
	assert.NotNil(t, TaskRepositoryModule)
}

func TestMockTaskRepository_Create(t *testing.T) {
	repo := NewMockTaskRepository()
	ctx := context.Background()

	taskID := uuid.New()
	userID := uuid.New()
	newTask := &task.Task{
		BaseEntity: entity.BaseEntity{
			Id:        taskID,
			CreatedAt: time.Now(),
		},
		Title:       "Test Task",
		Description: "Test Description",
		Status:      task.StatusTodo,
		Priority:    task.PriorityHigh,
		UserID:      userID,
	}

	created, err := repo.Create(ctx, newTask)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, taskID, created.Id)
	assert.Equal(t, "Test Task", created.Title)
}

func TestMockTaskRepository_FindById(t *testing.T) {
	repo := NewMockTaskRepository()
	ctx := context.Background()

	taskID := uuid.New()
	userID := uuid.New()
	newTask := &task.Task{
		BaseEntity: entity.BaseEntity{
			Id:        taskID,
			CreatedAt: time.Now(),
		},
		Title:  "Test Task",
		Status: task.StatusTodo,
		UserID: userID,
	}

	repo.Create(ctx, newTask)

	// Test find existing
	found, err := repo.FindById(ctx, taskID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "Test Task", found.Title)

	// Test find non-existing
	notFound, err := repo.FindById(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestMockTaskRepository_UpdateById(t *testing.T) {
	repo := NewMockTaskRepository()
	ctx := context.Background()

	taskID := uuid.New()
	userID := uuid.New()
	newTask := &task.Task{
		BaseEntity: entity.BaseEntity{
			Id:        taskID,
			CreatedAt: time.Now(),
		},
		Title:  "Original Title",
		Status: task.StatusTodo,
		UserID: userID,
	}

	repo.Create(ctx, newTask)

	// Update task
	updatedTask := &task.Task{
		Title:       "Updated Title",
		Description: "Updated Description",
		Status:      task.StatusInProgress,
		Priority:    task.PriorityMedium,
	}

	updated, err := repo.UpdateById(ctx, taskID, updatedTask)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, task.StatusInProgress, updated.Status)

	// Test update non-existing
	_, err = repo.UpdateById(ctx, uuid.New(), updatedTask)
	assert.Error(t, err)
}

func TestMockTaskRepository_FindAll(t *testing.T) {
	repo := NewMockTaskRepository()
	ctx := context.Background()

	userID := uuid.New()
	status := task.StatusTodo
	priority := task.PriorityHigh

	// Create test tasks
	task1 := &task.Task{
		BaseEntity: entity.BaseEntity{Id: uuid.New()},
		Title:      "Task 1",
		Status:     task.StatusTodo,
		Priority:   task.PriorityHigh,
		UserID:     userID,
	}
	task2 := &task.Task{
		BaseEntity: entity.BaseEntity{Id: uuid.New()},
		Title:      "Task 2",
		Status:     task.StatusInProgress,
		Priority:   task.PriorityMedium,
		UserID:     userID,
	}

	repo.Create(ctx, task1)
	repo.Create(ctx, task2)

	// Test find all without filter
	all, err := repo.FindAll(ctx, nil)
	assert.NoError(t, err)
	assert.Len(t, all, 2)

	// Test find all with status filter
	filter := &task.TaskFilter{
		Status: &status,
	}
	filtered, err := repo.FindAll(ctx, filter)
	assert.NoError(t, err)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "Task 1", filtered[0].Title)

	// Test find all with priority filter
	filter = &task.TaskFilter{
		Priority: &priority,
	}
	filtered, err = repo.FindAll(ctx, filter)
	assert.NoError(t, err)
	assert.Len(t, filtered, 1)

	// Test find all with user filter
	filter = &task.TaskFilter{
		UserID: &userID,
	}
	filtered, err = repo.FindAll(ctx, filter)
	assert.NoError(t, err)
	assert.Len(t, filtered, 2)
}

func TestMockTaskRepository_FindByUserId(t *testing.T) {
	repo := NewMockTaskRepository()
	ctx := context.Background()

	userID1 := uuid.New()
	userID2 := uuid.New()

	task1 := &task.Task{
		BaseEntity: entity.BaseEntity{Id: uuid.New()},
		Title:      "User 1 Task",
		UserID:     userID1,
	}
	task2 := &task.Task{
		BaseEntity: entity.BaseEntity{Id: uuid.New()},
		Title:      "User 2 Task",
		UserID:     userID2,
	}

	repo.Create(ctx, task1)
	repo.Create(ctx, task2)

	// Find tasks for user 1
	tasks, err := repo.FindByUserId(ctx, userID1, nil)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "User 1 Task", tasks[0].Title)
}

func TestMockTaskRepository_DeleteById(t *testing.T) {
	repo := NewMockTaskRepository()
	ctx := context.Background()

	taskID := uuid.New()
	userID := uuid.New()
	newTask := &task.Task{
		BaseEntity: entity.BaseEntity{Id: taskID},
		Title:      "Task to Delete",
		UserID:     userID,
	}

	repo.Create(ctx, newTask)

	// Verify task exists
	found, _ := repo.FindById(ctx, taskID)
	assert.NotNil(t, found)

	// Delete task
	err := repo.DeleteById(ctx, taskID, userID)
	assert.NoError(t, err)

	// Verify task is deleted
	notFound, _ := repo.FindById(ctx, taskID)
	assert.Nil(t, notFound)
}

func TestMockTaskRepository_MarkAsCompleted(t *testing.T) {
	repo := NewMockTaskRepository()
	ctx := context.Background()

	taskID := uuid.New()
	userID := uuid.New()
	newTask := &task.Task{
		BaseEntity: entity.BaseEntity{Id: taskID},
		Title:      "Task to Complete",
		Status:     task.StatusTodo,
		UserID:     userID,
	}

	repo.Create(ctx, newTask)

	// Mark as completed
	err := repo.MarkAsCompleted(ctx, taskID, userID)
	assert.NoError(t, err)

	// Verify status
	completed, _ := repo.FindById(ctx, taskID)
	assert.Equal(t, task.StatusDone, completed.Status)

	// Test complete non-existing
	err = repo.MarkAsCompleted(ctx, uuid.New(), userID)
	assert.Error(t, err)
}

func TestMockTaskRepository_FindTasksNearDeadline(t *testing.T) {
	repo := NewMockTaskRepository()
	ctx := context.Background()

	deadlineSoon := time.Now().Add(2 * time.Hour)
	deadlineFar := time.Now().Add(48 * time.Hour)

	task1 := &task.Task{
		BaseEntity: entity.BaseEntity{Id: uuid.New()},
		Title:      "Urgent Task",
		Status:     task.StatusTodo,
		Deadline:   &deadlineSoon,
	}
	task2 := &task.Task{
		BaseEntity: entity.BaseEntity{Id: uuid.New()},
		Title:      "Future Task",
		Status:     task.StatusTodo,
		Deadline:   &deadlineFar,
	}
	task3 := &task.Task{
		BaseEntity: entity.BaseEntity{Id: uuid.New()},
		Title:      "Completed Task",
		Status:     task.StatusDone,
		Deadline:   &deadlineSoon,
	}

	repo.Create(ctx, task1)
	repo.Create(ctx, task2)
	repo.Create(ctx, task3)

	// Find tasks within 24 hours
	tasks, err := repo.FindTasksNearDeadline(ctx, 24)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Urgent Task", tasks[0].Title)
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
