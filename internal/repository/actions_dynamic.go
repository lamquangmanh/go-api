package repository

import "context"

// DynamicActionListParams contains runtime-built SQL and args for action list/count queries.
type DynamicActionListParams struct {
	ListSQL   string
	ListArgs  []any
	CountSQL  string
	CountArgs []any
}

// ListActionsByDynamicQuery executes dynamic list/count queries and maps results into Action models.
func (q *Queries) ListActionsByDynamicQuery(ctx context.Context, arg DynamicActionListParams) ([]*Action, int64, error) {
	rows, err := q.db.Query(ctx, arg.ListSQL, arg.ListArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]*Action, 0)
	for rows.Next() {
		var item Action
		if err := rows.Scan(
			&item.ActionID,
			&item.ResourceID,
			&item.Name,
			&item.Description,
			&item.RequestType,
			&item.Url,
			&item.Method,
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
