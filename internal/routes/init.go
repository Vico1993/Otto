package routes

import (
	"net/http"

	"github.com/Vico1993/Otto/internal/handlers"
	v2 "github.com/Vico1993/Otto/internal/repository/v2"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.GET("/ping", handlers.Ping)
	r.GET("/test", func(c *gin.Context) {
		repo := v2.SChatRepository{}

		c.JSON(http.StatusOK, gin.H{
			"chats": repo.GetAll(),
		})
	})
}
