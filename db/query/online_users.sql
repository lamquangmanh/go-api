-- name: InsertOnlineUser :one
INSERT INTO online_users (user_id, socket_id, device_info, current_page_url, created_user_id, updated_user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING online_user_id, user_id, socket_id, device_info, current_page_url, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: GetOnlineUser :one
SELECT online_user_id, user_id, socket_id, device_info, current_page_url, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM online_users
WHERE online_user_id = $1 AND deleted_at IS NULL;

-- name: DeleteOnlineUser :one
UPDATE online_users
SET deleted_at = NOW(), deleted_user_id = $2
WHERE online_user_id = $1 AND deleted_at IS NULL
RETURNING online_user_id;

-- name: DeleteOnlineUserBySocket :exec
DELETE FROM online_users WHERE user_id = $1 AND socket_id = $2;

-- name: ListOnlineUsersByUser :many
SELECT online_user_id, user_id, socket_id, device_info, current_page_url, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM online_users
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountOnlineUsersByUser :one
SELECT COUNT(*) FROM online_users WHERE user_id = $1 AND deleted_at IS NULL;
