package feed

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/service"
	"github.com/Vico1993/Otto/internal/utils"
	"github.com/mmcdole/gofeed"
)

// Parsed one RSS feed to extract some information
func ParsedFeed(uri string) error {
	url, _ := url.Parse(uri)
	telegram := service.NewTelegramService()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(uri)
	if err != nil {
		return errors.New("Couldn't parsed " + url.Host + ": " + err.Error())
	}

	log.Default().Println("Waiking up")

	for _, item := range feed.Items {
		// If the category doesn't match with the interest tags
		match := isCategoriesAndTagsMatch(item.Categories)
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

		// Exclude medium from notification
		if url.Host != "medium.com" {
			telegram.TelegramUpdateTyping(true)
			telegram.TelegramPostMessage(item.Link)
			telegram.TelegramPostMessage("Sounds interesting: #" + strings.Join(match, ", #"))
			telegram.TelegramUpdateTyping(false)
		}
	}

	return nil
}

// find if a list of categories is in tags
// and return the list of tags present in the categories
func isCategoriesAndTagsMatch(categories []string) []string {
	match := []string{}
	for _, category := range categories {
		if utils.InSlice(strings.ToLower(category), tags) {
			match = append(match, strings.ToLower(category))
		}
	}

	return match
}
