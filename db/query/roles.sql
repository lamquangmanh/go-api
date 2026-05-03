-- name: InsertRole :one
INSERT INTO roles (name, description, module_id, created_user_id, updated_user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING role_id, name, description, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: GetRole :one
SELECT role_id, name, description, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM roles
WHERE role_id = $1 AND deleted_at IS NULL;

-- name: GetRoleByName :one
SELECT role_id, name, description, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM roles
WHERE name = $1 AND deleted_at IS NULL;

-- name: GetRolesByIDs :many
SELECT role_id FROM roles WHERE role_id = ANY($1::uuid[]) AND deleted_at IS NULL;

-- name: UpdateRole :one
UPDATE roles
SET name = $2, description = $3, module_id = $4, updated_at = NOW(), updated_user_id = $5
WHERE role_id = $1 AND deleted_at IS NULL
RETURNING role_id, name, description, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: DeleteRole :one
UPDATE roles
SET deleted_at = NOW(), deleted_user_id = $2
WHERE role_id = $1 AND deleted_at IS NULL
RETURNING role_id;

-- name: ListRoles :many
SELECT role_id, name, description, module_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM roles
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountRoles :one
SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL;
