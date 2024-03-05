package handlers

import (
	"html/template"
	"io"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/models"
	"github.com/go-chi/chi/v5"
)


func SetupSourceHandler(rt *chi.Mux, addr string, tmpl *template.Template) {
	// Wiki routing
	rt.Get("/s/{title}", func(w http.ResponseWriter, r *http.Request) {
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
		tmpl.Execute(w, models.WebPage{
			Title:   titleParam,
			Content: template.HTML(content),
			Theme:   "common",
		})
	})
}
