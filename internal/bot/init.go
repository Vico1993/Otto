package bot

import (
	"log"
	"os"

	"github.com/Vico1993/Otto/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	service.NewTelegramService().TelegramPostMessage(
		os.Getenv("TELEGRAM_USER_CHAT_ID"),
		`*Upgrade complete*! Ready to be even smarter and funnier than before. 🤖 🚀 ✨`,
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handleUpdates(update, bot)
	}
}
