-- name: InsertResource :one
INSERT INTO resources (name, module_id, created_user_id, updated_user_id)
VALUES ($1, $2, $3, $4)
RETURNING resource_id, name, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: GetResource :one
SELECT resource_id, name, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM resources
WHERE resource_id = $1 AND deleted_at IS NULL;

-- name: UpdateResource :one
UPDATE resources
SET name = $2, module_id = $3, updated_at = NOW(), updated_user_id = $4
WHERE resource_id = $1 AND deleted_at IS NULL
RETURNING resource_id, name, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: DeleteResource :one
UPDATE resources
SET deleted_at = NOW(), deleted_user_id = $2
WHERE resource_id = $1 AND deleted_at IS NULL
RETURNING resource_id;

-- name: ListResources :many
SELECT resource_id, name, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM resources
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountResources :one
SELECT COUNT(*) FROM resources WHERE deleted_at IS NULL;

-- name: ListResourcesByModule :many
SELECT resource_id, name, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM resources
WHERE module_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
