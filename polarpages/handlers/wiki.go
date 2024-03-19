package handlers

import (
	"html/template"
	"io"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/models"
	"github.com/go-chi/chi/v5"
)

func mpRedir(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/wiki/main_page", http.StatusSeeOther)
}

func SetupWikiHandler(rt *chi.Mux, addr string) {
	// Initalize templates
	tmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/wiki.html"))

	// Redirects to main page
	rt.Get("/", mpRedir)
	rt.Get("/wiki", mpRedir)

	// Not found (404) handler
	rt.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		tmpl.Execute(w, models.WebPage{
			Title:   "Page Not Found (404)",
			Content: "Page not found. Try checking your spelling if otherwise",
			Theme:   "common",
		})
	})

	// Wiki routing
	rt.Get("/wiki/{title}", func(w http.ResponseWriter, r *http.Request) {
		titleParam := chi.URLParam(r, "title")

		// Get page content
		res, err := http.Get(addr + "/api/wiki/" + titleParam)
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
		}{models.SessionData{LoggedIn: false}, models.WebPage{
			Title:   titleParam,
			Content: template.HTML(content),
			Theme:   "common",
		}})
	})
}
