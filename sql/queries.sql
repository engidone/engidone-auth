-- queries.sql

-- name: ListCredentials :many
SELECT id, user_id, password, created_at, updated_at FROM credentials;

-- name: GetCredential :one
SELECT id, user_id, password, created_at, updated_at FROM credentials WHERE user_id = $1 AND password = $2;

-- name: InsertCredential :one
INSERT INTO credentials (user_id, password) VALUES ($1, $2) RETURNING id, user_id, password, created_at, updated_at;

-- name: UpdatePassword :exec
UPDATE credentials
SET password = $1, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $2;

-- name: GetRecoveryCode :one
SELECT code FROM recovery_codes WHERE code = $1 AND is_valid = TRUE AND expires_at >= CURRENT_TIMESTAMP;

-- name: InsertRecoveryCode :one
INSERT INTO recovery_codes (user_id, code, is_valid, expires_at) VALUES ($1, $2, $3, $4) RETURNING id, user_id, code, is_valid, expires_at, created_at;

-- name: GetUserByRefreshToken :one
SELECT user_id
  FROM refresh_tokens
  WHERE refresh_token = $1;

-- name: InsertOrUpdateRefreshToken :one
INSERT INTO refresh_tokens (user_id, refresh_token)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE SET
    refresh_token = EXCLUDED.refresh_token
RETURNING *;