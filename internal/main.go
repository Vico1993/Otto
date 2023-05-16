package main

import (
	"fmt"

	"github.com/subosito/gotenv"

	"github.com/Vico1993/Otto/internal/bot"
	"github.com/Vico1993/Otto/internal/cron"
	"github.com/Vico1993/Otto/internal/database"
)

func main() {
	fmt.Println("test lll")

	// load .env file if any otherwise use env set
	_ = gotenv.Load()

	// Load the database
	database.Init()

	// Initialisation of the cron
	cron.Init()

	// bot
	bot.Init()
}
