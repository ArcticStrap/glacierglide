package polarpages

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type webPage struct {
	Title   string
	Content template.HTML
	Theme   string
}

var tmpl *template.Template

func Setup(rt *chi.Mux, addr string) {
	// Parse templates
	tmpl = template.Must(template.ParseFiles("polarpages/templates/index.html"))

	// Load skin assets
	sfs := http.FileServer(http.Dir("polarpages/skins"))
	rt.Handle("/skins/*", http.StripPrefix("/skins/", sfs))

	// Not found (404) handler
	rt.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		tmpl.Execute(w, webPage{
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
		tmpl.Execute(w, webPage{
			Title:   titleParam,
			Content: template.HTML(content),
			Theme:   "common",
		})
	})

	log.Println("PolarPages initalized")
}
