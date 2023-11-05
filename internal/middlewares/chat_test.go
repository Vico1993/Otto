package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNotValidChatId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mocksChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mocksChatRepository

	mocksChatRepository.On("GetByTelegramChatId", "foo").Return(nil)

	r := gin.Default()
	r.Use(ValidChat())

	r.GET("/chats/:chatid", func(c *gin.Context) {})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest("GET", "/chats/foo", nil))

	mocksChatRepository.AssertCalled(t, "GetByTelegramChatId", "foo")
	mocksChatRepository.AssertNotCalled(t, "GetByTelegramChatIdAndThreadId")

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "Invalid the status code "+strconv.Itoa(http.StatusBadRequest))
}

func TestValidChatIdButNotThreadId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mocksChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mocksChatRepository

	mocksChatRepository.On("GetByTelegramChatIdAndThreadId", "foo", "threadId").Return(nil)

	r := gin.Default()
	r.Use(ValidChat())

	r.GET("/chats/:chatid/:threadid", func(c *gin.Context) {})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest("GET", "/chats/foo/threadId", nil))

	mocksChatRepository.AssertCalled(t, "GetByTelegramChatIdAndThreadId", "foo", "threadId")
	mocksChatRepository.AssertNotCalled(t, "GetByTelegramChatId", "foo")

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "Invalid the status code "+strconv.Itoa(http.StatusBadRequest))
}
