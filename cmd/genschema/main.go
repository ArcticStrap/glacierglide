package main

import (
	"log"

	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/utils/environment"
)

func main() {
	// Load .env variables
	err := environment.Load()
	if err != nil {
		log.Fatalf("Error loading environment variables: %s", err)
	}

	db, err := data.ConnectToPostgresDatabase()
	if err != nil {
		log.Fatalf("Error connecting to data base: %s", err)
	} else {
		log.Println("Successfully connected to database.")
	}

	if err = db.CreateTables(); err != nil {
		log.Fatalf("Failed to create table: %s", err)
	} else {
		log.Println("Successfully created databse schemas")
	}

	db.Close()
	log.Println("Database closed.")
}
