package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAskValidQuery(t *testing.T) {
	message := "Test query"

	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/chat/completions", r.URL.String(), "Unexpected URL")

		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Content type should be JSON")
		assert.Equal(t, "Bearer FOO", r.Header.Get("Authorization"), "Baerer token should be set")

		// Read the response body into a byte slice
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var request gptRequest
		err = json.Unmarshal(bodyBytes, &request)
		if err != nil {
			fmt.Println("TEST1")
			fmt.Println(err.Error())
		}

		assert.Len(t, request.Messages, 1, "Request should contain 1 message")
		assert.Equal(t, request.Messages[0].Content, message, "Message should be equal to the content")

		fakeResponse := GptAskResponse{
			Id:      "123",
			Object:  "test",
			Created: 12341,
			Model:   "model",
		}
		byteResponse, _ := json.Marshal(fakeResponse)

		// Respond with a success status code
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(byteResponse)
	}))
	defer server.Close()

	s := &ChatGPTService{
		baseUrl: server.URL,
		token:   "FOO",
		client:  &http.Client{},
	}
	_, err := s.Ask(message)

	assert.Nil(t, err, "Shouldn't have any error returned")
}

func TestServerReturn500error(t *testing.T) {
	message := "Test query"

	// Set up a mock server to receive the HTTP POST request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a success status code
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	s := &ChatGPTService{
		baseUrl: server.URL,
		token:   "FOO",
		client:  &http.Client{},
	}
	_, err := s.Ask(message)

	assert.NotNil(t, err, "Error shoud be nil")
}

func TestNewService(t *testing.T) {
	os.Setenv("OPENAI_TOKEN", "FOO")

	s := &ChatGPTService{
		baseUrl: "https://api.openai.com/v1",
		token:   "FOO",
		client:  &http.Client{},
	}

	assert.Equal(t, s, NewChatGPTService(), "Both service should be equal")
}
