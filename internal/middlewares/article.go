package middlewares

import (
	"net/http"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
)

// Article Mildleware to make sure Article exist
func ValidArticle() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("articleid")
		article := repository.Article.GetOne(uuid)

		if article == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request parameters"})
			return
		}

		// Set article for next time
		c.Set("article", article)

		c.Next()
	}
}
