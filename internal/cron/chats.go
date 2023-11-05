package cron

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Vico1993/Otto/internal/repository"
)

// Function that will check if need to reset job for chats
func checkChat() {
	fmt.Println("Checking Chats")

	// Get All Chats
	chatsList := repository.Chat.GetAll()

	jobs, err := Scheduler.FindJobsByTag(chatsTag)
	// No job found but we have chats
	// OR if we have more or less chats than before
	if (err != nil && len(chatsList) > 0) || (len(chatsList) != len(jobs)) {
		fmt.Println("Need Reset Chats")
		chats(chatsList)
	} else {
		fmt.Println("No need to reset")
	}
}

func chats(chatsList []*repository.DBChat) {
	err := Scheduler.RemoveByTag(chatsTag)
	if err != nil {
		fmt.Println("ChatJob - Couldn't reset chats")
	}

	n := 1
	for _, chat := range chatsList {
		// Copy val to be sure it's not overrited with the next iteration
		chat := chat

		// Start at different time to avoid parsing all feed at the same time
		when := getDelay(len(chatsList)) * n
		fmt.Println("ChatJob - Adding Job -> " + chat.TelegramChatId)
		_, err := Scheduler.Every(1).
			Hour().
			Tag(chatsTag).
			StartAt(time.Now().Add(time.Duration(when) * time.Minute)).
			Do(func() {
				fmt.Println("ChatJob - Start : " + chat.TelegramChatId)

				articles := repository.Article.GetByChatAndTime(chat.Id)
				parsedArticles(articles, chat)

				// Update the time parsed
				repository.Chat.UpdateParsed(chat.Id)

				fmt.Println("ChatJob - End")
			})

		if err != nil {
			fmt.Println("ChatJob - Error initiate the cron for: " + chat.TelegramChatId + " - " + err.Error())
		}

		n += 1
	}
}

// Parsed list of Articles, and found if needed to send notifications
func parsedArticles(articles []*repository.DBArticle, chat *repository.DBChat) {
	nbArticle := len(articles)
	if nbArticle == 0 {
		return
	}

	fmt.Println("ChatJob - Articles found:" + strconv.Itoa(nbArticle))

	telegram.TelegramUpdateTyping(chat.TelegramChatId, true)
	for _, article := range articles {
		article := article

		host := article.Source
		u, err := url.Parse(article.Source)
		if err == nil {
			host = u.Host
		}

		threaId := ""
		if chat.TelegramThreadId != nil {
			threaId = *chat.TelegramThreadId
		}

		telegram.TelegramPostMessage(
			chat.TelegramChatId,
			threaId,
			BuildMessage(
				article.Title,
				host,
				article.Author,
				article.Tags,
				article.Link,
			),
		)
	}
	telegram.TelegramUpdateTyping(chat.TelegramChatId, false)
}
