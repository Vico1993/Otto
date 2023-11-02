package repository

import (
	"github.com/stretchr/testify/mock"
)

// MOCKS
type MocksFeedRepository struct {
	mock.Mock
}

func (m *MocksFeedRepository) Create(url string) *DBFeed {
	args := m.Called(url)
	return args.Get(0).(*DBFeed)
}

func (m *MocksFeedRepository) Delete(uuid string) bool {
	args := m.Called(uuid)
	return args.Get(0).(bool)
}

func (m *MocksFeedRepository) GetOne(uuid string) *DBFeed {
	args := m.Called(uuid)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBFeed)
}

func (m *MocksFeedRepository) GetAll() []*DBFeed {
	args := m.Called()
	return args.Get(0).([]*DBFeed)
}

func (m *MocksFeedRepository) GetAllActive() []*DBFeed {
	args := m.Called()
	return args.Get(0).([]*DBFeed)
}

func (m *MocksFeedRepository) GetByChatId(uuid string) []*DBFeed {
	args := m.Called(uuid)
	return args.Get(0).([]*DBFeed)
}

func (m *MocksFeedRepository) LinkChatAndFeed(feedId string, chatId string) bool {
	args := m.Called(feedId, chatId)
	return args.Get(0).(bool)
}

func (m *MocksFeedRepository) UnLinkChatAndFeed(feedId string, chatId string) bool {
	args := m.Called(feedId, chatId)
	return args.Get(0).(bool)
}

func (m *MocksFeedRepository) DisableFeed(feedId string) bool {
	args := m.Called(feedId)
	return args.Get(0).(bool)
}
