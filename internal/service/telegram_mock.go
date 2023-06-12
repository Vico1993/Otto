package service

import "github.com/stretchr/testify/mock"

type MocksTelegramService struct {
	mock.Mock
}

func (m *MocksTelegramService) TelegramPostMessage(text string) {
	m.Called(text)
}

func (m *MocksTelegramService) TelegramUpdateTyping(val bool) {
	m.Called(val)
}

func (m *MocksTelegramService) GetChatId() string {
	args := m.Called()
	return args.String(0)
}

func (m *MocksTelegramService) GetBaseUrl() string {
	args := m.Called()
	return args.String(0)
}
