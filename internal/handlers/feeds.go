package handlers

import (
	"net/http"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
)

type feedCreatePost struct {
	Url string `json:"url" binding:"required"`
}

// Road to create a Feed
func CreateFeed(c *gin.Context) {
	var json feedCreatePost
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feed := repository.Feed.Create(json.Url)

	c.JSON(http.StatusOK, gin.H{"feed": feed})
}

// Retrieve all Feed
func GetAllFeeds(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"feeds": repository.Feed.GetAll()})
}

// Retrieve one Feed
func GetFeed(c *gin.Context) {
	feed := c.MustGet("feed").(*repository.DBFeed)

	c.JSON(http.StatusOK, gin.H{"feed": feed})
}

// Road to delete a Feed
func DeleteFeed(c *gin.Context) {
	feed := c.MustGet("feed").(*repository.DBFeed)

	c.JSON(http.StatusOK, gin.H{"deleted": repository.Feed.Delete(feed.Id)})
}

// Get articles link to a feed
func GetFeedArticles(c *gin.Context) {
	feed := c.MustGet("feed").(*repository.DBFeed)

	c.JSON(http.StatusOK, gin.H{"articles": repository.Article.GetByFeedId(feed.Id)})
}
