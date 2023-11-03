package oldlibrary

import (
	"context"
	"database/sql"
	"embed"
	oldqueries "github.com/klaital/library/storage/oldlibrary/queries"
)

//go:embed migrations/*.sql
var Migrations embed.FS

type Storer struct {
	db      *sql.DB
	queries *oldqueries.Queries
}

func New(db *sql.DB) *Storer {
	return &Storer{
		db:      db,
		queries: oldqueries.New(db),
	}
}

func (s *Storer) ListItems() ([]oldqueries.Items, error) {
	return s.queries.ListItems(context.Background())
}
