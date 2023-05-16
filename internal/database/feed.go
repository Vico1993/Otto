package database

import "time"

type Feed struct {
	Url            string
	ChatId         string
	Added          time.Time
	LastTimeParsed time.Time
}

func NewFeed(
	url string,
	chatId string,
) *Feed {
	feed := Feed{}
	feed.Url = url
	feed.ChatId = chatId
	feed.Added = time.Now()

	return &feed
}
