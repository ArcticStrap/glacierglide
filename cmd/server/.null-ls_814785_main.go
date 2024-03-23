package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ArcticStrap/glacierglide/api/edit"
	"github.com/ArcticStrap/glacierglide/api/history"
	"github.com/ArcticStrap/glacierglide/api/source"
	"github.com/ArcticStrap/glacierglide/api/user"
	"github.com/ArcticStrap/glacierglide/api/wiki"
	"github.com/ArcticStrap/glacierglide/api/wikido"
	"github.com/ArcticStrap/glacierglide/appsignals"
	"github.com/ArcticStrap/glacierglide/data"
	"github.com/ArcticStrap/glacierglide/polarpages"
	"github.com/ArcticStrap/glacierglide/utils/environment"
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

	// Setup api
	rt.Route("/api", func(apiroute chi.Router) {
		// Initalize subrouters
		wikido.SetupDoRoute(apiroute, &db)
		wiki.SetupWikiRoute(apiroute, &db)
		edit.SetupEditRoute(apiroute, &db, sc)
		history.SetupHistoryRoute(apiroute, &db)
		source.SetupSourceRoute(apiroute, &db)
		user.SetupUserRoute(apiroute, &db, sc)
	})

	// Setup client (if permitted)
	if os.Getenv("WITHPOLARP") != "" {
		hostAddr := "http"
		if os.Getenv("DEV") == "" {
			hostAddr += "s"
		}
		hostAddr += "://" + os.Getenv("ADDR")
		polarpages.Setup(rt, hostAddr)
	}

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
	log.Println("Running on " + os.Getenv("ADDR"))
	if os.Getenv("DEV") == "" {
		http.ListenAndServeTLS(os.Getenv("ADDR"), "certs/cert.pem", "certs/key.pem", rt)
	} else {
		log.Println("(MODE: DEBUG)")
		http.ListenAndServe(os.Getenv("ADDR"), rt)
	}
}
