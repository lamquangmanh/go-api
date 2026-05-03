package repository

import "context"

// DynamicModuleListParams contains runtime-built SQL and args for module list/count queries.
type DynamicModuleListParams struct {
	ListSQL   string
	ListArgs  []any
	CountSQL  string
	CountArgs []any
}

// ListModulesByDynamicQuery executes dynamic list/count queries and maps results into Module models.
func (q *Queries) ListModulesByDynamicQuery(ctx context.Context, arg DynamicModuleListParams) ([]*Module, int64, error) {
	rows, err := q.db.Query(ctx, arg.ListSQL, arg.ListArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]*Module, 0)
	for rows.Next() {
		var item Module
		if err := rows.Scan(
			&item.ModuleID,
			&item.Name,
			&item.Description,
			&item.Url,
			&item.Icon,
			&item.ProductID,
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
