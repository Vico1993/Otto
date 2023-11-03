package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"

	"github.com/Vico1993/Otto/internal/cron"
	"github.com/Vico1993/Otto/internal/database"
	"github.com/Vico1993/Otto/internal/middlewares"
	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/routes"
	"github.com/Vico1993/Otto/internal/service"
	"github.com/Vico1993/Otto/internal/utils"
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

	// Start cron exec
	cron.Scheduler.StartAsync()

	// Notify update if chat present
	if os.Getenv("TELEGRAM_USER_CHAT_ID") != "" {
		version := utils.RetrieveVersion()

		service.NewTelegramService().TelegramPostMessage(
			os.Getenv("TELEGRAM_USER_CHAT_ID"),
			`🚀 🚀 [CRON-API] Version: *`+version+`* Succesfully deployed . 🚀 🚀`,
		)
	}

	r := gin.Default()

	// Error Middleware
	r.Use(middlewares.Error())

	// Init routes
	routes.Init(r)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
