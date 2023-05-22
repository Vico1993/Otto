package database

import "time"

type Feed struct {
	Url            string
	ChatId         string
	CreatedAt      time.Time
	LastTimeParsed time.Time
}

func NewFeed(
	url string,
	chatId string,
) *Feed {
	feed := Feed{}
	feed.Url = url
	feed.ChatId = chatId
	feed.CreatedAt = time.Now()

	return &feed
}
