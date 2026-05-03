-- name: InsertAction :one
INSERT INTO actions (resource_id, name, description, request_type, url, method, created_user_id, updated_user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING action_id, resource_id, name, description, request_type, url, method, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: GetAction :one
SELECT action_id, resource_id, name, description, request_type, url, method, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM actions
WHERE action_id = $1 AND deleted_at IS NULL;

-- name: UpdateAction :one
UPDATE actions
SET resource_id = $2, name = $3, description = $4, request_type = $5, url = $6, method = $7, updated_at = NOW(), updated_user_id = $8
WHERE action_id = $1 AND deleted_at IS NULL
RETURNING action_id, resource_id, name, description, request_type, url, method, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: DeleteAction :one
UPDATE actions
SET deleted_at = NOW(), deleted_user_id = $2
WHERE action_id = $1 AND deleted_at IS NULL
RETURNING action_id;

-- name: ListActions :many
SELECT action_id, resource_id, name, description, request_type, url, method, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM actions
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountActions :one
SELECT COUNT(*) FROM actions WHERE deleted_at IS NULL;

-- name: ListActionsByResource :many
SELECT action_id, resource_id, name, description, request_type, url, method, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM actions
WHERE resource_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
