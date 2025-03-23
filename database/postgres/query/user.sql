-- name: GetUser :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = $2, password = $3 WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: DeleteUser :one
UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL RETURNING *;
