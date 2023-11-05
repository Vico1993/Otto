package service

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
		assert.Equal(t, "/sendMessage", r.URL.String(), "Unexpected URL")

		data := url.Values{}
		data.Set("chat_id", "123")
		data.Set("text", "Test message")

		assert.Equal(t, data.Encode(), "chat_id=123&text=Test+message", "Body is not matching the expected body")

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := TelegramService{
		baseUrl: server.URL,
	}

	s.TelegramPostMessage("123", "", "Test message")
}

func TestTelegramPostMessageWithThreadId(t *testing.T) {
	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/sendMessage", r.URL.String(), "Unexpected URL")

		data := url.Values{}
		data.Set("chat_id", "123")
		data.Set("reply_to_message_id", "124")
		data.Set("text", "Test message")

		assert.Equal(t, data.Encode(), "chat_id=123&reply_to_message_id=124&text=Test+message", "Body is not matching the expected body")

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := TelegramService{
		baseUrl: server.URL,
	}

	s.TelegramPostMessage("123", "124", "Test message")
}

func TestTelegramSetTypingToTrue(t *testing.T) {
	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/sendChatAction", r.URL.String(), "Unexpected URL")

		data := url.Values{}
		data.Set("chat_id", "123")
		data.Set("action", "typing")

		assert.Equal(t, data.Encode(), "action=typing&chat_id=123", "Body is not matching the expected body")

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := TelegramService{
		baseUrl: server.URL,
	}

	s.TelegramUpdateTyping("123", true)
}

func TestTelegramSetTypingToFalse(t *testing.T) {
	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/sendChatAction", r.URL.String(), "Unexpected URL")

		data := url.Values{}
		data.Set("chat_id", "123")

		assert.Equal(t, data.Encode(), "chat_id=123", "Body is not matching the expected body")

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := TelegramService{
		baseUrl: server.URL,
	}

	s.TelegramUpdateTyping("123", false)
}

func TestBuildData(t *testing.T) {
	data := buildData(map[string]string{
		"test": "foo",
	})

	assert.True(t, data.Has("test"), "The key foo should be set to true")
	assert.Equal(t, data.Encode(), "test=foo", "The key test should be equal to foo, and only 1 key should be set")
}

func TestNewTelegramService(t *testing.T) {
	os.Setenv("TELEGRAM_BOT_TOKEN", "FOO")

	service := NewTelegramService()

	assert.Equal(t, service.GetBaseUrl(), "https://api.telegram.org/botFOO", "Url is should include FOO at the end")
}
