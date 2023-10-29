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

	w.Header().Set("Content-Type", "text/html")
	err = svc.htmlTemplates.ExecuteTemplate(w, "ListLocations.html", viewData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("failed to render template", "error", err.Error())
	}
}
