package handler

import (
	"net/http"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
)

// Retrieve all Chat
func GetAllChat(c *gin.Context) {
	chats := repository.Chat.GetAll()

	c.JSON(http.StatusOK, chats)
}

// Retrieve chat with its id
func GetChatById(c *gin.Context) {
	chat := repository.Chat.FindByChatId(c.Param("id"))

	if chat == nil {
		c.JSON(http.StatusNotFound, nil)
	}

	c.JSON(http.StatusOK, chat)
}

// Retrieve all feeds with chat id
func GetChatFeeds(c *gin.Context) {
	chat := repository.Chat.FindByChatId(c.Param("id"))

	if chat != nil {
		c.JSON(http.StatusOK, chat.Feeds)
	}

	c.JSON(http.StatusNotFound, nil)
}

type postFeed struct {
	Url string `json:"url"`
}

// Add feed to the chat
func CreateFeed(c *gin.Context) {
	var data postFeed
	err := c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	added := repository.Chat.PushNewFeed(data.Url, c.Param("id"))

	if added {
		c.JSON(http.StatusOK, repository.Chat.FindByChatId(c.Param("id")).Feeds)
	} else {
		c.JSON(http.StatusInternalServerError, nil)
	}
}
