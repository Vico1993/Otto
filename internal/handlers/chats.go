package handlers

import (
	"net/http"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
)

type chatCreatePost struct {
	ChatId string   `json:"chat_id" binding:"required"`
	UserId string   `json:"user_id"`
	Tags   []string `json:"tags" binding:"required"`
}

type createChatTagsPost struct {
	Tags []string `json:"tags" binding:"required"`
}

// Road to create a Chat
func CreateChat(c *gin.Context) {
	var json chatCreatePost
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat := repository.Chat.Create(json.ChatId, json.UserId, json.Tags)

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

// Road to delete a Chat
func DeleteChat(c *gin.Context) {
	chat := c.MustGet("chat").(*repository.DBChat)

	c.JSON(http.StatusOK, gin.H{"deleted": repository.Chat.Delete(chat.Id)})
}

// Retrieve all feeds link to a Chat
func GetChatFeeds(c *gin.Context) {
	chat := c.MustGet("chat").(*repository.DBChat)

	c.JSON(http.StatusOK, gin.H{"feeds": repository.Feed.GetByChatId(chat.Id)})
}

// Link feed to a chat
func CreateChatFeed(c *gin.Context) {
	chat := c.MustGet("chat").(*repository.DBChat)
	feed := c.MustGet("feed").(*repository.DBFeed)

	c.JSON(http.StatusOK, gin.H{"added": repository.Feed.LinkChatAndFeed(feed.Id, chat.Id)})
}

// Delete link between chat and feed
func DeleteChatFeed(c *gin.Context) {
	chat := c.MustGet("chat").(*repository.DBChat)
	feed := c.MustGet("feed").(*repository.DBFeed)

	c.JSON(http.StatusOK, gin.H{"deleted": repository.Feed.UnLinkChatAndFeed(feed.Id, chat.Id)})
}

// Retrieve all tags link to a Chat
func GetChatTags(c *gin.Context) {
	chat := c.MustGet("chat").(*repository.DBChat)

	c.JSON(http.StatusOK, gin.H{"tags": chat.Tags})
}

// Add tags to a chat
func CreateChatTag(c *gin.Context) {
	chat := c.MustGet("chat").(*repository.DBChat)

	var json createChatTagsPost
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTagList := chat.Tags
	newTagList = append(newTagList, json.Tags...)

	repository.Chat.UpdateTags(chat.Id, newTagList)

	c.JSON(http.StatusOK, gin.H{"tags": newTagList})
}

// Delete tags to a chat
func DeleteChatTag(c *gin.Context) {
	chat := c.MustGet("chat").(*repository.DBChat)
	tag := c.Param("tag")

	tagFound := false
	for k, t := range chat.Tags {
		if t == tag {
			tagFound = true
			chat.Tags = append(chat.Tags[:k], chat.Tags[k+1:]...)
		}
	}

	if !tagFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag not found"})
		return
	}

	repository.Chat.UpdateTags(chat.Id, chat.Tags)

	c.JSON(http.StatusOK, gin.H{"tags": chat.Tags})
}

// Update Parsed value for Chat
func ParsedChat(c *gin.Context) {
	chat := c.MustGet("chat").(*repository.DBChat)

	repository.Chat.UpdateParsed(chat.Id)

	c.JSON(http.StatusNoContent, gin.H{})
}
