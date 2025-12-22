package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"taskhub/config"
	"taskhub/internal/domains/task"
	"taskhub/pkg/db"
	"taskhub/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

var TaskRepositoryModule = fx.Module(
	"task-repo",
	fx.Provide(NewTaskRepository),
)

type TaskRepository struct {
	conn   *sql.DB
	logger *logger.Logger
}

func NewTaskRepository(config *config.Config, logger *logger.Logger) *TaskRepository {
	conn := db.NewDB(config).GetConnection()
	return &TaskRepository{
		conn:   conn,
		logger: logger,
	}
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Task) (*task.Task, error) {
	query := `INSERT INTO tasks (id, title, description, status, priority, deadline, user_id, created_at, created_by)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var id uuid.UUID
	err := r.conn.QueryRowContext(ctx, query,
		t.Id, t.Title, t.Description, t.Status, t.Priority, t.Deadline, t.UserID, t.CreatedAt, t.CreatedBy,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	t.Id = id
	return t, nil
}

func (r *TaskRepository) UpdateById(ctx context.Context, id uuid.UUID, t *task.Task) (*task.Task, error) {
	query := `UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, deadline = $5, updated_at = $6, updated_by = $7
              WHERE id = $8`

	result, err := r.conn.ExecContext(ctx, query,
		t.Title, t.Description, t.Status, t.Priority, t.Deadline, t.UpdateAt, t.UpdateBy, id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	t.Id = id
	return t, nil
}

func (r *TaskRepository) FindById(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	query := `SELECT id, title, description, status, priority, deadline, user_id, created_at, created_by, updated_at, updated_by
              FROM tasks WHERE id = $1 AND deleted_at IS NULL`

	var t task.Task
	var deadline, updatedAt sql.NullTime
	var updatedBy sql.NullString

	err := r.conn.QueryRowContext(ctx, query, id).Scan(
		&t.Id, &t.Title, &t.Description, &t.Status, &t.Priority, &deadline, &t.UserID, &t.CreatedAt, &t.CreatedBy, &updatedAt, &updatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if deadline.Valid {
		t.Deadline = &deadline.Time
	}
	if updatedAt.Valid {
		t.UpdateAt = &updatedAt.Time
	}
	if updatedBy.Valid {
		uid, _ := uuid.Parse(updatedBy.String)
		t.UpdateBy = &uid
	}

	return &t, nil
}

func (r *TaskRepository) FindAll(ctx context.Context, filter *task.TaskFilter) ([]*task.Task, error) {
	query := `SELECT id, title, description, status, priority, deadline, user_id, created_at, created_by, updated_at, updated_by
              FROM tasks WHERE deleted_at IS NULL`

	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		conditions := []string{}

		if filter.Status != nil {
			conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, *filter.Status)
			argIndex++
		}

		if filter.Priority != nil {
			conditions = append(conditions, fmt.Sprintf("priority = $%d", argIndex))
			args = append(args, *filter.Priority)
			argIndex++
		}

		if filter.UserID != nil {
			conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
			args = append(args, *filter.UserID)
			argIndex++
		}

		if filter.Deadline != nil {
			conditions = append(conditions, fmt.Sprintf("deadline <= $%d", argIndex))
			args = append(args, *filter.Deadline)
			argIndex++
		}

		if len(conditions) > 0 {
			query += " AND " + strings.Join(conditions, " AND ")
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		var t task.Task
		var deadline, updatedAt sql.NullTime
		var updatedBy sql.NullString

		err := rows.Scan(
			&t.Id, &t.Title, &t.Description, &t.Status, &t.Priority, &deadline, &t.UserID, &t.CreatedAt, &t.CreatedBy, &updatedAt, &updatedBy,
		)
		if err != nil {
			return nil, err
		}

		if deadline.Valid {
			t.Deadline = &deadline.Time
		}
		if updatedAt.Valid {
			t.UpdateAt = &updatedAt.Time
		}
		if updatedBy.Valid {
			uid, _ := uuid.Parse(updatedBy.String)
			t.UpdateBy = &uid
		}

		tasks = append(tasks, &t)
	}

	return tasks, nil
}

func (r *TaskRepository) FindByUserId(ctx context.Context, userID uuid.UUID, filter *task.TaskFilter) ([]*task.Task, error) {
	if filter == nil {
		filter = &task.TaskFilter{}
	}
	filter.UserID = &userID
	return r.FindAll(ctx, filter)
}

func (r *TaskRepository) DeleteById(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE tasks SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2`

	result, err := r.conn.ExecContext(ctx, query, userID, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TaskRepository) MarkAsCompleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE tasks SET status = $1, updated_at = NOW(), updated_by = $2 WHERE id = $3`

	result, err := r.conn.ExecContext(ctx, query, task.StatusDone, userID, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TaskRepository) FindTasksNearDeadline(ctx context.Context, hoursAhead int) ([]*task.Task, error) {
	query := `SELECT id, title, description, status, priority, deadline, user_id, created_at, created_by, updated_at, updated_by
              FROM tasks
              WHERE deleted_at IS NULL
              AND status != $1
              AND deadline IS NOT NULL
              AND deadline <= NOW() + INTERVAL '1 hour' * $2
              ORDER BY deadline ASC`

	rows, err := r.conn.QueryContext(ctx, query, task.StatusDone, hoursAhead)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		var t task.Task
		var deadline, updatedAt sql.NullTime
		var updatedBy sql.NullString

		err := rows.Scan(
			&t.Id, &t.Title, &t.Description, &t.Status, &t.Priority, &deadline, &t.UserID, &t.CreatedAt, &t.CreatedBy, &updatedAt, &updatedBy,
		)
		if err != nil {
			return nil, err
		}

		if deadline.Valid {
			t.Deadline = &deadline.Time
		}
		if updatedAt.Valid {
			t.UpdateAt = &updatedAt.Time
		}
		if updatedBy.Valid {
			uid, _ := uuid.Parse(updatedBy.String)
			t.UpdateBy = &uid
		}

		tasks = append(tasks, &t)
	}

	return tasks, nil
}
