-- name: CreateUser :exec
INSERT INTO users (username, email, password) VALUES (?, ?, ?);

-- name: GetUserByID :one
SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?;

-- name: GetUserByUsername :one
SELECT id, username, email, created_at, updated_at FROM users WHERE username = ?;

-- name: GetUserByEmail :one
SELECT id, username, email, created_at, updated_at FROM users WHERE email = ?;

-- name: VerifyCredentials :one
SELECT id, username, email, created_at, updated_at FROM users WHERE username = ? AND password = ?;

-- name: UpdateUser :exec
UPDATE users
SET username = ?, email = ?, password = COALESCE(?, password), updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: ListUsers :many
SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC;