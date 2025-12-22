package app

import (
	"context"
	"encoding/json"
	"taskhub/internal/domains/notification"
	"taskhub/internal/domains/task"
	taskrepo "taskhub/internal/domains/task/repo"
	"taskhub/pkg/logger"
	natsconn "taskhub/pkg/nats"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

var NotificationServiceModule = fx.Module(
	"notification-service",
	fx.Provide(NewNotificationService),
)

const (
	SubjectTaskCreated  = "task.created"
	SubjectTaskUpdated  = "task.updated"
	SubjectTaskReminder = "task.reminder"
)

type NotificationService struct {
	logger   *logger.Logger
	nats     *natsconn.Nats
	taskRepo *taskrepo.TaskRepository
}

func NewNotificationService(logger *logger.Logger, nats *natsconn.Nats, taskRepo *taskrepo.TaskRepository) *NotificationService {
	return &NotificationService{
		logger:   logger,
		nats:     nats,
		taskRepo: taskRepo,
	}
}

type TaskEvent struct {
	EventType string     `json:"event_type"`
	TaskID    uuid.UUID  `json:"task_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Title     string     `json:"title"`
	Deadline  *time.Time `json:"deadline,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (s *NotificationService) PublishTaskCreated(ctx context.Context, t *task.Task) error {
	event := &TaskEvent{
		EventType: SubjectTaskCreated,
		TaskID:    t.Id,
		UserID:    t.UserID,
		Title:     t.Title,
		Deadline:  t.Deadline,
		CreatedAt: time.Now(),
	}

	return s.publishEvent(SubjectTaskCreated, event)
}

func (s *NotificationService) PublishTaskUpdated(ctx context.Context, t *task.Task) error {
	event := &TaskEvent{
		EventType: SubjectTaskUpdated,
		TaskID:    t.Id,
		UserID:    t.UserID,
		Title:     t.Title,
		Deadline:  t.Deadline,
		CreatedAt: time.Now(),
	}

	return s.publishEvent(SubjectTaskUpdated, event)
}

func (s *NotificationService) publishEvent(subject string, event *TaskEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.nats.Publish(subject, data)
}

type ReminderNotification struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Deadline  time.Time `json:"deadline"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *NotificationService) CheckAndSendReminders(ctx context.Context, hoursAhead int) ([]*notification.Notification, error) {
	tasks, err := s.taskRepo.FindTasksNearDeadline(ctx, hoursAhead)
	if err != nil {
		return nil, err
	}

	var notifications []*notification.Notification
	for _, t := range tasks {
		reminder := &ReminderNotification{
			ID:        uuid.New(),
			TaskID:    t.Id,
			UserID:    t.UserID,
			Title:     t.Title,
			Message:   "Task deadline is approaching",
			Deadline:  *t.Deadline,
			CreatedAt: time.Now(),
		}

		data, err := json.Marshal(reminder)
		if err != nil {
			s.logger.Error("failed to marshal reminder", "error", err)
			continue
		}

		if err := s.nats.Publish(SubjectTaskReminder, data); err != nil {
			s.logger.Error("failed to publish reminder", "error", err)
			continue
		}

		s.logger.Info("reminder sent for task", "task_id", t.Id, "user_id", t.UserID)
	}

	return notifications, nil
}

func (s *NotificationService) SubscribeToReminders(ctx context.Context, handler func(*ReminderNotification)) error {
	return s.nats.Subscribe(SubjectTaskReminder, func(data []byte) {
		var reminder ReminderNotification
		if err := json.Unmarshal(data, &reminder); err != nil {
			s.logger.Error("failed to unmarshal reminder", "error", err)
			return
		}
		handler(&reminder)
	})
}
