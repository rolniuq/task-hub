package app

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskEvent(t *testing.T) {
	deadline := time.Now().Add(24 * time.Hour)
	event := &TaskEvent{
		EventType: SubjectTaskCreated,
		TaskID:    uuid.New(),
		UserID:    uuid.New(),
		Title:     "Test Task",
		Deadline:  &deadline,
		CreatedAt: time.Now(),
	}

	assert.Equal(t, SubjectTaskCreated, event.EventType)
	assert.NotEqual(t, uuid.Nil, event.TaskID)
	assert.NotEqual(t, uuid.Nil, event.UserID)
	assert.Equal(t, "Test Task", event.Title)
	assert.NotNil(t, event.Deadline)
	assert.False(t, event.CreatedAt.IsZero())
}

func TestReminderNotification(t *testing.T) {
	deadline := time.Now().Add(1 * time.Hour)
	reminder := &ReminderNotification{
		ID:        uuid.New(),
		TaskID:    uuid.New(),
		UserID:    uuid.New(),
		Title:     "Urgent Task",
		Message:   "Task deadline is approaching",
		Deadline:  deadline,
		CreatedAt: time.Now(),
	}

	assert.NotEqual(t, uuid.Nil, reminder.ID)
	assert.NotEqual(t, uuid.Nil, reminder.TaskID)
	assert.NotEqual(t, uuid.Nil, reminder.UserID)
	assert.Equal(t, "Urgent Task", reminder.Title)
	assert.Equal(t, "Task deadline is approaching", reminder.Message)
	assert.False(t, reminder.Deadline.IsZero())
	assert.False(t, reminder.CreatedAt.IsZero())
}

func TestSubjectConstants(t *testing.T) {
	assert.Equal(t, "task.created", SubjectTaskCreated)
	assert.Equal(t, "task.updated", SubjectTaskUpdated)
	assert.Equal(t, "task.reminder", SubjectTaskReminder)
}

func TestTaskEventWithoutDeadline(t *testing.T) {
	event := &TaskEvent{
		EventType: SubjectTaskUpdated,
		TaskID:    uuid.New(),
		UserID:    uuid.New(),
		Title:     "Task Without Deadline",
		Deadline:  nil,
		CreatedAt: time.Now(),
	}

	assert.Equal(t, SubjectTaskUpdated, event.EventType)
	assert.Nil(t, event.Deadline)
}

func TestNotificationServiceModule(t *testing.T) {
	assert.NotNil(t, NotificationServiceModule)
}

func TestNewNotificationService_Nil(t *testing.T) {
	service := NewNotificationService(nil, nil, nil)
	assert.NotNil(t, service)
	assert.Nil(t, service.logger)
	assert.Nil(t, service.nats)
	assert.Nil(t, service.taskRepo)
}
