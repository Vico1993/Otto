package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI interface {
	Send(msg tgbotapi.Chattable) (tgbotapi.Message, error)
}

// Push a message in the conversation
func postInConv(bot BotAPI, message tgbotapi.Message, text string, reply bool) {
	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		text,
	)

	if reply {
		msg.ReplyToMessageID = message.MessageID
	}

	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println()
		fmt.Println(err.Error())

		logger := log.Default()
		logger.Println("---------------------------")
		logger.Println("Error")
		logger.Println("Couldn't speak in the conversation")
		logger.Println("Tried to say ", text)
		logger.Println(err.Error())
		logger.Println("---------------------------")
	}
}
