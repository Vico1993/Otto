package cron

import (
	"errors"
	"net/url"
	"strings"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/utils"
	"github.com/mmcdole/gofeed"
)

// Parsed one RSS feed to extract some information
func fetchFeed(uri string, tags []string) (int, error) {
	url, _ := url.Parse(uri)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(uri)
	if err != nil {
		return 0, errors.New("Couldn't parsed " + url.Host + ": " + err.Error())
	}

	articleFound := 0
	for _, item := range feed.Items {
		// If the category doesn't match with the interest tags
		match := isCategoriesAndTagsMatch(item.Categories, tags)
		if len(match) == 0 {
			continue
		}

		// Looking into the DB to find if it's a new article...
		article := repository.Article.Find("title", item.Title)
		if article != nil {
			continue
		}

		repository.Article.Create(
			item.Title,
			item.Published,
			item.Link,
			feed.Title,
			item.Categories...,
		)

		articleFound += 1

		telegram.TelegramUpdateTyping(true)
		telegram.TelegramPostMessage(
			BuildMessage(
				item.Title,
				feed.Title,
				item.Author.Name,
				match,
				item.Link,
			),
		)
		telegram.TelegramUpdateTyping(false)
	}

	return articleFound, nil
}

// find if a list of categories is in tags
// and return the list of tags present in the categories
func isCategoriesAndTagsMatch(categories []string, tags []string) []string {
	match := []string{}
	for _, category := range categories {
		if utils.InSlice(strings.ToLower(category), tags) {
			match = append(match, strings.ToLower(category))
		}
	}

	return match
}
