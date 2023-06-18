package service

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ITelegramService interface {
	TelegramPostMessage(chatId string, text string)
	TelegramUpdateTyping(chatId string, val bool)
	GetBaseUrl() string
}

type TelegramService struct {
	baseUrl string
}

// TODO: Update this service to have dynamic TELEGRAM_USER_CHAT_ID
func NewTelegramService() ITelegramService {
	return &TelegramService{
		// Replace token in the URL
		baseUrl: "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

func (s *TelegramService) TelegramPostMessage(chatId string, text string) {
	data := buildData(map[string]string{
		"chat_id":    chatId,
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

func (s *TelegramService) TelegramUpdateTyping(chatId string, val bool) {
	data := buildData(map[string]string{
		"chat_id": chatId,
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

func (s *TelegramService) GetBaseUrl() string {
	return s.baseUrl
}

func buildData(params map[string]string) url.Values {
	data := url.Values{}

	for key, val := range params {
		data.Set(key, val)
	}

	return data
}
