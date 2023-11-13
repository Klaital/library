/*
 * Use remote data sources to load data for a UPC or ISBN, then save the item to the database.
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/klaital/library/datasources/gbooks"
	"github.com/klaital/library/datasources/upcdatabasedotorg"
	"github.com/klaital/library/storage/library"
	"os"
)

type config struct {
	locationId int64
	code       string
	codeType   string
}

func main() {
	cfg := config{}
	flag.Int64Var(&cfg.locationId, "l", 0, "Location ID. Must already exist.")
	flag.StringVar(&cfg.code, "code", "", "Code to look up.")
	flag.StringVar(&cfg.codeType, "type", "ISBN", "Type of code to look up.")
	flag.Parse()

	// Validate input
	if cfg.locationId == 0 {
		fmt.Printf("Location ID is required parameter\n")
		os.Exit(1)
	}

	var err error
	var item *library.Item
	if cfg.codeType == "ISBN" {
		isbnClient := gbooks.New("") // TODO: do I need an API key if I'm just querying ISBN data, not managing user bookshelves?
		item, err = isbnClient.LookupIsbn(context.Background(), cfg.code)
		if err != nil {
			fmt.Printf("Failed to look up ISBN: %s\n", err.Error())
			os.Exit(1)
		}
	} else if cfg.codeType == "UPC" {
		apiKey := os.Getenv("UPCDATABASEDOTORG_KEY")
		if apiKey == "" {
			fmt.Printf("Need API key for upcdatabase.org. Set as env var UPCDATABASEDOTORG_KEY")
			os.Exit(1)
		}
		upcClient := upcdatabasedotorg.New(apiKey)
		item, err = upcClient.LookupUpc(cfg.code)
		if err != nil {
			fmt.Printf("Failed to look up UPC: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Printf("Error: Unknown code type %s\n", cfg.codeType)
		os.Exit(1)
	}

	fmt.Printf("Title: %s\n", item.Title)

	// TODO: actually write a new entry to the DB

}
