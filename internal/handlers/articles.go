package handlers

import (
	"net/http"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
)

type articleCreatePost struct {
	FeedId string   `json:"feed_id" binding:"required"`
	Title  string   `json:"title" binding:"required"`
	Source string   `json:"source" binding:"required"`
	Link   string   `json:"link" binding:"required"`
	Author string   `json:"author" binding:"required"`
	Tags   []string `form:"tags" binding:"required"`
}

// Retrieve all article
func GetAllArticles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"articles": repository.Article.GetAll()})
}

// Retrieve one article
func GetArticle(c *gin.Context) {
	article := c.MustGet("article").(*repository.DBArticle)

	c.JSON(http.StatusOK, gin.H{"article": article})
}

// Delete one article
func DeleteArticle(c *gin.Context) {
	article := c.MustGet("article").(*repository.DBArticle)

	c.JSON(http.StatusOK, gin.H{"deleted": repository.Article.Delete(article.FeedId)})
}

// Create article
func CreateArticle(c *gin.Context) {
	var json articleCreatePost
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := repository.Article.Create(json.FeedId, json.Title, json.Source, json.Author, json.Link, json.Tags)

	c.JSON(http.StatusOK, gin.H{"article": article})
}
