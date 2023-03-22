package feed

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/utils"
	"github.com/mmcdole/gofeed"
)

// Tags interested in
var tags []string = []string{
	"btc",
	"bitcoin",
	"vechain",
	"apple",
	"aapl",
	"finance",
	"crypto",
	"crypto.com",
	"cro",
	"banks",
	"binance",
}

// Base of feed need to look at
var listOfFeeds []string = []string{
	"https://techcrunch.com/feed/",
	"https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml",
	"https://rss.nytimes.com/services/xml/rss/nyt/PersonalTech.xml",
	"https://dev.to/rss",
	"https://feeds.feedburner.com/CoingeckoBuzz",
	"https://coinjournal.net/news/feed/",
	"https://coinjournal.net/news/category/events/feed/",
	"https://medium.com/feed/tag/crypto",
	"https://medium.com/feed/tag/tech",
	"https://rss.nytimes.com/services/xml/rss/nyt/YourMoney.xml",
}

// Parsed one RSS feed to extract some information
func parsedFeed(uri string) error {
	url, _ := url.Parse(uri)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(uri)
	if err != nil {
		return errors.New("Couldn't parsed " + url.Host + ": " + err.Error())
	}

	fmt.Println("Start parsing : " + url.Host)
	start := time.Now()

	var newItem int = 0
	for _, item := range feed.Items {
		article := repository.FindArticleByTitle(item.Title)

		if article != nil {
			continue
		}

		repository.CreateArticle(item.Title, item.Published, item.Link, feed.Title, item.Categories...)
		newItem += 1
	}

	if newItem > 1 {
		elapsed := time.Since(start)
		utils.TelegramPostMessage("Insert " + strconv.Itoa(newItem) + " new articles, took me : " + elapsed.String())
	} else {
		fmt.Println("Nothing to aggregated")
	}

	return nil
}

// Parsed all articles from the constant
func PullNewArticles() {
	listOfFeeds := getList()

	fmt.Println("Starting aggregeting data, " + strconv.Itoa(len(listOfFeeds)) + " of feeds to analyzed")
	start := time.Now()

	for _, feed := range listOfFeeds {
		err := parsedFeed(feed)
		if err != nil {
			utils.TelegramPostMessage(err.Error())
		}
	}

	elapsed := time.Since(start)
	utils.TelegramPostMessage("Done aggregating, took me : " + elapsed.String())
}

// Return list of feed to watch
func getList() []string {
	return append(
		buildMediumFeedBasedOnTag(),
		listOfFeeds...,
	)
}
