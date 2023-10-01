package cron

import (
	"errors"
	"net/url"
	"strings"

	textrank "github.com/DavidBelicza/TextRank/v2"
	v2 "github.com/Vico1993/Otto/internal/repository/v2"
	"github.com/Vico1993/Otto/internal/utils"
	"github.com/mmcdole/gofeed"
)

// Extracting the Parsing url to a new var
var parseUrl = func(url string) (*gofeed.Feed, error) {
	return gofeedParser.ParseURL(url)
}

var (
	// Default Rule for parsing.
	rule = textrank.NewDefaultRule()
	// Default Language for filtering stop words.
	language = textrank.NewDefaultLanguage()
	// Default algorithm for ranking text.
	algorithmDef = textrank.NewDefaultAlgorithm()
)

type parser struct {
	url  string
	tags []string
}

// Parse an url to retrieve a list of articles matching the list of tags
// Will return an error OR the list of article found
func (p *parser) execute(articleRepository v2.IArticleRepository, feedId string) error {
	url, _ := url.Parse(p.url)

	feed, err := parseUrl(p.url)
	if err != nil {
		return errors.New("Couldn't parsed " + url.Host + ": " + err.Error())
	}

	for _, item := range feed.Items {
		item := item

		itemTags := item.Categories
		if len(itemTags) == 0 {
			itemTags = p.findTagFromTitle(item.Title)
		}

		// Looking into the DB to find if it's a new article...
		articleFound := articleRepository.GetByTitle(item.Title)
		if articleFound != nil {
			continue
		}

		author := "Unknown"
		if len(item.Authors) > 0 {
			author = item.Authors[0].Name
		}

		articleRepository.Create(
			feedId,
			item.Title,
			feed.Link,
			author,
			item.Link,
			itemTags,
		)
	}

	return nil
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

// Extract important word from the title
func (p *parser) findTagFromTitle(title string) []string {
	// TextRank object
	tr := textrank.NewTextRank()
	// Add text.
	tr.Populate(title, language, rule)
	// Run the ranking.
	tr.Ranking(algorithmDef)

	// Get all words order by weight.
	words := textrank.FindSingleWords(tr)

	var tags []string
	for _, word := range words {
		tags = append(tags, word.Word)
	}

	return tags
}
