package bot

import (
	"bytes"
	"errors"
	"log"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBotCmd struct {
	mock.Mock
}

func (c *mockBotCmd) GetCmdString() string {
	return "init"
}

func (c *mockBotCmd) Execute(bot BotAPI, message tgbotapi.Message)            {}
func (c *mockBotCmd) Reply(bot BotAPI, message tgbotapi.Message, data string) {}

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

func TestPostInConvBotSendError(t *testing.T) {
	// Create a new buffer to capture the log output
	var buf bytes.Buffer

	// Set the output of the logger to the buffer
	log.SetOutput(&buf)

	bot := new(MockBot)

	msg := tgbotapi.NewMessage(123, "Hello, world!")
	message := tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 123,
		},
		MessageID: 456,
	}

	bot.On("Send", msg).Return(tgbotapi.Message{}, errors.New("STOOOOOP"))

	// Call the function being tested
	postInConv(bot, message, msg.Text, false)

	// Assert that bot.Send was called with the correct argument
	bot.AssertCalled(t, "Send", msg)

	assert.Contains(t, buf.String(), "STOOOOOP")
	assert.Contains(t, buf.String(), "Couldn't speak in the conversation")
}

func TestIfNotValidCommandUse(t *testing.T) {
	ListCmd = []BotCmd{}

	res := isValidCommand("SUPER_COMMAND_USENR_JDFAJFAJ")

	assert.Nil(t, res, "This command should not exist")
}

func TestValidCommand(t *testing.T) {
	ListCmd = []BotCmd{
		new(mockBotCmd),
	}

	res := isValidCommand("init")

	assert.NotNil(t, res, "The init command should be present")
}
