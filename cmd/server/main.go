package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArcticStrap/glacierglide/api"
	"github.com/ArcticStrap/glacierglide/appsignals"
	"github.com/ArcticStrap/glacierglide/data"
	"github.com/ArcticStrap/glacierglide/middleware"
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

	rt := http.NewServeMux()

	// Initalize the app signal system
	sc := appsignals.NewSignalConnector()

	// Setup api
	api.Setup(rt, &db, sc)

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

	// Setup middleware
	mStack := middleware.CreateStack(
		middleware.Logging,
	)

	// Start up web server

	server := http.Server{
		Addr:    os.Getenv("ADDR"),
		Handler: mStack(rt),
	}

	log.Println("Running on " + os.Getenv("ADDR"))
	if os.Getenv("DEV") == "" {
		server.ListenAndServeTLS("certs/cert.pem", "certs/key.pem")
	} else {
		log.Println("(MODE: DEBUG)")
		log.Fatal(server.ListenAndServe())
	}
}
