package repository

import (
	"github.com/Vico1993/Otto/internal/database"
	"github.com/stretchr/testify/mock"
)

type MocksChatRep struct {
	*mock.Mock
}

func (m *MocksChatRep) GetAll() []*database.Chat {
	m.Called()
	return []*database.Chat{}
}

func (m *MocksChatRep) FindByChatId(chatId string) *database.Chat {
	m.Called(chatId)
	return database.NewChat("123", 123, []database.Feed{})
}

func (m *MocksChatRep) UpdateFeedCheckForUrl(url string, articleFound int, chat *database.Chat) bool {
	m.Called(url, articleFound, chat)
	return true
}

func (m *MocksChatRep) PushNewFeed(url string, chat *database.Chat) bool {
	m.Called(url, chat)
	return true
}

func (m *MocksChatRep) Create(chatid string, userid int64, tags []string, feeds []string) *database.Chat {
	m.Called(chatid, userid, tags, feeds)
	return database.NewChat("123", 123, []database.Feed{})
}
