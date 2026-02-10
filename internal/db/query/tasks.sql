-- name: CreateTask :one
INSERT INTO tasks (
    id, title, description, status, priority, deadline, user_id, created_at, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateTask :one
UPDATE tasks SET
    title = $1,
    description = $2,
    status = $3,
    priority = $4,
    deadline = $5,
    updated_at = $6,
    updated_by = $7
WHERE id = $8 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteTask :execrows
UPDATE tasks SET
    deleted_at = NOW(),
    deleted_by = $1
WHERE id = $2 AND deleted_at IS NULL;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListTasksByUser :many
SELECT * FROM tasks
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListTasksByStatus :many
SELECT * FROM tasks
WHERE status = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListTasksByPriority :many
SELECT * FROM tasks
WHERE priority = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: MarkTaskAsCompleted :execrows
UPDATE tasks SET
    status = 'done',
    updated_at = NOW(),
    updated_by = $1
WHERE id = $2 AND deleted_at IS NULL;

-- name: GetTasksNearDeadline :many
SELECT * FROM tasks
WHERE deleted_at IS NULL
AND status != 'done'
AND deadline IS NOT NULL
AND deadline <= NOW() + INTERVAL '1 hour' * $1
ORDER BY deadline ASC;
