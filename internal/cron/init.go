package cron

import (
	"fmt"
	"math"
	"time"

	"github.com/Vico1993/Otto/internal/service"
	"github.com/go-co-op/gocron"
	"github.com/mmcdole/gofeed"
)

var Scheduler = gocron.NewScheduler(time.UTC)
var gofeedParser = gofeed.NewParser()
var telegram = service.NewTelegramService()

var mainTag = "main"
var feedsTag = "feed"
var chatsTag = "chat"

func Init() {
	_, err := Scheduler.Every(1).Hour().Tag(mainTag).Do(func() {
		// Feeds
		checkFeed()
		checkChat()
	})

	if err != nil {
		fmt.Println("Couldn't initiate the main job - " + err.Error())
	}
}

// Calculate the delay between each job base on the number of feed
// Each feed need to be check once an hour
func getDelay(numberOfFeed int) int {
	return int(math.Round(float64(60 / numberOfFeed)))
}
