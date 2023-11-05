package middlewares

import (
	"net/http"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
)

// Chat Mildleware to make sure Chat exist
func ValidChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		chatId := c.Param("chatid")
		threadId := c.Param("threadid")

		var chat *repository.DBChat
		if threadId != "" {
			chat = repository.Chat.GetByTelegramChatIdAndThreadId(chatId, threadId)
		} else {
			chat = repository.Chat.GetByTelegramChatId(chatId)
		}

		if chat == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request parameters"})
			return
		}

		// Set chat for next time
		c.Set("chat", chat)

		c.Next()
	}
}
