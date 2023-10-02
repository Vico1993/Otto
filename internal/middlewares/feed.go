package middlewares

import (
	"net/http"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
)

// Feed Mildleware to make sure Feed exist
func ValidFeed() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("feedid")
		feed := repository.Feed.GetOne(uuid)

		if feed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request parameters"})
			return
		}

		// Set chat for next time
		c.Set("feed", feed)

		c.Next()
	}
}
