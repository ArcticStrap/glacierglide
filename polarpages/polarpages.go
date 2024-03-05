package polarpages

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/handlers"
	"github.com/go-chi/chi/v5"
)

var tmpl *template.Template

func Setup(rt *chi.Mux, addr string) {
	// Parse templates
	tmpl = template.Must(template.ParseFiles("polarpages/templates/index.html"))
	eTmpl := template.Must(template.ParseFiles("polarpages/templates/edit.html"))
	sTmpl := template.Must(template.ParseFiles("polarpages/templates/source.html"))

	// Load skin assets
	sfs := http.FileServer(http.Dir("polarpages/skins"))
	rt.Handle("/skins/*", http.StripPrefix("/skins/", sfs))

	// Setup handlers
	handlers.SetupWikiHandler(rt, addr, tmpl)
	handlers.SetupEditHandler(rt, addr, eTmpl)
	handlers.SetupSourceHandler(rt, addr, sTmpl)

	log.Println("PolarPages initalized")
}
