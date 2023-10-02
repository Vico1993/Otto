package routes

import (
	"github.com/Vico1993/Otto/internal/handlers"
	"github.com/gin-gonic/gin"
)

// Init all the routes for the API
func Init(r *gin.Engine) {
	r.GET("/ping", handlers.Ping)

	// Load Route link to Chat
	chatsRoute(r)

	// Load Route link to Feed
	feedsRoute(r)

	// Load Route link to Article
	articlesRoute(r)
}
