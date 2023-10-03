package routes

import (
	"github.com/Vico1993/Otto/internal/handlers"
	"github.com/Vico1993/Otto/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// Enable all the routes relates to apple
func articlesRoute(r *gin.Engine) {
	articles := r.Group("/articles")
	{
		articles.GET("/", handlers.GetAllArticles)
		articles.POST("/", handlers.CreateArticle)

		articleId := articles.Group("/:articleid", middlewares.ValidArticle())
		{
			articleId.GET("/", handlers.GetArticle)
			articleId.DELETE("/", handlers.DeleteArticle)
		}
	}

}
