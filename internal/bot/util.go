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

	sendMessage(bot, msg)
}

// Just send a Message
func sendMessage(bot BotAPI, message tgbotapi.MessageConfig) {
	_, err := bot.Send(message)
	if err != nil {
		fmt.Println()
		fmt.Println(err.Error())

		logger := log.Default()
		logger.Println("---------------------------")
		logger.Println("Error")
		logger.Println("Couldn't speak in the conversation")
		logger.Println("Tried to say ", message.Text)
		logger.Println(err.Error())
		logger.Println("---------------------------")
	}
}

// Check if the command receive is valid
func isValidCommand(cmdString string) *BotCmd {
	for _, command := range ListCmd {
		if command.GetCmdString() == cmdString {
			return &command
		}
	}

	return nil
}
