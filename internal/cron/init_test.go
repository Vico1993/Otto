package cron

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/service"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDelayWith1Element(t *testing.T) {
	res := getDelay(1)

	assert.Equal(t, 60, res, "If only one element is passed a delay of 60 should be returned")
}

func TestDelayWith20Element(t *testing.T) {
	res := getDelay(20)

	assert.Equal(t, 3, res, "If 20 element are sent, should return 3")
}

func TestResetCronForChat(t *testing.T) {
	chat := database.NewChat("1234", 123, nil)
	_, _ = scheduler.Every(1).Tag(chat.ChatId).Week().Do(func() {
		fmt.Println("Test job")
	})

	assert.Len(t, scheduler.Jobs(), 1, "We should have only 1 job present at the initialisation")

	resetCronForChatId(chat)

	assert.Len(t, scheduler.Jobs(), 0, "After the delete we should have 0 jobs")
}

func TestAddJobForChat(t *testing.T) {
	feed := database.NewFeed("https://google.com")
	chat := database.NewChat("1234", 123, []database.Feed{
		*feed,
	})

	startJobForChat(chat)

	jobs := scheduler.Jobs()
	assert.Len(t, jobs, 1, "We should have only 1 job")

	job := jobs[0]
	assert.Len(t, job.Tags(), 1, "Only one tag should be attached to the job")
	assert.Equal(t, job.Tags()[0], chat.ChatId, "The first and only tag should be equal to our chatid")

	_ = scheduler.RemoveByTag(chat.ChatId)
}

func TestSetupCronForChat(t *testing.T) {
	feed := database.NewFeed("https://google.com")
	chat := database.NewChat("1234", 123, []database.Feed{
		*feed,
	})
	_, _ = scheduler.Every(1).Tag(chat.ChatId).Week().Do(func() {
		fmt.Println("Test job")
	})
	_, _ = scheduler.Every(1).Tag(chat.ChatId).Week().Do(func() {
		fmt.Println("Test job2")
	})

	fmt.Println(len(scheduler.Jobs()))

	assert.Len(t, scheduler.Jobs(), 2, "We should have only 2 job present at the initialisation")

	setupCronForChat(chat)

	jobs := scheduler.Jobs()
	assert.Len(t, jobs, 1, "We should have only 1 job now for that chat")

	job := jobs[0]
	assert.Len(t, job.Tags(), 1, "Only one tag should be attached to the job")
	assert.Equal(t, job.Tags()[0], chat.ChatId, "The first and only tag should be equal to our chatid")

	_ = scheduler.RemoveByTag(chat.ChatId)
}

func TestJobExecuteReturnError(t *testing.T) {
	feed := database.NewFeed("https://test.com/feed")

	oldParseUrl := parseUrl
	defer func() { parseUrl = oldParseUrl }()

	parseUrl = func(url string) (*gofeed.Feed, error) {
		return nil, errors.New("Failling...")
	}

	err := job(feed, database.NewChat("124", 123, []database.Feed{*feed}, "tag1", "tag2"))

	fmt.Println(err)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("Couldn't parsed test.com: Failling...").Error(), err.Error(), "Error message should be the same")
	}
}

func TestJobExecuteNoArticlesFound(t *testing.T) {
	feed := database.NewFeed("https://test.com/feed")

	item := &gofeed.Item{
		Title:     "Super Title for an Article",
		Published: "2023-06-04",
		Link:      "https://test.com/article-1",
		Categories: []string{
			"tag5", "tag8",
		},
		Authors: []*gofeed.Person{
			{
				Name: "Victor",
			},
		},
	}

	f := &gofeed.Feed{
		Title: "Super Test",
		Items: []*gofeed.Item{
			item,
		},
	}

	oldParseUrl := parseUrl
	defer func() { parseUrl = oldParseUrl }()

	parseUrl = func(url string) (*gofeed.Feed, error) {
		return f, nil
	}

	articleRepositoryMock := new(repository.MocksArticleRep)
	repository.Article = articleRepositoryMock

	chatRepositoryMock := new(repository.MocksChatRep)
	repository.Chat = chatRepositoryMock

	telegramServiceMock := new(service.MocksTelegramService)
	telegram = telegramServiceMock

	err := job(feed, database.NewChat("124", 123, []database.Feed{*feed}, "tag1", "tag2"))

	assert.Nil(t, err)

	articleRepositoryMock.AssertNotCalled(t, "Find")
	articleRepositoryMock.AssertNotCalled(t, "Create")

	telegramServiceMock.AssertNotCalled(t, "TelegramPostMessage")
	telegramServiceMock.AssertNotCalled(t, "TelegramUpdateTyping")
}

func TestJobExecuteArticleFound(t *testing.T) {
	feed := database.NewFeed("https://test.com/feed")

	item := &gofeed.Item{
		Title:     "Super Title for an Article",
		Published: "2023-06-04",
		Link:      "https://test.com/article-1",
		Categories: []string{
			"tag5", "tag8",
		},
		Authors: []*gofeed.Person{
			{
				Name: "Victor",
			},
		},
	}

	f := &gofeed.Feed{
		Title: "Super Test",
		Items: []*gofeed.Item{
			item,
		},
	}

	chat := database.NewChat("124", 123, []database.Feed{*feed}, "tag1", "tag5")

	oldParseUrl := parseUrl
	defer func() { parseUrl = oldParseUrl }()

	parseUrl = func(url string) (*gofeed.Feed, error) {
		return f, nil
	}

	articleRepositoryMock := new(repository.MocksArticleRep)
	repository.Article = articleRepositoryMock

	chatRepositoryMock := new(repository.MocksChatRep)
	repository.Chat = chatRepositoryMock

	// Article repository return nil on Find
	articleRepositoryMock.On("Find", "title", item.Title).Return(nil)

	articleExpected := database.NewArticle(
		item.Title,
		item.Published,
		item.Link,
		f.Title,
		item.Authors[0].Name,
		[]string{"tag5"},
		item.Categories...,
	)

	// Article repository return an article on Create
	articleRepositoryMock.On("Create",
		item.Title,
		item.Published,
		item.Link,
		f.Title,
		item.Authors[0].Name,
		[]string{"tag5"},
		item.Categories,
	).Return(articleExpected)

	telegramServiceMock := new(service.MocksTelegramService)
	telegram = telegramServiceMock

	telegramServiceMock.On("TelegramUpdateTyping", chat.ChatId, true).Return()
	telegramServiceMock.On("TelegramUpdateTyping", chat.ChatId, false).Return()
	telegramServiceMock.On("TelegramPostMessage", chat.ChatId, mock.Anything).Return()

	chatRepositoryMock.On("UpdateFeedCheckForUrl", feed.Url, 1, chat.ChatId).Return(true)

	err := job(feed, chat)

	assert.Nil(t, err)

	articleRepositoryMock.AssertCalled(t, "Find", "title", item.Title)
	articleRepositoryMock.AssertCalled(t, "Create", item.Title, item.Published, item.Link, f.Title, item.Authors[0].Name, []string{"tag5"}, item.Categories)

	telegramServiceMock.AssertCalled(t, "TelegramPostMessage", chat.ChatId, mock.Anything)
	telegramServiceMock.AssertCalled(t, "TelegramUpdateTyping", chat.ChatId, mock.Anything)

	chatRepositoryMock.AssertCalled(t, "UpdateFeedCheckForUrl", feed.Url, 1, chat.ChatId)
}

func TestMainJobNoNewFeed(t *testing.T) {
	feed := database.NewFeed("https://google.com")
	chat := database.NewChat("1234", 123, []database.Feed{
		*feed,
	})
	_, _ = scheduler.Every(1).Tag(chat.ChatId).Week().Do(func() {
		fmt.Println("Test job")
	})

	chatRepositoryMock := new(repository.MocksChatRep)
	repository.Chat = chatRepositoryMock

	chatRepositoryMock.On("GetAll").Return([]*database.Chat{chat})

	jobs, _ := scheduler.FindJobsByTag(chat.ChatId)
	assert.Len(t, jobs, 1, "We should have only 1 job")

	mainJob()

	jobs, _ = scheduler.FindJobsByTag(chat.ChatId)
	assert.Len(t, jobs, 1, "We should have only 1 job even after the main job")

	_ = scheduler.RemoveByTag(chat.ChatId)
}

func TestMainJobNoJobFoundForChat(t *testing.T) {
	feed := database.NewFeed("https://google.com")
	chat := database.NewChat("1234", 123, []database.Feed{
		*feed,
	})

	chatRepositoryMock := new(repository.MocksChatRep)
	repository.Chat = chatRepositoryMock

	chatRepositoryMock.On("GetAll").Return([]*database.Chat{chat})

	_, err := scheduler.FindJobsByTag(chat.ChatId)
	assert.Error(t, err, "gocron: no jobs found with given tag")

	mainJob()

	jobs, _ := scheduler.FindJobsByTag(chat.ChatId)
	assert.Len(t, jobs, 1, "We should have only 1 job")

	_ = scheduler.RemoveByTag(chat.ChatId)
}

func TestMainJobAddNewFeed(t *testing.T) {
	feed := database.NewFeed("https://google.com")
	feed2 := database.NewFeed("https://google2.com")
	chat := database.NewChat("1234", 123, []database.Feed{
		*feed, *feed2,
	})

	_, _ = scheduler.Every(1).Tag(chat.ChatId).Week().Do(func() {
		fmt.Println("Test job")
	})

	chatRepositoryMock := new(repository.MocksChatRep)
	repository.Chat = chatRepositoryMock

	chatRepositoryMock.On("GetAll").Return([]*database.Chat{chat})

	jobs, _ := scheduler.FindJobsByTag(chat.ChatId)
	assert.Len(t, jobs, 1, "We should have only 1 job")

	mainJob()

	jobs, _ = scheduler.FindJobsByTag(chat.ChatId)
	assert.Len(t, jobs, 2, "We should have only 2 job")

	_ = scheduler.RemoveByTag(chat.ChatId)
}
