package repo

import (
	"context"
	"database/sql"
	"taskhub/config"
	"taskhub/internal/db"
	"taskhub/internal/domains/task"
	db2 "taskhub/pkg/db"
	"taskhub/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

var TaskRepositoryModule = fx.Module(
	"task-repo",
	fx.Provide(NewTaskRepository),
)

type TaskRepository struct {
	queries *db.Queries
	logger  *logger.Logger
}

func NewTaskRepository(config *config.Config, logger *logger.Logger) *TaskRepository {
	conn := db2.NewDB(config).GetConnection()
	queries := db.New(conn)
	return &TaskRepository{
		queries: queries,
		logger:  logger,
	}
}

// Helper function to convert domain Task to sqlc CreateTaskParams
func toCreateTaskParams(t *task.Task) db.CreateTaskParams {
	params := db.CreateTaskParams{
		ID:        t.Id,
		Title:     t.Title,
		Status:    string(t.Status),
		Priority:  string(t.Priority),
		UserID:    t.UserID,
		CreatedAt: t.CreatedAt,
		CreatedBy: t.CreatedBy,
	}

	if t.Description != "" {
		params.Description = sql.NullString{String: t.Description, Valid: true}
	}
	if t.Deadline != nil {
		params.Deadline = sql.NullTime{Time: *t.Deadline, Valid: true}
	}

	return params
}

// Helper function to convert sqlc Task to domain Task
func toDomainTask(t db.Task) *task.Task {
	result := &task.Task{
		Title:    t.Title,
		Status:   task.TaskStatus(t.Status),
		Priority: task.TaskPriority(t.Priority),
		UserID:   t.UserID,
	}

	// Set BaseEntity fields
	result.Id = t.ID
	result.CreatedAt = t.CreatedAt
	result.CreatedBy = t.CreatedBy

	if t.Description.Valid {
		result.Description = t.Description.String
	}
	if t.Deadline.Valid {
		result.Deadline = &t.Deadline.Time
	}
	if t.UpdatedAt.Valid {
		result.UpdateAt = &t.UpdatedAt.Time
	}
	if t.UpdatedBy.Valid {
		result.UpdateBy = &t.UpdatedBy.UUID
	}

	return result
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Task) (*task.Task, error) {
	params := toCreateTaskParams(t)

	created, err := r.queries.CreateTask(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainTask(created), nil
}

func (r *TaskRepository) UpdateById(ctx context.Context, id uuid.UUID, t *task.Task) (*task.Task, error) {
	params := db.UpdateTaskParams{
		ID:       id,
		Title:    t.Title,
		Status:   string(t.Status),
		Priority: string(t.Priority),
	}

	if t.Description != "" {
		params.Description = sql.NullString{String: t.Description, Valid: true}
	}
	if t.Deadline != nil {
		params.Deadline = sql.NullTime{Time: *t.Deadline, Valid: true}
	}
	if t.UpdateAt != nil {
		params.UpdatedAt = sql.NullTime{Time: *t.UpdateAt, Valid: true}
	}
	if t.UpdateBy != nil {
		params.UpdatedBy = uuid.NullUUID{UUID: *t.UpdateBy, Valid: true}
	}

	updated, err := r.queries.UpdateTask(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return toDomainTask(updated), nil
}

func (r *TaskRepository) FindById(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	found, err := r.queries.GetTask(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return toDomainTask(found), nil
}

func (r *TaskRepository) FindAll(ctx context.Context, filter *task.TaskFilter) ([]*task.Task, error) {
	var dbTasks []db.Task
	var err error

	// sqlc generates specific queries for different filters
	// We use the most appropriate query based on the filter
	if filter != nil {
		switch {
		case filter.UserID != nil:
			dbTasks, err = r.queries.ListTasksByUser(ctx, *filter.UserID)
		case filter.Status != nil:
			dbTasks, err = r.queries.ListTasksByStatus(ctx, string(*filter.Status))
		case filter.Priority != nil:
			dbTasks, err = r.queries.ListTasksByPriority(ctx, string(*filter.Priority))
		default:
			dbTasks, err = r.queries.ListTasks(ctx)
		}
	} else {
		dbTasks, err = r.queries.ListTasks(ctx)
	}

	if err != nil {
		return nil, err
	}

	tasks := make([]*task.Task, len(dbTasks))
	for i, t := range dbTasks {
		tasks[i] = toDomainTask(t)
	}

	return tasks, nil
}

func (r *TaskRepository) FindByUserId(ctx context.Context, userID uuid.UUID, filter *task.TaskFilter) ([]*task.Task, error) {
	dbTasks, err := r.queries.ListTasksByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]*task.Task, len(dbTasks))
	for i, t := range dbTasks {
		tasks[i] = toDomainTask(t)
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteById(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	params := db.DeleteTaskParams{
		ID:        id,
		DeletedBy: uuid.NullUUID{UUID: userID, Valid: true},
	}

	rowsAffected, err := r.queries.DeleteTask(ctx, params)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TaskRepository) MarkAsCompleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	params := db.MarkTaskAsCompletedParams{
		ID:        id,
		UpdatedBy: uuid.NullUUID{UUID: userID, Valid: true},
	}

	rowsAffected, err := r.queries.MarkTaskAsCompleted(ctx, params)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TaskRepository) FindTasksNearDeadline(ctx context.Context, hoursAhead int) ([]*task.Task, error) {
	dbTasks, err := r.queries.GetTasksNearDeadline(ctx, hoursAhead)
	if err != nil {
		return nil, err
	}

	tasks := make([]*task.Task, len(dbTasks))
	for i, t := range dbTasks {
		tasks[i] = toDomainTask(t)
	}

	return tasks, nil
}
