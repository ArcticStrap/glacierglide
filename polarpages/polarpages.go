package polarpages

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/handlers"
	"github.com/go-chi/chi/v5"
)

func Setup(rt *chi.Mux, addr string) {
	// Parse templates
	wTmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/wiki.html"))
	eTmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/edit.html"))
	sTmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/source.html"))

	// Load skin assets
	sfs := http.FileServer(http.Dir("polarpages/skins"))
	rt.Handle("/skins/*", http.StripPrefix("/skins/", sfs))

	// Setup handlers
	handlers.SetupWikiHandler(rt, addr, wTmpl)
	handlers.SetupEditHandler(rt, addr, eTmpl)
	handlers.SetupSourceHandler(rt, addr, sTmpl)

	log.Println("PolarPages initalized")
}
