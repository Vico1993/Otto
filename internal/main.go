package main

import (
	"github.com/subosito/gotenv"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/Vico1993/Otto/internal/repository"
)

func main() {
	// load .env file if any otherwise use env set
	_ = gotenv.Load()

	// Load the database
	database.Init()

	// Load repository
	repository.Init()

	// Initialisation of the cron
	// cron.Init()

	// bot
	// bot.Init()
}
