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

func TestNotValidFeedId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockFeedRepository := new(repository.MocksFeedRepository)
	repository.Feed = mockFeedRepository

	mockFeedRepository.On("GetOne", "foo").Return(nil)

	r := gin.Default()
	r.Use(ValidFeed())

	r.GET("/feeds/:feedid", func(c *gin.Context) {})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest("GET", "/feeds/foo", nil))

	mockFeedRepository.AssertCalled(t, "GetOne", "foo")

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "Invalid the status code "+strconv.Itoa(http.StatusBadRequest))
}
