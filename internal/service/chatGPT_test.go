package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAskValidQuery(t *testing.T) {
	os.Setenv("OPENAI_TOKEN", "FOO")

	message := "Test query"

	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/chat/completions", r.URL.String(), "Unexpected URL")

		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Content type should be JSON")
		assert.Equal(t, "Bearer FOO", r.Header.Get("Authorization"), "Baerer token should be set")

		body, _ := r.GetBody()

		// Read the response body into a byte slice
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			panic(err)
		}

		var request gptRequest
		_ = json.Unmarshal(bodyBytes, &request)

		assert.Len(t, request.Messages, 1, "Request should contain 1 message")
		assert.Equal(t, request.Messages[0].Content, message, "Message should be equal to the content")

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := NewChatGPTService()
	s.Ask(message)
}
