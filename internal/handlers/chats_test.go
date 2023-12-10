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

func TestCreateChat(t *testing.T) {
	type Response struct {
		Chat *repository.DBChat `json:"chat"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	telegramThreadId := "245"
	chatExpected := repository.DBChat{
		Id:               uuid.New().String(),
		TelegramChatId:   "124",
		TelegramUserId:   nil,
		TelegramThreadId: &telegramThreadId,
		Tags:             []string{"test1", "test2"},
	}

	mockChatRepository.On("GetByTelegramChatIdAndThreadId", "124", "245").Return(nil)
	mockChatRepository.On("Create", chatExpected.TelegramChatId, "", *chatExpected.TelegramThreadId, chatExpected.Tags).Return(&chatExpected)

	content := map[string]interface{}{
		"chat_id":   "124",
		"thread_id": "245",
		"tags":      []string{"test1", "test2"},
	}

	utils.MockPostRequest(ctx, content, false)

	CreateChat(ctx)

	mockChatRepository.AssertCalled(t, "GetByTelegramChatIdAndThreadId", "124", "245")
	mockChatRepository.AssertNotCalled(t, "GetByTelegramChatId")
	mockChatRepository.AssertCalled(t, "Create", chatExpected.TelegramChatId, "", *chatExpected.TelegramThreadId, chatExpected.Tags)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Chat)
	assert.Empty(t, res.Chat.TelegramUserId, "Parameter is optional, and should be empty")
}

func TestCreateChatMissingRequiredField(t *testing.T) {
	type Response struct {
		Error string `json:"error"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	content := map[string]interface{}{
		"tags": []string{"test1", "test2"},
	}

	utils.MockPostRequest(ctx, content, false)

	CreateChat(ctx)

	mockChatRepository.AssertNotCalled(t, "Create")

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "StatusCode should be 400 Bad Request")
	assert.Equal(t, "Key: 'chatCreatePost.ChatId' Error:Field validation for 'ChatId' failed on the 'required' tag", res.Error)
}

func TestCreateChatTelegramChatIdAlreadyUsed(t *testing.T) {
	type Response struct {
		Error string `json:"error"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	content := map[string]interface{}{
		"chat_id": "124",
		"tags":    []string{"test1", "test2"},
	}

	chat := repository.DBChat{
		Id:               uuid.New().String(),
		TelegramChatId:   "124",
		TelegramUserId:   nil,
		TelegramThreadId: nil,
		Tags:             []string{"test1", "test2"},
	}

	utils.MockPostRequest(ctx, content, false)
	mockChatRepository.On("GetByTelegramChatId", "124").Return(&chat)

	CreateChat(ctx)

	mockChatRepository.AssertCalled(t, "GetByTelegramChatId", "124")
	mockChatRepository.AssertNotCalled(t, "GetByTelegramChatIdAndThreadId")
	mockChatRepository.AssertNotCalled(t, "Create")

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "StatusCode should be 400 Bad Request")
	assert.Equal(t, "Chat id already used", res.Error)
}

func TestCreateChatTelegramChatIdAlreadyUsedBuThreadIdNotUsed(t *testing.T) {
	type Response struct {
		Chat *repository.DBChat `json:"chat"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	telegramThreadId := "245"
	newTelegramThreadId := "240"
	chatExpected := repository.DBChat{
		Id:               uuid.New().String(),
		TelegramChatId:   "124",
		TelegramUserId:   nil,
		TelegramThreadId: &telegramThreadId,
		Tags:             []string{"test1", "test2"},
	}

	mockChatRepository.On("GetByTelegramChatIdAndThreadId", "124", "240").Return(nil)
	mockChatRepository.On("Create", chatExpected.TelegramChatId, "", newTelegramThreadId, chatExpected.Tags).Return(&chatExpected)

	content := map[string]interface{}{
		"chat_id":   "124",
		"thread_id": newTelegramThreadId,
		"tags":      []string{"test1", "test2"},
	}

	utils.MockPostRequest(ctx, content, false)

	CreateChat(ctx)

	mockChatRepository.AssertCalled(t, "GetByTelegramChatIdAndThreadId", "124", "240")
	mockChatRepository.AssertNotCalled(t, "GetByTelegramChatId")
	mockChatRepository.AssertCalled(t, "Create", chatExpected.TelegramChatId, "", newTelegramThreadId, chatExpected.Tags)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Chat)
	assert.Empty(t, res.Chat.TelegramUserId, "Parameter is optional, and should be empty")
}

func TestDeleteChat(t *testing.T) {
	type Response struct {
		Deleted bool `json:"deleted"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)

	mockChatRepository.On("Delete", chatExpected.Id).Return(true)

	DeleteChat(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.True(t, res.Deleted, "Deleted should be true")
}

func TestGetChatFeeds(t *testing.T) {
	type Response struct {
		Feeds []chatFeedsResponse `json:"feeds"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mocksFeedRepository := new(repository.MocksFeedRepository)
	repository.Feed = mocksFeedRepository

	feedExpected := repository.DBFeed{
		Id:       uuid.New().String(),
		Url:      "https://google.com",
		Disabled: false,
	}

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)

	mocksFeedRepository.On("GetByChatId", chatExpected.Id).Return([]*repository.DBFeed{&feedExpected})

	GetChatFeeds(ctx)

	mocksFeedRepository.AssertCalled(t, "GetByChatId", chatExpected.Id)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Feeds)
	assert.Len(t, res.Feeds, 1)
	assert.Equal(t, []chatFeedsResponse{{Id: feedExpected.Id, Url: feedExpected.Url}}, res.Feeds, "The result should contain the feed url and the feed id")
}

func TestCreateChatFeeds(t *testing.T) {
	type Response struct {
		Added bool `json:"added"`
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mocksFeedRepository := new(repository.MocksFeedRepository)
	repository.Feed = mocksFeedRepository

	feedExpected := repository.DBFeed{
		Id:       uuid.New().String(),
		Url:      "https://google.com",
		Disabled: false,
	}

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)
	ctx.Set("feed", &feedExpected)

	mocksFeedRepository.On("LinkChatAndFeed", feedExpected.Id, chatExpected.Id).Return(true)

	CreateChatFeed(ctx)

	mocksFeedRepository.AssertCalled(t, "LinkChatAndFeed", feedExpected.Id, chatExpected.Id)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.True(t, res.Added, "Added should be true")
}

func TestDeleteChatFeeds(t *testing.T) {
	type Response struct {
		Deleted bool `json:"deleted"`
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mocksFeedRepository := new(repository.MocksFeedRepository)
	repository.Feed = mocksFeedRepository

	feedExpected := repository.DBFeed{
		Id:       uuid.New().String(),
		Url:      "https://google.com",
		Disabled: false,
	}

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)
	ctx.Set("feed", &feedExpected)

	mocksFeedRepository.On("UnLinkChatAndFeed", feedExpected.Id, chatExpected.Id).Return(true)

	DeleteChatFeed(ctx)

	mocksFeedRepository.AssertCalled(t, "UnLinkChatAndFeed", feedExpected.Id, chatExpected.Id)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.True(t, res.Deleted, "Deleted should be true")
}

func TestGetChatTags(t *testing.T) {
	type Response struct {
		Tags []string `json:"tags"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)

	GetChatTags(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Tags)
	assert.Len(t, res.Tags, 2)
}

func TestCreateChatTag(t *testing.T) {
	type Response struct {
		Tags []string `json:"tags"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{},
	}

	ctx.Set("chat", &chatExpected)

	content := map[string]interface{}{
		"tags": []string{"test1", "test2"},
	}

	utils.MockPostRequest(ctx, content, false)

	mockChatRepository.On("UpdateTags", chatExpected.Id, []string{"test1", "test2"}).Return(true)

	CreateChatTag(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Tags)
	assert.Len(t, res.Tags, 2)
}

func TestCreateChatTagMissingTags(t *testing.T) {
	type Response struct {
		Error string `json:"error"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{},
	}

	ctx.Set("chat", &chatExpected)

	content := map[string]interface{}{}

	utils.MockPostRequest(ctx, content, false)

	CreateChatTag(ctx)

	mockChatRepository.AssertNotCalled(t, "UpdateTags")

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "StatusCode should be 400 Bad Request")
}

func TestDeleteChatTag(t *testing.T) {
	type Response struct {
		Tags []string `json:"tags"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)
	ctx.AddParam("tag", "test1")

	mockChatRepository.On("UpdateTags", chatExpected.Id, []string{"test2"}).Return(true)

	DeleteChatTag(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Tags)
	assert.Len(t, res.Tags, 1)
}

func TestDeleteChatTagTagNotFound(t *testing.T) {
	type Response struct {
		Error string `json:"error"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)
	ctx.AddParam("tag", "test3")

	DeleteChatTag(ctx)

	mockChatRepository.AssertNotCalled(t, "UpdateTags")

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "StatusCode should be 400 Bad Request")
	assert.Equal(t, "Tag not found", res.Error)
}

func TestParsedChat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)

	mockChatRepository.On("UpdateParsed", chatExpected.Id).Return(true)

	ParsedChat(ctx)

	mockChatRepository.AssertCalled(t, "UpdateParsed", chatExpected.Id)

	assert.Equal(t, http.StatusNoContent, recorder.Result().StatusCode, "StatusCode should be 204 no content")
}

func TestGetAll(t *testing.T) {
	type Response struct {
		Chats []*repository.DBChat `json:"chats"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	mockChatRepository := new(repository.MocksChatRepository)
	repository.Chat = mockChatRepository

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{"test1", "test2"},
	}

	mockChatRepository.On("GetAll").Return([]*repository.DBChat{&chatExpected})

	GetAllChats(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	mockChatRepository.AssertCalled(t, "GetAll")
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Chats)
	assert.Len(t, res.Chats, 1, "Chats array should contain 1 element")
}

func TestGetLatestArticleFromChat(t *testing.T) {
	type Response struct {
		Articles []*repository.DBArticle `json:"articles"`
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: uuid.New().String(),
		Title:  "Super Article",
		Source: "The Truth",
		Author: "Unknown",
		Link:   "https://thetruth.com",
		Tags:   []string{"not", "true"},
	}

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: nil,
		Tags:           []string{},
	}

	mockArticleRepository := new(repository.MocksArticleRepository)
	repository.Article = mockArticleRepository

	ctx.Set("chat", &chatExpected)

	mockArticleRepository.On("GetByChatAndTime", chatExpected.Id).Return([]*repository.DBArticle{&articleExpected})

	GetLatestArticleFromChat(ctx)

	mockArticleRepository.AssertCalled(t, "GetByChatAndTime", chatExpected.Id)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Articles)
	assert.Len(t, res.Articles, 1)
}
