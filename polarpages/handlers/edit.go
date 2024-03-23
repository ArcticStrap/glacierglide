package handlers

import (
	"html/template"
	"io"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/models"
)

func SetupEditHandler(rt *http.ServeMux, addr string) {
	// Inialize templates
	tmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/edit.html", "polarpages/templates/pagenav.html"))

	// Edit routing
	rt.HandleFunc("GET /e/{title}", func(w http.ResponseWriter, r *http.Request) {
		titleParam := r.PathValue("title")

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
		}{UserSession(r), models.WebPage{
			Title:   titleParam,
			Content: template.HTML(content),
			Theme:   "common",
		}, models.WebModes{PageMode: "edit"}})
	})
}
