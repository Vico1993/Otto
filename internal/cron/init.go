package cron

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"time"

	"github.com/Vico1993/Otto/internal/feed"
	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/service"
	"github.com/go-co-op/gocron"
)

// TODO: Move this code into something CHAT_ID specific... maybe at /init ?
var scheduler = gocron.NewScheduler(time.UTC)

// Initialisation of the cronjob at the start of the program
func Init() {
	listOfFeed := repository.Feed.FindByChatId(os.Getenv("TELEGRAM_USER_CHAT_ID"))
	telegram := service.NewTelegramService()

	n := 1
	for _, feedurl := range listOfFeed {
		// Copy val to be sure it's not overrited with the next iteration
		rul := feedurl.Url
		url, _ := url.Parse(rul)

		// Start at different time to avoid parsing all feed at the same time
		when := getDelay(len(listOfFeed)) * n

		_, err := scheduler.Every(1).
			Hour().
			Tag(os.Getenv("TELEGRAM_USER_CHAT_ID")).
			StartAt(time.Now().Add(time.Duration(when) * time.Minute)).
			Do(func() {
				err := feed.ParsedFeed(rul)
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

	fmt.Println("Cron job initied!")

	// Start executing cron Async
	// For now..
	scheduler.StartAsync()
}

func getDelay(numberOfFeed int) int {
	return int(math.Round(float64(60 / numberOfFeed)))
}
