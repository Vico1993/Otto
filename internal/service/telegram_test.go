package service

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
		chatId:  "123",
		baseUrl: server.URL,
	}

	s.TelegramPostMessage("Test message")
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
		chatId:  "123",
		baseUrl: server.URL,
	}

	s.TelegramUpdateTyping(true)
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
		chatId:  "123",
		baseUrl: server.URL,
	}

	s.TelegramUpdateTyping(false)
}

func TestBuildData(t *testing.T) {
	data := buildData(map[string]string{
		"test": "foo",
	})

	assert.True(t, data.Has("test"), "The key foo should be set to true")
	assert.Equal(t, data.Encode(), "test=foo", "The key test should be equal to foo, and only 1 key should be set")
}
