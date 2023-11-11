package library

import (
	"github.com/klaital/library/storage/library/queries"
	"time"
)

type Item struct {
	ID                  uint64
	LocationID          uint64
	Code                string
	CodeType            string
	CodeSource          string
	Title               string
	TitleTranslated     string
	TitleTransliterated string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Location struct {
	ID    uint64
	Name  string
	Notes string
	Items []Item
}

func LocationFromQueries(ql queries.Location) Location {
	var l Location
	l.ID = uint64(ql.ID)
	if ql.Notes.Valid {
		l.Notes = ql.Notes.String
	}
	l.Name = ql.Name
	return l
}

func ItemFromQueries(qi queries.ListItemsForLocationRow) Item {
	var item Item
	item.ID = uint64(qi.ID)
	item.Code = qi.Code
	item.CodeType = qi.CodeType
	item.CodeSource = qi.CodeSource
	item.Title = qi.Title
	if qi.TitleTranslated.Valid {
		item.TitleTranslated = qi.TitleTranslated.String
	}
	if qi.TitleTransliterated.Valid {
		item.TitleTransliterated = qi.TitleTransliterated.String
	}
	item.CreatedAt = qi.CreatedAt
	item.UpdatedAt = qi.UpdatedAt
	return item
}

func ItemFromGetItemRow(qi queries.GetItemRow) Item {
	var item Item
	item.ID = uint64(qi.ID)
	item.Code = qi.Code
	item.CodeType = qi.CodeType
	item.CodeSource = qi.CodeSource
	item.Title = qi.Title
	if qi.TitleTranslated.Valid {
		item.TitleTranslated = qi.TitleTranslated.String
	}
	if qi.TitleTransliterated.Valid {
		item.TitleTransliterated = qi.TitleTransliterated.String
	}
	item.CreatedAt = qi.CreatedAt
	item.UpdatedAt = qi.UpdatedAt
	return item
}
