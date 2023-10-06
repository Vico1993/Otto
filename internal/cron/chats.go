package cron

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/utils"
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
		matched := isCategoriesAndTagsMatch(chat.Tags, article.Tags)
		if len(matched) == 0 {
			continue
		}

		host := article.Source
		u, err := url.Parse(article.Source)
		if err == nil {
			host = u.Host
		}

		telegram.TelegramPostMessage(
			chat.TelegramChatId,
			BuildMessage(
				article.Title,
				host,
				article.Author,
				matched,
				article.Link,
			),
		)
	}
	telegram.TelegramUpdateTyping(chat.TelegramChatId, false)
}

// find if a list of categories is in tags
// and return the list of tags present in the categories
func isCategoriesAndTagsMatch(chatTags []string, articleCategories []string) []string {
	match := []string{}
	for _, category := range chatTags {
		if utils.InSlice(strings.ToLower(category), articleCategories) {
			match = append(match, strings.ToLower(category))
		}
	}

	return match
}
