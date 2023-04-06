package bot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

type MockBot struct {
	mock.Mock
}

func (m *MockBot) Send(msg tgbotapi.Chattable) (tgbotapi.Message, error) {
	args := m.Called(msg)
	return tgbotapi.Message{}, args.Error(1)
}

func TestPostInConvWithoutReply(t *testing.T) {
	bot := new(MockBot)

	msg := tgbotapi.NewMessage(123, "Hello, world!")
	message := tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 123,
		},
		MessageID: 456,
	}

	bot.On("Send", msg).Return(tgbotapi.Message{}, nil)

	// Call the function being tested
	postInConv(bot, message, msg.Text, false)

	// Assert that bot.Send was called with the correct argument
	bot.AssertCalled(t, "Send", msg)
}

func TestPostInConvWithReply(t *testing.T) {
	bot := new(MockBot)

	msg := tgbotapi.NewMessage(123, "Hello, world!")
	msg.ReplyToMessageID = 456
	message := tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 123,
		},
		MessageID: 456,
	}

	bot.On("Send", msg).Return(tgbotapi.Message{}, nil)

	// Call the function being tested
	postInConv(bot, message, msg.Text, true)

	// Assert that bot.Send was called with the correct argument
	bot.AssertCalled(t, "Send", msg)
}
