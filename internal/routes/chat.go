package routes

import (
	"github.com/Vico1993/Otto/internal/handlers"
	"github.com/Vico1993/Otto/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// Enable all the routes relates to apple
func chatsRoute(r *gin.Engine) {
	chats := r.Group("/chats")
	{
		chats.POST("/", handlers.CreateChat)

		chatId := chats.Group("/:chatid", middlewares.ValidChat())
		{
			chatId.DELETE("/", handlers.DeleteChat)
			chatId.GET("/parsed", handlers.ParsedChat)

			feeds := chatId.Group("/feeds")
			{
				feeds.GET("/", handlers.GetChatFeeds)
				feeds.POST("/:feedid", middlewares.ValidFeed(), handlers.CreateChatFeed)
				feeds.DELETE("/:feedid", middlewares.ValidFeed(), handlers.DeleteChatFeed)
			}

			tags := chatId.Group("/tags")
			{
				tags.GET("/", handlers.GetChatTags)
				tags.POST("/", handlers.CreateChatTag)
				tags.DELETE("/:tag", handlers.DeleteChatTag)
			}
		}
	}

}
