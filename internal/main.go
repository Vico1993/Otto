package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/Vico1993/Otto/internal/middlewares"
	v2 "github.com/Vico1993/Otto/internal/repository/v2"
	"github.com/Vico1993/Otto/internal/routes"
)

func main() {
	// load .env file if any otherwise use env set
	_ = gotenv.Load()

	// Load the database
	database.Init()

	// Load repository
	v2.Init()

	r := gin.Default()

	// Error Middleware
	r.Use(middlewares.Error())

	// Init routes
	routes.Init(r)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Load repository
	// repository.Init()

	// service.NewTelegramService().TelegramPostMessage(
	// 	os.Getenv("TELEGRAM_USER_CHAT_ID"),
	// 	`*Upgrade complete*! Ready to be even smarter and funnier than before. 🤖 🚀 ✨`,
	// )

	// Initialisation of the cron
	// cron.Init()
}
