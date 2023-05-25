package cron

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"time"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/service"
	"github.com/go-co-op/gocron"
)

// TODO: Move this code into something CHAT_ID specific... maybe at /init ?
var scheduler = gocron.NewScheduler(time.UTC)
var telegram = service.NewTelegramService()

// Initialisation of the cronjob at the start of the program
func Init() {
	chatIds := repository.Feed.GetDistinctChatId()

	// Load all chatId and set each cron
	for _, chatId := range chatIds {
		ResetCronForChatId(fmt.Sprint(chatId))
	}

	// Start executing cron Async
	// For now..
	scheduler.StartAsync()
}

// Calculate the delay between each job base on the number of feed
// Each feed need to be check once an hour
func getDelay(numberOfFeed int) int {
	return int(math.Round(float64(60 / numberOfFeed)))
}

// Will load all feed for one chat id
// Calculate delay
// Start the cron again
func ResetCronForChatId(chatId string) {
	// Remove the previous one if any
	err := scheduler.RemoveByTag(chatId)
	if err != nil {
		fmt.Println("Couldn't clean tag " + chatId + " - " + err.Error())
	}

	listOfFeed := repository.Feed.FindByChatId(chatId)

	n := 1
	for _, feedurl := range listOfFeed {
		// Copy val to be sure it's not overrited with the next iteration
		rul := feedurl.Url
		url, _ := url.Parse(rul)

		// Start at different time to avoid parsing all feed at the same time
		when := getDelay(len(listOfFeed)) * n

		_, err := scheduler.Every(1).
			Hour().
			Tag(chatId).
			StartAt(time.Now().Add(time.Duration(when) * time.Minute)).
			Do(func() {
				err := parsedFeed(rul)
				if err != nil {
					telegram.TelegramPostMessage("Couldn't checked: *" + url.Host + "*-> _" + err.Error() + "_")
					return
				}

				// Update last time checked
				repository.Feed.SetLastTimeCheck(rul, os.Getenv("TELEGRAM_USER_CHAT_ID"))
			})

		if err != nil {
			fmt.Println("Couldn't initiate the cron for: " + url.Host + " - " + err.Error())
		}

		n += 1
	}
}
