package data

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestDataCreaton(t *testing.T) {
	// Load .env variables
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading environment variables: %s", err)
	}

	db, err := ConnectToPostgresDatabase()
	if err != nil {
		t.Fatalf("Error connecting to data base: %s", err)
	} else {
		t.Log("Successfully connected to database.")
	}

	if err = db.CreateTables(); err != nil {
		t.Fatalf("Failed to create table: %s", err)
	}

	db.Close()
	t.Log("Database disconnected.")
}
