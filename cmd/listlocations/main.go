package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/klaital/library/storage/library"
	"log/slog"
	"os"
)

func main() {
	// Connect to DB
	dbPath := os.Getenv("DB_FILE")
	if dbPath == "" {
		dbPath = "dev.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("Failed to connect to db", "err", err, "path", dbPath)
		panic("failed to connect to db")
	}

	// Initialize the storage layer
	libraryStorer := library.New(db)

	locs, err := libraryStorer.GetLocations(context.Background())
	if err != nil {
		slog.Error("Failed to fetch locations from DB", "err", err)
		os.Exit(1)
	}
	if len(locs) == 0 {
		fmt.Printf("No locations found")
		os.Exit(0)
	}

	for _, l := range locs {
		fmt.Printf("------\nID: %d\nName: %s\nNotes: %s\nEntries: %d\n",
			l.ID, l.Name, l.Notes, len(l.Items))
	}
}
