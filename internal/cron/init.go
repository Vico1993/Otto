package cron

import (
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/Vico1993/Otto/internal/feed"
	"github.com/Vico1993/Otto/internal/service"
	"github.com/go-co-op/gocron"
)

func Init() {
	s := gocron.NewScheduler(time.UTC)
	listOfFeed := feed.GetList()
	telegram := service.NewTelegramService()

	n := 1
	for _, feedurl := range listOfFeed {
		// Copy val to be sure it's not overrited with the next iteration
		rul := feedurl

		url, _ := url.Parse(rul)

		// Start at different time to avoid parsing all feed at the same time
		when := getDelay(listOfFeed) * n

		_, err := s.Every(1).
			Hour().
			Tag(url.Host).
			StartAt(time.Now().Add(time.Duration(when) * time.Minute)).
			Do(func() {
				err := feed.ParsedFeed(rul)
				if err != nil {
					telegram.TelegramPostMessage("Couldn't checked: *" + url.Host + "*-> _" + err.Error() + "_")
				}
			})

		if err != nil {
			fmt.Println("Couldn't initiate the cron for: " + url.Host + " - " + err.Error())
		}

		n += 1
	}

	fmt.Println("Cron job initied!")

	// Start executing cron Async
	// For now..
	s.StartAsync()
}

func getDelay(listOfFeed []string) int {
	return int(math.Round(float64(60 / len(listOfFeed))))
}
