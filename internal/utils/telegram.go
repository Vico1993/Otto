package utils

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

var telegramBaseURL = "https://api.telegram.org/bot<TOKEN>"

// Post a message in the telegram chat
func TelegramPostMessage(text string) {
	uri := strings.Replace(telegramBaseURL+"/sendMessage", "<TOKEN>", os.Getenv("TELEGRAM_BOT_TOKEN"), 1)

	data := url.Values{}
	data.Set("chat_id", os.Getenv("TELEGRAM_USER_CHAT_ID"))
	data.Set("text", text)

	_, err := http.Post(uri, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
}

func TelegramUpdateTyping(val bool) {
	uri := strings.Replace(telegramBaseURL+"/sendChatAction", "<TOKEN>", os.Getenv("TELEGRAM_BOT_TOKEN"), 1)

	data := url.Values{}
	data.Set("chat_id", os.Getenv("TELEGRAM_USER_CHAT_ID"))
	if val {
		data.Set("action", "typing")
	}

	_, err := http.Post(uri, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
}
