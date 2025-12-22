package task

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	deadline := time.Now().Add(24 * time.Hour)

	input := &Task{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    PriorityHigh,
		Deadline:    &deadline,
	}

	newTask := NewTask(ctx, input, userID)

	assert.NotNil(t, newTask)
	assert.NotEqual(t, uuid.Nil, newTask.Id)
	assert.Equal(t, "Test Task", newTask.Title)
	assert.Equal(t, "Test Description", newTask.Description)
	assert.Equal(t, StatusTodo, newTask.Status)
	assert.Equal(t, PriorityHigh, newTask.Priority)
	assert.NotNil(t, newTask.Deadline)
	assert.Equal(t, userID, newTask.UserID)
	assert.Equal(t, userID, newTask.CreatedBy)
	assert.False(t, newTask.CreatedAt.IsZero())
}

func TestNewTask_WithoutDeadline(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	input := &Task{
		Title:       "Task Without Deadline",
		Description: "No deadline set",
		Priority:    PriorityMedium,
	}

	newTask := NewTask(ctx, input, userID)

	assert.NotNil(t, newTask)
	assert.Nil(t, newTask.Deadline)
	assert.Equal(t, PriorityMedium, newTask.Priority)
}

func TestMarkAsCompleted(t *testing.T) {
	userID := uuid.New()
	task := &Task{
		Status: StatusTodo,
	}

	task.MarkAsCompleted(userID)

	assert.Equal(t, StatusDone, task.Status)
	assert.NotNil(t, task.UpdateAt)
	assert.NotNil(t, task.UpdateBy)
	assert.Equal(t, userID, *task.UpdateBy)
}

func TestMarkAsInProgress(t *testing.T) {
	userID := uuid.New()
	task := &Task{
		Status: StatusTodo,
	}

	task.MarkAsInProgress(userID)

	assert.Equal(t, StatusInProgress, task.Status)
	assert.NotNil(t, task.UpdateAt)
	assert.NotNil(t, task.UpdateBy)
	assert.Equal(t, userID, *task.UpdateBy)
}

func TestTaskStatus_Constants(t *testing.T) {
	assert.Equal(t, TaskStatus("todo"), StatusTodo)
	assert.Equal(t, TaskStatus("in_progress"), StatusInProgress)
	assert.Equal(t, TaskStatus("done"), StatusDone)
}

func TestTaskPriority_Constants(t *testing.T) {
	assert.Equal(t, TaskPriority("low"), PriorityLow)
	assert.Equal(t, TaskPriority("medium"), PriorityMedium)
	assert.Equal(t, TaskPriority("high"), PriorityHigh)
}

func TestTaskFilter(t *testing.T) {
	status := StatusTodo
	priority := PriorityHigh
	userID := uuid.New()
	deadline := time.Now()

	filter := &TaskFilter{
		Status:   &status,
		Priority: &priority,
		UserID:   &userID,
		Deadline: &deadline,
	}

	assert.Equal(t, StatusTodo, *filter.Status)
	assert.Equal(t, PriorityHigh, *filter.Priority)
	assert.Equal(t, userID, *filter.UserID)
	assert.NotNil(t, filter.Deadline)
}

func TestTaskFilter_Empty(t *testing.T) {
	filter := &TaskFilter{}

	assert.Nil(t, filter.Status)
	assert.Nil(t, filter.Priority)
	assert.Nil(t, filter.UserID)
	assert.Nil(t, filter.Deadline)
}

func TestTask_Fields(t *testing.T) {
	deadline := time.Now().Add(48 * time.Hour)
	userID := uuid.New()

	task := &Task{
		Title:       "Full Task",
		Description: "Complete task with all fields",
		Status:      StatusInProgress,
		Priority:    PriorityHigh,
		Deadline:    &deadline,
		UserID:      userID,
	}

	assert.Equal(t, "Full Task", task.Title)
	assert.Equal(t, "Complete task with all fields", task.Description)
	assert.Equal(t, StatusInProgress, task.Status)
	assert.Equal(t, PriorityHigh, task.Priority)
	assert.NotNil(t, task.Deadline)
	assert.Equal(t, userID, task.UserID)
}

func TestTask_StatusTransitions(t *testing.T) {
	userID := uuid.New()

	task := &Task{
		Status: StatusTodo,
	}
	assert.Equal(t, StatusTodo, task.Status)

	task.MarkAsInProgress(userID)
	assert.Equal(t, StatusInProgress, task.Status)

	task.MarkAsCompleted(userID)
	assert.Equal(t, StatusDone, task.Status)
}

func TestTask_MultipleStatusUpdates(t *testing.T) {
	userID := uuid.New()

	task := &Task{
		Status: StatusTodo,
	}

	task.MarkAsInProgress(userID)
	firstUpdate := task.UpdateAt

	time.Sleep(10 * time.Millisecond)

	task.MarkAsCompleted(userID)
	secondUpdate := task.UpdateAt

	assert.True(t, secondUpdate.After(*firstUpdate))
}
