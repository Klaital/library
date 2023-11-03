// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: items.sql

package queries

import (
	"context"
	"database/sql"
	"time"
)

const createItem = `-- name: CreateItem :one
INSERT INTO items (location_id, code, code_type, code_source,
                   title, title_translated, title_transliterated)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id
`

type CreateItemParams struct {
	LocationID          int64
	Code                string
	CodeType            string
	CodeSource          string
	Title               string
	TitleTranslated     sql.NullString
	TitleTransliterated sql.NullString
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createItem,
		arg.LocationID,
		arg.Code,
		arg.CodeType,
		arg.CodeSource,
		arg.Title,
		arg.TitleTranslated,
		arg.TitleTransliterated,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const listItemsForLocation = `-- name: ListItemsForLocation :many
SELECT id, code, code_type, code_source,
       title, title_translated, title_transliterated,
       created_at, updated_at
FROM items
WHERE location_id=?
`

type ListItemsForLocationRow struct {
	ID                  int64
	Code                string
	CodeType            string
	CodeSource          string
	Title               string
	TitleTranslated     sql.NullString
	TitleTransliterated sql.NullString
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (q *Queries) ListItemsForLocation(ctx context.Context, locationID int64) ([]ListItemsForLocationRow, error) {
	rows, err := q.db.QueryContext(ctx, listItemsForLocation, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListItemsForLocationRow
	for rows.Next() {
		var i ListItemsForLocationRow
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.CodeType,
			&i.CodeSource,
			&i.Title,
			&i.TitleTranslated,
			&i.TitleTransliterated,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}