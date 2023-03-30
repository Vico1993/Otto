package main

import (
	"github.com/subosito/gotenv"

	"github.com/Vico1993/Otto/internal/cron"
	"github.com/Vico1993/Otto/internal/database"
)

func main() {
	// load .env file
	gotenv.Load()

	// Load the database
	database.Init()

	// Initialisation of the cron
	cron.Init()
}
