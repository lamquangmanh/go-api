-- name: InsertUserRole :one
INSERT INTO user_roles (user_id, role_id, created_user_id, updated_user_id)
VALUES ($1, $2, $3, $4)
RETURNING user_role_id, user_id, role_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: GetUserRole :one
SELECT user_role_id, user_id, role_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM user_roles
WHERE user_role_id = $1 AND deleted_at IS NULL;

-- name: UpdateUserRole :one
UPDATE user_roles
SET user_id = $2, role_id = $3, updated_at = NOW(), updated_user_id = $4
WHERE user_role_id = $1 AND deleted_at IS NULL
RETURNING user_role_id, user_id, role_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id;

-- name: DeleteUserRole :one
UPDATE user_roles
SET deleted_at = NOW(), deleted_user_id = $2
WHERE user_role_id = $1 AND deleted_at IS NULL
RETURNING user_role_id;

-- name: ListUserRoles :many
SELECT user_role_id, user_id, role_id, created_at, created_user_id, updated_at, updated_user_id, deleted_at, deleted_user_id
FROM user_roles
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUserRoles :one
SELECT COUNT(*) FROM user_roles WHERE deleted_at IS NULL;

-- name: ListRolesByUser :many
SELECT r.role_id, r.name, r.description, r.module_id, r.created_at, r.created_user_id, r.updated_at, r.updated_user_id, r.deleted_at, r.deleted_user_id
FROM roles r
JOIN user_roles ur ON ur.role_id = r.role_id
WHERE ur.user_id = $1 AND ur.deleted_at IS NULL AND r.deleted_at IS NULL
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteUserRoleMappingsByUserID :exec
UPDATE user_roles SET deleted_at = NOW() WHERE user_id = $1 AND deleted_at IS NULL;
