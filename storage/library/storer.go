package library

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
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

func (s *Storer) GetAllItems(ctx context.Context) ([]Item, error) {
	// TODO: optimize this into a single query. sqlc is currently generating a separate row struct for ListAllItems, which would require maintaining another mapper.
	l, err := s.GetLocations(ctx)
	if err != nil {
		return nil, err
	}

	// Flatten all of the items into a single slice
	items := make([]Item, 0)
	for _, loc := range l {
		for j := range loc.Items {
			items = append(items, loc.Items[j])
		}
	}

	return items, nil
}

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

func (s *Storer) DescribeItem(ctx context.Context, itemId uint64) (*Item, error) {
	row, err := s.queries.GetItem(ctx, int64(itemId))
	if err != nil {
		return nil, fmt.Errorf("fetching item data: %w", err)
	}
	i := ItemFromGetItemRow(row)
	return &i, nil
}

func (s *Storer) UpdateItem(item Item) error {
	return ErrNotImplemented
}

func (s *Storer) DeleteItem(itemId uint64) error {
	return ErrNotImplemented
}

func (s *Storer) MoveItem(ctx context.Context, itemId int64, locationId int64) error {
	err := s.queries.MoveItem(ctx, queries.MoveItemParams{
		LocationID: locationId,
		ID:         itemId,
	})
	if err != nil {
		return fmt.Errorf("updating item: %w", err)
	}
	return nil
}

func (s *Storer) BulkItems(ctx context.Context, codeType string, codes []string) (map[string][]*Item, error) {
	rows, err := s.queries.GetItems(ctx, queries.GetItemsParams{
		CodeType: codeType,
		Codes:    codes,
	})
	resp := make(map[string][]*Item, 0)
	if err == sql.ErrNoRows {
		return resp, nil
	}
	if err != nil {
		slog.Error("Failed to query for items", "err", err)
		return nil, fmt.Errorf("fetching bulk items: %w", err)
	}

	for _, row := range rows {
		itemsForCode := resp[row.Code]
		if itemsForCode == nil {
			itemsForCode = make([]*Item, 0)
		}
		itemsForCode = append(itemsForCode, &Item{
			ID:                  uint64(row.ID),
			LocationID:          uint64(row.LocationID),
			Code:                row.Code,
			CodeType:            row.CodeType,
			CodeSource:          row.CodeSource,
			Title:               row.Title,
			TitleTranslated:     row.TitleTranslated.String,
			TitleTransliterated: row.TitleTransliterated.String,
			CreatedAt:           row.CreatedAt,
			UpdatedAt:           row.UpdatedAt,
		})

		// Update the storage map
		resp[row.Code] = itemsForCode
	}

	// Success!
	return resp, nil
}
