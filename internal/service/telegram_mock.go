package service

import "github.com/stretchr/testify/mock"

type MocksTelegramService struct {
	mock.Mock
}

func (m *MocksTelegramService) TelegramPostMessage(chatId string, text string) {
	m.Called(chatId, text)
}

func (m *MocksTelegramService) TelegramUpdateTyping(chatId string, val bool) {
	m.Called(chatId, val)
}

func (m *MocksTelegramService) TelegramCreateTopic(chatId string, name string) {
	m.Called(chatId, name)
}

func (m *MocksTelegramService) GetBaseUrl() string {
	args := m.Called()
	return args.String(0)
}
