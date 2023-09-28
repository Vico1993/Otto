package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Enpoint to make sure the API is alive
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "0.1",
	})
}
