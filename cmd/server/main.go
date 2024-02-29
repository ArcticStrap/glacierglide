package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ChaosIsFramecode/horinezumi/api/edit"
	"github.com/ChaosIsFramecode/horinezumi/api/history"
	"github.com/ChaosIsFramecode/horinezumi/api/source"
	"github.com/ChaosIsFramecode/horinezumi/api/user"
	"github.com/ChaosIsFramecode/horinezumi/api/wiki"
	"github.com/ChaosIsFramecode/horinezumi/api/wikido"
	"github.com/ChaosIsFramecode/horinezumi/appsignals"
	"github.com/ChaosIsFramecode/horinezumi/data"
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

	rt := chi.NewRouter()

	// Initalize the app signal system
	sc := appsignals.NewSignalConnector()

	// Use logger
	rt.Use(middleware.Logger)

	rt.Route("/api", func(apiroute chi.Router) {
		// Initalize subrouters
		wikido.SetupDoRoute(apiroute, &db)
		wiki.SetupWikiRoute(apiroute, &db)
		edit.SetupEditRoute(apiroute, &db, sc)
		history.SetupHistoryRoute(apiroute, &db)
		source.SetupSourceRoute(apiroute, &db)
		user.SetupUserRoute(apiroute, &db, sc)
	})

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
