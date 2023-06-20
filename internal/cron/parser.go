package cron

import (
	"errors"
	"net/url"
	"strings"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/utils"
	"github.com/mmcdole/gofeed"
)

// Extracting the Parsing url to a new var
var parseUrl = func(url string) (*gofeed.Feed, error) {
	return gofeedParser.ParseURL(url)
}

type parser struct {
	url  string
	tags []string
}

type parseResult struct {
	FeedTitle string
	Articles  []*database.Article
}

// Parse an url to retrieve a list of articles matching the list of tags
// Will return an error OR the list of article found
func (p *parser) execute(articleRepository repository.IArticleRepository) (*parseResult, error) {
	url, _ := url.Parse(p.url)

	feed, err := parseUrl(p.url)
	if err != nil {
		return nil, errors.New("Couldn't parsed " + url.Host + ": " + err.Error())
	}

	articles := []*database.Article{}
	for _, item := range feed.Items {
		// If the category doesn't match with the interest tags
		match := p.isCategoriesAndTagsMatch(item.Categories)
		if len(match) == 0 {
			continue
		}

		// Looking into the DB to find if it's a new article...
		articleFound := articleRepository.Find("title", item.Title)
		if articleFound != nil {
			continue
		}

		articles = append(articles,
			articleRepository.Create(
				item.Title,
				item.Published,
				item.Link,
				feed.Title,
				item.Authors[0].Name,
				match,
				item.Categories...,
			),
		)
	}

	return &parseResult{
		FeedTitle: feed.Title,
		Articles:  articles,
	}, nil
}

// find if a list of categories is in tags
// and return the list of tags present in the categories
func (p *parser) isCategoriesAndTagsMatch(categories []string) []string {
	match := []string{}
	for _, category := range categories {
		if utils.InSlice(strings.ToLower(category), p.tags) {
			match = append(match, strings.ToLower(category))
		}
	}

	return match
}
