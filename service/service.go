package service

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/klaital/library/datasources/gbooks"
	"github.com/klaital/library/datasources/upcdatabasedotorg"
	"github.com/klaital/library/storage/library"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//go:embed templates/*.html
var htmlTemplateData embed.FS

type Service struct {
	LibraryStorage    *library.Storer
	htmlTemplates     *template.Template
	GoogleBooksClient *gbooks.Client
	UpcDatabaseClient *upcdatabasedotorg.Client
}

func New(storer *library.Storer) *Service {
	tmpl, err := template.ParseFS(htmlTemplateData, "templates/*.html")
	if err != nil {
		slog.Error("Failed to parse html templates", "error", err.Error())
		panic("failed to parse html templates")
	}
	svc := &Service{
		LibraryStorage:    storer,
		htmlTemplates:     tmpl,
		GoogleBooksClient: gbooks.New(""),
	}
	apiKey := os.Getenv("UPCDATABASEDOTORG_KEY")
	if apiKey == "" {
		slog.Error("No API key given for upcdatabase.org. Expected to be set in env var UPCDATABASEDOTORG_KEY")
	} else {
		svc.UpcDatabaseClient = upcdatabasedotorg.New(apiKey)
	}

	return svc
}

func readJson(r *http.Request, target any) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	err = json.Unmarshal(b, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	// Success!
	return nil
}

func (svc *Service) HandleListLocations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//params := httprouter.ParamsFromContext(r.Context())
	locations, err := svc.LibraryStorage.GetLocations(r.Context())
	if err != nil {
		slog.Error("failed to fetch locations", "err", err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := json.Marshal(locations)
	if err != nil {
		slog.Error("failed to marshal locations to JSON", "error", err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

func (svc *Service) HandleCreateLocation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var loc library.Location
	err := readJson(r, &loc)
	if err != nil {
		slog.Error("read request body", "err", err)
		w.WriteHeader(400)
		return
	}

	newId, err := svc.LibraryStorage.CreateLocation(r.Context(), loc)
	if err != nil {
		slog.Debug("Failed to create location", "err", err)
		w.WriteHeader(500)
		return
	}
	loc.ID = newId

	b, err := json.Marshal(loc)
	if err != nil {
		slog.Error("Failed to marshal response", "err", err)
		w.WriteHeader(500)
		return
	}
	w.Write(b)
	w.WriteHeader(200)
}

func (svc *Service) HandleCreateItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var item library.Item
	err := readJson(r, &item)
	if err != nil {
		slog.Debug("read request body", "err", err)
		w.WriteHeader(400)
		return
	}
	locIdRaw := params.ByName("locationId")
	locId, err := strconv.ParseUint(locIdRaw, 10, 64)
	if err != nil {
		slog.Debug("invalid location id", "err", err)
		w.WriteHeader(404)
		return
	}
	_, err = svc.LibraryStorage.CreateItem(r.Context(), locId, item)
	if err != nil {
		slog.Debug("failed to create item", "err", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func (svc *Service) HandleGetItemsForLocation(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//params := httprouter.ParamsFromContext(r.Context())
	locIdRaw := params.ByName("locationId")
	locId, err := strconv.ParseUint(locIdRaw, 10, 64)
	if err != nil {
		slog.Debug("invalid location id", "err", err)
		w.WriteHeader(404)
		return
	}

	items, err := svc.LibraryStorage.ListItemsForLocation(r.Context(), locId)
	if err != nil {
		slog.Error("failed to fetch items", "err", err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := json.Marshal(items)
	if err != nil {
		slog.Error("failed to marshal locations to JSON", "error", err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

type CodeLookupRequest struct {
	Code string
	Type string
}

func (svc *Service) HandleCodeLookup(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	codeType := params.ByName("type")
	if len(codeType) == 0 {
		w.Write([]byte("no code type specified"))
		w.WriteHeader(400)
		return
	}

	code := params.ByName("code")
	if len(code) == 0 {
		w.Write([]byte("no code specified"))
		w.WriteHeader(400)
		return
	}

	if strings.ToLower(codeType) == "isbn" {
		itemData, err := svc.GoogleBooksClient.LookupIsbn(r.Context(), code)
		if err != nil {
			slog.Error("Failed to look up ISBN", "err", err, "ISBN", code)
			w.WriteHeader(500)
			return
		}
		b, err := json.Marshal(itemData)
		if err != nil {
			slog.Error("Failed to marshal response", "err", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(b)
		return
	} else if strings.ToLower(codeType) == "upc" {
		if svc.UpcDatabaseClient == nil {
			slog.Error("No upcdatabase.com client configured. Unable to look up UPC data")
			w.WriteHeader(500)
			return
		}
		itemData, err := svc.UpcDatabaseClient.LookupUpc(code)
		if err != nil {
			slog.Error("Failed to look up UPC", "err", err, "UPC", code)
			w.WriteHeader(500)
			return
		}
		b, err := json.Marshal(itemData)
		if err != nil {
			slog.Error("Failed to marshal response", "err", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(b)
		return

	} else {
		slog.Debug("Unknown code type", "type", codeType, "code", code)
		w.WriteHeader(400)
		w.Write([]byte("Unknown Code Type"))
		return
	}
}

func (svc *Service) HandleMoveItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	newLocation := params.ByName("locId")
	itemIdRaw := params.ByName("itemId")
	if len(newLocation) == 0 {
		w.WriteHeader(400)
		return
	}
	if len(itemIdRaw) == 0 {
		w.WriteHeader(400)
		return
	}

	itemId, err := strconv.ParseInt(itemIdRaw, 10, 0)
	if err != nil {
		slog.Error("Failed to parse Item ID", "rawId", itemIdRaw, "err", err)
		w.WriteHeader(400)
		return
	}
	locationId, err := strconv.ParseInt(newLocation, 10, 0)
	if err != nil {
		slog.Error("Failed to parse Location ID", "rawId", newLocation, "err", err)
		w.WriteHeader(400)
		return
	}

	err = svc.LibraryStorage.MoveItem(r.Context(), itemId, locationId)
	if err != nil {
		slog.Error("Failed to move item", "err", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
