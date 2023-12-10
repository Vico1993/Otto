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
		chats.GET("/", handlers.GetAllChats)

		chatId := chats.Group("/:chatid", middlewares.ValidChat())
		addChatGroup(chatId)

		threadId := chats.Group("/:chatid/:threadid", middlewares.ValidChat())
		addChatGroup(threadId)
	}

}

func addChatGroup(group *gin.RouterGroup) {
	group.DELETE("/", handlers.DeleteChat)
	group.GET("/parsed", handlers.ParsedChat)

	feeds := group.Group("/feeds")
	{
		feeds.GET("/", handlers.GetChatFeeds)
		feeds.POST("/:feedid", middlewares.ValidFeed(), handlers.CreateChatFeed)
		feeds.DELETE("/:feedid", middlewares.ValidFeed(), handlers.DeleteChatFeed)
	}

	tags := group.Group("/tags")
	{
		tags.GET("/", handlers.GetChatTags)
		tags.POST("/", handlers.CreateChatTag)
		tags.DELETE("/:tag", handlers.DeleteChatTag)
	}

	articles := group.Group("/articles")
	{
		articles.GET("/latest", handlers.GetLatestArticleFromChat)
	}
}
