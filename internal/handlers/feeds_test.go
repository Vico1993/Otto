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

func TestGetAllFeeds(t *testing.T) {
	type Response struct {
		Feeds []*repository.DBFeed `json:"feeds"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mocksFeedRepository := new(repository.MocksFeedRepository)
	repository.Feed = mocksFeedRepository

	feedExpected := repository.DBFeed{
		Id:  uuid.New().String(),
		Url: "https://google.com",
	}

	mocksFeedRepository.On("GetAll").Return([]*repository.DBFeed{&feedExpected})

	GetAllFeeds(ctx)

	mocksFeedRepository.AssertCalled(t, "GetAll")

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Feeds)
	assert.Len(t, res.Feeds, 1)
}

func TestGetOneFeed(t *testing.T) {
	type Response struct {
		Feed *repository.DBFeed `json:"feed"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	feedId := uuid.New().String()
	feedExpected := repository.DBFeed{
		Id:  feedId,
		Url: "https://google.com",
	}

	ctx.Set("feed", &feedExpected)

	GetFeed(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Feed)
}

func TestDeleteFeed(t *testing.T) {
	type Response struct {
		Deleted bool `json:"deleted"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mocksFeedRepository := new(repository.MocksFeedRepository)
	repository.Feed = mocksFeedRepository

	feedExpected := repository.DBFeed{
		Id:  uuid.New().String(),
		Url: "https://google.com",
	}

	ctx.Set("feed", &feedExpected)

	mocksFeedRepository.On("Delete", feedExpected.Id).Return(true)

	DeleteFeed(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.True(t, res.Deleted, "Deleted should be true")
}

func TestCreateFeed(t *testing.T) {
	type Response struct {
		Feed *repository.DBFeed `json:"feed"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mocksFeedRepository := new(repository.MocksFeedRepository)
	repository.Feed = mocksFeedRepository

	feedExpected := repository.DBFeed{
		Id:  uuid.New().String(),
		Url: "https://google.com",
	}

	mocksFeedRepository.On("Create", feedExpected.Url).Return(&feedExpected)

	content := map[string]interface{}{
		"url": feedExpected.Url,
	}

	utils.MockPostRequest(ctx, content, false)

	CreateFeed(ctx)

	mocksFeedRepository.AssertCalled(t, "Create", feedExpected.Url)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Feed)
}

func TestCreateFeedMissingRequiereParam(t *testing.T) {
	type Response struct {
		Error string `json:"error"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mocksFeedRepository := new(repository.MocksFeedRepository)
	repository.Feed = mocksFeedRepository

	content := map[string]interface{}{}

	utils.MockPostRequest(ctx, content, false)

	CreateFeed(ctx)

	mocksFeedRepository.AssertNotCalled(t, "Create")

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "StatusCode should be 400 Bad Request")
	assert.Equal(t, "Key: 'feedCreatePost.Url' Error:Field validation for 'Url' failed on the 'required' tag", res.Error)
}

func TestGetFeedsArticle(t *testing.T) {
	type Response struct {
		Articles []*repository.DBArticle `json:"articles"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	feedId := uuid.New().String()
	feedExpected := repository.DBFeed{
		Id:  feedId,
		Url: "https://google.com",
	}

	ctx.Set("feed", &feedExpected)

	mockArticleRepository := new(repository.MocksArticleRepository)
	repository.Article = mockArticleRepository

	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  "Super Article",
		Source: "The Truth",
		Author: "Unknown",
		Link:   "https://thetruth.com",
		Tags:   []string{"not", "true"},
	}

	mockArticleRepository.On("GetByFeedId", feedId).Return([]*repository.DBArticle{&articleExpected})

	GetFeedArticles(ctx)

	mockArticleRepository.AssertCalled(t, "GetByFeedId", feedId)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Articles)
	assert.Len(t, res.Articles, 1)

}
