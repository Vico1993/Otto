package database

import (
	"time"
)

type Article struct {
	Title     string
	Published string
	Parsed    time.Time
	Link      string
	Source    string
	Tags      []string
}

func NewArticle(title string, published string, link string, source string, tags ...string) *Article {
	return &Article{
		Title:     title,
		Published: published,
		Parsed:    time.Now(),
		Link:      link,
		Source:    source,
		Tags:      tags,
	}
}
