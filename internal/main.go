package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/Vico1993/Otto/internal/cron"
	"github.com/Vico1993/Otto/internal/database"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		return
	}

	// Load the database
	database.Init()

	// Initialisation of the cron
	cron.Init()
}
