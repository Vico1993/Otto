package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/Vico1993/Otto/internal/utils"
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
	assert.Len(t, res.Articles, 1)
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

func TestDeleteArticle(t *testing.T) {
	type Response struct {
		Deleted bool `json:"deleted"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockArticleRepository := new(repository.MocksArticleRepository)
	repository.Article = mockArticleRepository

	articleId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     articleId,
		FeedId: uuid.New().String(),
		Title:  "Super Article",
		Source: "The Truth",
		Author: "Unknown",
		Link:   "https://thetruth.com",
		Tags:   []string{"not", "true"},
	}

	ctx.Set("article", &articleExpected)

	mockArticleRepository.On("Delete", articleId).Return(true)

	DeleteArticle(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.True(t, res.Deleted, "Deleted should be true")
}

func TestCreateArticle(t *testing.T) {
	type Response struct {
		Article *repository.DBArticle `json:"article"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockArticleRepository := new(repository.MocksArticleRepository)
	repository.Article = mockArticleRepository

	feedId := uuid.New().String()
	title := "Super Article"
	source := "The Truth"
	author := "Unknown"
	link := "https://thetruth.com"
	tags := []string{"not", "true"}

	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  title,
		Source: source,
		Author: author,
		Link:   link,
		Tags:   tags,
	}

	mockArticleRepository.On("Create", feedId, title, source, author, link, tags).Return(&articleExpected)

	content := map[string]interface{}{
		"feed_id": feedId,
		"title":   title,
		"source":  source,
		"author":  author,
		"link":    link,
		"tags":    tags,
	}

	utils.MockPostRequest(ctx, content, false)

	CreateArticle(ctx)

	mockArticleRepository.AssertCalled(t, "Create", feedId, title, source, author, link, tags)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Article)
}

func TestCreateArticleMissingRequiereParam(t *testing.T) {
	type Response struct {
		Error string `json:"error"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockArticleRepository := new(repository.MocksArticleRepository)
	repository.Article = mockArticleRepository

	feedId := uuid.New().String()
	source := "The Truth"
	author := "Unknown"
	link := "https://thetruth.com"
	tags := []string{"not", "true"}
	content := map[string]interface{}{
		"feed_id": feedId,
		"source":  source,
		"author":  author,
		"link":    link,
		"tags":    tags,
	}

	utils.MockPostRequest(ctx, content, false)

	CreateArticle(ctx)

	mockArticleRepository.AssertNotCalled(t, "Create")

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "StatusCode should be 400 Bad Request")
	assert.Equal(t, "Key: 'articleCreatePost.Title' Error:Field validation for 'Title' failed on the 'required' tag", res.Error)
}
