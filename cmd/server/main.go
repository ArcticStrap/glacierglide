package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/subroutes/edit"
	"github.com/ChaosIsFramecode/horinezumi/subroutes/history"
	"github.com/ChaosIsFramecode/horinezumi/subroutes/user"
	"github.com/ChaosIsFramecode/horinezumi/subroutes/wiki"
)

func main() {
	// Load .env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables: %s", err)
	}

	// Connect to our database
	db, err := data.ConnectToPostgresDatabase()
	if err != nil {
		log.Fatalf("Error connecting to data base: %s", err)
	} else {
		log.Printf("Successfully connected to database")
	}
	defer db.Close()

	if err = db.CreateTables(); err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}

	rt := chi.NewRouter()

	// Use logger
	rt.Use(middleware.Logger)

	// Redirect root path to main page
	rt.Get("/", http.RedirectHandler("/wiki/Main_Page", http.StatusSeeOther).ServeHTTP)

	// Initalize subrouters
	wiki.SetupWikiroute(rt, &db)
	edit.SetupEditRoute(rt, &db)
	history.SetupHistoryRoute(rt, &db)
	user.SetupUserRoute(rt, &db)

	log.Println("Running on " + os.Getenv("ADDR"))

	if os.Getenv("DEV") == "" {
		http.ListenAndServeTLS(os.Getenv("ADDR"), "certs/cert.pem", "certs/key.pem", rt)
	} else {
		log.Println("(MODE: DEBUG)")
		http.ListenAndServe(os.Getenv("ADDR"), rt)
	}

}
