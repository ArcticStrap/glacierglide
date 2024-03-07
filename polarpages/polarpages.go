package polarpages

import (
	"log"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/handlers"
	"github.com/go-chi/chi/v5"
)

func Setup(rt *chi.Mux, addr string) {
	// Parse templates

	// Load skin assets
	sfs := http.FileServer(http.Dir("polarpages/skins"))
	rt.Handle("/skins/*", http.StripPrefix("/skins/", sfs))

	// Setup handlers
	handlers.SetupWikiHandler(rt, addr)
	handlers.SetupEditHandler(rt, addr)
	handlers.SetupSourceHandler(rt, addr)
	handlers.SetupUserHandler(rt)

	log.Println("PolarPages initalized")
}
