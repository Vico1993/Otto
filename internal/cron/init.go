package cron

import (
	"fmt"
	"math"
	"net/url"
	"time"

	v2 "github.com/Vico1993/Otto/internal/repository/v2"
	"github.com/go-co-op/gocron"
	"github.com/mmcdole/gofeed"
)

var scheduler = gocron.NewScheduler(time.UTC)
var gofeedParser = gofeed.NewParser()

var feedsTag = "feed"
var mainTag = "main"

func Init() {
	_, err := scheduler.Every(1).Hour().Tag(mainTag).Do(func() {
		// Feeds
		checkResetFeed()
	})

	if err != nil {
		fmt.Println("Couldn't initiate the main job - " + err.Error())
	}
}

// Function that will check if need to reset job for feeds
func checkResetFeed() {
	fmt.Println("Checking Feeds")

	// Get All Feeds
	feedsList := v2.Feed.GetAll()

	jobs, err := scheduler.FindJobsByTag(feedsTag)
	// No job found but we have feeds
	// OR if we have more or less feed than before
	if (err != nil && len(feedsList) > 0) || (len(feedsList) != len(jobs)) {
		fmt.Println("Need Reset Feeds")
		feeds(feedsList)
	} else {
		fmt.Println("No need to reset")
	}
}

// Parsing feeds
func feeds(feeds []*v2.DBFeed) {
	feedsTag := "feed"

	err := scheduler.RemoveByTag(feedsTag)
	if err != nil {
		fmt.Println("FeedJob - Couldn't reset feed")
	}

	n := 1
	for _, feed := range feeds {
		// Copy val to be sure it's not overrited with the next iteration
		feed := feed
		url, _ := url.Parse(feed.Url)

		// Start at different time to avoid parsing all feed at the same time
		when := getDelay(len(feeds)) * n

		fmt.Println("FeedJob - Adding Job -> " + feed.Url)
		_, err := scheduler.Every(1).
			Hour().
			Tag(feedsTag).
			StartAt(time.Now().Add(time.Duration(when) * time.Minute)).
			Do(func() {
				fmt.Println("FeedJob - Start : " + feed.Url)

				parser := &parser{
					url:  feed.Url,
					tags: []string{},
				}

				err := parser.execute(v2.Article, feed.Id)

				if err != nil {
					fmt.Println("FeedJob - Error Parsing feed id " + feed.Id + " -> " + feed.Url + " : " + err.Error())
				}

				fmt.Println("FeedJob - End : " + feed.Url)
			})

		if err != nil {
			fmt.Println("FeedJob - Error initiate the cron for: " + url.Host + " - " + err.Error())
		}

		n += 1
	}
}

// Calculate the delay between each job base on the number of feed
// Each feed need to be check once an hour
func getDelay(numberOfFeed int) int {
	return int(math.Round(float64(60 / numberOfFeed)))
}
