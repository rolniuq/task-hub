-- name: CreateUser :one
INSERT INTO users (
    id, name, email, password, created_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET
    name = $1,
    email = $2,
    updated_at = $3
WHERE id = $4
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1;
