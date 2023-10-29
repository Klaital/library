package library

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"github.com/klaital/library/storage/library/queries"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

//go:embed migrations/*.sql
var Migrations embed.FS

type Storer struct {
	db      *sql.DB
	queries *queries.Queries
}

func New(db *sql.DB) *Storer {
	return &Storer{
		db:      db,
		queries: queries.New(db),
	}
}

var ErrNotImplemented = errors.New("operation not implemented")

func (s *Storer) GetLocations(ctx context.Context) ([]Location, error) {
	l, err := s.queries.ListLocations(ctx)
	if err != nil {
		return nil, err
	}

	locations := make([]Location, len(l))
	for i, ql := range l {
		locations[i] = LocationFromQueries(ql)
		// Load the items for each location
		items, err := s.queries.ListItemsForLocation(ctx, ql.ID)
		if err != nil {
			slog.Error("failed to load items for location", "Location", ql.ID, "error", err.Error())
			continue
		} else {
			locations[i].Items = make([]Item, len(items))
			for j, qi := range items {
				locations[i].Items[j] = ItemFromQueries(qi)
			}
		}
	}

	return locations, nil
}

func (s *Storer) DescribeLocation(locationId uint64) (*Location, error) {
	return nil, ErrNotImplemented
}

func (s *Storer) CreateLocation(ctx context.Context, location Location) (newId uint64, err error) {
	params := queries.CreateLocationParams{
		Name: location.Name,
	}
	if len(location.Notes) > 0 {
		params.Notes = sql.NullString{
			String: location.Notes,
			Valid:  true,
		}
	}
	i, err := s.queries.CreateLocation(ctx, params)
	if err != nil {
		return 0, err
	}
	// Success!
	return uint64(i), nil
}

func (s *Storer) UpdateLocation(location Location) error {
	return ErrNotImplemented
}

func (s *Storer) DeleteLocation(locationId uint64) error {
	return ErrNotImplemented
}

func (s *Storer) CreateItem(ctx context.Context, locationId uint64, item Item) (newId uint64, err error) {
	params := queries.CreateItemParams{
		LocationID: int64(locationId),
		Code:       item.Code,
		CodeType:   item.CodeType,
		CodeSource: item.CodeSource,
		Title:      item.Title,
	}
	if len(item.TitleTranslated) > 0 {
		params.TitleTranslated = sql.NullString{
			String: item.TitleTranslated,
			Valid:  true,
		}
	}
	if len(item.TitleTransliterated) > 0 {
		params.TitleTransliterated = sql.NullString{
			String: item.TitleTransliterated,
			Valid:  true,
		}
	}
	id, err := s.queries.CreateItem(ctx, params)
	if err != nil {
		return 0, err
	}
	// Success!
	return uint64(id), nil
}

func (s *Storer) ListItemsForLocation(ctx context.Context, locationId uint64) ([]Item, error) {
	rows, err := s.queries.ListItemsForLocation(ctx, int64(locationId))
	if err != nil {
		return nil, err
	}
	items := make([]Item, len(rows))
	for i, row := range rows {
		items[i] = ItemFromQueries(row)
		items[i].LocationID = locationId
	}
	return items, nil
}

func (s *Storer) DescribeItem(itemId uint64) (*Item, error) {
	return nil, ErrNotImplemented
}

func (s *Storer) UpdateItem(item Item) error {
	return ErrNotImplemented
}

func (s *Storer) DeleteItem(itemId uint64) error {
	return ErrNotImplemented
}
