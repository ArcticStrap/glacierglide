package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"slices"

	"github.com/ArcticStrap/glacierglide/polarpages/models"
	"github.com/go-chi/chi/v5"
)

func SetupHistoryHandler(rt *chi.Mux, addr string) {
	// Initalize templates
	tmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/history.html", "polarpages/templates/pagenav.html"))

	// Source routing
	rt.Get("/h/{title}", func(w http.ResponseWriter, r *http.Request) {
		titleParam := chi.URLParam(r, "title")

		// Get page history
		res, err := http.Get(addr + "/api/h/" + titleParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		// Read response
		var PageHistory []models.PageEdit

		err = json.NewDecoder(res.Body).Decode(&PageHistory)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Reverse order to show latest at top
		slices.Reverse(PageHistory)

		// Parse template
		tmpl.Execute(w, struct {
			models.SessionData
			models.WebPage
			History []models.PageEdit
			models.WebModes
		}{UserSession(r), models.WebPage{
			Title: titleParam,
			Theme: "common",
		}, PageHistory, models.WebModes{PageMode: "history"}})
	})
}
