package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTelegramPostMessage(t *testing.T) {
	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "https://api.telegram.org/bot"+os.Getenv("TELEGRAM_BOT_TOKEN")+"/sendMessage?chat_id="+os.Getenv("TELEGRAM_USER_CHAT_ID")+"&text=Test+message", r.URL.String(), "Unexpected URL")

		if r.Body != http.NoBody {
			t.Errorf("unexpected request body: got %v, want %v", r.Body, http.NoBody)
		}

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Set the mock server's URL as the base URL for the Telegram API
	telegramBaseURL = server.URL + "/bot<TOKEN>"

	// Call the function with a test message
	TelegramPostMessage("Test message")

	// The test passes if the function sends an HTTP POST request to the mock server with the correct URL and body
}
