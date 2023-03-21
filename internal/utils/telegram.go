package utils

import (
	"net/http"
	"os"
	"strings"
)

var telegramBaseURL = "https://api.telegram.org/bot<TOKEN>/sendMessage"

// Post a message in the telegram chat
func TelegramPostMessage(text string) {
	url := strings.Replace(telegramBaseURL, "<TOKEN>", os.Getenv("TELEGRAM_BOT_TOKEN"), 1)

	_, err := http.Post(url+"?chat_id="+os.Getenv("TELEGRAM_USER_CHAT_ID")+"&text="+text, "", nil)
	if err != nil {
		panic(err)
	}
}
