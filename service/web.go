package service

import (
	"github.com/julienschmidt/httprouter"
	"github.com/klaital/library/storage/library"
	"log/slog"
	"net/http"
)

func (svc *Service) WebListLocations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	viewData := struct {
		Locations []library.Location
	}{}
	//params := httprouter.ParamsFromContext(r.Context())
	locations, err := svc.LibraryStorage.GetLocations(r.Context())
	if err != nil {
		slog.Error("failed to fetch locations", "err", err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	viewData.Locations = locations

	// Validate the titles' encoding

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.Header().Set("Content-Encoding", "utf-8")
	err = svc.htmlTemplates.ExecuteTemplate(w, "ListLocations.html", viewData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to render template", "error", err.Error())
	}
}

func (svc *Service) WebAllItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	viewData := struct {
		Items []library.Item
	}{}
	//params := httprouter.ParamsFromContext(r.Context())
	items, err := svc.LibraryStorage.GetAllItems(r.Context())
	if err != nil {
		slog.Error("failed to fetch items", "err", err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	viewData.Items = items

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.Header().Set("Content-Encoding", "utf-8")
	err = svc.htmlTemplates.ExecuteTemplate(w, "ListItems.html", viewData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to render template", "error", err.Error())
	}
}
