package repository

import "context"

// DynamicUserListParams contains runtime-built SQL and args for user list/count queries.
type DynamicUserListParams struct {
	ListSQL   string
	ListArgs  []any
	CountSQL  string
	CountArgs []any
}

// ListUsersByDynamicQuery executes dynamic list/count queries and maps results into User models.
func (q *Queries) ListUsersByDynamicQuery(ctx context.Context, arg DynamicUserListParams) ([]*User, int64, error) {
	rows, err := q.db.Query(ctx, arg.ListSQL, arg.ListArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]*User, 0)
	for rows.Next() {
		var item User
		if err := rows.Scan(
			&item.UserID,
			&item.UserName,
			&item.Email,
			&item.Password,
			&item.Phone,
			&item.Avatar,
			&item.Status,
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
