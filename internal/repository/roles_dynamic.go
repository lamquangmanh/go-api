package repository

import "context"

// DynamicRoleListParams contains runtime-built SQL and args for role list/count queries.
type DynamicRoleListParams struct {
	ListSQL   string
	ListArgs  []any
	CountSQL  string
	CountArgs []any
}

// ListRolesByDynamicQuery executes dynamic list/count queries and maps results into Role models.
func (q *Queries) ListRolesByDynamicQuery(ctx context.Context, arg DynamicRoleListParams) ([]*Role, int64, error) {
	rows, err := q.db.Query(ctx, arg.ListSQL, arg.ListArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]*Role, 0)
	for rows.Next() {
		var item Role
		if err := rows.Scan(
			&item.RoleID,
			&item.Name,
			&item.Description,
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
