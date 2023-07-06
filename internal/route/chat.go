package route

import (
	"github.com/Vico1993/Otto/internal/handler"
	"github.com/gin-gonic/gin"
)

// Enable all the routes relates to chat
func chatRoutes(r *gin.Engine) {
	chatGroup := r.Group("/chat")

	chatGroup.GET("/", handler.GetAllChat)
	chatGroup.GET("/:id", handler.GetChatById)
	chatGroup.GET("/:id/feeds", handler.GetChatFeeds)
	chatGroup.POST("/:id/feeds", handler.CreateFeed)
}
