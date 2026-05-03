package repository

import "context"

// DynamicResourceListParams contains runtime-built SQL and args for resource list/count queries.
type DynamicResourceListParams struct {
	ListSQL   string
	ListArgs  []any
	CountSQL  string
	CountArgs []any
}

// ListResourcesByDynamicQuery executes dynamic list/count queries and maps results into Resource models.
func (q *Queries) ListResourcesByDynamicQuery(ctx context.Context, arg DynamicResourceListParams) ([]*Resource, int64, error) {
	rows, err := q.db.Query(ctx, arg.ListSQL, arg.ListArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]*Resource, 0)
	for rows.Next() {
		var item Resource
		if err := rows.Scan(
			&item.ResourceID,
			&item.Name,
			&item.ModuleID,
			&item.CreatedAt,
			&item.CreatedUserID,
			&item.UpdatedAt,
			&item.UpdatedUserID,
			&item.DeletedAt,
			&item.DeletedUserID,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, &item)
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
