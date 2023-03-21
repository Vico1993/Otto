# Otto

## Setup

1. Need a telegram bot, look at the great telegram [documentation](https://core.telegram.org/bots)

2. Create a `.env` file with these keys:

```bash
# Bot token given by the @BotFather
TELEGRAM_BOT_TOKEN=<TELEGRAM_BOT>

# Used by the bot to know which chat to interact with
TELEGRAM_USER_CHAT_ID=<CHAT_ID>

#MONGO srv url you want to use with username + password for the connection
MONGO_URI=<URI>
```

## List of potential feed

TECH

-   https://techcrunch.com/feed/
-   https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml
-   https://rss.nytimes.com/services/xml/rss/nyt/PersonalTech.xml
-   https://dev.to/rss

CRYPTO

-   https://feeds.feedburner.com/CoingeckoBuzz
-   https://coinjournal.net/feeds/
-   https://medium.com/feed/tag/crypto
-   https://medium.com/feed/tag/tech

MONEY

-   https://rss.nytimes.com/services/xml/rss/nyt/YourMoney.xml
