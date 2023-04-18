package bot

import (
	"github.com/Vico1993/Otto/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type startCmd struct{}

// Start command string
func (cmd *startCmd) GetCmdString() string {
	return "start"
}

// Start execution process
// TODO: Implement different dialog + support different lang
func (cmd *startCmd) Execute(bot BotAPI, message tgbotapi.Message) {
	intro := "Hello! My name is Otto and I am a bot created to help you find the latest RSS feeds. \n With my advanced search algorithms, I can help you discover a wide range of RSS feeds from various sources. \nJust let me know what topics interest you, and I'll take care of the rest!"

	// If user didn't provide firstName
	user := message.From.FirstName
	if user == "" {
		user = message.From.UserName
	}

	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		intro+"\n\n"+user+" I don't know you yet, nice to meet you.\n Do you mind if I store some of your information. \n What I use you will ask? \n Just your id, firstname, lastname, your username,  if you are bot ",
	)

	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
			tgbotapi.NewInlineKeyboardButtonData("No", "no"),
		),
	)

	sendMessage(bot, msg)
}

func (cmd *startCmd) Reply(bot BotAPI, message tgbotapi.Message, data string) {
	if data == "yes" {
		// Initialisation
		repository.User.Create(
			message.Chat.ID,
			message.From.ID,
			message.From.FirstName,
			message.From.LastName,
			message.From.UserName,
			message.From.LanguageCode,
			message.From.IsBot,
			false,
		)

		// TODO: Send the list of commands
		postInConv(bot, message, "Thank you! let's work together", false)
	} else if data == "no" {
		postInConv(bot, message, "Sad but I understand your point, I will not save your information but I will not be able to help you.", false)
	}
}
