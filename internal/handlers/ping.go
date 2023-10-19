package handlers

import (
	"net/http"

	"github.com/Vico1993/Otto/internal/utils"
	"github.com/gin-gonic/gin"
)

// Enpoint to make sure the API is alive
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"version": utils.RetrieveVersion(),
	})
}
