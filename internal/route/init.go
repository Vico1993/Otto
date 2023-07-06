package route

import (
	"github.com/Vico1993/Otto/internal/handler"
	"github.com/gin-gonic/gin"
)

// Init routes
func Init(r *gin.Engine) {
	r.GET("/ping", handler.Ping)
}
