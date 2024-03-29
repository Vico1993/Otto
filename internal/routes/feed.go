package routes

import (
	"github.com/Vico1993/Otto/internal/handlers"
	"github.com/Vico1993/Otto/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// Enable all the routes relates to apple
func feedsRoute(r *gin.Engine) {
	feeds := r.Group("/feeds")
	{
		feeds.GET("/", handlers.GetAllFeeds)
		feeds.GET("/active", handlers.GetAllActiveFeeds)
		feeds.POST("/", handlers.CreateFeed)

		feedId := feeds.Group("/:feedid", middlewares.ValidFeed())
		{
			feedId.GET("/", handlers.GetFeed)
			feedId.DELETE("/", handlers.DeleteFeed)
			feedId.PUT("/disable", handlers.DisableFeed)

			feeds := feedId.Group("/articles")
			{
				feeds.GET("/", handlers.GetFeedArticles)
			}
		}
	}

}
