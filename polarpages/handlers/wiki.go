package handlers

import (
	"html/template"
	"io"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/models"
)

func SetupWikiHandler(rt *http.ServeMux, addr string) {
	// Initalize templates
	tmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/wiki.html", "polarpages/templates/pagenav.html"))

	// Redirects to main page
	rt.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			tmpl.Execute(w, struct {
				models.SessionData
				models.WebPage
				models.WebModes
			}{UserSession(r), models.WebPage{
				Title:   "Page Not Found (404)",
				Content: "Page not found. Try checking your spelling if otherwise",
				Theme:   "common",
			}, models.WebModes{PageMode: "read"}})
			return
		} else {
			http.Redirect(w, r, "/wiki/main_page", http.StatusSeeOther)
		}
	})
	rt.HandleFunc("GET /wiki", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/wiki/main_page", http.StatusSeeOther)
	})

	// Wiki routing
	rt.HandleFunc("GET /wiki/{title}", func(w http.ResponseWriter, r *http.Request) {
		titleParam := r.PathValue("title")

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
		err = tmpl.Execute(w, struct {
			models.SessionData
			models.WebPage
			models.WebModes
		}{UserSession(r), models.WebPage{
			Title:   titleParam,
			Content: template.HTML(content),
			Theme:   "common",
		}, models.WebModes{PageMode: "read"}})
	})

}
