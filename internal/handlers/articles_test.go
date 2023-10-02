package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllArticles(t *testing.T) {
	type Response struct {
		Articles []*repository.DBArticle `json:"articles"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockArticleRepository := new(repository.MocksArticleRepository)
	repository.Article = mockArticleRepository

	feedId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  "Super Article",
		Source: "The Truth",
		Author: "Unknown",
		Link:   "https://thetruth.com",
		Tags:   []string{"not", "true"},
	}

	mockArticleRepository.On("GetAll").Return([]*repository.DBArticle{&articleExpected})

	GetAllArticles(ctx)

	mockArticleRepository.AssertCalled(t, "GetAll")

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Articles)
}

func TestGetOneArticle(t *testing.T) {
	type Response struct {
		Article *repository.DBArticle `json:"article"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	feedId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  "Super Article",
		Source: "The Truth",
		Author: "Unknown",
		Link:   "https://thetruth.com",
		Tags:   []string{"not", "true"},
	}

	ctx.Set("article", &articleExpected)

	GetArticle(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Article)
}
