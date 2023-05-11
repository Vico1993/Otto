package service

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TelegramService struct {
	baseUrl string
	chatId  string
}

func NewTelegramService() *TelegramService {
	return &TelegramService{
		// Replace token in the URL
		baseUrl: "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN"),
		chatId:  os.Getenv("TELEGRAM_USER_CHAT_ID"),
	}
}

func (s *TelegramService) TelegramPostMessage(text string) {
	data := buildData(map[string]string{
		"chat_id":    s.chatId,
		"text":       text,
		"parse_mode": "markdown",
	})

	_, err := http.Post(
		s.baseUrl+"/sendMessage",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		panic(err)
	}
}

func (s *TelegramService) TelegramUpdateTyping(val bool) {
	data := buildData(map[string]string{
		"chat_id": s.chatId,
	})

	if val {
		data.Set("action", "typing")
	}

	_, err := http.Post(
		s.baseUrl+"/sendChatAction",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		panic(err)
	}
}

func buildData(params map[string]string) url.Values {
	data := url.Values{}

	for key, val := range params {
		data.Set(key, val)
	}

	return data
}
