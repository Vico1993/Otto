package feed

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
	"https://rss.nytimes.com/services/xml/rss/nyt/YourMoney.xml",
	"https://waxy.org/feed/",
	"https://news.ycombinator.com/rss",
}

// Return list of feed to watch
func GetList() []string {
	return append(
		buildMediumFeedBasedOnTag(),
		listOfFeeds...,
	)
}
