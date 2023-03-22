package utils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTelegramPostMessage(t *testing.T) {
	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bot/sendMessage", r.URL.String(), "Unexpected URL")

		data := url.Values{}
		data.Set("chat_id", os.Getenv("TELEGRAM_USER_CHAT_ID"))
		data.Set("text", "Test message")

		assert.Equal(t, data.Encode(), "chat_id=&text=Test+message", "Body is not matching the expected body")

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Set the mock server's URL as the base URL for the Telegram API
	telegramBaseURL = server.URL + "/bot<TOKEN>"

	// Call the function with a test message
	TelegramPostMessage("Test message")
}

func TestTelegramSetTypingToTrue(t *testing.T) {
	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bot/sendChatAction", r.URL.String(), "Unexpected URL")

		data := url.Values{}
		data.Set("chat_id", os.Getenv("TELEGRAM_USER_CHAT_ID"))
		data.Set("action", "typing")

		assert.Equal(t, data.Encode(), "action=typing&chat_id=", "Body is not matching the expected body")

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Set the mock server's URL as the base URL for the Telegram API
	telegramBaseURL = server.URL + "/bot<TOKEN>"

	TelegramUpdateTyping(true)
}

func TestTelegramSetTypingToFalse(t *testing.T) {
	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bot/sendChatAction", r.URL.String(), "Unexpected URL")

		data := url.Values{}
		data.Set("chat_id", os.Getenv("TELEGRAM_USER_CHAT_ID"))

		assert.Equal(t, data.Encode(), "chat_id=", "Body is not matching the expected body")

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Set the mock server's URL as the base URL for the Telegram API
	telegramBaseURL = server.URL + "/bot<TOKEN>"

	TelegramUpdateTyping(false)
}
