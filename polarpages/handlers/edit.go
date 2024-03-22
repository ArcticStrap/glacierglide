package handlers

import (
	"html/template"
	"io"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/models"
	"github.com/go-chi/chi/v5"
)

func SetupEditHandler(rt *chi.Mux, addr string) {
	// Inialize templates
	tmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/edit.html", "polarpages/templates/pagenav.html"))

	// Edit routing
	rt.Get("/e/{title}", func(w http.ResponseWriter, r *http.Request) {
		titleParam := chi.URLParam(r, "title")

		// Get page content
		res, err := http.Get(addr + "/api/s/" + titleParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		// Read response
		content, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse template
		tmpl.Execute(w, struct {
			models.SessionData
			models.WebPage
			models.WebModes
		}{models.SessionData{LoggedIn: false}, models.WebPage{
			Title:   titleParam,
			Content: template.HTML(content),
			Theme:   "common",
		}, models.WebModes{PageMode: "edit"}})
	})
}
