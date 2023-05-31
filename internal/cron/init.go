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
)

var scheduler = gocron.NewScheduler(time.UTC)
var telegram = service.NewTelegramService()

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
}

// Will setup cron job for that chat
func SetupCronForChat(chat *database.Chat) {
	resetCronForChatId(chat)

	startJobForChat(chat)
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
	for _, feedurl := range chat.Feeds {
		// Copy val to be sure it's not overrited with the next iteration
		rul := feedurl.Url
		url, _ := url.Parse(rul)

		// Start at different time to avoid parsing all feed at the same time
		when := getDelay(len(chat.Feeds)) * n

		_, err := scheduler.Every(1).
			Hour().
			Tag(chat.ChatId).
			StartAt(time.Now().Add(time.Duration(when) * time.Minute)).
			Do(func() {
				err := parsedFeed(rul, chat.Tags)
				if err != nil {
					telegram.TelegramPostMessage("Couldn't checked: *" + url.Host + "*-> _" + err.Error() + "_")
					return
				}

				// Update last time checked
				repository.Chat.SetLastTimeCheckForUrl(rul, chat)
			})

		if err != nil {
			fmt.Println("Couldn't initiate the cron for: " + url.Host + " - " + err.Error())
		}

		n += 1
	}
}

// Calculate the delay between each job base on the number of feed
// Each feed need to be check once an hour
func getDelay(numberOfFeed int) int {
	return int(math.Round(float64(60 / numberOfFeed)))
}
