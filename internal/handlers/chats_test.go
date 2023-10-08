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

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: "",
		Tags:           []string{"test1", "test2"},
	}

	mockChatRepository.On("GetByTelegramChatId", "124").Return(nil)
	mockChatRepository.On("Create", chatExpected.TelegramChatId, "", chatExpected.Tags).Return(&chatExpected)

	content := map[string]interface{}{
		"chat_id": "124",
		"tags":    []string{"test1", "test2"},
	}

	utils.MockPostRequest(ctx, content, false)

	CreateChat(ctx)

	mockChatRepository.AssertCalled(t, "GetByTelegramChatId", "124")
	mockChatRepository.AssertCalled(t, "Create", chatExpected.TelegramChatId, "", chatExpected.Tags)

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
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: "",
		Tags:           []string{"test1", "test2"},
	}

	utils.MockPostRequest(ctx, content, false)
	mockChatRepository.On("GetByTelegramChatId", "124").Return(&chat)

	CreateChat(ctx)

	mockChatRepository.AssertCalled(t, "GetByTelegramChatId", "124")
	mockChatRepository.AssertNotCalled(t, "Create")

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "StatusCode should be 400 Bad Request")
	assert.Equal(t, "Chat id already used", res.Error)
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
		TelegramUserId: "",
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

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: "",
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)

	mocksFeedRepository.On("GetByChatId", chatExpected.Id).Return([]string{feedExpected.Url})

	GetChatFeeds(ctx)

	mocksFeedRepository.AssertCalled(t, "GetByChatId", chatExpected.Id)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.NotEmpty(t, res.Feeds)
	assert.Len(t, res.Feeds, 1)
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
		Id:  uuid.New().String(),
		Url: "https://google.com",
	}

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: "",
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
		Id:  uuid.New().String(),
		Url: "https://google.com",
	}

	chatExpected := repository.DBChat{
		Id:             uuid.New().String(),
		TelegramChatId: "124",
		TelegramUserId: "",
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
		TelegramUserId: "",
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
		TelegramUserId: "",
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
		TelegramUserId: "",
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
		TelegramUserId: "",
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
		TelegramUserId: "",
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
		TelegramUserId: "",
		Tags:           []string{"test1", "test2"},
	}

	ctx.Set("chat", &chatExpected)

	mockChatRepository.On("UpdateParsed", chatExpected.Id).Return(true)

	ParsedChat(ctx)

	mockChatRepository.AssertCalled(t, "UpdateParsed", chatExpected.Id)

	assert.Equal(t, http.StatusNoContent, recorder.Result().StatusCode, "StatusCode should be 204 no content")
}
