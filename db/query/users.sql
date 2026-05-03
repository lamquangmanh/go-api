-- name: InsertUser :one
INSERT INTO users (user_name, email, password, phone, avatar, status, created_user_id, updated_user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING user_id, user_name, email, password, phone, avatar, status, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: GetUser :one
SELECT user_id, user_name, email, password, phone, avatar, status, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM users
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT user_id, user_name, email, password, phone, avatar, status, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT user_id, user_name, email, password, phone, avatar, status, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM users
WHERE user_name = $1 AND deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE users
SET user_name = $2, email = $3, password = $4, phone = $5, avatar = $6, status = $7, updated_at = NOW(), updated_user_id = $8
WHERE user_id = $1 AND deleted_at IS NULL
RETURNING user_id, user_name, email, password, phone, avatar, status, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: DeleteUser :one
UPDATE users
SET deleted_at = NOW(), deleted_user_id = $2
WHERE user_id = $1 AND deleted_at IS NULL
RETURNING user_id;

-- name: ListUsers :many
SELECT user_id, user_name, email, password, phone, avatar, status, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;
