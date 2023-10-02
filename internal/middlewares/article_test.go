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

func TestNotValidArticleId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockArticleRepository := new(repository.MocksArticleRepository)
	repository.Article = mockArticleRepository

	mockArticleRepository.On("GetOne", "foo").Return(nil)

	r := gin.Default()
	r.Use(ValidArticle())

	r.GET("/articles/:articleid", func(c *gin.Context) {})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest("GET", "/articles/foo", nil))

	mockArticleRepository.AssertCalled(t, "GetOne", "foo")

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "Invalid the status code "+strconv.Itoa(http.StatusBadRequest))
}
