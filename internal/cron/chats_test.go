package cron

import (
	"fmt"
	"testing"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCheckChatNoChatsInDB(t *testing.T) {
	chatRepositoryMock := new(repository.MocksChatRepository)
	repository.Chat = chatRepositoryMock

	chatRepositoryMock.On("GetAll").Return([]*repository.DBChat{})

	checkChat()

	jobs, _ := Scheduler.FindJobsByTag(chatsTag)
	tags := Scheduler.GetAllTags()

	assert.Len(t, tags, 0)
	chatRepositoryMock.AssertCalled(t, "GetAll")
	assert.Nil(t, jobs, "Should have no job link to "+chatsTag)
}

func TestCheckChatNoNeedToAddSameNumber(t *testing.T) {
	_, _ = Scheduler.Every(1).Tag(chatsTag).Week().Do(func() {
		fmt.Println("Test job")
	})

	chatRepositoryMock := new(repository.MocksChatRepository)
	repository.Chat = chatRepositoryMock

	chat1 := repository.DBChat{
		Id:             "12",
		TelegramChatId: "12314",
		Tags:           []string{"tag1", "tag2"},
	}

	chatRepositoryMock.On("GetAll").Return([]*repository.DBChat{&chat1})

	jobs, _ := Scheduler.FindJobsByTag(chatsTag)
	assert.Len(t, jobs, 1, "Should have 1 job in the queue at start")

	checkChat()

	jobs, _ = Scheduler.FindJobsByTag(chatsTag)
	tags := Scheduler.GetAllTags()

	assert.Len(t, tags, 1)
	assert.Equal(t, []string{"chat"}, tags, "Should have 1 tags call chats")
	chatRepositoryMock.AssertCalled(t, "GetAll")
	assert.Len(t, jobs, 1, "Should have 1 job in the queue after call")

	// Reset jobs
	_ = Scheduler.RemoveByTag(chatsTag)
}

func TestCheckChatAddJob(t *testing.T) {
	_, _ = Scheduler.Every(1).Tag(chatsTag).Week().Do(func() {
		fmt.Println("Test job")
	})

	chatRepositoryMock := new(repository.MocksChatRepository)
	repository.Chat = chatRepositoryMock

	chat1 := repository.DBChat{
		Id:             "12",
		TelegramChatId: "12314",
		Tags:           []string{"tag1", "tag2"},
	}

	chat2 := repository.DBChat{
		Id:             "13",
		TelegramChatId: "12315",
		Tags:           []string{"tag1", "tag2"},
	}

	chatRepositoryMock.On("GetAll").Return([]*repository.DBChat{&chat1, &chat2})

	jobs, _ := Scheduler.FindJobsByTag(chatsTag)
	assert.Len(t, jobs, 1, "Should have 1 job in the queue at start")

	checkChat()

	jobs, _ = Scheduler.FindJobsByTag(chatsTag)
	tags := Scheduler.GetAllTags()

	assert.Len(t, tags, 2)
	assert.Equal(t, []string{"chat", "chat"}, tags, "Should have 2 tags call chats")
	chatRepositoryMock.AssertCalled(t, "GetAll")
	assert.Len(t, jobs, 2, "Should have 2 job in the queue after call")

	// Reset jobs
	_ = Scheduler.RemoveByTag(chatsTag)
}

func TestParsedArticlesNoArticlesFound(t *testing.T) {
	chat := repository.DBChat{
		Id:             "12",
		TelegramChatId: "12314",
		Tags:           []string{"tag3"},
	}

	mockTelegramService := new(service.MocksTelegramService)
	telegram = mockTelegramService

	parsedArticles([]*repository.DBArticle{}, &chat)

	mockTelegramService.AssertNotCalled(t, "TelegramUpdateTyping", chat.TelegramChatId, true)
	mockTelegramService.AssertNotCalled(t, "TelegramUpdateTyping", chat.TelegramChatId, false)
	mockTelegramService.AssertNotCalled(t, "TelegramPostMessage", chat.TelegramChatId, "MESSAGE")
}

func TestParsedArticlesWith1Article(t *testing.T) {
	article := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: uuid.New().String(),
		Title:  "Super Title",
		Source: "Title",
		Author: "Unknown",
		Link:   "https://super.com/title",
		Tags:   []string{"tag1", "tag2"},
	}

	chat := repository.DBChat{
		Id:               "12",
		TelegramChatId:   "12314",
		TelegramThreadId: nil,
		Tags:             []string{"tag2"},
	}

	templates = []string{"MESSAGE"}

	mockTelegramService := new(service.MocksTelegramService)
	telegram = mockTelegramService
	mockTelegramService.On("TelegramUpdateTyping", chat.TelegramChatId, true).Return()
	mockTelegramService.On("TelegramUpdateTyping", chat.TelegramChatId, false).Return()
	mockTelegramService.On("TelegramPostMessage", chat.TelegramChatId, "", "MESSAGE").Return()

	parsedArticles([]*repository.DBArticle{&article}, &chat)

	mockTelegramService.AssertCalled(t, "TelegramUpdateTyping", chat.TelegramChatId, true)
	mockTelegramService.AssertCalled(t, "TelegramUpdateTyping", chat.TelegramChatId, false)
	mockTelegramService.AssertCalled(t, "TelegramPostMessage", chat.TelegramChatId, "", "MESSAGE")
}

func TestParsedArticlesWith1ArticleAnd1ThreadId(t *testing.T) {
	article := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: uuid.New().String(),
		Title:  "Super Title",
		Source: "Title",
		Author: "Unknown",
		Link:   "https://super.com/title",
		Tags:   []string{"tag1", "tag2"},
	}

	threadId := "134"
	chat := repository.DBChat{
		Id:               "12",
		TelegramChatId:   "12314",
		TelegramThreadId: &threadId,
		Tags:             []string{"tag2"},
	}

	templates = []string{"MESSAGE"}

	mockTelegramService := new(service.MocksTelegramService)
	telegram = mockTelegramService
	mockTelegramService.On("TelegramUpdateTyping", chat.TelegramChatId, true).Return()
	mockTelegramService.On("TelegramUpdateTyping", chat.TelegramChatId, false).Return()
	mockTelegramService.On("TelegramPostMessage", chat.TelegramChatId, threadId, "MESSAGE").Return()

	parsedArticles([]*repository.DBArticle{&article}, &chat)

	mockTelegramService.AssertCalled(t, "TelegramUpdateTyping", chat.TelegramChatId, true)
	mockTelegramService.AssertCalled(t, "TelegramUpdateTyping", chat.TelegramChatId, false)
	mockTelegramService.AssertCalled(t, "TelegramPostMessage", chat.TelegramChatId, threadId, "MESSAGE")
}
