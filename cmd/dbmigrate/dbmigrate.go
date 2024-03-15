package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/ArcticStrap/glacierglide/utils/environment"
)

func main() {
	// Load .env variables
	err := environment.Load()
	if err != nil {
		log.Fatalf("Error loading environment variables: %s\n", err)
		return
	}

	// Extract arguments
	command := os.Args[1]
	var version string
	if len(os.Args) >= 3 {
		version = os.Args[2]
	}
	url := os.Getenv("PAGEDATAURL")

	cmd := exec.Command("./cmd/dbmigrate/migrate", "-database", url, "-path", "./migrations")

	switch command {
	case "up":
		cmd.Args = append(cmd.Args, "up")
	case "down":
		cmd.Args = append(cmd.Args, "down")
	case "force":
		cmd.Args = append(cmd.Args, "force", version)
	default:
		log.Fatalf("Invalid command. Available commands: up, down, force\n")
		return
	}

	// Execute
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing command: %s\n", err)
		return
	}

	// Print output
	log.Println(string(output))
}
