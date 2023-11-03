package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/klaital/library/storage/library"
	"github.com/klaital/library/storage/oldlibrary"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

func main() {
	ctx := context.Background()

	// Connect to the old DB
	oldDbPath := "old.db"
	oldDb, err := sql.Open("sqlite3", oldDbPath)
	if err != nil {
		slog.Error("Failed to connect to old DB", "error", err.Error())
		panic(err)
	}
	// Run DB migrations to ensure schema is up-to-date
	migrationsDir, err := iofs.New(library.Migrations, "migrations")
	if err != nil {
		slog.Error("Failed to load db migrations dir", "err", err)
		panic("failed to load db migrations")
	}
	m, err := migrate.NewWithSourceInstance("iofs", migrationsDir, fmt.Sprintf("sqlite3://%s", oldDbPath))
	if err != nil {
		slog.Error("Failed to prepare migration driver", "err", err)
		panic("failed to prepare migration driver")
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		slog.Debug("migrations not needed")
	} else if err != nil {
		slog.Error("Failed to execute db migrations", "err", err)
		panic("Failed to execute db migrations")
	} else {
		slog.Info("database migration successful")
	}
	oldStorer := oldlibrary.New(oldDb)

	// Connect to the new DB
	dbPath := "dev.db"
	newDb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("Failed to connect to new DB", "error", err.Error())
		panic(err)
	}
	// Run DB migrations to ensure schema is up-to-date
	migrationsDir, err = iofs.New(library.Migrations, "migrations")
	if err != nil {
		slog.Error("Failed to load db migrations dir", "err", err)
		panic("failed to load db migrations")
	}
	m, err = migrate.NewWithSourceInstance("iofs", migrationsDir, fmt.Sprintf("sqlite3://%s", dbPath))
	if err != nil {
		slog.Error("Failed to prepare migration driver", "err", err)
		panic("failed to prepare migration driver")
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		slog.Debug("migrations not needed")
	} else if err != nil {
		slog.Error("Failed to execute db migrations", "err", err)
		panic("Failed to execute db migrations")
	} else {
		slog.Info("database migration successful")
	}
	newStorer := library.New(newDb)

	// Examine new db.
	existingLocations, err := newStorer.GetLocations(ctx)
	slog.Info("Fetched existing locations", "count", len(existingLocations))
	for _, loc := range existingLocations {
		fmt.Printf("%d\t%s\n", len(loc.Items), loc.Name)
	}

	// Query all items from the old database
	oldItems, err := oldStorer.ListItems()
	if err != nil {
		slog.Error("Failed to fetch old data", "err", err)
		panic(err)
	}
	slog.Info("Fetched all old data", "count", len(oldItems))

	// Iterate once to create all of the Locations
	for _, archivalItem := range oldItems {
		locationName := "unknown"
		if archivalItem.Location.Valid {
			locationName = archivalItem.Location.String
		}
		id, err := newStorer.CreateLocation(ctx, library.Location{
			Name:  locationName,
			Notes: "archived at 6001, then moved.",
		})
		if err == nil {
			slog.Info("Created Location", "loc", locationName, "locId", id)
		}
	}

	// Cache the locations by name, since the old data uses a string ID
	existingLocations, err = newStorer.GetLocations(ctx)
	slog.Info("Fetched existing locations", "count", len(existingLocations))
	for _, loc := range existingLocations {
		fmt.Printf("%d\t%s\n", len(loc.Items), loc.Name)
	}
	locationsByName := make(map[string]library.Location, 0)
	for i, loc := range existingLocations {
		locationsByName[loc.Name] = existingLocations[i]
	}

	// Iterate the archive set a second time to migrate the item data
	for _, archivalItem := range oldItems {
		newItem := library.Item{
			ID:                  uint64(archivalItem.ID),
			LocationID:          0,
			Code:                archivalItem.Code,
			CodeType:            "",
			CodeSource:          "",
			Title:               "",
			TitleTranslated:     "",
			TitleTransliterated: "",
			CreatedAt:           archivalItem.Created.Time,
			UpdatedAt:           archivalItem.Created.Time,
		}
		if archivalItem.Location.Valid {
			newItem.LocationID = locationsByName[archivalItem.Location.String].ID
		} else {
			slog.Error("No valid location for archived item", "id", archivalItem.ID, "title", archivalItem.Title.String)
			continue
		}
		if archivalItem.CodeType.Valid {
			newItem.CodeType = archivalItem.CodeType.String
		}
		if archivalItem.DataSource.Valid {
			newItem.CodeSource = archivalItem.DataSource.String
		}
		if archivalItem.Title.Valid {
			newItem.Title = archivalItem.Title.String
		} else {
			slog.Error("No valid title for archived item", "id", archivalItem.ID, "code", archivalItem.Code)
			continue
		}
		if archivalItem.TitleTranslated.Valid {
			newItem.TitleTranslated = archivalItem.TitleTranslated.String
		}
		newId, err := newStorer.CreateItem(ctx, newItem.LocationID, newItem)
		if err != nil {
			slog.Error("Failed to create item", "err", err, "id", archivalItem.ID, "title", archivalItem.Title.String)
		} else {
			fmt.Printf("Imported #%d: %s (@%s)\n", newId, newItem.Title, archivalItem.Location.String)
		}
	}
}
