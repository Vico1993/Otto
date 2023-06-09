package repository

import (
	"github.com/Vico1993/Otto/internal/database"
	"github.com/stretchr/testify/mock"
)

type MocksChatRep struct {
	mock.Mock
}

func (m *MocksChatRep) GetAll() []*database.Chat {
	m.Called()
	return []*database.Chat{}
}

func (m *MocksChatRep) FindByChatId(chatId string) *database.Chat {
	args := m.Called(chatId)
	return args.Get(0).(*database.Chat)
}

func (m *MocksChatRep) UpdateFeedCheckForUrl(url string, articleFound int, chat *database.Chat) bool {
	args := m.Called(url, articleFound, chat)
	return args.Bool(0)
}

func (m *MocksChatRep) PushNewFeed(url string, chat *database.Chat) bool {
	args := m.Called(url, chat)
	return args.Bool(0)
}

func (m *MocksChatRep) Create(chatid string, userid int64, tags []string, feeds []string) *database.Chat {
	args := m.Called(chatid, userid, tags, feeds)
	return args.Get(0).(*database.Chat)
}
