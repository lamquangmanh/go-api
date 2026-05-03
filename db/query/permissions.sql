-- name: InsertPermission :one
INSERT INTO permissions (role_id, resource_id, action_id, created_user_id, updated_user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING permission_id, role_id, resource_id, action_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: GetPermission :one
SELECT permission_id, role_id, resource_id, action_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM permissions
WHERE permission_id = $1 AND deleted_at IS NULL;

-- name: UpdatePermission :one
UPDATE permissions
SET role_id = $2, resource_id = $3, action_id = $4, updated_at = NOW(), updated_user_id = $5
WHERE permission_id = $1 AND deleted_at IS NULL
RETURNING permission_id, role_id, resource_id, action_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: DeletePermission :one
UPDATE permissions
SET deleted_at = NOW(), deleted_user_id = $2
WHERE permission_id = $1 AND deleted_at IS NULL
RETURNING permission_id;

-- name: ListPermissions :many
SELECT permission_id, role_id, resource_id, action_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM permissions
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPermissions :one
SELECT COUNT(*) FROM permissions WHERE deleted_at IS NULL;

-- name: ListPermissionsByRole :many
SELECT permission_id, role_id, resource_id, action_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM permissions
WHERE role_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
