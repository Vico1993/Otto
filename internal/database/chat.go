package database

import "time"

type Feed struct {
	Url            string
	Tags           []string
	ArticleFound   int
	CreatedAt      time.Time
	LastTimeParsed time.Time
}

func NewFeed(url string) *Feed {
	feed := Feed{}
	feed.Url = url
	feed.ArticleFound = 0
	feed.CreatedAt = time.Now()

	return &feed
}

type Chat struct {
	ChatId    string
	UserId    int64
	Tags      []string
	Feeds     []Feed
	CreatedAt time.Time
}

func NewChat(
	chatid string,
	userid int64,
	feeds []Feed,
	tags ...string,
) *Chat {
	return &Chat{
		ChatId:    chatid,
		UserId:    userid,
		Feeds:     feeds,
		Tags:      tags,
		CreatedAt: time.Now(),
	}
}
