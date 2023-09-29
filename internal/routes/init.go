package routes

import (
	"github.com/Vico1993/Otto/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.GET("/ping", handlers.Ping)
}
