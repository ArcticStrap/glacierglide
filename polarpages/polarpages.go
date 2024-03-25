package polarpages

import (
	"log"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/handlers"
)

func Setup(rt *http.ServeMux, addr string) {
	// Setup handlers
	handlers.SetupWikiHandler(rt, addr)
	handlers.SetupEditHandler(rt, addr)
	handlers.SetupHistoryHandler(rt, addr)
	handlers.SetupSourceHandler(rt, addr)
	handlers.SetupUserHandler(rt)

	// Load skin assets
	rt.Handle("GET /skins/*", http.StripPrefix("/skins/", http.FileServer(http.Dir("polarpages/skins"))))

	log.Println("PolarPages initalized")
}
