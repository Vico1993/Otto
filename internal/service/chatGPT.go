package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type IChatGPTService interface {
	Ask(query string) *GptAskResponse
}

type ChatGPTService struct {
	baseUrl string
	token   string
}

func NewChatGPTService() IChatGPTService {
	return &ChatGPTService{
		baseUrl: "https://api.openai.com/v1/",
		token:   os.Getenv("OPENAI_TOKEN"),
	}
}

// build request to open ai api
func (s *ChatGPTService) buildRequest(body []byte, path string) *http.Request {
	req, err := http.NewRequest(
		"POST",
		s.baseUrl+path,
		nil,
	)

	if err != nil {
		fmt.Println("Error building request:", err.Error())
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	return req
}

// Ask question to ChatGPT
func (s *ChatGPTService) Ask(query string) *GptAskResponse {
	body := s.buildReqBody(query)
	if body == nil {
		return nil
	}

	req := s.buildRequest(body, "chat/completions")

	// Create a new HTTP client
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed")
		return nil
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response GptAskResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		fmt.Println("Response parse error", err.Error())

		return nil
	}

	return &response
}

// Marshal the user object into a JSON-encoded byte slice
func (s *ChatGPTService) buildReqBody(query string) []byte {
	request := newGPTRequest([]GptMessage{
		{
			Role:    "user",
			Content: query,
		},
	})

	body, err := json.Marshal(request)

	if err != nil {
		fmt.Println("Request parse error")
		return nil
	}

	return body
}

// Message sent by ChatGpt
type GptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request structure for Chat Gpt
type gptRequest struct {
	Model       string       `json:"model"`
	Messages    []GptMessage `json:"messages"`
	Temperature float32      `json:"temperature"`
	MaxTokens   int          `json:"max_tokens"`
}

// Create a New GPT Request
func newGPTRequest(messages []GptMessage) *gptRequest {
	return &gptRequest{
		Model:       "gpt-3.5-turbo-0301",
		Messages:    messages,
		Temperature: 1,
		MaxTokens:   1000,
	}
}

// Choices made by Chat Gpt
type gptChoice struct {
	Message      GptMessage `json:"message"`
	Index        int8       `json:"index"`
	FinishReason string     `json:"finish_reason"`
}

// Response schema from the Completions endpoint
type GptAskResponse struct {
	Id      string      `json:"id"`
	Object  string      `json:"object"`
	Created int         `json:"created"`
	Model   string      `json:"model"`
	Choices []gptChoice `json:"choices,omitempty"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
