package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
)

// Article Mildleware to make sure Article exist
func ValidArticle() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("articleid")
		article := repository.Article.GetOne(uuid)

		fmt.Println(article, uuid)

		if article == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request parameters"})
			return
		}

		// Set article for next time
		c.Set("article", article)

		c.Next()
	}
}
