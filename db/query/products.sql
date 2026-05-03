-- name: InsertProduct :one
INSERT INTO products (name, description, url, created_user_id, updated_user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProduct :one
SELECT *
FROM products
WHERE product_id = $1 AND deleted_at IS NULL;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, url = $4, updated_at = NOW(), updated_user_id = $5
WHERE product_id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProduct :one
UPDATE products
SET deleted_at = NOW(), deleted_user_id = $2
WHERE product_id = $1 AND deleted_at IS NULL
RETURNING product_id;

-- name: ListProducts :many
SELECT *
FROM products
WHERE
	deleted_at IS NULL
	AND (sqlc.arg(name)::VARCHAR = '' OR name ILIKE CONCAT('%', sqlc.arg(name)::VARCHAR, '%'))
	AND (sqlc.arg(description)::VARCHAR = '' OR description ILIKE CONCAT('%', sqlc.arg(description)::VARCHAR, '%'))
	AND (COALESCE(array_length(sqlc.arg(product_ids)::UUID[], 1), 0) = 0 OR product_id = ANY(sqlc.arg(product_ids)::UUID[]))
ORDER BY
	CASE WHEN sqlc.arg(sort_field_1)::TEXT = 'product_id' AND sqlc.arg(sort_order_1)::TEXT = 'ASC' THEN product_id END ASC,
	CASE WHEN sqlc.arg(sort_field_1)::TEXT = 'product_id' AND sqlc.arg(sort_order_1)::TEXT = 'DESC' THEN product_id END DESC,
	CASE WHEN sqlc.arg(sort_field_1)::TEXT = 'name' AND sqlc.arg(sort_order_1)::TEXT = 'ASC' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_field_1)::TEXT = 'name' AND sqlc.arg(sort_order_1)::TEXT = 'DESC' THEN name END DESC,
	CASE WHEN sqlc.arg(sort_field_1)::TEXT = 'created_at' AND sqlc.arg(sort_order_1)::TEXT = 'ASC' THEN created_at END ASC,
	CASE WHEN sqlc.arg(sort_field_1)::TEXT = 'created_at' AND sqlc.arg(sort_order_1)::TEXT = 'DESC' THEN created_at END DESC,
	CASE WHEN sqlc.arg(sort_field_1)::TEXT = 'updated_at' AND sqlc.arg(sort_order_1)::TEXT = 'ASC' THEN updated_at END ASC,
	CASE WHEN sqlc.arg(sort_field_1)::TEXT = 'updated_at' AND sqlc.arg(sort_order_1)::TEXT = 'DESC' THEN updated_at END DESC,

	CASE WHEN sqlc.arg(sort_field_2)::TEXT = 'product_id' AND sqlc.arg(sort_order_2)::TEXT = 'ASC' THEN product_id END ASC,
	CASE WHEN sqlc.arg(sort_field_2)::TEXT = 'product_id' AND sqlc.arg(sort_order_2)::TEXT = 'DESC' THEN product_id END DESC,
	CASE WHEN sqlc.arg(sort_field_2)::TEXT = 'name' AND sqlc.arg(sort_order_2)::TEXT = 'ASC' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_field_2)::TEXT = 'name' AND sqlc.arg(sort_order_2)::TEXT = 'DESC' THEN name END DESC,
	CASE WHEN sqlc.arg(sort_field_2)::TEXT = 'created_at' AND sqlc.arg(sort_order_2)::TEXT = 'ASC' THEN created_at END ASC,
	CASE WHEN sqlc.arg(sort_field_2)::TEXT = 'created_at' AND sqlc.arg(sort_order_2)::TEXT = 'DESC' THEN created_at END DESC,
	CASE WHEN sqlc.arg(sort_field_2)::TEXT = 'updated_at' AND sqlc.arg(sort_order_2)::TEXT = 'ASC' THEN updated_at END ASC,
	CASE WHEN sqlc.arg(sort_field_2)::TEXT = 'updated_at' AND sqlc.arg(sort_order_2)::TEXT = 'DESC' THEN updated_at END DESC,

	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'product_id' AND sqlc.arg(sort_order_3)::TEXT = 'ASC' THEN product_id END ASC,
	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'product_id' AND sqlc.arg(sort_order_3)::TEXT = 'DESC' THEN product_id END DESC,
	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'name' AND sqlc.arg(sort_order_3)::TEXT = 'ASC' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'name' AND sqlc.arg(sort_order_3)::TEXT = 'DESC' THEN name END DESC,
	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'created_at' AND sqlc.arg(sort_order_3)::TEXT = 'ASC' THEN created_at END ASC,
	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'created_at' AND sqlc.arg(sort_order_3)::TEXT = 'DESC' THEN created_at END DESC,
	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'updated_at' AND sqlc.arg(sort_order_3)::TEXT = 'ASC' THEN updated_at END ASC,
	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'updated_at' AND sqlc.arg(sort_order_3)::TEXT = 'DESC' THEN updated_at END DESC,
	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'url' AND sqlc.arg(sort_order_3)::TEXT = 'ASC' THEN url END ASC,
	CASE WHEN sqlc.arg(sort_field_3)::TEXT = 'url' AND sqlc.arg(sort_order_3)::TEXT = 'DESC' THEN url END DESC,
	created_at DESC
LIMIT sqlc.arg(p_limit) OFFSET sqlc.arg(p_offset);

-- name: CountProducts :one
SELECT COUNT(*)
FROM products
WHERE
	deleted_at IS NULL
	AND (sqlc.arg(name)::VARCHAR = '' OR name ILIKE CONCAT('%', sqlc.arg(name)::VARCHAR, '%'))
	AND (sqlc.arg(description)::VARCHAR = '' OR description ILIKE CONCAT('%', sqlc.arg(description)::VARCHAR, '%'))
	AND (COALESCE(array_length(sqlc.arg(product_ids)::UUID[], 1), 0) = 0 OR product_id = ANY(sqlc.arg(product_ids)::UUID[]));
