package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/julienschmidt/httprouter"
	"github.com/klaital/library/service"
	"github.com/klaital/library/storage/library"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	var err error
	var libraryStorer *library.Storer

	loggerOptions := &slog.HandlerOptions{Level: slog.LevelDebug}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, loggerOptions)))

	//
	// Prepare the DB
	//
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
	slog.Info("established connection to database", "db", dbPath)

	// Run DB migrations to ensure schema is up-to-date
	migrationsDir, err := iofs.New(library.Migrations, "migrations")
	if err != nil {
		slog.Error("Failed to load db migrations dir", "err", err)
		panic("failed to load db migrations")
	}

	//driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	//if err != nil {
	//	slog.Error("Failed to prepare migration driver", "err", err)
	//	panic("failed to prepare migration driver")
	//}
	m, err := migrate.NewWithSourceInstance("iofs", migrationsDir, fmt.Sprintf("sqlite3://%s", dbPath))
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

	// Initialize the storage layer
	libraryStorer = library.New(db)

	// Initialize the service layer
	svc := service.New(libraryStorer)

	//
	// Prepare the HTTP server
	//
	router := httprouter.New()

	// JSON APIs
	router.GET("/api/locations", svc.HandleListLocations)
	router.POST("/api/locations", svc.HandleCreateLocation)
	router.GET("/api/locations/:locationId/items", svc.HandleGetItemsForLocation)
	router.POST("/api/locations/:locationId/items", svc.HandleCreateItem)
	router.PUT("/api/items/:itemId/relocate/:locationId", svc.HandleMoveItem)
	router.GET("/api/code", svc.HandleCodeLookup)

	// Web UI
	router.ServeFiles("/web/*filepath", http.Dir("web/build"))

	slog.Info("Listening for HTTP requests on :8080")

	corsCfg := cors.New(cors.Options{
		AllowedOrigins:             nil,
		AllowOriginFunc:            nil,
		AllowOriginRequestFunc:     nil,
		AllowOriginVaryRequestFunc: nil,
		AllowedMethods:             []string{"GET", "POST", "HEAD"},
		AllowedHeaders:             []string{"Accept", "Content-Type"},
		ExposedHeaders:             nil,
		MaxAge:                     0,
		AllowCredentials:           false,
		AllowPrivateNetwork:        false,
		OptionsPassthrough:         false,
		OptionsSuccessStatus:       0,
		Debug:                      false,
	})

	corsHandler := corsCfg.Handler(router)
	http.ListenAndServe(":8080", corsHandler)
}
