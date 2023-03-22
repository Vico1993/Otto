package feed

import (
	"fmt"
	"strconv"
	"time"
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
	"ethereum",
	"eth",
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

// Parsed all articles from the constant
func PullNewArticles() {
	listOfFeeds := getList()

	fmt.Println("Starting aggregeting data, " + strconv.Itoa(len(listOfFeeds)) + " of feeds to analyzed")
	start := time.Now()

	for _, feed := range listOfFeeds {
		err := parsedFeed(feed)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	elapsed := time.Since(start)
	fmt.Println("Done aggregating, took me : " + elapsed.String())
}

// Return list of feed to watch
func getList() []string {
	return append(
		buildMediumFeedBasedOnTag(),
		listOfFeeds...,
	)
}
