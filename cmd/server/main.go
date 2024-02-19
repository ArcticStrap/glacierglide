package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ChaosIsFramecode/horinezumi/appsignals"
	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/subroutes/edit"
	"github.com/ChaosIsFramecode/horinezumi/subroutes/history"
	"github.com/ChaosIsFramecode/horinezumi/subroutes/user"
	"github.com/ChaosIsFramecode/horinezumi/subroutes/wiki"
	"github.com/ChaosIsFramecode/horinezumi/subroutes/wikido"
	"github.com/ChaosIsFramecode/horinezumi/utils/environment"
)

func main() {
	// Load .env variables
	err := environment.Load()
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

  // Initalize the app signal system
  sc := appsignals.NewSignalConnector()

	// Use logger
	rt.Use(middleware.Logger)

	// Redirect root path to main page
	rt.Get("/", http.RedirectHandler("/wiki/Main_Page", http.StatusSeeOther).ServeHTTP)

	// Initalize subrouters
	wikido.SetupDoRoute(rt, &db)
	wiki.SetupWikiRoute(rt, &db)
	edit.SetupEditRoute(rt, &db)
	history.SetupHistoryRoute(rt, &db)
	user.SetupUserRoute(rt, &db,sc)

	log.Println("Running on " + os.Getenv("ADDR"))

	// Setup signal handling for cleanup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Closing server and database")

		db.Close()
		os.Exit(0)
	}()

	// Start up web server
	if os.Getenv("DEV") == "" {
		http.ListenAndServeTLS(os.Getenv("ADDR"), "certs/cert.pem", "certs/key.pem", rt)
	} else {
		log.Println("(MODE: DEBUG)")
		http.ListenAndServe(os.Getenv("ADDR"), rt)
	}
}
