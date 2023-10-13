package cron

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Vico1993/Otto/internal/repository"
)

// Function that will check if need to reset job for feeds
func checkFeed() {
	fmt.Println("Checking Feeds")

	// Get All Feeds
	feedsList := repository.Feed.GetAllActive()

	jobs, err := Scheduler.FindJobsByTag(feedsTag)
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
func feeds(feeds []*repository.DBFeed) {
	err := Scheduler.RemoveByTag(feedsTag)
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
		_, err := Scheduler.Every(1).
			Hour().
			Tag(feedsTag).
			StartAt(time.Now().Add(time.Duration(when) * time.Minute)).
			Do(func() {
				fmt.Println("FeedJob - Start : " + feed.Url)

				parser := &parser{
					url:  feed.Url,
					tags: []string{},
				}

				err := parser.execute(repository.Article, feed.Id)

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
