package cron

import (
	"fmt"
	"time"

	"github.com/Vico1993/Otto/internal/feed"
	"github.com/go-co-op/gocron"
)

func Init() {
	s := gocron.NewScheduler(time.UTC)

	// Check every hours
	_, err := s.Every(1).Hour().Do(func() {
		feed.PullNewArticles()
	})
	if err != nil {
		fmt.Println("Couldn't run the Check Message")
	}

	// Start executing cron Async
	// For now..
	s.StartBlocking()

	fmt.Println("Cron job initied!")
}
