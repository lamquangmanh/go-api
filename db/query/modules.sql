-- name: InsertModule :one
INSERT INTO modules (name, description, url, icon, product_id, created_user_id, updated_user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING module_id, name, description, url, icon, product_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: GetModule :one
SELECT module_id, name, description, url, icon, product_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM modules
WHERE module_id = $1 AND deleted_at IS NULL;

-- name: UpdateModule :one
UPDATE modules
SET name = $2, description = $3, url = $4, icon = $5, product_id = $6, updated_at = NOW(), updated_user_id = $7
WHERE module_id = $1 AND deleted_at IS NULL
RETURNING module_id, name, description, url, icon, product_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: DeleteModule :one
UPDATE modules
SET deleted_at = NOW(), deleted_user_id = $2
WHERE module_id = $1 AND deleted_at IS NULL
RETURNING module_id;

-- name: ListModules :many
SELECT module_id, name, description, url, icon, product_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM modules
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountModules :one
SELECT COUNT(*) FROM modules WHERE deleted_at IS NULL;
