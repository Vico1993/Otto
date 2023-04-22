package bot

import (
	"log"
	"os"

	"github.com/Vico1993/Otto/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotCmd interface {
	GetCmdString() string
	Execute(bot BotAPI, message tgbotapi.Message)
	Reply(bot BotAPI, message tgbotapi.Message, data string)
}

var ListCmd []BotCmd

func Init() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	service.NewTelegramService().TelegramPostMessage("Just received an updates!")
	initCommand()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handleUpdates(update, bot)
	}
}

// Initiate list of command
func initCommand() {
	ListCmd = append(ListCmd, &startCmd{})
}
