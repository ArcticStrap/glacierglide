package main

import (
	"log"
	"net/http"

	"github.com/ChaosIsFramecode/horinezumi/wikiroute"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	rt := chi.NewRouter()

	// Use logger
	rt.Use(middleware.Logger)

	// Redirect root path to main page
	rt.Get("/", http.RedirectHandler("/wiki/Main_Page", http.StatusSeeOther).ServeHTTP)

	// Wiki sub router
	wikiroute.SetupWikiroute(rt)

	log.Println("Running on 127.0.0.1:8080")
	http.ListenAndServe(":8080", rt)
}
