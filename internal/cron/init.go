package cron

import (
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/Vico1993/Otto/internal/feed"
	"github.com/Vico1993/Otto/internal/utils"
	"github.com/go-co-op/gocron"
)

func Init() {
	s := gocron.NewScheduler(time.UTC)
	listOfFeed := feed.GetList()

	n := 0
	for _, feedurl := range listOfFeed {
		url, _ := url.Parse(feedurl)

		// Start at different time to avoid parsing all feed at the same time
		when := getDelay(listOfFeed) * n
		_, _ = s.Every(1).Hour().Tag(url.Host).StartAt(time.Now().Add(time.Duration(when) * time.Second)).Do(func() {
			err := feed.ParsedFeed(feedurl)
			if err != nil {
				utils.TelegramPostMessage("Couldn't checked: " + url.Host + "-> " + err.Error())
			}
		})

		n += 1
	}

	fmt.Println("Cron job initied!")

	// Start executing cron Async
	// For now..
	s.StartBlocking()
}

func getDelay(listOfFeed []string) int {
	return int(math.Round(float64(60 / len(listOfFeed))))
}
