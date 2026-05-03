package repository

import "context"

// DynamicProductListParams contains runtime-built SQL and args for product list/count queries.
type DynamicProductListParams struct {
	ListSQL   string
	ListArgs  []any
	CountSQL  string
	CountArgs []any
}

// ListProductsByDynamicQuery executes dynamic list/count queries and maps results into Product models.
func (q *Queries) ListProductsByDynamicQuery(ctx context.Context, arg DynamicProductListParams) ([]*Product, int64, error) {
	rows, err := q.db.Query(ctx, arg.ListSQL, arg.ListArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]*Product, 0)
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Description,
			&product.Url,
			&product.CreatedAt,
			&product.CreatedUserID,
			&product.UpdatedAt,
			&product.UpdatedUserID,
			&product.DeletedAt,
			&product.DeletedUserID,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, &product)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	row := q.db.QueryRow(ctx, arg.CountSQL, arg.CountArgs...)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
