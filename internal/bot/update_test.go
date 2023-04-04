package bot

import (
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestShouldActNotActBecauseInGroup(t *testing.T) {
	res := shouldAct(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				Type: "group",
			},
		},
	})

	assert.False(t, res, "Should be false as the type is group and not private ")
}

func TestShouldActNotActBecauseNotCorrectChatId(t *testing.T) {
	os.Setenv("TELEGRAM_USER_CHAT_ID", "132")

	res := shouldAct(tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID:   1233,
				Type: "private",
			},
		},
	})

	assert.False(t, res, "Should be false as the chat id doesn't match: 132")
}

func TestShouldActNotActBecauseNotMessageOrCallBack(t *testing.T) {
	os.Setenv("TELEGRAM_USER_CHAT_ID", "1233")

	res := shouldAct(tgbotapi.Update{
		CallbackQuery: nil,
		Message:       nil,
		EditedMessage: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID:   1233,
				Type: "private",
			},
		},
	})

	assert.False(t, res, "Should be false as Message and CallbackQuery are nil")
}
