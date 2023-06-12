package cron

import (
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/service"
	"github.com/go-co-op/gocron"
	"github.com/mmcdole/gofeed"
)

var scheduler = gocron.NewScheduler(time.UTC)
var telegram = service.NewTelegramService()
var gofeedParser = gofeed.NewParser()

// Initialisation of the cronjob at the start of the program
func Init() {
	chats := repository.Chat.GetAll()

	// Load all chat and set each cron
	for _, chat := range chats {
		SetupCronForChat(
			chat,
		)
	}

	// Start executing cron Async
	// For now..
	scheduler.StartAsync()

	fmt.Println("Cron ready for all chats!!")
}

// Will setup cron job for that chat
func SetupCronForChat(chat *database.Chat) {
	fmt.Println("Start cleaning cron for: " + chat.ChatId)
	resetCronForChatId(chat)

	fmt.Println("Reinitilisation cron for: " + chat.ChatId)
	startJobForChat(chat)

	fmt.Println("Cron setup for: " + chat.ChatId)
}

// Delete all previous tasks in the cron link to that chat id
func resetCronForChatId(chat *database.Chat) {
	// Remove the previous one if any
	err := scheduler.RemoveByTag(chat.ChatId)
	if err != nil {
		fmt.Println("Couldn't clean tag " + chat.ChatId + " - " + err.Error())
	}
}

// Will calculate the delay between each feed, and add then to the scheduler
func startJobForChat(chat *database.Chat) {
	n := 1
	for _, feed := range chat.Feeds {
		// Copy val to be sure it's not overrited with the next iteration
		rul := feed.Url
		url, _ := url.Parse(rul)

		// Start at different time to avoid parsing all feed at the same time
		when := getDelay(len(chat.Feeds)) * n

		_, err := scheduler.Every(1).
			Hour().
			Tag(chat.ChatId).
			StartAt(time.Now().Add(time.Duration(when) * time.Minute)).
			Do(func() {
				err := job(feed, chat)
				if err != nil {
					telegram.TelegramPostMessage("Couldn't checked: *" + url.Host + "*-> _" + err.Error() + "_")
				}
			})

		if err != nil {
			fmt.Println("Couldn't initiate the cron for: " + url.Host + " - " + err.Error())
		}

		n += 1
	}
}

// Job to execute
func job(feed database.Feed, chat *database.Chat) error {
	parser := &parser{
		url:  feed.Url,
		tags: append(feed.Tags, chat.Tags...),
	}

	result, err := parser.execute(repository.Article)
	if err != nil {
		return err
	}

	if len(result.articles) == 0 {
		return nil
	}

	for _, article := range result.articles {
		telegram.TelegramUpdateTyping(true)
		telegram.TelegramPostMessage(
			BuildMessage(
				article.Title,
				result.feedTitle,
				article.Author,
				article.MatchingTags,
				article.Link,
			),
		)
		telegram.TelegramUpdateTyping(false)
	}

	// Update feed after check
	repository.Chat.UpdateFeedCheckForUrl(feed.Url, len(result.articles), chat)

	return nil
}

// Calculate the delay between each job base on the number of feed
// Each feed need to be check once an hour
func getDelay(numberOfFeed int) int {
	return int(math.Round(float64(60 / numberOfFeed)))
}
