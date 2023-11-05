package repository

import "github.com/stretchr/testify/mock"

type MocksChatRepository struct {
	mock.Mock
}

func (m *MocksChatRepository) Create(telegramChatId string, telegramUserId string, telegramThreadId string, tags []string) *DBChat {
	args := m.Called(telegramChatId, telegramUserId, telegramThreadId, tags)
	return args.Get(0).(*DBChat)
}

func (m *MocksChatRepository) Delete(uuid string) bool {
	args := m.Called(uuid)
	return args.Get(0).(bool)
}

func (m *MocksChatRepository) GetOne(uuid string) *DBChat {
	args := m.Called(uuid)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBChat)
}

func (m *MocksChatRepository) GetByTelegramChatId(chatId string) *DBChat {
	args := m.Called(chatId)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBChat)
}

func (m *MocksChatRepository) GetByTelegramChatIdAndThreadId(chatId string, threadId string) *DBChat {
	args := m.Called(chatId, threadId)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBChat)
}

func (m *MocksChatRepository) GetAll() []*DBChat {
	args := m.Called()
	return args.Get(0).([]*DBChat)
}

func (m *MocksChatRepository) GetByChatId(uuid string) []string {
	args := m.Called(uuid)
	return args.Get(0).([]string)
}

func (m *MocksChatRepository) UpdateTags(uuid string, tags []string) bool {
	args := m.Called(uuid, tags)
	return args.Get(0).(bool)
}

func (m *MocksChatRepository) UpdateParsed(uuid string) bool {
	args := m.Called(uuid)
	return args.Get(0).(bool)
}
